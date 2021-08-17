package http

import (
	"bytes"
	"encoding/json"
	"github.com/gustavoppereira/llv/agent/configuration"
	"io"
	"net/http"
)

const (
	CreateInstance = "/instance"
)

//goland:noinspection GoUnhandledErrorResult
func post(path string, payload Payload, headers map[string]string) ([]byte, error) {
	var req, err = http.NewRequest("POST", getFullPath(path), bytes.NewBuffer(payload.Content()))
	if err != nil {
		return nil, err
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	return body, err
}

func getFullPath(path string) string {
	return configuration.GetServerHost() + path
}

func CallCreateInstance(token string, instanceName string) (*CreateInstanceResponse, error) {
	var payload = newCreateInstanceRequest(instanceName)

	resp, err := post(CreateInstance, payload, map[string]string{
		"Token": token,
	})
	if err != nil {
		return nil, err
	}
	return parseCreateInstanceResponse(resp)
}

func marshal(value interface{}) ([]byte, error) {
	return json.Marshal(value)
}
