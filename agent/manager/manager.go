package manager

import (
	"github.com/gustavoppereira/llv/agent/configuration"
	"github.com/gustavoppereira/llv/agent/event"
	"github.com/gustavoppereira/llv/agent/http"
	"github.com/gustavoppereira/llv/agent/mqtt"
	"github.com/gustavoppereira/llv/agent/watcher"
	"log"
	"time"
)

const WatcherTickerDuration = 5 * time.Second

type Manager struct {
	mqttClient         *mqtt.Client
	agentConfiguration *configuration.AgentConfiguration
	watcher            *watcher.FileWatcher

	watcherTicker *time.Ticker
}

func NewManager() *Manager {
	return &Manager{}
}

func (m *Manager) StartTicker() {
	ticker := time.NewTicker(WatcherTickerDuration)
	go func() {
		<-ticker.C
		// Ticker ended without ping event. Stop and cleaning up watcher resources
		m.watcher.Cleanup()
		// Set back agent state to Enabled
		m.UpdateAgentState(configuration.Enabled)
	}()
}

func (m *Manager) TickWatcher() {
	m.watcherTicker.Reset(WatcherTickerDuration)
}

func (m *Manager) SetFileWatcher(fileWatcher *watcher.FileWatcher) {
	m.watcher = fileWatcher
}

func (m *Manager) UpdateAgentState(state configuration.State) {
	m.agentConfiguration.State = state
}

func (m *Manager) getGlobalMessageHandler() mqtt.GlobalMessageHandler {
	return func(topic string, event event.Event) {
		handler, err := GetEventHandler(event)
		if err != nil {
			log.Println(err)
			return
		}
		handler.OnEventReceived(*m, event)
	}
}

func (m *Manager) StartAgent(token string) {
	// 1 - Read from local config to see if it needs to create a new one or use an already created
	m.initializeConfiguration(token)

	m.connectToRemoteBroker()
}

func (m *Manager) initializeConfiguration(token string) {
	var hasLocalConfiguration = configuration.HasConfiguration()
	if !hasLocalConfiguration {
		// Create instance on remote server and retrieve broker connection info
		createInstanceResponse := m.callServerToCreateInstance(token)
		brokerConfiguration := configuration.NewBrokerConfiguration(createInstanceResponse.BrokerPort,
			createInstanceResponse.BrokerPassword)

		m.agentConfiguration = configuration.NewEnabled(token, brokerConfiguration)
		configuration.SaveConfiguration(m.agentConfiguration)
	} else {
		var err error
		m.agentConfiguration, err = configuration.GetConfiguration()
		if err != nil {
			log.Fatalf("Error loading stored configuration: %v\n", err)
		}
	}
}

func (m *Manager) connectToRemoteBroker() *mqtt.Client {
	if m.agentConfiguration == nil {
		log.Fatalf("Agent configuration initialization required before connecting to broker.\n")
	}

	options := mqtt.NewConnectionOptions(configuration.GetServerHost(), m.agentConfiguration.Broker.Port,
		"ClientID", configuration.GetInstanceName(), m.agentConfiguration.Broker.Password)

	client, err := mqtt.Connect(options, m.getGlobalMessageHandler())
	if err != nil {
		log.Fatalf("Error connecting to server's broker: %v\n", err)
	}
	return client
}

func (m *Manager) callServerToCreateInstance(token string) *http.CreateInstanceResponse {
	response, err := http.CallCreateInstance(token, configuration.GetInstanceName())
	if err != nil {
		log.Fatalf("Could not connect with server to fetch configuration and initialize instance: %v", err)
	}
	return response
}
