package model

import "net/http"

type Value struct {
	Index int    `json:"index"`
	Temp  string `json:"temp"`
	Hum   string `json:"hum"`
}

func (v *Value) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}
