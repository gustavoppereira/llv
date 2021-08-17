package manager

import (
	"github.com/gustavoppereira/llv/agent/configuration"
	"github.com/gustavoppereira/llv/agent/event"
	"github.com/gustavoppereira/llv/agent/http"
	"github.com/gustavoppereira/llv/agent/mqtt"
	"github.com/gustavoppereira/llv/agent/watcher"
	"log"
)

type Manager struct {
	mqttClient         *mqtt.Client
	agentConfiguration *configuration.AgentConfiguration
	watcher            *watcher.FileWatcher
}

func (m *Manager) SetFileWatcher(fileWatcher *watcher.FileWatcher) {
	m.watcher = fileWatcher
}

func (m *Manager) UpdateAgentState(state configuration.State) {
	m.agentConfiguration.State = state
}

var globalMessageHandler mqtt.GlobalMessageHandler = func(topic string, event event.Event) {
	log.Printf("Received message %v on topic %v\n", event, topic)
}

func StartAgent(token string) *Manager {
	// 1 - Read from local config to see if need to create a new one or use already created
	var hasLocalConfiguration = configuration.HasConfiguration()
	var config *configuration.AgentConfiguration
	if !hasLocalConfiguration {
		config = createConfiguration(token)
	} else {
		var err error
		config, err = configuration.GetConfiguration()
		if err != nil {
			log.Fatalf("Error loading stored configuration: %v\n", err)
		}
	}
	mqttClient := connectToRemoteBroker(config)

	return &Manager{
		mqttClient:         mqttClient,
		agentConfiguration: config,
	}
}

func connectToRemoteBroker(agentConfiguration *configuration.AgentConfiguration) *mqtt.Client {
	options := mqtt.NewConnectionOptions(configuration.GetServerHost(), agentConfiguration.Broker.Port,
		"ClientID", configuration.GetInstanceName(), agentConfiguration.Broker.Password)

	client, err := mqtt.Connect(options, globalMessageHandler)
	if err != nil {
		log.Fatalf("Error connecting to server's broker: %v\n", err)
	}
	return client
}

func createConfiguration(token string) *configuration.AgentConfiguration {
	createInstanceResponse := callServerToCreateInstance(token)
	brokerConfiguration := configuration.NewBrokerConfiguration(createInstanceResponse.BrokerPort,
		createInstanceResponse.BrokerPassword)

	config := configuration.NewEnabled(token, brokerConfiguration)
	configuration.SaveConfiguration(config)

	return config
}

func callServerToCreateInstance(token string) *http.CreateInstanceResponse {
	response, err := http.CallCreateInstance(token, configuration.GetInstanceName())
	if err != nil {
		log.Fatalf("Could not connect with server to fetch configuration and initialize instance: %v", err)
	}
	return response
}
