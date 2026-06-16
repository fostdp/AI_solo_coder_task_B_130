package craft_comparator

import (
	"context"
	"testing"
)

func TestCalcMachaInterception(t *testing.T) {
	c := NewCraftComparator(context.Background())

	t.Run("0个杩槎截流效率为0", func(t *testing.T) {
		result := c.CalcMachaInterception(0, 0.05)
		if result != 0 {
			t.Errorf("CalcMachaInterception(0, 0.05) = %f, 期望 0", result)
		}
	})

	t.Run("17个杩槎清代效率到达上限0.85", func(t *testing.T) {
		result := c.CalcMachaInterception(17, 0.05)
		if result != 0.85 {
			t.Errorf("CalcMachaInterception(17, 0.05) = %f, 期望 0.85", result)
		}
	})

	t.Run("10个杩槎效率正常计算", func(t *testing.T) {
		result := c.CalcMachaInterception(10, 0.05)
		expected := 10 * 0.05
		if result != expected {
			t.Errorf("CalcMachaInterception(10, 0.05) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("效率不超过上限0.85", func(t *testing.T) {
		result := c.CalcMachaInterception(100, 0.05)
		if result > 0.85 {
			t.Errorf("CalcMachaInterception(100, 0.05) = %f, 超过上限 0.85", result)
		}
	})
}

func TestCalcBambooStability(t *testing.T) {
	c := NewCraftComparator(context.Background())

	t.Run("0个竹笼稳定性为基础值0.7", func(t *testing.T) {
		result := c.CalcBambooStability(0, 0.7)
		if result != 0.7 {
			t.Errorf("CalcBambooStability(0, 0.7) = %f, 期望 0.7", result)
		}
	})

	t.Run("10个竹笼稳定性到达上限1.0", func(t *testing.T) {
		result := c.CalcBambooStability(10, 0.7)
		if result != 1.0 {
			t.Errorf("CalcBambooStability(10, 0.7) = %f, 期望 1.0", result)
		}
	})

	t.Run("5个竹笼稳定性正常计算", func(t *testing.T) {
		result := c.CalcBambooStability(5, 0.7)
		expected := 0.7 + float64(5)*0.03
		if result != expected {
			t.Errorf("CalcBambooStability(5, 0.7) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("稳定性不超过上限1.0", func(t *testing.T) {
		result := c.CalcBambooStability(50, 0.7)
		if result > 1.0 {
			t.Errorf("CalcBambooStability(50, 0.7) = %f, 超过上限 1.0", result)
		}
	})
}

func TestCompareEfficiency(t *testing.T) {
	c := NewCraftComparator(context.Background())

	t.Run("正常效率对比", func(t *testing.T) {
		result := c.CompareEfficiency(20, 20, 0.05, 0.04)

		qingEff := result["qing_efficiency"]
		libingEff := result["libing_efficiency"]
		improvementRatio := result["improvement_ratio"]

		expectedQing := 20.0 * 0.05
		expectedLibing := 20.0 * 0.04

		if qingEff != expectedQing {
			t.Errorf("qing_efficiency = %f, 期望 %f", qingEff, expectedQing)
		}
		if libingEff != expectedLibing {
			t.Errorf("libing_efficiency = %f, 期望 %f", libingEff, expectedLibing)
		}

		expectedRatio := expectedQing / expectedLibing
		if improvementRatio != expectedRatio {
			t.Errorf("improvement_ratio = %f, 期望 %f", improvementRatio, expectedRatio)
		}
	})

	t.Run("李冰效率为0时improvement_ratio为0", func(t *testing.T) {
		result := c.CompareEfficiency(20, 0, 0.05, 0.04)
		improvementRatio := result["improvement_ratio"]
		if improvementRatio != 0 {
			t.Errorf("improvement_ratio = %f, 期望 0", improvementRatio)
		}
	})

	t.Run("清代和李冰数量都为0", func(t *testing.T) {
		result := c.CompareEfficiency(0, 0, 0.05, 0.04)
		qingEff := result["qing_efficiency"]
		libingEff := result["libing_efficiency"]
		improvementRatio := result["improvement_ratio"]

		if qingEff != 0 {
			t.Errorf("qing_efficiency = %f, 期望 0", qingEff)
		}
		if libingEff != 0 {
			t.Errorf("libing_efficiency = %f, 期望 0", libingEff)
		}
		if improvementRatio != 0 {
			t.Errorf("improvement_ratio = %f, 期望 0", improvementRatio)
		}
	})
}

func TestValidateArchaeologicalParams(t *testing.T) {
	c := NewCraftComparator(context.Background())

	t.Run("包含三个字段", func(t *testing.T) {
		params := `{"source_archaeology":"site_A","uncertainty_range":"±5%","experimental_method":"radiocarbon"}`
		hasSource, hasUncertainty, hasMethod := c.ValidateArchaeologicalParams(params)

		if !hasSource {
			t.Error("hasSource 应为 true")
		}
		if !hasUncertainty {
			t.Error("hasUncertainty 应为 true")
		}
		if !hasMethod {
			t.Error("hasMethod 应为 true")
		}
	})

	t.Run("只包含source_archaeology", func(t *testing.T) {
		params := `{"source_archaeology":"site_A"}`
		hasSource, hasUncertainty, hasMethod := c.ValidateArchaeologicalParams(params)

		if !hasSource {
			t.Error("hasSource 应为 true")
		}
		if hasUncertainty {
			t.Error("hasUncertainty 应为 false")
		}
		if hasMethod {
			t.Error("hasMethod 应为 false")
		}
	})

	t.Run("无效JSON返回全为false", func(t *testing.T) {
		params := `invalid json`
		hasSource, hasUncertainty, hasMethod := c.ValidateArchaeologicalParams(params)

		if hasSource {
			t.Error("hasSource 应为 false")
		}
		if hasUncertainty {
			t.Error("hasUncertainty 应为 false")
		}
		if hasMethod {
			t.Error("hasMethod 应为 false")
		}
	})

	t.Run("空JSON返回全为false", func(t *testing.T) {
		params := `{}`
		hasSource, hasUncertainty, hasMethod := c.ValidateArchaeologicalParams(params)

		if hasSource {
			t.Error("hasSource 应为 false")
		}
		if hasUncertainty {
			t.Error("hasUncertainty 应为 false")
		}
		if hasMethod {
			t.Error("hasMethod 应为 false")
		}
	})
}
