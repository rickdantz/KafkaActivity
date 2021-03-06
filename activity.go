package KafkaActivity

// Imports all of the flowGo binaries
import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/optiopay/kafka"
	"github.com/optiopay/kafka/proto"
	"github.com/Shopify/sarama"
)

// log is the default package logger
var log = logger.GetLogger("activity-tibco-kafka")

// Construct Input Params
const (
	topic     = "topic"
	message   = "message"
	partition = 0
)

//NBCU Kafka Dev Servers
var kafkaAddrs = []string{"ushapld00119la:9092", "ushapld00119la:9092"}

// KafkaActivity is a Kafka Activity implementation
type KafkaActivity struct {
	metadata *activity.Metadata
	syncProducerMap *map[string]sarama.SyncProducer
}

// init create & register activity
// func init() {
// 	md := activity.NewMetadata(jsonMetadata)
//	activity.Register(&KafkaActivity{metadata: md})
// }

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	log.Debug("KafkaActivity NewActivity")
	pKafkaActivity := &KafkaActivity{metadata: metadata}
	producers := make(map[string]sarama.SyncProducer)
	pKafkaActivity.syncProducerMap = &producers
	return pKafkaActivity
}


// Metadata implements activity.Activity.Metadata
func (a *KafkaActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *KafkaActivity) Eval(context activity.Context) (done bool, err error) {

	topicInput := context.GetInput(topic).(string)

	messageInput := context.GetInput(message).(string)

	conf := kafka.NewBrokerConf("NBCU-FloGo-Client")
	conf.AllowTopicCreation = true

	// connect to kafka cluster
	broker, err := kafka.Dial(kafkaAddrs, conf)
	if err != nil {
		log.Error("cannot connect to kafka cluster:", err)
	}
	defer broker.Close()

	// Connect & Send Message to Kafka
	producer := broker.Producer(kafka.NewProducerConf())

	msg := &proto.Message{Value: []byte(messageInput)}

	resp, err := producer.Produce(topicInput, partition, msg)

	if err != nil {
		log.Error("Error sending message to Kafka broker:", err)
	}

	// if log.IsEnabledFor(log.DEBUG) {
	log.Debug("Response:", resp)
	// }

	return true, nil
}
