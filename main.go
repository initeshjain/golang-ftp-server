// main.go
package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"ftp-server/handlers"
	"ftp-server/models"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

var config models.Config

var logger *log.Logger

var db *sql.DB

func main() {
	initLogger()

	logger.Println("loading config...")
	loadConfig()

	logger.Println("init database...")
	initDatabase()

	// Adding logger, db instance, and config instance into context
	ctx := context.WithValue(context.Background(), "logger", logger)
	ctx = context.WithValue(ctx, "db", db)
	ctx = context.WithValue(ctx, "config", config)

	router := gin.Default()
	router.Use(ContextMiddleware(ctx))

	router.POST("/upload", handlers.HandleUpload)
	router.GET("/get/:filename", handlers.HandleGet)
	router.DELETE("/delete/:filename", handlers.HandleDelete)

	logger.Println("Server is running on :8080")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}

// ContextMiddleware Middleware to handle and pass context to request handlers
func ContextMiddleware(ctx context.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Add context to gin request context
		c.Set("ctx", ctx)
		c.Next()
	}
}

// Loading config from config.json file
func loadConfig() {
	// Load configuration from a file (config.json)
	configFile, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Error opening config file:", err)
	}
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			logger.Fatal("Error loading config file: ", err.Error())
		}
	}(configFile)

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Fatal("Error decoding config file:", err)
	}
}

// Creating instance of logger
func initLogger() {
	logFile, err := os.OpenFile("ftp_service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	logger = log.New(logFile, "FTP Service: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Logger initialized")
}

// Creating instance of database
func initDatabase() {
	var err error
	db, err = sql.Open("sqlite3", config.DbPath)
	if err != nil {
		logger.Fatal("Error opening database:", err)
	}

	// Create files table if not exists
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS files (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            filename TEXT,
            path TEXT
        )
    `)
	if err != nil {
		logger.Fatal("Error creating files table:", err)
	}
}
