package cmd

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/SmilingAli3n/crud/pkg/auth"
	"github.com/SmilingAli3n/crud/pkg/cache"
	"github.com/SmilingAli3n/crud/pkg/repos"
	"github.com/SmilingAli3n/crud/pkg/response"
)

var c = cache.New(time.Minute)
var once sync.Once

func init() {
	go once.Do(func() {
		c.Init()
	})
}

func RunServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ticket/", ticketHandler)
	mux.HandleFunc("/tickets", allTicketsHandler)
	http.ListenAndServe(":8080", mux)
	fmt.Println("Server started")
}

func allTicketsHandler(w http.ResponseWriter, req *http.Request) {
	if !auth.Authorized(req) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	if req.URL.Path != "/tickets" {
		log.Print("Request URL must not contain any params")
		return
	}
	resp := response.New()
	defer resp.Send(w)
	if req.Method == http.MethodPost {
		repos.CreateTicket(req, resp)
	} else if req.Method == http.MethodGet {
		repos.GetAllTickets(req, resp, c)
	} else {
		resp.Status = fmt.Sprintf("Method %v is not supported", req.Method)
		resp.StatusCode = http.StatusMethodNotAllowed
		return
	}
}
func ticketHandler(w http.ResponseWriter, req *http.Request) {
	resp := response.New()
	defer resp.Send(w)
	if !auth.Authorized(req) {
		resp.StatusCode = http.StatusUnauthorized
		return
	}
	url := strings.Trim(req.URL.Path, "/")
	urlParts := strings.Split(url, "/")
	if len(urlParts) < 2 {
		resp.Status = "expect /ticket/<id> in task handler"
		resp.StatusCode = http.StatusBadRequest
		return
	}

	id, err := strconv.Atoi(urlParts[1])
	if err != nil {
		resp.Status = err.Error()
		resp.StatusCode = http.StatusBadRequest
		return
	}
	if req.Method == http.MethodPut {
		repos.UpdateTicket(req, resp, int64(id))
	} else if req.Method == http.MethodDelete {
		repos.DeleteTicketById(req, resp, int64(id))
	} else if req.Method == http.MethodGet {
		repos.GetTicketById(req, resp, c, int64(id))
	} else {
		resp.Status = fmt.Sprintf("expect method GET or DELETE at /ticket/<id>, got %v", req.Method)
		resp.StatusCode = http.StatusMethodNotAllowed
		return
	}

}
