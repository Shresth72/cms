package main

import (
	"encoding/json"
	"net/http"

	"github.com/IBM/sarama"
)

// Handlers
type SendMessageRequest struct {
	Message string `json:"message"`
}

func (app *application) sendMessageToKafka(w http.ResponseWriter, r *http.Request) {
  var body SendMessageRequest
	err := app.readJSON(w, r, &body)
	if err != nil {
		app.badRequestError(w, r, err)
		return
	}

  message := body.Message
  topic := app.cfg.kafka.topic
  err = app.produceMessage(topic, message)
  if err != nil {
    app.serverError(w, r, err, "failed to store data in kafka")
    return
  }

	err = app.writeJSON(w, r, http.StatusOK, envelope{
		"message": "message uploaded successfully",
	})

	if err != nil {
		app.serverError(w, r, err)
	}
}

// Utils
type EncodingMessage struct {
  Title string `json:"title"`
  Url string `json:"url"`
}

func (app *application) pushVideoForEncodingToKafka(title string, url string) error {
  msg := &EncodingMessage{
    Title: title,
    Url: url,
  }
  messageBytes, err := json.Marshal(msg)
  if err != nil {
    return err
  }
  message := string(messageBytes)

  topic := app.cfg.kafka.topic
  err = app.produceMessage(topic, message)
  if err != nil {
    return err
  }

  return nil
} 

func (app *application) produceMessage(topic, message string) error {
  msg := &sarama.ProducerMessage{
    Topic: topic,
    Value: sarama.StringEncoder(message),
  }

  partition, offset, err := app.producer.SendMessage(msg)
  if err != nil {
    return err
  }

  app.logger.Info().Msgf("Data is successfully stored with unique identifier %v/%d/%d", topic, partition, offset)
  return nil
}
