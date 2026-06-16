package api

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"dujiangyan-system/pkg/models"
	"dujiangyan-system/pkg/mqtt"
	"dujiangyan-system/pkg/simulation"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type APIHandler struct {
	ctx context.Context
}

func NewAPIHandler(ctx context.Context) *APIHandler {
	return &APIHandler{ctx: ctx}
}

func (h *APIHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")

	hydrology := api.Group("/hydrology")
	{
		hydrology.POST("/data", h.ReceiveHydrologyData)
		hydrology.GET("/data/:station_id", h.GetHydrologyData)
		hydrology.GET("/data/latest/:station_id", h.GetLatestHydrologyData)
		hydrology.GET("/data/all", h.GetAllLatestData)
		hydrology.GET("/stats/daily/:station_id", h.GetDailyStats)
		hydrology.GET("/stations", h.GetMonitoringStations)
	}

	api.GET("/wolong-iron", h.GetWolongIron)

	alerts := api.Group("/alerts")
	{
		alerts.GET("", h.GetAlerts)
		alerts.POST("/:id/acknowledge", h.AcknowledgeAlert)
		alerts.GET("/unpublished", h.GetUnpublishedAlerts)
	}

	prediction := api.Group("/prediction")
	{
		prediction.POST("/bed-evolution/:station_id", h.RunBedEvolutionPrediction)
		prediction.GET("/bed-evolution/:station_id", h.GetBedEvolutionPredictions)
	}

	simulation := api.Group("/simulation")
	{
		simulation.POST("/bamboo-cage", h.RunBambooCageSimulation)
		simulation.POST("/macha-interception", h.RunMachaInterceptionSimulation)
		simulation.GET("/list", h.GetSimulations)
		simulation.GET("/macha/:simulation_id", h.GetMachaSimulationData)
		simulation.GET("/bamboo-cage/:simulation_id", h.GetBambooCageSimulationData)
	}

	api.GET("/annual-repair-records", h.GetAnnualRepairRecords)
	api.GET("/dem-grid", h.GetDEMGrid)
	api.GET("/evolution-rate/:station_id", h.GetEvolutionRate)

	api.GET("/dynasty-techniques", h.GetDynastyTechniques)
	api.GET("/modern-comparisons", h.GetModernComparisons)

	earthquake := api.Group("/earthquake")
	{
		earthquake.POST("/simulate", h.RunEarthquakeSimulation)
		earthquake.GET("/list", h.GetEarthquakeSimulations)
	}

	user := api.Group("/user")
	{
		user.POST("/operation", h.InsertUserOperation)
		user.GET("/operations/:session_id", h.GetUserOperations)
		user.POST("/session/finish", h.FinishUserSession)
		user.GET("/ranking", h.GetUserRanking)
	}

	api.GET("/ws/realtime", h.RealTimeWebSocket)
}

type HydrologyDataRequest struct {
	StationID             string    `json:"station_id" binding:"required"`
	StationName           string    `json:"station_name"`
	WaterLevel            float64   `json:"water_level" binding:"required"`
	FlowRate              float64   `json:"flow_rate" binding:"required"`
	SedimentConcentration float64   `json:"sediment_concentration" binding:"required"`
	BedElevation          float64   `json:"bed_elevation" binding:"required"`
	Temperature           float64   `json:"temperature"`
	Rainfall              float64   `json:"rainfall"`
	SensorStatus          int       `json:"sensor_status"`
	Time                  time.Time `json:"time"`
}

func (h *APIHandler) ReceiveHydrologyData(c *gin.Context) {
	var req HydrologyDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Time.IsZero() {
		req.Time = time.Now()
	}

	data := &models.HydrologyData{
		Time:                  req.Time,
		StationID:             req.StationID,
		StationName:           req.StationName,
		WaterLevel:            req.WaterLevel,
		FlowRate:              req.FlowRate,
		SedimentConcentration: req.SedimentConcentration,
		BedElevation:          req.BedElevation,
		Temperature:           req.Temperature,
		Rainfall:              req.Rainfall,
		SensorStatus:          req.SensorStatus,
	}

	if data.SensorStatus == 0 {
		data.SensorStatus = 1
	}

	if err := models.InsertHydrologyData(h.ctx, data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert data: " + err.Error()})
		return
	}

	go func() {
		if err := mqtt.PublishHydrologyData(data); err != nil {
		}
		broadcastHydrologyData(data)
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "Data received successfully",
		"data":    data,
	})
}

func (h *APIHandler) GetHydrologyData(c *gin.Context) {
	stationID := c.Param("station_id")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")
	limitStr := c.Query("limit")

	var startTime, endTime time.Time
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			startTime = time.Now().AddDate(0, -1, 0)
		}
	} else {
		startTime = time.Now().AddDate(0, -1, 0)
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			endTime = time.Now()
		}
	} else {
		endTime = time.Now()
	}

	limit := 1000
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	data, err := models.GetHydrologyData(h.ctx, stationID, startTime, endTime, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"station_id": stationID,
		"count":      len(data),
		"data":       data,
	})
}

func (h *APIHandler) GetLatestHydrologyData(c *gin.Context) {
	stationID := c.Param("station_id")

	data, err := models.GetLatestHydrologyData(h.ctx, stationID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No data found"})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *APIHandler) GetAllLatestData(c *gin.Context) {
	data, err := models.GetAllLatestHydrologyData(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(data),
		"data":  data,
	})
}

func (h *APIHandler) GetDailyStats(c *gin.Context) {
	stationID := c.Param("station_id")
	startTimeStr := c.Query("start_time")
	endTimeStr := c.Query("end_time")

	startTime := time.Now().AddDate(0, -1, 0)
	endTime := time.Now()
	var err error

	if startTimeStr != "" {
		startTime, err = time.Parse(time.RFC3339, startTimeStr)
		if err != nil {
			startTime = time.Now().AddDate(0, -1, 0)
		}
	}

	if endTimeStr != "" {
		endTime, err = time.Parse(time.RFC3339, endTimeStr)
		if err != nil {
			endTime = time.Now()
		}
	}

	stats, err := models.GetDailyStats(h.ctx, stationID, startTime, endTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"station_id": stationID,
		"count":      len(stats),
		"data":       stats,
	})
}

func (h *APIHandler) GetMonitoringStations(c *gin.Context) {
	stations, err := models.GetMonitoringStations(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(stations),
		"data":  stations,
	})
}

func (h *APIHandler) GetWolongIron(c *gin.Context) {
	data, err := models.GetWolongIron(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(data),
		"data":  data,
	})
}

func (h *APIHandler) GetAlerts(c *gin.Context) {
	acknowledgedStr := c.Query("acknowledged")
	limitStr := c.Query("limit")

	var acknowledged *bool
	if acknowledgedStr != "" {
		ack, _ := strconv.ParseBool(acknowledgedStr)
		acknowledged = &ack
	}

	limit := 100
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	alerts, err := models.GetAlerts(h.ctx, acknowledged, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(alerts),
		"data":  alerts,
	})
}

func (h *APIHandler) AcknowledgeAlert(c *gin.Context) {
	alertID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid alert ID"})
		return
	}

	var req struct {
		AcknowledgedBy string `json:"acknowledged_by" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := models.AcknowledgeAlert(h.ctx, alertID, req.AcknowledgedBy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Alert acknowledged successfully"})
}

func (h *APIHandler) GetUnpublishedAlerts(c *gin.Context) {
	alerts, err := models.GetUnpublishedAlerts(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(alerts),
		"data":  alerts,
	})
}

type BedEvolutionRequest struct {
	Years int `json:"years"`
}

func (h *APIHandler) RunBedEvolutionPrediction(c *gin.Context) {
	stationID := c.Param("station_id")
	years, _ := strconv.Atoi(c.DefaultQuery("years", "10"))

	var req BedEvolutionRequest
	if err := c.ShouldBindJSON(&req); err == nil && req.Years > 0 {
		years = req.Years
	}

	results, err := simulation.PredictBedEvolution(h.ctx, stationID, years)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var baseElevation float64 = 726.5
	latestData, err := models.GetLatestHydrologyData(h.ctx, stationID)
	if err == nil && latestData.BedElevation > 0 {
		baseElevation = latestData.BedElevation
	}

	monthlyPredictions := make([]map[string]interface{}, 0, years*12)
	var avgAnnualDeposition, avgAnnualErosion float64

	for yearIdx, annualResult := range results {
		for month := 0; month < 12; month++ {
			monthFraction := float64(yearIdx) + float64(month)/12.0
			seasonalFactor := math.Sin(monthFraction*2*math.Pi/1.0) * 0.02
			
			monthlyDeposition := annualResult.Deposition / 12.0 * (1 + seasonalFactor)
			monthlyErosion := annualResult.Erosion / 12.0 * (1 - seasonalFactor)
			elevationChange := annualResult.PredictedElevation - baseElevation
			elevationChange += (float64(month) / 12.0) * (annualResult.NetChange)

			predDate := time.Now().AddDate(yearIdx, month, 0)
			
			monthlyPredictions = append(monthlyPredictions, map[string]interface{}{
				"prediction_date":      predDate,
				"bed_elevation_change": elevationChange,
				"predicted_elevation":  annualResult.PredictedElevation,
				"erosion_rate":         monthlyErosion,
				"deposition_rate":      monthlyDeposition,
				"sediment_accumulation": elevationChange * 1000,
				"confidence":           annualResult.Confidence,
			})

			avgAnnualDeposition += monthlyDeposition
			avgAnnualErosion += monthlyErosion
		}
	}

	avgAnnualDeposition = avgAnnualDeposition / float64(years) * 12
	avgAnnualErosion = avgAnnualErosion / float64(years) * 12
	finalElevation := results[len(results)-1].PredictedElevation
	
	riskLevel := "低"
	elevationDiff := finalElevation - baseElevation
	if elevationDiff > 0.3 {
		riskLevel = "高"
	} else if elevationDiff > 0.15 {
		riskLevel = "中"
	}

	c.JSON(http.StatusOK, gin.H{
		"station_id":               stationID,
		"years":                    years,
		"model":                    "Sediment Transport Model v1.0",
		"base_elevation":           baseElevation,
		"predictions":              monthlyPredictions,
		"average_annual_deposition": avgAnnualDeposition,
		"average_annual_erosion":    avgAnnualErosion,
		"final_elevation":          finalElevation,
		"risk_level":               riskLevel,
		"annual_data":              results,
	})
}

func (h *APIHandler) GetBedEvolutionPredictions(c *gin.Context) {
	stationID := c.Param("station_id")

	results, err := models.GetBedEvolutionPredictions(h.ctx, stationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"station_id": stationID,
		"count":      len(results),
		"data":       results,
	})
}

type BambooCageSimulationRequest struct {
	Location   string `json:"location" binding:"required"`
	CageCount  int    `json:"cage_count" binding:"required,min=1,max=100"`
	SimName    string `json:"simulation_name"`
	CreatedBy  string `json:"created_by"`
}

func (h *APIHandler) RunBambooCageSimulation(c *gin.Context) {
	var req BambooCageSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SimName == "" {
		req.SimName = "竹笼装石仿真 - " + time.Now().Format("2006-01-02 15:04:05")
	}

	simRecord := &models.AnnualRepairSimulation{
		SimulationName: req.SimName,
		SimulationType: "bamboo_cage",
		Status:         "running",
		CreatedBy:      req.CreatedBy,
	}

	simID, err := models.InsertAnnualRepairSimulation(h.ctx, simRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cages, err := simulation.SimulateBambooCagePlacement(h.ctx, simID, req.Location, req.CageCount)
	if err != nil {
		models.UpdateSimulationStatus(h.ctx, simID, "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := map[string]interface{}{
		"cages":      cages,
		"cage_count": len(cages),
	}
	models.UpdateSimulationResult(h.ctx, simID, result)

	c.JSON(http.StatusOK, gin.H{
		"simulation_id": simID,
		"location":      req.Location,
		"cage_count":    len(cages),
		"data":          cages,
	})
}

type MachaInterceptionRequest struct {
	Location          string  `json:"location" binding:"required"`
	MachaCount        int     `json:"macha_count" binding:"required,min=1,max=50"`
	InitialFlowRate   float64 `json:"initial_flow_rate" binding:"required,min=0"`
	InitialWaterLevel float64 `json:"initial_water_level" binding:"required,min=0"`
	SimName           string  `json:"simulation_name"`
	CreatedBy         string  `json:"created_by"`
}

func (h *APIHandler) RunMachaInterceptionSimulation(c *gin.Context) {
	var req MachaInterceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.SimName == "" {
		req.SimName = "杩槎截流仿真 - " + time.Now().Format("2006-01-02 15:04:05")
	}

	simRecord := &models.AnnualRepairSimulation{
		SimulationName: req.SimName,
		SimulationType: "macha_interception",
		Status:         "running",
		CreatedBy:      req.CreatedBy,
	}

	simID, err := models.InsertAnnualRepairSimulation(h.ctx, simRecord)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	machas, interceptionData, err := simulation.SimulateMachaInterception(
		h.ctx, simID, req.Location, req.MachaCount,
		req.InitialFlowRate, req.InitialWaterLevel,
	)
	if err != nil {
		models.UpdateSimulationStatus(h.ctx, simID, "failed")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result := map[string]interface{}{
		"machas":             machas,
		"interception_data":  interceptionData,
		"macha_count":        len(machas),
		"final_efficiency":   interceptionData[len(interceptionData)-1].InterceptionEfficiency,
		"final_flow_rate":    interceptionData[len(interceptionData)-1].FlowRateAfter,
	}
	models.UpdateSimulationResult(h.ctx, simID, result)

	c.JSON(http.StatusOK, gin.H{
		"simulation_id": simID,
		"location":      req.Location,
		"macha_count":   len(machas),
		"machas":        machas,
		"interception":  interceptionData,
	})
}

func (h *APIHandler) GetSimulations(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	simulations, err := models.GetAnnualRepairSimulations(h.ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(simulations),
		"data":  simulations,
	})
}

func (h *APIHandler) GetMachaSimulationData(c *gin.Context) {
	simID, err := strconv.Atoi(c.Param("simulation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid simulation ID"})
		return
	}

	data, err := models.GetMachaSimulationData(h.ctx, simID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulation_id": simID,
		"count":         len(data),
		"data":          data,
	})
}

func (h *APIHandler) GetBambooCageSimulationData(c *gin.Context) {
	simID, err := strconv.Atoi(c.Param("simulation_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid simulation ID"})
		return
	}

	data, err := models.GetBambooCageSimulationData(h.ctx, simID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulation_id": simID,
		"count":         len(data),
		"data":          data,
	})
}

func (h *APIHandler) GetAnnualRepairRecords(c *gin.Context) {
	records, err := models.GetAnnualRepairRecords(h.ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(records),
		"data":  records,
	})
}

func (h *APIHandler) GetDEMGrid(c *gin.Context) {
	centerX, _ := strconv.ParseFloat(c.DefaultQuery("center_x", "0"), 64)
	centerY, _ := strconv.ParseFloat(c.DefaultQuery("center_y", "0"), 64)
	gridSize, _ := strconv.ParseFloat(c.DefaultQuery("grid_size", "100"), 64)
	resolution, _ := strconv.ParseFloat(c.DefaultQuery("resolution", "5"), 64)
	baseElevation, _ := strconv.ParseFloat(c.DefaultQuery("base_elevation", "726.5"), 64)

	grid := simulation.GenerateDEMGrid(h.ctx, centerX, centerY, gridSize, resolution, baseElevation)

	c.JSON(http.StatusOK, gin.H{
		"grid_size":     gridSize,
		"resolution":    resolution,
		"base_elevation": baseElevation,
		"dimensions":    len(grid),
		"data":          grid,
	})
}

func (h *APIHandler) GetEvolutionRate(c *gin.Context) {
	stationID := c.Param("station_id")

	data, err := models.GetHydrologyData(
		h.ctx, stationID,
		time.Now().AddDate(0, -1, 0),
		time.Now(),
		8760,
	)
	if err != nil || len(data) < 2 {
		c.JSON(http.StatusOK, gin.H{
			"station_id": stationID,
			"annual_deposition": 0,
			"annual_erosion":    0,
			"annual_net_change": 0,
			"has_data":          false,
		})
		return
	}

	deposition, erosion, netChange := simulation.CalculateEvolutionRate(data)

	c.JSON(http.StatusOK, gin.H{
		"station_id":         stationID,
		"annual_deposition":  deposition,
		"annual_erosion":     erosion,
		"annual_net_change":  netChange,
		"has_data":           true,
		"data_points":        len(data),
	})
}

var wsClients = make(map[*websocket.Conn]bool)
var wsBroadcast = make(chan interface{})

func (h *APIHandler) RealTimeWebSocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	wsClients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(wsClients, ws)
			break
		}
	}
}

func broadcastHydrologyData(data *models.HydrologyData) {
	msg := map[string]interface{}{
		"type": "hydrology",
		"data": data,
	}
	wsBroadcast <- msg
}

func broadcastAlert(alert *models.Alert) {
	msg := map[string]interface{}{
		"type": "alert",
		"data": alert,
	}
	wsBroadcast <- msg
}

func StartWebSocketBroadcaster(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-wsBroadcast:
				for client := range wsClients {
					payload, _ := json.Marshal(msg)
					err := client.WriteMessage(websocket.TextMessage, payload)
					if err != nil {
						client.Close()
						delete(wsClients, client)
					}
				}
			}
		}
	}()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

func SetupStaticFiles(r *gin.Engine) {
	frontendPath := os.Getenv("FRONTEND_PATH")
	if frontendPath == "" {
		frontendPath = "./frontend"
	}
	r.Static("/frontend", frontendPath)
	r.GET("/", func(c *gin.Context) {
		c.File(frontendPath + "/index.html")
	})
}

func (h *APIHandler) GetDynastyTechniques(c *gin.Context) {
	dynasty := c.Query("dynasty")
	category := c.Query("category")

	data, err := models.GetDynastyTechniques(h.ctx, dynasty, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(data),
		"data":  data,
	})
}

func (h *APIHandler) GetModernComparisons(c *gin.Context) {
	category := c.Query("category")

	data, err := models.GetModernComparisons(h.ctx, category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(data),
		"data":  data,
	})
}

type EarthquakeSimulationRequest struct {
	SimulationName  string  `json:"simulation_name" binding:"required"`
	Magnitude       float64 `json:"magnitude" binding:"required,min=0,max=10"`
	EpicenterX      float64 `json:"epicenter_x" binding:"required"`
	EpicenterY      float64 `json:"epicenter_y" binding:"required"`
	EpicenterZ      float64 `json:"epicenter_z"`
	FocalMechanism  string  `json:"focal_mechanism"`
	DurationSeconds float64 `json:"duration_seconds" binding:"required,min=1"`
	CreatedBy       string  `json:"created_by"`
}

type StructureLocation struct {
	X float64
	Y float64
	Z float64
}

func (h *APIHandler) RunEarthquakeSimulation(c *gin.Context) {
	var req EarthquakeSimulationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var yuzuiDynamic, feishayanDynamic, baopingkouDynamic, renzidiDynamic float64

	if req.FocalMechanism == "" {
		req.FocalMechanism = "strike-slip"
	}

	pga := math.Pow(10, 0.5*req.Magnitude-3)

	yuzuiLoc := StructureLocation{X: 100, Y: 50, Z: 730}
	feishayanLoc := StructureLocation{X: 150, Y: 30, Z: 728}
	baopingkouLoc := StructureLocation{X: 200, Y: 60, Z: 725}
	renzidiLoc := StructureLocation{X: 80, Y: 80, Z: 727}

	// ============= 新增：Newmark-β 动力时程积分 & 并行计算 =============
	dynamicEnabled := req.Magnitude >= 6.0 // 6级以上启用动力分析
	parallelWorkers := 4                   // 4个结构并行计算

	if dynamicEnabled {
		type dynResult struct {
			name         string
			dispMax      float64
			accMax       float64
			baseShear    float64
			damageDynamic float64
		}

		structures := []struct {
			name string
			loc  StructureLocation
			mass float64 // 结构质量(t)
			stiff float64 // 刚度(N/m)
			damp  float64 // 阻尼比
		}{
			{"yuzui",     yuzuiLoc,     2.5e6, 1.2e9, 0.05},
			{"feishayan", feishayanLoc, 1.8e6, 8.5e8, 0.05},
			{"baopingkou", baopingkouLoc, 3.0e6, 1.5e9, 0.05},
			{"renzidi",   renzidiLoc,   1.2e6, 6.0e8, 0.05},
		}

		// Newmark-β 方法参数
		beta := 0.25
		gamma := 0.5
		dt := 0.005 // 5ms时间步
		nSteps := int(req.DurationSeconds / dt)

		// 生成人工地震波（作为激励加速度时程）
		earthquakeWave := make([]float64, nSteps)
		for i := 0; i < nSteps; i++ {
			t := float64(i) * dt
			freq := 2.0 + 3.0*rand.Float64() // 2-5Hz
			envelope := 1.0
			if t < 2.0 {
				envelope = t / 2.0
			} else if t > req.DurationSeconds-2.0 {
				envelope = (req.DurationSeconds - t) / 2.0
			} else if t > req.DurationSeconds {
				envelope = 0
			}
			earthquakeWave[i] = pga * 9.81 * envelope * (0.6*math.Sin(2*math.Pi*freq*t) + 0.4*rand.NormFloat64())
		}

		// goroutine 并行计算4个结构的动力响应
		results := make(chan dynResult, parallelWorkers)
		var wg sync.WaitGroup
		wg.Add(parallelWorkers)

		for _, s := range structures {
			go func(name string, loc StructureLocation, mass, stiff, damp float64) {
				defer wg.Done()
				// Newmark-β 积分
				disp := make([]float64, nSteps)
				vel := make([]float64, nSteps)
				acc := make([]float64, nSteps)
				omega := math.Sqrt(stiff / mass) // 自振圆频率
				kPrime := stiff + gamma/(beta*dt)*damp + mass/(beta*dt*dt)

				disp[0] = 0
				vel[0] = 0
				acc[0] = earthquakeWave[0]

				maxDisp := 0.0
				maxAcc := 0.0

				for i := 1; i < nSteps; i++ {
					// 预测位移速度
					predDisp := disp[i-1] + dt*vel[i-1] + 0.5*dt*dt*(1-2*beta)*acc[i-1]
					predVel := vel[i-1] + dt*(1-gamma)*acc[i-1]
					predAcc := earthquakeWave[i]

					// 修正
					deltaDisp := (predAcc - damp*predVel - stiff*predDisp) / kPrime
					disp[i] = predDisp + deltaDisp
					vel[i] = predVel + gamma/beta*deltaDisp/dt
					acc[i] = predAcc + (gamma/beta - 1)*deltaDisp/(dt*dt)

					if math.Abs(disp[i]) > maxDisp {
						maxDisp = math.Abs(disp[i])
					}
					if math.Abs(acc[i]) > maxAcc {
						maxAcc = math.Abs(acc[i])
					}
				}

				// 计算结构损伤
				dispLimit := stiff / (omega * omega * mass) * 0.1 // 允许位移的10%
				damageDynamic := math.Min(1.0, maxDisp/dispLimit)

				results <- dynResult{
					name:         name,
					dispMax:      maxDisp,
					accMax:       maxAcc,
					baseShear:    maxAcc * mass,
					damageDynamic: damageDynamic,
				}
			}(s.name, s.loc, s.mass, s.stiff, s.damp)
		}

		wg.Wait()
		close(results)

		// 收集并行计算结果
		for r := range results {
			switch r.name {
			case "yuzui":
				yuzuiDynamic = r.damageDynamic
			case "feishayan":
				feishayanDynamic = r.damageDynamic
			case "baopingkou":
				baopingkouDynamic = r.damageDynamic
			case "renzidi":
				renzidiDynamic = r.damageDynamic
			}
		}
	}
	// ============= 新增结束 =============

	calcDistance := func(loc StructureLocation) float64 {
		dx := loc.X - req.EpicenterX
		dy := loc.Y - req.EpicenterY
		dz := loc.Z - req.EpicenterZ
		return math.Sqrt(dx*dx + dy*dy + dz*dz)
	}

	calcAttenuation := func(distance float64) float64 {
		return math.Exp(-0.001 * distance)
	}

	yuzuiDist := calcDistance(yuzuiLoc)
	feishayanDist := calcDistance(feishayanLoc)
	baopingkouDist := calcDistance(baopingkouLoc)
	renzidiDist := calcDistance(renzidiLoc)

	yuzuiDamage := pga * calcAttenuation(yuzuiDist)
	feishayanDamage := pga * calcAttenuation(feishayanDist)
	baopingkouDamage := pga * calcAttenuation(baopingkouDist)
	renzidiDamage := pga * calcAttenuation(renzidiDist)

	// 合并静态和动力损伤（取较大值）
	if yuzuiDynamic > 0 {
		yuzuiDamage = math.Max(yuzuiDamage, yuzuiDynamic)
	}
	if feishayanDynamic > 0 {
		feishayanDamage = math.Max(feishayanDamage, feishayanDynamic)
	}
	if baopingkouDynamic > 0 {
		baopingkouDamage = math.Max(baopingkouDamage, baopingkouDynamic)
	}
	if renzidiDynamic > 0 {
		renzidiDamage = math.Max(renzidiDamage, renzidiDynamic)
	}

	maxDamage := math.Max(math.Max(yuzuiDamage, feishayanDamage), math.Max(baopingkouDamage, renzidiDamage))
	safetyAssessment := "safe"
	if maxDamage > 0.6 {
		safetyAssessment = "danger"
	} else if maxDamage > 0.3 {
		safetyAssessment = "caution"
	}

	timeStep := 0.01
	numSteps := int(req.DurationSeconds / timeStep)
	timeSeries := make([]map[string]interface{}, 0, numSteps)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numSteps; i++ {
		t := float64(i) * timeStep
		amplitude := pga * math.Exp(-0.5*t) * (1 + 0.3*r.Float64())
		acceleration := amplitude * math.Sin(2*math.Pi*2*t) * (1 + 0.2*r.Float64())
		timeSeries = append(timeSeries, map[string]interface{}{
			"time":         t,
			"acceleration": acceleration,
			"displacement": acceleration * 0.1,
		})
	}
	timeSeriesJSON, _ := json.Marshal(timeSeries)

	structureDamage := map[string]interface{}{
		"yuzui": map[string]interface{}{
			"damage_level":  yuzuiDamage,
			"cracks":        int(yuzuiDamage * 10),
			"settlement_mm": yuzuiDamage * 50,
		},
		"feishayan": map[string]interface{}{
			"damage_level":  feishayanDamage,
			"cracks":        int(feishayanDamage * 10),
			"settlement_mm": feishayanDamage * 50,
		},
		"baopingkou": map[string]interface{}{
			"damage_level":  baopingkouDamage,
			"cracks":        int(baopingkouDamage * 10),
			"settlement_mm": baopingkouDamage * 50,
		},
		"renzidi": map[string]interface{}{
			"damage_level":  renzidiDamage,
			"cracks":        int(renzidiDamage * 10),
			"settlement_mm": renzidiDamage * 50,
		},
	}
	structureDamageJSON, _ := json.Marshal(structureDamage)

	bankCollapse := make([]map[string]interface{}, 0)
	if maxDamage > 0.3 {
		collapseCount := int(maxDamage * 20)
		for i := 0; i < collapseCount; i++ {
			bankCollapse = append(bankCollapse, map[string]interface{}{
				"position_x":  100 + r.Float64()*200,
				"position_y":  20 + r.Float64()*80,
				"volume_m3":   10 + r.Float64()*100,
				"collapse_at": r.Float64() * req.DurationSeconds,
			})
		}
	}
	bankCollapseJSON, _ := json.Marshal(bankCollapse)

	sedimentDisturbance := pga * 0.5
	flowPathChange := pga * 0.3
	waterDiversionChange := pga * 0.2

	sim := &models.EarthquakeSimulation{
		SimulationName:       req.SimulationName,
		Magnitude:            req.Magnitude,
		EpicenterX:           req.EpicenterX,
		EpicenterY:           req.EpicenterY,
		EpicenterZ:           req.EpicenterZ,
		FocalMechanism:       req.FocalMechanism,
		PGA:                  pga,
		DurationSeconds:      req.DurationSeconds,
		TimeSeriesData:       string(timeSeriesJSON),
		StructureDamage:      string(structureDamageJSON),
		YuzuiDamage:          yuzuiDamage,
		FeishayanDamage:      feishayanDamage,
		BaopingkouDamage:     baopingkouDamage,
		RenzidiDamage:        renzidiDamage,
		BankCollapse:         string(bankCollapseJSON),
		SedimentDisturbance:  sedimentDisturbance,
		FlowPathChange:       flowPathChange,
		WaterDiversionChange: waterDiversionChange,
		SafetyAssessment:     safetyAssessment,
		CreatedBy:            req.CreatedBy,
	}

	id, err := models.InsertEarthquakeSimulation(h.ctx, sim)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sim.ID = id

	c.JSON(http.StatusOK, gin.H{
		"id":                     id,
		"simulation_name":        sim.SimulationName,
		"magnitude":              sim.Magnitude,
		"pga":                    sim.PGA,
		"epicenter":              map[string]float64{"x": sim.EpicenterX, "y": sim.EpicenterY, "z": sim.EpicenterZ},
		"focal_mechanism":        sim.FocalMechanism,
		"duration_seconds":       sim.DurationSeconds,
		"structure_damage":       structureDamage,
		"yuzui_damage":           sim.YuzuiDamage,
		"feishayan_damage":       sim.FeishayanDamage,
		"baopingkou_damage":      sim.BaopingkouDamage,
		"renzidi_damage":         sim.RenzidiDamage,
		"bank_collapse":          bankCollapse,
		"sediment_disturbance":   sim.SedimentDisturbance,
		"flow_path_change":       sim.FlowPathChange,
		"water_diversion_change": sim.WaterDiversionChange,
		"safety_assessment":      sim.SafetyAssessment,
		"time_series_data":       timeSeries,
		"created_by":             sim.CreatedBy,
		"dynamic_analysis_enabled": dynamicEnabled,
		"parallel_workers":       parallelWorkers,
	})
}

func (h *APIHandler) GetEarthquakeSimulations(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 50
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	data, err := models.GetEarthquakeSimulations(h.ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(data),
		"data":  data,
	})
}

type UserOperationRequest struct {
	SessionID             string  `json:"session_id" binding:"required"`
	UserNickname          string  `json:"user_nickname"`
	OperationType         string  `json:"operation_type" binding:"required"`
	ObjectType            string  `json:"object_type"`
	PositionX             float64 `json:"position_x"`
	PositionY             float64 `json:"position_y"`
	PositionZ             float64 `json:"position_z"`
	RotationAngle         float64 `json:"rotation_angle"`
	ObjectParams          string  `json:"object_params"`
	OperationOrder        int     `json:"operation_order" binding:"required"`
	InterceptionEfficiency float64 `json:"interception_efficiency"`
	StabilityScore        float64 `json:"stability_score"`
	DredgingVolume        float64 `json:"dredging_volume"`
}

func (h *APIHandler) InsertUserOperation(c *gin.Context) {
	var req UserOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	op := &models.UserRepairOperation{
		SessionID:             req.SessionID,
		UserNickname:          req.UserNickname,
		OperationType:         req.OperationType,
		ObjectType:            req.ObjectType,
		PositionX:             req.PositionX,
		PositionY:             req.PositionY,
		PositionZ:             req.PositionZ,
		RotationAngle:         req.RotationAngle,
		ObjectParams:          req.ObjectParams,
		OperationOrder:        req.OperationOrder,
		InterceptionEfficiency: req.InterceptionEfficiency,
		StabilityScore:        req.StabilityScore,
		DredgingVolume:        req.DredgingVolume,
		CompletionStatus:      "in_progress",
	}

	id, err := models.InsertUserOperation(h.ctx, op)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	op.ID = id

	c.JSON(http.StatusOK, gin.H{
		"message": "Operation recorded successfully",
		"data":    op,
	})
}

func (h *APIHandler) GetUserOperations(c *gin.Context) {
	sessionID := c.Param("session_id")

	data, err := models.GetUserOperations(h.ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"session_id": sessionID,
		"count":      len(data),
		"data":       data,
	})
}

type FinishSessionRequest struct {
	SessionID        string  `json:"session_id" binding:"required"`
	TotalScore       float64 `json:"total_score" binding:"required"`
	CompletionStatus string  `json:"completion_status" binding:"required"`
	Achievement      string  `json:"achievement"`
	DurationSeconds  int     `json:"duration_seconds" binding:"required,min=0"`
}

func (h *APIHandler) FinishUserSession(c *gin.Context) {
	var req FinishSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.UpdateUserSessionScore(h.ctx, req.SessionID, req.TotalScore, req.CompletionStatus, req.Achievement, req.DurationSeconds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Session finished successfully",
		"session_id":       req.SessionID,
		"total_score":      req.TotalScore,
		"completion_status": req.CompletionStatus,
		"achievement":      req.Achievement,
		"duration_seconds": req.DurationSeconds,
	})
}

func (h *APIHandler) GetUserRanking(c *gin.Context) {
	limitStr := c.Query("limit")
	limit := 100
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	data, err := models.GetUserScoreRanking(h.ctx, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": len(data),
		"data":  data,
	})
}
