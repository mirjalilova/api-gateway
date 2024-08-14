package handlers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

// @Summary Create a new Lifestyle Record
// @Description Create a new Lifestyle Record
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param lifestyle body pb.LifestyleCreate true "body for LifestyleCreate"
// @Success 200 {object} string "Record created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /lifestyle [post]
func (h *Handlers) CreateLifestyle(c *gin.Context) {
	var req pb.LifestyleCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Clients.LifeStyle.Create(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to create Lifestyle: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Lifestyle created successfully")
	c.JSON(200, gin.H{"message": "Record created successfully"})
}

// @Summary Update Lifestyle Record
// @Description Update Lifestyle Record
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Lifestyle Record ID"
// @Param lifestyle body pb.LifestyleUpdateBody true "body for LifestyleUpdate"
// @Success 200 {object} string "Record updated successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /lifestyle/{id} [put]
func (h *Handlers) UpdateLifestyle(c *gin.Context) {
	id := c.Query("id")

	var body pb.LifestyleUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req := &pb.LifestyleUpdate{
		Id:            id,
		DataType:      body.DataType,
		DataValue:     body.DataValue,
		HeartRate:     body.HeartRate,
		BloodPressure: body.BloodPressure,
		Weight:        body.Weight,
		Height:        body.Height,
		Temperature:   body.Temperature,
		StressLevel:   body.StressLevel,
		SleepDuration: body.SleepDuration,
	}

	_, err := h.Clients.LifeStyle.Update(context.Background(), req)
	if err != nil {
		slog.Error("Failed to update Lifestyle: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Lifestyle updated successfully")
	c.JSON(200, gin.H{"message": "Record updated successfully"})
}

// @Summary Delete Lifestyle Record
// @Description Delete Lifestyle Record
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Lifestyle Record ID"
// @Success 200 {object} string "Record deleted successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /lifestyle/{id} [delete]
func (h *Handlers) DeleteLifestyle(c *gin.Context) {
	id := c.Query("id")

	_, err := h.Clients.LifeStyle.Delete(context.Background(), &pb.GetById{Id: id})
	if err != nil {
		slog.Error("Failed to delete Lifestyle: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Lifestyle deleted successfully")
	c.JSON(200, gin.H{"message": "Record deleted successfully"})
}

// @Summary Get Lifestyle Record
// @Description Get Lifestyle Record
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Lifestyle Record ID"
// @Success 200 {object} pb.LifestyleRes
// @Failure 400 {object} string "error": "error message"
// @Router /lifestyle/{id} [get]
func (h *Handlers) GetLifestyle(c *gin.Context) {
	var req pb.GetById

	id := c.Query("id")
	req.Id = id

	res, err := h.Clients.LifeStyle.Get(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to get Lifestyle: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Lifestyle retrieved: %+v", res)
	c.JSON(200, res)
}

// @Summary List Lifestyle Records
// @Description List Lifestyle Records
// @Tags Lifestyle
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "User ID"
// @Param data_type query string false "Data Type"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllLifestyleRes
// @Failure 400 {object} string "error": "error message"
// @Router /lifestyle [get]
func (h *Handlers) ListLifestyle(c *gin.Context) {
	var req pb.GetAllLifestyleReq

	userId := c.Query("user_id")
	dataType := c.Query("data_type")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.UserId = userId
	req.DataType = dataType

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

	res, err := h.Clients.LifeStyle.List(context.Background(), &req)
	if err != nil {
		slog.Error("Error fetching lifestyle records", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Lifestyle Records retrieved: %+v", res)
	c.JSON(200, res)
}
