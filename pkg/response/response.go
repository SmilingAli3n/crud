package response

import (
	"encoding/json"
	"log"
	"net/http"
	// "fmt"
	//"github.com/SmilingAli3n/crud/entities"
)

type Response struct {
	http.Response
	//Result string `json:"result"`
	//Code int
	Entities []interface{} `json:"entities"`
}

func New() *Response {
	return &Response{}
}

func (r *Response) Send(w http.ResponseWriter) {
	if r.StatusCode == 0 { // something's wrong
		r.StatusCode = http.StatusInternalServerError
		w.WriteHeader(r.StatusCode)
		log.Print("Status code was not set")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(tk)
	if len(r.Entities) == 0 {
		w.WriteHeader(r.StatusCode)
	} else {
		respBin, err := json.Marshal(r.Entities)
		if err != nil {
			r.StatusCode = http.StatusInternalServerError // 500
			w.WriteHeader(r.StatusCode)
			log.Print(err)
		}
		_, err = w.Write(respBin)
		if err != nil {
			r.StatusCode = http.StatusInternalServerError // 500
			w.WriteHeader(r.StatusCode)
			log.Print(err)
		}
	}
}
