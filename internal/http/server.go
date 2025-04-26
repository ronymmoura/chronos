package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/ronymmoura/chronos/internal/db"
	"github.com/ronymmoura/chronos/internal/db/repositories"
	"github.com/ronymmoura/chronos/internal/models"
	"github.com/ronymmoura/chronos/internal/util"

	"github.com/gorilla/mux"
)

type Server struct {
	DB      *db.DB
	Config  *util.Config
	Router  *mux.Router
	Channel chan bool
}

func NewServer(db *db.DB, config *util.Config, router *mux.Router, c chan bool) *Server {
	return &Server{
		DB:      db,
		Config:  config,
		Router:  router,
		Channel: c,
	}
}

func (s *Server) RegisterRoutes() {
	s.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello there!")
		s.Channel <- true
	})

	s.Router.HandleFunc("/process", s.createProcess).Methods(http.MethodPost)
	s.Router.HandleFunc("/process", s.getAllProcesses).Methods(http.MethodGet)
	s.Router.HandleFunc("/process/{id}", s.getProcess).Methods(http.MethodGet)
}

func (s *Server) createProcess(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			http.Error(w, "No file submitted", http.StatusBadRequest)
		} else {
			http.Error(w, "Failed to parse file", http.StatusBadRequest)
		}
		return
	}
	defer file.Close()

	folder := uuid.New().String()
	path := filepath.Join("uploads", folder)

	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		http.Error(w, "Failed to create directory", http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(path, header.Filename)

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err = io.Copy(f, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	envPath := filepath.Join(path, ".env")
	err = os.WriteFile(envPath, []byte(r.FormValue("env")), 0755)
	if err != nil {
		http.Error(w, "Failed to create .env file", http.StatusInternalServerError)
		return
	}

	executeEverySecs, _ := strconv.Atoi(r.FormValue("execute_every_secs"))

	process := &models.Process{
		Name:             r.FormValue("name"),
		Description:      r.FormValue("description"),
		Env:              r.FormValue("env"),
		Path:             filePath,
		ExecuteEverySecs: executeEverySecs,
	}

	repo := repositories.NewProcessRepository(s.DB)
	newProcess, err := repo.Insert(process)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create process: %v", err), http.StatusInternalServerError)
		return
	}

	s.Channel <- true

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProcess)
}

func (s *Server) getProcess(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid process ID", http.StatusBadRequest)
		return
	}

	// Fetch the process from the database
	repo := repositories.NewProcessRepository(s.DB)
	process, err := repo.GetByID(int32(id))
	if err != nil {
		http.Error(w, "Process not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(process)
}

func (s *Server) getAllProcesses(w http.ResponseWriter, r *http.Request) {
	// Fetch all processes from the database
	repo := repositories.NewProcessRepository(s.DB)
	processes, err := repo.SelectAll()
	if err != nil {
		http.Error(w, "Failed to fetch processes", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(processes)
}
