package main

import (
	"fmt"
	"net/http"
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

func main() {
	config := util.LoadConfig(".env")
	db := db.Connect(config.ConnectionString)
	router := mux.NewRouter()

	server := httpServer.NewServer(db, config, router)
	server.RegisterRoutes()

	app := &App{
		Server:      server,
		ProcessList: []*models.Process{},
	}

	app.RegisterProcessList()

	go (func() {
		lastRun := time.Now()
		// nextRun := lastRun.Add(time.Second * 2)
		count := 1

		for {
			for _, process := range app.ProcessList {
				// fmt.Println(process.Name)
				// fmt.Println(count)
				// fmt.Println(lastRun)
				// fmt.Println(nextRun)
				// fmt.Println(time.Now().Sub(nextRun))
				time.Sleep(time.Second * time.Duration(process.ExecuteEverySecs))
				// fmt.Println(time.Now())
				// fmt.Println(time.Now().Sub(nextRun))
				// fmt.Println("-----------------")

				//if time.Now().Sub(nextRun) >= 0 {
				fmt.Printf("Running %s...\n", process.Name)
				count++
				lastRun = time.Now()
				// nextRun := lastRun.Add(time.Second * time.Duration(process.ExecuteEverySecs))
				repositories.NewProcessRepository(app.Server.DB).UpdateStatus(process.ID, lastRun)
				//}
			}
		}
	})()

	http.ListenAndServe(":8080", app.Server.Router)
}

func (app *App) RegisterProcessList() error {
	processList, err := repositories.NewProcessRepository(app.Server.DB).SelectAll()
	if err != nil {
		return err
	}

	app.ProcessList = processList
	return nil
}
