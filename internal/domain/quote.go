package domain

import "fmt"

type Quote struct {
	Text   string `json:"text"`
	Author string `json:"author"`
}

func (q Quote) String() string {
	return fmt.Sprintf("text: %s, author: %s", q.Text, q.Author)
}

func (q Quote) Bytes() []byte {
	return []byte(q.String())
}
