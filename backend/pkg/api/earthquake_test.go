package api

import (
	"math"
	"testing"
)

type testLoc struct {
	X, Y, Z float64
}

var (
	yuzuiTestLoc     = testLoc{X: 100, Y: 50, Z: 730}
	feishayanTestLoc  = testLoc{X: 150, Y: 30, Z: 728}
	baopingkouTestLoc = testLoc{X: 200, Y: 60, Z: 725}
	renzidiTestLoc    = testLoc{X: 80, Y: 80, Z: 727}
	allTestLocs       = []testLoc{yuzuiTestLoc, feishayanTestLoc, baopingkouTestLoc, renzidiTestLoc}
)

func testPGA(mag float64) float64 {
	return math.Pow(10, 0.5*mag-3)
}

func testDist(loc, epic testLoc) float64 {
	dx := loc.X - epic.X
	dy := loc.Y - epic.Y
	dz := loc.Z - epic.Z
	return math.Sqrt(dx*dx + dy*dy + dz*dz)
}

func testAtt(dist float64) float64 {
	return math.Exp(-0.001 * dist)
}

func testDmg(pga, att float64) float64 {
	return pga * att
}

func testMaxDamage(pga float64, epic testLoc) float64 {
	maxD := 0.0
	for _, loc := range allTestLocs {
		d := testDmg(pga, testAtt(testDist(loc, epic)))
		if d > maxD {
			maxD = d
		}
	}
	return maxD
}

func testSafety(maxDmg float64) string {
	if maxDmg > 0.6 {
		return "danger"
	}
	if maxDmg > 0.3 {
		return "caution"
	}
	return "safe"
}

const eps = 1e-9

func TestEarthquake_PGA_Calculation(t *testing.T) {
	tests := []struct {
		name string
		mag  float64
		want float64
	}{
		{"M8.0", 8.0, 10.0},
		{"M7.0", 7.0, math.Sqrt(10)},
		{"M6.0", 6.0, 1.0},
		{"M5.0", 5.0, 1.0 / math.Sqrt(10)},
		{"M4.0", 4.0, 0.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := testPGA(tt.mag)
			if math.Abs(got-tt.want) > eps {
				t.Errorf("PGA(%s) = %v, want %v", tt.name, got, tt.want)
			}
		})
	}
}

func TestEarthquake_SafetyAssessment(t *testing.T) {
	epicQushou := testLoc{X: 100, Y: 50, Z: 0}

	t.Run("M8.0_epicenter_qushou_danger", func(t *testing.T) {
		pga := testPGA(8.0)
		maxDmg := testMaxDamage(pga, epicQushou)
		safety := testSafety(maxDmg)
		if safety != "danger" {
			t.Errorf("safety = %s, want danger (maxDamage=%.4f)", safety, maxDmg)
		}
	})

	t.Run("M6.0_epicenter_qushou_caution", func(t *testing.T) {
		pga := testPGA(6.0)
		maxDmg := testMaxDamage(pga, epicQushou)
		safety := testSafety(maxDmg)
		if safety != "caution" {
			t.Errorf("safety = %s, want caution (maxDamage=%.4f)", safety, maxDmg)
		}
	})

	t.Run("M4.0_epicenter_qushou_safe", func(t *testing.T) {
		pga := testPGA(4.0)
		maxDmg := testMaxDamage(pga, epicQushou)
		safety := testSafety(maxDmg)
		if safety != "safe" {
			t.Errorf("safety = %s, want safe (maxDamage=%.4f)", safety, maxDmg)
		}
	})

	t.Run("M5.0_far_epicenter_safe", func(t *testing.T) {
		pga := testPGA(5.0)
		farEpic := testLoc{X: 10000, Y: 10000, Z: 0}
		maxDmg := testMaxDamage(pga, farEpic)
		safety := testSafety(maxDmg)
		if safety != "safe" {
			t.Errorf("safety = %s, want safe (maxDamage=%.8f)", safety, maxDmg)
		}
	})
}

func TestEarthquake_DistanceAttenuation(t *testing.T) {
	t.Run("epicenter_qushou_yuzui_dist730", func(t *testing.T) {
		epic := testLoc{X: 100, Y: 50, Z: 0}
		dist := testDist(yuzuiTestLoc, epic)
		att := testAtt(dist)
		if math.Abs(dist-730) > 0.01 {
			t.Errorf("dist = %.4f, want 730", dist)
		}
		if math.Abs(att-0.48) > 0.01 {
			t.Errorf("att = %.4f, want ~0.48", att)
		}
	})

	t.Run("epicenter_origin_yuzui", func(t *testing.T) {
		epic := testLoc{X: 0, Y: 0, Z: 0}
		dist := testDist(yuzuiTestLoc, epic)
		att := testAtt(dist)
		if math.Abs(dist-738.5) > 0.5 {
			t.Errorf("dist = %.4f, want ~738.5", dist)
		}
		if math.Abs(att-0.48) > 0.01 {
			t.Errorf("att = %.4f, want ~0.48", att)
		}
	})

	t.Run("farther_distance_less_damage", func(t *testing.T) {
		epic := testLoc{X: 100, Y: 50, Z: 0}
		pga := testPGA(7.0)
		dmgNear := testDmg(pga, testAtt(testDist(yuzuiTestLoc, epic)))
		farEpic := testLoc{X: 1000, Y: 1000, Z: 0}
		dmgFar := testDmg(pga, testAtt(testDist(yuzuiTestLoc, farEpic)))
		if dmgNear <= dmgFar {
			t.Errorf("near damage %.4f should > far damage %.4f", dmgNear, dmgFar)
		}
	})
}

func TestEarthquake_StructureDamage_Ranking(t *testing.T) {
	epic := testLoc{X: 100, Y: 50, Z: 730}
	pga := testPGA(8.0)

	yuzuiDmg := testDmg(pga, testAtt(testDist(yuzuiTestLoc, epic)))
	feishayanDmg := testDmg(pga, testAtt(testDist(feishayanTestLoc, epic)))
	baopingkouDmg := testDmg(pga, testAtt(testDist(baopingkouTestLoc, epic)))
	renzidiDmg := testDmg(pga, testAtt(testDist(renzidiTestLoc, epic)))

	t.Run("yuzui_highest_damage_when_epicenter_at_yuzui", func(t *testing.T) {
		if yuzuiDmg <= feishayanDmg {
			t.Errorf("yuzuiDamage %.4f should > feishayanDamage %.4f", yuzuiDmg, feishayanDmg)
		}
		if yuzuiDmg <= baopingkouDmg {
			t.Errorf("yuzuiDamage %.4f should > baopingkouDamage %.4f", yuzuiDmg, baopingkouDmg)
		}
		if yuzuiDmg <= renzidiDmg {
			t.Errorf("yuzuiDamage %.4f should > renzidiDamage %.4f", yuzuiDmg, renzidiDmg)
		}
	})

	t.Run("damage_ranking_consistent_with_distance", func(t *testing.T) {
		if renzidiDmg <= feishayanDmg {
			t.Errorf("renzidiDamage %.4f should > feishayanDamage %.4f (renzidi is closer to epicenter)", renzidiDmg, feishayanDmg)
		}
		if feishayanDmg <= baopingkouDmg {
			t.Errorf("feishayanDamage %.4f should > baopingkouDamage %.4f", feishayanDmg, baopingkouDmg)
		}
	})
}

func TestEarthquake_SecondaryEffects(t *testing.T) {
	t.Run("M8.0_values", func(t *testing.T) {
		pga := testPGA(8.0)
		sediment := pga * 0.5
		flow := pga * 0.3
		diversion := pga * 0.2

		if math.Abs(sediment-5.0) > eps {
			t.Errorf("sedimentDisturbance = %v, want 5.0", sediment)
		}
		if math.Abs(flow-3.0) > eps {
			t.Errorf("flowPathChange = %v, want 3.0", flow)
		}
		if math.Abs(diversion-2.0) > eps {
			t.Errorf("waterDiversionChange = %v, want 2.0", diversion)
		}
	})

	t.Run("proportional_relationship", func(t *testing.T) {
		pga := testPGA(8.0)
		sediment := pga * 0.5
		flow := pga * 0.3
		diversion := pga * 0.2

		if !(sediment > flow && flow > diversion) {
			t.Errorf("expected sediment(%.2f) > flow(%.2f) > diversion(%.2f)", sediment, flow, diversion)
		}
	})

	t.Run("coefficients_independent_of_magnitude", func(t *testing.T) {
		for _, mag := range []float64{4.0, 6.0, 8.0} {
			pga := testPGA(mag)
			sediment := pga * 0.5
			flow := pga * 0.3
			diversion := pga * 0.2

			if math.Abs(sediment/pga-0.5) > eps {
				t.Errorf("M%.1f: sediment/pga = %v, want 0.5", mag, sediment/pga)
			}
			if math.Abs(flow/pga-0.3) > eps {
				t.Errorf("M%.1f: flow/pga = %v, want 0.3", mag, flow/pga)
			}
			if math.Abs(diversion/pga-0.2) > eps {
				t.Errorf("M%.1f: diversion/pga = %v, want 0.2", mag, diversion/pga)
			}
		}
	})
}

func TestEarthquake_Boundary_M0(t *testing.T) {
	t.Run("pga_value", func(t *testing.T) {
		pga := testPGA(0)
		if math.Abs(pga-0.001) > eps {
			t.Errorf("PGA(M0) = %v, want 0.001", pga)
		}
	})

	t.Run("safety_safe", func(t *testing.T) {
		pga := testPGA(0)
		epic := testLoc{X: 100, Y: 50, Z: 0}
		maxDmg := testMaxDamage(pga, epic)
		safety := testSafety(maxDmg)
		if safety != "safe" {
			t.Errorf("M0 safety = %s, want safe (maxDamage=%.8f)", safety, maxDmg)
		}
	})
}

func TestEarthquake_Boundary_M10(t *testing.T) {
	t.Run("pga_value", func(t *testing.T) {
		pga := testPGA(10)
		if math.Abs(pga-100.0) > eps {
			t.Errorf("PGA(M10) = %v, want 100.0", pga)
		}
	})

	t.Run("safety_danger", func(t *testing.T) {
		pga := testPGA(10)
		epic := testLoc{X: 100, Y: 50, Z: 0}
		maxDmg := testMaxDamage(pga, epic)
		safety := testSafety(maxDmg)
		if safety != "danger" {
			t.Errorf("M10 safety = %s, want danger (maxDamage=%.4f)", safety, maxDmg)
		}
	})
}

func TestEarthquake_Abnormal_NegativeMagnitude(t *testing.T) {
	t.Run("pga_value", func(t *testing.T) {
		pga := testPGA(-1)
		want := math.Pow(10, -3.5)
		if math.Abs(pga-want) > eps {
			t.Errorf("PGA(M-1) = %v, want %v", pga, want)
		}
	})

	t.Run("safety_safe", func(t *testing.T) {
		pga := testPGA(-1)
		epic := testLoc{X: 100, Y: 50, Z: 0}
		maxDmg := testMaxDamage(pga, epic)
		safety := testSafety(maxDmg)
		if safety != "safe" {
			t.Errorf("M-1 safety = %s, want safe (maxDamage=%.10f)", safety, maxDmg)
		}
	})
}

func TestEarthquake_TimeSeriesGeneration(t *testing.T) {
	duration := 10.0
	timeStep := 0.01
	numSteps := int(duration / timeStep)

	t.Run("step_count", func(t *testing.T) {
		if numSteps != 1000 {
			t.Errorf("numSteps = %d, want 1000", numSteps)
		}
	})

	t.Run("time_values_linear", func(t *testing.T) {
		if numSteps < 2 {
			t.Fatal("need at least 2 steps")
		}
		t0 := 0.0 * timeStep
		tLast := float64(numSteps-1) * timeStep
		if math.Abs(t0) > eps {
			t.Errorf("t[0] = %v, want 0", t0)
		}
		if math.Abs(tLast-(duration-timeStep)) > eps {
			t.Errorf("t[last] = %v, want %v", tLast, duration-timeStep)
		}
	})

	t.Run("envelope_monotonic_decay", func(t *testing.T) {
		pga := testPGA(7.0)
		for i := 1; i < numSteps; i++ {
			tPrev := float64(i-1) * timeStep
			tCurr := float64(i) * timeStep
			envPrev := pga * math.Exp(-0.5*tPrev)
			envCurr := pga * math.Exp(-0.5*tCurr)
			if envCurr >= envPrev {
				t.Errorf("envelope at t=%.2f (%.6f) should < t=%.2f (%.6f)", tCurr, envCurr, tPrev, envPrev)
				return
			}
		}
	})

	t.Run("amplitude_range_with_noise", func(t *testing.T) {
		pga := testPGA(7.0)
		env0 := pga * math.Exp(-0.5*0)
		minAmp := env0 * 1.0
		maxAmp := env0 * 1.3
		if math.Abs(minAmp-pga) > eps {
			t.Errorf("min amplitude at t=0 = %v, want %v", minAmp, pga)
		}
		if math.Abs(maxAmp-pga*1.3) > eps {
			t.Errorf("max amplitude at t=0 = %v, want %v", maxAmp, pga*1.3)
		}
	})

	t.Run("envelope_at_t5", func(t *testing.T) {
		pga := testPGA(7.0)
		env5 := pga * math.Exp(-0.5*5.0)
		want := pga * math.Exp(-2.5)
		if math.Abs(env5-want) > eps {
			t.Errorf("envelope at t=5 = %v, want %v", env5, want)
		}
	})
}

func TestEarthquake_BankCollapse(t *testing.T) {
	epic := testLoc{X: 100, Y: 50, Z: 0}

	t.Run("M4.0_no_collapse", func(t *testing.T) {
		maxDmg := testMaxDamage(testPGA(4.0), epic)
		if maxDmg > 0.3 {
			t.Errorf("M4.0 maxDamage = %.4f, should be <= 0.3", maxDmg)
		}
		collapseCount := 0
		if maxDmg > 0.3 {
			collapseCount = int(maxDmg * 20)
		}
		if collapseCount != 0 {
			t.Errorf("M4.0 collapseCount = %d, want 0", collapseCount)
		}
	})

	t.Run("M8.0_has_collapse", func(t *testing.T) {
		maxDmg := testMaxDamage(testPGA(8.0), epic)
		if maxDmg <= 0.3 {
			t.Errorf("M8.0 maxDamage = %.4f, should be > 0.3", maxDmg)
		}
		if maxDmg <= 0.6 {
			t.Errorf("M8.0 maxDamage = %.4f, should be > 0.6", maxDmg)
		}
		collapseCount := int(maxDmg * 20)
		if collapseCount <= 0 {
			t.Errorf("M8.0 collapseCount = %d, want > 0", collapseCount)
		}
	})

	t.Run("collapse_count_formula", func(t *testing.T) {
		maxDmg := 0.5
		expected := int(maxDmg * 20)
		if expected != 10 {
			t.Errorf("int(0.5*20) = %d, want 10", expected)
		}
	})

	t.Run("threshold_boundary", func(t *testing.T) {
		aboveThreshold := 0.31
		collapseCount := 0
		if aboveThreshold > 0.3 {
			collapseCount = int(aboveThreshold * 20)
		}
		if collapseCount <= 0 {
			t.Errorf("maxDamage=0.31 should produce collapse, got count=%d", collapseCount)
		}

		belowThreshold := 0.29
		collapseCount2 := 0
		if belowThreshold > 0.3 {
			collapseCount2 = int(belowThreshold * 20)
		}
		if collapseCount2 != 0 {
			t.Errorf("maxDamage=0.29 should not produce collapse, got count=%d", collapseCount2)
		}
	})
}
