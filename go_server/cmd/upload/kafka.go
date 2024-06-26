package main

import (
	"time"

	"github.com/IBM/sarama"
)

func (app *application) newDataCollector(brokerList []string, version sarama.KafkaVersion) sarama.SyncProducer  {
  config := sarama.NewConfig()
  config.Version = version
  config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
  // TODO: Increase in-sync replicas with `min.insync.replicas`
  config.Producer.Retry.Max = 10
  config.Producer.Return.Successes = true

  // TODO: Add TLS Configuration
  // tlsConfig := createTlsConfiguration()
  // if tlsConfig != nil {
  //   config.Net.TLS.Config = tlsConfig
  //   config.Net.TLS.Enable = true
  // }
  
  producer, err := sarama.NewSyncProducer(brokerList, config)
  if err != nil {
    app.logger.Panic().Err(err).Msg("Failed to start Sarama Producer")
  }
  return producer
}

func (app *application) newLogProducer(brokerList []string, version sarama.KafkaVersion) sarama.AsyncProducer {
  config := sarama.NewConfig()
  config.Version = version
  config.Producer.RequiredAcks = sarama.WaitForLocal // only wait for leader to ack 
  config.Producer.Compression = sarama.CompressionSnappy
  config.Producer.Flush.Frequency = 500 * time.Millisecond

  producer, err := sarama.NewAsyncProducer(brokerList, config)
  if err != nil {
    app.logger.Panic().Err(err).Msg("Failed to start Sarama LogProducer")
  }
  
  // Log Error if not able to produce messages
  go func ()  {
    for err := range producer.Errors() {
      app.logger.Err(err).Msg("Failed to write access log entry")
    }
  }()

  return producer
}
