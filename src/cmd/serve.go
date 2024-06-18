package cmd

import (
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/app/handler"
	"github.com/mjedari/health-checker/app/services/healthsrv"
	"github.com/mjedari/health-checker/app/services/tasksrv"
	"github.com/mjedari/health-checker/infra/storage"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serving application",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
		start()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}

func start() {
	fmt.Println("Starting...")
	mux := http.NewServeMux()

	newMySQL, _ := storage.NewMySQL(config.Config.MySQL)
	memory := storage.NewInMemory()
	//redis, _ := storage.NewRedis(config.Config.Redis)
	newTaskService := tasksrv.NewTaskService(memory)
	healthService := healthsrv.NewHealthService(newMySQL, newTaskService)
	hh := handler.NewHealthHandler(healthService)

	// set routes
	mux.HandleFunc("GET /endpoint", hh.Index)
	mux.HandleFunc("POST /endpoint", hh.Create)
	mux.HandleFunc("GET /endpoint/{id}/start", hh.Start)
	mux.HandleFunc("GET /endpoint/{id}/stop", hh.Stop)
	mux.HandleFunc("DELETE /endpoint/{id}", hh.Delete)

	// todo: based on config
	err := http.ListenAndServe("localhost:8080", mux)
	if err != nil {
		log.Fatal("could not start server: ", err.Error())
	}

	// wait to end program by os interrupt or kill signal
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
	fmt.Println("\nShutting down...")
}
