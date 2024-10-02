package kafkaapp

import (
	"apisvr/app/am"
	"fmt"
	"log"
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
// 	recvtm = time.Now()
// }

func (k *KafkaHandler) ConnectConsume() error {

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

	log.Println("#1 connectK")
	k.client, err = sarama.NewClient(brokers, k.config)
	if err != nil {
		log.Printf("Kafka 클라이언트 생성 실패: %v\n", err)
		return err
	}
	log.Println("#2 connectK")

	k.consumer, err = sarama.NewConsumerFromClient(k.client)
	if err != nil {
		log.Fatalf("Kafka 컨슈머 생성 실패: %v", err)
		return err
	}
	log.Println("#3 connectK")

	// 특정 토픽의 파티션 컨슈머 생성
	k.partitionConsumer, err = k.consumer.ConsumePartition("test_topic", 0, sarama.OffsetNewest)
	if err != nil {
		log.Fatalf("파티션 컨슈머 생성 실패: %v", err)
		return err
	}
	log.Println("#4 connectK")

	k.ResetCommEnv()
	k.isconnected = true
	return nil
}

func (k *KafkaHandler) ManageLine() bool {
	if !k.isconnected {
		am.Kfklog.Print(2, "Kfk ManageRx connected..")
		k.ConnectConsume()
	}

	return true
}

func (k *KafkaHandler) ManageRx() {

	if k.isconnected {
		select {
		case message := <-k.partitionConsumer.Messages():
			am.Kfklog.Print(2, "received message: %v", message)
			k.recvtm = time.Now()
		default:
			//am.Kfklog.Print(2, "Not received message")
		}
	}
	curtm := time.Now()
	curSec := curtm.Unix()
	prevRecvSec := k.recvtm.Unix()

	if (curSec - prevRecvSec) > 5 {
		am.Kfklog.Warn("There aren't any received message within 5sec elapsed sed[%d]", (curSec - prevRecvSec))
		// close and connect

		k.CloseConsume()
	}

}

func (k *KafkaHandler) CloseConsume() {
	am.Kfklog.Print(2, "Close..#1")
	k.clearClient()
	//am.Kfklog.Print(2, "Close..#2")
	k.clearConsumer()
	//am.Kfklog.Print(2, "Close..#3")
	k.clearPartitionConsumer()
	//am.Kfklog.Print(2, "Close..#4")
	k.isconnected = false
}

// func (k *KafkaHandler) clearClient() {
// 	if k.client != nil && !k.client.Closed() {
// 		k.client.Close()
// 		am.Kfklog.Print(2, "clearClient")
// 	}
// 	k.client = nil
// }
// func (k *KafkaHandler) clearConsumer() {

// 	if k.consumer != nil {
// 		k.consumer.Close()
// 		am.Kfklog.Print(2, "clearConsumer")
// 	}
// 	k.consumer = nil
// }

// func (k *KafkaHandler) clearPartitionConsumer() {
// 	if k.partitionConsumer != nil {
// 		k.partitionConsumer.Close()
// 		am.Kfklog.Print(2, "clearPartitionConsumer")
// 	}
// 	k.partitionConsumer = nil
// }
