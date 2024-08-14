package handlers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

// @Summary Create Wearable Data
// @Description Create Wearable Data
// @Tags WearableData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wearabledata body pb.WearableDataCreate true "body for WearableDataCreate"
// @Success 200 {object} string "Wearable data created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /wearable-data [post]
func (h *Handlers) CreateWearableData(c *gin.Context) {
	var req pb.WearableDataCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Clients.WearableData.Create(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to create WearableData: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("WearableData created successfully")
	c.JSON(200, gin.H{"message": "Wearable data created successfully"})
}

// @Summary Update Wearable Data
// @Description Update Wearable Data
// @Tags WearableData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Wearable Data ID"
// @Param wearabledata body pb.WearableDataUpdateBody true "body for WearableDataUpdate"
// @Success 200 {object} string "Wearable data updated successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /wearable-data/{id} [put]
func (h *Handlers) UpdateWearableData(c *gin.Context) {
	id := c.Query("id")

	var body pb.WearableDataUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req := &pb.WearableDataUpdate{
		Id:                id,
		DeviceType:        body.DeviceType,
		StepCount:         body.StepCount,
		CaloriesBurned:    body.CaloriesBurned,
		DistanceTraveled:  body.DistanceTraveled,
		HeartRate:         body.HeartRate,
		SleepDuration:     body.SleepDuration,
		WorkoutType:       body.WorkoutType,
		RecordedTimestamp: body.RecordedTimestamp,
	}

	_, err := h.Clients.WearableData.Update(context.Background(), req)
	if err != nil {
		slog.Error("Failed to update WearableData: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("WearableData updated successfully")
	c.JSON(200, gin.H{"message": "Wearable data updated successfully"})
}

// @Summary Delete Wearable Data
// @Description Delete Wearable Data
// @Tags WearableData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Wearable Data ID"
// @Success 200 {object} string "Wearable data deleted successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /wearable-data/{id} [delete]
func (h *Handlers) DeleteWearableData(c *gin.Context) {
	id := c.Query("id")

	_, err := h.Clients.WearableData.Delete(context.Background(), &pb.GetById{Id: id})
	if err != nil {
		slog.Error("Failed to delete WearableData: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("WearableData deleted successfully")
	c.JSON(200, gin.H{"message": "Wearable data deleted successfully"})
}

// @Summary Get Wearable Data
// @Description Get Wearable Data
// @Tags WearableData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Wearable Data ID"
// @Success 200 {object} pb.WearableDataRes
// @Failure 400 {object} string "error": "error message"
// @Router /wearable-data/{id} [get]
func (h *Handlers) GetWearableData(c *gin.Context) {
	var req pb.GetById

	id := c.Query("id")
	req.Id = id

	res, err := h.Clients.WearableData.Get(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to get WearableData: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("WearableData retrieved: %+v", res)
	c.JSON(200, res)
}

// @Summary List Wearable Data Records
// @Description List Wearable Data Records
// @Tags WearableData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "User ID"
// @Param recorded_timestamp query string false "Recorded Timestamp"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllWearableDataRes
// @Failure 400 {object} string "error": "error message"
// @Router /wearable-data [get]
func (h *Handlers) ListWearableData(c *gin.Context) {
	var req pb.GetAllWearableDataReq

	userId := c.Query("user_id")
	recorded_timestamp := c.Query("recorded_timestamp")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.UserId = userId
	req.RecordedTimestamp = recorded_timestamp

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

	res, err := h.Clients.WearableData.List(context.Background(), &req)
	if err != nil {
		slog.Error("Error fetching wearable data records", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("WearableData Records retrieved: %+v", res)
	c.JSON(200, res)
}
