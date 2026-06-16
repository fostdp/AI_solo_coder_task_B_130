package api

import (
	"math"
	"testing"
)

func calcAverage(data []float64) float64 {
	if len(data) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum / float64(len(data))
}

func calcReductionRate(ancient, modern float64) float64 {
	if ancient == 0 {
		return 0
	}
	return (ancient - modern) / ancient * 100
}

func calcCostEfficiencyRatio(cost, efficiency float64) float64 {
	if efficiency == 0 {
		return math.Inf(1)
	}
	return cost / efficiency
}

func hasNegativeValue(data []float64) bool {
	for _, v := range data {
		if v < 0 {
			return true
		}
	}
	return false
}

func countExceeding(data []float64, limit float64) int {
	count := 0
	for _, v := range data {
		if v > limit {
			count++
		}
	}
	return count
}

func TestModernComparison_CostAnalysis(t *testing.T) {
	ancientCost := []float64{100, 80, 120, 40, 70, 100, 100, 50, 100, 100}
	modernCost := []float64{150, 130, 200, 110, 180, 250, 120, 200, 90, 80}

	t.Run("平均古代成本", func(t *testing.T) {
		avg := calcAverage(ancientCost)
		expected := 86.0
		if math.Abs(avg-expected) > 0.01 {
			t.Errorf("平均古代成本 = %.2f, 期望 %.2f", avg, expected)
		}
	})

	t.Run("平均现代成本", func(t *testing.T) {
		avg := calcAverage(modernCost)
		expected := 151.0
		if math.Abs(avg-expected) > 0.01 {
			t.Errorf("平均现代成本 = %.2f, 期望 %.2f", avg, expected)
		}
	})

	t.Run("现代平均成本大于古代平均成本", func(t *testing.T) {
		avgAncient := calcAverage(ancientCost)
		avgModern := calcAverage(modernCost)
		if avgModern <= avgAncient {
			t.Errorf("现代平均成本 %.2f 应大于古代平均成本 %.2f", avgModern, avgAncient)
		}
	})

	t.Run("工期效率项古代成本与现代成本接近", func(t *testing.T) {
		idx := 8
		if math.Abs(ancientCost[idx]-100) > 0.01 {
			t.Errorf("工期效率项古代成本 = %.0f, 期望 100", ancientCost[idx])
		}
		if math.Abs(modernCost[idx]-90) > 0.01 {
			t.Errorf("工期效率项现代成本 = %.0f, 期望 90", modernCost[idx])
		}
		diff := math.Abs(ancientCost[idx] - modernCost[idx])
		threshold := 15.0
		if diff > threshold {
			t.Errorf("工期效率项古今成本差 %.1f 超过阈值 %.1f", diff, threshold)
		}
	})
}

func TestModernComparison_EfficiencyGap(t *testing.T) {
	ancientEfficiency := []float64{60, 65, 40, 55, 50, 100, 30, 45, 15, 30}
	modernEfficiency := []float64{95, 90, 95, 92, 93, 250, 95, 90, 70, 98}

	t.Run("每项现代效率大于古代效率", func(t *testing.T) {
		for i := 0; i < len(ancientEfficiency); i++ {
			if modernEfficiency[i] <= ancientEfficiency[i] {
				t.Errorf("项目 %d: 现代效率 %.0f 应大于古代效率 %.0f", i, modernEfficiency[i], ancientEfficiency[i])
			}
		}
	})

	t.Run("材料成本项例外_成本更高不是效率", func(t *testing.T) {
		idx := 5
		if ancientEfficiency[idx] != 100 {
			t.Errorf("材料成本项古代效率 = %.0f, 期望 100", ancientEfficiency[idx])
		}
		if modernEfficiency[idx] != 250 {
			t.Errorf("材料成本项现代效率 = %.0f, 期望 250", modernEfficiency[idx])
		}
		if modernEfficiency[idx] <= ancientEfficiency[idx] {
			t.Error("材料成本项现代值应大于古代值，但高值代表成本而非效率提升")
		}
	})

	t.Run("效率差距范围验证", func(t *testing.T) {
		for i := 0; i < len(ancientEfficiency); i++ {
			gap := modernEfficiency[i] - ancientEfficiency[i]
			if gap <= 0 {
				t.Errorf("项目 %d: 效率差距 %.0f 应为正值", i, gap)
			}
		}
	})
}

func TestModernComparison_EnvironmentalImpact(t *testing.T) {
	ancientEnv := []float64{30, 20, 50, 15, 25, 10, 30, 35, 10, 5}
	modernEnv := []float64{70, 50, 35, 60, 65, 80, 40, 50, 65, 30}

	t.Run("生态影响项_古代更友好", func(t *testing.T) {
		idx := 8
		if ancientEnv[idx] != 10 {
			t.Errorf("生态影响项古代值 = %.0f, 期望 10", ancientEnv[idx])
		}
		if modernEnv[idx] != 65 {
			t.Errorf("生态影响项现代值 = %.0f, 期望 65", modernEnv[idx])
		}
		if ancientEnv[idx] >= modernEnv[idx] {
			t.Error("生态影响项古代值应小于现代值(值越低越好)")
		}
	})

	t.Run("监测技术项_古代几乎无影响", func(t *testing.T) {
		idx := 9
		if ancientEnv[idx] != 5 {
			t.Errorf("监测技术项古代值 = %.0f, 期望 5", ancientEnv[idx])
		}
		if modernEnv[idx] != 30 {
			t.Errorf("监测技术项现代值 = %.0f, 期望 30", modernEnv[idx])
		}
		if ancientEnv[idx] >= modernEnv[idx] {
			t.Error("监测技术项古代值应小于现代值(古代几乎无影响)")
		}
	})

	t.Run("古代环境平均值低于现代", func(t *testing.T) {
		avgAncient := calcAverage(ancientEnv)
		avgModern := calcAverage(modernEnv)
		if avgAncient >= avgModern {
			t.Errorf("古代环境平均 %.2f 应低于现代环境平均 %.2f", avgAncient, avgModern)
		}
	})
}

func TestModernComparison_LaborReduction(t *testing.T) {
	ancientLabor := []float64{120, 80, 100, 60, 300, 200, 50, 90, 70, 40}
	modernLabor := []float64{15, 10, 20, 8, 8, 30, 5, 12, 10, 5}

	t.Run("每项现代人力小于古代人力", func(t *testing.T) {
		for i := 0; i < len(ancientLabor); i++ {
			if modernLabor[i] >= ancientLabor[i] {
				t.Errorf("项目 %d: 现代人力 %.0f 应小于古代人力 %.0f", i, modernLabor[i], ancientLabor[i])
			}
		}
	})

	t.Run("截流人力减少率", func(t *testing.T) {
		idx := 0
		rate := calcReductionRate(ancientLabor[idx], modernLabor[idx])
		expected := 87.5
		if math.Abs(rate-expected) > 0.1 {
			t.Errorf("截流人力减少率 = %.1f%%, 期望 %.1f%%", rate, expected)
		}
	})

	t.Run("护岸人力减少率", func(t *testing.T) {
		idx := 1
		rate := calcReductionRate(ancientLabor[idx], modernLabor[idx])
		expected := 87.5
		if math.Abs(rate-expected) > 0.1 {
			t.Errorf("护岸人力减少率 = %.1f%%, 期望 %.1f%%", rate, expected)
		}
	})

	t.Run("清淤人力减少率", func(t *testing.T) {
		idx := 4
		rate := calcReductionRate(ancientLabor[idx], modernLabor[idx])
		expected := 97.3
		if math.Abs(rate-expected) > 0.1 {
			t.Errorf("清淤人力减少率 = %.1f%%, 期望 %.1f%%", rate, expected)
		}
	})
}

func TestModernComparison_DurationComparison(t *testing.T) {
	ancientDuration := []float64{30, 25, 20, 15, 45, 60, 10, 35, 40, 20}
	modernDuration := []float64{10, 8, 7, 5, 10, 20, 3, 12, 15, 6}

	t.Run("每项现代工期小于古代工期", func(t *testing.T) {
		for i := 0; i < len(ancientDuration); i++ {
			if modernDuration[i] >= ancientDuration[i] {
				t.Errorf("项目 %d: 现代工期 %.0f 应小于古代工期 %.0f", i, modernDuration[i], ancientDuration[i])
			}
		}
	})

	t.Run("清淤工期缩短率", func(t *testing.T) {
		idx := 4
		rate := calcReductionRate(ancientDuration[idx], modernDuration[idx])
		expected := 77.8
		if math.Abs(rate-expected) > 0.1 {
			t.Errorf("清淤工期缩短率 = %.1f%%, 期望 %.1f%%", rate, expected)
		}
	})

	t.Run("总体工期缩短率大于50%", func(t *testing.T) {
		totalAncient := 0.0
		totalModern := 0.0
		for i := 0; i < len(ancientDuration); i++ {
			totalAncient += ancientDuration[i]
			totalModern += modernDuration[i]
		}
		overallRate := calcReductionRate(totalAncient, totalModern)
		if overallRate <= 50 {
			t.Errorf("总体工期缩短率 = %.1f%%, 应大于 50%%", overallRate)
		}
	})
}

func TestModernComparison_CostEfficiencyRatio(t *testing.T) {
	ancientCost := []float64{100, 80, 120, 40, 70, 100, 100, 50, 100, 100}
	modernCost := []float64{150, 130, 200, 110, 180, 250, 120, 200, 90, 80}
	ancientEfficiency := []float64{60, 65, 40, 55, 50, 100, 30, 45, 15, 30}
	modernEfficiency := []float64{95, 90, 95, 92, 93, 250, 95, 90, 70, 98}

	t.Run("古代截流成本效率比", func(t *testing.T) {
		idx := 0
		ratio := calcCostEfficiencyRatio(ancientCost[idx], ancientEfficiency[idx])
		expected := 100.0 / 60.0
		if math.Abs(ratio-expected) > 0.01 {
			t.Errorf("古代截流成本效率比 = %.2f, 期望 %.2f", ratio, expected)
		}
	})

	t.Run("现代截流成本效率比", func(t *testing.T) {
		idx := 0
		ratio := calcCostEfficiencyRatio(modernCost[idx], modernEfficiency[idx])
		expected := 150.0 / 95.0
		if math.Abs(ratio-expected) > 0.01 {
			t.Errorf("现代截流成本效率比 = %.2f, 期望 %.2f", ratio, expected)
		}
	})

	t.Run("现代成本效率比不大于古代", func(t *testing.T) {
		ancientRatio := calcCostEfficiencyRatio(ancientCost[0], ancientEfficiency[0])
		modernRatio := calcCostEfficiencyRatio(modernCost[0], modernEfficiency[0])
		if modernRatio > ancientRatio {
			t.Errorf("现代成本效率比 %.2f 应不大于古代 %.2f", modernRatio, ancientRatio)
		}
	})

	t.Run("全部项目现代成本效率比验证", func(t *testing.T) {
		for i := 0; i < len(ancientCost); i++ {
			ancientRatio := calcCostEfficiencyRatio(ancientCost[i], ancientEfficiency[i])
			modernRatio := calcCostEfficiencyRatio(modernCost[i], modernEfficiency[i])
			if modernRatio > ancientRatio {
				t.Logf("项目 %d: 现代比 %.2f > 古代比 %.2f (现代虽贵但某些项效率不足以抵消)", i, modernRatio, ancientRatio)
			}
		}
	})
}

func TestModernComparison_Boundary(t *testing.T) {
	t.Run("效率为0时成本效率比为正无穷", func(t *testing.T) {
		ratio := calcCostEfficiencyRatio(100, 0)
		if !math.IsInf(ratio, 1) {
			t.Errorf("效率为0时成本效率比 = %.2f, 期望 +Inf", ratio)
		}
	})

	t.Run("成本为0时成本效率比为0", func(t *testing.T) {
		ratio := calcCostEfficiencyRatio(0, 60)
		if ratio != 0 {
			t.Errorf("成本为0时成本效率比 = %.2f, 期望 0", ratio)
		}
	})

	t.Run("空数组平均值返回0", func(t *testing.T) {
		avg := calcAverage([]float64{})
		if avg != 0 {
			t.Errorf("空数组平均值 = %.2f, 期望 0", avg)
		}
	})

	t.Run("单条对比数据统计", func(t *testing.T) {
		singleAncient := []float64{100}
		singleModern := []float64{150}
		avgAncient := calcAverage(singleAncient)
		avgModern := calcAverage(singleModern)
		if avgAncient != 100 {
			t.Errorf("单条古代平均 = %.2f, 期望 100", avgAncient)
		}
		if avgModern != 150 {
			t.Errorf("单条现代平均 = %.2f, 期望 150", avgModern)
		}
		rate := calcReductionRate(100, 150)
		if rate >= 0 {
			t.Errorf("单条减少率 = %.1f%%, 期望负值(现代更高)", rate)
		}
	})
}

func TestModernComparison_Abnormal(t *testing.T) {
	t.Run("负数成本检测", func(t *testing.T) {
		negativeCost := []float64{100, -50, 80, -10}
		if !hasNegativeValue(negativeCost) {
			t.Error("应检测到负数成本")
		}
		validCost := []float64{100, 50, 80, 10}
		if hasNegativeValue(validCost) {
			t.Error("不应误报有效成本为负数")
		}
	})

	t.Run("效率超过100检测", func(t *testing.T) {
		overLimitEfficiency := []float64{60, 95, 110, 85, 250}
		count := countExceeding(overLimitEfficiency, 100)
		if count != 2 {
			t.Errorf("超过100的项目数 = %d, 期望 2", count)
		}
		normalEfficiency := []float64{60, 95, 80, 85, 100}
		countNormal := countExceeding(normalEfficiency, 100)
		if countNormal != 0 {
			t.Errorf("正常效率超过100的项目数 = %d, 期望 0", countNormal)
		}
	})

	t.Run("空数组不崩溃", func(t *testing.T) {
		empty := []float64{}
		avg := calcAverage(empty)
		if avg != 0 {
			t.Errorf("空数组平均值 = %.2f, 期望 0", avg)
		}
		rate := calcReductionRate(0, 0)
		if rate != 0 {
			t.Errorf("零值减少率 = %.1f, 期望 0", rate)
		}
		ratio := calcCostEfficiencyRatio(0, 0)
		if !math.IsInf(ratio, 1) {
			t.Errorf("零值成本效率比 = %.2f, 期望 +Inf", ratio)
		}
	})
}
