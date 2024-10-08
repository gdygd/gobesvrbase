package am

import (
	"github.com/gdygd/goglib"
)

// ------------------------------------------------------------------------------
// struct
// ------------------------------------------------------------------------------
type appVariable struct {
	Version string

	NetAddr string
	NetPort int

	DbHost       string
	DbPort       int
	DbUser       string
	DbName       string
	DbPasswd     string
	DbVer        string
	Dbms         string
	DbDriverName string

	HttpPort     int
	Sslcrt       string
	Sslkey       string
	Sslcertpem   string
	Domain       string
	Https        string // "yes" or "no"
	Alloworigins []string

	DebugLv   int
	KafkaMode string // "yes or no"	// yes : kafka app run, else not run
}

// ---------------------------------------------------------------------------
// Global variable
// ---------------------------------------------------------------------------
// Log
// ---------------------------------------------------------------------------
var Applog *goglib.OLog2 = goglib.InitLogEnv("./log", "apiapp", 2)
var Netlog *goglib.OLog2 = goglib.InitLogEnv("./log", "netapp", 2)
var Kfklog *goglib.OLog2 = goglib.InitLogEnv("./log", "kafkaapp", 2)

// ---------------------------------------------------------------------------
// App Variable
// ---------------------------------------------------------------------------
var AppVar appVariable
