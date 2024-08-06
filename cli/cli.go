package main

import (
	"besvrbase/cli/cmd"
	"besvrbase/sharedObj"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"unsafe"

	"github.com/gdygd/goshm/shmlinux"
)

// ------------------------------------------------------------------------------
// const
// ------------------------------------------------------------------------------
const MemorySize = 1024 * 100 // 10kb
const systemini = "./system.ini"
const skey = 0x1234 // shared memory key

// ------------------------------------------------------------------------------
// local
// ------------------------------------------------------------------------------
var SharedMem *sharedObj.SharedMemory
var terminate bool = false

var shminst *shmlinux.Linuxshm

// ------------------------------------------------------------------------------
// sigHandler
// ------------------------------------------------------------------------------
func sigHandler(chSig chan os.Signal) {
	for {
		signal := <-chSig
		fmt.Printf("Accept Signal %d", signal)
		switch signal {
		case syscall.SIGHUP: // 터미널 연결 끊겼을경우
			fmt.Printf("SIGHUP(%d)\n", signal)
			terminate = true
			clearEnv()
		case syscall.SIGINT:
			fmt.Printf("SIGINT(%d)\n", signal)
			terminate = true
			clearEnv()
		case syscall.SIGTERM:
			fmt.Printf("SIGTERM(%d)\n", signal)
			terminate = true
			clearEnv()
		case syscall.SIGKILL:
			fmt.Printf("SIGKILL(%d)\n", signal)
			terminate = true
			clearEnv()
		default:
			fmt.Printf("Unknown signal(%d)\n", signal)
			terminate = true
			clearEnv()
		}
	}
}

// ------------------------------------------------------------------------------
// initSignal
// ------------------------------------------------------------------------------
func initSignal() {
	fmt.Printf("iniSignal...\n")
	// signal handler
	ch_signal := make(chan os.Signal, 10)
	signal.Notify(ch_signal, syscall.SIGSEGV, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGUSR1)
	go sigHandler(ch_signal)
}

// ------------------------------------------------------------------------------
// initMemory
// ------------------------------------------------------------------------------
func initMemory2() bool {
	shminst := shmlinux.NewLinuxShm()
	shminst.InitShm(sharedObj.MEM_KEY, sharedObj.MEM_SIZE)

	err := shminst.CreateShm()
	if err != nil {
		fmt.Println("initMemory2 CreateShm err : ", err)
	}
	err = shminst.AttachShm()
	if err != nil {
		fmt.Println("initMemory2 AttachShm err : ", err)
	}
	SharedMem = (*sharedObj.SharedMemory)(unsafe.Pointer(shminst.Addr))
	return true
}

// ------------------------------------------------------------------------------
// initEnv
// ------------------------------------------------------------------------------
func initEnv() bool {
	initSignal()

	// if !initMemory() {
	// 	fmt.Printf("Share memory open fail..\n")
	// 	return false
	// }

	if !initMemory2() {
		fmt.Printf("Share memory open fail..\n")
		return false
	}

	return true
}

// ------------------------------------------------------------------------------
// clearEnv
// ------------------------------------------------------------------------------
func clearEnv() {

	// detach memory
	fmt.Printf("[cli] memory detach (1)\n")

	if shminst != nil {
		err := shminst.DeleteShm()
		if err != nil {
			fmt.Println("clearEnv.. DeleteShm err:", err)
		}
	}

	fmt.Printf("[cli] memory detach (2)\n")
	os.Exit(0)
}

func main() {
	fmt.Println("Command Line Interface :)\n")

	var initOk bool = false
	initOk = initEnv()

	// input command
	cmdStr := make([]string, 100)
	command := make([]interface{}, 100)
	for i := range cmdStr {
		command[i] = &cmdStr[i]
	}

	cli := cmd.NewCLI()
	cli.InitialMessage()

	if initOk {
		cli.SetShmMemory(SharedMem)
	}

	//defer clearEnv()
	for !cli.Exit && !terminate {

		if !initOk || terminate {
			fmt.Println("Exit Command Line Interface bye~~ :)\n")
			break
		}

		if terminate {
			fmt.Println("terminate Command Line Interface bye~~ :)\n")
			break
		}

		fmt.Printf("CLI>>")
		count, _ := fmt.Scanln(command...)

		if count > 9 {
			fmt.Println("Invalid command")
			continue
		}

		if count == 0 {
			continue
		}

		cli.SetCommand(cmdStr[0:count])
		cli.PrintCmd()

		cli.Run()

		if cli.Terminate {
			fmt.Println("cli process terminate")
			break
		}
	}
	clearEnv()
}
