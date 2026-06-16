package models

import (
	"time"
)

type HydrologyData struct {
	Time                  time.Time `json:"time" db:"time"`
	StationID             string    `json:"station_id" db:"station_id"`
	StationName           string    `json:"station_name" db:"station_name"`
	WaterLevel            float64   `json:"water_level" db:"water_level"`
	FlowRate              float64   `json:"flow_rate" db:"flow_rate"`
	SedimentConcentration float64   `json:"sediment_concentration" db:"sediment_concentration"`
	BedElevation          float64   `json:"bed_elevation" db:"bed_elevation"`
	Temperature           float64   `json:"temperature" db:"temperature"`
	Rainfall              float64   `json:"rainfall" db:"rainfall"`
	SensorStatus          int       `json:"sensor_status" db:"sensor_status"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type WolongIron struct {
	ID            int       `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Location      string    `json:"location" db:"location"`
	Elevation     float64   `json:"elevation" db:"elevation"`
	Description   string    `json:"description" db:"description"`
	InstalledYear int       `json:"installed_year" db:"installed_year"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type MonitoringStation struct {
	ID              int       `json:"id" db:"id"`
	StationID       string    `json:"station_id" db:"station_id"`
	Name            string    `json:"name" db:"name"`
	LocationLat     float64   `json:"location_lat" db:"location_lat"`
	LocationLng     float64   `json:"location_lng" db:"location_lng"`
	ReachName       string    `json:"reach_name" db:"reach_name"`
	BedrockElevation float64  `json:"bedrock_elevation" db:"bedrock_elevation"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}

type Alert struct {
	ID                  int        `json:"id" db:"id"`
	AlertTime           time.Time  `json:"alert_time" db:"alert_time"`
	AlertType           string     `json:"alert_type" db:"alert_type"`
	AlertLevel          string     `json:"alert_level" db:"alert_level"`
	StationID           string     `json:"station_id" db:"station_id"`
	Message             string     `json:"message" db:"message"`
	BedElevation        float64    `json:"bed_elevation" db:"bed_elevation"`
	WolongIronElevation float64    `json:"wolong_iron_elevation" db:"wolong_iron_elevation"`
	ExceededValue       float64    `json:"exceeded_value" db:"exceeded_value"`
	Acknowledged        bool       `json:"acknowledged" db:"acknowledged"`
	AcknowledgedAt      *time.Time `json:"acknowledged_at,omitempty" db:"acknowledged_at"`
	AcknowledgedBy      string     `json:"acknowledged_by,omitempty" db:"acknowledged_by"`
	MqttPublished       bool       `json:"mqtt_published" db:"mqtt_published"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
}

type BedEvolutionPrediction struct {
	ID                        int       `json:"id" db:"id"`
	StationID                 string    `json:"station_id" db:"station_id"`
	PredictionDate            time.Time `json:"prediction_date" db:"prediction_date"`
	ForecastHorizonMonths     int       `json:"forecast_horizon_months" db:"forecast_horizon_months"`
	PredictedBedElevation     float64   `json:"predicted_bed_elevation" db:"predicted_bed_elevation"`
	PredictedSedimentDeposition float64 `json:"predicted_sediment_deposition" db:"predicted_sediment_deposition"`
	PredictedErosion          float64   `json:"predicted_erosion" db:"predicted_erosion"`
	ModelVersion              string    `json:"model_version" db:"model_version"`
	Confidence                float64   `json:"confidence" db:"confidence"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
}

type AnnualRepairSimulation struct {
	ID             int64           `json:"id" db:"id"`
	SimulationName string          `json:"simulation_name" db:"simulation_name"`
	SimulationType string          `json:"simulation_type" db:"simulation_type"`
	StartTime      *time.Time      `json:"start_time,omitempty" db:"start_time"`
	EndTime        *time.Time      `json:"end_time,omitempty" db:"end_time"`
	Status         string          `json:"status" db:"status"`
	Parameters     SimulationParams `json:"parameters" db:"parameters"`
	Result         interface{}     `json:"result" db:"result"`
	CreatedBy      string          `json:"created_by" db:"created_by"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

type SimulationParams struct {
	Location     string    `json:"location"`
	StartTime    time.Time `json:"start_time"`
	DurationDays int       `json:"duration_days"`
	GridSize     float64   `json:"grid_size"`
	TimeStep     int       `json:"time_step"`
	Parameters   map[string]interface{} `json:"parameters,omitempty"`
}

type MachaInterceptionData struct {
	Time                  time.Time `json:"time" db:"time"`
	SimulationID          int       `json:"simulation_id" db:"simulation_id"`
	PositionX             float64   `json:"position_x" db:"position_x"`
	PositionY             float64   `json:"position_y" db:"position_y"`
	WaterLevelBefore      float64   `json:"water_level_before" db:"water_level_before"`
	WaterLevelAfter       float64   `json:"water_level_after" db:"water_level_after"`
	FlowRateBefore        float64   `json:"flow_rate_before" db:"flow_rate_before"`
	FlowRateAfter         float64   `json:"flow_rate_after" db:"flow_rate_after"`
	InterceptionEfficiency float64  `json:"interception_efficiency" db:"interception_efficiency"`
	MachaCount            int       `json:"macha_count" db:"macha_count"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type BambooCageData struct {
	ID                    int       `json:"id" db:"id"`
	SimulationID          int       `json:"simulation_id" db:"simulation_id"`
	CageID                string    `json:"cage_id" db:"cage_id"`
	PositionX             float64   `json:"position_x" db:"position_x"`
	PositionY             float64   `json:"position_y" db:"position_y"`
	PositionZ             float64   `json:"position_z" db:"position_z"`
	StoneCount            int       `json:"stone_count" db:"stone_count"`
	CageDiameter          float64   `json:"cage_diameter" db:"cage_diameter"`
	CageLength            float64   `json:"cage_length" db:"cage_length"`
	Porosity              float64   `json:"porosity" db:"porosity"`
	StabilityCoefficient  float64   `json:"stability_coefficient" db:"stability_coefficient"`
	DepositionHeight      float64   `json:"deposition_height" db:"deposition_height"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type AnnualRepairRecord struct {
	ID                 int       `json:"id" db:"id"`
	RepairYear         int       `json:"repair_year" db:"repair_year"`
	StartDate          string    `json:"start_date" db:"start_date"`
	EndDate            string    `json:"end_date" db:"end_date"`
	Location           string    `json:"location" db:"location"`
	RepairType         string    `json:"repair_type" db:"repair_type"`
	BambooCageCount    int       `json:"bamboo_cage_count" db:"bamboo_cage_count"`
	MachaCount         int       `json:"macha_count" db:"macha_count"`
	DredgingVolume     float64   `json:"dredging_volume" db:"dredging_volume"`
	BedElevationBefore float64   `json:"bed_elevation_before" db:"bed_elevation_before"`
	BedElevationAfter  float64   `json:"bed_elevation_after" db:"bed_elevation_after"`
	Notes              string    `json:"notes" db:"notes"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
}

type DailyStats struct {
	Bucket          time.Time `json:"bucket" db:"bucket"`
	StationID       string    `json:"station_id" db:"station_id"`
	AvgWaterLevel   float64   `json:"avg_water_level" db:"avg_water_level"`
	MaxWaterLevel   float64   `json:"max_water_level" db:"max_water_level"`
	MinWaterLevel   float64   `json:"min_water_level" db:"min_water_level"`
	AvgFlowRate     float64   `json:"avg_flow_rate" db:"avg_flow_rate"`
	MaxFlowRate     float64   `json:"max_flow_rate" db:"max_flow_rate"`
	AvgSediment     float64   `json:"avg_sediment" db:"avg_sediment"`
	MaxSediment     float64   `json:"max_sediment" db:"max_sediment"`
	AvgBedElevation float64   `json:"avg_bed_elevation" db:"avg_bed_elevation"`
	RecordCount     int       `json:"record_count" db:"record_count"`
}

type DEMGrid struct {
	GridX     int     `json:"grid_x"`
	GridY     int     `json:"grid_y"`
	Elevation float64 `json:"elevation"`
	WaterDepth float64 `json:"water_depth"`
}

type DynastyRepairTechnique struct {
	ID               int       `json:"id" db:"id"`
	Dynasty          string    `json:"dynasty" db:"dynasty"`
	DynastyName      string    `json:"dynasty_name" db:"dynasty_name"`
	EraYears         string    `json:"era_years" db:"era_years"`
	TechniqueName    string    `json:"technique_name" db:"technique_name"`
	Category         string    `json:"category" db:"category"`
	Materials        string    `json:"materials" db:"materials"`
	Tools            string    `json:"tools" db:"tools"`
	Procedures       string    `json:"procedures" db:"procedures"`
	MachaParams      string    `json:"macha_params,omitempty" db:"macha_params"`
	BambooParams     string    `json:"bamboo_params,omitempty" db:"bamboo_params"`
	DredgingParams   string    `json:"dredging_params,omitempty" db:"dredging_params"`
	EfficiencyScore  float64   `json:"efficiency_score" db:"efficiency_score"`
	LaborCost        int       `json:"labor_cost" db:"labor_cost"`
	DurationDays     int       `json:"duration_days" db:"duration_days"`
	CostSilverLiang  float64   `json:"cost_silver_liang" db:"cost_silver_liang"`
	HistoricalNotes     string    `json:"historical_notes" db:"historical_notes"`
	ReferenceSources    string    `json:"reference_sources" db:"reference_sources"`
	ArchaeologicalSource string   `json:"source_archaeology,omitempty" db:"-"`
	UncertaintyRange    string    `json:"uncertainty_range,omitempty" db:"-"`
	ExperimentalMethod  string    `json:"experimental_method,omitempty" db:"-"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

type ModernRepairComparison struct {
	ID                  int       `json:"id" db:"id"`
	ComparisonItem      string    `json:"comparison_item" db:"comparison_item"`
	Category            string    `json:"category" db:"category"`
	AncientMethod       string    `json:"ancient_method" db:"ancient_method"`
	ModernMethod        string    `json:"modern_method" db:"modern_method"`
	AncientEfficiency   float64   `json:"ancient_efficiency" db:"ancient_efficiency"`
	ModernEfficiency    float64   `json:"modern_efficiency" db:"modern_efficiency"`
	AncientCost         float64   `json:"ancient_cost" db:"ancient_cost"`
	ModernCost          float64   `json:"modern_cost" db:"modern_cost"`
	AncientEnvImpact    float64   `json:"ancient_env_impact" db:"ancient_env_impact"`
	ModernEnvImpact     float64   `json:"modern_env_impact" db:"modern_env_impact"`
	AncientEquipment    string    `json:"ancient_equipment" db:"ancient_equipment"`
	ModernEquipment     string    `json:"modern_equipment" db:"modern_equipment"`
	LaborAncient        int       `json:"labor_ancient" db:"labor_ancient"`
	LaborModern         int       `json:"labor_modern" db:"labor_modern"`
	DurationAncientDays int       `json:"duration_ancient_days" db:"duration_ancient_days"`
	DurationModernDays  int       `json:"duration_modern_days" db:"duration_modern_days"`
	Unit                string    `json:"unit" db:"unit"`
	StandardReference   string    `json:"standard_reference" db:"standard_reference"`
	StandardCode        string    `json:"standard_code" db:"standard_code"`
	Description         string    `json:"description" db:"description"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
}

type EarthquakeSimulation struct {
	ID                    int       `json:"id" db:"id"`
	SimulationName        string    `json:"simulation_name" db:"simulation_name"`
	Magnitude             float64   `json:"magnitude" db:"magnitude"`
	EpicenterX            float64   `json:"epicenter_x" db:"epicenter_x"`
	EpicenterY            float64   `json:"epicenter_y" db:"epicenter_y"`
	EpicenterZ            float64   `json:"epicenter_z" db:"epicenter_z"`
	FocalMechanism        string    `json:"focal_mechanism" db:"focal_mechanism"`
	PGA                   float64   `json:"pga" db:"pga"`
	DurationSeconds       float64   `json:"duration_seconds" db:"duration_seconds"`
	TimeSeriesData        string    `json:"time_series_data,omitempty" db:"time_series_data"`
	StructureDamage       string    `json:"structure_damage,omitempty" db:"structure_damage"`
	YuzuiDamage           float64   `json:"yuzui_damage" db:"yuzui_damage"`
	FeishayanDamage       float64   `json:"feishayan_damage" db:"feishayan_damage"`
	BaopingkouDamage      float64   `json:"baopingkou_damage" db:"baopingkou_damage"`
	RenzidiDamage         float64   `json:"renzidi_damage" db:"renzidi_damage"`
	BankCollapse          string    `json:"bank_collapse,omitempty" db:"bank_collapse"`
	SedimentDisturbance   float64   `json:"sediment_disturbance" db:"sediment_disturbance"`
	FlowPathChange        float64   `json:"flow_path_change" db:"flow_path_change"`
	WaterDiversionChange  float64   `json:"water_diversion_change" db:"water_diversion_change"`
	SafetyAssessment      string    `json:"safety_assessment" db:"safety_assessment"`
	CreatedBy             string    `json:"created_by,omitempty" db:"created_by"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}

type UserRepairOperation struct {
	ID                    int       `json:"id" db:"id"`
	SessionID             string    `json:"session_id" db:"session_id"`
	UserNickname          string    `json:"user_nickname,omitempty" db:"user_nickname"`
	OperationType         string    `json:"operation_type" db:"operation_type"`
	ObjectType            string    `json:"object_type,omitempty" db:"object_type"`
	PositionX             float64   `json:"position_x,omitempty" db:"position_x"`
	PositionY             float64   `json:"position_y,omitempty" db:"position_y"`
	PositionZ             float64   `json:"position_z,omitempty" db:"position_z"`
	RotationAngle         float64   `json:"rotation_angle,omitempty" db:"rotation_angle"`
	ObjectParams          string    `json:"object_params,omitempty" db:"object_params"`
	OperationOrder        int       `json:"operation_order" db:"operation_order"`
	SimulationResult      string    `json:"simulation_result,omitempty" db:"simulation_result"`
	InterceptionEfficiency float64  `json:"interception_efficiency" db:"interception_efficiency"`
	StabilityScore        float64   `json:"stability_score" db:"stability_score"`
	DredgingVolume        float64   `json:"dredging_volume" db:"dredging_volume"`
	CompletionStatus      string    `json:"completion_status" db:"completion_status"`
	TotalScore            float64   `json:"total_score" db:"total_score"`
	Achievement           string    `json:"achievement,omitempty" db:"achievement"`
	DurationSeconds       int       `json:"duration_seconds" db:"duration_seconds"`
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
}
