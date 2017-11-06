package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/lib/pq"
)

// Server - The server object
type Server struct {
	Router *chi.Mux
	DB     *sql.DB
}

// Init - Initializes the Server
func (s *Server) Init(user, password, dbname string) {
	connection := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	s.DB, err = sql.Open("postgres", connection)
	checkError(err)

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.InitRoutes()
}

// InitRoutes - Initializes the Server's route
func (s *Server) InitRoutes() {
	s.Router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("dcu"))
	})

	s.Router.Route("/factions", func(r chi.Router) {
		r.Get("/", s.ListFactions)
		r.Route("/{factionId}", func(r chi.Router) {
			r.Get("/", GetFaction)
		})
	})

	s.Router.Route("/characters", func(r chi.Router) {
		r.Get("/", ListCharacters)
		r.Route("/{characterId}", func(r chi.Router) {
			r.Get("/", GetCharacter)
		})
	})
}

// Run - Start the Server
func (s *Server) Run() {
	log.Fatal(http.ListenAndServe(":8080", s.Router))
}

// ListFactions - lists available factions
func (s *Server) ListFactions(w http.ResponseWriter, r *http.Request) {
	factions, err := getFactions(s.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, factions)
}

// GetFaction - all characters that belong to a specific faction id
func GetFaction(w http.ResponseWriter, r *http.Request) {

}

// ListCharacters - list all characters
func ListCharacters(w http.ResponseWriter, r *http.Request) {
	//if err := render.RenderList()
}

// GetCharacter - get character details given a character id
func GetCharacter(w http.ResponseWriter, r *http.Request) {

}
