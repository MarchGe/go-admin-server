package mq

// Message demo message
type Message struct {
	Id   string `json:"id"`
	Data any    `json:"data"`
}
