package handlers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

// @Summary Create a new Notification
// @Description Create a new Notification
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param notification body pb.Notification true "body for Notification"
// @Success 200 {object} string "Notification created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /notifications [post]
func (h *Handlers) CreateNotification(c *gin.Context) {
	var req pb.Notification
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Clients.Notification.Create(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to create Notification: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Notification created successfully")
	c.JSON(200, gin.H{"message": "Notification created successfully"})
}

// @Summary Get Notification by ID
// @Description Get Notification by ID
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Notification ID"
// @Success 200 {object} pb.Notification
// @Failure 400 {object} string "error": "error message"
// @Router /notifications/{id} [get]
func (h *Handlers) GetNotification(c *gin.Context) {
	id := c.Query("id")
	req := &pb.GetById{Id: id}

	res, err := h.Clients.Notification.Get(context.Background(), req)
	if err != nil {
		slog.Error("Failed to get Notification: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Notification retrieved successfully")
	c.JSON(200, res)
}

// @Summary List Notifications
// @Description List Notifications
// @Tags Notifications
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param receiver query string false "Receiver ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllNotificationRes
// @Failure 400 {object} string "error": "error message"
// @Router /notifications [get]
func (h *Handlers) ListNotifications(c *gin.Context) {
	var req pb.GetAllNotificationReq

	req.Receiver = c.Query("receiver")
	limit := c.Query("limit")
	offset := c.Query("offset")

	limitValue := 10
	offsetValue := 0

	if limit != "" {
		parsedLimit, err := strconv.Atoi(limit)
		if err != nil {
			slog.Error("Invalid limit value", err)
			c.JSON(400, gin.H{"error": "Invalid limit value"})
			return
		}
		limitValue = parsedLimit
	}

	if offset != "" {
		parsedOffset, err := strconv.Atoi(offset)
		if err != nil {
			slog.Error("Invalid offset value", err)
			c.JSON(400, gin.H{"error": "Invalid offset value"})
			return
		}
		offsetValue = parsedOffset
	}

	req.Filter = &pb.Filter{
		Limit:  int32(limitValue),
		Offset: int32(offsetValue),
	}

	res, err := h.Clients.Notification.List(context.Background(), &req)
	if err != nil {
		slog.Error("Error fetching Notifications", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Notifications retrieved successfully")
	c.JSON(200, res)
}
