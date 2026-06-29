package message

type Entity struct {
	Source  string     `json:"source"`
	Key     string     `json:"key"`
	ToUsers []string   `json:"to_users"`
	MsgTime int64      `json:"msg_time"`
	Msgs    []*Message `json:"msgs"`
	Extra   string     `json:"extra"`
	MsgId   int64      `json:"msg_id"`
}

type Message struct {
	Terminals []int64 `json:"terminals"`
	Category  string  `json:"category"`
	Template  string  `json:"template"`
	Body      string  `json:"body"`
	Version   int64   `json:"version"`
	Ext       string  `json:"ext"`
	Nopop     bool    `json:"nopop"`
}
