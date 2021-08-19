package manager

import (
	"errors"
	"fmt"
	"github.com/gustavoppereira/llv/agent/configuration"
	"github.com/gustavoppereira/llv/agent/event"
	fileWatcher "github.com/gustavoppereira/llv/agent/watcher"
	"log"
)

type EventHandler interface {
	OnEventReceived(manager Manager, event event.Event)
}

type UpdateAgentEvent struct {
}

func (u UpdateAgentEvent) OnEventReceived(manager Manager, event event.Event) {
}

type PingLoggingEvent struct {
}

func (p PingLoggingEvent) OnEventReceived(manager Manager, event event.Event) {
	agentState := manager.agentConfiguration.State
	switch agentState {
	case configuration.Enabled:
		watcher := fileWatcher.NewFileWatcher(configuration.GetLogFilePath(), func(logLine string) {
			manager.mqttClient.Publish("someTopic", logLine)
		})
		err := watcher.Watch()
		if err != nil {
			log.Fatalf("Error watching file watcher: %v\n", err)
		}
		manager.SetFileWatcher(watcher)
		manager.StartTicker()
	case configuration.Logging:
		// Keep ticking watcher to maintain tail channel
		manager.TickWatcher()
	}
}

func GetEventHandler(e event.Event) (EventHandler, error) {
	switch t := e.Type; t {
	case event.UpdateAgent:
		return UpdateAgentEvent{}, nil
	case event.PingLogging:
		return PingLoggingEvent{}, nil
	default:
		return nil, errors.New(fmt.Sprintf("No handler defined for event type %v", t))
	}
}
