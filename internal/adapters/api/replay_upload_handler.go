package api

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"go-rest-test/internal/infrastructure/auth"
	"html/template"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"
)

type ReplayUploadHandler struct {
	s3Client        *s3.Client
	repo            repository.Repository[entities.Replay]
	uploadStateRepo repository.Repository[entities.MultipartUpload]
}

func NewReplayUploadHandler(s3Client *s3.Client,
	repo repository.Repository[entities.Replay],
	uploadStateRepo repository.Repository[entities.MultipartUpload]) ReplayUploadHandler {
	return ReplayUploadHandler{s3Client: s3Client, repo: repo, uploadStateRepo: uploadStateRepo}
}

func (h ReplayUploadHandler) InitUploadHandler(c *gin.Context) {
	var req struct {
		FileName   string `json:"fileName"`
		FileSize   int64  `json:"fileSize"`
		TotalParts int    `json:"totalParts"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, _ := c.Get(auth.UserIDContextKey)
	userIDStr := userID.(string)

	// Generate S3 key
	key := fmt.Sprintf("uploads/%s-%s", time.Now().Format("20060102-150405"), req.FileName)
	bucketName := "fg-analyzer-replay-uploads"

	// Create multipart upload in S3
	createInput := &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
	}

	result, err := h.s3Client.CreateMultipartUpload(context.Background(), createInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize upload"})
		return
	}

	// Store upload state
	uploadInfo := entities.MultipartUpload{
		ID:             *result.UploadId,
		S3Key:          key,
		Status:         entities.MultipartUploadStatusInProgress,
		S3Bucket:       bucketName,
		FileName:       req.FileName,
		FileSize:       req.FileSize,
		TotalParts:     req.TotalParts,
		CompletedParts: make(map[int]string),
		UserID:         userIDStr,
		StartedAt:      time.Now(),
	}

	_, err = h.uploadStateRepo.Create(c.Request.Context(), uploadInfo)
	// Save upload state
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save upload state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"uploadId": *result.UploadId,
		"key":      key,
	})
}

func (h ReplayUploadHandler) GetUploadPartURL(c *gin.Context) {
	uploadID := c.Query("uploadId")
	partNumberStr := c.Query("partNumber")

	// Parse part number
	partNumber, err := strconv.ParseInt(partNumberStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid part number"})
		return
	}

	uploadInfo, err := h.uploadStateRepo.Get(c.Request.Context(), uploadID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	// Generate presigned URL for this part
	presignClient := s3.NewPresignClient(h.s3Client)
	presignedUrl, err := presignClient.PresignUploadPart(context.Background(), &s3.UploadPartInput{
		Bucket:     aws.String(uploadInfo.S3Bucket),
		Key:        aws.String(uploadInfo.S3Key),
		PartNumber: aws.Int32(int32(partNumber)),
		UploadId:   aws.String(uploadInfo.ID),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = time.Hour * 24 // URL valid for 24 hours
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate upload URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"url": presignedUrl.URL,
	})
}

func (h ReplayUploadHandler) CompletePart(c *gin.Context) {
	var req struct {
		UploadID   string `json:"uploadId"`
		PartNumber int    `json:"partNumber"`
		ETag       string `json:"etag"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	uploadInfo, err := h.uploadStateRepo.Get(c.Request.Context(), req.UploadID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	// Update completed parts
	uploadInfo.CompletedParts[req.PartNumber] = req.ETag

	// Save updated state
	if err := h.uploadStateRepo.Update(c.Request.Context(), uploadInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update upload state"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Part completed"})
}

func (h ReplayUploadHandler) CompleteUpload(c *gin.Context) {
	var req struct {
		UploadID string `json:"uploadId"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	uploadInfo, err := h.uploadStateRepo.Get(c.Request.Context(), req.UploadID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload not found"})
		return
	}

	// Prepare completed parts list
	var completedParts []types.CompletedPart
	for i := 1; i <= uploadInfo.TotalParts; i++ {
		etag, exists := uploadInfo.CompletedParts[i]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Part %d not completed", i)})
			return
		}
		completedParts = append(completedParts, types.CompletedPart{
			PartNumber: aws.Int32(int32(i)),
			ETag:       aws.String(etag),
		})
	}

	// Complete multipart upload
	_, err = h.s3Client.CompleteMultipartUpload(context.Background(), &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(uploadInfo.S3Bucket),
		Key:             aws.String(uploadInfo.S3Key),
		UploadId:        aws.String(uploadInfo.ID),
		MultipartUpload: &types.CompletedMultipartUpload{Parts: completedParts},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete upload"})
		return
	}

	// Create replay entity
	replayEntity := entities.Replay{
		UserID:     uploadInfo.UserID,
		UploadedAt: time.Now(),
		S3Bucket:   uploadInfo.S3Bucket,
		S3Path:     uploadInfo.S3Key,
		S3FileName: uploadInfo.FileName,
		S3FileSize: uploadInfo.FileSize,
	}

	newReplay, err := h.repo.Create(c.Request.Context(), replayEntity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create replay record"})
		return
	}

	// Clean up upload state
	if err := h.uploadStateRepo.Delete(c.Request.Context(), uploadInfo.ID); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to clean up upload state: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Upload completed successfully",
		"replayId": newReplay.ID,
	})
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

	userId, exists := c.Get(auth.UserIDContextKey)
	if !exists {
		h.sendError(c, http.StatusInternalServerError, "Failed to retrieve user ID", isHtmx, tmpl)
		return
	}
	userIdStr, ok := userId.(string)
	if !ok {
		h.sendError(c, http.StatusInternalServerError, "User ID is not a valid string", isHtmx, tmpl)
		return
	}
	fmt.Println("User ID: ", userIdStr)

	replayEntity := entities.Replay{
		UserID:     userIdStr,
		UploadedAt: time.Now(),
		S3Bucket:   bucketName,
		S3Path:     uploadKey,
		S3FileName: header.Filename,
		S3FileSize: header.Size,
	}
	newReplay, err := h.repo.Create(context.Background(), replayEntity)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "Failed to create replay: "+err.Error(), isHtmx, tmpl)
		return
	}
	fmt.Println("Uploaded replay: ", newReplay.ID)

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
