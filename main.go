package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"slices"
	"time"

	"github.com/gorilla/mux"
	"github.com/ronymmoura/chronos/internal/db"
	"github.com/ronymmoura/chronos/internal/db/repositories"
	httpServer "github.com/ronymmoura/chronos/internal/http"
	"github.com/ronymmoura/chronos/internal/models"
	"github.com/ronymmoura/chronos/internal/util"
)

type App struct {
	Server      *httpServer.Server
	ProcessList []*models.Process
}

type Action struct {
	op string
	id int32
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credential", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, X-Auth-Token, Authorization")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	c := make(chan bool, 1)

	config := util.LoadConfig(".env")
	db := db.Connect(config.ConnectionString)
	router := mux.NewRouter()

	server := httpServer.NewServer(db, config, router, c)
	server.RegisterRoutes()

	app := &App{
		Server:      server,
		ProcessList: []*models.Process{},
	}

	err := app.RegisterProcessList()
	if err != nil {
		fmt.Printf("Error registering process list: %v\n", err)
		return
	}

	go func() {
		//processRepo := repositories.NewProcessRepository(app.Server.DB)
		processRunRepo := repositories.NewProcessRunRepository(app.Server.DB)
		processRunLogRepository := repositories.NewProcessRunLogRepository(app.Server.DB)
		processRepo := repositories.NewProcessRepository(processRunRepo.DB)

		runningList := []int32{}
		runningListChannel := make(chan *Action)

		go (func() {
			for a := range runningListChannel {
				id := a.id

				if a.op == "add" {
					runningList = append(runningList, id)
				} else if a.op == "remove" {
					idx := slices.Index(runningList, id)
					if idx != -1 {
						runningList = slices.Delete(runningList, idx, idx+1)
					}
				}

				fmt.Printf("Running list: %v\n", runningList)
			}
		})()

		for {
			select {
			case <-app.Server.Channel:
				err := app.RegisterProcessList()
				if err != nil {
					fmt.Printf("Error registering process list: %v\n", err)
					return
				}
				fmt.Println("Yay")
				continue
			default:
				for _, process := range app.ProcessList {
					if !slices.Contains(runningList, process.ID) {
						runningListChannel <- &Action{op: "add", id: process.ID}
						go runProcess(process, runningListChannel, processRepo, processRunRepo, processRunLogRepository)
					}
				}
			}
		}
	}()

	http.ListenAndServe(":8080", corsMiddleware(app.Server.Router))
}

func (app *App) RegisterProcessList() error {
	processList, err := repositories.NewProcessRepository(app.Server.DB).SelectActives()
	if err != nil {
		return err
	}

	app.ProcessList = processList
	return nil
}

func runProcess(process *models.Process, runningListChannel chan *Action, processRepo *repositories.ProcessRepository, processRunRepo *repositories.ProcessRunRepository, processRunLogRepo *repositories.ProcessRunLogRepository) {
	defer func() {
		runningListChannel <- &Action{op: "remove", id: process.ID}
	}()

	err := processRepo.UpdateRunning(process.ID, true)
	if err != nil {
		fmt.Printf("Error updating process running status: %v\n", err)
		return
	}

	fmt.Printf("Running %s...\n", process.Name)
	startTime := time.Now()

	filePath := "./" + filepath.Base(process.Path)
	dir := filepath.Dir(process.Path)
	cmd := exec.Command(filePath)
	cmd.Dir = dir

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error creating stdout pipe: %v\n", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Printf("Error creating stderr pipe: %v\n", err)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error starting command: %v\n", err)
	}

	errors := []string{}
	scannerErr := bufio.NewScanner(stderr)
	for scannerErr.Scan() {
		errors = append(errors, scannerErr.Text())
	}

	logs := []string{}
	scannerOut := bufio.NewScanner(stdout)
	for scannerOut.Scan() {
		logs = append(logs, scannerOut.Text())
	}

	fmt.Printf("Finished running process %s\n", process.Name)
	endTime := time.Now()
	processRunId, err := processRunRepo.Insert(&models.ProcessRun{
		ProcessID: process.ID,
		StartedAt: startTime,
		EndedAt:   endTime,
		Success:   len(errors) == 0,
	})
	if err != nil {
		fmt.Printf("Error inserting process run: %v\n", err)
	}

	for _, log := range logs {
		_, err := processRunLogRepo.Insert(&models.ProcessRunLog{
			ProcessRunID: processRunId,
			LogTime:      time.Now(),
			Message:      log,
			Type:         "info",
		})
		if err != nil {
			fmt.Printf("Error inserting process run log: %v\n", err)
		}
	}

	for _, log := range errors {
		_, err := processRunLogRepo.Insert(&models.ProcessRunLog{
			ProcessRunID: processRunId,
			LogTime:      time.Now(),
			Message:      log,
			Type:         "error",
		})
		if err != nil {
			fmt.Printf("Error inserting process run log: %v\n", err)
		}
	}
	err = processRepo.UpdateRunning(process.ID, false)
	if err != nil {
		fmt.Printf("Error updating process running status: %v\n", err)
		return
	}

	time.Sleep(time.Duration(process.ExecuteEverySecs) * time.Second)
}
