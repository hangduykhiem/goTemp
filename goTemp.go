package main

//curl -X POST -d '{"index":1,"temp":"23","hum":"43"}' http://localhost:3333/value

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/pressly/chi/render"
	"goTemp/dbRepo"
	"goTemp/model"
	"log"
	"net/http"
	"strconv"
)

type App struct {
	db dbRepo.TempDb
}

func main() {

	// Initdb
	app, err := initApp()

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

	r.Mount("/value", ValueController(app))

	http.ListenAndServe(":3333", r)
}

func initApp() (*App,error) {

	postgres, err := sql.Open("postgres", "postgresql://test:test@localhost:5432/testdb")
	db := &dbRepo.PostgesTempDb{postgres}
	if err != nil {
		return nil, err;
	}

	return &App{db:db}, nil
}

func ValueController(app *App) chi.Router {

	r := chi.NewRouter()

	r.Route("/", func(r chi.Router) {

		r.Post("/", app.CreateValueRoute)
		r.Get("/", app.ListValuesRoute)
		r.Get("/:sensorId", app.GetValueRoute)
	})

	return r
}

func NewValueListResponse(articles []*model.Value) []render.Renderer {
	list := []render.Renderer{}
	for _, article := range articles {
		list = append(list, article)
	}
	return list
}

func (app *App) CreateValueRoute(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	data := &model.Value{}
	err := decoder.Decode(data)
	if err != nil {
		panic(err)
	}

	err = app.db.CreateValue(data)

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

func (app *App) GetValueRoute(w http.ResponseWriter, r *http.Request) {
	articleID := chi.URLParam(r, "sensorId")
	i, err := strconv.Atoi(articleID)
	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404: Not found"))
		return
	}

	valueList, err := app.db.GetValue(i)

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404: Not found"))
	}

	if err := render.RenderList(w, r, NewValueListResponse(valueList)); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500: Service unavailable"))
	}
}

func (app *App) ListValuesRoute(w http.ResponseWriter, r *http.Request) {

	valueList, err := app.db.ListValues()

	if err != nil {
		w.WriteHeader(404)
		w.Write([]byte("404: Not found"))
	}

	if err := render.RenderList(w, r, NewValueListResponse(valueList)); err != nil {
		w.WriteHeader(500)
		w.Write([]byte("500: Service unavailable"))
	}
}
