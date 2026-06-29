package event

type Entity struct {
	Topic         string `json:"topic"`
	Operation     string `json:"operation"`
	Time          int64  `json:"time"`
	Nonce         string `json:"nonce"`
	Signature     string `json:"signature"`
	Data          string `json:"data"`
}
