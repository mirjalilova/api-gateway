package handlers

import (
	"context"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

// @Summary Create a new Genetic Data record
// @Description Create a new Genetic Data record
// @Tags GeneticData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param geneticdata body pb.GeneticDataCreate true "body for GeneticDataCreate"
// @Success 200 {object} string "Record created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /genetic-data [post]
func (h *Handlers) CreateGeneticData(c *gin.Context) {
	var req pb.GeneticDataCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Clients.GeneticData.Create(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to create Genetic Data: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Genetic Data created successfully")
	c.JSON(200, gin.H{"message": "Record created successfully"})
}

// @Summary Update Genetic Data
// @Description Update Genetic Data
// @Tags GeneticData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Genetic Data ID"
// @Param geneticdata body pb.GeneticDataUpdate true "body for GeneticDataUpdate"
// @Success 200 {object} string "Record updated successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /genetic-data/{id} [put]
func (h *Handlers) UpdateGeneticData(c *gin.Context) {
	id := c.Query("id")

	var body pb.GeneticDataUpdate
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	body.Id = id
	_, err := h.Clients.GeneticData.Update(context.Background(), &body)
	if err != nil {
		slog.Error("Failed to update Genetic Data: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Genetic Data updated successfully")
	c.JSON(200, gin.H{"message": "Record updated successfully"})
}

// @Summary Delete Genetic Data
// @Description Delete Genetic Data
// @Tags GeneticData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Genetic Data ID"
// @Success 200 {object} string "Record deleted successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /genetic-data/{id} [delete]
func (h *Handlers) DeleteGeneticData(c *gin.Context) {
	id := c.Query("id")

	_, err := h.Clients.GeneticData.Delete(context.Background(), &pb.GetById{Id: id})
	if err != nil {
		slog.Error("Failed to delete Genetic Data: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Genetic Data deleted successfully")
	c.JSON(200, gin.H{"message": "Record deleted successfully"})
}

// @Summary Get Genetic Data by ID
// @Description Get Genetic Data by ID
// @Tags GeneticData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Genetic Data ID"
// @Success 200 {object} pb.GeneticDataRes
// @Failure 400 {object} string "error": "error message"
// @Router /genetic-data/{id} [get]
func (h *Handlers) GetGeneticData(c *gin.Context) {
	id := c.Query("id")
	req := &pb.GetById{Id: id}

	res, err := h.Clients.GeneticData.Get(context.Background(), req)
	if err != nil {
		slog.Error("Failed to get Genetic Data: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Genetic Data retrieved successfully")
	c.JSON(200, res)
}

// @Summary List Genetic Data records
// @Description List Genetic Data records
// @Tags GeneticData
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "User ID"
// @Param doctor_id query string false "Doctor ID"
// @Param analysis_date query string false "Analysis Date"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllGeneticDateRes
// @Failure 400 {object} string "error": "error message"
// @Router /genetic-data [get]
func (h *Handlers) ListGeneticData(c *gin.Context) {
	var req pb.GetAllGeneticDateReq

	req.UserId = c.Query("user_id")
	req.DoctorId = c.Query("doctor_id")
	req.AnalysisDate = c.Query("analysis_date")
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

	res, err := h.Clients.GeneticData.List(context.Background(), &req)
	if err != nil {
		slog.Error("Error fetching Genetic Data records", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Genetic Data records retrieved successfully")
	c.JSON(200, res)
}
