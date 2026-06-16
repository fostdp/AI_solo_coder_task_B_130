package earthquake_simulator

import (
	"context"
	"math"
	"testing"
)

func TestCalcPGA(t *testing.T) {
	s := NewEarthquakeSimulator(context.Background(), 4)

	t.Run("M8.0地震PGA为10.0", func(t *testing.T) {
		result := s.CalcPGA(8.0)
		expected := 10.0
		if math.Abs(result-expected) > 1e-9 {
			t.Errorf("CalcPGA(8.0) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("M6.0地震PGA为1.0", func(t *testing.T) {
		result := s.CalcPGA(6.0)
		expected := 1.0
		if math.Abs(result-expected) > 1e-9 {
			t.Errorf("CalcPGA(6.0) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("M4.0地震PGA为0.1", func(t *testing.T) {
		result := s.CalcPGA(4.0)
		expected := 0.1
		if math.Abs(result-expected) > 1e-9 {
			t.Errorf("CalcPGA(4.0) = %f, 期望 %f", result, expected)
		}
	})

	t.Run("震级越大PGA越大", func(t *testing.T) {
		pga5 := s.CalcPGA(5.0)
		pga7 := s.CalcPGA(7.0)
		pga9 := s.CalcPGA(9.0)

		if pga5 >= pga7 {
			t.Errorf("M5.0 PGA %f 应小于 M7.0 PGA %f", pga5, pga7)
		}
		if pga7 >= pga9 {
			t.Errorf("M7.0 PGA %f 应小于 M9.0 PGA %f", pga7, pga9)
		}
	})
}

func TestCalcDamage(t *testing.T) {
	s := NewEarthquakeSimulator(context.Background(), 4)

	t.Run("近处破坏大于远处破坏", func(t *testing.T) {
		nearLoc := StructureLocation{X: 0, Y: 0, Z: 0}
		farLoc := StructureLocation{X: 1000, Y: 1000, Z: 0}
		pga := 1.0

		nearDamage := s.CalcDamage(pga, nearLoc, 0, 0)
		farDamage := s.CalcDamage(pga, farLoc, 0, 0)

		if nearDamage <= farDamage {
			t.Errorf("近处破坏 %f 应大于远处破坏 %f", nearDamage, farDamage)
		}
	})

	t.Run("破坏值在0到1之间", func(t *testing.T) {
		loc := StructureLocation{X: 100, Y: 100, Z: 0}
		pga := 100.0

		damage := s.CalcDamage(pga, loc, 0, 0)

		if damage < 0 {
			t.Errorf("破坏值 %f 不能小于 0", damage)
		}
		if damage > 1.0 {
			t.Errorf("破坏值 %f 不能大于 1.0", damage)
		}
	})

	t.Run("PGA为0时破坏为0", func(t *testing.T) {
		loc := StructureLocation{X: 100, Y: 100, Z: 0}
		damage := s.CalcDamage(0, loc, 0, 0)

		if damage != 0 {
			t.Errorf("PGA为0时破坏 = %f, 期望 0", damage)
		}
	})
}

func TestAssessSafety(t *testing.T) {
	s := NewEarthquakeSimulator(context.Background(), 4)

	t.Run("破坏大于0.6返回danger", func(t *testing.T) {
		result := s.AssessSafety(0.7, 0.5, 0.4, 0.3)
		if result != "danger" {
			t.Errorf("AssessSafety(0.7, 0.5, 0.4, 0.3) = %s, 期望 danger", result)
		}
	})

	t.Run("破坏0.45返回caution", func(t *testing.T) {
		result := s.AssessSafety(0.45, 0.3, 0.2, 0.1)
		if result != "caution" {
			t.Errorf("AssessSafety(0.45, 0.3, 0.2, 0.1) = %s, 期望 caution", result)
		}
	})

	t.Run("破坏0.2返回safe", func(t *testing.T) {
		result := s.AssessSafety(0.2, 0.1, 0.15, 0.05)
		if result != "safe" {
			t.Errorf("AssessSafety(0.2, 0.1, 0.15, 0.05) = %s, 期望 safe", result)
		}
	})

	t.Run("最大破坏决定安全等级", func(t *testing.T) {
		result := s.AssessSafety(0.1, 0.1, 0.1, 0.7)
		if result != "danger" {
			t.Errorf("最大破坏0.7时 = %s, 期望 danger", result)
		}
	})
}

func TestCalcSecondaryEffects(t *testing.T) {
	s := NewEarthquakeSimulator(context.Background(), 4)

	t.Run("sediment为maxDamage的0.5倍", func(t *testing.T) {
		sediment, flow, diversion := s.CalcSecondaryEffects(1.0)

		expectedSediment := 0.5
		if math.Abs(sediment-expectedSediment) > 1e-9 {
			t.Errorf("sediment = %f, 期望 %f", sediment, expectedSediment)
		}

		expectedFlow := 0.3
		if math.Abs(flow-expectedFlow) > 1e-9 {
			t.Errorf("flow = %f, 期望 %f", flow, expectedFlow)
		}

		expectedDiversion := 0.2
		if math.Abs(diversion-expectedDiversion) > 1e-9 {
			t.Errorf("diversion = %f, 期望 %f", diversion, expectedDiversion)
		}
	})

	t.Run("maxDamage为0时全部为0", func(t *testing.T) {
		sediment, flow, diversion := s.CalcSecondaryEffects(0)

		if sediment != 0 {
			t.Errorf("sediment = %f, 期望 0", sediment)
		}
		if flow != 0 {
			t.Errorf("flow = %f, 期望 0", flow)
		}
		if diversion != 0 {
			t.Errorf("diversion = %f, 期望 0", diversion)
		}
	})
}

func TestGenerateTimeSeries(t *testing.T) {
	s := NewEarthquakeSimulator(context.Background(), 4)

	t.Run("时间序列长度正确", func(t *testing.T) {
		duration := 5.0
		timeStep := 0.01
		expectedLength := int(duration / timeStep)

		series := s.GenerateTimeSeries(1.0, duration)

		if len(series) != expectedLength {
			t.Errorf("时间序列长度 = %d, 期望 %d", len(series), expectedLength)
		}
	})

	t.Run("时间序列值在合理范围", func(t *testing.T) {
		pga := 1.0
		duration := 2.0

		series := s.GenerateTimeSeries(pga, duration)

		for i, val := range series {
			t.Logf("series[%d] = %f", i, val)
			if i > 20 {
				break
			}
		}

		for i, val := range series {
			if math.IsNaN(val) {
				t.Errorf("时间序列第 %d 个值为 NaN", i)
			}
			if math.IsInf(val, 0) {
				t.Errorf("时间序列第 %d 个值为 Inf", i)
			}
		}
	})

	t.Run("持续时间为0返回空序列", func(t *testing.T) {
		series := s.GenerateTimeSeries(1.0, 0)
		if len(series) != 0 {
			t.Errorf("持续时间为0时序列长度 = %d, 期望 0", len(series))
		}
	})
}

func TestNewmarkBetaIntegration(t *testing.T) {
	s := NewEarthquakeSimulator(context.Background(), 4)

	t.Run("简单单自由度系统基本验证", func(t *testing.T) {
		mass := 1000.0
		stiffness := 100000.0
		dampingRatio := 0.05
		dt := 0.01

		wave := make([]float64, 100)
		for i := 0; i < 50; i++ {
			wave[i] = 1.0 * math.Sin(2*math.Pi*2*float64(i)*dt)
		}

		maxDisp, maxAcc := s.NewmarkBetaIntegration(mass, stiffness, dampingRatio, wave, dt)

		if math.IsNaN(maxDisp) {
			t.Error("maxDisp 为 NaN")
		}
		if math.IsNaN(maxAcc) {
			t.Error("maxAcc 为 NaN")
		}
		if math.IsInf(maxDisp, 0) {
			t.Error("maxDisp 为 Inf")
		}
		if math.IsInf(maxAcc, 0) {
			t.Error("maxAcc 为 Inf")
		}

		if maxDisp <= 0 {
			t.Errorf("maxDisp = %f, 应大于 0", maxDisp)
		}
		if maxAcc <= 0 {
			t.Errorf("maxAcc = %f, 应大于 0", maxAcc)
		}
	})

	t.Run("零输入波响应为零", func(t *testing.T) {
		mass := 1000.0
		stiffness := 100000.0
		dampingRatio := 0.05
		dt := 0.01

		wave := make([]float64, 100)

		maxDisp, maxAcc := s.NewmarkBetaIntegration(mass, stiffness, dampingRatio, wave, dt)

		if maxDisp != 0 {
			t.Errorf("零输入时 maxDisp = %f, 期望 0", maxDisp)
		}
		if maxAcc != 0 {
			t.Errorf("零输入时 maxAcc = %f, 期望 0", maxAcc)
		}
	})

	t.Run("单步脉冲响应", func(t *testing.T) {
		mass := 1.0
		stiffness := 1.0
		dampingRatio := 0.0
		dt := 0.01

		wave := []float64{1.0, 0, 0, 0, 0}

		maxDisp, maxAcc := s.NewmarkBetaIntegration(mass, stiffness, dampingRatio, wave, dt)

		if maxDisp <= 0 {
			t.Errorf("maxDisp = %f, 应大于 0", maxDisp)
		}
		if maxAcc <= 0 {
			t.Errorf("maxAcc = %f, 应大于 0", maxAcc)
		}
	})
}
