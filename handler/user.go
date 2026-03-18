package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	appkafka "github.com/mamh1019/lambda-go-server-boilerplate/kafka"
	"github.com/mamh1019/lambda-go-server-boilerplate/user"
)

type UserHandler struct {
	repo     *user.Repository
	producer *appkafka.Producer
}

func NewUserHandler(repo *user.Repository, producer *appkafka.Producer) *UserHandler {
	return &UserHandler{
		repo:     repo,
		producer: producer,
	}
}

func (h *UserHandler) List(c *gin.Context) {
	users, err := h.repo.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users})
}

func (h *UserHandler) Get(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	u, err := h.repo.Get(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if u == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

type createUserRequest struct {
	UserName int `json:"user_name" binding:"required"`
	Coin     int `json:"coin"`
	Jewel    int `json:"jewel"`
}

func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := h.repo.Create(c.Request.Context(), req.UserName, req.Coin, req.Jewel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.producer.PublishUserCreated(c.Request.Context(), id); err != nil {
		log.Printf("failed to publish user created event: %v", err)
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": id})
}

type updateUserRequest struct {
	UserName int `json:"user_name" binding:"required"`
	Coin     int `json:"coin"`
	Jewel    int `json:"jewel"`
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.repo.Update(c.Request.Context(), id, req.UserName, req.Coin, req.Jewel); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.repo.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
