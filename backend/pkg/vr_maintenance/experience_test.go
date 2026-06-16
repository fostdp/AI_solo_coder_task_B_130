package vr_maintenance

import (
	"context"
	"math"
	"testing"

	"dujiangyan-system/pkg/models"
)

func TestCalcMachaEfficiency(t *testing.T) {
	v := NewVRMaintenance(context.Background())

	t.Run("0个杩槎效率为0", func(t *testing.T) {
		result := v.CalcMachaEfficiency(0, MachaEffPerUnitQing)
		if result != 0 {
			t.Errorf("CalcMachaEfficiency(0, %f) = %f, 期望 0", MachaEffPerUnitQing, result)
		}
	})

	t.Run("13个清代杩槎效率", func(t *testing.T) {
		result := v.CalcMachaEfficiency(13, MachaEffPerUnitQing)
		expected := 13.0 * MachaEffPerUnitQing
		if result != expected {
			t.Errorf("CalcMachaEfficiency(13, %f) = %f, 期望 %f", MachaEffPerUnitQing, result, expected)
		}
	})

	t.Run("20个清代杩槎受上限约束0.85", func(t *testing.T) {
		result := v.CalcMachaEfficiency(20, MachaEffPerUnitQing)
		if result != MachaEffUpperLimit {
			t.Errorf("CalcMachaEfficiency(20, %f) = %f, 期望上限 %f", MachaEffPerUnitQing, result, MachaEffUpperLimit)
		}
	})

	t.Run("效率不超过上限", func(t *testing.T) {
		result := v.CalcMachaEfficiency(100, MachaEffPerUnitQing)
		if result > MachaEffUpperLimit {
			t.Errorf("效率 %f 超过上限 %f", result, MachaEffUpperLimit)
		}
	})
}

func TestCalcStability(t *testing.T) {
	v := NewVRMaintenance(context.Background())

	t.Run("0石时稳定性为初始值0.7", func(t *testing.T) {
		result := v.CalcStability(BambooBaseStability, 0)
		if result != BambooBaseStability {
			t.Errorf("CalcStability(%f, 0) = %f, 期望 %f", BambooBaseStability, result, BambooBaseStability)
		}
	})

	t.Run("填充10石稳定性计算", func(t *testing.T) {
		result := v.CalcStability(BambooBaseStability, 10)
		expected := BambooBaseStability + float64(10)*BambooStabPerStone - float64(10)/10.0*BambooStabDecayPer
		if result != expected {
			t.Errorf("CalcStability(%f, 10) = %f, 期望 %f", BambooBaseStability, result, expected)
		}
	})

	t.Run("稳定性不低于下限0.5", func(t *testing.T) {
		result := v.CalcStability(0.3, 100)
		if result < BambooStabLowerLimit {
			t.Errorf("稳定性 %f 低于下限 %f", result, BambooStabLowerLimit)
		}
	})

	t.Run("稳定性不超过上限1.0", func(t *testing.T) {
		result := v.CalcStability(0.95, 10)
		if result > 1.0 {
			t.Errorf("稳定性 %f 超过上限 1.0", result)
		}
	})
}

func TestCalcDredgeProgress(t *testing.T) {
	v := NewVRMaintenance(context.Background())

	t.Run("未淘淤时床高728，到卧铁距离1.88m", func(t *testing.T) {
		ops := []models.UserRepairOperation{}
		currentBed, toWolong := v.CalcDredgeProgress(ops)

		expectedBed := 728.0
		if currentBed != expectedBed {
			t.Errorf("currentBed = %f, 期望 %f", currentBed, expectedBed)
		}

		expectedToWolong := 728.0 - WolongElevation
		if math.Abs(toWolong-expectedToWolong) > 1e-9 {
			t.Errorf("toWolong = %f, 期望 %f", toWolong, expectedToWolong)
		}
	})

	t.Run("淘1.88m到达卧铁", func(t *testing.T) {
		ops := []models.UserRepairOperation{
			{OperationType: "dredge", DredgingVolume: 1.88},
		}
		currentBed, toWolong := v.CalcDredgeProgress(ops)

		expectedBed := 728.0 - 1.88
		if math.Abs(currentBed-expectedBed) > 1e-9 {
			t.Errorf("currentBed = %f, 期望 %f", currentBed, expectedBed)
		}

		if math.Abs(toWolong) > 1e-9 {
			t.Errorf("toWolong = %f, 期望接近 0", toWolong)
		}
	})

	t.Run("非淘淤操作不影响进度", func(t *testing.T) {
		ops := []models.UserRepairOperation{
			{OperationType: "place_macha", DredgingVolume: 0.5},
			{OperationType: "place_bamboo", DredgingVolume: 0.3},
		}
		currentBed, _ := v.CalcDredgeProgress(ops)

		if currentBed != 728.0 {
			t.Errorf("currentBed = %f, 期望 728.0", currentBed)
		}
	})

	t.Run("多次淘淤累加", func(t *testing.T) {
		ops := []models.UserRepairOperation{
			{OperationType: "dredge", DredgingVolume: 0.5},
			{OperationType: "dredge", DredgingVolume: 0.3},
			{OperationType: "dredge", DredgingVolume: 0.2},
		}
		currentBed, _ := v.CalcDredgeProgress(ops)

		expectedBed := 728.0 - 1.0
		if math.Abs(currentBed-expectedBed) > 1e-9 {
			t.Errorf("currentBed = %f, 期望 %f", currentBed, expectedBed)
		}
	})
}

func TestCalcTotalScore(t *testing.T) {
	v := NewVRMaintenance(context.Background())

	t.Run("正常情况约580分", func(t *testing.T) {
		score := v.CalcTotalScore(0.5, 0.65, 1.0, 600)
		expected := 0.5*ScoreWeightIntercept*100 + 0.65*ScoreWeightStability*100 + 1.0*ScoreWeightDredging*100

		if math.Abs(score-expected) > 1e-9 {
			t.Errorf("CalcTotalScore = %f, 期望 %f", score, expected)
		}
	})

	t.Run("零分情况", func(t *testing.T) {
		score := v.CalcTotalScore(0, 0, 0, 600)
		if score != 0 {
			t.Errorf("CalcTotalScore(0, 0, 0, 600) = %f, 期望 0", score)
		}
	})

	t.Run("时间越短分数越高", func(t *testing.T) {
		scoreFast := v.CalcTotalScore(0.5, 0.5, 0.5, 100)
		scoreSlow := v.CalcTotalScore(0.5, 0.5, 0.5, 500)

		if scoreFast <= scoreSlow {
			t.Errorf("短时间得分 %f 应大于长时间得分 %f", scoreFast, scoreSlow)
		}
	})

	t.Run("超过600秒无时间奖励", func(t *testing.T) {
		score1 := v.CalcTotalScore(0.5, 0.5, 0.5, 600)
		score2 := v.CalcTotalScore(0.5, 0.5, 0.5, 1000)

		if score1 != score2 {
			t.Errorf("超过600秒后得分应相同，但 %f != %f", score1, score2)
		}
	})
}

func TestCheckAchievements(t *testing.T) {
	v := NewVRMaintenance(context.Background())

	t.Run("首放杩槎成就", func(t *testing.T) {
		ops := []models.UserRepairOperation{
			{OperationType: "place_macha"},
		}
		achievements := v.CheckAchievements(ops, 0.3, 0.7, false)

		hasFirstMacha := false
		for _, a := range achievements {
			if a == "first_macha" {
				hasFirstMacha = true
				break
			}
		}
		if !hasFirstMacha {
			t.Error("应获得 first_macha 成就")
		}
	})

	t.Run("10个竹笼成就", func(t *testing.T) {
		ops := make([]models.UserRepairOperation, 10)
		for i := 0; i < 10; i++ {
			ops[i] = models.UserRepairOperation{OperationType: "place_bamboo"}
		}
		achievements := v.CheckAchievements(ops, 0.3, 0.7, false)

		hasTenBamboo := false
		for _, a := range achievements {
			if a == "ten_bamboo" {
				hasTenBamboo = true
				break
			}
		}
		if !hasTenBamboo {
			t.Error("应获得 ten_bamboo 成就")
		}
	})

	t.Run("50%截流效率成就", func(t *testing.T) {
		ops := []models.UserRepairOperation{}
		achievements := v.CheckAchievements(ops, 0.5, 0.7, false)

		hasFifty := false
		for _, a := range achievements {
			if a == "fifty_percent_intercept" {
				hasFifty = true
				break
			}
		}
		if !hasFifty {
			t.Error("应获得 fifty_percent_intercept 成就")
		}
	})

	t.Run("80%截流效率成就", func(t *testing.T) {
		ops := []models.UserRepairOperation{}
		achievements := v.CheckAchievements(ops, 0.8, 0.7, false)

		hasEighty := false
		for _, a := range achievements {
			if a == "eighty_percent_intercept" {
				hasEighty = true
				break
			}
		}
		if !hasEighty {
			t.Error("应获得 eighty_percent_intercept 成就")
		}
	})

	t.Run("淘至卧铁成就", func(t *testing.T) {
		ops := []models.UserRepairOperation{}
		achievements := v.CheckAchievements(ops, 0.3, 0.7, true)

		hasDredgeToWolong := false
		for _, a := range achievements {
			if a == "dredge_to_wolong" {
				hasDredgeToWolong = true
				break
			}
		}
		if !hasDredgeToWolong {
			t.Error("应获得 dredge_to_wolong 成就")
		}
	})

	t.Run("完成全部修复成就", func(t *testing.T) {
		ops := make([]models.UserRepairOperation, 11)
		ops[0] = models.UserRepairOperation{OperationType: "place_macha"}
		for i := 1; i < 11; i++ {
			ops[i] = models.UserRepairOperation{OperationType: "place_bamboo"}
		}
		achievements := v.CheckAchievements(ops, 0.85, 0.8, true)

		hasComplete := false
		for _, a := range achievements {
			if a == "complete_repair" {
				hasComplete = true
				break
			}
		}
		if !hasComplete {
			t.Error("应获得 complete_repair 成就")
		}
	})

	t.Run("空操作无成就", func(t *testing.T) {
		ops := []models.UserRepairOperation{}
		achievements := v.CheckAchievements(ops, 0.3, 0.7, false)

		if len(achievements) != 0 {
			t.Errorf("应无成就，但获得 %d 个", len(achievements))
		}
	})
}

func TestValidateOperation(t *testing.T) {
	v := NewVRMaintenance(context.Background())

	t.Run("正常操作验证通过", func(t *testing.T) {
		result := v.ValidateOperation("place_macha", 1, "session_001", 500, 500)
		if !result {
			t.Error("正常操作应返回 true")
		}
	})

	t.Run("无效操作类型返回false", func(t *testing.T) {
		result := v.ValidateOperation("invalid_op", 1, "session_001", 500, 500)
		if result {
			t.Error("无效操作类型应返回 false")
		}
	})

	t.Run("负顺序返回false", func(t *testing.T) {
		result := v.ValidateOperation("place_macha", -1, "session_001", 500, 500)
		if result {
			t.Error("负顺序应返回 false")
		}
	})

	t.Run("空sessionID返回false", func(t *testing.T) {
		result := v.ValidateOperation("place_macha", 1, "", 500, 500)
		if result {
			t.Error("空sessionID应返回 false")
		}
	})

	t.Run("posX越界返回false", func(t *testing.T) {
		result := v.ValidateOperation("place_macha", 1, "session_001", -1, 500)
		if result {
			t.Error("posX为负应返回 false")
		}

		result = v.ValidateOperation("place_macha", 1, "session_001", 1001, 500)
		if result {
			t.Error("posX超过1000应返回 false")
		}
	})

	t.Run("posY越界返回false", func(t *testing.T) {
		result := v.ValidateOperation("place_macha", 1, "session_001", 500, -1)
		if result {
			t.Error("posY为负应返回 false")
		}

		result = v.ValidateOperation("place_macha", 1, "session_001", 500, 1001)
		if result {
			t.Error("posY超过1000应返回 false")
		}
	})

	t.Run("边界值验证通过", func(t *testing.T) {
		result := v.ValidateOperation("place_bamboo", 0, "session_001", 0, 0)
		if !result {
			t.Error("边界值(0,0,0)应返回 true")
		}

		result = v.ValidateOperation("dredge", 100, "session_001", 1000, 1000)
		if !result {
			t.Error("边界值(1000,1000)应返回 true")
		}
	})

	t.Run("remove_macha操作有效", func(t *testing.T) {
		result := v.ValidateOperation("remove_macha", 2, "session_001", 300, 400)
		if !result {
			t.Error("remove_macha 应返回 true")
		}
	})

	t.Run("remove_bamboo操作有效", func(t *testing.T) {
		result := v.ValidateOperation("remove_bamboo", 3, "session_001", 300, 400)
		if !result {
			t.Error("remove_bamboo 应返回 true")
		}
	})
}
