package client

// Letter is a letter that is sent to deadletter
type Letter struct {
	Method  string              `json:"method"`
	URL     string              `json:"url"`
	Body    []byte              `json:"body"`
	Headers map[string][]string `json:"headers"`
}

// DeadLetter save request to somewhere to ensure consistency
type DeadLetter interface {
	Save(letter *Letter) error
}
