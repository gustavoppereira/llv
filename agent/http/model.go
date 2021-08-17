package http

import (
	"encoding/json"
	"log"
)

type Payload interface {
	Content() []byte
}

type CreateInstanceRequest struct {
	instanceName string
}

type CreateInstanceResponse struct {
	BrokerPort     int
	BrokerPassword string
}

func (c *CreateInstanceRequest) Content() []byte {
	content, err := json.Marshal(c)
	if err != nil {
		log.Printf("Could not parse CreateInstanceRequest from: %v\n", c)
		panic(err)
	}
	return content
}

func newCreateInstanceRequest(instanceName string) *CreateInstanceRequest {
	return &CreateInstanceRequest{
		instanceName: instanceName,
	}
}

func parseCreateInstanceResponse(content []byte) (*CreateInstanceResponse, error) {
	response := CreateInstanceResponse{}
	var err = json.Unmarshal(content, &response)

	return &response, err
}
