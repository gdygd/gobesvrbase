package kafkaapp

import (
	"time"

	//"github.com/Shopify/sarama"
	"github.com/IBM/sarama"
	dt "github.com/gdygd/goglib/datastructure"
)

const (
	PT_TEST          = 0x99
	STATE_INTERVAL   = 1000
	CONNECT_INTERVAL = 3000
)

const (
	ComRecvBuf    = 1
	ComSendBuf    = 100
	CHBUF_NET_PRC = 300

	MAX_BUFLEN = 102400
)

// ---------------------------------------------------------------------------
// Recevied Kafka message for sending msgapp
// ---------------------------------------------------------------------------
type KafkakMSG struct {
	Id     int
	Code   byte
	Data   []byte    // 통신데이터
	RecvTm time.Time // 수신 시간
}

type MessageQue struct {
	Code      byte
	M_message []byte
	Seq       byte
}

type KafkaBroker struct {
	addr string
	port int
}

type KafkaHandler struct {
	config            *sarama.Config
	client            sarama.Client
	producer          sarama.SyncProducer
	consumer          sarama.Consumer
	partitionConsumer sarama.PartitionConsumer
	broker            KafkaBroker

	msg_q *dt.SlQueye[MessageQue] // message for producer

	m_message   []byte
	isconnected bool // broker server connect state

	stateTimer   time.Time
	connectTimer time.Time

	recvtm time.Time
}

func MakeKfkHandler(addr string, port int) *KafkaHandler {
	k := &KafkaHandler{}
	k.broker.addr = addr
	k.broker.port = port

	return k
}
