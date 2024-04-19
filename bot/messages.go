package bot

type EventPayload struct {
	Op int    `json:"op,omitempty"`
	S  int    `json:"s,omitempty"`
	T  string `json:"t,omitempty"`
}

type HelloMessage struct {
	EventPayload
	D struct {
		HeartbeatInterval int `json:"heartbeat_interval,omitempty"`
	} `json:"d,omitempty"`
}

type HeartbeatMessage struct {
	EventPayload
	D int `json:"d,omitempty"`
}
