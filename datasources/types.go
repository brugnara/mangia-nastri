package datasources

import "net/http"

type Hash string

func (h Hash) String() string {
	return string(h)
}

type Request struct {
	Method string      `json:"method"`
	URL    string      `json:"url"`
	Header http.Header `json:"headers" gorm:"serializer:json"`
	Body   string      `json:"body"`
}

type Response struct {
	Header http.Header `json:"headers" gorm:"serializer:json"`
	Body   string      `json:"body"`
	Status int         `json:"status"`
}

type Payload struct {
	CreationTime int64    `json:"creationTime"`
	Request      Request  `json:"request" gorm:"serializer:json"`
	Response     Response `json:"response" gorm:"serializer:json"`
}
