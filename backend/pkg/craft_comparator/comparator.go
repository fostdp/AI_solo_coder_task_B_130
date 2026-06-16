package craft_comparator

import (
	"context"
	"encoding/json"
	"math"

	"dujiangyan-system/pkg/models"
)

type CraftComparator struct {
	ctx context.Context
}

func NewCraftComparator(ctx context.Context) *CraftComparator {
	return &CraftComparator{ctx: ctx}
}

func (c *CraftComparator) GetTechniques(dynasty, category string) ([]models.DynastyRepairTechnique, error) {
	return models.GetDynastyTechniques(c.ctx, dynasty, category)
}

func (c *CraftComparator) CalcMachaInterception(count int, efficiencyPerUnit float64) float64 {
	result := float64(count) * efficiencyPerUnit
	return math.Min(result, 0.85)
}

func (c *CraftComparator) CalcBambooStability(count int, baseStability float64) float64 {
	result := baseStability + float64(count)*0.03
	if result > 1.0 {
		return 1.0
	}
	return result
}

func (c *CraftComparator) CompareEfficiency(qingCount, libingCount int, qingEff, libingEff float64) map[string]float64 {
	qingEfficiency := float64(qingCount) * qingEff
	libingEfficiency := float64(libingCount) * libingEff

	improvementRatio := 0.0
	if libingEfficiency > 0 {
		improvementRatio = qingEfficiency / libingEfficiency
	}

	return map[string]float64{
		"qing_efficiency":   qingEfficiency,
		"libing_efficiency": libingEfficiency,
		"improvement_ratio": improvementRatio,
	}
}

func (c *CraftComparator) ValidateArchaeologicalParams(paramsJSON string) (hasSource bool, hasUncertainty bool, hasMethod bool) {
	var params map[string]interface{}
	if err := json.Unmarshal([]byte(paramsJSON), &params); err != nil {
		return false, false, false
	}

	_, hasSource = params["source_archaeology"]
	_, hasUncertainty = params["uncertainty_range"]
	_, hasMethod = params["experimental_method"]
	return
}
