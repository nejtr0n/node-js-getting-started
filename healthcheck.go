///bin/true; exec /usr/bin/env go run "$0" "$@"
package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/lib/pq"
)

var (
	errNodeJsIsDown   = errors.New("nodejs app is down")
	errPostgresIsDown = errors.New("postgres is down")
	errElasticIsDown  = errors.New("elastic is down")
)

func main() {
	nodejsAppUrl := getEnv("APP_NODEJS_URL", "http://localhost:8000")
	postgresDsn := getEnv("APP_POSTGRES_DSN", "postgres://user:pass@localhost:5432/db?sslmode=disable")
	elasticUrl := getEnv("APP_ELASTIC_URL", "http://localhost:9200")

	isNodejsAppAliveErr := isNodejsAppAlive(nodejsAppUrl)
	if isNodejsAppAliveErr != nil {
		log.Fatal(isNodejsAppAliveErr)
	}
	log.Println(fmt.Sprintf("nodejs app on address %s is alive!", nodejsAppUrl))

	isPostgresAliveErr := isPostgresAlive(postgresDsn)
	if isPostgresAliveErr != nil {
		log.Fatal(isPostgresAliveErr)
	}
	log.Println(fmt.Sprintf("postgres database on address %s is up!", postgresDsn))

	isElasticAliveErr := isElasticAlive(elasticUrl)
	if isElasticAliveErr != nil {
		log.Fatal(isElasticAliveErr)
	}
	log.Println(fmt.Sprintf("elasticsearch on address %s is up!", elasticUrl))

	os.Exit(0)
}

func isNodejsAppAlive(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%v: %w", err, errNodeJsIsDown)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%v: %w", err, errNodeJsIsDown)
	}

	return nil
}

func isPostgresAlive(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("%v: %w", err, errPostgresIsDown)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("%v: %w", err, errPostgresIsDown)
	}

	return nil
}

func isElasticAlive(url string) error {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{url},
	})
	if err != nil {
		return fmt.Errorf("%v: %w", err, errElasticIsDown)
	}

	res, err := es.Info()
	if err != nil {
		return fmt.Errorf("%v: %w", err, errElasticIsDown)
	}

	if res.IsError() {
		return fmt.Errorf("%s: %w", res.String(), errElasticIsDown)
	}

	return nil
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
