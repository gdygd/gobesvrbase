package kafkaapp

import (
	"apisvr/app/am"
	"time"

	"github.com/gdygd/goglib"
)

func initKafkaProducer(k *KafkaHandler) bool {

	return true
}

func clearKfkProducer(k *KafkaHandler) {
	am.Kfklog.Print(2, "clearKfkProducer...")
	k.CloseProduce()
}

// ------------------------------------------------------------------------------
// THRKafkaProduce
// ------------------------------------------------------------------------------
func THRKafkaProduce(t *goglib.Thread, chThrStop chan bool, arg1, arg2, arg3 interface{}) {
	am.Kfklog.Print(2, "THRKafkaProduce start...")

	k, ok := arg1.(*KafkaHandler)

	var initOk bool = false
	initOk = initKafkaProducer(k)
	defer clearKfkProducer(k)

	if !ok {
		am.Kfklog.Error("THRKafkaProduce arg error... ")
		return
	}
	if !initOk {
		am.Kfklog.Error("initKafkaProducer fail... ")
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
		}
		if terminate {
			break
		}

		k.ManageTx()

		t.RunBase.UpdateRunInfo()

		time.Sleep(time.Millisecond * 1000)

	}
	am.Kfklog.Print(2, "THRKafkaProduce exit..")
}
