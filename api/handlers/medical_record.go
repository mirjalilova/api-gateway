package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

// @Summary Create a new Medical Record
// @Description Create a new Medical Record
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param medicalrecord body pb.MedicalRecordCreate true "body for MedicalRecordCreate"
// @Success 200 {object} string "Record created successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /medical-records [post]
func (h *Handlers) CreateMedicalRecord(c *gin.Context) {
	var req pb.MedicalRecordCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	input, err := json.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal JSON: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	err = h.Producer.ProduceMessages("cre-med", input)
	if err != nil {
		slog.Error("Failed to produce message: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}
	
	// _, err := h.Clients.MedicalRecords.Create(context.Background(), &req)
	// if err != nil {
	// 	slog.Error("Failed to create Athlete: ", err)
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	slog.Info("Medical Record created: %+v", err)
	c.JSON(200, gin.H{"message": "Record created successfully"})
}

// @Summary Update Medical Record
// @Description Update Medical Record
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Medical Record ID"
// @Param Athlete body pb.MedicalRecordUpdateBody true "body for Medical Record Update Body"
// @Success 200 {object} string "Record updated successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /medical-records/{id} [put]
func (h *Handlers) UpdateMedicalRecord(c *gin.Context) {
	id := c.Query("id")

	var body pb.MedicalRecordUpdateBody
	if err := c.ShouldBindJSON(&body); err != nil {
		slog.Error("Failed to bind JSON: ", err)
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	req := &pb.MedicalRecordUpdate{
		Id:          id,
		RecordType:  body.RecordType,
		RecordDate:  body.RecordDate,
		Description: body.Description,
		DoctorId:    "a8d52519-950f-4903-96a1-9c90354cc197",
		Attachments: body.Attachments,
	}

	input, err := json.Marshal(req)
	if err != nil {
		slog.Error("Failed to marshal JSON: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	err = h.Producer.ProduceMessages("upd-med", input)
	if err != nil {
		slog.Error("Failed to produce message: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	// _, err := h.Clients.MedicalRecords.Update(context.Background(), req)
	// if err != nil {
	// 	slog.Error("Failed to update Medical Record: ", err)
	// 	c.JSON(500, gin.H{"error": err.Error()})
	// 	return
	// }

	slog.Info("Medical Record updated successfully")
	c.JSON(200, gin.H{"message": "Record updated successfully"})
}

// @Summary Delete Medical Record
// @Description Delete Medical Record
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Medical Record ID"
// @Success 200 {object} string "Record deleted successfully"
// @Failure 400 {object} string "error": "error message"
// @Router /medical-records/{id} [delete]
func (h *Handlers) DeleteMedicalRecord(c *gin.Context) {
	id := c.Query("id")

	_, err := h.Clients.MedicalRecords.Delete(context.Background(), &pb.GetById{Id: id})
	if err != nil {
		slog.Error("Failed to delete Medical Record: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Medical Record deleted successfully")
	c.JSON(200, gin.H{"message": "Record deleted successfully"})
}

// @Summary Get Medical Record with the given details
// @Description Get Medical Record with the given details
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id query string true "Medical Record ID"
// @Success 200 {object} pb.MedicalRecordRes
// @Failure 400 {object} string "error": "error message"
// @Router /medical-records/{id} [get]
func (h *Handlers) GetMedicalRecord(c *gin.Context) {
	var req pb.GetById

	id := c.Query("id")
	req.Id = id

	res, err := h.Clients.MedicalRecords.Get(context.Background(), &req)
	if err != nil {
		slog.Error("Failed to get Athlete: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Athlete retrieved: %+v", res)
	c.JSON(200, res)
}

// @Summary List Medical Records
// @Description List Medical Records
// @Tags MedicalRecords
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "User ID"
// @Param record_type query string false "Record Type"
// @Param record_date query string false "Record Date"
// @Param doctor_id query string false "Doctor ID"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} pb.GetAllMedicalRecordsRes
// @Failure 400 {object} string "error": "error message"
// @Router /medical-records [get]
func (h *Handlers) ListMedicalRecords(c *gin.Context) {
	var req pb.GetAllMedicalRecords

	userId := c.Query("user_id")
	recordType := c.Query("record_type")
	recordDate := c.Query("record_date")
	doctorId := c.Query("doctor_id")
	limit := c.Query("limit")
	offset := c.Query("offset")

	req.RecordType = recordType
	req.RecordDate = recordDate
	req.DoctorId = doctorId
	req.UserId = userId

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

	res, err := h.Clients.MedicalRecords.List(context.Background(), &req)
	if err != nil {
		slog.Error("Error fetching medical records", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	slog.Info("Medical Records retrieved: %+v", res)
	c.JSON(200, res)
}
