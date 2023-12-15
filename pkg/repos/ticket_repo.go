package repos

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/SmilingAli3n/crud/pkg/cache"
	"github.com/SmilingAli3n/crud/pkg/entities"
	"github.com/SmilingAli3n/crud/pkg/response"
)

func jsonToMap(r *http.Request) (map[string]interface{}, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func jsonToTicket(r *http.Request) (*entities.Ticket, error) {
	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}
	t := &entities.Ticket{}
	err = json.Unmarshal(b, t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

func GetTicketById(r *http.Request, resp *response.Response, c *cache.Cache, id int64) { //w http.ResponseWriter,
	tickets, err := c.Get("tickets")
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		return
	}
	if tck, ok := tickets[id]; ok {
		resp.StatusCode = http.StatusOK // 200
		resp.Status = "200 OK"
		resp.Entities = append(resp.Entities, tck.(entities.Ticket))
	} else {
		resp.StatusCode = http.StatusNotFound
		log.Print(err)
		return
	}
	/*
		200 OK if the response contains data
	    204 No Content if the response contains no data
	*/
}

func GetAllTickets(r *http.Request, resp *response.Response, c *cache.Cache) { //w http.ResponseWriter,
	tickets, err := c.Get("tickets")
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError // 500
		log.Print(err)
		return
	}
	resp.Status = "OK"
	for _, t := range tickets {
		resp.Entities = append(resp.Entities, t.(entities.Ticket))
	}
	resp.StatusCode = http.StatusOK // 200
}

func CreateTicket(r *http.Request, resp *response.Response) { //w http.ResponseWriter,
	t, err := jsonToTicket(r)
	if err != nil {
		resp.StatusCode = http.StatusInternalServerError
		log.Print(err)
	}
	if t.Time == "" {
		t.Time = "0001-01-01 00:00:00"
	}
	err = t.Create()
	if err != nil {
		if errors.Is(err, entities.ErrInternal) {
			resp.StatusCode = http.StatusInternalServerError // 500
		} else if errors.Is(err, entities.ErrNotFound) {
			resp.StatusCode = http.StatusNotFound
		}
		log.Print(err)
		return
	}
	/*
		201 Created for successful create operation
	    200 OK for successful read operation if the response contains data
	    204 No Content for successful read operation if the response contains NO data
	*/
	resp.StatusCode = http.StatusCreated // 201
}

func UpdateTicket(r *http.Request, resp *response.Response, id int64) { //w http.ResponseWriter,
	m, err := jsonToMap(r)
	if err != nil {
		resp.StatusCode = http.StatusBadRequest // ???
		log.Print(err)
		return
	}
	err = entities.UpdateTicket(id, m)
	if err != nil {
		if errors.Is(err, entities.ErrInternal) {
			resp.StatusCode = http.StatusInternalServerError // 500
		} else if errors.Is(err, entities.ErrNotFound) {
			resp.StatusCode = http.StatusNotFound
		}
		log.Print(err)
		return
	}
	/*
		200 OK if the response contains data
	    204 No Content if the response contains no data
	*/
	resp.StatusCode = http.StatusNoContent // 204
}

func DeleteTicketById(r *http.Request, resp *response.Response, id int64) {
	err := entities.DeleteTicket(id)
	if err != nil {
		if errors.Is(err, entities.ErrInternal) {
			resp.StatusCode = http.StatusInternalServerError // 500
		} else if errors.Is(err, entities.ErrNotFound) {
			resp.StatusCode = http.StatusNotFound
		}
		log.Print(err)
		return
	}
	resp.StatusCode = http.StatusNoContent // 204
}
