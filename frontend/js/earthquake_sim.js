(function(global) {
class EarthquakeSimulation {
    constructor() {
        this.canvas = null;
        this.ctx = null;
        this.running = false;
        this.animationId = null;
        this.timeStep = 0;
        this.simulationData = null;
        this.currentDamage = {
            yuzui: 0,
            feishayan: 0,
            baopingkou: 0,
            renzidi: 0
        };
        this.pga = 0;
        this.gauges = {};
        this.waveChart = null;
        this.cracks = [];
        this.shakeOffset = { x: 0, y: 0 };
        this.init();
    }

    init() {
        this.initControls();
        this.initSceneCanvas();
        this.initGauges();
        this.initWaveChart();
        this.bindEvents();
    }

    initControls() {
        const controlsContainer = document.querySelector('.earthquake-controls');
        if (!controlsContainer) return;

        const presetGroup = document.createElement('div');
        presetGroup.className = 'control-group';
        presetGroup.innerHTML = `
            <label>地震预设</label>
            <select id="eq-preset">
                <option value="">-- 选择预设 --</option>
                ${Object.entries(CONFIG.EARTHQUAKE_PRESETS).map(([key, preset]) => 
                    `<option value="${key}">${preset.name}</option>`
                ).join('')}
            </select>
        `;
        controlsContainer.insertBefore(presetGroup, controlsContainer.firstChild);

        const durationGroup = document.createElement('div');
        durationGroup.className = 'control-group';
        durationGroup.innerHTML = `
            <label>持续时间 (秒)</label>
            <input type="number" id="eq-duration" value="60" min="1" max="300" step="1">
        `;
        controlsContainer.insertBefore(durationGroup, controlsContainer.querySelector('#run-earthquake'));

        const latGroup = document.createElement('div');
        latGroup.className = 'control-group';
        latGroup.innerHTML = `
            <label>纬度 (°N)</label>
            <input type="number" id="eq-lat" value="31.0" min="0" max="90" step="0.1">
        `;
        controlsContainer.insertBefore(latGroup, controlsContainer.querySelector('#run-earthquake'));

        const lngGroup = document.createElement('div');
        lngGroup.className = 'control-group';
        lngGroup.innerHTML = `
            <label>经度 (°E)</label>
            <input type="number" id="eq-lng" value="103.6" min="0" max="180" step="0.1">
        `;
        controlsContainer.insertBefore(lngGroup, controlsContainer.querySelector('#run-earthquake'));

        const distanceInput = document.getElementById('eq-distance');
        if (distanceInput) {
            distanceInput.remove();
        }
        const mechanismInput = document.getElementById('eq-mechanism');
        if (mechanismInput && mechanismInput.parentElement) {
            mechanismInput.parentElement.remove();
        }
    }

    initSceneCanvas() {
        const sceneContainer = document.querySelector('.earthquake-scene');
        if (!sceneContainer) return;

        const oldContainer = document.getElementById('earthquake-3d-container');
        if (oldContainer) {
            oldContainer.remove();
        }

        const title = sceneContainer.querySelector('h4');
        if (title) {
            title.textContent = '都江堰渠首震害模拟';
        }

        this.canvas = document.createElement('canvas');
        this.canvas.id = 'earthquake-scene-canvas';
        this.canvas.className = 'sim-canvas';
        this.canvas.style.height = '450px';
        sceneContainer.appendChild(this.canvas);

        this.ctx = this.canvas.getContext('2d');
        this.resizeCanvas();

        this.drawScene();
    }

    resizeCanvas() {
        if (!this.canvas) return;
        const rect = this.canvas.parentElement.getBoundingClientRect();
        this.canvas.width = rect.width;
        this.canvas.height = 450;
        this.width = this.canvas.width;
        this.height = this.canvas.height;
    }

    initGauges() {
        const gaugeIds = ['gauge-yuzui', 'gauge-feishayan', 'gauge-baopingkou', 'gauge-renzidi'];
        gaugeIds.forEach(id => {
            const canvas = document.getElementById(id);
            if (canvas) {
                const key = id.replace('gauge-', '');
                this.gauges[key] = {
                    canvas: canvas,
                    ctx: canvas.getContext('2d'),
                    value: 0
                };
                this.resizeGauge(canvas);
                this.drawGauge(key, 0);
            }
        });
    }

    resizeGauge(canvas) {
        const rect = canvas.parentElement.getBoundingClientRect();
        canvas.width = rect.width || 200;
        canvas.height = 120;
    }

    drawGauge(key, value) {
        const gauge = this.gauges[key];
        if (!gauge) return;

        const { canvas, ctx } = gauge;
        const w = canvas.width;
        const h = canvas.height;
        const centerX = w / 2;
        const centerY = h * 0.85;
        const radius = Math.min(w, h * 1.8) * 0.4;

        ctx.clearRect(0, 0, w, h);

        const startAngle = Math.PI;
        const endAngle = 2 * Math.PI;

        ctx.beginPath();
        ctx.arc(centerX, centerY, radius, startAngle, endAngle);
        ctx.strokeStyle = 'rgba(255, 255, 255, 0.1)';
        ctx.lineWidth = 15;
        ctx.lineCap = 'round';
        ctx.stroke();

        const valueAngle = startAngle + (endAngle - startAngle) * Utils.clamp(value, 0, 1);
        const gradient = ctx.createLinearGradient(centerX - radius, centerY, centerX + radius, centerY);
        gradient.addColorStop(0, '#00ff88');
        gradient.addColorStop(0.5, '#ffaa00');
        gradient.addColorStop(1, '#ff4444');

        ctx.beginPath();
        ctx.arc(centerX, centerY, radius, startAngle, valueAngle);
        ctx.strokeStyle = gradient;
        ctx.lineWidth = 15;
        ctx.lineCap = 'round';
        ctx.stroke();

        const tickCount = 11;
        for (let i = 0; i < tickCount; i++) {
            const t = i / (tickCount - 1);
            const angle = startAngle + (endAngle - startAngle) * t;
            const innerR = radius - 8;
            const outerR = radius + 8;
            const x1 = centerX + Math.cos(angle) * innerR;
            const y1 = centerY + Math.sin(angle) * innerR;
            const x2 = centerX + Math.cos(angle) * outerR;
            const y2 = centerY + Math.sin(angle) * outerR;

            ctx.beginPath();
            ctx.moveTo(x1, y1);
            ctx.lineTo(x2, y2);
            ctx.strokeStyle = 'rgba(255, 255, 255, 0.3)';
            ctx.lineWidth = 1;
            ctx.stroke();
        }

        const needleAngle = valueAngle;
        const needleLength = radius - 20;
        const needleX = centerX + Math.cos(needleAngle) * needleLength;
        const needleY = centerY + Math.sin(needleAngle) * needleLength;

        ctx.beginPath();
        ctx.moveTo(centerX, centerY);
        ctx.lineTo(needleX, needleY);
        ctx.strokeStyle = '#ffffff';
        ctx.lineWidth = 2;
        ctx.lineCap = 'round';
        ctx.stroke();

        ctx.beginPath();
        ctx.arc(centerX, centerY, 5, 0, Math.PI * 2);
        ctx.fillStyle = '#ffffff';
        ctx.fill();

        ctx.fillStyle = '#ffffff';
        ctx.font = 'bold 18px Arial';
        ctx.textAlign = 'center';
        ctx.fillText(Math.round(value * 100) + '%', centerX, centerY - 10);

        const valueEl = document.getElementById('damage-' + key);
        if (valueEl) {
            valueEl.textContent = Math.round(value * 100) + ' %';
        }
    }

    initWaveChart() {
        const resultsContainer = document.querySelector('.earthquake-results');
        if (!resultsContainer) return;

        const waveCard = document.createElement('div');
        waveCard.className = 'result-card wide';
        waveCard.style.marginTop = '20px';
        waveCard.innerHTML = `
            <h4>地震波时程曲线 (P波 / S波)</h4>
            <canvas id="earthquake-wave-chart"></canvas>
        `;
        resultsContainer.appendChild(waveCard);

        const ctx = document.getElementById('earthquake-wave-chart');
        if (!ctx) return;

        const labels = [];
        for (let i = 0; i <= 10; i += 0.1) {
            labels.push(i.toFixed(1) + 's');
        }

        this.waveChart = new Chart(ctx, {
            type: 'line',
            data: {
                labels: labels,
                datasets: [
                    {
                        label: 'P波 (纵波)',
                        data: new Array(labels.length).fill(0),
                        borderColor: '#00d4ff',
                        backgroundColor: 'rgba(0, 212, 255, 0.1)',
                        tension: 0.3,
                        fill: true,
                        pointRadius: 0,
                        borderWidth: 2
                    },
                    {
                        label: 'S波 (横波)',
                        data: new Array(labels.length).fill(0),
                        borderColor: '#ff66cc',
                        backgroundColor: 'rgba(255, 102, 204, 0.1)',
                        tension: 0.3,
                        fill: true,
                        pointRadius: 0,
                        borderWidth: 2
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        grid: { color: 'rgba(255,255,255,0.1)' },
                        ticks: { color: '#aaa' },
                        title: {
                            display: true,
                            text: '加速度 (gal)',
                            color: '#8892a0'
                        }
                    },
                    x: {
                        grid: { color: 'rgba(255,255,255,0.1)' },
                        ticks: { color: '#aaa', maxTicksLimit: 11 }
                    }
                },
                plugins: {
                    legend: { labels: { color: '#ccc' } }
                },
                animation: {
                    duration: 0
                }
            }
        });
    }

    generateWaveData(magnitude, duration) {
        const labels = [];
        const pWaveData = [];
        const sWaveData = [];
        const totalPoints = 101;
        const maxAccel = Math.pow(10, magnitude - 4) * 20;

        for (let i = 0; i < totalPoints; i++) {
            const t = i * 0.1;
            labels.push(t.toFixed(1) + 's');

            const pWaveStart = 0.5;
            const pWaveDecay = Math.exp(-Math.max(0, t - pWaveStart) * 0.3);
            const pWaveAmplitude = t >= pWaveStart ? Math.sin((t - pWaveStart) * 8) * maxAccel * 0.4 * pWaveDecay : 0;
            const pWaveNoise = t >= pWaveStart ? (Math.random() - 0.5) * maxAccel * 0.2 * pWaveDecay : 0;
            pWaveData.push(pWaveAmplitude + pWaveNoise);

            const sWaveStart = 1.5;
            const sWaveDecay = Math.exp(-Math.max(0, t - sWaveStart) * 0.2);
            const sWaveAmplitude = t >= sWaveStart ? Math.sin((t - sWaveStart) * 5) * maxAccel * 0.8 * sWaveDecay : 0;
            const sWaveNoise = t >= sWaveStart ? (Math.random() - 0.5) * maxAccel * 0.4 * sWaveDecay : 0;
            sWaveData.push(sWaveAmplitude + sWaveNoise);
        }

        return { labels, pWaveData, sWaveData, maxAccel };
    }

    updateWaveChart(pWaveData, sWaveData) {
        if (!this.waveChart) return;
        this.waveChart.data.datasets[0].data = pWaveData;
        this.waveChart.data.datasets[1].data = sWaveData;
        this.waveChart.update('none');
    }

    bindEvents() {
        const presetSelect = document.getElementById('eq-preset');
        if (presetSelect) {
            presetSelect.addEventListener('change', (e) => {
                this.applyPreset(e.target.value);
            });
        }

        const runBtn = document.getElementById('run-earthquake');
        if (runBtn) {
            runBtn.addEventListener('click', () => {
                this.runSimulation();
            });
        }

        window.addEventListener('resize', () => {
            this.resizeCanvas();
            Object.values(this.gauges).forEach(g => this.resizeGauge(g.canvas));
            this.drawScene();
            Object.keys(this.gauges).forEach(k => this.drawGauge(k, this.currentDamage[k] || 0));
        });
    }

    applyPreset(presetKey) {
        if (!presetKey || !CONFIG.EARTHQUAKE_PRESETS[presetKey]) return;

        const preset = CONFIG.EARTHQUAKE_PRESETS[presetKey];

        const magnitudeInput = document.getElementById('eq-magnitude');
        const depthInput = document.getElementById('eq-depth');
        const latInput = document.getElementById('eq-lat');
        const lngInput = document.getElementById('eq-lng');
        const durationInput = document.getElementById('eq-duration');

        if (magnitudeInput) magnitudeInput.value = preset.magnitude;
        if (depthInput) depthInput.value = preset.depth_km;
        if (latInput) latInput.value = preset.lat;
        if (lngInput) lngInput.value = preset.lng;
        if (durationInput) durationInput.value = preset.duration_sec;
    }

    getParams() {
        return {
            magnitude: parseFloat(document.getElementById('eq-magnitude')?.value || 7.0),
            depth_km: parseFloat(document.getElementById('eq-depth')?.value || 15),
            lat: parseFloat(document.getElementById('eq-lat')?.value || 31.0),
            lng: parseFloat(document.getElementById('eq-lng')?.value || 103.6),
            duration_sec: parseInt(document.getElementById('eq-duration')?.value || 60)
        };
    }

    async runSimulation() {
        this.stop();
        this.timeStep = 0;
        this.cracks = [];
        this.currentDamage = { yuzui: 0, feishayan: 0, baopingkou: 0, renzidi: 0 };
        Object.keys(this.gauges).forEach(k => this.drawGauge(k, 0));

        const params = this.getParams();

        try {
            const result = await API.earthquake.simulate(params);
            this.processSimulationResult(result, params);
        } catch (error) {
            console.error('地震模拟API失败，使用Mock数据:', error);
            const mockResult = this.generateMockResult(params);
            this.processSimulationResult(mockResult, params);
        }
    }

    generateMockResult(params) {
        const { magnitude } = params;
        let baseDamage;

        if (magnitude >= 8.0) {
            baseDamage = 0.7 + Math.random() * 0.2;
        } else if (magnitude >= 7.0) {
            baseDamage = 0.5 + Math.random() * 0.2;
        } else if (magnitude >= 6.5) {
            baseDamage = 0.3 + Math.random() * 0.2;
        } else if (magnitude >= 5.0) {
            baseDamage = 0.1 + Math.random() * 0.2;
        } else {
            baseDamage = Math.random() * 0.1;
        }

        const variance = 0.1;
        return {
            yuzui_damage: Utils.clamp(baseDamage + (Math.random() - 0.5) * variance, 0, 1),
            feishayan_damage: Utils.clamp(baseDamage * 0.9 + (Math.random() - 0.5) * variance, 0, 1),
            baopingkou_damage: Utils.clamp(baseDamage * 0.7 + (Math.random() - 0.5) * variance, 0, 1),
            renzidi_damage: Utils.clamp(baseDamage * 0.85 + (Math.random() - 0.5) * variance, 0, 1),
            pga: Math.pow(10, magnitude - 4) * 15,
            safety_assessment: this.calculateSafetyAssessment(baseDamage),
            repair_duration_days: Math.round(baseDamage * 180),
            repair_cost_10k_yuan: Math.round(baseDamage * 5000),
            recommendations: this.generateRecommendations(baseDamage)
        };
    }

    calculateSafetyAssessment(damage) {
        if (damage < 0.2) return 'safe';
        if (damage < 0.5) return 'caution';
        return 'danger';
    }

    generateRecommendations(damage) {
        if (damage < 0.2) {
            return '结构整体安全，建议进行常规巡查，重点检查附属设施。';
        } else if (damage < 0.5) {
            return '建议立即开展专业检测，对受损部位进行临时加固，限制重载通行，制定修复方案。';
        } else {
            return '结构存在重大安全隐患，建议立即关闭相关区域，启动应急预案，组织专家论证抢险方案。';
        }
    }

    processSimulationResult(result, params) {
        this.simulationData = result;
        this.pga = result.pga || Math.pow(10, params.magnitude - 4) * 15;

        const waveData = this.generateWaveData(params.magnitude, params.duration_sec);
        this.updateWaveChart(waveData.pWaveData, waveData.sWaveData);
        this.pga = waveData.maxAccel;

        this.generateCracks(Math.max(
            result.yuzui_damage,
            result.feishayan_damage,
            result.baopingkou_damage,
            result.renzidi_damage
        ));

        this.updateSafetyAssessment(result);

        this.targetDamage = {
            yuzui: result.yuzui_damage || 0,
            feishayan: result.feishayan_damage || 0,
            baopingkou: result.baopingkou_damage || 0,
            renzidi: result.renzidi_damage || 0
        };

        this.running = true;
        this.animate();
    }

    generateCracks(damageLevel) {
        this.cracks = [];
        const crackCount = Math.floor(damageLevel * 80) + 10;

        const regions = [
            { xMin: 0.35, xMax: 0.45, yMin: 0.35, yMax: 0.55 },
            { xMin: 0.5, xMax: 0.65, yMin: 0.5, yMax: 0.7 },
            { xMin: 0.7, xMax: 0.8, yMin: 0.3, yMax: 0.5 },
            { xMin: 0.55, xMax: 0.7, yMin: 0.25, yMax: 0.4 }
        ];

        for (let i = 0; i < crackCount; i++) {
            const region = regions[Math.floor(Math.random() * regions.length)];
            const baseX = Utils.lerp(region.xMin, region.xMax, Math.random());
            const baseY = Utils.lerp(region.yMin, region.yMax, Math.random());
            const points = [{ x: baseX, y: baseY }];
            const segments = 3 + Math.floor(Math.random() * 5);
            let x = baseX;
            let y = baseY;
            const angle = Math.random() * Math.PI * 2;

            for (let j = 0; j < segments; j++) {
                const segLen = (0.02 + Math.random() * 0.05) * (0.5 + damageLevel);
                const segAngle = angle + (Math.random() - 0.5) * 1.2;
                x += Math.cos(segAngle) * segLen;
                y += Math.sin(segAngle) * segLen;
                x = Utils.clamp(x, 0, 1);
                y = Utils.clamp(y, 0, 1);
                points.push({ x, y });
            }

            this.cracks.push({
                points,
                width: (1 + Math.random() * 3) * (0.5 + damageLevel),
                opacity: 0.4 + Math.random() * 0.6
            });
        }
    }

    updateSafetyAssessment(result) {
        const levelEl = document.getElementById('safety-level');
        const durationEl = document.getElementById('repair-duration');
        const costEl = document.getElementById('repair-cost');
        const recEl = document.getElementById('safety-recommendation');
        const assessmentCard = document.getElementById('safety-assessment');

        const level = result.safety_assessment || 'safe';
        const levelText = { safe: '安全', caution: '注意', danger: '危险' }[level];
        const levelColor = { safe: '#00ff88', caution: '#ffaa00', danger: '#ff4444' }[level];

        if (levelEl) {
            levelEl.textContent = levelText;
            levelEl.style.color = levelColor;
            levelEl.style.fontWeight = 'bold';
            levelEl.style.fontSize = '20px';
        }

        if (durationEl) {
            durationEl.textContent = (result.repair_duration_days || 0) + ' 天';
        }

        if (costEl) {
            costEl.textContent = (result.repair_cost_10k_yuan || 0) + ' 万元';
        }

        if (recEl) {
            recEl.textContent = result.recommendations || '--';
        }

        if (assessmentCard) {
            assessmentCard.style.borderLeft = `4px solid ${levelColor}`;
            assessmentCard.style.background = levelColor.replace(')', ', 0.05)').replace('rgb', 'rgba').replace('#', '');
            const bgMap = {
                safe: 'rgba(0, 255, 136, 0.05)',
                caution: 'rgba(255, 170, 0, 0.05)',
                danger: 'rgba(255, 68, 68, 0.05)'
            };
            assessmentCard.style.background = bgMap[level];
        }
    }

    animate() {
        if (!this.running) return;

        this.timeStep++;

        const progress = Math.min(this.timeStep / 180, 1);
        Object.keys(this.targetDamage).forEach(key => {
            this.currentDamage[key] = Utils.lerp(this.currentDamage[key] || 0, this.targetDamage[key], 0.03);
            this.drawGauge(key, this.currentDamage[key]);
        });

        if (progress < 1) {
            const shakeIntensity = this.pga / 200;
            this.shakeOffset.x = (Math.random() - 0.5) * this.pga * 0.4;
            this.shakeOffset.y = (Math.random() - 0.5) * this.pga * 0.4;
        } else {
            this.shakeOffset.x *= 0.9;
            this.shakeOffset.y *= 0.9;
        }

        this.drawScene();

        if (progress < 1 || Math.abs(this.shakeOffset.x) > 0.1 || Math.abs(this.shakeOffset.y) > 0.1) {
            this.animationId = requestAnimationFrame(() => this.animate());
        } else {
            this.running = false;
        }
    }

    drawScene() {
        if (!this.ctx) return;

        const ctx = this.ctx;
        const w = this.width;
        const h = this.height;

        ctx.save();
        ctx.clearRect(0, 0, w, h);

        ctx.translate(this.shakeOffset.x, this.shakeOffset.y);

        this.drawBackground(ctx, w, h);
        this.drawRiverSystem(ctx, w, h);
        this.drawStructures(ctx, w, h);
        this.drawCracks(ctx, w, h);
        this.drawLabels(ctx, w, h);

        ctx.restore();
    }

    drawBackground(ctx, w, h) {
        const gradient = ctx.createLinearGradient(0, 0, w, h);
        gradient.addColorStop(0, '#1a2a4a');
        gradient.addColorStop(0.5, '#2d4a6a');
        gradient.addColorStop(1, '#1a3a2a');
        ctx.fillStyle = gradient;
        ctx.fillRect(0, 0, w, h);

        ctx.strokeStyle = 'rgba(0, 212, 255, 0.08)';
        ctx.lineWidth = 1;
        for (let x = 0; x <= w; x += 40) {
            ctx.beginPath();
            ctx.moveTo(x, 0);
            ctx.lineTo(x, h);
            ctx.stroke();
        }
        for (let y = 0; y <= h; y += 40) {
            ctx.beginPath();
            ctx.moveTo(0, y);
            ctx.lineTo(w, y);
            ctx.stroke();
        }

        ctx.fillStyle = CONFIG.COLORS.terrain;
        ctx.globalAlpha = 0.3;
        ctx.beginPath();
        ctx.moveTo(0, h * 0.15);
        ctx.quadraticCurveTo(w * 0.1, h * 0.05, w * 0.2, h * 0.12);
        ctx.quadraticCurveTo(w * 0.35, h * 0.02, w * 0.5, h * 0.1);
        ctx.quadraticCurveTo(w * 0.65, h * 0.03, w * 0.8, h * 0.08);
        ctx.quadraticCurveTo(w * 0.9, h * 0.04, w, h * 0.15);
        ctx.lineTo(w, 0);
        ctx.lineTo(0, 0);
        ctx.closePath();
        ctx.fill();
        ctx.globalAlpha = 1;
    }

    drawRiverSystem(ctx, w, h) {
        ctx.lineCap = 'round';
        ctx.lineJoin = 'round';

        const waterGradient = ctx.createLinearGradient(0, h * 0.3, 0, h * 0.9);
        waterGradient.addColorStop(0, '#0066aa');
        waterGradient.addColorStop(1, '#00a8ff');

        ctx.strokeStyle = waterGradient;
        ctx.lineWidth = 70;
        ctx.beginPath();
        ctx.moveTo(-20, h * 0.5);
        ctx.bezierCurveTo(w * 0.25, h * 0.45, w * 0.3, h * 0.55, w * 0.4, h * 0.5);
        ctx.stroke();

        ctx.strokeStyle = 'rgba(0, 212, 255, 0.3)';
        ctx.lineWidth = 74;
        ctx.beginPath();
        ctx.moveTo(-20, h * 0.5);
        ctx.bezierCurveTo(w * 0.25, h * 0.45, w * 0.3, h * 0.55, w * 0.4, h * 0.5);
        ctx.stroke();

        ctx.strokeStyle = waterGradient;
        ctx.lineWidth = 45;
        ctx.beginPath();
        ctx.moveTo(w * 0.4, h * 0.5);
        ctx.bezierCurveTo(w * 0.5, h * 0.42, w * 0.6, h * 0.38, w * 0.72, h * 0.4);
        ctx.bezierCurveTo(w * 0.8, h * 0.42, w * 0.9, h * 0.38, w + 20, h * 0.35);
        ctx.stroke();

        ctx.strokeStyle = 'rgba(0, 212, 255, 0.3)';
        ctx.lineWidth = 49;
        ctx.beginPath();
        ctx.moveTo(w * 0.4, h * 0.5);
        ctx.bezierCurveTo(w * 0.5, h * 0.42, w * 0.6, h * 0.38, w * 0.72, h * 0.4);
        ctx.bezierCurveTo(w * 0.8, h * 0.42, w * 0.9, h * 0.38, w + 20, h * 0.35);
        ctx.stroke();

        ctx.strokeStyle = waterGradient;
        ctx.lineWidth = 50;
        ctx.beginPath();
        ctx.moveTo(w * 0.4, h * 0.5);
        ctx.bezierCurveTo(w * 0.5, h * 0.6, w * 0.6, h * 0.62, w * 0.7, h * 0.6);
        ctx.bezierCurveTo(w * 0.8, h * 0.58, w * 0.9, h * 0.65, w + 20, h * 0.7);
        ctx.stroke();

        ctx.strokeStyle = 'rgba(0, 212, 255, 0.3)';
        ctx.lineWidth = 54;
        ctx.beginPath();
        ctx.moveTo(w * 0.4, h * 0.5);
        ctx.bezierCurveTo(w * 0.5, h * 0.6, w * 0.6, h * 0.62, w * 0.7, h * 0.6);
        ctx.bezierCurveTo(w * 0.8, h * 0.58, w * 0.9, h * 0.65, w + 20, h * 0.7);
        ctx.stroke();
    }

    drawStructures(ctx, w, h) {
        this.drawYuzui(ctx, w * 0.4, h * 0.5, this.currentDamage.yuzui);
        this.drawFeishayan(ctx, w * 0.56, h * 0.5, this.currentDamage.feishayan);
        this.drawBaopingkou(ctx, w * 0.75, h * 0.38, this.currentDamage.baopingkou);
        this.drawRenzidi(ctx, w * 0.62, h * 0.32, this.currentDamage.renzidi);
        this.drawJingangdi(ctx, w, h);
    }

    drawYuzui(ctx, cx, cy, damage) {
        ctx.save();

        const damageColor = Utils.getElevationColor(1 - damage, 0, 1);

        ctx.fillStyle = CONFIG.COLORS.stone;
        ctx.strokeStyle = '#666';
        ctx.lineWidth = 2;

        ctx.beginPath();
        ctx.moveTo(cx - 45, cy - 5);
        ctx.quadraticCurveTo(cx - 55, cy, cx - 45, cy + 5);
        ctx.lineTo(cx - 20, cy + 20);
        ctx.lineTo(cx + 50, cy + 15);
        ctx.quadraticCurveTo(cx + 65, cy, cx + 50, cy - 15);
        ctx.lineTo(cx - 20, cy - 20);
        ctx.closePath();
        ctx.fill();
        ctx.stroke();

        if (damage > 0.2) {
            ctx.fillStyle = damageColor;
            ctx.globalAlpha = damage * 0.5;
            ctx.fill();
            ctx.globalAlpha = 1;
        }

        ctx.fillStyle = '#555';
        for (let i = 0; i < 8; i++) {
            const px = cx - 35 + i * 12;
            const py = cy - 10 + Math.sin(i) * 5;
            ctx.beginPath();
            ctx.arc(px, py, 3, 0, Math.PI * 2);
            ctx.fill();
        }

        ctx.restore();
    }

    drawFeishayan(ctx, cx, cy, damage) {
        ctx.save();

        const damageColor = Utils.getElevationColor(1 - damage, 0, 1);

        ctx.fillStyle = CONFIG.COLORS.sand;
        ctx.strokeStyle = '#a08060';
        ctx.lineWidth = 2;

        ctx.beginPath();
        ctx.moveTo(cx - 5, cy - 25);
        ctx.lineTo(cx + 55, cy - 30);
        ctx.lineTo(cx + 60, cy + 10);
        ctx.lineTo(cx, cy + 15);
        ctx.closePath();
        ctx.fill();
        ctx.stroke();

        if (damage > 0.2) {
            ctx.fillStyle = damageColor;
            ctx.globalAlpha = damage * 0.5;
            ctx.fill();
            ctx.globalAlpha = 1;
        }

        ctx.strokeStyle = CONFIG.COLORS.stone;
        ctx.lineWidth = 4;
        ctx.beginPath();
        ctx.moveTo(cx - 5, cy - 8);
        ctx.lineTo(cx + 55, cy - 13);
        ctx.stroke();

        ctx.restore();
    }

    drawBaopingkou(ctx, cx, cy, damage) {
        ctx.save();

        const damageColor = Utils.getElevationColor(1 - damage, 0, 1);

        ctx.fillStyle = '#7a6a5a';
        ctx.strokeStyle = '#5a4a3a';
        ctx.lineWidth = 2;

        ctx.beginPath();
        ctx.moveTo(cx - 25, cy - 40);
        ctx.lineTo(cx - 10, cy - 50);
        ctx.lineTo(cx + 15, cy - 45);
        ctx.lineTo(cx + 30, cy - 30);
        ctx.lineTo(cx + 25, cy + 10);
        ctx.lineTo(cx - 30, cy + 5);
        ctx.closePath();
        ctx.fill();
        ctx.stroke();

        if (damage > 0.2) {
            ctx.fillStyle = damageColor;
            ctx.globalAlpha = damage * 0.5;
            ctx.fill();
            ctx.globalAlpha = 1;
        }

        const waterGradient = ctx.createLinearGradient(cx - 15, cy - 35, cx - 15, cy);
        waterGradient.addColorStop(0, '#0088cc');
        waterGradient.addColorStop(1, '#00a8ff');
        ctx.fillStyle = waterGradient;
        ctx.beginPath();
        ctx.moveTo(cx - 18, cy - 30);
        ctx.quadraticCurveTo(cx - 5, cy - 38, cx + 10, cy - 32);
        ctx.lineTo(cx + 15, cy);
        ctx.lineTo(cx - 22, cy - 3);
        ctx.closePath();
        ctx.fill();

        ctx.restore();
    }

    drawRenzidi(ctx, cx, cy, damage) {
        ctx.save();

        const damageColor = Utils.getElevationColor(1 - damage, 0, 1);

        ctx.fillStyle = CONFIG.COLORS.terrain;
        ctx.strokeStyle = '#5c4a3a';
        ctx.lineWidth = 2;

        ctx.beginPath();
        ctx.moveTo(cx - 40, cy + 35);
        ctx.lineTo(cx, cy - 30);
        ctx.lineTo(cx + 40, cy + 35);
        ctx.lineTo(cx + 30, cy + 45);
        ctx.lineTo(cx, cy - 15);
        ctx.lineTo(cx - 30, cy + 45);
        ctx.closePath();
        ctx.fill();
        ctx.stroke();

        if (damage > 0.2) {
            ctx.fillStyle = damageColor;
            ctx.globalAlpha = damage * 0.5;
            ctx.fill();
            ctx.globalAlpha = 1;
        }

        ctx.strokeStyle = CONFIG.COLORS.stone;
        ctx.lineWidth = 3;
        ctx.setLineDash([5, 3]);
        ctx.beginPath();
        ctx.moveTo(cx - 35, cy + 30);
        ctx.lineTo(cx, cy - 20);
        ctx.lineTo(cx + 35, cy + 30);
        ctx.stroke();
        ctx.setLineDash([]);

        ctx.restore();
    }

    drawJingangdi(ctx, w, h) {
        ctx.fillStyle = CONFIG.COLORS.terrainDark;
        ctx.strokeStyle = '#3a2a1a';
        ctx.lineWidth = 2;

        ctx.beginPath();
        ctx.moveTo(w * 0.4, h * 0.5);
        ctx.lineTo(w * 0.5, h * 0.42);
        ctx.lineTo(w * 0.52, h * 0.48);
        ctx.lineTo(w * 0.42, h * 0.56);
        ctx.closePath();
        ctx.fill();
        ctx.stroke();
    }

    drawCracks(ctx, w, h) {
        if (this.cracks.length === 0) return;

        const avgDamage = (this.currentDamage.yuzui + this.currentDamage.feishayan + 
                          this.currentDamage.baopingkou + this.currentDamage.renzidi) / 4;

        if (avgDamage < 0.05) return;

        ctx.save();

        this.cracks.forEach(crack => {
            if (crack.points.length < 2) return;

            ctx.strokeStyle = `rgba(180, 20, 20, ${crack.opacity * Utils.clamp(avgDamage * 1.5, 0, 1)})`;
            ctx.lineWidth = crack.width;
            ctx.lineCap = 'round';
            ctx.lineJoin = 'round';

            ctx.beginPath();
            ctx.moveTo(crack.points[0].x * w, crack.points[0].y * h);
            for (let i = 1; i < crack.points.length; i++) {
                ctx.lineTo(crack.points[i].x * w, crack.points[i].y * h);
            }
            ctx.stroke();
        });

        ctx.restore();
    }

    drawLabels(ctx, w, h) {
        ctx.font = 'bold 14px Microsoft YaHei, sans-serif';
        ctx.textAlign = 'center';

        const labels = [
            { text: '鱼嘴', x: w * 0.4, y: h * 0.5 - 40, color: '#00d4ff' },
            { text: '飞沙堰', x: w * 0.58, y: h * 0.5 - 45, color: '#ffaa00' },
            { text: '宝瓶口', x: w * 0.75, y: h * 0.38 - 60, color: '#00ff88' },
            { text: '人字堤', x: w * 0.62, y: h * 0.32 - 50, color: '#ff66cc' },
            { text: '内江', x: w * 0.65, y: h * 0.38, color: '#00a8ff' },
            { text: '外江', x: w * 0.8, y: h * 0.7, color: '#00a8ff' },
            { text: '岷江', x: w * 0.15, y: h * 0.42, color: '#00a8ff' }
        ];

        labels.forEach(label => {
            ctx.fillStyle = 'rgba(0, 0, 0, 0.7)';
            const metrics = ctx.measureText(label.text);
            const padX = 8;
            const padY = 4;
            ctx.fillRect(
                label.x - metrics.width / 2 - padX,
                label.y - 12 - padY,
                metrics.width + padX * 2,
                20 + padY * 2
            );

            ctx.fillStyle = label.color;
            ctx.fillText(label.text, label.x, label.y);
        });
    }

    stop() {
        this.running = false;
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
            this.animationId = null;
        }
    }

    reset() {
        this.stop();
        this.timeStep = 0;
        this.cracks = [];
        this.currentDamage = { yuzui: 0, feishayan: 0, baopingkou: 0, renzidi: 0 };
        this.shakeOffset = { x: 0, y: 0 };
        Object.keys(this.gauges).forEach(k => this.drawGauge(k, 0));
        this.drawScene();

        if (this.waveChart) {
            this.waveChart.data.datasets[0].data = new Array(101).fill(0);
            this.waveChart.data.datasets[1].data = new Array(101).fill(0);
            this.waveChart.update('none');
        }
    }
}

let earthquakeSim = null;

function initEarthquakeSimulation() {
    if (!earthquakeSim) {
        earthquakeSim = new EarthquakeSimulation();
    }
}

document.addEventListener('DOMContentLoaded', () => {
    if (document.querySelector('.earthquake-container')) {
        initEarthquakeSimulation();
    }
});

global.EarthquakeSimulator = EarthquakeSimulation;
global.initEarthquakeSimulation = initEarthquakeSimulation;
})(window);
