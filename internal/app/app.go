package app

import (
	"log"
	"os"
	"os/signal"
	"server-calendar/internal/worker"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"server-calendar/cfg"
	"server-calendar/internal/app/httpserver"
	"server-calendar/internal/handler"
	log2 "server-calendar/internal/log"
	"server-calendar/internal/service"
	"server-calendar/internal/storage"
)

func Run(path string) {
	config, err := cfg.NewConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	log2.SetLogrus(config.Log)

	stg := storage.NewStorage()
	svc := service.NewCalendarService(stg)

	asyncLogger := log2.NewAsyncLogger(100)
	asyncLogger.Start()
	defer asyncLogger.Stop()

	archiveWorker := worker.NewArchiveWorker(
		5*time.Minute,
		stg.ArchiveOldEvents,
	)
	archiveWorker.Start()
	defer archiveWorker.Stop()

	router := handler.NewRouter(svc)

	handlerWithLogs := log2.LoggerMiddleware(asyncLogger)(router)

	srv := httpserver.New(
		handlerWithLogs,
		httpserver.Port(config.Port),
		httpserver.ReadTimeout(5*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.ShutdownTimeout(5*time.Second),
	)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logrus.Infof("Received signal: %v", sig)
		if err = srv.Shutdown(); err != nil {
			logrus.Errorf("Failed to shutdown server: %v", err)
		}
	case errNotify := <-srv.Notify():
		logrus.Errorf("Server exited with error: %v", errNotify)
	}
}
