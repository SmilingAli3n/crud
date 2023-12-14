package cache

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/SmilingAli3n/crud/pkg/entities"
)

type Cache struct {
	Entities map[string]*Entity // DB tables
	ttl      time.Duration
	Expires  time.Time
	sync.Mutex
}

type Entity struct {
	Records map[int64]interface{}
	sync.Mutex
}

func New(ttl time.Duration) *Cache {
	return &Cache{
		Entities: make(map[string]*Entity),
		ttl:      ttl,
		Expires:  time.Now().Add(ttl),
	}
}

func (c *Cache) Set(key string, wg *sync.WaitGroup) {
	c.Lock()
	c.Entities[key] = &Entity{
		Records: make(map[int64]interface{}),
	}
	defer wg.Done()
	c.Unlock()
}

func (e *Entity) SetRecord(id int64, value interface{}, wg *sync.WaitGroup) {
	e.Lock()
	e.Records[id] = value
	defer wg.Done()
	e.Unlock()
}

func (c *Cache) Get(key string) (map[int64]interface{}, error) {
	if c.Expired() {
		c.Renew()
	}

	c.Lock()
	defer c.Unlock()
	val, ok := c.Entities[key]
	if !ok {
		return nil, fmt.Errorf("No %s entity found in cache", key)
	} else {
		return val.Records, nil
	}
}

/*
func (e *Entity) GetRecord(id int64) (Record, error) {
	e.Lock()
	defer e.Unlock()
	val, ok := e.Record[id]
	if !ok {
		return nil, errors.New("No such key")
	}
	return val, nil
}
*/

func (c *Cache) Delete(key string) {
	c.Lock()
	delete(c.Entities, key)
	c.Unlock()
}

func (e *Entity) DeleteRecord(id int64) {
	e.Lock()
	delete(e.Records, id)
	e.Unlock()
}

func (c *Cache) Expired() bool {
	if time.Now().Before(c.Expires) {
		return false
	} else {
		return true
	}
}

func (c *Cache) Init() {
	c.Renew()
	/*for {
		if (c.Expired()) {
			break
		}
	}
	go func() {
		c.Renew()
		time.Sleep(2 * time.Second)
	}()*/
}

func (c *Cache) Renew() {
	tickets, err := entities.GetAllTickets()
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	c.Set("tickets", wg)
	wg.Wait()

	for _, val := range tickets {
		ticket := val
		wg.Add(1)
		go c.Entities["tickets"].SetRecord(ticket.Id, ticket, wg)
	}
	defer wg.Wait()

	return
}
