package mqtt

import (
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gustavoppereira/llv/agent/event"
	"log"
)

const DefaultQuiesce = 2000

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	log.Println("Connected!")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	log.Printf("Connection lost: %v", err)
}

type GlobalMessageHandler func(string, event.Event)

type ConnectionOptions struct {
	brokerHost string
	port       int
	clientID   string
	username   string
	password   string
}

func NewConnectionOptions(brokerHost string, port int, clientID string, username string, password string) *ConnectionOptions {
	return &ConnectionOptions{
		brokerHost: brokerHost,
		port:       port,
		clientID:   clientID,
		username:   username,
		password:   password,
	}
}

func (c *ConnectionOptions) NewClientOptions() *mqtt.ClientOptions {
	var clientOptions = mqtt.NewClientOptions()

	clientOptions.AddBroker(fmt.Sprintf("tcp://%s:%d", c.brokerHost, c.port))

	// Agent ID Randomly generated on Setup Phase
	clientOptions.SetClientID(c.clientID)
	clientOptions.SetUsername(c.username)
	clientOptions.SetPassword(c.password)

	return clientOptions
}

type Client struct {
	client *mqtt.Client
}

func (c *Client) Disconnect() {
	(*c.client).Disconnect(DefaultQuiesce)
}

func (c *Client) Unsubscribe(topic string) {
	(*c.client).Unsubscribe(topic)
}

func (c *Client) Publish(topic string, message string) {
	(*c.client).Publish(topic, 0, false, message)
}

func Connect(options *ConnectionOptions, messageHandler GlobalMessageHandler) (*Client, error) {

	clientOptions := options.NewClientOptions()

	clientOptions.SetDefaultPublishHandler(func(client mqtt.Client, message mqtt.Message) {
		messageHandler(message.Topic(), parseMqttEventFromMessage(message))
	})
	clientOptions.SetOnConnectHandler(connectHandler)
	clientOptions.SetConnectionLostHandler(connectLostHandler)

	mqttClient := mqtt.NewClient(clientOptions)
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{
		client: &mqttClient,
	}, nil
}

func parseMqttEventFromMessage(message mqtt.Message) event.Event {
	mqttEvent := event.Event{}
	err := json.Unmarshal(message.Payload(), &mqttEvent)
	if err != nil {
		log.Printf("Could not parse MQTT event: %v\n", err)
	}
	return mqttEvent
}
