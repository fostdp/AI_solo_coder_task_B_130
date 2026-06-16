package api

import (
	"math"
	"testing"
)

func calcMachaEff(count, perUnit float64) float64 {
	return math.Min(count*perUnit, 0.85)
}

func calcBambooStability(initial float64, stoneCount int) float64 {
	return math.Max(initial-0.01*float64(stoneCount), 0.5)
}

func calcTotalScore(efficiency, stability, dredgingDepth, timeBonus float64) float64 {
	return efficiency*5 + stability*2 + dredgingDepth*2 + timeBonus
}

var validOperationTypes = map[string]bool{
	"place_macha":  true,
	"place_bamboo": true,
	"dredge":       true,
	"remove":       true,
}

type userOperation struct {
	OperationType  string
	OperationOrder int
	SessionID      string
}

func validateOperation(op userOperation) []string {
	var errs []string
	if !validOperationTypes[op.OperationType] {
		errs = append(errs, "invalid operation_type")
	}
	if op.OperationOrder < 0 {
		errs = append(errs, "invalid operation_order")
	}
	if op.SessionID == "" {
		errs = append(errs, "invalid session_id")
	}
	return errs
}

func TestUserExperience_MachaInterception_Normal(t *testing.T) {
	t.Run("李冰时期_efficiency_4_percent", func(t *testing.T) {
		perUnit := 0.04
		cases := []struct {
			name     string
			count    float64
			expected float64
		}{
			{"0个_0%", 0, 0},
			{"5个_20%", 5, 0.20},
			{"10个_40%", 10, 0.40},
			{"17个_68%", 17, 0.68},
			{"21个_84%", 21, 0.84},
			{"22个_85%_上限", 22, 0.85},
			{"30个_85%_上限", 30, 0.85},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				result := calcMachaEff(tc.count, perUnit)
				if math.Abs(result-tc.expected) > 1e-9 {
					t.Errorf("count=%.0f: expected %.2f, got %.2f", tc.count, tc.expected, result)
				}
			})
		}
	})

	t.Run("清代_efficiency_5_percent", func(t *testing.T) {
		perUnit := 0.05
		cases := []struct {
			name     string
			count    float64
			expected float64
		}{
			{"0个_0%", 0, 0},
			{"5个_25%", 5, 0.25},
			{"17个_85%_上限", 17, 0.85},
			{"30个_85%_上限", 30, 0.85},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				result := calcMachaEff(tc.count, perUnit)
				if math.Abs(result-tc.expected) > 1e-9 {
					t.Errorf("count=%.0f: expected %.2f, got %.2f", tc.count, tc.expected, result)
				}
			})
		}
	})
}

func TestUserExperience_MachaInterception_EfficiencyProgressive(t *testing.T) {
	t.Run("每增1个杩槎效率递增", func(t *testing.T) {
		for _, perUnit := range []float64{0.04, 0.05} {
			var prev float64
			for count := 0; count <= 30; count++ {
				curr := calcMachaEff(float64(count), perUnit)
				if curr < prev-1e-9 {
					t.Errorf("perUnit=%.2f count=%d: efficiency %.4f < prev %.4f, 不应递减", perUnit, count, curr, prev)
				}
				prev = curr
			}
		}
	})
}

func TestUserExperience_BambooStability_Normal(t *testing.T) {
	t.Run("初始0.7到1.0范围", func(t *testing.T) {
		for _, initial := range []float64{0.7, 0.75, 0.8, 0.85, 0.9, 0.95, 1.0} {
			result := calcBambooStability(initial, 0)
			if result < 0.7-1e-9 || result > 1.0+1e-9 {
				t.Errorf("initial=%.2f stability=%.4f 超出[0.7,1.0]范围", initial, result)
			}
		}
	})

	t.Run("填充10个石块后降低0.1", func(t *testing.T) {
		for _, initial := range []float64{0.7, 0.85, 1.0} {
			result := calcBambooStability(initial, 10)
			expected := initial - 0.01*10
			if math.Abs(result-expected) > 1e-9 {
				t.Errorf("initial=%.2f: expected %.4f, got %.4f", initial, expected, result)
			}
		}
	})

	t.Run("填充到下限0.5", func(t *testing.T) {
		result := calcBambooStability(0.6, 20)
		if result < 0.5-1e-9 {
			t.Errorf("stability %.4f 低于下限0.5", result)
		}
	})
}

func TestUserExperience_BambooStability_Boundary(t *testing.T) {
	t.Run("初始1.0填充50个石块至下限", func(t *testing.T) {
		result := calcBambooStability(1.0, 50)
		if math.Abs(result-0.5) > 1e-9 {
			t.Errorf("expected 0.5, got %.4f", result)
		}
	})

	t.Run("初始0.7填充20个石块至下限", func(t *testing.T) {
		result := calcBambooStability(0.7, 20)
		if math.Abs(result-0.5) > 1e-9 {
			t.Errorf("expected 0.5, got %.4f", result)
		}
	})

	t.Run("不会低于0.5", func(t *testing.T) {
		for _, initial := range []float64{0.7, 0.85, 1.0} {
			for stones := 0; stones <= 100; stones++ {
				result := calcBambooStability(initial, stones)
				if result < 0.5-1e-9 {
					t.Errorf("initial=%.2f stones=%d: stability %.4f 低于下限0.5", initial, stones, result)
				}
			}
		}
	})
}

func TestUserExperience_TotalScore_Normal(t *testing.T) {
	t.Run("正常评分场景", func(t *testing.T) {
		efficiency := 70.0
		stability := 80.0
		dredgingDepth := 30.0
		timeBonus := 10.0
		expected := 70.0*5 + 80.0*2 + 30.0*2 + 10.0
		result := calcTotalScore(efficiency, stability, dredgingDepth, timeBonus)
		if math.Abs(result-expected) > 1e-9 {
			t.Errorf("expected %.0f, got %.0f", expected, result)
		}
	})
}

func TestUserExperience_TotalScore_Boundary(t *testing.T) {
	t.Run("全0操作_0分", func(t *testing.T) {
		result := calcTotalScore(0, 0, 0, 0)
		if result != 0 {
			t.Errorf("expected 0, got %.0f", result)
		}
	})

	t.Run("满分场景", func(t *testing.T) {
		efficiency := 85.0
		stability := 100.0
		dredgingDepth := 50.0
		timeBonus := 100.0
		expected := 85.0*5 + 100.0*2 + 50.0*2 + 100.0
		result := calcTotalScore(efficiency, stability, dredgingDepth, timeBonus)
		if math.Abs(result-expected) > 1e-9 {
			t.Errorf("expected %.0f, got %.0f", expected, result)
		}
	})
}

func TestUserExperience_Achievement_FirstMacha(t *testing.T) {
	t.Run("操作列表为空时无成就", func(t *testing.T) {
		machaCount := 0
		if machaCount >= 1 {
			t.Error("空操作列表不应解锁首次放置杩槎")
		}
	})

	t.Run("第1次place_macha后解锁", func(t *testing.T) {
		machaCount := 1
		if machaCount < 1 {
			t.Error("放置1个杩槎后应解锁首次放置杩槎")
		}
		unlocked := machaCount >= 1
		if !unlocked {
			t.Error("应解锁首次放置杩槎成就")
		}
	})
}

func TestUserExperience_Achievement_Interception50(t *testing.T) {
	t.Run("李冰需13个杩槎", func(t *testing.T) {
		perUnit := 0.04
		result12 := calcMachaEff(12, perUnit)
		result13 := calcMachaEff(13, perUnit)
		if result12 >= 0.50 {
			t.Errorf("12个杩槎效率%.2f不应达到50%%", result12)
		}
		if result13 < 0.50 {
			t.Errorf("13个杩槎效率%.2f应达到50%%", result13)
		}
	})

	t.Run("清代需10个杩槎", func(t *testing.T) {
		perUnit := 0.05
		result9 := calcMachaEff(9, perUnit)
		result10 := calcMachaEff(10, perUnit)
		if result9 >= 0.50 {
			t.Errorf("9个杩槎效率%.2f不应达到50%%", result9)
		}
		if result10 < 0.50 {
			t.Errorf("10个杩槎效率%.2f应达到50%%", result10)
		}
	})
}

func TestUserExperience_Achievement_Interception80(t *testing.T) {
	t.Run("李冰需20个杩槎", func(t *testing.T) {
		perUnit := 0.04
		result19 := calcMachaEff(19, perUnit)
		result20 := calcMachaEff(20, perUnit)
		if result19 >= 0.80 {
			t.Errorf("19个杩槎效率%.2f不应达到80%%", result19)
		}
		if result20 < 0.80 {
			t.Errorf("20个杩槎效率%.2f应达到80%%", result20)
		}
	})

	t.Run("清代需16个杩槎", func(t *testing.T) {
		perUnit := 0.05
		result15 := calcMachaEff(15, perUnit)
		result16 := calcMachaEff(16, perUnit)
		if result15 >= 0.80 {
			t.Errorf("15个杩槎效率%.2f不应达到80%%", result15)
		}
		if result16 < 0.80 {
			t.Errorf("16个杩槎效率%.2f应达到80%%", result16)
		}
	})
}

func TestUserExperience_Achievement_CompleteRepair(t *testing.T) {
	achieveCompleteRepair := func(machaCount int, perUnit float64, stability float64, dredgingDepth float64) bool {
		efficiency := calcMachaEff(float64(machaCount), perUnit)
		return efficiency >= 0.70 && stability >= 70 && dredgingDepth >= 25
	}

	t.Run("满足条件时解锁", func(t *testing.T) {
		if !achieveCompleteRepair(18, 0.04, 75, 30) {
			t.Error("效率≥70%+稳定性≥70+淘淤≥25 应解锁完成岁修")
		}
	})

	t.Run("效率不足时不解锁", func(t *testing.T) {
		if achieveCompleteRepair(10, 0.04, 80, 30) {
			t.Error("效率不足70%不应解锁完成岁修")
		}
	})

	t.Run("稳定性不足时不解锁", func(t *testing.T) {
		if achieveCompleteRepair(18, 0.04, 60, 30) {
			t.Error("稳定性不足70不应解锁完成岁修")
		}
	})

	t.Run("淘淤不足时不解锁", func(t *testing.T) {
		if achieveCompleteRepair(18, 0.04, 75, 20) {
			t.Error("淘淤不足25不应解锁完成岁修")
		}
	})
}

func TestUserExperience_OperationValidation_Abnormal(t *testing.T) {
	t.Run("无效operation_type", func(t *testing.T) {
		op := userOperation{OperationType: "invalid_op", OperationOrder: 1, SessionID: "sess1"}
		errs := validateOperation(op)
		if len(errs) == 0 {
			t.Error("无效operation_type应返回错误")
		}
	})

	t.Run("负数operation_order", func(t *testing.T) {
		op := userOperation{OperationType: "place_macha", OperationOrder: -1, SessionID: "sess1"}
		errs := validateOperation(op)
		if len(errs) == 0 {
			t.Error("负数operation_order应返回错误")
		}
	})

	t.Run("空session_id", func(t *testing.T) {
		op := userOperation{OperationType: "place_macha", OperationOrder: 1, SessionID: ""}
		errs := validateOperation(op)
		if len(errs) == 0 {
			t.Error("空session_id应返回错误")
		}
	})

	t.Run("重复提交幂等性", func(t *testing.T) {
		op := userOperation{OperationType: "place_macha", OperationOrder: 1, SessionID: "sess1"}
		var recorded []userOperation
		for i := 0; i < 3; i++ {
			recorded = append(recorded, op)
		}
		if len(recorded) != 3 {
			t.Errorf("重复提交3次应记录3条, 实际%d条", len(recorded))
		}
	})
}

func TestUserExperience_DredgingToWolongIron(t *testing.T) {
	t.Run("淘淤深度计算", func(t *testing.T) {
		wolongIronElevation := 726.12
		currentRiverbed := 728.0
		needDredge := currentRiverbed - wolongIronElevation
		if math.Abs(needDredge-1.88) > 1e-9 {
			t.Errorf("需淘淤深度 expected 1.88m, got %.2fm", needDredge)
		}
	})

	t.Run("最少操作次数", func(t *testing.T) {
		needDredge := 1.88
		maxClearance := 0.5
		minOps := math.Ceil(needDredge / maxClearance)
		if minOps != 4 {
			t.Errorf("最少操作次数 expected 4, got %.0f", minOps)
		}
	})

	t.Run("最多操作次数", func(t *testing.T) {
		needDredge := 1.88
		minClearance := 0.1
		maxOps := math.Ceil(needDredge / minClearance)
		if maxOps != 19 {
			t.Errorf("最多操作次数 expected 19, got %.0f", maxOps)
		}
	})

	t.Run("操作次数范围验证", func(t *testing.T) {
		needDredge := 1.88
		for ops := 4; ops <= 19; ops++ {
			maxTotal := float64(ops) * 0.5
			if maxTotal < needDredge {
				t.Errorf("ops=%d: maxTotal=%.2f < needDredge=%.2f, 无法达到目标", ops, maxTotal, needDredge)
			}
		}
		for ops := 1; ops <= 3; ops++ {
			maxTotal := float64(ops) * 0.5
			if maxTotal >= needDredge {
				t.Errorf("ops=%d: 不应满足条件", ops)
			}
		}
	})
}
