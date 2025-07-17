package domain

type Message struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}
