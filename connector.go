package connector

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const dateformat = "2006-01-02 15:04:05"

var (
	db *sql.DB //database connection
)

//prepared database queries
var selectAll, selectFilter, update, _delete, insert *sql.Stmt

func Open(driverName, dataSourceName, url string) error {
	var err error
	db, err = sql.Open(driverName, dataSourceName) //database initialization
	if err != nil {
		return err
	}

	selectAll, err = db.Prepare(
		`SELECT id, start_date, end_date, text
		 FROM events
	`)
	if err != nil {
		return err
	}
	selectFilter, err = db.Prepare(
		`SELECT id, start_date, end_date, text
		 FROM events
		 WHERE start_date > ? AND start_date < ?
		 OR end_date > ? AND end_date < ?
	`)
	if err != nil {
		return err
	}
	update, err = db.Prepare(
		`UPDATE events
		 SET start_date = ?, end_date = ?, text = ?
		 WHERE id = ?
	`)
	if err != nil {
		return err
	}
	insert, err = db.Prepare(
		`INSERT INTO events (start_date, end_date, text)
		 VALUES (?, ?, ?)
	`)
	if err != nil {
		return err
	}
	_delete, err = db.Prepare(
		`DELETE
		 FROM events
		 WHERE id = ?
	`)
	if err != nil {
		return err
	}
//, _delete, insert

	if len(url) == 0 {
		url = "/connector"
	}
	http.HandleFunc(url, dhtmlxHandler)
	return nil
}

type ActionResult struct {
	Data []Action `json:"data"`
}

type Action struct {
	Type string `json:"type"`
	Sid  string `json:"sid"`
	Tid  string `json:"tid"`
}

func dhtmlxHandler(res http.ResponseWriter, req *http.Request) {
	ids := req.FormValue("ids")
	var err error
	if len(ids) == 0 { //SELECT
		from := req.FormValue("from")
		if len(from) == 0 {
			rows, err := selectAll.Query()
			selectResult(res, rows, err)
		} else {
			to := req.FormValue("to")
			rows, err := selectFilter.Query(from, to, from, to)
			selectResult(res, rows, err)
		}
	} else {
		var b bytes.Buffer
		getField := func(id, postfix string) string {
			b.Reset()
			b.WriteString(id)
			b.WriteString(postfix)
			return req.FormValue(b.String())
		}
		actions := make([]Action, 0)
		var action Action
		for _, id := range(strings.Split(ids, ",")) {
			status := getField(id, "_!nativeeditor_status")
			oldId := getField(id, "_id")
			newId := oldId
			if status == "deleted" {
				_, err = _delete.Exec(id)
			} else {
				start_date := getField(id, "_start_date")
				end_date := getField(id, "_end_date")
				text := getField(id, "_text")
				if status == "inserted" {
					r, err := insert.Exec(start_date, end_date, text)
					if err != nil {
						fmt.Println("insert error: ", err)
						return
					}
					newId64, err := r.LastInsertId()
					newId = strconv.Itoa(int(newId64))
				} else if status == "updated" {
					_, err = update.Exec(start_date, end_date, text, id)
				} else {
					return
				}
			}
			if err != nil {
				fmt.Println("action error: ", err)
				return
			}
			action.Type = status
			action.Sid = oldId
			action.Tid = newId
			actions = append(actions, action)
		}
		by, err := json.Marshal(ActionResult{actions})
		if err != nil {
			fmt.Println("marshal error", err)
			return
		}
		fmt.Fprint(res, string(by))
	}
}

type SelectResult struct {
	Data []Event `json:"data"`
}

type Event struct {
	Id         int    `json:"id"`
	Start_date string `json:"start_date"`
	End_date   string `json:"end_date"`
	Text       string `json:"text"`
}

func selectResult(res http.ResponseWriter, rows *sql.Rows, err error) {
	if err != nil {
		fmt.Println("error getting rows", err)
		return
	}
	defer rows.Close()

	events := make([]Event, 0)
	var event Event
	var start_date, end_date time.Time

	for rows.Next() {
		err = rows.Scan(&event.Id, &start_date, &end_date, &event.Text)
		if err != nil {
			fmt.Println("error scanning rows", err)
			return
		}
		event.Start_date = start_date.Format(dateformat)
		event.End_date = end_date.Format(dateformat)
		events = append(events, event)
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("rows error", err)
		return
	}

	b, err := json.Marshal(SelectResult{events})
	if err != nil {
		fmt.Println("marshal error", err)
		return
	}
	fmt.Fprint(res, string(b))
}
