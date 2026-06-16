(function(global) {
class UserExperienceSimulation {
    constructor(canvasId) {
        this.canvas = document.getElementById(canvasId);
        this.ctx = this.canvas.getContext('2d');
        this.running = false;
        this.animationId = null;
        this.timeStep = 0;

        this.sessionId = null;
        this.nickname = '';
        this.startTime = null;
        this.elapsedTime = 0;
        this.timerInterval = null;

        this.placedObjects = [];
        this.dredgeAreas = [];
        this.selectedTool = null;
        this.selectedObject = null;
        this.isDragging = false;
        this.dragOffset = { x: 0, y: 0 };

        // ===== 竹笼编织引导系统 =====
        this.tutorialMode = true;  // 默认开启新手教程
        this.bambooWeavingSteps = [
            {
                id: 1,
                title: '第一步：选竹',
                description: '选取3年生慈竹或毛竹，直径10-15cm，直度好，无虫蛀',
                canvasHint: '点击左侧竹料堆选择合适的竹子',
                targetAction: 'select_bamboo',
                completed: false
            },
            {
                id: 2,
                title: '第二步：破篾',
                description: '用篾刀将竹材破成宽约6分(2cm)、厚约1分(3mm)的竹篾',
                canvasHint: '按住鼠标左键沿竹子纵向拖动破篾',
                targetAction: 'cut_bamboo',
                completed: false
            },
            {
                id: 3,
                title: '第三步：编笼',
                description: '采用"六角编"技法，先编底、再编身、最后收编口部，形成圆筒形',
                canvasHint: '点击竹篾交叉点进行编织，形成圆筒形状',
                targetAction: 'weave_cage',
                completed: false
            },
            {
                id: 4,
                title: '第四步：装石',
                description: '将卵石装入竹笼，大面朝下，级配填充，每笼约装卵石80-120个',
                canvasHint: '拖动卵石到竹笼中进行装填',
                targetAction: 'fill_stones',
                completed: false
            },
            {
                id: 5,
                title: '第五步：封口',
                description: '用木楔加固笼口，篾条收紧扎牢，确保运输和投放时不散架',
                canvasHint: '点击笼口进行封口加固',
                targetAction: 'seal_cage',
                completed: false
            },
            {
                id: 6,
                title: '第六步：投放',
                description: '竹笼运至指定位置，船载定位，人工辅助抛投入水',
                canvasHint: '点击投放按钮将竹笼放入江中',
                targetAction: 'deploy_cage',
                completed: false
            }
        ];
        this.currentBambooStep = 0;
        this.guideOverlayActive = false;
        this.showTutorial = true;
        this.tutorialCompleted = false;

        this.operationCount = 0;
        this.interceptionEfficiency = 0;
        this.stabilityScore = 0;
        this.totalScore = 0;
        this.dredgeDepth = 0;

        this.achievements = {
            firstMacha: false,
            tenBamboo: false,
            interception50: false,
            interception80: false,
            dredgeToWolong: false,
            completeRepair: false
        };
        this.achievementQueue = [];

        this.waterParticles = [];

        this.wolongIronY = 0;
        this.riverbedBaseY = 0;

        if (!document.getElementById('tutorial-style')) {
            const style = document.createElement('style');
            style.id = 'tutorial-style';
            style.textContent = `
                .tutorial-overlay { position: fixed; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0,0,0,0.5); z-index: 9999; display: flex; align-items: center; justify-content: center; }
                .tutorial-card { background: white; padding: 24px; border-radius: 12px; max-width: 420px; box-shadow: 0 10px 40px rgba(0,0,0,0.3); animation: slideUp 0.3s ease; }
                .tutorial-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
                .tutorial-step { background: #3498db; color: white; padding: 4px 12px; border-radius: 12px; font-size: 14px; font-weight: 600; }
                .tutorial-close { background: none; border: none; font-size: 24px; cursor: pointer; color: #999; }
                .tutorial-hint { background: #fff9e6; padding: 12px; border-radius: 8px; margin: 16px 0; border-left: 4px solid #f39c12; font-size: 14px; }
                .tutorial-progress { display: flex; gap: 8px; justify-content: center; margin-top: 16px; }
                .tutorial-dot { width: 12px; height: 12px; border-radius: 50%; background: #ddd; }
                .tutorial-dot.active { background: #3498db; transform: scale(1.2); }
                .tutorial-dot.completed { background: #2ecc71; }
                @keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
            `;
            document.head.appendChild(style);
        }

        this.initCanvas();
        this.initSession();
        this.initWaterParticles();
    }

    initCanvas() {
        const rect = this.canvas.parentElement.getBoundingClientRect();
        this.canvas.width = rect.width;
        this.canvas.height = 500;
        this.width = this.canvas.width;
        this.height = this.canvas.height;

        this.riverbedBaseY = this.height * 0.75;
        this.wolongIronY = this.height * 0.82;
        this.waterSurfaceY = this.height * 0.45;
    }

    initSession() {
        let savedSession = localStorage.getItem('dujiangyan_ux_session');
        if (savedSession) {
            try {
                const data = JSON.parse(savedSession);
                this.sessionId = data.sessionId;
            } catch (e) {
                this.sessionId = this.generateSessionId();
            }
        } else {
            this.sessionId = this.generateSessionId();
        }

        const savedNickname = localStorage.getItem('dujiangyan_ux_nickname');
        if (savedNickname) {
            this.nickname = savedNickname;
            const nicknameInput = document.getElementById('user-nickname');
            if (nicknameInput) nicknameInput.value = savedNickname;
        }

        this.saveSessionToStorage();
        this.updateSessionDisplay();
    }

    generateSessionId() {
        return 'UX-' + Date.now().toString(36) + '-' + Math.random().toString(36).substr(2, 9).toUpperCase();
    }

    saveSessionToStorage() {
        localStorage.setItem('dujiangyan_ux_session', JSON.stringify({
            sessionId: this.sessionId,
            timestamp: Date.now()
        }));
    }

    saveNickname(nickname) {
        this.nickname = nickname;
        localStorage.setItem('dujiangyan_ux_nickname', nickname);
    }

    updateSessionDisplay() {
        const sessionIdEl = document.getElementById('session-id');
        if (sessionIdEl) sessionIdEl.textContent = this.sessionId;
    }

    initWaterParticles() {
        this.waterParticles = [];
        for (let i = 0; i < 80; i++) {
            this.waterParticles.push({
                x: Math.random() * this.width,
                y: this.waterSurfaceY + Math.random() * (this.riverbedBaseY - this.waterSurfaceY - 20),
                vx: 1 + Math.random() * 2,
                vy: (Math.random() - 0.5) * 0.5,
                size: 2 + Math.random() * 3,
                opacity: 0.4 + Math.random() * 0.4
            });
        }
    }

    start() {
        if (this.running) return;
        this.running = true;
        this.startTime = Date.now() - this.elapsedTime;
        this.startTimer();
        this.animate();
    }

    stop() {
        this.running = false;
        this.stopTimer();
        if (this.animationId) {
            cancelAnimationFrame(this.animationId);
        }
    }

    startTimer() {
        if (this.timerInterval) clearInterval(this.timerInterval);
        this.timerInterval = setInterval(() => {
            this.elapsedTime = Date.now() - this.startTime;
            this.updateTimerDisplay();
        }, 1000);
    }

    stopTimer() {
        if (this.timerInterval) {
            clearInterval(this.timerInterval);
            this.timerInterval = null;
        }
    }

    updateTimerDisplay() {
        const seconds = Math.floor(this.elapsedTime / 1000);
        const mins = Math.floor(seconds / 60);
        const secs = seconds % 60;
        const timerEl = document.getElementById('timer');
        if (timerEl) {
            timerEl.textContent = `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
        }
    }

    animate() {
        if (!this.running) return;

        this.ctx.clearRect(0, 0, this.width, this.height);
        this.timeStep++;

        this.drawBackground();
        this.drawRiverbed();
        this.drawWater();
        this.updateAndDrawWaterParticles();
        this.drawWolongIron();
        this.drawDredgeAreas();
        this.drawPlacedObjects();

        if (this.selectedObject) {
            this.drawSelectionIndicator(this.selectedObject);
        }

        this.processAchievementQueue();

        this.animationId = requestAnimationFrame(() => this.animate());
    }

    drawBackground() {
        const gradient = this.ctx.createLinearGradient(0, 0, 0, this.height);
        gradient.addColorStop(0, '#1a2a4a');
        gradient.addColorStop(0.4, '#2d4a6a');
        gradient.addColorStop(1, '#3d5a7a');
        this.ctx.fillStyle = gradient;
        this.ctx.fillRect(0, 0, this.width, this.height);

        this.ctx.strokeStyle = 'rgba(0, 212, 255, 0.08)';
        this.ctx.lineWidth = 1;
        for (let i = 0; i < this.width; i += 40) {
            this.ctx.beginPath();
            this.ctx.moveTo(i, 0);
            this.ctx.lineTo(i, this.height);
            this.ctx.stroke();
        }
        for (let i = 0; i < this.height; i += 40) {
            this.ctx.beginPath();
            this.ctx.moveTo(0, i);
            this.ctx.lineTo(this.width, i);
            this.ctx.stroke();
        }
    }

    drawRiverbed() {
        this.ctx.fillStyle = CONFIG.COLORS.terrain;
        this.ctx.beginPath();
        this.ctx.moveTo(0, this.riverbedBaseY);

        for (let x = 0; x <= this.width; x += 15) {
            let yOffset = Math.sin(x * 0.015) * 8 + Math.sin(x * 0.04) * 4;
            let adjustedY = this.riverbedBaseY + yOffset + this.dredgeDepth;
            this.ctx.lineTo(x, adjustedY);
        }

        this.ctx.lineTo(this.width, this.height);
        this.ctx.lineTo(0, this.height);
        this.ctx.closePath();
        this.ctx.fill();

        this.ctx.fillStyle = CONFIG.COLORS.sand;
        for (let x = 0; x < this.width; x += 12) {
            for (let y = this.riverbedBaseY + 5 + this.dredgeDepth; y < this.height; y += 8) {
                if (Math.random() > 0.65) {
                    this.ctx.beginPath();
                    this.ctx.arc(x + Math.random() * 8, y + Math.random() * 4, 1 + Math.random() * 1.5, 0, Math.PI * 2);
                    this.ctx.fill();
                }
            }
        }
    }

    drawWater() {
        const waterGradient = this.ctx.createLinearGradient(0, this.waterSurfaceY, 0, this.riverbedBaseY);
        waterGradient.addColorStop(0, 'rgba(0, 168, 255, 0.7)');
        waterGradient.addColorStop(0.5, 'rgba(0, 136, 220, 0.6)');
        waterGradient.addColorStop(1, 'rgba(0, 102, 170, 0.5)');
        this.ctx.fillStyle = waterGradient;

        this.ctx.beginPath();
        this.ctx.moveTo(0, this.waterSurfaceY);
        for (let x = 0; x <= this.width; x += 10) {
            const waveY = this.waterSurfaceY + Math.sin(x * 0.03 + this.timeStep * 0.05) * 3;
            this.ctx.lineTo(x, waveY);
        }
        this.ctx.lineTo(this.width, this.riverbedBaseY + this.dredgeDepth);
        this.ctx.lineTo(0, this.riverbedBaseY + this.dredgeDepth);
        this.ctx.closePath();
        this.ctx.fill();

        this.ctx.strokeStyle = 'rgba(255, 255, 255, 0.3)';
        this.ctx.lineWidth = 1;
        this.ctx.beginPath();
        for (let x = 0; x <= this.width; x += 10) {
            const waveY = this.waterSurfaceY + Math.sin(x * 0.03 + this.timeStep * 0.05) * 3;
            if (x === 0) this.ctx.moveTo(x, waveY);
            else this.ctx.lineTo(x, waveY);
        }
        this.ctx.stroke();
    }

    updateAndDrawWaterParticles() {
        this.waterParticles.forEach(particle => {
            particle.x += particle.vx;
            particle.y += particle.vy;
            particle.y += Math.sin(this.timeStep * 0.08 + particle.x * 0.05) * 0.3;

            if (particle.x > this.width) {
                particle.x = -10;
                particle.y = this.waterSurfaceY + 10 + Math.random() * (this.riverbedBaseY - this.waterSurfaceY - 30);
            }
            if (particle.y < this.waterSurfaceY) particle.y = this.waterSurfaceY + 5;
            if (particle.y > this.riverbedBaseY + this.dredgeDepth - 5) {
                particle.y = this.riverbedBaseY + this.dredgeDepth - 10;
            }

            this.placedObjects.forEach(obj => {
                if (obj.type === 'macha' && this.checkParticleMachaCollision(particle, obj)) {
                    particle.vx *= 0.3;
                }
            });

            this.ctx.fillStyle = `rgba(0, 212, 255, ${particle.opacity})`;
            this.ctx.beginPath();
            this.ctx.arc(particle.x, particle.y, particle.size, 0, Math.PI * 2);
            this.ctx.fill();
        });
    }

    checkParticleMachaCollision(particle, macha) {
        return particle.x > macha.x - macha.width / 2 - 10 &&
               particle.x < macha.x + macha.width / 2 + 10 &&
               particle.y > macha.y - 5 &&
               particle.y < macha.y + macha.height;
    }

    drawWolongIron() {
        this.ctx.fillStyle = CONFIG.COLORS.wolongIron;
        this.ctx.strokeStyle = '#808080';
        this.ctx.lineWidth = 2;

        const ironWidth = 40;
        const ironHeight = 8;
        const ironX = this.width * 0.5;
        const ironY = this.wolongIronY;

        this.ctx.fillRect(ironX - ironWidth / 2, ironY, ironWidth, ironHeight);
        this.ctx.strokeRect(ironX - ironWidth / 2, ironY, ironWidth, ironHeight);

        this.ctx.fillStyle = '#ffffff';
        this.ctx.font = '11px Arial';
        this.ctx.textAlign = 'center';
        this.ctx.fillText('卧铁', ironX, ironY - 5);

        this.ctx.setLineDash([5, 5]);
        this.ctx.strokeStyle = 'rgba(192, 192, 192, 0.4)';
        this.ctx.beginPath();
        this.ctx.moveTo(0, ironY);
        this.ctx.lineTo(this.width, ironY);
        this.ctx.stroke();
        this.ctx.setLineDash([]);
    }

    drawDredgeAreas() {
        this.dredgeAreas.forEach(area => {
            this.ctx.fillStyle = 'rgba(139, 115, 85, 0.6)';
            this.ctx.beginPath();
            this.ctx.ellipse(area.x, area.y, area.width, area.height, 0, 0, Math.PI * 2);
            this.ctx.fill();

            this.ctx.strokeStyle = 'rgba(90, 70, 50, 0.8)';
            this.ctx.lineWidth = 2;
            this.ctx.stroke();
        });
    }

    drawPlacedObjects() {
        this.placedObjects.forEach(obj => {
            if (obj.type === 'macha') {
                this.drawMacha(obj);
            } else if (obj.type === 'bamboo') {
                this.drawBambooCage(obj);
            }
        });
    }

    drawMacha(obj) {
        this.ctx.save();
        this.ctx.translate(obj.x, obj.y);
        this.ctx.rotate(obj.angle || 0);

        this.ctx.strokeStyle = CONFIG.COLORS.wood;
        this.ctx.lineWidth = 3;
        this.ctx.fillStyle = 'rgba(139, 69, 19, 0.2)';

        this.ctx.beginPath();
        this.ctx.moveTo(0, 0);
        this.ctx.lineTo(-obj.width / 2, obj.height);
        this.ctx.lineTo(obj.width / 2, obj.height);
        this.ctx.closePath();
        this.ctx.fill();
        this.ctx.stroke();

        this.ctx.beginPath();
        this.ctx.moveTo(-obj.width / 3, obj.height * 0.5);
        this.ctx.lineTo(obj.width / 3, obj.height * 0.5);
        this.ctx.stroke();

        this.ctx.fillStyle = CONFIG.COLORS.bamboo;
        this.ctx.fillRect(-obj.width / 2 - 5, obj.height - 6, obj.width + 10, 10);

        this.ctx.fillStyle = '#8b4513';
        this.ctx.beginPath();
        this.ctx.arc(-obj.width / 3, obj.height * 0.5, 3, 0, Math.PI * 2);
        this.ctx.fill();
        this.ctx.beginPath();
        this.ctx.arc(obj.width / 3, obj.height * 0.5, 3, 0, Math.PI * 2);
        this.ctx.fill();

        this.ctx.restore();
    }

    drawBambooCage(obj) {
        this.ctx.save();
        this.ctx.translate(obj.x, obj.y);

        this.ctx.strokeStyle = CONFIG.COLORS.bamboo;
        this.ctx.lineWidth = 2;
        this.ctx.fillStyle = 'rgba(107, 142, 35, 0.35)';

        this.ctx.beginPath();
        this.roundRect(-obj.width / 2, -obj.height / 2, obj.width, obj.height, 4);
        this.ctx.fill();
        this.ctx.stroke();

        this.ctx.strokeStyle = 'rgba(107, 142, 35, 0.7)';
        this.ctx.beginPath();
        this.ctx.moveTo(-obj.width / 2, 0);
        this.ctx.lineTo(obj.width / 2, 0);
        this.ctx.stroke();

        this.ctx.beginPath();
        this.ctx.moveTo(0, -obj.height / 2);
        this.ctx.lineTo(0, obj.height / 2);
        this.ctx.stroke();

        this.ctx.fillStyle = CONFIG.COLORS.stone;
        for (let i = 0; i < 5; i++) {
            const sx = -obj.width / 2 + 8 + (i % 3) * (obj.width / 3);
            const sy = -obj.height / 2 + 8 + Math.floor(i / 3) * (obj.height / 2);
            this.ctx.beginPath();
            this.ctx.arc(sx, sy, 2 + Math.random(), 0, Math.PI * 2);
            this.ctx.fill();
        }

        this.ctx.restore();
    }

    roundRect(x, y, w, h, r) {
        this.ctx.moveTo(x + r, y);
        this.ctx.lineTo(x + w - r, y);
        this.ctx.quadraticCurveTo(x + w, y, x + w, y + r);
        this.ctx.lineTo(x + w, y + h - r);
        this.ctx.quadraticCurveTo(x + w, y + h, x + w - r, y + h);
        this.ctx.lineTo(x + r, y + h);
        this.ctx.quadraticCurveTo(x, y + h, x, y + h - r);
        this.ctx.lineTo(x, y + r);
        this.ctx.quadraticCurveTo(x, y, x + r, y);
    }

    drawSelectionIndicator(obj) {
        this.ctx.save();
        this.ctx.strokeStyle = '#ffff00';
        this.ctx.lineWidth = 2;
        this.ctx.setLineDash([5, 3]);

        if (obj.type === 'macha') {
            this.ctx.strokeRect(
                obj.x - obj.width / 2 - 8,
                obj.y - 8,
                obj.width + 16,
                obj.height + 16
            );
        } else if (obj.type === 'bamboo') {
            this.ctx.strokeRect(
                obj.x - obj.width / 2 - 5,
                obj.y - obj.height / 2 - 5,
                obj.width + 10,
                obj.height + 10
            );
        }

        this.ctx.setLineDash([]);
        this.ctx.restore();
    }

    selectTool(tool) {
        this.selectedTool = tool;
        this.selectedObject = null;

        document.querySelectorAll('.tool-btn').forEach(btn => {
            btn.classList.remove('active');
            if (btn.dataset.tool === tool) {
                btn.classList.add('active');
            }
        });

        const toolNames = { macha: '杩槎', bamboo: '竹笼', dredge: '淘淤' };
        const currentToolEl = document.getElementById('current-tool');
        if (currentToolEl) {
            currentToolEl.textContent = tool ? toolNames[tool] || '未选择' : '未选择';
        }

        this.canvas.style.cursor = tool ? 'crosshair' : 'default';
    }

    handleCanvasClick(e) {
        if (!this.running) return;

        const pos = this.getCanvasPosition(e);

        if (!this.selectedTool) return;
        if (this.selectedObject) return;

        if (this.selectedTool === 'macha') {
            this.placeMacha(pos.x, pos.y);
        } else if (this.selectedTool === 'bamboo') {
            this.placeBambooCage(pos.x, pos.y);
        } else if (this.selectedTool === 'dredge') {
            this.performDredge(pos.x, pos.y);
        }
    }

    handleMouseDown(e) {
        if (!this.running) return;

        const pos = this.getCanvasPosition(e);

        const clickedObject = this.findObjectAtPosition(pos.x, pos.y);
        if (clickedObject) {
            this.selectedObject = clickedObject;
            this.isDragging = true;
            this.dragOffset.x = pos.x - clickedObject.x;
            this.dragOffset.y = pos.y - clickedObject.y;
            this.canvas.style.cursor = 'move';
        }
    }

    handleMouseMove(e) {
        if (!this.running) return;

        const pos = this.getCanvasPosition(e);

        if (this.isDragging && this.selectedObject) {
            this.selectedObject.x = pos.x - this.dragOffset.x;
            this.selectedObject.y = pos.y - this.dragOffset.y;
            this.clampObjectPosition(this.selectedObject);
        }
    }

    handleMouseUp(e) {
        if (this.isDragging && this.selectedObject) {
            this.isDragging = false;
            this.calculateObjectScores(this.selectedObject);
            this.updateScores();
            this.recordOperation('drag', this.selectedObject.type, {
                x: this.selectedObject.x,
                y: this.selectedObject.y
            });
            this.canvas.style.cursor = this.selectedTool ? 'crosshair' : 'default';
        }
    }

    getCanvasPosition(e) {
        const rect = this.canvas.getBoundingClientRect();
        return {
            x: e.clientX - rect.left,
            y: e.clientY - rect.top
        };
    }

    findObjectAtPosition(x, y) {
        for (let i = this.placedObjects.length - 1; i >= 0; i--) {
            const obj = this.placedObjects[i];
            if (this.isPointInObject(x, y, obj)) {
                return obj;
            }
        }
        return null;
    }

    isPointInObject(x, y, obj) {
        if (obj.type === 'macha') {
            return x > obj.x - obj.width / 2 - 5 &&
                   x < obj.x + obj.width / 2 + 5 &&
                   y > obj.y - 5 &&
                   y < obj.y + obj.height + 5;
        } else if (obj.type === 'bamboo') {
            return x > obj.x - obj.width / 2 &&
                   x < obj.x + obj.width / 2 &&
                   y > obj.y - obj.height / 2 &&
                   y < obj.y + obj.height / 2;
        }
        return false;
    }

    clampObjectPosition(obj) {
        obj.x = Math.max(30, Math.min(this.width - 30, obj.x));
        if (obj.type === 'macha') {
            obj.y = Math.max(this.waterSurfaceY + 10, Math.min(this.riverbedBaseY + this.dredgeDepth - obj.height + 10, obj.y));
        } else if (obj.type === 'bamboo') {
            obj.y = Math.max(this.waterSurfaceY + 20, Math.min(this.riverbedBaseY + this.dredgeDepth - obj.height / 2, obj.y));
        }
    }

    placeMacha(x, y) {
        const macha = {
            id: 'macha_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5),
            type: 'macha',
            x: x,
            y: Math.min(y, this.riverbedBaseY + this.dredgeDepth - 50),
            width: 35,
            height: 55,
            angle: (Math.random() - 0.5) * 0.2,
            interception_efficiency: 0,
            stability_score: 0
        };

        this.clampObjectPosition(macha);
        this.placedObjects.push(macha);
        this.calculateObjectScores(macha);
        this.operationCount++;
        this.updateScores();
        this.checkAchievements();

        this.recordOperation('place', 'macha', { x: macha.x, y: macha.y });

        if (!this.achievements.firstMacha) {
            this.achievements.firstMacha = true;
            this.queueAchievement('首次放置杩槎！');
        }
    }

    placeBambooCage(x, y) {
        if (this.tutorialMode && !this.tutorialCompleted) {
            const success = this.advanceBambooStep('deploy_cage');
            if (!success) {
                this.showTutorialHint('请先完成编织前序步骤！');
                return;
            }
        }

        const bamboo = {
            id: 'bamboo_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5),
            type: 'bamboo',
            x: x,
            y: Math.min(y, this.riverbedBaseY + this.dredgeDepth - 15),
            width: 40,
            height: 28,
            interception_efficiency: 0,
            stability_score: 0
        };

        this.clampObjectPosition(bamboo);
        this.placedObjects.push(bamboo);
        this.calculateObjectScores(bamboo);
        this.operationCount++;
        this.updateScores();
        this.checkAchievements();

        this.recordOperation('place', 'bamboo', { x: bamboo.x, y: bamboo.y });

        const bambooCount = this.placedObjects.filter(o => o.type === 'bamboo').length;
        if (bambooCount >= 10 && !this.achievements.tenBamboo) {
            this.achievements.tenBamboo = true;
            this.queueAchievement('放置了10个竹笼！');
        }
    }

    performDredge(x, y) {
        const area = {
            id: 'dredge_' + Date.now() + '_' + Math.random().toString(36).substr(2, 5),
            x: x,
            y: Math.max(y, this.riverbedBaseY),
            width: 35,
            height: 15
        };

        this.dredgeAreas.push(area);
        this.dredgeDepth = Math.min(this.dredgeDepth + 3, 40);
        this.operationCount++;
        this.updateScores();
        this.checkAchievements();

        this.recordOperation('dredge', 'dredge', {
            x: area.x,
            y: area.y,
            dredge_depth: this.dredgeDepth
        });

        if (this.dredgeDepth >= 30 && !this.achievements.dredgeToWolong) {
            this.achievements.dredgeToWolong = true;
            this.queueAchievement('淘淤至卧铁深度！');
        }
    }

    calculateObjectScores(obj) {
        const riverCenterX = this.width / 2;
        const distFromCenter = Math.abs(obj.x - riverCenterX);
        const centerBonus = Math.max(0, 1 - distFromCenter / (this.width / 2));

        const riverbedY = this.riverbedBaseY + this.dredgeDepth;
        const heightRatio = Math.max(0, 1 - (obj.y - this.waterSurfaceY) / (riverbedY - this.waterSurfaceY));

        if (obj.type === 'macha') {
            const baseEfficiency = 4 + Math.random();
            obj.interception_efficiency = baseEfficiency * (0.6 + centerBonus * 0.4) * (0.5 + heightRatio * 0.5);
            obj.stability_score = 50 + centerBonus * 30 + (1 - heightRatio) * 20;
        } else if (obj.type === 'bamboo') {
            obj.interception_efficiency = 1 + centerBonus * 1.5;
            obj.stability_score = 60 + (1 - heightRatio) * 30 + centerBonus * 10;
        }

        obj.interception_efficiency = Math.max(0, Math.min(100, obj.interception_efficiency));
        obj.stability_score = Math.max(0, Math.min(100, obj.stability_score));
    }

    updateScores() {
        let totalEfficiency = 0;
        let totalStability = 0;
        let machaCount = 0;
        let bambooCount = 0;

        this.placedObjects.forEach(obj => {
            totalEfficiency += obj.interception_efficiency || 0;
            totalStability += obj.stability_score || 0;
            if (obj.type === 'macha') machaCount++;
            if (obj.type === 'bamboo') bambooCount++;
        });

        this.interceptionEfficiency = Math.min(85, totalEfficiency);

        const totalObjects = machaCount + bambooCount;
        if (totalObjects > 0) {
            const bambooStabilityBonus = bambooCount * 3;
            this.stabilityScore = Math.min(100, (totalStability / totalObjects) + bambooStabilityBonus);
        } else {
            this.stabilityScore = 0;
        }

        const dredgeScore = this.dredgeDepth * 2;
        const efficiencyScore = this.interceptionEfficiency * 5;
        const stabilityScoreComponent = this.stabilityScore * 2;
        const timeBonus = Math.max(0, 300 - Math.floor(this.elapsedTime / 1000));
        this.totalScore = Math.floor(efficiencyScore + stabilityScoreComponent + dredgeScore + timeBonus);

        this.updateScoreDisplay();
        this.checkAchievements();
    }

    updateScoreDisplay() {
        const opCountEl = document.getElementById('operation-count');
        if (opCountEl) opCountEl.textContent = this.operationCount;

        const efficiencyEl = document.getElementById('interception-efficiency');
        if (efficiencyEl) efficiencyEl.textContent = this.interceptionEfficiency.toFixed(1) + ' %';

        const stabilityEl = document.getElementById('stability-score');
        if (stabilityEl) stabilityEl.textContent = this.stabilityScore.toFixed(1);

        const scoreEl = document.getElementById('realtime-score');
        if (scoreEl) scoreEl.textContent = this.totalScore;

        const completionEl = document.getElementById('completion-rate');
        if (completionEl) {
            const completion = Math.min(100,
                (this.interceptionEfficiency / 85) * 50 +
                (this.stabilityScore / 100) * 25 +
                (this.dredgeDepth / 40) * 25
            );
            completionEl.textContent = completion.toFixed(0) + '%';
        }
    }

    checkAchievements() {
        if (this.interceptionEfficiency >= 50 && !this.achievements.interception50) {
            this.achievements.interception50 = true;
            this.queueAchievement('截流效率超过50%！');
        }

        if (this.interceptionEfficiency >= 80 && !this.achievements.interception80) {
            this.achievements.interception80 = true;
            this.queueAchievement('截流效率超过80%！');
        }

        if (this.interceptionEfficiency >= 70 &&
            this.stabilityScore >= 70 &&
            this.dredgeDepth >= 25 &&
            !this.achievements.completeRepair) {
            this.achievements.completeRepair = true;
            this.queueAchievement('完成岁修工程！');
        }
    }

    queueAchievement(message) {
        this.achievementQueue.push({
            message: message,
            startTime: this.timeStep,
            duration: 180
        });
        this.saveAchievements();
        this.updateAchievementDisplay();
    }

    updateAchievementDisplay() {
        const achievementMap = {
            firstMacha: 'ach-firstMacha',
            tenBamboo: 'ach-tenBamboo',
            interception50: 'ach-interception50',
            interception80: 'ach-interception80',
            dredgeToWolong: 'ach-dredgeToWolong',
            completeRepair: 'ach-completeRepair'
        };

        Object.keys(achievementMap).forEach(key => {
            const el = document.getElementById(achievementMap[key]);
            if (el) {
                if (this.achievements[key]) {
                    el.classList.add('unlocked');
                } else {
                    el.classList.remove('unlocked');
                }
            }
        });
    }

    processAchievementQueue() {
        if (this.achievementQueue.length === 0) return;

        const achievement = this.achievementQueue[0];
        const elapsed = this.timeStep - achievement.startTime;

        if (elapsed > achievement.duration) {
            this.achievementQueue.shift();
            return;
        }

        const alpha = elapsed < 30 ? elapsed / 30 :
                      elapsed > achievement.duration - 30 ? (achievement.duration - elapsed) / 30 : 1;

        this.ctx.save();
        this.ctx.globalAlpha = alpha;

        const text = '🏆 ' + achievement.message;
        this.ctx.font = 'bold 20px Arial';
        const textWidth = this.ctx.measureText(text).width;
        const boxWidth = textWidth + 60;
        const boxHeight = 50;
        const boxX = (this.width - boxWidth) / 2;
        const boxY = 20;

        this.ctx.fillStyle = 'rgba(255, 200, 0, 0.9)';
        this.ctx.strokeStyle = '#ffaa00';
        this.ctx.lineWidth = 3;
        this.ctx.beginPath();
        this.roundRect(boxX, boxY, boxWidth, boxHeight, 10);
        this.ctx.fill();
        this.ctx.stroke();

        this.ctx.fillStyle = '#5c3a00';
        this.ctx.textAlign = 'center';
        this.ctx.textBaseline = 'middle';
        this.ctx.fillText(text, this.width / 2, boxY + boxHeight / 2);

        this.ctx.restore();
    }

    saveAchievements() {
        localStorage.setItem('dujiangyan_ux_achievements', JSON.stringify(this.achievements));
    }

    loadAchievements() {
        const saved = localStorage.getItem('dujiangyan_ux_achievements');
        if (saved) {
            try {
                this.achievements = Object.assign(this.achievements, JSON.parse(saved));
            } catch (e) {}
        }
    }

    async recordOperation(action, objectType, data) {
        const operation = {
            session_id: this.sessionId,
            nickname: this.nickname,
            action: action,
            object_type: objectType,
            data: data,
            timestamp: new Date().toISOString(),
            operation_number: this.operationCount
        };

        try {
            await API.userExperience.saveOperation(operation);
        } catch (error) {
            console.log('API不可用，保存操作到本地');
            this.saveOperationLocal(operation);
        }
    }

    saveOperationLocal(operation) {
        let operations = [];
        const saved = localStorage.getItem('dujiangyan_ux_operations');
        if (saved) {
            try {
                operations = JSON.parse(saved);
            } catch (e) {}
        }
        operations.push(operation);
        localStorage.setItem('dujiangyan_ux_operations', JSON.stringify(operations));
    }

    async finishSession() {
        this.stop();

        const sessionData = {
            session_id: this.sessionId,
            nickname: this.nickname || '匿名用户',
            total_score: this.totalScore,
            interception_efficiency: this.interceptionEfficiency,
            stability_score: this.stabilityScore,
            operation_count: this.operationCount,
            dredge_depth: this.dredgeDepth,
            elapsed_time: this.elapsedTime,
            object_count: this.placedObjects.length,
            achievements: this.achievements,
            completed_at: new Date().toISOString()
        };

        try {
            await API.userExperience.finishSession(sessionData);
        } catch (error) {
            console.log('API不可用，保存会话到本地');
            this.saveSessionLocal(sessionData);
        }

        await this.showLeaderboard();
    }

    saveSessionLocal(sessionData) {
        let sessions = [];
        const saved = localStorage.getItem('dujiangyan_ux_sessions');
        if (saved) {
            try {
                sessions = JSON.parse(saved);
            } catch (e) {}
        }
        sessions.push(sessionData);
        sessions.sort((a, b) => b.total_score - a.total_score);
        localStorage.setItem('dujiangyan_ux_sessions', JSON.stringify(sessions.slice(0, 100)));
    }

    async showLeaderboard() {
        let rankingData = [];

        try {
            const response = await API.userExperience.getRanking(20);
            rankingData = response.ranking || response || [];
        } catch (error) {
            console.log('API不可用，加载本地排行榜');
            rankingData = this.getLocalRanking();
        }

        this.renderLeaderboard(rankingData);
    }

    getLocalRanking() {
        const saved = localStorage.getItem('dujiangyan_ux_sessions');
        if (!saved) return this.getMockRanking();

        try {
            const sessions = JSON.parse(saved);
            if (sessions.length === 0) return this.getMockRanking();
            return sessions.slice(0, 20);
        } catch (e) {
            return this.getMockRanking();
        }
    }

    getMockRanking() {
        const names = ['李冰后人', '治水能手', '岁修匠人', '竹笼大师', '杩槎专家', '河道守护', '淘淤先锋', '卧铁守望者'];
        return names.map((name, i) => ({
            nickname: name,
            total_score: 1000 - i * 80 - Math.floor(Math.random() * 50),
            elapsed_time: (120 + i * 30) * 1000,
            interception_efficiency: 80 - i * 3
        }));
    }

    renderLeaderboard(data) {
        const tbody = document.getElementById('leaderboard-body');
        if (!tbody) return;

        tbody.innerHTML = '';

        if (!data || data.length === 0) {
            tbody.innerHTML = '<tr><td colspan="4" style="text-align:center;color:#888;">暂无记录</td></tr>';
            return;
        }

        data.forEach((item, index) => {
            const tr = document.createElement('tr');
            const seconds = Math.floor((item.elapsed_time || 0) / 1000);
            const mins = Math.floor(seconds / 60);
            const secs = seconds % 60;
            const timeStr = `${mins}:${secs.toString().padStart(2, '0')}`;

            tr.innerHTML = `
                <td>${index + 1}</td>
                <td>${item.nickname || '匿名'}</td>
                <td class="score-value">${item.total_score || 0}</td>
                <td>${timeStr}</td>
            `;

            if (item.nickname === this.nickname && this.nickname) {
                tr.style.backgroundColor = 'rgba(0, 212, 255, 0.15)';
            }

            tbody.appendChild(tr);
        });
    }

    reset() {
        this.stop();
        this.placedObjects = [];
        this.dredgeAreas = [];
        this.selectedTool = null;
        this.selectedObject = null;
        this.isDragging = false;

        this.tutorialMode = true;
        this.currentBambooStep = 0;
        this.bambooWeavingSteps.forEach(s => s.completed = false);
        this.guideOverlayActive = false;
        this.showTutorial = true;
        this.tutorialCompleted = false;

        this.operationCount = 0;
        this.interceptionEfficiency = 0;
        this.stabilityScore = 0;
        this.totalScore = 0;
        this.dredgeDepth = 0;
        this.elapsedTime = 0;
        this.achievementQueue = [];

        this.sessionId = this.generateSessionId();
        this.saveSessionToStorage();
        this.updateSessionDisplay();

        document.querySelectorAll('.tool-btn').forEach(btn => btn.classList.remove('active'));
        const currentToolEl = document.getElementById('current-tool');
        if (currentToolEl) currentToolEl.textContent = '未选择';

        this.updateScoreDisplay();
        const timerEl = document.getElementById('timer');
        if (timerEl) timerEl.textContent = '00:00';

        this.ctx.clearRect(0, 0, this.width, this.height);
    }

    startBambooTutorial() {
        this.tutorialMode = true;
        this.currentBambooStep = 0;
        this.bambooWeavingSteps.forEach(s => s.completed = false);
        this.guideOverlayActive = true;
        this.showTutorial = true;
        this.tutorialCompleted = false;
        this.updateTutorialUI();
    }

    advanceBambooStep(actionType) {
        if (!this.tutorialMode || this.currentBambooStep >= this.bambooWeavingSteps.length) return;
        const currentStep = this.bambooWeavingSteps[this.currentBambooStep];
        if (actionType === currentStep.targetAction) {
            currentStep.completed = true;
            this.currentBambooStep++;
            this.addAchievement('bamboo_weaving_progress', `完成${this.currentBambooStep}/6步`);
            if (this.currentBambooStep >= this.bambooWeavingSteps.length) {
                this.tutorialCompleted = true;
                this.guideOverlayActive = false;
                this.addAchievement('bamboo_master', '完整完成竹笼编织');
            }
            this.updateTutorialUI();
            return true;
        }
        return false;
    }

    updateTutorialUI() {
        const guideEl = document.getElementById('bamboo-guide-overlay');
        if (!guideEl && this.guideOverlayActive) {
            const overlay = document.createElement('div');
            overlay.id = 'bamboo-guide-overlay';
            overlay.className = 'tutorial-overlay';
            document.body.appendChild(overlay);
        }
        if (guideEl) {
            if (this.guideOverlayActive && this.currentBambooStep < this.bambooWeavingSteps.length) {
                guideEl.style.display = 'flex';
                const step = this.bambooWeavingSteps[this.currentBambooStep];
                const isLastStep = step.targetAction === 'deploy_cage';
                guideEl.innerHTML = `
                    <div class="tutorial-card">
                        <div class="tutorial-header">
                            <span class="tutorial-step">${step.id}/6</span>
                            <button class="tutorial-close">×</button>
                        </div>
                        <h4>${step.title}</h4>
                        <p>${step.description}</p>
                        <div class="tutorial-hint">💡 ${step.canvasHint}</div>
                        <div class="tutorial-progress">
                            ${this.bambooWeavingSteps.map((s, i) => 
                                `<div class="tutorial-dot ${s.completed ? 'completed' : i === this.currentBambooStep ? 'active' : ''}"></div>`
                            ).join('')}
                        </div>
                        ${!isLastStep ? `<button class="tutorial-next-btn" style="margin-top: 16px; width: 100%; padding: 12px; background: #3498db; color: white; border: none; border-radius: 8px; font-size: 16px; cursor: pointer;">下一步</button>` : `<div style="margin-top: 16px; padding: 12px; background: #e8f5e9; border-radius: 8px; text-align: center; color: #2e7d32;">请在画布上点击放置竹笼完成最后一步</div>`}
                    </div>
                `;
                const closeBtn = guideEl.querySelector('.tutorial-close');
                if (closeBtn) {
                    closeBtn.onclick = () => {
                        this.guideOverlayActive = false;
                        guideEl.style.display = 'none';
                    };
                }
                const nextBtn = guideEl.querySelector('.tutorial-next-btn');
                if (nextBtn) {
                    nextBtn.onclick = () => {
                        this.advanceBambooStep(step.targetAction);
                    };
                }
            } else if (this.tutorialCompleted) {
                guideEl.style.display = 'none';
            }
        }
    }

    addAchievement(key, message) {
        if (!this.achievements[key]) {
            this.achievements[key] = true;
            this.queueAchievement(message);
        }
    }

    showTutorialHint(message) {
        this.queueAchievement(message);
    }
}

let userExperienceSim = null;

function initUserExperience() {
    if (!userExperienceSim) {
        userExperienceSim = new UserExperienceSimulation('experience-canvas');
    }
    return userExperienceSim;
}

document.addEventListener('DOMContentLoaded', () => {
    const canvas = document.getElementById('experience-canvas');
    if (!canvas) return;

    const sim = initUserExperience();

    const nicknameInput = document.getElementById('user-nickname');
    if (nicknameInput) {
        nicknameInput.addEventListener('input', (e) => {
            sim.saveNickname(e.target.value.trim());
        });
    }

    document.querySelectorAll('.tool-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const tool = btn.dataset.tool;
            if (sim.selectedTool === tool) {
                sim.selectTool(null);
            } else {
                sim.selectTool(tool);
            }
        });
    });

    canvas.addEventListener('mousedown', (e) => sim.handleMouseDown(e));
    canvas.addEventListener('mousemove', (e) => sim.handleMouseMove(e));
    canvas.addEventListener('mouseup', (e) => sim.handleMouseUp(e));
    canvas.addEventListener('mouseleave', (e) => sim.handleMouseUp(e));
    canvas.addEventListener('click', (e) => sim.handleCanvasClick(e));

    const finishBtn = document.getElementById('finish-session');
    if (finishBtn) {
        finishBtn.addEventListener('click', () => {
            sim.finishSession();
        });
    }

    const resetBtn = document.getElementById('reset-experience');
    if (resetBtn) {
        resetBtn.addEventListener('click', () => {
            sim.reset();
            sim.start();
        });

        const tutorialBtn = document.createElement('button');
        tutorialBtn.id = 'start-tutorial';
        tutorialBtn.className = 'btn-secondary';
        tutorialBtn.textContent = '编织教程';
        tutorialBtn.onclick = () => sim.startBambooTutorial();
        resetBtn.parentNode.insertBefore(tutorialBtn, resetBtn.nextSibling);
    }

    const startBtn = document.getElementById('start-experience');
    if (startBtn) {
        startBtn.addEventListener('click', () => {
            sim.start();
        });
    }

    if (!sim.running) {
        sim.start();
    }

    sim.startBambooTutorial();

    sim.loadAchievements();
    sim.updateAchievementDisplay();
    sim.showLeaderboard();
});

global.VRMaintenance = UserExperienceSimulation;
global.initUserExperience = initUserExperience;
})(window);
