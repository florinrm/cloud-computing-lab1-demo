package service

import (
	"database/sql"
	"exercise2/gateway"
	"exercise2/repository"
	"github.com/emicklei/go-restful/v3"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
)

const (
	pgHostKey     = "PGHOST"
	pgPortKey     = "PGPORT"
	pgDatabaseKey = "PGDATABASE"
	pgUserKey     = "PGUSER"
	pgPasswordKey = "PGPASSWORD"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) StartWebService() {
	ws := new(restful.WebService)
	restful.Add(ws)

	db, err := createDBConnection()
	if err != nil {
		return
	}

	// create logger
	l := logrus.New()

	api := gateway.NewAPI(repository.NewBookRepository(db, l))
	api.RegisterRoutes(ws)

	log.Printf("Started serving on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func createDBConnection() (*sql.DB, error) {
	sqlConnection, err := repository.CreatePostgresConnection(
		os.Getenv(pgHostKey),
		os.Getenv(pgPortKey),
		os.Getenv(pgDatabaseKey),
		os.Getenv(pgUserKey),
		os.Getenv(pgPasswordKey),
		"disable",
	)

	if err != nil {
		logrus.WithError(err).Errorf("failed to create database connection")
		return nil, err
	}

	return sqlConnection, nil
}
