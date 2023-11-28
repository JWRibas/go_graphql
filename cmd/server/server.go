package main

import (
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go_graphql/graph"
	"go_graphql/internal/database"
	"log"
	"net/http"
	"os"
)

const defaultPort = "8080"

func LoadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	LoadEnvVars()

	dsn := os.Getenv("ROACHDB")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to open a DB connection: %s", err))
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("Failed to ping DB: %s", err))
	}
	defer db.Close()

	clientDB := database.ClientDB{DB: db}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	resolver := &graph.Resolver{
		ClientDB: &clientDB,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
