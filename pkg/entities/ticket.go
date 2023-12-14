package entities

import (
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	// "database/sql"
	"github.com/SmilingAli3n/crud/pkg/db"

	_ "github.com/go-sql-driver/mysql"
)

type Ticket struct {
	//Opened   bool   `json:"opened"`
	Priority int    `json:"priority"`
	Id       int64  `json:"id"`
	Text     string `json:"text"`
	//UserId sql.NullInt64
	//ClientId sql.NullInt64
	Time string 	`json:"time"`
	//Unread bool `json:"unread"`
	//CreatorId sql.NullInt64
}

var (
	ErrNotFound = errors.New("Ticket not found")
	ErrInternal = errors.New("Internal server error")
)

func GetAllTickets() ([]Ticket, error) {
	db, err := db.GetInstance()
	if err != nil {
		log.Print(err.Error())
		return nil, ErrInternal
	}
	defer db.Close()
	rows, err := db.Query("SELECT id, text, priority, time FROM tickets ORDER BY id")
	if err != nil {
		log.Print(err.Error())
		return nil, ErrInternal
	}
	tckts := make([]Ticket, 0)
	for rows.Next() {
		var t Ticket
		if err := rows.Scan(&t.Id, &t.Text, &t.Priority, &t.Time); err != nil {
			log.Print(err.Error())
			return nil, ErrInternal
		}
		tckts = append(tckts, t)
	}
	return tckts, nil
}

func (t *Ticket) Create() error {
	db, err := db.GetInstance()
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	defer db.Close()
	//fmt.Printf("Time is %s\n", t.Time)
	res, err := db.Exec("INSERT INTO tickets (text, priority, time) VALUES (?, ?, ?)", t.Text, t.Priority, t.Time) //, t.UserId, t.ClientId, t.Time, t.Unread, t.CreatorId)
	if err != nil {
		fmt.Println(err.Error())
		return ErrInternal
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	if id == 0 {
		log.Print(err.Error())
		return ErrInternal
	}
	return nil
}

func DeleteTicket(id int64) error {
	db, err := db.GetInstance()
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	defer db.Close()
	res, err := db.Exec("DELETE FROM tickets WHERE id = ?", id)
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	if rows == 0 {
		log.Print(err.Error())
		return ErrInternal
	}
	return nil
}

func UpdateTicket(id int64, m map[string]interface{}) error {
	db, err := db.GetInstance()
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	defer db.Close()
	//
	rows, err := db.Query("SELECT * FROM tickets WHERE id=?", id)
	if err != nil {
		log.Print(err.Error())
		return ErrInternal
	}
	if !rows.Next() {
		return ErrNotFound
	}
	keys := make([]string, 0)
	for key := range m {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	vals := make([]interface{}, 0)
	for _, key := range keys {
		vals = append(vals, m[string(key)])
	}
	var sb strings.Builder
	sb.WriteString("UPDATE tickets SET")
	for _, val := range keys {
		sb.WriteString(fmt.Sprintf(" %s=?,", val))
	}
	qStr := strings.TrimRight(sb.String(), ",")
	qStr += " WHERE id=?"
	vals = append(vals, id)
	_, err = db.Exec(qStr, vals...)
	if err != nil {
		fmt.Println(err.Error())
		return ErrInternal
	}
	/*
		rows, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rows == 0 {
			return errors.New("Could not update ticket")
		}
	*/
	return nil
}
