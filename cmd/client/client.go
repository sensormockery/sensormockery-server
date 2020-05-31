package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	apiV1 "github.com/sensormockery/sensormockery-server/pkg/api/v1"
	"github.com/sensormockery/sensormockery-server/pkg/stream"
	"github.com/streadway/amqp"
)

const (
	// BrokerURL is the broker url.
	BrokerURL = "amqp://wlkwdorl:qzRONejPv8GmBiW5-S2TNyYbrTCNL71U@roedeer.rmq.cloudamqp.com/wlkwdorl"
	// APIURL is an env var for the server url.
	APIURL = "API_URL"
)

func main() {
	// send subcribe request to server
	id := sendSubscribeRequest()

	// set up connection with message broker
	conn, channel, dpChan := setUpMessageBroker(id)
	defer conn.Close()
	defer channel.Close()

	// read data points from message brocker queue
	readDataPoints(dpChan)

	// send unsibscribe request to server
	sendUnsubscribeRequest(id)
}

func sendSubscribeRequest() int {
	stream := &apiV1.Stream{
		WaveType:   "sine",
		Sensor:     "accelerometer",
		NoiseCoeff: 0.1,
		BrokerURL:  BrokerURL,
	}

	req, err := json.Marshal(stream)
	if err != nil {
		handleError(err, "Cannot create object")
	}

	apiV1URL := fmt.Sprintf("%s%s", os.Getenv(APIURL), apiV1.APIPrefix)

	resp, err := http.Post(apiV1URL+apiV1.CreateStreamPath, "application/json", bytes.NewBuffer(req))
	if err != nil {
		handleError(err, "Error sending stream creation request")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		handleError(err, "Error reading response")
	}
	defer resp.Body.Close()

	createStreamResp := &apiV1.CreateStreamResp{}
	if err := json.Unmarshal(body, createStreamResp); err != nil {
		handleError(err, "Error reading response")
	}

	return createStreamResp.ID
}

func setUpMessageBroker(id int) (*amqp.Connection, *amqp.Channel, <-chan amqp.Delivery) {
	// connect to message broker
	conn, err := amqp.Dial(BrokerURL)
	handleError(err, "Can't connect to AMQP")

	// set up channel
	amqpChannel, err := conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	// set up data points queue
	queue, err := amqpChannel.QueueDeclare(fmt.Sprintf("stream-%d", id), false, true, false, false, nil)
	handleError(err, "Could not connect to queue")

	// register a consumer
	dataPointsChannel, err := amqpChannel.Consume(
		queue.Name,
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	handleError(err, "Could not register consumer")

	return conn, amqpChannel, dataPointsChannel
}

func readDataPoints(dataPointsChannel <-chan amqp.Delivery) {
	// read each sensor data point
	log.Printf("Client ready to receive sensor data, PID: %d", os.Getpid())
	dataPointsCount := 0
	for dataPointMessage := range dataPointsChannel {
		if dataPointsCount == 100 {
			return
		}

		// parse data point message from queue
		dataPoint := &stream.DataPoint{}
		err := json.Unmarshal(dataPointMessage.Body, dataPoint)
		handleError(err, "Error decoding JSON")

		log.Printf("Received data point: %f", dataPoint.Data)

		// confirm that sensor task is successfully performed
		if err := dataPointMessage.Ack(false); err != nil {
			log.Printf("Error acknowledging data point : %s", err)
		}

		dataPointsCount++
	}
}

func sendUnsubscribeRequest(id int) {
	reqBody := &apiV1.CreateStreamResp{
		ID: id,
	}

	reqBodyAsJSON, err := json.Marshal(reqBody)
	handleError(err, "Error building request body")

	apiV1URL := fmt.Sprintf("%s%s", os.Getenv(APIURL), apiV1.APIPrefix)

	client := &http.Client{}

	req, err := http.NewRequest("DELETE", apiV1URL+apiV1.DeleteStreamPath, strings.NewReader(string(reqBodyAsJSON)))
	handleError(err, "Error creating request")

	resp, err := client.Do(req)
	handleError(err, "Error sending request")
	defer resp.Body.Close()
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
