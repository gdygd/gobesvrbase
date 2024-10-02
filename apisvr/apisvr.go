package main

import (
	"besvrbase/sharedObj"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/gdygd/goglib"

	"apisvr/app/am"
	"apisvr/app/dbapp"
	"apisvr/app/dbapp/mdb"
	"apisvr/app/httpapp"
	"apisvr/app/kafkaapp"
	"apisvr/app/msgapp"
	"apisvr/app/netapp"
	"apisvr/app/objdb"

	"github.com/gorilla/handlers"

	"github.com/gdygd/goshm/shmlinux"
	"gopkg.in/ini.v1"
)

// ------------------------------------------------------------------------------
// Local
// ------------------------------------------------------------------------------
var shminst *shmlinux.Linuxshm = nil

var dbInst dbapp.DBHandler
var tcpInst *netapp.CommHandler = nil
var kfkproducInst *kafkaapp.KafkaHandler = nil
var kfkconsumInst *kafkaapp.KafkaHandler = nil

var httpInst *httpapp.HttpAppHandler = nil
var process *goglib.Process
var terminate bool = false

var logmngTimer time.Time
var tmupdTimer time.Time

var thrmsgprc *goglib.Thread = nil       // msg process thread
var thrnetclient *goglib.Thread = nil    // net client thread
var thrkafkaConsum *goglib.Thread = nil  // kafka consumer thread
var thrkafkaProduce *goglib.Thread = nil // kafka producer thread

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const sysenvini = "./sys_env.ini"
const systemini = "./system.ini" // shared memory id
const versionpath = "./version.txt"

const LOGMNG_INTERVAL = 1000 * 600 // 10 minute
const SYSSTTCHECK_INTERVAL = 1000  // 1sec
const LOGFILE_STORE_DATE = 30      // 30일만 로그파일 보존
const TMUPD_INTERVAL = 1000

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	for {
		signal := <-chSig
		am.Applog.Always("[apisvr]]Accept Signal %d", signal)
		switch signal {
		case syscall.SIGHUP: // 터미널 연결 끊겼을경우
			am.Applog.Always("SIGHUP(%d)\n", signal)
			//terminate = true
		case syscall.SIGINT:
			am.Applog.Always("SIGINT(%d)\n", signal)
			//terminate = true
		case syscall.SIGTERM:
			am.Applog.Always("SIGTERM(%d)\n", signal)
			//clearEnv()
			terminate = true
		case syscall.SIGKILL:
			am.Applog.Always("SIGKILL(%d)\n", signal)
			terminate = true
		default:
			am.Applog.Always("Unknown signal(%d)\n", signal)
			terminate = true
		}
	}
}

// ------------------------------------------------------------------------------
// checkSystemState
// ------------------------------------------------------------------------------
func checkSystemState() {

	// server time update
	if goglib.CheckElapsedTime(&tmupdTimer, TMUPD_INTERVAL) {
		curtm := time.Now()
		utcsecs := curtm.Unix()
		objdb.SysInfo.SetSvrTime(utcsecs)
	}
}

// ------------------------------------------------------------------------------
// manageDebug
// ------------------------------------------------------------------------------
func manageDebug() {
	debugLv := process.DebugLv
	am.Applog.SetLevel(debugLv)

	//dbapp.SetlogLevel(debugLv)
}

// ------------------------------------------------------------------------------
// initVariable
// ------------------------------------------------------------------------------
func initVariable() {
	am.Applog.Always("[apisvr]initVariable..")

	logmngTimer = time.Now()
}

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	am.Applog.Always("[apisvr]initEnvVaiable..")

	cfg, err := ini.Load(sysenvini)
	if err != nil {
		am.Applog.Error("fail to read sysenvini.ini %v", err)
		return false
	}

	am.AppVar.NetAddr = cfg.Section("NETINFO").Key("addr").String()
	am.AppVar.NetPort, _ = cfg.Section("NETINFO").Key("port").Int()

	am.AppVar.DbHost = cfg.Section("DATABASE").Key("host").String()
	am.AppVar.DbPort, _ = cfg.Section("DATABASE").Key("port").Int()
	am.AppVar.DbUser = cfg.Section("DATABASE").Key("user").String()
	am.AppVar.DbName = cfg.Section("DATABASE").Key("dbname").String()
	am.AppVar.DbPasswd = cfg.Section("DATABASE").Key("passwd").String()

	am.AppVar.HttpPort, _ = cfg.Section("HTTP").Key("port").Int()
	am.AppVar.Domain = cfg.Section("HTTP").Key("domain").String()

	am.AppVar.Https = cfg.Section("HTTP").Key("https").String()
	am.AppVar.Sslcrt = cfg.Section("HTTP").Key("sslcrt").String()
	am.AppVar.Sslkey = cfg.Section("HTTP").Key("sslkey").String()
	am.AppVar.Sslcertpem = cfg.Section("HTTP").Key("certpem").String()
	strAlloworigin := cfg.Section("HTTP").Key("alloworigins").String()
	arralloworgin := strings.Split(strAlloworigin, ",")
	am.AppVar.Alloworigins = []string{}
	for _, allolworigin := range arralloworgin {
		if am.AppVar.Https != "yes" {
			am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("http://%s", allolworigin))
			am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("https://%s", allolworigin))
		} else {
			am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("http://%s", allolworigin))
			am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("https://%s", allolworigin))
		}
	}
	if am.AppVar.Https != "yes" {
		am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("http://%s:%d", am.AppVar.Domain, am.AppVar.HttpPort))
		am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("https://%s:%d", am.AppVar.Domain, am.AppVar.HttpPort))
	} else {
		am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("http://%s:%d", am.AppVar.Domain, am.AppVar.HttpPort))
		am.AppVar.Alloworigins = append(am.AppVar.Alloworigins, fmt.Sprintf("https://%s:%d", am.AppVar.Domain, am.AppVar.HttpPort))
	}

	am.AppVar.DebugLv, _ = cfg.Section("OPRMODE").Key("debuglv").Int()
	am.AppVar.KafkaMode = cfg.Section("OPRMODE").Key("kafka").String()

	am.Applog.Always("[apisvr]App Environment variable")

	am.Applog.Always(" \t [Database]")
	am.Applog.Always(" \t host:%s", am.AppVar.DbHost)
	am.Applog.Always(" \t port:%d", am.AppVar.DbPort)
	am.Applog.Always(" \t user:%s", am.AppVar.DbUser)
	am.Applog.Always(" \t pw:%s", am.AppVar.DbPasswd)
	am.Applog.Always(" \t db  :%s", am.AppVar.DbName)

	am.Applog.Always(" \t [Http]")
	am.Applog.Always(" \t port:%d", am.AppVar.HttpPort)
	am.Applog.Always(" \t domain:%s", am.AppVar.Domain)
	am.Applog.Always(" \t AllowOrigins:%v", am.AppVar.Alloworigins)

	am.Applog.Always(" \t [OPRMODE]")
	am.Applog.Always(" \t debuglv:%d", am.AppVar.DebugLv)
	am.Applog.Always(" \t kafka:%s", am.AppVar.KafkaMode)

	// read version file
	dat, err := os.ReadFile(versionpath)
	strversion := fmt.Sprintf("%v", string(dat))
	am.Applog.Always(" \t VERSION:%s", strversion)

	return true
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {
	am.Applog.Always("[apisvr]iniSignal..")
	// signal handler
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)
	go sigHandler(ch_signal)
}

// ------------------------------------------------------------------------------
// isAlreadyProcess
// ------------------------------------------------------------------------------
func isAlreadyProcess() bool {
	// 프로세스 러닝상태 확인
	// 이미 프로세스가 동작중이면 false

	am.Applog.Info("[apisvr]isAlreadyProcess...")

	var isrunning bool = false

	am.Applog.Print(3, "prcdesc : %v", process.PrcName)

	prcnm := process.PrcName // 코어프로세스
	cmdstr := fmt.Sprintf(`ps -ef | grep %s | grep -v grep`, prcnm)
	cmd := exec.Command("bash", "-c", cmdstr)
	output, _ := cmd.CombinedOutput()
	stroutput := string(output)

	am.Applog.Print(3, "cmd output : %v %d", stroutput, len(stroutput))

	if len(stroutput) == 0 {
		isrunning = false
	} else {
		isrunning = true
	}

	return isrunning
}

// ------------------------------------------------------------------------------
// initMemory
// ------------------------------------------------------------------------------
func initMemory() bool {
	shminst := shmlinux.NewLinuxShm()
	//shminst.InitShm(skey, MemorySize)
	shminst.InitShm(sharedObj.MEM_KEY, sharedObj.MEM_SIZE)

	err := shminst.CreateShm()
	if err != nil {
		am.Applog.Error("initMemory CreateShm err : ", err)
	}
	err = shminst.AttachShm()
	if err != nil {
		am.Applog.Error("initMemory AttachShm err : ", err)
	}
	objdb.SharedMem = (*sharedObj.SharedMemory)(unsafe.Pointer(shminst.Addr))
	return true
}

// ------------------------------------------------------------------------------
// initProcess
// ------------------------------------------------------------------------------
func initProcess() bool {
	am.Applog.Always("[apisvr]initProcess..")
	process = &objdb.SharedMem.Process[sharedObj.PRC_IDX_PRC01]
	return true
}

// ------------------------------------------------------------------------------
// RegisterProcess
// ------------------------------------------------------------------------------
func RegisterProcess() bool {

	am.Applog.Always("[apisvr]RegisterProcess..")
	process.RunBase.Active = true
	// register process id
	process.RegisterPid(os.Getpid())

	am.Applog.Always(" \t process ID : %d %d", process.GetPid(), os.Getpid())
	return true

}

// ------------------------------------------------------------------------------
// initDebugLevel
// ------------------------------------------------------------------------------
func initDebugLevel() {
	process.DebugLv = am.Applog.GetLevel()
}

// ------------------------------------------------------------------------------
// initDatabase
// ------------------------------------------------------------------------------
func initDatabase() bool {
	am.Applog.Always("[apisvr]initDatabase..")

	dbInst = mdb.NewMariadbHandler(am.AppVar.DbUser, am.AppVar.DbPasswd, am.AppVar.DbName, am.AppVar.DbHost)

	db, dbErr := dbInst.Open()
	defer func() {
		dbInst.Close(db)
	}()
	if dbErr != nil {
		am.Applog.Error("Database open fail(1)..%v", dbErr)
		return false
	}

	if dbInst == nil {
		am.Applog.Error("Database open fail(2)..")
		return false
	}

	return true
}

// ------------------------------------------------------------------------------
// InitObjectDB
// ------------------------------------------------------------------------------
func InitObjectDB() {
	am.Applog.Always("[apisvr]initObjectDB..")

}

// ------------------------------------------------------------------------------
// initObject
// ------------------------------------------------------------------------------
func initObject() {
	am.Applog.Always("[apisvr]initObject..")
}

// ------------------------------------------------------------------------------
// initHttp
// ------------------------------------------------------------------------------
func initHttp() bool {
	am.Applog.Always("[apisvr]initHttp..")

	//httpInst
	if dbInst != nil {
		httpInst = httpapp.MakeHandler(dbInst)
	} else {
		am.Applog.Info("initHttp fail.. db instance not exist")
		return false
	}

	go func() {
		var err error

		// cors policy option
		headersOK := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials", "Set-Cookie"})
		originsOK := handlers.AllowedOrigins(am.AppVar.Alloworigins)
		methodsOK := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "DELETE", "PUT"})

		// https options
		credential := handlers.AllowCredentials()

		port := fmt.Sprintf(":%d", am.AppVar.HttpPort)

		if am.AppVar.Https != "yes" {
			err = http.ListenAndServe(port, handlers.CORS(originsOK, headersOK, methodsOK)(httpInst.Handler))
		} else {
			err = http.ListenAndServeTLS(port, am.AppVar.Sslcrt, am.AppVar.Sslkey, handlers.CORS(originsOK, headersOK, methodsOK, credential)(httpInst.Handler))
		}

		if err != nil {
			panic(err)
			terminate = true
		}
	}()

	return true
}

// ------------------------------------------------------------------------------
// initServer
// ------------------------------------------------------------------------------
func initServer() bool {

	am.Applog.Always("[apisvr]initServer..")
	return true
}

// ------------------------------------------------------------------------------
// initNetwork
// ------------------------------------------------------------------------------
func initNetwork() bool {
	fmt.Println("initNetwork..")
	//  init network env
	// network instance
	tcpInst = netapp.MakeNetHandler("apiapp", am.AppVar.NetPort, am.AppVar.NetAddr)

	// connect network
	if tcpInst != nil {
		ok := tcpInst.Open()

		if !ok {
			am.Applog.Error("server connect fail.. %d %s", am.AppVar.NetPort, am.AppVar.NetAddr)
			//return false
			return true
		}

	} else {
		am.Applog.Error("Tcp instance nil..")
		//return false
		return true
	}

	kfkproducInst = kafkaapp.MakeKfkHandler("127.0.0.1", 9092)
	kfkconsumInst = kafkaapp.MakeKfkHandler("127.0.0.1", 9092)

	return true
}

// ------------------------------------------------------------------------------
// initThread
// ------------------------------------------------------------------------------
func initThread() {
	am.Applog.Always("[apisvr]initThread..")

	thrmsgprc = goglib.NewThread()
	thrmsgprc.Init(msgapp.THRPrcmsg, 10)

	thrnetclient = goglib.NewThread()
	thrnetclient.Init(netapp.THRNetclient, 10, tcpInst)

	thrkafkaConsum = goglib.NewThread()
	thrkafkaConsum.Init(kafkaapp.THRKafkaConsum, 10, kfkconsumInst)

	thrkafkaProduce = goglib.NewThread()
	thrkafkaProduce.Init(kafkaapp.THRKafkaProduce, 10, kfkproducInst)
}

// ------------------------------------------------------------------------------
// appRun
// ------------------------------------------------------------------------------
func appRun() {
	am.Applog.Always("[apisvr]appRun..")

	am.Applog.Always("[apisvr]message process thread start..(1)")
	thrmsgprc.Start()
	am.Applog.Always("[apisvr]message process thread start..(2)")

	// am.Applog.Always("[apisvr]net client process thread start..(1)")
	// thrnetclient.Start()
	// am.Applog.Always("[apisvr]net client process thread start..(2)")

	am.Applog.Always("[apisvr]kafka consumer thread start..(1)")
	thrkafkaConsum.Start()
	am.Applog.Always("[apisvr]kafka consumer thread start..(2)")

	// am.Applog.Always("[apisvr]kafka consumer thread start..(1)")
	// thrkafkaProduce.Start()
	// am.Applog.Always("[apisvr]kafka producer thread start..(2)")
}

// ------------------------------------------------------------------------------
// manageRoutine
// ------------------------------------------------------------------------------
func manageRoutine() {

	if thrmsgprc.RunBase.Active {
		var state int = 0

		thrmsgprc.IsRunning(&state)

		switch state {
		case goglib.RST_ABNOMAL:
			am.Applog.Warn("RST_ABNOMAL(thrmsgprc)")
			thrmsgprc.Kill()
			thrmsgprc.Start()
		default:
			break
		}
	}

	if thrnetclient.RunBase.Active {
		var state int = 0

		thrnetclient.IsRunning(&state)

		switch state {
		case goglib.RST_ABNOMAL:
			am.Applog.Warn("RST_ABNOMAL(thrnetclient)")
			thrnetclient.Kill()
			thrnetclient.Start()
		default:
			break
		}
	}

	if thrkafkaConsum.RunBase.Active {
		var state int = 0

		thrkafkaConsum.IsRunning(&state)

		switch state {
		case goglib.RST_ABNOMAL:
			am.Applog.Warn("RST_ABNOMAL(thrkafkaConsum)")
			thrkafkaConsum.Kill()
			thrkafkaConsum.Start()
		default:
			break
		}
	}
	if thrkafkaProduce.RunBase.Active {
		var state int = 0

		thrkafkaProduce.IsRunning(&state)

		switch state {
		case goglib.RST_ABNOMAL:
			am.Applog.Warn("RST_ABNOMAL(thrkafkaProduce)")
			thrkafkaProduce.Kill()
			thrkafkaProduce.Start()
		default:
			break
		}
	}
}

// ------------------------------------------------------------------------------
// manageLogfile
// ------------------------------------------------------------------------------
func manageLogfile() {

	if !goglib.CheckElapsedTime(&logmngTimer, LOGMNG_INTERVAL) {
		return
	}

	var path = "./log"

	finfos, err := goglib.ReadDir(path)

	if err != nil {
		am.Applog.Error("[apisvr]manageLogfile err.. %v", err)
		return
	}

	nowdt := time.Now()
	nowsec := nowdt.Unix()

	const RefDay = 86400 * LOGFILE_STORE_DATE

	for idx, f := range finfos {

		fsec := f.ModTime.Unix()

		if (nowsec - fsec) >= RefDay {
			fmt.Printf("(%d) %v \n", idx, f.FileName)
			err := os.RemoveAll(path + "/" + f.FileName)
			if err != nil {
				am.Applog.Error("Logfile remove err.. %v", err)
			}
		}
	}
}

func SetDebugLv() {
	am.Applog.Info("[apisvr]SetDebugLv...")
	process.SetDebugLv(am.AppVar.DebugLv)
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {

	// read initial file
	initEnvVaiable()

	initVariable()

	initSignal()

	if !initMemory() { // create memory
		am.Applog.Error("Share memory created fail..")
		return false
	}

	if !initProcess() {
		am.Applog.Error("Process initialize fail..")
		return false
	}

	if !RegisterProcess() {
		am.Applog.Error("Process RegisterProcess fail..")
		return false
	}

	initNetwork()
	initThread()

	if !initDatabase() {
		am.Applog.Error("Database initial fail..")
		//return false
	}

	initDebugLevel()

	initObject()
	initHttp()

	if !initServer() {
		am.Applog.Error("initServer initial fail..")
		return false
	}

	// set logging level
	SetDebugLv()

	return true
}

// ------------------------------------------------------------------------------
// clearEnv
// ------------------------------------------------------------------------------
func clearEnv() {
	am.Applog.Always("clearEnv quit [apisvr]...")

	// thread kill
	if thrmsgprc != nil && thrmsgprc.RunBase.Active {
		thrmsgprc.Kill()
	}
	if thrnetclient != nil && thrnetclient.RunBase.Active {
		thrnetclient.Kill() // thread종료시 network close
	}
	if thrkafkaConsum != nil && thrkafkaConsum.RunBase.Active {
		thrkafkaConsum.Kill() // thread종료시
	}

	if thrkafkaProduce != nil && thrkafkaProduce.RunBase.Active {
		thrkafkaProduce.Kill() // thread종료시
	}

	// clear pid
	process.Deregister(os.Getpid())
	am.Applog.Always("clearEnv pid (%d) (%d) (%v)", os.Getpid(), process.RunBase.ID, process.RunBase.Active)
	process.RunBase.Active = false

	// 공유메모리 확인하기 위해 잠시 대기
	time.Sleep(time.Millisecond * 200)

	// detach memory
	if shminst != nil {
		err := shminst.DeleteShm()
		if err != nil {
			am.Applog.Always("clearEnv.. DeleteShm err:", err)
		}
	}

	am.Applog.Always("Process quit [apisvr] [%d]! bye~~:)", os.Getpid())

	//로그파일 close
	am.Applog.Fileclose()

	// DBlogfile close
	//dbapp.LogClose()

	os.Exit(0)
}

// ------------------------------------------------------------------------------
// main
// ------------------------------------------------------------------------------
func main() {
	am.Applog.Info("Process start [apisvr]")

	var initOk bool = false
	initOk = initEnv()

	if initOk {
		appRun()
	}

	defer clearEnv()

	var num int = 0
	for {
		if !initOk || terminate {
			break
		}

		if num%1000 == 0 {
			am.Applog.Print(2, "api main run.. [%d]", num)
		}
		num++

		// check SystemState
		checkSystemState()

		manageDebug()
		manageRoutine()
		manageLogfile()

		process.RunBase.UpdateRunInfo()
		time.Sleep(time.Millisecond * 3000)
	}

	am.Applog.Info("Process end.. [apisvr]")
}
