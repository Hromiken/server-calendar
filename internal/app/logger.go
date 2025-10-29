package app

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// SetLogrus инициализация logrus
func SetLogrus(level string) {
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrusLevel)
	}
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})
	logrus.SetOutput(os.Stdout)
}

// LoggerMiddleware логирует каждый HTTP-запрос (метод, URL, время выполнения)
func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)

		logrus.Infof("%s %s %v", r.Method, r.URL.Path, duration)
	})
}
