package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/demeesterdev/todo-service/internal/argon2id"
	"github.com/demeesterdev/todo-service/pkg/authorization"
	authorizationEps "github.com/demeesterdev/todo-service/pkg/authorization/endpoints"
	authorizationTrsp "github.com/demeesterdev/todo-service/pkg/authorization/transport"
	"github.com/go-kit/log"
)

const (
	defaultHTTPPort = "8082"
	defaultDBtarget = ":memory:"
)

func main() {
	var (
		logger   log.Logger
		httpAddr = net.JoinHostPort("localhost", envString("HTTP_PORT", defaultHTTPPort))
		dbTarget = envString("DB_PATH_AUTH", defaultDBtarget)
	)

	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	service, err := authorization.NewSqliteDBService(dbTarget, argon2id.DefaultConfig)
	if err != nil {
		panic(err)
	}

	var (
		eps         = authorizationEps.MakeServerEndpoints(service)
		httpHandler = authorizationTrsp.MakeHTTPHandler(eps, log.With(logger, "component", "HTTP"))
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
