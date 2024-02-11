package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpMsg struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func (m *HttpMsg) toBytes() []byte {
	data, err := json.Marshal(m)

	if err != nil {
		log.Println(err)
	}

	return data
}

func (m *HttpMsg) Write(w http.ResponseWriter) {
	w.WriteHeader(m.Code)
	w.Write(m.toBytes())
}

var ErroHttpMsgPageNotFound HttpMsg = HttpMsg{
	Msg:  "Erro Page Not Found",
	Code: http.StatusNotFound,
}

var ErroHttpMsgMethodNotAllowed HttpMsg = HttpMsg{
	Msg:  "Erro Method Not Allowed",
	Code: http.StatusMethodNotAllowed,
}
