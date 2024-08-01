package cmd

import (
	"besvrbase/sharedObj"
	"fmt"
	"os"
	"strconv"
	"time"
)

const versionpath = "./version.txt"

// ----------------------------------------------------------------------
// Process start message
// ----------------------------------------------------------------------
func (c *CLI) InitialMessage() {
	fmt.Println("VERSION:")
	fmt.Println("	1.0")
	fmt.Println("\n")

	fmt.Println("USAGE:")
	fmt.Println("$ CLI>>[COMMAND]")
	fmt.Println("\n")

	fmt.Println("COMMANDS:")
	fmt.Println("	help")
	fmt.Println("	system")
	fmt.Println("	version")
	fmt.Println("	process")
	fmt.Println("	debug")
	fmt.Println("	exit")
	fmt.Println("	terminate")
	fmt.Println("\n")
}

func helpDebug() {
	fmt.Println(" #process debug info and set debug mode")
	fmt.Println("USAGE: ")
	fmt.Println("	$ debug {process} {lv}")
	fmt.Println("		ex) $ debug processname 3")

}

func helpExit() {
	fmt.Println("cli exit")
}

func helpTerminate() {
	fmt.Println("osvms process terminate")
}

func (c *CLI) help() {
	arglen := len(c.args)
	if arglen > 1 {
		fmt.Println("Invalid argument")
		return
	}

	if arglen == 0 {
		c.InitialMessage()
		return
	}

	arg := c.args[0]
	if arg == "exit" {
		helpExit()
	} else if arg == "debug" {
		helpDebug()
	} else if arg == "terminate" {
		helpTerminate()
	} else {
		fmt.Println("Invalid argument")
	}
}

func (c *CLI) showSystemInfo() {
	// process name, debug lv
	fmt.Printf("========================================================================\n")
	fmt.Printf("%-5s| %-3s| %-20s| \n", "PWR", "DB", "TIME")

	fmt.Printf("------------------------------------------------------------------------\n")
	ptrPrc := c.shmMem.System

	svrTime := time.Unix(int64(ptrPrc.SvrUtc), 0)
	strSvrtime := fmt.Sprintf("%04d-%02d-%02d  %02d:%02d:%02d", svrTime.Year(), svrTime.Month(), svrTime.Day(), svrTime.Hour(), svrTime.Minute(), svrTime.Second())

	fmt.Printf("%-5v| %-3d| %-20s| \n", ptrPrc.Terminate, ptrPrc.DbSvrComm, strSvrtime)

	fmt.Printf("========================================================================\n")
}

func (c *CLI) showVersionInfo() {
	// process name, debug lv
	dat, err := os.ReadFile(versionpath)
	strversion := fmt.Sprintf("%v", string(dat))
	if err != nil {
		strversion = ""
	}

	fmt.Printf("========================================================================\n")
	fmt.Printf("%-15s| \n", "VERSION")
	fmt.Printf("------------------------------------------------------------------------\n")
	fmt.Printf("%-15s| \n", strversion)

	fmt.Printf("========================================================================\n")
}

func (c *CLI) showProcessInfo() {
	// process name, debug lv
	fmt.Printf("===================================================\n")
	fmt.Printf("%-5s| %-10s| %-10s| %-3s| \n", "ACT", "PID", "PNAME", "lv")
	fmt.Printf("---------------------------------------------------\n")

	for idx := 0; idx < sharedObj.MAX_PROCESS; idx++ {
		ptrPrc := c.shmMem.Process[idx]
		var strName string = "*"
		strName = ptrPrc.GetPNameByArr()
		fmt.Printf("%-5v| %-10d| %-10s| %-3d| \n", ptrPrc.RunBase.Active, ptrPrc.RunBase.ID, strName, ptrPrc.DebugLv)
	}
	fmt.Printf("===================================================\n")
}

func (c *CLI) system() {
	arglen := len(c.args)
	if arglen > 1 {
		fmt.Println("Invalid exit command")
		return
	}

	c.showSystemInfo()
}

func (c *CLI) version() {
	arglen := len(c.args)
	if arglen > 1 {
		fmt.Println("Invalid exit command")
		return
	}

	c.showVersionInfo()
}

func (c *CLI) process() {
	arglen := len(c.args)
	if arglen > 1 {
		fmt.Println("Invalid exit command")
		return
	}

	c.showProcessInfo()
}

func (c *CLI) debug() {
	// debug prcname level

	arglen := len(c.args)
	if arglen != 2 {
		fmt.Println("Invalid debug command")
		return
	}

	var prcName string = c.args[0]
	var strLv string = c.args[1]
	lv, err := strconv.Atoi(strLv)
	if err != nil {
		fmt.Println("Invalid debug level")
		return
	}
	if lv < 0 || lv > 9 {
		fmt.Println("Invalid debug level")
		return
	}

	for idx := 0; idx < sharedObj.MAX_PROCESS; idx++ {
		ptrPrc := &c.shmMem.Process[idx]
		strName := ptrPrc.GetPNameByArr()

		if strName == prcName {
			fmt.Println("change debug level ", strName)
			ptrPrc.SetDebugLv(lv)
			break
		}
	}
}

func (c *CLI) exit() {
	arglen := len(c.args)
	if arglen > 1 {
		fmt.Println("Invalid exit command")
		return
	}
	c.Exit = true
}

func (c *CLI) terminate() {
	arglen := len(c.args)
	if arglen > 1 {
		fmt.Println("Invalid terminate command")
		return
	}

	c.shmMem.System.Terminate = true
	c.Terminate = true
}

func charArrToStr(data []byte) string {
	var s string = ""

	for i := 0; i < len(data); i++ {
		if data[i] == 0x00 {
			continue
		}

		ch := fmt.Sprintf("%c", data[i])
		s += ch
	}

	return s
}
