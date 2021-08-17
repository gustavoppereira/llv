package event

type Type string

const (
	UpdateAgent = Type("UPDATE_AGENT")
	PingLogging = Type("PING_LOGGING")
)

type Event struct {
	Type    Type
	Payload []byte
}
