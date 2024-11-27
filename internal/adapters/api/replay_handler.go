package api

import (
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"go-rest-test/internal/core/entities"
	"go-rest-test/internal/core/repository"
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

func (h *ReplayHandler) StreamReplay(c *gin.Context) {
	//id := c.Param("id")

	// Get file from storage
	//file, err := h.s3Client.GetReplayFile(id)
	//if err != nil {
	//	c.Status(http.StatusNotFound)
	//	return
	//}
	//
	//// Set streaming headers
	//c.Header("Content-Type", "video/mp4")
	//c.Header("Accept-Ranges", "bytes")
	//
	//// Handle range requests for seeking
	//http.ServeContent(c.Writer, c.Request, "", time.Time{}, file)
}
