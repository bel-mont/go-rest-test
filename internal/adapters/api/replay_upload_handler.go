package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type ReplayUploadHandler struct {
	s3Client *s3.Client
}

func NewReplayUploadHandler(s3Client *s3.Client) ReplayUploadHandler {
	return ReplayUploadHandler{s3Client: s3Client}
}

func (h ReplayUploadHandler) UploadHandler(c *gin.Context) {
	tmpl, err := template.ParseFiles("web/components/files/upload-status.gohtml")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load upload status template"})
		return
	}
	// Check if it's an HTMX request
	isHtmx := c.GetHeader("HX-Request") == "true"

	// Get the title
	title := c.PostForm("title")
	if title == "" {
		h.sendError(c, http.StatusBadRequest, "Title is required", isHtmx, tmpl)
		return
	}

	// Get the file
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "File upload failed: "+err.Error(), isHtmx, tmpl)
		return
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			h.sendError(c, http.StatusInternalServerError, "Failed to close file: "+err.Error(), isHtmx, tmpl)
		}
	}(file)

	// Validate file size (e.g., 100MB limit)
	const maxSize = 100 << 20 // 100MB
	if header.Size > maxSize {
		h.sendError(c, http.StatusBadRequest, "File too large (max 100MB)", isHtmx, tmpl)
		return
	}

	// Read file into buffer
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		h.sendError(c, http.StatusInternalServerError, "Failed to process file: "+err.Error(), isHtmx, tmpl)
		return
	}

	// Upload to S3
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	bucketName := "fg-analyzer-replay-uploads"
	uploadKey := fmt.Sprintf("uploads/%s-%s", time.Now().Format("20060102-150405"), header.Filename)

	_, err = h.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(uploadKey),
		Body:        bytes.NewReader(buf.Bytes()),
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "Failed to upload file: "+err.Error(), isHtmx, tmpl)
		return
	}

	h.sendSuccess(c, header.Filename, isHtmx, tmpl)
}

func (h ReplayUploadHandler) sendError(c *gin.Context, status int, message string, isHtmx bool, tmpl *template.Template) {
	if isHtmx {
		err := tmpl.ExecuteTemplate(c.Writer, "partials/upload-status.gohtml", gin.H{
			"success": false,
			"message": message,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send error message"})
			return
		}
	} else {
		c.JSON(status, gin.H{"error": message})
	}
}

func (h ReplayUploadHandler) sendSuccess(c *gin.Context, filename string, isHtmx bool, tmpl *template.Template) {
	message := fmt.Sprintf("Successfully uploaded %s", filename)
	if isHtmx {
		err := tmpl.ExecuteTemplate(c.Writer, "components/upload-status.gohtml", gin.H{
			"success": true,
			"message": message,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send success message"})
			return
		}
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": message,
			"file":    filename,
		})
	}
}
