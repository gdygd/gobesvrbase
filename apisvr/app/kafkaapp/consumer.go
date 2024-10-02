package kafkaapp

import (
	"apisvr/app/am"
	"time"

	"github.com/gdygd/goglib"
)

func initKafkaConsumer(k *KafkaHandler) bool {

	return true
}

func clearKfkConsumer(k *KafkaHandler) {
	am.Kfklog.Print(2, "clearKfkConsumer...")
	k.CloseConsume()
}

// ------------------------------------------------------------------------------
// THRKafkaConsum
// ------------------------------------------------------------------------------
func THRKafkaConsum(t *goglib.Thread, chThrStop chan bool, arg1, arg2, arg3 interface{}) {
	am.Netlog.Print(2, "THRKafkaConsum start...")

	k, ok := arg1.(*KafkaHandler)

	var initOk bool = false
	initOk = initKafkaConsumer(k)
	defer clearKfkConsumer(k)

	if !ok {
		am.Kfklog.Error("THRKafkaConsum arg error... ")
		return
	}
	if !initOk {
		am.Kfklog.Error("THRKafkaConsum fail... ")
		return
	}

	var terminate = false
	for {

		select {
		case <-chThrStop:
			am.Netlog.Always("[1]thread stop,, Client tx thread quit..")
			terminate = true

			break
		default:
			//
			k.ManageLine()
			k.ManageRx()
		}
		if terminate {
			break
		}

		//am.Kfklog.Print(2, "run...#1")
		// k.ManageRx()

		t.RunBase.UpdateRunInfo()

		//am.Kfklog.Print(2, "run...#2")
		time.Sleep(time.Millisecond * 100)

	}
	am.Kfklog.Print(2, "THRKafkaConsum exit..")
}
