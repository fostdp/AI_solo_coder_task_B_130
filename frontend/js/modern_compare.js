let modernComparisonData = [];
let modernComparisonGrouped = {};
let modernCardCharts = {};

function generateMockModernComparisons() {
    const categories = Object.keys(CONFIG.COMPARISON_CATEGORIES);
    const mockItems = [
        {
            name: '截流',
            category: 'technology',
            ancient_method: '杩槎截流、竹笼装石、人工打桩、经验判断截流时机',
            modern_method: '钢板桩围堰、液压截流机、GPS定位、水力模型仿真',
            ancient_efficiency: 45,
            modern_efficiency: 95,
            ancient_cost: 100,
            modern_cost: 250,
            ancient_environment: 70,
            modern_environment: 85,
            ancient_manpower: 90,
            modern_manpower: 25,
            ancient_duration: 85,
            modern_duration: 30,
            ancient_precision: 35,
            modern_precision: 92
        },
        {
            name: '护岸',
            category: 'technology',
            ancient_method: '竹笼护脚、条石护坡、夯土加固、木桩固基',
            modern_method: '钢筋混凝土挡墙、土工格栅、生态护坡、防渗膜',
            ancient_efficiency: 40,
            modern_efficiency: 90,
            ancient_cost: 90,
            modern_cost: 220,
            ancient_environment: 60,
            modern_environment: 90,
            ancient_manpower: 85,
            modern_manpower: 20,
            ancient_duration: 80,
            modern_duration: 35,
            ancient_precision: 40,
            modern_precision: 88
        },
        {
            name: '清淤',
            category: 'efficiency',
            ancient_method: '人工淘淤、竹篮搬运、卧铁标记、季节性作业',
            modern_method: '绞吸式挖泥船、高压水枪、管道输送、GPS精准定位',
            ancient_efficiency: 35,
            modern_efficiency: 92,
            ancient_cost: 85,
            modern_cost: 180,
            ancient_environment: 55,
            modern_environment: 75,
            ancient_manpower: 95,
            modern_manpower: 15,
            ancient_duration: 90,
            modern_duration: 25,
            ancient_precision: 30,
            modern_precision: 85
        },
        {
            name: '监测',
            category: 'technology',
            ancient_method: '水尺观测、人工巡查、卧铁比对、经验判断',
            modern_method: '自动化传感器、卫星遥感、无人机巡检、AI预警系统',
            ancient_efficiency: 30,
            modern_efficiency: 98,
            ancient_cost: 80,
            modern_cost: 280,
            ancient_environment: 80,
            modern_environment: 95,
            ancient_manpower: 88,
            modern_manpower: 10,
            ancient_duration: 70,
            modern_duration: 5,
            ancient_precision: 25,
            modern_precision: 96
        },
        {
            name: '材料',
            category: 'material',
            ancient_method: '竹、木、卵石、条石、石灰、黏土',
            modern_method: '钢筋混凝土、高分子材料、合金钢、环保复合材料',
            ancient_efficiency: 50,
            modern_efficiency: 88,
            ancient_cost: 85,
            modern_cost: 200,
            ancient_environment: 85,
            modern_environment: 70,
            ancient_manpower: 75,
            modern_manpower: 30,
            ancient_duration: 65,
            modern_duration: 40,
            ancient_precision: 45,
            modern_precision: 90
        },
        {
            name: '渗漏',
            category: 'material',
            ancient_method: '黏土封堵、夯土防渗、竹席挡水、经验堵漏',
            modern_method: '灌浆堵漏、防渗膜、化学注浆、超声波检测',
            ancient_efficiency: 38,
            modern_efficiency: 93,
            ancient_cost: 95,
            modern_cost: 260,
            ancient_environment: 70,
            modern_environment: 80,
            ancient_manpower: 82,
            modern_manpower: 22,
            ancient_duration: 75,
            modern_duration: 28,
            ancient_precision: 32,
            modern_precision: 94
        },
        {
            name: '堤防',
            category: 'technology',
            ancient_method: '夯土筑堤、竹笼护岸、条石压顶、人工夯实',
            modern_method: '机械化碾压、土工合成材料、防渗墙、位移监测',
            ancient_efficiency: 42,
            modern_efficiency: 91,
            ancient_cost: 105,
            modern_cost: 240,
            ancient_environment: 65,
            modern_environment: 82,
            ancient_manpower: 92,
            modern_manpower: 18,
            ancient_duration: 88,
            modern_duration: 32,
            ancient_precision: 38,
            modern_precision: 89
        },
        {
            name: '人力',
            category: 'efficiency',
            ancient_method: '征调民夫、简单工具、经验施工、人海战术',
            modern_method: '专业队伍、大型机械、标准化作业、技能培训',
            ancient_efficiency: 33,
            modern_efficiency: 87,
            ancient_cost: 110,
            modern_cost: 190,
            ancient_environment: 75,
            modern_environment: 78,
            ancient_manpower: 98,
            modern_manpower: 12,
            ancient_duration: 92,
            modern_duration: 28,
            ancient_precision: 28,
            modern_precision: 86
        },
        {
            name: '工期',
            category: 'efficiency',
            ancient_method: '枯水期作业、季节性限制、天气依赖、人工调度',
            modern_method: '全年施工、围堰挡水、全天候作业、进度管理系统',
            ancient_efficiency: 36,
            modern_efficiency: 94,
            ancient_cost: 120,
            modern_cost: 300,
            ancient_environment: 68,
            modern_environment: 72,
            ancient_manpower: 90,
            modern_manpower: 25,
            ancient_duration: 95,
            modern_duration: 20,
            ancient_precision: 30,
            modern_precision: 91
        },
        {
            name: '生态',
            category: 'environment',
            ancient_method: '尊重自然、因势利导、材料环保、生态影响小',
            modern_method: '生态修复、环境评估、生物多样性保护、绿色施工',
            ancient_efficiency: 55,
            modern_efficiency: 85,
            ancient_cost: 80,
            modern_cost: 150,
            ancient_environment: 90,
            modern_environment: 88,
            ancient_manpower: 70,
            modern_manpower: 35,
            ancient_duration: 60,
            modern_duration: 45,
            ancient_precision: 50,
            modern_precision: 93
        }
    ];

    return mockItems.map((item, index) => ({
        id: index + 1,
        ...item
    }));
}

function initCategoryDropdown() {
    const select = document.getElementById('modern-category');
    if (!select) return;

    Object.entries(CONFIG.COMPARISON_CATEGORIES).forEach(([key, value]) => {
        const option = document.createElement('option');
        option.value = key;
        option.textContent = value;
        select.appendChild(option);
    });

    select.addEventListener('change', (e) => {
        filterAndRenderComparisons(e.target.value);
    });
}

function groupByCategory(data) {
    const grouped = {};
    data.forEach(item => {
        if (!grouped[item.category]) {
            grouped[item.category] = [];
        }
        grouped[item.category].push(item);
    });
    return grouped;
}

function filterAndRenderComparisons(category) {
    let filteredData;
    if (category === 'all') {
        filteredData = modernComparisonData;
    } else {
        filteredData = modernComparisonData.filter(item => item.category === category);
    }
    renderComparisonCards(filteredData);
    updateEfficiencyBarChart(filteredData);
}

function renderComparisonCards(data) {
    const container = document.getElementById('modern-cards');
    if (!container) return;

    Object.values(modernCardCharts).forEach(chart => {
        if (chart) chart.destroy();
    });
    modernCardCharts = {};

    container.innerHTML = '';

    data.forEach((item, index) => {
        const card = document.createElement('div');
        card.className = 'compare-card';
        card.innerHTML = `
            <h4>${item.name} <span style="font-size:12px;color:#8892a0;font-weight:normal;">[${CONFIG.COMPARISON_CATEGORIES[item.category] || item.category}]</span></h4>
            <div class="compare-row">
                <div class="compare-col ancient">
                    <h5>古代</h5>
                    <p>${item.ancient_method}</p>
                </div>
                <div class="compare-arrow">→</div>
                <div class="compare-col modern">
                    <h5>现代</h5>
                    <p>${item.modern_method}</p>
                </div>
            </div>
            <canvas id="card-chart-${item.id}"></canvas>
        `;
        container.appendChild(card);

        setTimeout(() => {
            createCardChart(item);
        }, 50);
    });
}

function createCardChart(item) {
    const ctx = document.getElementById(`card-chart-${item.id}`);
    if (!ctx) return;

    const chart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: ['效率', '成本', '环境', '人力', '工期'],
            datasets: [
                {
                    label: '古代',
                    data: [
                        item.ancient_efficiency,
                        item.ancient_cost,
                        item.ancient_environment,
                        item.ancient_manpower,
                        item.ancient_duration
                    ],
                    backgroundColor: 'rgba(255, 170, 0, 0.6)',
                    borderColor: '#ffaa00',
                    borderWidth: 1
                },
                {
                    label: '现代',
                    data: [
                        item.modern_efficiency,
                        item.modern_cost,
                        item.modern_environment,
                        item.modern_manpower,
                        item.modern_duration
                    ],
                    backgroundColor: 'rgba(0, 212, 255, 0.6)',
                    borderColor: '#00d4ff',
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

    modernCardCharts[item.id] = chart;
}

function createRadarChart(data) {
    const ctx = document.getElementById('modern-radar-chart');
    if (!ctx) return;

    if (window.modernRadarChart) {
        window.modernRadarChart.destroy();
    }

    const avgAncient = {
        technology: data.length ? data.reduce((s, i) => s + i.ancient_efficiency, 0) / data.length : 40,
        cost: data.length ? data.reduce((s, i) => s + i.ancient_cost, 0) / data.length : 95,
        efficiency: data.length ? data.reduce((s, i) => s + i.ancient_manpower, 0) / data.length : 85,
        environment: data.length ? data.reduce((s, i) => s + i.ancient_environment, 0) / data.length : 70,
        duration: data.length ? data.reduce((s, i) => s + i.ancient_duration, 0) / data.length : 80,
        precision: data.length ? data.reduce((s, i) => s + i.ancient_precision, 0) / data.length : 35
    };

    const avgModern = {
        technology: data.length ? data.reduce((s, i) => s + i.modern_efficiency, 0) / data.length : 92,
        cost: data.length ? data.reduce((s, i) => s + i.modern_cost, 0) / data.length : 230,
        efficiency: data.length ? data.reduce((s, i) => s + i.modern_manpower, 0) / data.length : 20,
        environment: data.length ? data.reduce((s, i) => s + i.modern_environment, 0) / data.length : 82,
        duration: data.length ? data.reduce((s, i) => s + i.modern_duration, 0) / data.length : 30,
        precision: data.length ? data.reduce((s, i) => s + i.modern_precision, 0) / data.length : 90
    };

    window.modernRadarChart = new Chart(ctx, {
        type: 'radar',
        data: {
            labels: ['技术效率', '成本投入', '人力需求', '环境影响', '工期耗时', '施工精度'],
            datasets: [
                {
                    label: '古代',
                    data: [
                        avgAncient.technology,
                        avgAncient.cost,
                        avgAncient.efficiency,
                        avgAncient.environment,
                        avgAncient.duration,
                        avgAncient.precision
                    ],
                    borderColor: '#ffaa00',
                    backgroundColor: 'rgba(255, 170, 0, 0.2)',
                    pointBackgroundColor: '#ffaa00',
                    pointBorderColor: '#fff',
                    pointHoverBackgroundColor: '#fff',
                    pointHoverBorderColor: '#ffaa00'
                },
                {
                    label: '现代',
                    data: [
                        avgModern.technology,
                        avgModern.cost,
                        avgModern.efficiency,
                        avgModern.environment,
                        avgModern.duration,
                        avgModern.precision
                    ],
                    borderColor: '#00d4ff',
                    backgroundColor: 'rgba(0, 212, 255, 0.2)',
                    pointBackgroundColor: '#00d4ff',
                    pointBorderColor: '#fff',
                    pointHoverBackgroundColor: '#fff',
                    pointHoverBorderColor: '#00d4ff'
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
                    grid: { color: 'rgba(255,255,255,0.1)' },
                    angleLines: { color: 'rgba(255,255,255,0.1)' },
                    pointLabels: { color: '#aaa', font: { size: 12 } },
                    ticks: { color: '#aaa', backdropColor: 'transparent' }
                }
            },
            plugins: {
                legend: { labels: { color: '#ccc' } }
            }
        }
    });
}

function updateEfficiencyBarChart(data) {
    const ctx = document.getElementById('modern-efficiency-chart');
    if (!ctx) return;

    if (window.modernEfficiencyChart) {
        window.modernEfficiencyChart.destroy();
    }

    window.modernEfficiencyChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: data.map(item => item.name),
            datasets: [
                {
                    label: '古代效率 (%)',
                    data: data.map(item => item.ancient_efficiency),
                    backgroundColor: 'rgba(255, 170, 0, 0.6)',
                    borderColor: '#ffaa00',
                    borderWidth: 1
                },
                {
                    label: '现代效率 (%)',
                    data: data.map(item => item.modern_efficiency),
                    backgroundColor: 'rgba(0, 212, 255, 0.6)',
                    borderColor: '#00d4ff',
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

async function loadModernComparisons() {
    try {
        const data = await API.dynasty.getModernComparisons();
        if (data && data.length > 0) {
            modernComparisonData = data;
        } else {
            modernComparisonData = generateMockModernComparisons();
        }
    } catch (error) {
        console.error('加载古今对比数据失败:', error);
        modernComparisonData = generateMockModernComparisons();
    }

    modernComparisonGrouped = groupByCategory(modernComparisonData);
    renderComparisonCards(modernComparisonData);
    createRadarChart(modernComparisonData);
    updateEfficiencyBarChart(modernComparisonData);
}

document.addEventListener('DOMContentLoaded', () => {
    initCategoryDropdown();
    loadModernComparisons();
});
