package era_comparator

import (
	"context"
	"testing"

	"dujiangyan-system/pkg/models"
)

func TestReductionRate(t *testing.T) {
	c := NewEraComparator(context.Background())

	t.Run("正常成本减少率计算", func(t *testing.T) {
		result := c.ReductionRate(100, 20)
		expected := 80.0
		if result != expected {
			t.Errorf("ReductionRate(100, 20) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("ancient为0时除零保护返回0", func(t *testing.T) {
		result := c.ReductionRate(0, 50)
		if result != 0 {
			t.Errorf("ReductionRate(0, 50) = %f, 期望 0", result)
		}
	})

	t.Run("modern等于ancient时减少率为0", func(t *testing.T) {
		result := c.ReductionRate(100, 100)
		if result != 0 {
			t.Errorf("ReductionRate(100, 100) = %f, 期望 0", result)
		}
	})

	t.Run("modern大于ancient时减少率为负", func(t *testing.T) {
		result := c.ReductionRate(50, 100)
		if result >= 0 {
			t.Errorf("ReductionRate(50, 100) = %f, 应为负数", result)
		}
	})
}

func TestCostEfficiencyRatio(t *testing.T) {
	c := NewEraComparator(context.Background())

	t.Run("正常成本效率比计算", func(t *testing.T) {
		result := c.CostEfficiencyRatio(50, 100)
		expected := 2.0
		if result != expected {
			t.Errorf("CostEfficiencyRatio(50, 100) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("cost为0时返回0", func(t *testing.T) {
		result := c.CostEfficiencyRatio(0, 100)
		if result != 0 {
			t.Errorf("CostEfficiencyRatio(0, 100) = %f, 期望 0", result)
		}
	})

	t.Run("efficiency为0时返回0", func(t *testing.T) {
		result := c.CostEfficiencyRatio(50, 0)
		if result != 0 {
			t.Errorf("CostEfficiencyRatio(50, 0) = %f, 期望 0", result)
		}
	})
}

func TestOverallComparison(t *testing.T) {
	c := NewEraComparator(context.Background())

	t.Run("空数组返回零值", func(t *testing.T) {
		result := c.OverallComparison([]models.ModernRepairComparison{})

		avgAncientCost := result["avg_ancient_cost"].(float64)
		avgModernCost := result["avg_modern_cost"].(float64)
		avgCostReduction := result["avg_cost_reduction"].(float64)
		avgEfficiencyImprovement := result["avg_efficiency_improvement"].(float64)

		if avgAncientCost != 0 {
			t.Errorf("avg_ancient_cost = %f, 期望 0", avgAncientCost)
		}
		if avgModernCost != 0 {
			t.Errorf("avg_modern_cost = %f, 期望 0", avgModernCost)
		}
		if avgCostReduction != 0 {
			t.Errorf("avg_cost_reduction = %f, 期望 0", avgCostReduction)
		}
		if avgEfficiencyImprovement != 0 {
			t.Errorf("avg_efficiency_improvement = %f, 期望 0", avgEfficiencyImprovement)
		}
	})

	t.Run("单条数据正常计算", func(t *testing.T) {
		comparisons := []models.ModernRepairComparison{
			{
				AncientCost:      100,
				ModernCost:       20,
				AncientEfficiency: 50,
				ModernEfficiency:  100,
			},
		}

		result := c.OverallComparison(comparisons)

		avgAncientCost := result["avg_ancient_cost"].(float64)
		avgModernCost := result["avg_modern_cost"].(float64)

		if avgAncientCost != 100 {
			t.Errorf("avg_ancient_cost = %f, 期望 100", avgAncientCost)
		}
		if avgModernCost != 20 {
			t.Errorf("avg_modern_cost = %f, 期望 20", avgModernCost)
		}
	})

	t.Run("多条数据平均值计算", func(t *testing.T) {
		comparisons := []models.ModernRepairComparison{
			{AncientCost: 100, ModernCost: 20},
			{AncientCost: 200, ModernCost: 40},
		}

		result := c.OverallComparison(comparisons)

		avgAncientCost := result["avg_ancient_cost"].(float64)
		avgModernCost := result["avg_modern_cost"].(float64)

		if avgAncientCost != 150 {
			t.Errorf("avg_ancient_cost = %f, 期望 150", avgAncientCost)
		}
		if avgModernCost != 30 {
			t.Errorf("avg_modern_cost = %f, 期望 30", avgModernCost)
		}
	})
}

func TestValidateStandardRef(t *testing.T) {
	c := NewEraComparator(context.Background())

	t.Run("全部有单位标准规范", func(t *testing.T) {
		comparisons := []models.ModernRepairComparison{
			{Unit: "m³", StandardReference: "GB50201-2014", StandardCode: "GB50201"},
			{Unit: "t", StandardReference: "SL252-2000", StandardCode: "SL252"},
		}

		withUnit, withStandard, withCode := c.ValidateStandardRef(comparisons)

		if withUnit != 2 {
			t.Errorf("withUnit = %d, 期望 2", withUnit)
		}
		if withStandard != 2 {
			t.Errorf("withStandard = %d, 期望 2", withStandard)
		}
		if withCode != 2 {
			t.Errorf("withCode = %d, 期望 2", withCode)
		}
	})

	t.Run("空数组返回0", func(t *testing.T) {
		withUnit, withStandard, withCode := c.ValidateStandardRef([]models.ModernRepairComparison{})

		if withUnit != 0 {
			t.Errorf("withUnit = %d, 期望 0", withUnit)
		}
		if withStandard != 0 {
			t.Errorf("withStandard = %d, 期望 0", withStandard)
		}
		if withCode != 0 {
			t.Errorf("withCode = %d, 期望 0", withCode)
		}
	})

	t.Run("部分有字段", func(t *testing.T) {
		comparisons := []models.ModernRepairComparison{
			{Unit: "m³", StandardReference: "", StandardCode: ""},
			{Unit: "", StandardReference: "GB50201", StandardCode: ""},
			{Unit: "", StandardReference: "", StandardCode: "SL252"},
		}

		withUnit, withStandard, withCode := c.ValidateStandardRef(comparisons)

		if withUnit != 1 {
			t.Errorf("withUnit = %d, 期望 1", withUnit)
		}
		if withStandard != 1 {
			t.Errorf("withStandard = %d, 期望 1", withStandard)
		}
		if withCode != 1 {
			t.Errorf("withCode = %d, 期望 1", withCode)
		}
	})
}
