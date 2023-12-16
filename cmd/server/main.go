package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"github.com/salmanbao/server/config"
	"github.com/salmanbao/server/controllers"

	dbConn "github.com/salmanbao/server/db/sqlc"
)

var (
	server *gin.Engine
	db     *dbConn.Queries
	conn   *pgx.Conn

	UserController controllers.UserController
)

func init() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	conn, err := pgx.Connect(context.Background(), config.PostgresURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	db = dbConn.New(conn)

	fmt.Println("PostgreSQL connected successfully...")

	UserController = *controllers.NewUserController(db)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	router := server.Group("/api")

	router.POST("/users", UserController.CreateUser)
	router.POST("/users/generateotp", UserController.GenerateOTP)
	router.POST("/users/verifyotp", UserController.VerifyOTP)

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Welcome to Golang with PostgreSQL"})
	})
	defer conn.Close(context.Background())

	log.Fatal(server.Run(":" + config.Port))
}
