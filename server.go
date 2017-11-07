package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
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

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})

	s.Router = chi.NewRouter()
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(cors.Handler)
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
			r.Get("/", s.GetFaction)
		})
	})

	s.Router.Route("/characters", func(r chi.Router) {
		r.Get("/", s.ListCharacters)
		r.Route("/{characterId}", func(r chi.Router) {
			r.Get("/", s.GetCharacter)
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
func (s *Server) GetFaction(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "factionId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid faction id")
	}
	f := faction{ID: id}
	err = f.getFaction(s.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, f)
}

// ListCharacters - list all characters
func (s *Server) ListCharacters(w http.ResponseWriter, r *http.Request) {
	characters, err := getCharacters(s.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, characters)
}

// GetCharacter - get character details given a character id
func (s *Server) GetCharacter(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "characterId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Invalid character id")
	}

	c := character{ID: id}
	if err := c.getCharacter(s.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Character not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, c)
}
