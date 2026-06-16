package vr_maintenance

import (
	"context"
	"math"

	"dujiangyan-system/pkg/models"
)

const (
	WolongElevation       = 726.12
	MachaEffPerUnitQing   = 0.05
	MachaEffPerUnitLiBing = 0.04
	MachaEffUpperLimit    = 0.85
	BambooBaseStability   = 0.7
	BambooStabPerStone    = 0.01
	BambooStabDecayPer    = 0.1
	BambooStabLowerLimit  = 0.5
	ScoreWeightIntercept  = 5.0
	ScoreWeightStability  = 2.0
	ScoreWeightDredging   = 2.0
)

type OperationRequest struct {
	SessionID              string
	UserNickname           string
	OperationType          string
	ObjectType             string
	PositionX              float64
	PositionY              float64
	PositionZ              float64
	RotationAngle          float64
	ObjectParams           string
	OperationOrder         int
	InterceptionEfficiency float64
	StabilityScore         float64
	DredgingVolume         float64
}

type SessionResult struct {
	SessionID        string
	TotalScore       float64
	CompletionStatus string
	Achievement      string
	DurationSeconds  int
}

type VRMaintenance struct {
	ctx context.Context
}

func NewVRMaintenance(ctx context.Context) *VRMaintenance {
	return &VRMaintenance{ctx: ctx}
}

func (v *VRMaintenance) CalcMachaEfficiency(count int, perUnit float64) float64 {
	return math.Min(float64(count)*perUnit, MachaEffUpperLimit)
}

func (v *VRMaintenance) CalcStability(baseStability float64, filledStones int) float64 {
	stability := baseStability + float64(filledStones)*BambooStabPerStone - float64(filledStones)/10.0*BambooStabDecayPer
	return math.Max(BambooStabLowerLimit, math.Min(1.0, stability))
}

func (v *VRMaintenance) CalcDredgeProgress(ops []models.UserRepairOperation) (currentBed float64, toWolong float64) {
	var cumulativeDredge float64
	for _, op := range ops {
		if op.OperationType == "dredge" {
			cumulativeDredge += op.DredgingVolume
		}
	}
	currentBed = 728.0 - cumulativeDredge
	toWolong = currentBed - WolongElevation
	return
}

func (v *VRMaintenance) CalcTotalScore(interceptEff, stability, dredgeDepth float64, durationSec int) float64 {
	interceptPart := interceptEff * ScoreWeightIntercept * 100
	stabPart := stability * ScoreWeightStability * 100
	dredgePart := dredgeDepth * ScoreWeightDredging * 100
	timeBonus := math.Max(0, (600-float64(durationSec))/600) * 50
	return interceptPart + stabPart + dredgePart + timeBonus
}

func (v *VRMaintenance) CheckAchievements(ops []models.UserRepairOperation, interceptEff, stability float64, dredgeToWolong bool) []string {
	var machaCount, bambooCount int
	for _, op := range ops {
		switch op.OperationType {
		case "place_macha":
			machaCount++
		case "place_bamboo":
			bambooCount++
		}
	}

	var achievements []string

	hasFirstMacha := machaCount >= 1
	hasTenBamboo := bambooCount >= 10
	hasFiftyIntercept := interceptEff >= 0.5
	hasEightyIntercept := interceptEff >= 0.8

	if hasFirstMacha {
		achievements = append(achievements, "first_macha")
	}
	if hasTenBamboo {
		achievements = append(achievements, "ten_bamboo")
	}
	if hasFiftyIntercept {
		achievements = append(achievements, "fifty_percent_intercept")
	}
	if hasEightyIntercept {
		achievements = append(achievements, "eighty_percent_intercept")
	}
	if dredgeToWolong {
		achievements = append(achievements, "dredge_to_wolong")
	}

	if hasFirstMacha && hasTenBamboo && hasEightyIntercept && dredgeToWolong {
		achievements = append(achievements, "complete_repair")
	}

	return achievements
}

func (v *VRMaintenance) ValidateOperation(opType string, order int, sessionID string, posX, posY float64) bool {
	validTypes := map[string]bool{
		"place_macha":  true,
		"place_bamboo": true,
		"dredge":       true,
		"remove_macha": true,
		"remove_bamboo": true,
	}
	if !validTypes[opType] {
		return false
	}
	if order < 0 {
		return false
	}
	if sessionID == "" {
		return false
	}
	if posX < 0 || posX > 1000 {
		return false
	}
	if posY < 0 || posY > 1000 {
		return false
	}
	return true
}

func (v *VRMaintenance) RecordOperation(req OperationRequest) (*models.UserRepairOperation, error) {
	op := &models.UserRepairOperation{
		SessionID:              req.SessionID,
		UserNickname:           req.UserNickname,
		OperationType:          req.OperationType,
		ObjectType:             req.ObjectType,
		PositionX:              req.PositionX,
		PositionY:              req.PositionY,
		PositionZ:              req.PositionZ,
		RotationAngle:          req.RotationAngle,
		ObjectParams:           req.ObjectParams,
		OperationOrder:         req.OperationOrder,
		InterceptionEfficiency: req.InterceptionEfficiency,
		StabilityScore:         req.StabilityScore,
		DredgingVolume:         req.DredgingVolume,
		CompletionStatus:       "in_progress",
	}

	id, err := models.InsertUserOperation(v.ctx, op)
	if err != nil {
		return nil, err
	}
	op.ID = id
	return op, nil
}

func (v *VRMaintenance) GetOperations(sessionID string) ([]models.UserRepairOperation, error) {
	return models.GetUserOperations(v.ctx, sessionID)
}

func (v *VRMaintenance) FinishSession(result SessionResult) error {
	return models.UpdateUserSessionScore(v.ctx, result.SessionID, result.TotalScore, result.CompletionStatus, result.Achievement, result.DurationSeconds)
}

func (v *VRMaintenance) GetRanking(limit int) ([]models.UserRepairOperation, error) {
	return models.GetUserScoreRanking(v.ctx, limit)
}
