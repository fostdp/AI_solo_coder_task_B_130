package era_comparator

import (
	"context"

	"dujiangyan-system/pkg/models"
)

type EraComparator struct {
	ctx context.Context
}

func NewEraComparator(ctx context.Context) *EraComparator {
	return &EraComparator{ctx: ctx}
}

func (c *EraComparator) GetComparisons(category string) ([]models.ModernRepairComparison, error) {
	return models.GetModernComparisons(c.ctx, category)
}

func (c *EraComparator) ReductionRate(ancient, modern float64) float64 {
	if ancient == 0 {
		return 0
	}
	return (ancient - modern) / ancient * 100
}

func (c *EraComparator) CostEfficiencyRatio(cost, efficiency float64) float64 {
	if cost == 0 {
		return 0
	}
	return efficiency / cost
}

func (c *EraComparator) OverallComparison(comparisons []models.ModernRepairComparison) map[string]interface{} {
	n := len(comparisons)
	if n == 0 {
		return map[string]interface{}{
			"avg_ancient_cost":         0.0,
			"avg_modern_cost":          0.0,
			"avg_cost_reduction":       0.0,
			"avg_efficiency_improvement": 0.0,
		}
	}

	var totalAncientCost, totalModernCost float64
	var totalCostReduction, totalEfficiencyImprovement float64

	for _, comp := range comparisons {
		totalAncientCost += comp.AncientCost
		totalModernCost += comp.ModernCost
		totalCostReduction += c.ReductionRate(comp.AncientCost, comp.ModernCost)

		if comp.AncientEfficiency > 0 {
			totalEfficiencyImprovement += (comp.ModernEfficiency - comp.AncientEfficiency) / comp.AncientEfficiency * 100
		}
	}

	return map[string]interface{}{
		"avg_ancient_cost":           totalAncientCost / float64(n),
		"avg_modern_cost":            totalModernCost / float64(n),
		"avg_cost_reduction":         totalCostReduction / float64(n),
		"avg_efficiency_improvement": totalEfficiencyImprovement / float64(n),
	}
}

func (c *EraComparator) ValidateStandardRef(comparisons []models.ModernRepairComparison) (withUnit int, withStandard int, withCode int) {
	for _, comp := range comparisons {
		if comp.Unit != "" {
			withUnit++
		}
		if comp.StandardReference != "" {
			withStandard++
		}
		if comp.StandardCode != "" {
			withCode++
		}
	}
	return
}
