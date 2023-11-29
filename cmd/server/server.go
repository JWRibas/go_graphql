package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"go_graphql/graph"
	"go_graphql/internal/database"
	"log"
	"os"
	"time"
)

const defaultPort = "8080"
const ginContextKey = "GinContextKey"

func LoadEnvVars() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func DatabaseCon() *database.ClientDB {
	dsn := os.Getenv("ROACHDB")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to open a DB connection: %s", err))
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Sprintf("Failed to ping DB: %s", err))
	}

	clientDB := &database.ClientDB{DB: db}
	return clientDB
}

func GinContextToContextMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.WithValue(c.Request.Context(), ginContextKey, c)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

var clientDB *database.ClientDB
var resolver *graph.Resolver

func graphqlHandler() gin.HandlerFunc {
	if clientDB == nil {
		log.Fatal("client database is not initiated")
	}
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	LoadEnvVars()
	clientDB = DatabaseCon()
	resolver = &graph.Resolver{ClientDB: clientDB}

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(GinContextToContextMiddleware())
	r.POST("/query", graphqlHandler())
	r.GET("/", playgroundHandler())
	err := r.Run()
	if err != nil {
		log.Fatal("Error running server: ", err)
	}
}
