package cmd

import (
	"fmt"
	"github.com/mjedari/health-checker/app/config"
	"github.com/mjedari/health-checker/app/contract"
	"github.com/mjedari/health-checker/app/handler"
	"github.com/mjedari/health-checker/app/services/healthsrv"
	"github.com/mjedari/health-checker/app/services/httpsrv"
	"github.com/mjedari/health-checker/app/services/tasksrv"
	"github.com/mjedari/health-checker/domain"
	"github.com/mjedari/health-checker/infra/storage"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"log"
	"net"
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
	mux := http.NewServeMux()

	newMySQL, _ := storage.NewMySQL(config.Config.MySQL)

	cache := newCache()
	pool := domain.NewTaskPool()
	newTaskService := tasksrv.NewTaskService(cache, pool)
	clientService := httpsrv.NewHttpService(config.Config.Webhook)
	healthService := healthsrv.NewHealthService(clientService, newMySQL, newTaskService)
	hh := handler.NewHealthHandler(healthService)

	// set routes
	mux.HandleFunc("GET /endpoint", hh.Index)
	mux.HandleFunc("POST /endpoint", hh.Create)
	mux.HandleFunc("GET /endpoint/{id}/start", hh.Start)
	mux.HandleFunc("GET /endpoint/{id}/stop", hh.Stop)
	mux.HandleFunc("DELETE /endpoint/{id}", hh.Delete)

	// start web server
	runHTTPServer(mux)

	// wait to end program by os interrupt or kill signal
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
	fmt.Println("\nShutting down...")
}

func runHTTPServer(mux *http.ServeMux) {
	go func() {
		err := http.ListenAndServe(net.JoinHostPort(config.Config.Server.Host, config.Config.Server.Port), mux)
		if err != nil {
			log.Fatal("could not start server: ", err.Error())
		}

	}()

	logrus.WithField("HTTP_Host", config.Config.Server.Host).
		WithField("HTTP_Port", config.Config.Server.Port).
		Info("starting HTTP/REST health-checker...")
}

func newCache() contract.ICache {
	switch config.Config.Service.Cache {
	case "redis":
		redis, _ := storage.NewRedis(config.Config.Redis)
		return redis
	case "memory":
		return storage.NewInMemory()
	default:
		return storage.NewInMemory()
	}
}
