package handlers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

// @Summary Create Recommendation
// @Description Create Recommendation
// @Tags Recommendation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param recommendation body pb.RecommendationCreate true "body for RecommendationCreate"
// @Success 200 {object} string "Recommendation created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /recommendation [post]
func (h *Handlers) CreateRecommendation(c *gin.Context) {
	var req pb.RecommendationCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Clients.HealthRecommendations.Create(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to create Recommendation: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Recommendation created successfully")
	c.JSON(200, gin.H{"message": "Recommendation created successfully"})
}

// @Summary Update Recommendation
// @Description Update Recommendation
// @Tags Recommendation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Recommendation ID"
// @Param recommendation body pb.RecommendationUpdateBody true "body for RecommendationUpdate"
// @Success 200 {object} string "Recommendation updated successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /recommendation/{id} [put]
func (h *Handlers) UpdateRecommendation(c *gin.Context) {
	id := c.Query("id")

	var body pb.RecommendationUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req := &pb.RecommendationUpdate{
		Id:                 id,
		RecommendationType: body.RecommendationType,
		Description:        body.Description,
		Priority:           body.Priority,
	}

	_, err := h.Clients.HealthRecommendations.Update(context.Background(), req)
	if err != nil {
		slog.Error("Failed to update Recommendation: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Recommendation updated successfully")
	c.JSON(200, gin.H{"message": "Recommendation updated successfully"})
}

// @Summary Delete Recommendation
// @Description Delete Recommendation
// @Tags Recommendation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Recommendation ID"
// @Success 200 {object} string "Recommendation deleted successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /recommendation/{id} [delete]
func (h *Handlers) DeleteRecommendation(c *gin.Context) {
	id := c.Query("id")

	_, err := h.Clients.HealthRecommendations.Delete(context.Background(), &pb.GetById{Id: id})
	if err != nil {
		slog.Error("Failed to delete Recommendation: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Recommendation deleted successfully")
	c.JSON(200, gin.H{"message": "Recommendation deleted successfully"})
}

// @Summary Get Recommendation
// @Description Get Recommendation
// @Tags Recommendation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Recommendation ID"
// @Success 200 {object} pb.RecommendationRes
// @Failure 400 {object} string "error": "error message"
// @Router /recommendation/{id} [get]
func (h *Handlers) GetRecommendation(c *gin.Context) {
	var req pb.GetById

	id := c.Query("id")
	req.Id = id

	res, err := h.Clients.HealthRecommendations.Get(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to get Recommendation: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Recommendation retrieved: %+v", res)
	c.JSON(200, res)
}

// @Summary List Recommendations
// @Description List Recommendations
// @Tags Recommendation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "User ID"
// @Param recommendated_by query string false "Recommendated By"
// @Param priority query string false "Priority"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllRecommendationReq
// @Failure 400 {object} string "error": "error message"
// @Router /recommendation [get]
func (h *Handlers) ListRecommendations(c *gin.Context) {
	var req pb.GetAllRecommendationReq

	userId := c.Query("user_id")
	recommendated_by := c.Query("recommendated_by")
	priority := c.Query("priority")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.UserId = userId
	req.RecommendatedBy = recommendated_by

	if priority != "" {
		priorityInt, err := strconv.Atoi(priority)
		if err != nil {
			slog.Error("Invalid priority value", err)
			c.JSON(400, gin.H{"error": "Invalid priority value"})
			return
		}

		req.Priority = int32(priorityInt)
	} else {
		req.Priority = 0
	}

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

	res, err := h.Clients.HealthRecommendations.List(context.Background(), &req)
	if err != nil {
		slog.Error("Error fetching recommendations", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Recommendations retrieved: %+v", res)
	c.JSON(200, res)
}
