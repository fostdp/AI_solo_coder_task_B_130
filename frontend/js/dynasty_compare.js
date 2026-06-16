(function(global) {
let qinLiBingData = [];
let qingData = [];
let currentCategory = 'all';

const DYNASTY_MOCK = {
    qin_li_bing: {
        dynasty_name: '李冰时期（战国秦）',
        techniques: [
            {
                id: 1,
                name: '杩槎截流',
                category: 'macha',
                era: '公元前256年',
                materials: ['木材', '竹篾', '卵石', '黏土'],
                tools: ['斧头', '锯子', '竹刀', '夯锤'],
                steps: [
                    '选取三根六至八米长的圆木，将其一端削尖',
                    '用竹篾将三根圆木绑扎成三角架，称为"杩槎"',
                    '将杩槎立于江中，尖脚插入河床',
                    '在杩槎内侧堆放卵石竹笼增加稳定性',
                    '用黏土填充杩槎缝隙形成截流墙'
                ],
                efficiency_score: 65,
                labor_cost: 300,
                duration_days: 30,
                cost_silver: 800,
                historical_record: '《华阳国志·蜀志》载："冰乃壅江作堋，穿郫江、检江，别支流双过郡下"',
                references: ['史记·河渠书', '华阳国志', '灌县志']
            },
            {
                id: 2,
                name: '竹笼装石',
                category: 'bamboo',
                era: '公元前256年',
                materials: ['慈竹', '卵石', '竹篾'],
                tools: ['竹刀', '篾刀', '筐篓', '绳索'],
                steps: [
                    '选取三年以上慈竹，劈成宽约2厘米的竹篾',
                    '用竹篾编织成直径约1米、长约10米的圆筒形竹笼',
                    '在竹笼两端留口，中间装填卵石',
                    '用竹篾将两端封口，防止卵石散出',
                    '将竹笼层层堆叠于需要加固的堤岸'
                ],
                efficiency_score: 68,
                labor_cost: 250,
                duration_days: 25,
                cost_silver: 600,
                historical_record: '《元和郡县志》载："破竹为笼，圆径三尺，以石实中，累而壅水"',
                references: ['元和郡县志', '都江堰水利述要', '四川通志']
            },
            {
                id: 3,
                name: '古法淘淤',
                category: 'dredging',
                era: '公元前250年',
                materials: ['麻绳', '木橇', '竹筐'],
                tools: ['锄头', '铁锹', '挑担', '畚箕'],
                steps: [
                    '每年枯水期（农历十月至次年二月）进行淘淤',
                    '先用杩槎截流，将内江江水导入外江',
                    '工人进入干涸的河床，用锄头和铁锹挖掘淤积泥沙',
                    '将泥沙装入竹筐，人力挑运至外江',
                    '淘至卧铁标记处为止，确保河床深度达标'
                ],
                efficiency_score: 62,
                labor_cost: 500,
                duration_days: 60,
                cost_silver: 1200,
                historical_record: '《深淘滩，低作堰》："岁淘必及根，深则不能，浅则不可"',
                references: ['都江堰岁修章程', '治水三字经', '灌江备考']
            },
            {
                id: 4,
                name: '卧铁埋设',
                category: 'wolong',
                era: '公元前250年',
                materials: ['生铁', '木炭', '黏土'],
                tools: ['熔炉', '风箱', '铁锤', '量尺'],
                steps: [
                    '选定内江凤栖窝河段作为淘淤基准点',
                    '用生铁铸造长约一丈（约3.3米）的铁柱',
                    '在河床底部挖深坑，将卧铁埋入',
                    '卧铁顶部与规定淘淤深度齐平',
                    '记录埋设位置，作为后世岁修淘淤的标准'
                ],
                efficiency_score: 70,
                labor_cost: 150,
                duration_days: 15,
                cost_silver: 500,
                historical_record: '《灌县志》载："沿埋铁石于滩底，岁修以此为准，谓之卧铁"',
                references: ['灌县志', '都江堰工程见闻录', '四川水利志']
            }
        ]
    },
    qing: {
        dynasty_name: '清代',
        techniques: [
            {
                id: 101,
                name: '改良杩槎',
                category: 'macha',
                era: '清光绪年间（1875-1908）',
                materials: ['柏木', '楠竹', '卵石', '黏土', '条石'],
                tools: ['铁斧', '钢锯', '竹刀', '铁夯', '撬棍'],
                steps: [
                    '精选优质柏木，每根直径不小于30厘米',
                    '采用"三搭六"绑扎法，增加杩槎整体稳定性',
                    '杩槎底部增设横木，增大与河床接触面积',
                    '内侧先堆竹笼卵石，外层加砌条石加固',
                    '黏土中掺加石灰，增强防渗效果'
                ],
                efficiency_score: 82,
                labor_cost: 500,
                duration_days: 45,
                cost_silver: 1500,
                historical_record: '《都江堰水利述要》载："清代杩槎，用料更精，绑扎更固，截流功效倍于前代"',
                references: ['都江堰水利述要', '清代四川水利档案', '灌县乡土志']
            },
            {
                id: 102,
                name: '精工竹笼',
                category: 'bamboo',
                era: '清乾隆年间（1736-1795）',
                materials: ['楠竹', '卵石', '竹篾', '棕绳'],
                tools: ['钢刃竹刀', '篾刀', '绞绳器', '铁钩'],
                steps: [
                    '选用深山楠竹，竹龄五年以上，竹节均匀',
                    '竹篾劈为四层，层层叠压编织，增加笼体强度',
                    '笼身每米设五道竹箍，防止装填后鼓胀变形',
                    '卵石选用花岗岩，直径15-25厘米，大小搭配',
                    '叠放时采用"品字形"交错排列，增强整体稳定性'
                ],
                efficiency_score: 85,
                labor_cost: 450,
                duration_days: 40,
                cost_silver: 1200,
                historical_record: '《灌江备考》载："竹笼之制，至清而精，篾密石匀，虽大水不溃"',
                references: ['灌江备考', '清代岁修章程', '四川通志·水利志']
            },
            {
                id: 103,
                name: '冬季淘淤',
                category: 'dredging',
                era: '清嘉庆年间（1796-1820）',
                materials: ['麻绳', '木橇', '竹筐', '滑轨'],
                tools: ['铁锄', '钢锹', '独轮车', '绞车', '畚箕'],
                steps: [
                    '立冬后截流，比前代提前半月开工',
                    '河床铺设滑轨，用独轮车替代人力挑运',
                    '设置绞车滑轮组，将泥沙从深沟中提升',
                    '分区作业，每区设工长负责质量监督',
                    '淘淤至四根卧铁以下三寸，确保深度达标'
                ],
                efficiency_score: 80,
                labor_cost: 800,
                duration_days: 75,
                cost_silver: 2500,
                historical_record: '《都江堰岁修章程》载："每岁立冬兴工，立春蒇事，定额淘夫一千二百名"',
                references: ['都江堰岁修章程', '清代水利档案汇编', '蜀水考']
            },
            {
                id: 104,
                name: '卧铁补铸',
                category: 'wolong',
                era: '清同治三年（1864）',
                materials: ['灰口铁', '木炭', '耐火泥'],
                tools: ['反射炉', '鼓风机', '铁范', '卡尺', '水平仪'],
                steps: [
                    '核查前代卧铁位置，标记补铸点',
                    '采用反射炉熔铁，温度更高，铁质更纯',
                    '使用铁范铸造，卧铁尺寸精度大幅提高',
                    '埋设时使用水平仪校准，确保顶面水平',
                    '在卧铁上铸造年号和监造官员姓名'
                ],
                efficiency_score: 83,
                labor_cost: 200,
                duration_days: 20,
                cost_silver: 800,
                historical_record: '《灌县志》载："同治三年，续铸卧铁二根，埋于故址，以准淘淤"',
                references: ['灌县志（同治版）', '都江堰文物志', '四川金石志']
            }
        ]
    }
};

function initDynastyCompareDOM() {
    const container = document.getElementById('dynasty');
    if (!container) return;

    container.innerHTML = `
        <div class="dynasty-container">
            <h3>朝代工艺对比</h3>
            <div class="dynasty-controls">
                <div class="control-group">
                    <label>工艺分类</label>
                    <select id="dynasty-category">
                        <option value="all">全部</option>
                        <option value="macha">杩槎</option>
                        <option value="bamboo">竹笼</option>
                        <option value="dredging">淘淤</option>
                        <option value="wolong">卧铁</option>
                    </select>
                </div>
            </div>
            <div class="dynasty-columns">
                <div class="dynasty-column">
                    <h4>李冰时期（战国秦）</h4>
                    <div class="dynasty-cards" id="qin-cards"></div>
                </div>
                <div class="dynasty-column">
                    <h4>清代</h4>
                    <div class="dynasty-cards" id="qing-cards"></div>
                </div>
            </div>
            <div class="dynasty-charts">
                <div class="chart-card">
                    <h4>综合指标雷达图对比</h4>
                    <canvas id="dynasty-radar-chart"></canvas>
                </div>
                <div class="chart-card">
                    <h4>各类工艺效率评分对比</h4>
                    <canvas id="dynasty-bar-chart"></canvas>
                </div>
            </div>
        </div>
    `;
}

function parseJSONArray(value) {
    if (!value) return [];
    if (Array.isArray(value)) return value;
    try {
        if (typeof value === 'string') {
            const parsed = JSON.parse(value);
            return Array.isArray(parsed) ? parsed : [];
        }
    } catch (e) {
        return String(value).split(/[,，、;；]/).map(s => s.trim()).filter(Boolean);
    }
    return [];
}

function createTechniqueCard(technique, dynastyKey) {
    const materials = parseJSONArray(technique.materials);
    const tools = parseJSONArray(technique.tools);
    const steps = parseJSONArray(technique.steps);
    const references = parseJSONArray(technique.references);
    const categoryName = CONFIG.DYNAMIC_CATEGORIES[technique.category] || technique.category;

    const borderColor = dynastyKey === 'qin_li_bing' ? '#00a8ff' : '#ffaa00';
    const headerBg = dynastyKey === 'qin_li_bing' ? 'linear-gradient(135deg, rgba(0,168,255,0.2), rgba(0,168,255,0.05))' : 'linear-gradient(135deg, rgba(255,170,0,0.2), rgba(255,170,0,0.05))';

    return `
        <div class="technique-card" style="border-left: 4px solid ${borderColor};">
            <div class="card-header" style="background: ${headerBg};">
                <h5 class="technique-name">${technique.name || '未命名工艺'}</h5>
                <span class="technique-category">${categoryName}</span>
            </div>
            <div class="card-body">
                <div class="card-row">
                    <span class="row-label">年代：</span>
                    <span class="row-value">${technique.era || '--'}</span>
                </div>
                <div class="card-row">
                    <span class="row-label">材料：</span>
                    <span class="tags">
                        ${materials.map(m => `<span class="tag tag-material">${m}</span>`).join('')}
                    </span>
                </div>
                <div class="card-row">
                    <span class="row-label">工具：</span>
                    <span class="tags">
                        ${tools.map(t => `<span class="tag tag-tool">${t}</span>`).join('')}
                    </span>
                </div>
                <div class="card-section">
                    <span class="section-title">工序步骤：</span>
                    <ol class="steps-list">
                        ${steps.map(s => `<li>${s}</li>`).join('')}
                    </ol>
                </div>
                <div class="card-stats">
                    <div class="stat-mini">
                        <span class="stat-mini-label">效率评分</span>
                        <span class="stat-mini-value">${technique.efficiency_score || '--'}</span>
                    </div>
                    <div class="stat-mini">
                        <span class="stat-mini-label">人力</span>
                        <span class="stat-mini-value">${technique.labor_cost || '--'} 人</span>
                    </div>
                    <div class="stat-mini">
                        <span class="stat-mini-label">工期</span>
                        <span class="stat-mini-value">${technique.duration_days || '--'} 天</span>
                    </div>
                    <div class="stat-mini">
                        <span class="stat-mini-label">花费</span>
                        <span class="stat-mini-value">${technique.cost_silver || '--'} 两</span>
                    </div>
                </div>
                <div class="card-row">
                    <span class="row-label">历史记载：</span>
                </div>
                <div class="historical-record">${technique.historical_record || '暂无记载'}</div>
                <div class="card-row">
                    <span class="row-label">参考文献：</span>
                    <span class="tags">
                        ${references.map(r => `<span class="tag tag-reference">${r}</span>`).join('')}
                    </span>
                </div>
            </div>
        </div>
    `;
}

function renderCards() {
    const qinContainer = document.getElementById('qin-cards');
    const qingContainer = document.getElementById('qing-cards');
    if (!qinContainer || !qingContainer) return;

    const filteredQin = currentCategory === 'all'
        ? qinLiBingData
        : qinLiBingData.filter(t => t.category === currentCategory);

    const filteredQing = currentCategory === 'all'
        ? qingData
        : qingData.filter(t => t.category === currentCategory);

    qinContainer.innerHTML = filteredQin.length > 0
        ? filteredQin.map(t => createTechniqueCard(t, 'qin_li_bing')).join('')
        : '<div class="no-data">暂无该分类工艺数据</div>';

    qingContainer.innerHTML = filteredQing.length > 0
        ? filteredQing.map(t => createTechniqueCard(t, 'qing')).join('')
        : '<div class="no-data">暂无该分类工艺数据</div>';
}

function calculateRadarData(techniques) {
    if (!techniques || techniques.length === 0) {
        return [0, 0, 0, 0, 0];
    }

    const avgEfficiency = techniques.reduce((sum, t) => sum + (t.efficiency_score || 0), 0) / techniques.length;
    const avgLabor = techniques.reduce((sum, t) => sum + (t.labor_cost || 0), 0) / techniques.length;
    const avgDuration = techniques.reduce((sum, t) => sum + (t.duration_days || 0), 0) / techniques.length;
    const avgCost = techniques.reduce((sum, t) => sum + (t.cost_silver || 0), 0) / techniques.length;

    const laborScore = avgLabor > 0 ? Math.max(0, 100 - (avgLabor / 10)) : 0;
    const durationScore = avgDuration > 0 ? Math.max(0, 100 - (avgDuration / 1.2)) : 0;
    const costScore = avgCost > 0 ? Math.max(0, 100 - (avgCost / 30)) : 0;
    const innovationScore = avgEfficiency * 0.8 + 15;

    return [
        Math.min(100, Math.max(0, avgEfficiency)),
        Math.min(100, Math.max(0, laborScore)),
        Math.min(100, Math.max(0, durationScore)),
        Math.min(100, Math.max(0, costScore)),
        Math.min(100, Math.max(0, innovationScore))
    ];
}

function getCategoryEfficiency(techniques, category) {
    const filtered = techniques.filter(t => t.category === category);
    if (filtered.length === 0) return null;
    return filtered.reduce((sum, t) => sum + (t.efficiency_score || 0), 0) / filtered.length;
}

function renderCharts() {
    const radarLabels = ['效率', '人力指数', '工期指数', '成本指数', '创新度'];
    const qinRadar = calculateRadarData(qinLiBingData);
    const qingRadar = calculateRadarData(qingData);

    const radarCtx = document.getElementById('dynasty-radar-chart');
    if (radarCtx) {
        if (window.dynastyRadarChart) {
            window.dynastyRadarChart.destroy();
        }
        window.dynastyRadarChart = new Chart(radarCtx, {
            type: 'radar',
            data: {
                labels: radarLabels,
                datasets: [
                    {
                        label: '李冰时期',
                        data: qinRadar,
                        borderColor: '#00a8ff',
                        backgroundColor: 'rgba(0, 168, 255, 0.2)',
                        borderWidth: 2,
                        pointBackgroundColor: '#00a8ff',
                        pointBorderColor: '#fff',
                        pointRadius: 4
                    },
                    {
                        label: '清代',
                        data: qingRadar,
                        borderColor: '#ffaa00',
                        backgroundColor: 'rgba(255, 170, 0, 0.2)',
                        borderWidth: 2,
                        pointBackgroundColor: '#ffaa00',
                        pointBorderColor: '#fff',
                        pointRadius: 4
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    r: {
                        min: 0,
                        max: 100,
                        beginAtZero: true,
                        grid: { color: 'rgba(255,255,255,0.1)' },
                        angleLines: { color: 'rgba(255,255,255,0.1)' },
                        pointLabels: { color: '#ccc', font: { size: 12 } },
                        ticks: { color: '#aaa', backdropColor: 'transparent' }
                    }
                },
                plugins: {
                    legend: { labels: { color: '#ccc' } }
                }
            }
        });
    }

    const barCategories = ['macha', 'bamboo', 'dredging', 'wolong'];
    const barLabels = barCategories.map(c => CONFIG.DYNAMIC_CATEGORIES[c] || c);
    const qinBarData = barCategories.map(c => getCategoryEfficiency(qinLiBingData, c));
    const qingBarData = barCategories.map(c => getCategoryEfficiency(qingData, c));

    const barCtx = document.getElementById('dynasty-bar-chart');
    if (barCtx) {
        if (window.dynastyBarChart) {
            window.dynastyBarChart.destroy();
        }
        window.dynastyBarChart = new Chart(barCtx, {
            type: 'bar',
            data: {
                labels: barLabels,
                datasets: [
                    {
                        label: '李冰时期',
                        data: qinBarData,
                        backgroundColor: 'rgba(0, 168, 255, 0.6)',
                        borderColor: '#00a8ff',
                        borderWidth: 1
                    },
                    {
                        label: '清代',
                        data: qingBarData,
                        backgroundColor: 'rgba(255, 170, 0, 0.6)',
                        borderColor: '#ffaa00',
                        borderWidth: 1
                    }
                ]
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                scales: {
                    y: {
                        min: 0,
                        max: 100,
                        grid: { color: 'rgba(255,255,255,0.1)' },
                        ticks: { color: '#aaa' }
                    },
                    x: {
                        grid: { color: 'rgba(255,255,255,0.1)' },
                        ticks: { color: '#aaa' }
                    }
                },
                plugins: {
                    legend: { labels: { color: '#ccc' } }
                }
            }
        });
    }
}

async function loadDynastyData() {
    try {
        const qinData = await API.dynasty.getTechniques('qin_li_bing', null);
        qinLiBingData = Array.isArray(qinData) ? qinData : (qinData.techniques || []);
    } catch (error) {
        console.warn('加载李冰时期数据失败，使用mock数据:', error);
        qinLiBingData = DYNASTY_MOCK.qin_li_bing.techniques;
    }

    try {
        const qingResp = await API.dynasty.getTechniques('qing', null);
        qingData = Array.isArray(qingResp) ? qingResp : (qingResp.techniques || []);
    } catch (error) {
        console.warn('加载清代数据失败，使用mock数据:', error);
        qingData = DYNASTY_MOCK.qing.techniques;
    }
}

async function refreshByCategory() {
    const category = currentCategory === 'all' ? null : currentCategory;

    try {
        const qinData = await API.dynasty.getTechniques('qin_li_bing', category);
        qinLiBingData = Array.isArray(qinData) ? qinData : (qinData.techniques || []);
    } catch (error) {
        console.warn('加载李冰时期分类数据失败，使用mock过滤:', error);
        const allQin = DYNASTY_MOCK.qin_li_bing.techniques;
        qinLiBingData = category ? allQin.filter(t => t.category === category) : allQin;
    }

    try {
        const qingResp = await API.dynasty.getTechniques('qing', category);
        qingData = Array.isArray(qingResp) ? qingResp : (qingResp.techniques || []);
    } catch (error) {
        console.warn('加载清代分类数据失败，使用mock过滤:', error);
        const allQing = DYNASTY_MOCK.qing.techniques;
        qingData = category ? allQing.filter(t => t.category === category) : allQing;
    }

    renderCards();
    renderCharts();
}

document.addEventListener('DOMContentLoaded', async () => {
    initDynastyCompareDOM();

    const categorySelect = document.getElementById('dynasty-category');
    if (categorySelect) {
        categorySelect.addEventListener('change', (e) => {
            currentCategory = e.target.value;
            refreshByCategory();
        });
    }

    await loadDynastyData();
    renderCards();
    renderCharts();
});

window.initDynastyCompareDOM = initDynastyCompareDOM;
window.DynastyCompare = { init: initDynastyCompareDOM, loadData: loadDynastyData };
})(window);
