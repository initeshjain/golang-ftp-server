// main.go
package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path/filepath"
)

// Config holds configuration parameters
type Config struct {
	UploadDir string `json:"UploadDir"`
	DbPath    string `json:"DbPath"`
}

var config Config

var logger *log.Logger

var db *sql.DB

func main() {
	initLogger()

	logger.Println("loading config...")
	loadConfig()

	logger.Println("init database...")
	initDatabase()

	router := gin.Default()

	router.POST("/upload", handleUpload)
	router.GET("/get/:filename", handleGet)
	router.DELETE("/delete/:filename", handleDelete)

	logger.Println("Server is running on :8080")
	err := router.Run(":8080")
	if err != nil {
		return
	}
}

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

func initLogger() {
	logFile, err := os.OpenFile("ftp_service.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening log file:", err)
	}
	logger = log.New(logFile, "FTP Service: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println("Logger initialized")
}

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

func handleUpload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		logger.Println("Bad Request -", err)
		c.JSON(400, gin.H{"error": "Bad Request"})
		return
	}

	// Check if the file already exists in the database
	var existingFilename string
	err = db.QueryRow("SELECT filename FROM files WHERE filename = ?", file.Filename).Scan(&existingFilename)
	if err == nil {
		// File with the same name already exists
		logger.Println("File already exists with same name -", file.Filename)
		c.JSON(409, gin.H{"error": "File already exists with same name"})
		return
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Other database error
		logger.Println("Error checking file existence -", err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Save the file to the specified upload directory
	dst := filepath.Join(config.UploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Println("Error saving uploaded file -", err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Insert metadata into the database
	_, err = db.Exec("INSERT INTO files (filename, path) VALUES (?, ?)", file.Filename, dst)
	if err != nil {
		logger.Fatal("Error inserting metadata into database:", err)
	}

	logger.Println("File uploaded successfully -", file.Filename)
	c.JSON(200, gin.H{"message": "File uploaded successfully"})
}

func handleGet(c *gin.Context) {
	filename := c.Param("filename")

	// Retrieve metadata from the database
	var path string
	err := db.QueryRow("SELECT path FROM files WHERE filename = ?", filename).Scan(&path)
	if err != nil {
		logger.Println("File not found -", err)
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	logger.Println("File retrieved -", filename)
	c.File(path)
}

func handleDelete(c *gin.Context) {
	filename := c.Param("filename")

	// Retrieve metadata from the database
	var path string
	err := db.QueryRow("SELECT path FROM files WHERE filename = ?", filename).Scan(&path)
	if err != nil {
		logger.Println("File not found -", err)
		c.JSON(404, gin.H{"error": "File not found"})
		return
	}

	// Delete file from the file system
	if err := os.Remove(path); err != nil {
		logger.Println("Error deleting file -", err)
		c.JSON(500, gin.H{"error": "Internal Server Error"})
		return
	}

	// Delete metadata from the database
	_, err = db.Exec("DELETE FROM files WHERE filename = ?", filename)
	if err != nil {
		logger.Fatal("Error deleting metadata from database:", err)
	}

	logger.Println("File deleted successfully -", filename)
	c.JSON(200, gin.H{"message": "File deleted successfully"})
}
