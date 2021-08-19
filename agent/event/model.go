package event

type Type string

const (
	UpdateAgent Type = Type("UPDATE_AGENT")
	PingLogging Type = Type("PING_LOGGING")
)

type Event struct {
	Type        Type
	Payload     []byte
}
