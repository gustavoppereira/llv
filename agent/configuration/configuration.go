package configuration

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"os"
)

const (
	Enabled State = "ENABLED"
	Logging State = "LOGGING"

	FilePath = "/usr/var/llv/configuration.json"
)

type State string

type AgentConfiguration struct {
	Token  string
	Id     string
	State  State
	Broker BrokerConfiguration
}

type BrokerConfiguration struct {
	Port     int
	Password string
}

func NewBrokerConfiguration(port int, password string) BrokerConfiguration {
	return BrokerConfiguration{
		Port:     port,
		Password: password,
	}
}

type AgentConfigurationError struct {
	Err error
}

func (a AgentConfigurationError) Error() string {
	panic(a.Err)
}

func NewEnabled(token string, configuration BrokerConfiguration) *AgentConfiguration {
	return newConfiguration(Enabled, token, configuration)
}

func NewLogging(token string, configuration BrokerConfiguration) *AgentConfiguration {
	return newConfiguration(Logging, token, configuration)
}

func newConfiguration(state State, token string, configuration BrokerConfiguration) *AgentConfiguration {

	var id, _ = uuid.NewUUID()

	return &AgentConfiguration{
		Token:  token,
		Id:     id.String(),
		State:  state,
		Broker: configuration,
	}
}

func GetConfiguration() (*AgentConfiguration, error) {
	content, err := readFile(FilePath)

	if e, ok := err.(*os.PathError); ok && os.IsNotExist(e.Err) {
		log.Println("Could not find configuration file. Returning error...")
		return nil, AgentConfigurationError{
			Err: e,
		}
	}

	var agentConfiguration = AgentConfiguration{}
	err = json.Unmarshal(content, &agentConfiguration)
	if err != nil {
		panic(err)
	}
	return &agentConfiguration, nil
}

func HasConfiguration() bool {
	_, err := readFile(FilePath)

	e, ok := err.(*os.PathError)
	return err == nil || (ok && !os.IsNotExist(e.Err))
}

func SaveConfiguration(configuration *AgentConfiguration) {
	bytes, err := json.Marshal(configuration)
	if err != nil {
		panic(err)
	}
	err = saveFile(bytes, FilePath)
	if err != nil {
		log.Fatal(err)
	}
}
