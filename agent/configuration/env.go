package configuration

import "os"

type EnvironmentVariable string

const (
	InstanceName EnvironmentVariable = "LLV_INSTANCE_NAME"
	ServerHost   EnvironmentVariable = "LLV_SERVER_HOST"
	LogFilePath  EnvironmentVariable = "LLV_LOG_FILE_PATH"
)

func GetInstanceName() string {
	return os.Getenv(string(InstanceName))
}

func GetServerHost() string {
	return os.Getenv(string(ServerHost))
}

func GetLogFilePath() string {
	return os.Getenv(string(LogFilePath))
}
