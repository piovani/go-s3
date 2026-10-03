package main

import (
	"io"
	"log"
	"net/http"

	"context"

	"github.com/gin-gonic/gin"
	"github.com/piovani/go-class/internal/config"
	"github.com/piovani/go-class/internal/storage/storage"
)

var storageClient *storage.S3Client

func main() {
	cfg := config.Load()
	var err error

	storageClient, err = storage.NewS3Client(context.Background(), cfg)
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}

	r := gin.Default()

	r.POST("/upload", UploadHandler)
	r.GET("/list", ListHandler)
	r.GET("/download", DownloadHandler)
	r.DELETE("/delete", DeleteHandler)

	if err := r.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}

func UploadHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	url, err := storageClient.Upload(c.Request.Context(), fileHeader.Filename, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "upload successful",
		"url":     url,
	})
}

func ListHandler(c *gin.Context) {
	files, err := storageClient.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"files": files,
	})
}

func DownloadHandler(c *gin.Context) {
	fileName := c.Query("file")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	data, err := storageClient.Download(c.Request.Context(), fileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+fileName)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func DeleteHandler(c *gin.Context) {
	fileName := c.Query("file")
	if fileName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	if err := storageClient.Delete(c.Request.Context(), fileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "delete successful",
	})
}
