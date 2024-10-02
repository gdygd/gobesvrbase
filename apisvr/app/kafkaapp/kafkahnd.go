package kafkaapp

import (
	"time"

	dt "github.com/gdygd/goglib/datastructure"
)

func (k *KafkaHandler) SetCommEnv(broker KafkaBroker) {
	k.broker.addr = broker.addr
	k.broker.port = broker.port
}

// ---------------------------------------------------------------------------
// ResetCommEnv
// ---------------------------------------------------------------------------
func (k *KafkaHandler) ResetCommEnv() {

	k.m_message = []byte{}

	k.msg_q = dt.NewSlQueue[MessageQue]()
	k.stateTimer = time.Now()
	//k.connectTimer = time.Now()
	k.recvtm = time.Now()
}

func (k *KafkaHandler) clearClient() {
	//am.Kfklog.Print(2, "clearClient #1")
	if k.client != nil && !k.client.Closed() {
		k.client.Close()
		//am.Kfklog.Print(2, "clearClient #1")
	}
	k.client = nil
}
func (k *KafkaHandler) clearProducer() {
	//am.Kfklog.Print(2, "clearProducer #1")
	if k.producer != nil {
		k.producer.Close()
		//am.Kfklog.Print(2, "clearProducer #2")
	}
	k.producer = nil
}

func (k *KafkaHandler) clearConsumer() {

	//am.Kfklog.Print(2, "clearConsumer #1")
	if k.consumer != nil {
		k.consumer.Close()
		//am.Kfklog.Print(2, "clearConsumer #2")
	}
	k.consumer = nil
}

func (k *KafkaHandler) clearPartitionConsumer() {
	//am.Kfklog.Print(2, "clearPartitionConsumer#1")
	if k.partitionConsumer != nil {
		k.partitionConsumer.Close()
		//am.Kfklog.Print(2, "clearPartitionConsumer#2")
	}
	k.partitionConsumer = nil
}
