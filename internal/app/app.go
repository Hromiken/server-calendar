package app

import (
	"log"
	"os"
	"os/signal"
	"server-calendar/internal/storage"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"

	"server-calendar/cfg"
	"server-calendar/internal/app/httpserver"
	"server-calendar/internal/handler"
	"server-calendar/internal/service"
)

func Run(path string) {
	// config
	config, err := cfg.NewConfig(path)
	if err != nil {
		log.Fatal(err)
	}

	//logger
	SetLogrus(config.Log)

	// storage + service
	storage := storage.NewStorage()
	svc := service.NewCalendarService(storage)

	// router
	mux := handler.NewRouter(svc)
	muxWithLogs := LoggerMiddleware(mux)

	//server
	srv := httpserver.New(
		muxWithLogs,
		httpserver.Port(config.Port),
		httpserver.ReadTimeout(5*time.Second),
		httpserver.WriteTimeout(10*time.Second),
		httpserver.ShutdownTimeout(5*time.Second),
	)
	//graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		logrus.Infof("Received signal: %v", sig)
		e := srv.Shutdown()
		if e != nil {
			logrus.Errorf("Failed to shutdown server: %v", e)
		}
	case errNotify := <-srv.Notify():
		logrus.Errorf("Server exited with error: %v", errNotify)
	}

}
