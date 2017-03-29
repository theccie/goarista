// Copyright (C) 2017  Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package kafka

import (
	"expvar"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/Shopify/sarama"
	"github.com/aristanetworks/glog"
	"github.com/aristanetworks/goarista/monitor"
	"github.com/golang/protobuf/proto"
)

// MessageEncoder is an encoder interface
// which handles encoding proto.Message to sarama.ProducerMessage
type MessageEncoder interface {
	Encode(proto.Message) ([]*sarama.ProducerMessage, error)
	HandleSuccess(*sarama.ProducerMessage)
	HandleError(*sarama.ProducerError)
}

// BaseEncoder implements MessageEncoder interface
// and mainly handle monitoring
type BaseEncoder struct {
	// Used for monitoring
	histogram    *monitor.Histogram
	numSuccesses monitor.Uint
	numFailures  monitor.Uint
}

// counter counts the number Sysdb clients we have, and is used to guarantee that we
// always have a unique name exported to expvar
var counter uint32

// NewBaseEncoder returns a new base MessageEncoder
func NewBaseEncoder(typ string) *BaseEncoder {

	// Setup monitoring structures
	histName := "kafkaProducerHistogram_" + typ
	statsName := "messagesStats"
	if id := atomic.AddUint32(&counter, 1); id > 1 {
		histName = fmt.Sprintf("%s-%d", histName, id)
		statsName = fmt.Sprintf("%s-%d", statsName, id)
	}
	hist := monitor.NewHistogram(histName, 32, 0.3, 1000, 0)
	e := &BaseEncoder{
		histogram: hist,
	}

	statsMap := expvar.NewMap(statsName)
	statsMap.Set("successes", &e.numSuccesses)
	statsMap.Set("failures", &e.numFailures)

	return e
}

// Encode encodes the proto message to a sarama.ProducerMessage
func (e *BaseEncoder) Encode(message proto.Message) ([]*sarama.ProducerMessage,
	error) {
	// doesn't do anything, but keep it in order for BaseEncoder
	// to implement MessageEncoder interface
	return nil, nil
}

// HandleSuccess process the metadata of messages from kafka producer Successes channel
func (e *BaseEncoder) HandleSuccess(msg *sarama.ProducerMessage) {
	metadata := msg.Metadata.(Metadata)
	// TODO: Add a monotonic clock source when one becomes available
	e.histogram.UpdateLatencyValues(metadata.StartTime, time.Now())
	e.numSuccesses.Add(uint64(metadata.NumMessages))
}

// HandleError process the metadata of messages from kafka producer Errors channel
func (e *BaseEncoder) HandleError(msg *sarama.ProducerError) {
	metadata := msg.Msg.Metadata.(Metadata)
	// TODO: Add a monotonic clock source when one becomes available
	e.histogram.UpdateLatencyValues(metadata.StartTime, time.Now())
	glog.Errorf("Kafka Producer error: %s", msg.Error())
	e.numFailures.Add(uint64(metadata.NumMessages))

}