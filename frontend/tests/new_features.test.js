const assert = require('assert');

const EPS = 1e-9;

function approxEqual(a, b, eps = EPS) {
    return Math.abs(a - b) < eps;
}

// ==============================================================
// Test 1: 工艺对比验证截流效率 - 对应 dynasty_compare.js
// ==============================================================
function testDynastyMachaInterception() {
    const MACHA_EFFICIENCY_CAP = 0.85;
    const LIBING_EFF_PER_UNIT = 0.04;
    const QING_EFF_PER_UNIT = 0.05;

    function calcMachaEfficiency(count, effPerUnit) {
        if (count <= 0 || effPerUnit <= 0) return 0;
        return Math.min(count * effPerUnit, MACHA_EFFICIENCY_CAP);
    }

    console.log('\n=== 测试一：工艺对比 - 杩槎截流效率 ===');

    // 1.1 正常场景
    console.log('  [正常] 李冰时期截流效率...');
    assert.ok(approxEqual(calcMachaEfficiency(0, LIBING_EFF_PER_UNIT), 0), '0个杩槎=0');
    assert.ok(approxEqual(calcMachaEfficiency(5, LIBING_EFF_PER_UNIT), 0.20), '5个杩槎=20%');
    assert.ok(approxEqual(calcMachaEfficiency(10, LIBING_EFF_PER_UNIT), 0.40), '10个杩槎=40%');
    assert.ok(approxEqual(calcMachaEfficiency(17, LIBING_EFF_PER_UNIT), 0.68), '17个杩槎=68%');
    assert.ok(approxEqual(calcMachaEfficiency(22, LIBING_EFF_PER_UNIT), 0.85), '22个杩槎触达上限');
    console.log('    ✓ 李冰时期通过');

    console.log('  [正常] 清代时期截流效率...');
    assert.ok(approxEqual(calcMachaEfficiency(0, QING_EFF_PER_UNIT), 0), '0个=0');
    assert.ok(approxEqual(calcMachaEfficiency(5, QING_EFF_PER_UNIT), 0.25), '5个=25%');
    assert.ok(approxEqual(calcMachaEfficiency(17, QING_EFF_PER_UNIT), 0.85), '17个触达上限');
    assert.ok(approxEqual(calcMachaEfficiency(30, QING_EFF_PER_UNIT), 0.85), '30个仍为上限');
    console.log('    ✓ 清代时期通过');

    console.log('  [正常] 清代效率 ≥ 李冰效率...');
    for (let count = 0; count <= 30; count++) {
        const libing = calcMachaEfficiency(count, LIBING_EFF_PER_UNIT);
        const qing = calcMachaEfficiency(count, QING_EFF_PER_UNIT);
        assert.ok(qing >= libing, `${count}个杩槎：清代${(qing*100).toFixed(1)}% ≥ 李冰${(libing*100).toFixed(1)}%`);
    }
    console.log('    ✓ 对比通过');

    // 1.2 渐进性验证
    console.log('  [正常] 效率递增不回退...');
    let prevEff = 0;
    for (let count = 0; count <= 30; count++) {
        const currEff = calcMachaEfficiency(count, LIBING_EFF_PER_UNIT);
        assert.ok(currEff >= prevEff, `第${count}个杩槎效率回退：${prevEff}→${currEff}`);
        prevEff = currEff;
    }
    console.log('    ✓ 渐进性通过');

    // 1.3 边界
    console.log('  [边界] 0/1/100个杩槎...');
    assert.ok(approxEqual(calcMachaEfficiency(1, LIBING_EFF_PER_UNIT), 0.04));
    assert.ok(approxEqual(calcMachaEfficiency(1, QING_EFF_PER_UNIT), 0.05));
    assert.ok(approxEqual(calcMachaEfficiency(100, LIBING_EFF_PER_UNIT), 0.85));
    assert.ok(approxEqual(calcMachaEfficiency(100, QING_EFF_PER_UNIT), 0.85));
    console.log('    ✓ 边界通过');

    // 1.4 异常
    console.log('  [异常] 负数/0参数...');
    assert.ok(approxEqual(calcMachaEfficiency(-5, LIBING_EFF_PER_UNIT), 0), '负数杩槎→0');
    assert.ok(approxEqual(calcMachaEfficiency(10, 0), 0), 'efficiency=0→0');
    assert.ok(approxEqual(calcMachaEfficiency(-10, -0.01), 0), '负数→0');
    console.log('    ✓ 异常通过');
}

// ==============================================================
// Test 2: 跨时代对比维护成本 - 对应 modern_compare.js
// ==============================================================
function testModernMaintenanceCost() {
    const ancientCost = [100, 80, 120, 40, 70, 100, 100, 50, 100, 100];
    const modernCost = [150, 130, 200, 110, 180, 250, 120, 200, 90, 80];
    const ancientLabor = [120, 80, 300, 50, 150, 100, 200, 100, 10, 30];
    const modernLabor = [15, 10, 8, 20, 25, 120, 30, 20, 5, 8];
    const ancientDuration = [60, 45, 45, 30, 40, 30, 50, 45, 20, 30];
    const modernDuration = [15, 20, 10, 15, 20, 25, 30, 15, 10, 20];

    function avg(arr) { return arr.length ? arr.reduce((a,b)=>a+b,0)/arr.length : 0; }
    function reductionRate(oldVal, newVal) { return oldVal>0 ? (oldVal-newVal)/oldVal*100 : 0; }
    function costEfficiencyRatio(cost, efficiency) { return efficiency>0 ? cost/efficiency : Infinity; }

    console.log('\n=== 测试二：跨时代对比 - 维护成本 ===');

    // 2.1 成本对比
    console.log('  [正常] 整体成本分析...');
    const avgAncient = avg(ancientCost);
    const avgModern = avg(modernCost);
    console.log(`    古代平均成本：${avgAncient.toFixed(2)}，现代平均成本：${avgModern.toFixed(2)}`);
    assert.ok(avgModern > avgAncient, `现代平均${avgModern}应>古代${avgAncient}`);
    assert.ok(Math.abs(avgAncient - 86) < 1, `古代平均≈86，实际${avgAncient}`);
    assert.ok(Math.abs(avgModern - 151) < 1, `现代平均≈151，实际${avgModern}`);
    console.log('    ✓ 成本对比通过');

    // 2.2 人力减少率
    console.log('  [正常] 人力减少率验证...');
    const interceptionReduction = reductionRate(120, 15);
    const bankReduction = reductionRate(80, 10);
    const dredgeReduction = reductionRate(300, 8);
    console.log(`    截流减少${interceptionReduction.toFixed(1)}%，护岸减少${bankReduction.toFixed(1)}%，清淤减少${dredgeReduction.toFixed(1)}%`);
    assert.ok(interceptionReduction > 85, `截流减少率${interceptionReduction}%>85%`);
    assert.ok(Math.abs(interceptionReduction - 87.5) < 0.1);
    assert.ok(Math.abs(bankReduction - 87.5) < 0.1);
    assert.ok(Math.abs(dredgeReduction - 97.33) < 0.1);
    const materialLaborIdx = 5;
    for (let i = 0; i < ancientLabor.length; i++) {
        if (i === materialLaborIdx) continue;
        assert.ok(modernLabor[i] < ancientLabor[i], `项${i}：现代人力${modernLabor[i]}<古代${ancientLabor[i]}`);
    }
    assert.ok(modernLabor[materialLaborIdx] > ancientLabor[materialLaborIdx],
        '材料项特殊：现代材料技术人力成本更高');
    console.log('    ✓ 人力减少通过');

    // 2.3 工期缩短
    console.log('  [正常] 工期缩短验证...');
    let totalOld = 0, totalNew = 0;
    for (let i = 0; i < ancientDuration.length; i++) {
        totalOld += ancientDuration[i];
        totalNew += modernDuration[i];
        assert.ok(modernDuration[i] < ancientDuration[i], `项${i}：现代工期${modernDuration[i]}<古代${ancientDuration[i]}`);
    }
    const totalReduction = reductionRate(totalOld, totalNew);
    const dredgeShortening = reductionRate(45, 10);
    console.log(`    总体缩短${totalReduction.toFixed(1)}%，清淤缩短${dredgeShortening.toFixed(1)}%`);
    assert.ok(totalReduction > 50, `总体缩短率${totalReduction}%>50%`);
    assert.ok(Math.abs(dredgeShortening - 77.78) < 0.1);
    console.log('    ✓ 工期缩短通过');

    // 2.4 成本效率比
    console.log('  [正常] 成本效率比...');
    const ancientInterceptRatio = costEfficiencyRatio(100, 60);
    const modernInterceptRatio = costEfficiencyRatio(150, 95);
    console.log(`    截流：古代${ancientInterceptRatio.toFixed(2)} vs 现代${modernInterceptRatio.toFixed(2)}`);
    assert.ok(modernInterceptRatio <= ancientInterceptRatio, `现代成本效率比${modernInterceptRatio}≤古代${ancientInterceptRatio}`);
    console.log('    ✓ 成本效率比通过');

    // 2.5 边界
    console.log('  [边界] 极端参数...');
    assert.strictEqual(costEfficiencyRatio(0, 100), 0, 'cost=0→0');
    assert.strictEqual(costEfficiencyRatio(100, 0), Infinity, 'eff=0→∞');
    assert.strictEqual(avg([]), 0, '空数组→0');
    console.log('    ✓ 边界通过');

    // 2.6 异常
    console.log('  [异常] 异常数据检测...');
    const hasNegative = (arr) => arr.some(v => v < 0);
    assert.strictEqual(hasNegative(ancientCost), false);
    assert.strictEqual(hasNegative(modernCost), false);
    const countOver100 = (arr) => arr.filter(v => v > 100 && v < 200).length;
    assert.ok(countOver100(modernCost) > 0);
    console.log('    ✓ 异常检测通过');
}

// ==============================================================
// Test 3: 地震影响验证结构安全 - 对应 earthquake_sim.js
// ==============================================================
function testEarthquakeStructureSafety() {
    const EPS3 = 1e-3;
    const STRUCTURES = {
        yuzui: { X: 100, Y: 50, Z: 730 },
        feishayan: { X: 150, Y: 30, Z: 728 },
        baopingkou: { X: 200, Y: 60, Z: 725 },
        renzidi: { X: 80, Y: 80, Z: 727 }
    };

    function calcPGA(mag) { return Math.pow(10, 0.5 * mag - 3); }
    function calcDistance(loc, epic) {
        return Math.sqrt(
            Math.pow(loc.X - epic.X, 2) +
            Math.pow(loc.Y - epic.Y, 2) +
            Math.pow(loc.Z - epic.Z, 2)
        );
    }
    function calcAttenuation(dist) { return Math.exp(-0.001 * dist); }
    function calcDamage(pga, dist) { return pga * calcAttenuation(dist); }
    function calcAllDamages(pga, epic) {
        const result = {};
        for (const [name, loc] of Object.entries(STRUCTURES)) {
            result[name] = calcDamage(pga, calcDistance(loc, epic));
        }
        return result;
    }
    function safetyAssessment(damages) {
        const max = Math.max(...Object.values(damages));
        if (max > 0.6) return 'danger';
        if (max > 0.3) return 'caution';
        return 'safe';
    }

    console.log('\n=== 测试三：地震模拟 - 结构安全 ===');

    // 3.1 PGA计算
    console.log('  [正常] PGA计算...');
    assert.ok(approxEqual(calcPGA(8.0), 10.0, EPS3), 'M8.0 PGA=10');
    assert.ok(approxEqual(calcPGA(7.0), Math.sqrt(10), EPS3), 'M7.0 PGA=√10');
    assert.ok(approxEqual(calcPGA(6.0), 1.0, EPS3), 'M6.0 PGA=1');
    assert.ok(approxEqual(calcPGA(5.0), 1/Math.sqrt(10), EPS3), 'M5.0 PGA=1/√10');
    assert.ok(approxEqual(calcPGA(4.0), 0.1, EPS3), 'M4.0 PGA=0.1');
    console.log('    ✓ PGA通过');

    // 3.2 安全评估
    console.log('  [正常] 安全分级...');
    const qushou = { X: 100, Y: 50, Z: 0 };
    const m8Damages = calcAllDamages(calcPGA(8.0), qushou);
    const m6Damages = calcAllDamages(calcPGA(6.0), qushou);
    const m4Damages = calcAllDamages(calcPGA(4.0), qushou);
    assert.strictEqual(safetyAssessment(m8Damages), 'danger', 'M8.0→danger');
    assert.strictEqual(safetyAssessment(m6Damages), 'caution', 'M6.0→caution');
    assert.strictEqual(safetyAssessment(m4Damages), 'safe', 'M4.0→safe');
    console.log(`    M8.0 最大损伤=${Math.max(...Object.values(m8Damages)).toFixed(4)} → danger`);
    console.log(`    M6.0 最大损伤=${Math.max(...Object.values(m6Damages)).toFixed(4)} → caution`);
    console.log(`    M4.0 最大损伤=${Math.max(...Object.values(m4Damages)).toFixed(4)} → safe`);
    console.log('    ✓ 安全分级通过');

    // 3.3 距离衰减
    console.log('  [正常] 距离衰减...');
    const pga7 = calcPGA(7.0);
    const nearEpic = qushou;
    const farEpic = { X: 1000, Y: 1000, Z: 0 };
    const yuzuiNear = calcDamage(pga7, calcDistance(STRUCTURES.yuzui, nearEpic));
    const yuzuiFar = calcDamage(pga7, calcDistance(STRUCTURES.yuzui, farEpic));
    assert.ok(yuzuiNear > yuzuiFar, `近处损伤${yuzuiNear}>远处${yuzuiFar}`);
    console.log(`    鱼嘴 近处=${yuzuiNear.toFixed(4)} > 远处=${yuzuiFar.toFixed(8)}`);
    console.log('    ✓ 距离衰减通过');

    // 3.4 鱼嘴震中→鱼嘴损伤最大
    console.log('  [正常] 结构损伤排序...');
    const yuzuiEpic = { X: 100, Y: 50, Z: 730 };
    const damagesAtYuzui = calcAllDamages(calcPGA(8.0), yuzuiEpic);
    console.log('    震中在鱼嘴时：', Object.entries(damagesAtYuzui).map(([k,v])=>`${k}=${v.toFixed(4)}`).join(', '));
    const { yuzui: yDmg, feishayan: fDmg, baopingkou: bDmg, renzidi: rDmg } = damagesAtYuzui;
    assert.ok(yDmg > fDmg, `鱼嘴${yDmg}>飞沙堰${fDmg}`);
    assert.ok(yDmg > bDmg, `鱼嘴${yDmg}>宝瓶口${bDmg}`);
    assert.ok(yDmg > rDmg, `鱼嘴${yDmg}>人字堤${rDmg}`);
    console.log('    ✓ 损伤排序通过');

    // 3.5 次生效应
    console.log('  [正常] 次生效应比例...');
    const pga8 = calcPGA(8.0);
    assert.ok(approxEqual(pga8 * 0.5, 5.0));
    assert.ok(approxEqual(pga8 * 0.3, 3.0));
    assert.ok(approxEqual(pga8 * 0.2, 2.0));
    const sediment = pga8 * 0.5, flow = pga8 * 0.3, diversion = pga8 * 0.2;
    assert.ok(sediment > flow && flow > diversion, `sediment(${sediment}) > flow(${flow}) > diversion(${diversion})`);
    console.log('    ✓ 次生效应通过');

    // 3.6 边界
    console.log('  [边界] M0/M10级地震...');
    const m0Dmg = calcAllDamages(calcPGA(0), qushou);
    const m10Dmg = calcAllDamages(calcPGA(10), qushou);
    assert.strictEqual(safetyAssessment(m0Dmg), 'safe', 'M0→safe');
    assert.strictEqual(safetyAssessment(m10Dmg), 'danger', 'M10→danger');
    assert.ok(approxEqual(calcPGA(0), 0.001, EPS3), 'M0 PGA=0.001');
    assert.ok(approxEqual(calcPGA(10), 100.0, EPS3), 'M10 PGA=100');
    console.log('    ✓ 边界通过');

    // 3.7 异常
    console.log('  [异常] 负震级/远距震中...');
    const mNeg1Dmg = calcAllDamages(calcPGA(-1), qushou);
    assert.strictEqual(safetyAssessment(mNeg1Dmg), 'safe', 'M-1→safe');
    const farFarEpic = { X: 10000, Y: 10000, Z: 0 };
    const farDmg = calcAllDamages(calcPGA(8.0), farFarEpic);
    assert.strictEqual(safetyAssessment(farDmg), 'safe', '超远距离→safe');
    console.log(`    超远距离最大损伤=${Math.max(...Object.values(farDmg)).toExponential(4)}`);
    console.log('    ✓ 异常通过');

    // 3.8 堤岸坍塌阈值
    console.log('  [正常] 堤岸坍塌阈值...');
    const m4Max = Math.max(...Object.values(m4Damages));
    const m8Max = Math.max(...Object.values(m8Damages));
    assert.ok(m4Max <= 0.3, `M4.0 maxDamage(${m4Max.toFixed(4)}) ≤ 0.3 → 无坍塌`);
    assert.ok(m8Max > 0.6, `M8.0 maxDamage(${m8Max.toFixed(4)}) > 0.6 → 有坍塌`);
    const thresholdCase = 0.31;
    const collapseCount031 = thresholdCase > 0.3 ? Math.floor(thresholdCase * 20) : 0;
    assert.ok(collapseCount031 > 0, '0.31时应产生坍塌');
    const noCollapse = 0.29;
    const collapseCount029 = noCollapse > 0.3 ? Math.floor(noCollapse * 20) : 0;
    assert.strictEqual(collapseCount029, 0, '0.29时无坍塌');
    console.log('    ✓ 坍塌阈值通过');
}

// ==============================================================
// Test 4: 虚拟体验测试操作合理性 - 对应 user_experience.js
// ==============================================================
function testUserExperienceOperations() {
    const MACHA_EFF_CAP = 0.85;
    const WOLONG_ELEVATION = 726.12;

    function calcMachaEff(count, perUnit) {
        if (count <= 0 || perUnit <= 0) return 0;
        return Math.min(count * perUnit, MACHA_EFF_CAP);
    }
    function calcBambooStability(initial, stonesFilled) {
        const clampedInit = Math.max(0.5, Math.min(1.0, initial));
        const result = clampedInit - 0.01 * stonesFilled;
        return Math.max(0.5, result);
    }
    function calcTotalScore(efficiencyPct, stabilityScore, dredgeDepth, timeBonus) {
        return efficiencyPct * 5 + stabilityScore * 2 + dredgeDepth * 2 + timeBonus;
    }
    function checkAchievements(ops, finalEff, finalStab, finalDredge) {
        const ach = new Set();
        const machaOps = ops.filter(o => o.type === 'place_macha');
        const bambooCount = ops.filter(o => o.type === 'place_bamboo').length;
        if (machaOps.length >= 1) ach.add('first_macha');
        if (bambooCount >= 10) ach.add('ten_bamboo');
        if (finalEff >= 0.50) ach.add('interception_50');
        if (finalEff >= 0.80) ach.add('interception_80');
        if (finalDredge >= (728.0 - WOLONG_ELEVATION)) ach.add('dredge_to_wolong');
        if (finalEff >= 0.70 && finalStab >= 70 && finalDredge >= 25) ach.add('complete_repair');
        return ach;
    }
    function validateOperation(op) {
        const validTypes = ['place_macha', 'place_bamboo', 'dredge', 'remove_macha', 'remove_bamboo'];
        if (!op || typeof op !== 'object') return false;
        if (!op.session_id || typeof op.session_id !== 'string') return false;
        if (!validTypes.includes(op.operation_type)) return false;
        if (typeof op.operation_order !== 'number' || op.operation_order < 0) return false;
        if (op.position_x < 0 || op.position_x > 1000) return false;
        if (op.position_y < 0 || op.position_y > 1000) return false;
        return true;
    }

    console.log('\n=== 测试四：虚拟体验 - 操作合理性 ===');

    // 4.1 杩槎正常操作
    console.log('  [正常] 杩槎截流效率...');
    const libing = 0.04, qing = 0.05;
    assert.ok(approxEqual(calcMachaEff(0, libing), 0));
    assert.ok(approxEqual(calcMachaEff(13, libing), 0.52), '13李冰杩槎=52%(触达50%)');
    assert.ok(approxEqual(calcMachaEff(20, libing), 0.80), '20李冰=80%');
    assert.ok(approxEqual(calcMachaEff(22, libing), 0.85), '22李冰=上限');
    assert.ok(approxEqual(calcMachaEff(10, qing), 0.50), '10清代=50%');
    assert.ok(approxEqual(calcMachaEff(16, qing), 0.80), '16清代=80%');
    assert.ok(approxEqual(calcMachaEff(17, qing), 0.85), '17清代=上限');
    console.log('    ✓ 杩槎效率通过');

    // 4.2 竹笼稳定性
    console.log('  [正常] 竹笼稳定性...');
    assert.ok(calcBambooStability(0.85, 0) >= 0.7 && calcBambooStability(0.85, 0) <= 1.0, '初始范围');
    assert.ok(approxEqual(calcBambooStability(1.0, 10), 0.90), '10石块→1.0-0.1=0.90');
    assert.ok(approxEqual(calcBambooStability(1.0, 50), 0.50), '50石块→下限0.50');
    assert.ok(approxEqual(calcBambooStability(0.7, 20), 0.50), '20石块→下限0.50');
    assert.ok(approxEqual(calcBambooStability(1.0, 100), 0.50), '100石块→仍为0.50');
    console.log('    ✓ 竹笼稳定性通过');

    // 4.3 总评分
    console.log('  [正常] 总评分计算...');
    const scoreNormal = calcTotalScore(70, 80, 30, 10);
    assert.strictEqual(scoreNormal, 580, `70*5+80*2+30*2+10=580，实际${scoreNormal}`);
    const scoreZero = calcTotalScore(0, 0, 0, 0);
    assert.strictEqual(scoreZero, 0, '全0=0');
    const scoreMax = calcTotalScore(85, 100, 50, 100);
    assert.strictEqual(scoreMax, 825, `满分=425+200+100+100=825，实际${scoreMax}`);
    console.log(`    正常=580，零分=0，满分=825 ✓`);
    console.log('    ✓ 总评分通过');

    // 4.4 成就系统
    console.log('  [正常] 成就解锁...');

    const ach1 = checkAchievements([], 0, 0, 0);
    assert.strictEqual(ach1.size, 0, '无操作→0成就');

    const ach2 = checkAchievements([{ type: 'place_macha' }], 0.04, 0, 0);
    assert.ok(ach2.has('first_macha'), '第1个杩槎→解锁first_macha');
    assert.strictEqual(ach2.size, 1);

    const ach3 = checkAchievements(
        Array.from({length: 13}, () => ({ type: 'place_macha' })),
        0.52, 0, 0
    );
    assert.ok(ach3.has('first_macha'), '13杩槎含first_macha');
    assert.ok(ach3.has('interception_50'), '13杩槎=52%→interception_50');
    console.log('    13杩槎(李冰)：截流超50% ✓');

    const tenBambooOps = [
        { type: 'place_macha' },
        ...Array.from({length: 10}, () => ({ type: 'place_bamboo' }))
    ];
    const ach4 = checkAchievements(tenBambooOps, 0.04, 0, 0);
    assert.ok(ach4.has('ten_bamboo'), '10竹笼→ten_bamboo');
    console.log('    10个竹笼 ✓');

    const ach5 = checkAchievements(
        Array.from({length: 20}, () => ({ type: 'place_macha' })),
        0.80, 0, 0
    );
    assert.ok(ach5.has('interception_80'), '20杩槎=80%→interception_80');
    console.log('    截流超80% ✓');

    const completeOps = [
        ...Array.from({length: 18}, () => ({ type: 'place_macha' })),
        ...Array.from({length: 5}, () => ({ type: 'place_bamboo' }))
    ];
    const ach6 = checkAchievements(completeOps, 0.72, 75, 30);
    assert.ok(ach6.has('complete_repair'), '72%效率+75稳+30淘淤→完成岁修');
    console.log('    完成岁修成就 ✓');

    const ach7 = checkAchievements(completeOps, 0.69, 80, 30);
    assert.ok(!ach7.has('complete_repair'), '69%效率不满足→不解锁');
    const ach8 = checkAchievements(completeOps, 0.72, 69, 30);
    assert.ok(!ach8.has('complete_repair'), '69稳定不满足→不解锁');
    const ach9 = checkAchievements(completeOps, 0.72, 80, 24);
    assert.ok(!ach9.has('complete_repair'), '24淘淤不满足→不解锁');
    console.log('    成就条件严格校验 ✓');
    console.log('    ✓ 成就系统通过');

    // 4.5 操作有效性验证
    console.log('  [正常] 操作合理性验证...');
    const validOp = {
        session_id: 'test-sess-001',
        operation_type: 'place_macha',
        operation_order: 1,
        position_x: 300, position_y: 200, position_z: 0
    };
    assert.strictEqual(validateOperation(validOp), true, '合法操作应通过验证');
    console.log('    ✓ 合法操作通过');

    // 4.6 边界操作
    console.log('  [边界] 操作边界...');
    const zeroOrderOp = { ...validOp, operation_order: 0 };
    assert.strictEqual(validateOperation(zeroOrderOp), true, '第0步操作有效');
    const maxPosOp = { ...validOp, position_x: 1000, position_y: 1000 };
    assert.strictEqual(validateOperation(maxPosOp), true, '边界坐标有效');
    const allTypes = ['place_macha', 'place_bamboo', 'dredge', 'remove_macha', 'remove_bamboo'];
    for (const t of allTypes) {
        const op = { ...validOp, operation_type: t };
        assert.strictEqual(validateOperation(op), true, `类型 ${t} 应有效`);
    }
    console.log('    ✓ 边界操作通过');

    // 4.7 异常操作
    console.log('  [异常] 异常操作验证...');
    const badType = { ...validOp, operation_type: 'invalid_type' };
    assert.strictEqual(validateOperation(badType), false, '非法类型应无效');
    const negOrder = { ...validOp, operation_order: -1 };
    assert.strictEqual(validateOperation(negOrder), false, '负序号应无效');
    const emptySession = { ...validOp, session_id: '' };
    assert.strictEqual(validateOperation(emptySession), false, '空session应无效');
    const nullOp = null;
    assert.strictEqual(validateOperation(nullOp), false, 'null应无效');
    const outOfBounds = { ...validOp, position_x: 1001 };
    assert.strictEqual(validateOperation(outOfBounds), false, '越界坐标应无效');
    console.log('    ✓ 异常操作验证通过');

    // 4.8 淘淤至卧铁验证
    console.log('  [正常] 淘淤至卧铁验证...');
    const currentBed = 728.0;
    const needDredge = currentBed - WOLONG_ELEVATION;
    assert.ok(needDredge > 1.8 && needDredge < 1.9, `需淘淤${needDredge.toFixed(2)}m，应≈1.88`);
    const dredgePerOpMin = 0.1, dredgePerOpMax = 0.5;
    const minOps = Math.ceil(needDredge / dredgePerOpMax);
    const maxOps = Math.ceil(needDredge / dredgePerOpMin);
    console.log(`    最少${minOps}次，最多${maxOps}次操作可达卧铁`);
    assert.ok(minOps > 0 && maxOps < 20, '操作次数范围合理');
    console.log('    ✓ 淘淤验证通过');

    // 4.9 渐进式效率验证
    console.log('  [正常] 效率递增不回退...');
    let prev = 0;
    for (let i = 0; i <= 30; i++) {
        const curr = calcMachaEff(i, 0.04);
        assert.ok(curr >= prev, `第${i}步效率回退：${prev}→${curr}`);
        prev = curr;
    }
    console.log('    ✓ 效率递增通过');
}

// ==============================================================
// 测试汇总
// ==============================================================
function runAllTests() {
    let passed = 0, failed = 0;
    const tests = [
        ['工艺对比 - 截流效率', testDynastyMachaInterception],
        ['跨时代对比 - 维护成本', testModernMaintenanceCost],
        ['地震模拟 - 结构安全', testEarthquakeStructureSafety],
        ['虚拟体验 - 操作合理性', testUserExperienceOperations]
    ];

    console.log('='.repeat(60));
    console.log('都江堰岁修系统 - 新增功能测试套件');
    console.log('='.repeat(60));

    for (const [name, fn] of tests) {
        try {
            fn();
            console.log(`\n  ✓ ${name} 通过`);
            passed++;
        } catch (err) {
            console.log(`\n  ✗ ${name} 失败: ${err.message}`);
            console.error(err.stack);
            failed++;
        }
    }

    console.log('\n' + '='.repeat(60));
    console.log(`测试完成：通过 ${passed}，失败 ${failed}，共 ${tests.length} 项`);
    console.log('='.repeat(60));

    if (failed > 0) {
        process.exit(1);
    }
}

if (require.main === module) {
    runAllTests();
}

module.exports = {
    testDynastyMachaInterception,
    testModernMaintenanceCost,
    testEarthquakeStructureSafety,
    testUserExperienceOperations,
    runAllTests
};