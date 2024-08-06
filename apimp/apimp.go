package main

import (
	"besvrbase/sharedObj"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
	"unsafe"

	"github.com/gdygd/goglib"

	"github.com/gdygd/goshm/shmlinux"
	"gopkg.in/ini.v1"
)

// ------------------------------------------------------------------------------
// struct
// ------------------------------------------------------------------------------
type appVariable struct {
	DbHost   string
	DbPort   int
	DbUser   string
	DbName   string
	DbPasswd string
	DebugLv  int
}

// ------------------------------------------------------------------------------
// Glocal
// ------------------------------------------------------------------------------
var SharedMem *sharedObj.SharedMemory
var process *goglib.Process
var Mlog *goglib.OLog2 = goglib.InitLogEnv("./log", "apimp", 0)

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const systemini = "./system.ini"
const sysenvini = "./sys_env.ini"

// ------------------------------------------------------------------------------
// Local
// ------------------------------------------------------------------------------
var isAleayProcess bool = false
var PRC_DESC = []string{"apimp", "./apisvr"}
var prcArgv = [][]string{{""}, {""}}
var AppVar appVariable

var shminst *shmlinux.Linuxshm = nil

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	for {
		signal := <-chSig
		str := fmt.Sprintf("[apimp] Accept Signal : %d", signal)
		Mlog.Info("%s", str)
		switch signal {
		case syscall.SIGHUP:
			Mlog.Info("[apimp]SIGHUP(%d)\n", signal)
		case syscall.SIGINT:
			Mlog.Info("[apimp]SIGINT(%d)\n", signal)
			SharedMem.System.Terminate = true
		case syscall.SIGTERM:
			Mlog.Info("[apimp]SIGTERM(%d)\n", signal)
			SharedMem.System.Terminate = true
		case syscall.SIGKILL:
			Mlog.Info("[apimp]SIGKILL(%d)\n", signal)
			SharedMem.System.Terminate = true
		default:
			Mlog.Info("[apimp]Unknown signal(%d)\n", signal)
			SharedMem.System.Terminate = true
		}
	}
}

// ------------------------------------------------------------------------------
// initProcessDesc
// ------------------------------------------------------------------------------
func initProcessDesc() {
	Mlog.Always("initProcessDesc..")
	for idx := 0; idx < sharedObj.MAX_PROCESS; idx++ {
		SharedMem.Process[idx] = goglib.InitProcess(PRC_DESC[idx], prcArgv[idx])
	}
}

// ------------------------------------------------------------------------------
// initEnvVaiable
// ------------------------------------------------------------------------------
func initEnvVaiable() bool {
	Mlog.Always("[apimp]initEnvVaiable..")

	cfg, err := ini.Load(sysenvini)
	if err != nil {
		Mlog.Error("fail to read sysenvini.ini %v", err)
		return false
	}

	AppVar.DbHost = cfg.Section("DATABASE").Key("host").String()
	AppVar.DbPort, _ = cfg.Section("DATABASE").Key("port").Int()
	AppVar.DbUser = cfg.Section("DATABASE").Key("user").String()
	AppVar.DbName = cfg.Section("DATABASE").Key("dbname").String()
	AppVar.DbPasswd = cfg.Section("DATABASE").Key("passwd").String()

	AppVar.DebugLv, _ = cfg.Section("OPRMODE").Key("debuglv").Int()

	Mlog.Always("[apimp]App Environment variable")

	Mlog.Always(" \t [Database]")
	Mlog.Always(" \t host:%s", AppVar.DbHost)
	Mlog.Always(" \t port:%d", AppVar.DbPort)
	Mlog.Always(" \t db  :%s", AppVar.DbName)
	Mlog.Always(" \t user:%s", AppVar.DbUser)
	Mlog.Always(" \t pw:%s", AppVar.DbPasswd)

	Mlog.Always(" \t [Logging]")
	Mlog.Always(" \t DebufLv:%d", AppVar.DebugLv)

	return true
}

// ------------------------------------------------------------------------------
// initProcess
// ------------------------------------------------------------------------------
func initProcess() bool {
	// process initialize & start

	Mlog.Info("[apimp]initProcess")
	process = &SharedMem.Process[sharedObj.PRC_IDX_MAIN]

	for idx := 0; idx < sharedObj.MAX_PROCESS; idx++ {
		Mlog.Info(PRC_DESC[idx])
		if idx == sharedObj.PRC_IDX_MAIN {
			SharedMem.Process[idx].RunBase.Active = true
			continue
		} else {
			Mlog.Info(">>%v", PRC_DESC[idx])
			SharedMem.Process[idx].Start2()
		}
	}

	return true
}

// ------------------------------------------------------------------------------
// initMemory
// ------------------------------------------------------------------------------
func initMemory() bool {
	Mlog.Info("[apimp]initMemory...#1")
	shminst := shmlinux.NewLinuxShm()
	shminst.InitShm(sharedObj.MEM_KEY, sharedObj.MEM_SIZE)

	err := shminst.CreateShm()
	if err != nil {
		Mlog.Error("initMemory CreateShm err : ", err)
	}
	err = shminst.AttachShm()
	if err != nil {
		Mlog.Error("initMemory AttachShm err : ", err)
	}

	Mlog.Info("[apimp]initMemory...#2")

	SharedMem = (*sharedObj.SharedMemory)(unsafe.Pointer(shminst.Addr))
	SharedMem.System.Terminate = false
	Mlog.Info("[apimp]initMemory...#3")
	return true
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {
	Mlog.Info("[apimp]iniSignal")
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

	Mlog.Info("[apimp]isAlreadyProcess...")

	var isrunning bool = false

	Mlog.Print(3, "prcdesc : %v", PRC_DESC[1])

	prcnm := PRC_DESC[1] // 코어프로세스
	cmdstr := fmt.Sprintf(`ps -ef | grep %s | grep -v grep`, prcnm)
	cmd := exec.Command("bash", "-c", cmdstr)
	output, _ := cmd.CombinedOutput()
	stroutput := string(output)

	Mlog.Print(3, "cmd output : %v %d", stroutput, len(stroutput))

	if len(stroutput) == 0 {
		isrunning = false
	} else {
		isrunning = true
	}

	return isrunning
}

func SetDebugLv() {
	Mlog.Info("[apimp]SetDebugLv...")
	process.SetDebugLv(AppVar.DebugLv)
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {
	// check process running...
	if isAlreadyProcess() {
		Mlog.Error("[apimp]프로세스 is running ...")
		isAleayProcess = true
		return false
	}

	Mlog.Info("[apimp]initEnv ...")

	if !initEnvVaiable() {
		Mlog.Error("[apimp] 환경변수 초기화FAIL...")
		return false
	}

	initSignal()

	if !initMemory() {
		Mlog.Info("Share memory created fail..")
		return false
	}

	initProcessDesc() // process name and process command
	if !initProcess() {
		Mlog.Info("Process initialize fail..")
		return false
	}

	Mlog.Info("[apimp]initEnv ok")

	// Register PID
	process.RegisterPid(os.Getpid())
	process.RunBase.Active = true

	SetDebugLv()

	return true
}

func manageProcess() {

	for idx := 1; idx < sharedObj.MAX_PROCESS; idx++ {
		ptrPrc := &SharedMem.Process[idx]

		//Mlog.Info("manageProcess : %d %v %s", ptrPrc.RunBase.ID, ptrPrc.RunBase.Active, ptrPrc.PrcName)

		if ptrPrc.RunBase.Active {
			//Mlog.Info("Pid : ", ptrPrc.GetPid())

			// running check
			if !ptrPrc.RunBase.Active {
				continue
			}
			var state int = 0
			if ptrPrc.IsRunning(&state) {
				continue
			}

			switch state {
			case goglib.RST_ABNOMAL:
				Mlog.Warn("RST_ABNOMAL [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
				Mlog.Warn("Processkill [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
				ptrPrc.Kill()
			case goglib.RST_UNEXIST:
				Mlog.Warn("RST_UNEXIST [%s %d]", ptrPrc.PrcName, ptrPrc.GetPid())
				Mlog.Warn("Process start [%s]", ptrPrc.PrcName)
				ptrPrc.Start2()
			default:
				break
			}
		} else {
			Mlog.Print(1, "Not Active %s %d", ptrPrc.PrcName, ptrPrc.GetPid())
		}
	}
}

// ------------------------------------------------------------------------------
// clearEnv
// ------------------------------------------------------------------------------
func clearEnv() {
	if isAleayProcess {
		Mlog.Print(6, "[apimp] Is Aleady process")
		Mlog.Print(6, "[apimp] Process quit, byebye~ :)")
		return
	}

	// All sub process quit
	for idx := 1; idx < sharedObj.MAX_PROCESS; idx++ {
		ptrPrc := &SharedMem.Process[idx]
		if ptrPrc == nil {
			continue
		}

		if ptrPrc.RunBase.Active {
			Mlog.Always("Kill Process:[%v][%d] [%s]", ptrPrc.RunBase.Active, ptrPrc.RunBase.ID, ptrPrc.PrcName)
			ptrPrc.Kill()
		}
	}

	var isAllClose bool = true
	var prcIdx int = 1
	for {
		ptrPrc := &SharedMem.Process[prcIdx]
		if ptrPrc == nil {
			continue
		}
		//Mlog.Info("Child process state:[%v] [%d] [%s]", ptrPrc.RunBase.Active, ptrPrc.RunBase.ID, ptrPrc.PrcName)
		// child process yet active state..

		if ptrPrc.RunBase.ID != 0 {
			isAllClose = false
		}

		if prcIdx == sharedObj.MAX_PROCESS-1 {
			if isAllClose {
				break
			} else {
				isAllClose = true
				break
			}
		}

		prcIdx++
		time.Sleep(time.Millisecond * 100)
	}

	Mlog.Print(2, "All child process close(2)")
	time.Sleep(time.Millisecond * 500)

	// detach memory
	if shminst != nil {
		err := shminst.DeleteShm()
		if err != nil {
			Mlog.Always("clearEnv.. DeleteShm err:", err)
		}
	}

	Mlog.Print(2, "memery destroy")

	Mlog.Print(2, "[apimp] Process quit, byebye~ :)")

	// 로그파일 close
	//pmlog := &Mlog
	Mlog.Fileclose()
}

func checkLogDir() {
	// log directory 존재 확인
	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		// root dir 생성
		os.Mkdir("./log", os.ModePerm)
	}
}

func main() {
	var initOk bool = false
	checkLogDir()
	Mlog.Info("%s", "[apimp] Process start")

	initOk = initEnv()

	Mlog.Print(4, "Mng.. %v, %v", initOk, SharedMem.System.Terminate)

	for {
		if !initOk || SharedMem.System.Terminate {
			break
		}
		// manage process
		manageProcess()

		// update process run info
		process.RunBase.UpdateRunInfo()

		time.Sleep(time.Millisecond * 1000)
	}

	defer clearEnv()

	Mlog.Info("[apimp] Process end..")

}
