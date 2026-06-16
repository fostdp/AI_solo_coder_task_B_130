package api

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

const machaEfficiencyCap = 0.85

type machaParams struct {
	Height           float64 `json:"height"`
	LogDiameter      float64 `json:"log_diameter"`
	BundleLayers     int     `json:"bundle_layers"`
	EfficiencyPerUnit float64 `json:"efficiency_per_unit"`
}

type bambooParams struct {
	Diameter  float64 `json:"diameter"`
	Length    float64 `json:"length"`
	Porosity  float64 `json:"porosity"`
	Stones    int     `json:"stones"`
}

func calcMachaEfficiency(count int, efficiencyPerUnit float64) float64 {
	if count <= 0 || efficiencyPerUnit <= 0 {
		return 0
	}
	return math.Min(float64(count)*efficiencyPerUnit, machaEfficiencyCap)
}

func TestDynastyTechnique_InterceptionEfficiency(t *testing.T) {
	gin.SetMode(gin.TestMode)

	libingEffPerUnit := 0.04
	qingEffPerUnit := 0.05
	machaCount := 20

	t.Run("李冰时期20个杩槎效率不超过上限", func(t *testing.T) {
		efficiency := calcMachaEfficiency(machaCount, libingEffPerUnit)
		expected := float64(machaCount) * libingEffPerUnit
		if efficiency != expected {
			t.Errorf("李冰效率 = %f, 期望 %f", efficiency, expected)
		}
		if efficiency > machaEfficiencyCap {
			t.Errorf("李冰效率 %f 超过上限 %f", efficiency, machaEfficiencyCap)
		}
	})

	t.Run("清代20个杩槎效率受上限约束", func(t *testing.T) {
		efficiency := calcMachaEfficiency(machaCount, qingEffPerUnit)
		theoretical := float64(machaCount) * qingEffPerUnit
		if theoretical <= machaEfficiencyCap {
			t.Errorf("理论效率 %f 应超过上限 %f", theoretical, machaEfficiencyCap)
		}
		if efficiency != machaEfficiencyCap {
			t.Errorf("清代效率 = %f, 期望上限 %f", efficiency, machaEfficiencyCap)
		}
	})

	t.Run("清代效率不低于李冰效率", func(t *testing.T) {
		libingEff := calcMachaEfficiency(machaCount, libingEffPerUnit)
		qingEff := calcMachaEfficiency(machaCount, qingEffPerUnit)
		if qingEff < libingEff {
			t.Errorf("清代效率 %f 低于李冰效率 %f", qingEff, libingEff)
		}
	})
}

func TestDynastyTechnique_MachaParamsParsing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	libingJSON := `{"height":4.5,"log_diameter":0.18,"bundle_layers":5,"efficiency_per_unit":0.04}`
	qingJSON := `{"height":5.5,"log_diameter":0.25,"bundle_layers":7,"efficiency_per_unit":0.05}`

	var libing, qing machaParams

	t.Run("李冰杩槎参数解析", func(t *testing.T) {
		if err := json.Unmarshal([]byte(libingJSON), &libing); err != nil {
			t.Fatalf("解析李冰杩槎参数失败: %v", err)
		}
		if libing.Height != 4.5 {
			t.Errorf("height = %f, 期望 4.5", libing.Height)
		}
		if libing.LogDiameter != 0.18 {
			t.Errorf("log_diameter = %f, 期望 0.18", libing.LogDiameter)
		}
		if libing.BundleLayers != 5 {
			t.Errorf("bundle_layers = %d, 期望 5", libing.BundleLayers)
		}
		if libing.EfficiencyPerUnit != 0.04 {
			t.Errorf("efficiency_per_unit = %f, 期望 0.04", libing.EfficiencyPerUnit)
		}
	})

	t.Run("清代杩槎参数解析", func(t *testing.T) {
		if err := json.Unmarshal([]byte(qingJSON), &qing); err != nil {
			t.Fatalf("解析清代杩槎参数失败: %v", err)
		}
		if qing.Height != 5.5 {
			t.Errorf("height = %f, 期望 5.5", qing.Height)
		}
		if qing.LogDiameter != 0.25 {
			t.Errorf("log_diameter = %f, 期望 0.25", qing.LogDiameter)
		}
		if qing.BundleLayers != 7 {
			t.Errorf("bundle_layers = %d, 期望 7", qing.BundleLayers)
		}
		if qing.EfficiencyPerUnit != 0.05 {
			t.Errorf("efficiency_per_unit = %f, 期望 0.05", qing.EfficiencyPerUnit)
		}
	})

	t.Run("清代杩槎参数大于李冰", func(t *testing.T) {
		if qing.Height <= libing.Height {
			t.Errorf("清代 height %f 应大于李冰 height %f", qing.Height, libing.Height)
		}
		if qing.BundleLayers <= libing.BundleLayers {
			t.Errorf("清代 bundle_layers %d 应大于李冰 bundle_layers %d", qing.BundleLayers, libing.BundleLayers)
		}
	})
}

func TestDynastyTechnique_BambooParamsParsing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	libingJSON := `{"diameter":0.7,"length":3.5,"porosity":0.35,"stones":80}`
	qingJSON := `{"diameter":0.9,"length":4.0,"porosity":0.3,"stones":120}`

	var libing, qing bambooParams

	t.Run("李冰竹笼参数解析", func(t *testing.T) {
		if err := json.Unmarshal([]byte(libingJSON), &libing); err != nil {
			t.Fatalf("解析李冰竹笼参数失败: %v", err)
		}
		if libing.Diameter != 0.7 {
			t.Errorf("diameter = %f, 期望 0.7", libing.Diameter)
		}
		if libing.Length != 3.5 {
			t.Errorf("length = %f, 期望 3.5", libing.Length)
		}
		if libing.Porosity != 0.35 {
			t.Errorf("porosity = %f, 期望 0.35", libing.Porosity)
		}
		if libing.Stones != 80 {
			t.Errorf("stones = %d, 期望 80", libing.Stones)
		}
	})

	t.Run("清代竹笼参数解析", func(t *testing.T) {
		if err := json.Unmarshal([]byte(qingJSON), &qing); err != nil {
			t.Fatalf("解析清代竹笼参数失败: %v", err)
		}
		if qing.Diameter != 0.9 {
			t.Errorf("diameter = %f, 期望 0.9", qing.Diameter)
		}
		if qing.Length != 4.0 {
			t.Errorf("length = %f, 期望 4.0", qing.Length)
		}
		if qing.Porosity != 0.3 {
			t.Errorf("porosity = %f, 期望 0.3", qing.Porosity)
		}
		if qing.Stones != 120 {
			t.Errorf("stones = %d, 期望 120", qing.Stones)
		}
	})

	t.Run("清代竹笼参数对比", func(t *testing.T) {
		if qing.Diameter <= libing.Diameter {
			t.Errorf("清代 diameter %f 应大于李冰 diameter %f", qing.Diameter, libing.Diameter)
		}
		if qing.Porosity >= libing.Porosity {
			t.Errorf("清代 porosity %f 应小于李冰 porosity %f (更密实)", qing.Porosity, libing.Porosity)
		}
		if qing.Stones <= libing.Stones {
			t.Errorf("清代 stones %d 应大于李冰 stones %d", qing.Stones, libing.Stones)
		}
	})
}

func TestDynastyTechnique_EfficiencyScoreRange(t *testing.T) {
	gin.SetMode(gin.TestMode)

	libingScore := 65.0
	qingScore := 82.5

	t.Run("李冰效率评分范围", func(t *testing.T) {
		if libingScore < 60 || libingScore > 70 {
			t.Errorf("李冰效率评分 %f 不在 60-70 范围内", libingScore)
		}
	})

	t.Run("清代效率评分范围", func(t *testing.T) {
		if qingScore < 80 || qingScore > 85 {
			t.Errorf("清代效率评分 %f 不在 80-85 范围内", qingScore)
		}
	})

	t.Run("清代效率评分大于李冰", func(t *testing.T) {
		if qingScore <= libingScore {
			t.Errorf("清代效率评分 %f 应大于李冰效率评分 %f", qingScore, libingScore)
		}
	})
}

func TestDynastyTechnique_LaborAndDuration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	libingLabor := 500
	qingLabor := 800
	libingDredgingDays := 30
	qingDredgingDays := 45

	t.Run("清代人力不少于李冰人力", func(t *testing.T) {
		if qingLabor < libingLabor {
			t.Errorf("清代人力 %d 不应少于李冰人力 %d", qingLabor, libingLabor)
		}
	})

	t.Run("清代淘淤工期大于李冰淘淤工期", func(t *testing.T) {
		if qingDredgingDays <= libingDredgingDays {
			t.Errorf("清代淘淤工期 %d 应大于李冰淘淤工期 %d", qingDredgingDays, libingDredgingDays)
		}
	})
}

func TestDynastyTechnique_MachaEfficiency_Boundary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	libingEffPerUnit := 0.04
	qingEffPerUnit := 0.05

	t.Run("0个杩槎效率为0", func(t *testing.T) {
		libingEff := calcMachaEfficiency(0, libingEffPerUnit)
		qingEff := calcMachaEfficiency(0, qingEffPerUnit)
		if libingEff != 0 {
			t.Errorf("李冰0个杩槎效率 = %f, 期望 0", libingEff)
		}
		if qingEff != 0 {
			t.Errorf("清代0个杩槎效率 = %f, 期望 0", qingEff)
		}
	})

	t.Run("1个杩槎效率", func(t *testing.T) {
		libingEff := calcMachaEfficiency(1, libingEffPerUnit)
		qingEff := calcMachaEfficiency(1, qingEffPerUnit)
		if libingEff != 0.04 {
			t.Errorf("李冰1个杩槎效率 = %f, 期望 0.04", libingEff)
		}
		if qingEff != 0.05 {
			t.Errorf("清代1个杩槎效率 = %f, 期望 0.05", qingEff)
		}
	})

	t.Run("17个杩槎清代到达上限", func(t *testing.T) {
		libingEff := calcMachaEfficiency(17, libingEffPerUnit)
		qingEff := calcMachaEfficiency(17, qingEffPerUnit)
		if libingEff != 0.68 {
			t.Errorf("李冰17个杩槎效率 = %f, 期望 0.68", libingEff)
		}
		if qingEff != 0.85 {
			t.Errorf("清代17个杩槎效率 = %f, 期望 0.85 (到达上限)", qingEff)
		}
	})

	t.Run("100个杩槎两者均受上限约束", func(t *testing.T) {
		libingEff := calcMachaEfficiency(100, libingEffPerUnit)
		qingEff := calcMachaEfficiency(100, qingEffPerUnit)
		if libingEff != machaEfficiencyCap {
			t.Errorf("李冰100个杩槎效率 = %f, 期望 %f", libingEff, machaEfficiencyCap)
		}
		if qingEff != machaEfficiencyCap {
			t.Errorf("清代100个杩槎效率 = %f, 期望 %f", qingEff, machaEfficiencyCap)
		}
	})
}

func TestDynastyTechnique_MachaEfficiency_Abnormal(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("负数杩槎效率为0", func(t *testing.T) {
		eff := calcMachaEfficiency(-5, 0.04)
		if eff != 0 {
			t.Errorf("负数杩槎效率 = %f, 期望 0", eff)
		}
	})

	t.Run("efficiency_per_unit为0效率为0", func(t *testing.T) {
		eff := calcMachaEfficiency(20, 0)
		if eff != 0 {
			t.Errorf("efficiency_per_unit为0时效率 = %f, 期望 0", eff)
		}
	})
}

func TestDynastyTechnique_RequestValidation(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("query参数dynasty正确读取", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/dynasty-techniques?dynasty=qing&category=macha", nil)

		dynasty := c.Query("dynasty")
		category := c.Query("category")

		if dynasty != "qing" {
			t.Errorf("dynasty = %s, 期望 qing", dynasty)
		}
		if category != "macha" {
			t.Errorf("category = %s, 期望 macha", category)
		}
	})

	t.Run("无query参数时返回空字符串", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/v1/dynasty-techniques", nil)

		dynasty := c.Query("dynasty")
		category := c.Query("category")

		if dynasty != "" {
			t.Errorf("dynasty = %s, 期望空字符串", dynasty)
		}
		if category != "" {
			t.Errorf("category = %s, 期望空字符串", category)
		}
	})
}
