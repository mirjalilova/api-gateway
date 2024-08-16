package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	pb "github.com/mirjalilova/api-gateway/genproto/health_sync"
	"golang.org/x/exp/slog"
)

var cacheExpiration = 15 * time.Minute

// Cache the data in Redis
func cacheData(ctx context.Context, client *redis.Client, key string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return client.Set(ctx, key, jsonData, cacheExpiration).Err()
}

// Get the data from Redis
func getCachedData(ctx context.Context, client *redis.Client, key string, dest interface{}) error {
	jsonData, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(jsonData, dest)
}

// @Summary Health Monitoring Weekly Report
// @Description Health Monitoring Weekly Report
// @Tags Health Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string true "User ID"
// @Success 200 {object} pb.WeeklySummary
// @Failure 400 {object} string "error": "error message"
// @Router /health-monitoring/weekly-summary/{id} [get]
func (h *Handlers) HealthMonitoringWeeklyReport(c *gin.Context) {
	id := c.Query("user_id")
	cacheKey := "weekly:" + id
	ctx := context.Background()

	var res pb.WeeklySummary

	err := getCachedData(ctx, h.Redis, cacheKey, &res)
	if err == nil {
		slog.Info("Weekly data retrieved from cache")
		c.JSON(200, res)
		return
	}

	req := &pb.WeeklySummaryReq{UserId: id}
	resp, err := h.Clients.LifeStyle.GetWeeklySummary(ctx, req)
	if err != nil {
		slog.Error("Failed to get weekly summary: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res = *resp

	cacheData(ctx, h.Redis, cacheKey, res)

	slog.Info("Weekly summary retrieved and cached successfully")
	c.JSON(200, res)
}

// @Summary Health Monitoring Real Time Report
// @Description Health Monitoring Real Time Report
// @Tags Health Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string true "User ID"
// @Success 200 {object} pb.LifestyleRes
// @Failure 400 {object} string "error": "error message"
// @Router /health-monitoring/real-time/{id} [get]
func (h *Handlers) HealthMonitoringRealTime(c *gin.Context) {
	id := c.Query("user_id")
	cacheKey := "realtime:" + id
	ctx := context.Background()

	var res pb.LifestyleRes

	err := getCachedData(ctx, h.Redis, cacheKey, &res)
	if err == nil {
		slog.Info("Real-time data retrieved from cache")
		c.JSON(200, res)
		return
	}

	resp, err := h.Clients.WearableData.GetRealTimeMonitoringData(ctx, &pb.GetById{Id: id})
	if err != nil {
		slog.Error("Failed to get real-time data: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res = pb.LifestyleRes{
		HeartRate:    resp.HeartRate,
		SleepDuration: resp.SleepDuration,
		Temperature:  resp.Temperature,
	}

	cacheData(ctx, h.Redis, cacheKey, res)

	slog.Info("Real-time data retrieved and cached successfully")
	c.JSON(200, res)
}


// @Summary Health Monitoring Daily Report
// @Description Health Monitoring Daily Report
// @Tags Health Monitoring
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user_id query string true "User ID"
// @Success 200 {object} pb.LifestyleRes
// @Failure 400 {object} string "error": "error message"
// @Router /health-monitoring/daily-summary/{id} [get]
func (h *Handlers) HealthMonitoringDailyReport(c *gin.Context) {
	id := c.Query("user_id")
	cacheKey := "daily:" + id
	ctx := context.Background()

	var res pb.LifestyleRes

	err := getCachedData(ctx, h.Redis, cacheKey, &res)
	if err == nil {
		slog.Info("Daily data retrieved from cache")
		c.JSON(200, res)
		return
	}

	req := &pb.GetById{Id: id}
	resp, err := h.Clients.LifeStyle.Get(ctx, req)
	if err != nil {
		slog.Error("Failed to get daily summary: ", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	res = *resp

	cacheData(ctx, h.Redis, cacheKey, res)

	slog.Info("Daily summary retrieved and cached successfully")
	c.JSON(200, res)
}
