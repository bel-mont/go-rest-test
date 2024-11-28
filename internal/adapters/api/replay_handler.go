package api

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type ReplayHandler struct {
	s3Client *s3.Client
	repo     repository.Repository[entities.Replay]
}

func NewReplayHandler(s3Client *s3.Client, repo repository.Repository[entities.Replay]) ReplayHandler {
	return ReplayHandler{
		s3Client: s3Client,
		repo:     repo,
	}
}

const (
	// 2MB chunks - adjust based on your needs
	defaultChunkSize int64 = 2 * 1024 * 1024
)

func (h *ReplayHandler) StreamReplay(c *gin.Context) {
	id := c.Param("id")

	// Get replay info from repository
	replay, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		fmt.Printf("Error getting replay: %v\n", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Replay not found"})
		return
	}
	fmt.Printf("Found replay: %+v\n", replay)

	// Get file metadata
	headInput := &s3.HeadObjectInput{
		Bucket: aws.String(replay.S3Bucket),
		Key:    aws.String(replay.S3Path),
	}

	headOutput, err := h.s3Client.HeadObject(c.Request.Context(), headInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file info"})
		return
	}

	fileSize := *headOutput.ContentLength

	// Parse Range header
	var start, end int64
	rangeHeader := c.GetHeader("Range")
	if rangeHeader == "" {
		// If no range is specified, send the first chunk
		start = 0
		end = minClamp(defaultChunkSize-1, fileSize-1)
	} else {
		// Parse range header
		ranges := strings.Split(strings.TrimPrefix(rangeHeader, "bytes="), "-")
		if len(ranges) != 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid range header"})
			return
		}

		start, _ = strconv.ParseInt(ranges[0], 10, 64)
		if ranges[1] == "" {
			// If end is not specified, limit to chunk size from start
			end = minClamp(start+defaultChunkSize-1, fileSize-1)
		} else {
			requestedEnd, _ := strconv.ParseInt(ranges[1], 10, 64)
			// Limit the requested range to our chunk size
			end = minClamp(requestedEnd, start+defaultChunkSize-1)
		}

		if start > end || start < 0 || start >= fileSize {
			c.Header("Content-Range", fmt.Sprintf("bytes */%d", fileSize))
			c.Status(http.StatusRequestedRangeNotSatisfiable)
			return
		}
	}

	output, err := h.s3Client.GetObject(c.Request.Context(), &s3.GetObjectInput{
		Bucket: aws.String(replay.S3Bucket),
		Key:    aws.String(replay.S3Path),
		Range:  aws.String("bytes=0-11"), // Get first few bytes to check header
	})
	if err != nil {
		fmt.Printf("Error getting file header: %v\n", err)
		return
	}
	_, err = io.ReadAll(output.Body)
	if err != nil {
		fmt.Printf("Error reading file header: %v\n", err)
		return
	}
	err = output.Body.Close()
	if err != nil {
		fmt.Printf("Error closing file: %v\n", err)
		return
	}
	//fmt.Printf("File header bytes: %x\n", headerBytes)

	// Calculate content length for this chunk
	contentLength := end - start + 1

	// Get the chunk from S3
	input := &s3.GetObjectInput{
		Bucket: aws.String(replay.S3Bucket),
		Key:    aws.String(replay.S3Path),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", start, end)),
	}

	output, err = h.s3Client.GetObject(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file"})
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Printf("Error closing file: %v\n", err)
		}
	}(output.Body)

	// Set streaming-friendly headers
	c.Header("Content-Type", "video/mp4")
	c.Header("Content-Disposition", "inline")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Length", strconv.FormatInt(contentLength, 10))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("Connection", "keep-alive")

	// Set response status
	if rangeHeader != "" {
		c.Status(http.StatusPartialContent)
	} else {
		c.Status(http.StatusOK)
	}

	// Stream the content with a smaller buffer for more frequent updates
	buffer := make([]byte, 64*1024) // 64KB streaming buffer
	for {
		n, err := output.Body.Read(buffer)
		if err != nil && err != io.EOF {
			// Only log error message, not the buffer contents
			fmt.Printf("Error streaming file: %v\n", err)
			return
		}
		if n == 0 {
			break
		}

		// Disable any debug logging of the actual bytes
		if _, err := c.Writer.Write(buffer[:n]); err != nil {
			// Only log error message
			fmt.Printf("Error writing to response: %v\n", err)
			return
		}

		c.Writer.Flush()
	}
}

func minClamp(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
