package earthquake_simulator

import (
	"context"
	"encoding/json"
	"math"
	"math/rand"
	"sync"
	"time"

	"dujiangyan-system/pkg/models"
)

type StructureLocation struct {
	X float64
	Y float64
	Z float64
}

type SimulationRequest struct {
	SimulationName  string
	Magnitude       float64
	EpicenterX      float64
	EpicenterY      float64
	EpicenterZ      float64
	FocalMechanism  string
	DurationSeconds float64
	CreatedBy       string
}

type SimulationResult struct {
	Simulation             *models.EarthquakeSimulation
	DynamicAnalysisEnabled bool
	ParallelWorkers        int
}

type EarthquakeSimulator struct {
	ctx             context.Context
	parallelWorkers int
}

func NewEarthquakeSimulator(ctx context.Context, workers int) *EarthquakeSimulator {
	if workers <= 0 {
		workers = 4
	}
	return &EarthquakeSimulator{ctx: ctx, parallelWorkers: workers}
}

func (s *EarthquakeSimulator) CalcPGA(magnitude float64) float64 {
	return math.Pow(10, 0.5*magnitude-3)
}

func (s *EarthquakeSimulator) CalcDamage(pga float64, loc StructureLocation, epicX, epicY float64) float64 {
	distance := math.Sqrt((loc.X-epicX)*(loc.X-epicX) + (loc.Y-epicY)*(loc.Y-epicY))
	attenuation := math.Exp(-0.001 * distance)
	damage := pga * attenuation * 0.1
	return math.Min(damage, 1.0)
}

func (s *EarthquakeSimulator) AssessSafety(yuzui, feishayan, baopingkou, renzidi float64) string {
	maxDamage := math.Max(math.Max(yuzui, feishayan), math.Max(baopingkou, renzidi))
	if maxDamage > 0.6 {
		return "danger"
	}
	if maxDamage > 0.3 {
		return "caution"
	}
	return "safe"
}

func (s *EarthquakeSimulator) CalcSecondaryEffects(maxDamage float64) (sediment, flow, diversion float64) {
	sediment = maxDamage * 0.5
	flow = maxDamage * 0.3
	diversion = maxDamage * 0.2
	return
}

func (s *EarthquakeSimulator) GenerateTimeSeries(pga, duration float64) []float64 {
	timeStep := 0.01
	numSteps := int(duration / timeStep)
	wave := make([]float64, numSteps)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numSteps; i++ {
		t := float64(i) * timeStep
		amplitude := pga * math.Exp(-0.5*t) * (1 + 0.3*r.Float64())
		acceleration := amplitude * math.Sin(2*math.Pi*2*t) * (1 + 0.2*r.Float64())
		wave[i] = acceleration
	}
	return wave
}

func (s *EarthquakeSimulator) NewmarkBetaIntegration(mass, stiffness, dampingRatio float64, wave []float64, dt float64) (maxDisp, maxAcc float64) {
	beta := 0.25
	gamma := 0.5
	nSteps := len(wave)

	disp := make([]float64, nSteps)
	vel := make([]float64, nSteps)
	acc := make([]float64, nSteps)

	omega := math.Sqrt(stiffness / mass)
	c := 2 * dampingRatio * mass * omega
	kPrime := stiffness + gamma/(beta*dt)*c + mass/(beta*dt*dt)

	disp[0] = 0
	vel[0] = 0
	acc[0] = wave[0]

	maxDisp = 0
	maxAcc = 0

	for i := 1; i < nSteps; i++ {
		predDisp := disp[i-1] + dt*vel[i-1] + 0.5*dt*dt*(1-2*beta)*acc[i-1]
		predVel := vel[i-1] + dt*(1-gamma)*acc[i-1]
		predAcc := wave[i]

		deltaDisp := (predAcc - c*predVel - stiffness*predDisp) / kPrime
		disp[i] = predDisp + deltaDisp
		vel[i] = predVel + gamma/beta*deltaDisp/dt
		acc[i] = predAcc + (gamma/beta-1)*deltaDisp/(dt*dt)

		if math.Abs(disp[i]) > maxDisp {
			maxDisp = math.Abs(disp[i])
		}
		if math.Abs(acc[i]) > maxAcc {
			maxAcc = math.Abs(acc[i])
		}
	}

	return
}

func (s *EarthquakeSimulator) RunDynamicAnalysis(req SimulationRequest) map[string]float64 {
	pga := s.CalcPGA(req.Magnitude)

	type structParams struct {
		name  string
		mass  float64
		stiff float64
		damp  float64
	}

	structures := []structParams{
		{"yuzui", 2.5e6, 1.2e9, 0.05},
		{"feishayan", 1.8e6, 8.5e8, 0.05},
		{"baopingkou", 3.0e6, 1.5e9, 0.05},
		{"renzidi", 1.2e6, 6.0e8, 0.05},
	}

	dt := 0.005
	nSteps := int(req.DurationSeconds / dt)
	earthquakeWave := make([]float64, nSteps)
	for i := 0; i < nSteps; i++ {
		t := float64(i) * dt
		freq := 2.0 + 3.0*rand.Float64()
		envelope := 1.0
		if t < 2.0 {
			envelope = t / 2.0
		} else if t > req.DurationSeconds-2.0 {
			envelope = (req.DurationSeconds - t) / 2.0
		} else if t > req.DurationSeconds {
			envelope = 0
		}
		earthquakeWave[i] = pga * 9.81 * envelope * (0.6*math.Sin(2*math.Pi*freq*t) + 0.4*rand.NormFloat64())
	}

	results := make(map[string]float64)
	var mu sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(structures))

	for _, st := range structures {
		go func(name string, mass, stiff, damp float64) {
			defer wg.Done()
			maxDisp, _ := s.NewmarkBetaIntegration(mass, stiff, damp, earthquakeWave, dt)

			omega := math.Sqrt(stiff / mass)
			dispLimit := stiff / (omega * omega * mass) * 0.1
			damageDynamic := math.Min(1.0, maxDisp/dispLimit)

			mu.Lock()
			results[name] = damageDynamic
			mu.Unlock()
		}(st.name, st.mass, st.stiff, st.damp)
	}

	wg.Wait()
	return results
}

func (s *EarthquakeSimulator) Run(req SimulationRequest) (*SimulationResult, error) {
	if req.FocalMechanism == "" {
		req.FocalMechanism = "strike-slip"
	}

	pga := s.CalcPGA(req.Magnitude)

	yuzuiLoc := StructureLocation{X: 100, Y: 50, Z: 730}
	feishayanLoc := StructureLocation{X: 150, Y: 30, Z: 728}
	baopingkouLoc := StructureLocation{X: 200, Y: 60, Z: 725}
	renzidiLoc := StructureLocation{X: 80, Y: 80, Z: 727}

	yuzuiDamage := s.CalcDamage(pga, yuzuiLoc, req.EpicenterX, req.EpicenterY)
	feishayanDamage := s.CalcDamage(pga, feishayanLoc, req.EpicenterX, req.EpicenterY)
	baopingkouDamage := s.CalcDamage(pga, baopingkouLoc, req.EpicenterX, req.EpicenterY)
	renzidiDamage := s.CalcDamage(pga, renzidiLoc, req.EpicenterX, req.EpicenterY)

	dynamicEnabled := req.Magnitude >= 6.0

	if dynamicEnabled {
		dynResults := s.RunDynamicAnalysis(req)
		yuzuiDamage = math.Max(yuzuiDamage, dynResults["yuzui"])
		feishayanDamage = math.Max(feishayanDamage, dynResults["feishayan"])
		baopingkouDamage = math.Max(baopingkouDamage, dynResults["baopingkou"])
		renzidiDamage = math.Max(renzidiDamage, dynResults["renzidi"])
	}

	safetyAssessment := s.AssessSafety(yuzuiDamage, feishayanDamage, baopingkouDamage, renzidiDamage)

	maxDamage := math.Max(math.Max(yuzuiDamage, feishayanDamage), math.Max(baopingkouDamage, renzidiDamage))
	sedimentDisturbance, flowPathChange, waterDiversionChange := s.CalcSecondaryEffects(maxDamage)

	timeSeriesData := s.GenerateTimeSeries(pga, req.DurationSeconds)

	timeStep := 0.01
	timeSeriesJSON := make([]map[string]interface{}, len(timeSeriesData))
	for i, acc := range timeSeriesData {
		t := float64(i) * timeStep
		timeSeriesJSON[i] = map[string]interface{}{
			"time":         t,
			"acceleration": acc,
			"displacement": acc * 0.1,
		}
	}
	timeSeriesDataJSON, _ := json.Marshal(timeSeriesJSON)

	structureDamage := map[string]interface{}{
		"yuzui": map[string]interface{}{
			"damage_level":  yuzuiDamage,
			"cracks":        int(yuzuiDamage * 10),
			"settlement_mm": yuzuiDamage * 50,
		},
		"feishayan": map[string]interface{}{
			"damage_level":  feishayanDamage,
			"cracks":        int(feishayanDamage * 10),
			"settlement_mm": feishayanDamage * 50,
		},
		"baopingkou": map[string]interface{}{
			"damage_level":  baopingkouDamage,
			"cracks":        int(baopingkouDamage * 10),
			"settlement_mm": baopingkouDamage * 50,
		},
		"renzidi": map[string]interface{}{
			"damage_level":  renzidiDamage,
			"cracks":        int(renzidiDamage * 10),
			"settlement_mm": renzidiDamage * 50,
		},
	}
	structureDamageJSON, _ := json.Marshal(structureDamage)

	bankCollapse := "stable"
	if yuzuiDamage > 0.3 || feishayanDamage > 0.3 || baopingkouDamage > 0.3 || renzidiDamage > 0.3 {
		bankCollapse = "partial_collapse"
	}

	sim := &models.EarthquakeSimulation{
		SimulationName:       req.SimulationName,
		Magnitude:            req.Magnitude,
		EpicenterX:           req.EpicenterX,
		EpicenterY:           req.EpicenterY,
		EpicenterZ:           req.EpicenterZ,
		FocalMechanism:       req.FocalMechanism,
		PGA:                  pga,
		DurationSeconds:      req.DurationSeconds,
		TimeSeriesData:       string(timeSeriesDataJSON),
		StructureDamage:      string(structureDamageJSON),
		YuzuiDamage:          yuzuiDamage,
		FeishayanDamage:      feishayanDamage,
		BaopingkouDamage:     baopingkouDamage,
		RenzidiDamage:        renzidiDamage,
		BankCollapse:         bankCollapse,
		SedimentDisturbance:  sedimentDisturbance,
		FlowPathChange:       flowPathChange,
		WaterDiversionChange: waterDiversionChange,
		SafetyAssessment:     safetyAssessment,
		CreatedBy:            req.CreatedBy,
	}

	id, err := models.InsertEarthquakeSimulation(s.ctx, sim)
	if err != nil {
		return nil, err
	}
	sim.ID = id

	return &SimulationResult{
		Simulation:             sim,
		DynamicAnalysisEnabled: dynamicEnabled,
		ParallelWorkers:        s.parallelWorkers,
	}, nil
}
