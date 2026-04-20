package handlers

import (
	"net/http"
	"path/filepath"
	"time"
	"training-go/internal/repository"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
)

const bucket = "my-bucket"

// UploadFile godoc
// @Summary      Upload a file to object storage
// @Description  Upload a file via multipart/form-data to MinIO / Dell ObjectScale
// @Tags         files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file    true  "File to upload"
// @Param        path  formData  string  false "Custom object key/path (optional)"
// @Success      200   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /files/upload [post]
func UploadFileHandler(s3Client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		fileHeader, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
			return
		}

		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open file"})
			return
		}
		defer file.Close()

		objectKey := c.PostForm("path")
		if objectKey == "" {
			ext := filepath.Ext(fileHeader.Filename)
			name := fileHeader.Filename[:len(fileHeader.Filename)-len(ext)]
			objectKey = name + "_" + time.Now().Format("20060102_150405") + ext
		}

		contentType := fileHeader.Header.Get("Content-Type")
		if contentType == "" {
			contentType = "application/octet-stream"
		}

		if err := repository.UploadFileS3(s3Client, bucket, objectKey, file, fileHeader.Size, contentType); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload file: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"message":    "file uploaded successfully",
			"filename":   fileHeader.Filename,
			"object_key": objectKey,
			"size":       fileHeader.Size,
		})
	}
}

// GetFile godoc
// @Summary      Get file info / download URL
// @Description  Get a presigned download URL for a stored object
// @Tags         files
// @Produce      json
// @Param        key  query  string  true  "Object key"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files [get]
func GetFileHandler(s3Client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		objectKey := c.Query("key")
		if objectKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
			return
		}

		url, err := repository.GetPresignedURLs3(s3Client, bucket, objectKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate download URL: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":       "success",
			"object_key":   objectKey,
			"download_url": url,
			"expires_in":   "15 minutes",
		})
	}
}

// DeleteFile godoc
// @Summary      Delete a file
// @Description  Delete an object from storage by key
// @Tags         files
// @Produce      json
// @Param        key  query  string  true  "Object key"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /files [delete]
func DeleteFileHandler(s3Client *s3.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		objectKey := c.Query("key")
		if objectKey == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "key is required"})
			return
		}

		if err := repository.DeleteFileS3(s3Client, bucket, objectKey); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete file: " + err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "success",
			"message":    "file deleted successfully",
			"object_key": objectKey,
		})
	}
}
