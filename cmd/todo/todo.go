package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/demeesterdev/todo-service/pkg/todo"
	"github.com/go-kit/log"
)

const (
	defaultHTTPPort = "8081"
	defaultDBtarget = ":memory:"
)

func main() {
	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
		dbTarget = envString("DB_PATH_TODO", defaultDBtarget)
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	service, err := todo.NewSqliteDBService(dbTarget)
	if err != nil {
		panic(err)
	}

	var (
		eps         = todo.MakeServerEndpoints(service)
		httpHandler = todo.MakeHTTPHandler(eps, log.With(logger, "component", "HTTP"))
	)

	errs := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "SQL", "addr", dbTarget)
		logger.Log("transport", "HTTP", "addr", httpAddr)
		errs <- http.ListenAndServe(httpAddr, httpHandler)
	}()

	logger.Log("exit", <-errs)
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
