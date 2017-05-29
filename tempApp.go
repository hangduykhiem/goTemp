package main

//curl -X POST -d '{"index":1,"temp":"23","hum":"43"}' http://localhost:3333/value

import (
	"net/http"
	"github.com/pressly/chi"
	"database/sql"
	"log"
	_ "github.com/lib/pq"

	"github.com/pressly/chi/render"
	"github.com/pressly/chi/middleware"
	"encoding/json"
	"strconv"
)

var db *sql.DB

func main() {

	// Init
	var err error
	db, err = sql.Open("postgres", "postgresql://test:test@localhost:5432/testdb")
	if err != nil {
		log.Fatal(err)
	}

	// Init router
	r := chi.NewRouter()
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Logger)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		if err == nil {
			w.Write([]byte("Database initialized successfully"))
		} else {
			w.Write([]byte("Failed initializing database"))
		}

	})

	r.Mount("/value", ValueController())

	http.ListenAndServe(":3333", r)
}

func ValueController() chi.Router {

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		r.Post("/", CreateValue)
		r.Get("/", ListValue)
		r.Get("/:sensorId", GetValue)
	})

	return r
}

type Value struct {
	Index int `json:"index"`
	Temp  string `json:"temp"`
	Hum   string `json:"hum"`
}

func (v *Value) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewValueListResponse(articles []*Value) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, article)
	}
	return list
}

func CreateValue(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := &Value{}
	err := decoder.Decode(data)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS value(SensorId INT, Temp VARCHAR(50), Hum VARCHAR(50))")

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500: Service unavailable"))
		log.Panic(err)
		return
	}
	var sensorId int;

	err = db.QueryRow(`INSERT INTO value
	VALUES($1, $2, $3) RETURNING SensorId`, data.Index, data.Temp, data.Hum).Scan(&sensorId)

	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500: Service unavailable"))
		log.Panic(err)
		return
	} else {
		w.WriteHeader(200)
		w.Write([]byte("200: Accept"))
	}
}

func GetValue(w http.ResponseWriter, r *http.Request){
	articleID := chi.URLParam(r, "sensorId")
	i, err := strconv.Atoi(articleID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404: Not found"))
		return
	}
	rows, err := db.Query("SELECT * FROM value WHERE sensorid=$1", i)
	valueList := make([]*Value, 0)

	var index int
	var hum string
	var temp string

	if rows != nil {

		for rows.Next() {
			err := rows.Scan(&index, &temp, &hum)

			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("500: Service unavailable"))
				return
			}

			valueList = append(valueList, &Value{index, temp, hum})
		}

		err = rows.Err()

		if err != nil {
			w.Write([]byte("500: Service unavailable"))
		}

	} else {
		w.WriteHeader(404)
		w.Write([]byte("404: Not Found"))
	}

	if err := render.RenderList(w, r, NewValueListResponse(valueList)); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500: Service unavailable"))
	}
}

func ListValue(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM value")

	valueList := make([]*Value, 0)
	var index int
	var hum string
	var temp string

	if rows != nil {

		for rows.Next() {
			err := rows.Scan(&index, &temp, &hum)

			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("500: Service unavailable"))
				return
			}

			valueList = append(valueList, &Value{index, temp, hum})
		}

		err = rows.Err()

		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte("500: Service unavailable"))
		}

	} else {
		w.Write([]byte("Nothing found"))
	}

	if err := render.RenderList(w, r, NewValueListResponse(valueList)); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500: Service unavailable"))
	}
}