package handlers

import (
	"context"
	"database/sql"
	"errors"
	"ftp-server/models"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path/filepath"
)

func HandleUpload(c *gin.Context) {

	// Retrieve logger and db from context
	ctx := c.MustGet("ctx").(context.Context)
	logger := ctx.Value("logger").(*log.Logger)
	db := ctx.Value("db").(*sql.DB)
	config := ctx.Value("config").(models.Config)

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

func HandleGet(c *gin.Context) {
	// Retrieve logger and db from context
	ctx := c.MustGet("ctx").(context.Context)
	logger := ctx.Value("logger").(*log.Logger)
	db := ctx.Value("db").(*sql.DB)

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

func HandleDelete(c *gin.Context) {

	// Retrieve logger and db from context
	ctx := c.MustGet("ctx").(context.Context)
	logger := ctx.Value("logger").(*log.Logger)
	db := ctx.Value("db").(*sql.DB)

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
