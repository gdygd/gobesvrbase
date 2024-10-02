package kafkaapp

import (
	"apisvr/app/am"
	"fmt"
	"time"

	"github.com/IBM/sarama"
	"github.com/gdygd/goglib"
)

// func (k *KafkaHandler) SetCommEnv(broker KafkaBroker) {
// 	k.broker.addr = broker.addr
// 	k.broker.port = broker.port
// }

// // ---------------------------------------------------------------------------
// // ResetCommEnv
// // ---------------------------------------------------------------------------
// func (k *KafkaHandler) ResetCommEnv() {

// 	k.m_message = []byte{}

// 	k.msg_q = dt.NewSlQueue[MessageQue]()
// 	k.stateTimer = time.Now()
// }

func (k *KafkaHandler) ConnectProduce() error {
	var err error = nil

	am.Kfklog.Print(2, "Kfk Connect..#1")
	if !goglib.CheckElapsedTime(&k.connectTimer, CONNECT_INTERVAL) {
		return nil
	}

	am.Kfklog.Print(2, "Kfk Connect..#2")
	// Kafka 브로커 주소
	addr := fmt.Sprintf("%s:%d", k.broker.addr, k.broker.port)
	brokers := []string{addr}

	// Sarama 설정
	k.config = sarama.NewConfig()
	k.config.Producer.Return.Successes = true

	k.client, err = sarama.NewClient(brokers, k.config)
	if err != nil {
		am.Kfklog.Error("Kafka 클라이언트 생성 실패: %v\n", err)
		return err
	}

	k.producer, err = sarama.NewSyncProducerFromClient(k.client)
	if err != nil {
		am.Kfklog.Error("Kafka 프로듀서 생성 실패: %v\n", err)
		return err
	}

	am.Kfklog.Print(2, "Kfk Connect..#3")
	k.ResetCommEnv()
	// connect broker server
	k.isconnected = true
	return nil
}

func (k *KafkaHandler) ManageTx() bool {
	am.Kfklog.Print(2, "Kfk ManageTx #1")

	// check isconnected
	if !k.isconnected {
		// check connect interval and connect
		am.Kfklog.Print(2, "Kfk ManageTx #1.1")
		err := k.ConnectProduce()
		if err != nil {
			return false
		}
		return true
	}
	am.Kfklog.Print(2, "Kfk ManageTx #2")

	if goglib.CheckElapsedTime(&k.stateTimer, STATE_INTERVAL) {
		k.SendTest()
	}

	am.Kfklog.Print(2, "Kfk ManageTx #3")
	// get message and send message
	msg, isexist := k.msg_q.Pop()

	if isexist {
		isok := k.SendMsg(msg)

		if !isok {
			k.CloseProduce()
		}
	}
	am.Kfklog.Print(2, "Kfk ManageTx #4")

	// check send state

	// if state is false
	// 	>> Close

	return true
}

// ------------------------------------------------------------------------------
// SendMsg
// ------------------------------------------------------------------------------
func (k *KafkaHandler) SendMsg(msg MessageQue) bool {
	message := &sarama.ProducerMessage{
		Topic: "test_topic",
		//Value: sarama.StringEncoder("Hello, Kafka!"),
		Value: sarama.ByteEncoder(msg.M_message),
	}

	// 메시지 전송
	partition, offset, err := k.producer.SendMessage(message)
	if err != nil {
		am.Kfklog.Error("메시지 전송 실패: %v\n", err)
		return false
	}

	am.Kfklog.Print(2, "메시지 전송 성공: 파티션=%d, 오프셋=%d\n", partition, offset)

	return true
}

func (k *KafkaHandler) SendTest() {
	// send only time_t
	var info []byte = make([]byte, MAX_BUFLEN)
	var idx int = 0

	curtm := time.Now()
	utcsecs := curtm.Unix()

	// time
	goglib.SetNumber64(info, idx, utcsecs, 8, goglib.ED_BIG)
	idx += 8

	k.PushMsg(PT_TEST, info[:idx])
}

// ------------------------------------------------------------------------------
// PushMsg
// ------------------------------------------------------------------------------
func (k *KafkaHandler) PushMsg(code byte, buf []byte) {
	am.Netlog.Print(2, "PushMsg ..")

	var msg MessageQue = MessageQue{
		Code:      code,
		M_message: buf,
	}

	k.msg_q.Push(msg)
}

func (k *KafkaHandler) CloseProduce() {
	k.clearClient()
	k.clearProducer()
	k.isconnected = false
}

// func (k *KafkaHandler) clearClient() {
// 	if k.client != nil && !k.client.Closed() {
// 		k.client.Close()
// 	}
// 	k.client = nil
// }
// func (k *KafkaHandler) clearProducer() {
// 	if k.producer != nil {
// 		k.producer.Close()
// 	}
// 	k.producer = nil
// }
