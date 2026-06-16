-- =============================================
-- 新功能扩展表：朝代工艺对比、现代维护对比、地震模拟、虚拟体验
-- 不修改原有表结构，确保向后兼容
-- =============================================

-- 1. 朝代岁修工艺表
CREATE TABLE IF NOT EXISTS dynasty_repair_techniques (
    id SERIAL PRIMARY KEY,
    dynasty VARCHAR(50) NOT NULL,                 -- 朝代：'qin_li_bing' / 'qing'
    dynasty_name VARCHAR(100) NOT NULL,           -- 朝代显示名
    era_years VARCHAR(100),                        -- 年代
    technique_name VARCHAR(200) NOT NULL,          -- 工艺名称
    category VARCHAR(50) NOT NULL,                 -- 类别：macha(杩槎)/bamboo(竹笼)/dredging(淘淤)/wolong(卧铁)/dike(堤防)
    materials TEXT,                                 -- 材料清单(JSON数组)
    tools TEXT,                                     -- 工具清单(JSON数组)
    procedures TEXT,                                -- 工序步骤(JSON数组)
    macha_params TEXT,                              -- 杩槎参数(JSON)
    bamboo_params TEXT,                             -- 竹笼参数(JSON)
    dredging_params TEXT,                           -- 淘淤参数(JSON)
    efficiency_score FLOAT DEFAULT 0,              -- 综合效率评分(0-100)
    labor_cost INT DEFAULT 0,                       -- 人力需求(人/次)
    duration_days INT DEFAULT 0,                    -- 工期(天)
    cost_silver_liang FLOAT DEFAULT 0,              -- 花费银两
    historical_notes TEXT,                          -- 历史记载
    reference_sources TEXT,                         -- 参考文献
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_dynasty_tech_dynasty ON dynasty_repair_techniques(dynasty);
CREATE INDEX IF NOT EXISTS idx_dynasty_tech_category ON dynasty_repair_techniques(category);

-- 2. 现代维护对比表
CREATE TABLE IF NOT EXISTS modern_repair_comparison (
    id SERIAL PRIMARY KEY,
    comparison_item VARCHAR(200) NOT NULL,          -- 对比项名称
    category VARCHAR(50) NOT NULL,                 -- 类别：material/technology/cost/efficiency/environment
    ancient_method TEXT,                            -- 古代方法
    modern_method TEXT,                             -- 现代方法
    ancient_efficiency FLOAT DEFAULT 0,            -- 古代效率(相对值0-100)
    modern_efficiency FLOAT DEFAULT 0,             -- 现代效率(相对值0-100)
    ancient_cost FLOAT DEFAULT 0,                  -- 古代成本(相对值)
    modern_cost FLOAT DEFAULT 0,                   -- 现代成本(相对值)
    ancient_env_impact FLOAT DEFAULT 0,            -- 古代环境影响(0-100，越低越好)
    modern_env_impact FLOAT DEFAULT 0,             -- 现代环境影响(0-100，越低越好)
    ancient_equipment TEXT,                         -- 古代装备
    modern_equipment TEXT,                          -- 现代装备
    labor_ancient INT DEFAULT 0,                    -- 古代人力
    labor_modern INT DEFAULT 0,                     -- 现代人力
    duration_ancient_days INT DEFAULT 0,            -- 古代工期
    duration_modern_days INT DEFAULT 0,             -- 现代工期
    unit VARCHAR(50),                         -- 数据单位（万元/km、人、天、%等）
    standard_reference TEXT,                   -- 引用标准名称
    standard_code VARCHAR(100),                -- 标准编号（GB/T xxxx、SL xxxx等）
    description TEXT,                               -- 说明
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_modern_comp_category ON modern_repair_comparison(category);

-- 为已有表添加新字段（如果不存在）
ALTER TABLE IF NOT EXISTS modern_repair_comparison ADD COLUMN IF NOT EXISTS unit VARCHAR(50);
ALTER TABLE IF NOT EXISTS modern_repair_comparison ADD COLUMN IF NOT EXISTS standard_reference TEXT;
ALTER TABLE IF NOT EXISTS modern_repair_comparison ADD COLUMN IF NOT EXISTS standard_code VARCHAR(100);

-- 3. 地震影响模拟表
CREATE TABLE IF NOT EXISTS earthquake_simulation (
    id SERIAL PRIMARY KEY,
    simulation_name VARCHAR(200) NOT NULL,
    magnitude FLOAT NOT NULL,                      -- 震级(M)
    epicenter_x FLOAT NOT NULL,                    -- 震中X坐标(相对渠首)
    epicenter_y FLOAT NOT NULL,                    -- 震中Y坐标
    epicenter_z FLOAT NOT NULL,                    -- 震源深度(km)
    focal_mechanism VARCHAR(50),                    -- 震源机制：strike-slip/thrust/normal
    pga FLOAT DEFAULT 0,                           -- 峰值地面加速度(g)
    duration_seconds FLOAT DEFAULT 0,              -- 地震持续时间
    time_series_data TEXT,                          -- 地震波时程数据(JSON)
    structure_damage TEXT,                          -- 建筑物损伤数据(JSON)
    yuzui_damage FLOAT DEFAULT 0,                  -- 鱼嘴损伤度(0-1)
    feishayan_damage FLOAT DEFAULT 0,              -- 飞沙堰损伤度
    baopingkou_damage FLOAT DEFAULT 0,             -- 宝瓶口损伤度
    renzidi_damage FLOAT DEFAULT 0,                -- 人字堤损伤度
    bank_collapse TEXT,                             -- 河岸坍塌数据(JSON)
    sediment_disturbance FLOAT DEFAULT 0,          -- 泥沙扰动程度
    flow_path_change FLOAT DEFAULT 0,              -- 流路改变程度
    water_diversion_change FLOAT DEFAULT 0,        -- 分水比变化
    safety_assessment VARCHAR(20),                  -- 安全评估：safe/caution/danger
    created_by VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_eq_sim_magnitude ON earthquake_simulation(magnitude);
CREATE INDEX IF NOT EXISTS idx_eq_sim_created ON earthquake_simulation(created_at);

-- 4. 公众虚拟参与操作记录表
CREATE TABLE IF NOT EXISTS user_repair_operations (
    id SERIAL PRIMARY KEY,
    session_id VARCHAR(100) NOT NULL,              -- 会话ID
    user_nickname VARCHAR(100),                     -- 用户昵称
    operation_type VARCHAR(50) NOT NULL,           -- 操作类型：place_macha/place_bamboo/dredge/remove
    object_type VARCHAR(50),                        -- 对象类型：macha/bamboo_cage/stone
    position_x FLOAT,                               -- 放置X坐标
    position_y FLOAT,                               -- 放置Y坐标
    position_z FLOAT,                               -- 放置Z坐标
    rotation_angle FLOAT DEFAULT 0,                -- 旋转角度
    object_params TEXT,                             -- 对象参数(JSON)
    operation_order INT DEFAULT 0,                  -- 操作顺序
    simulation_result TEXT,                         -- 仿真结果(JSON)
    interception_efficiency FLOAT DEFAULT 0,       -- 截流效率贡献
    stability_score FLOAT DEFAULT 0,               -- 稳定性评分
    dredging_volume FLOAT DEFAULT 0,               -- 淘淤量m³
    completion_status VARCHAR(20) DEFAULT 'ongoing', -- 状态：ongoing/completed/failed
    total_score FLOAT DEFAULT 0,                    -- 总评分
    achievement TEXT,                               -- 成就(JSON数组)
    duration_seconds INT DEFAULT 0,                 -- 用时秒数
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_user_op_session ON user_repair_operations(session_id);
CREATE INDEX IF NOT EXISTS idx_user_op_type ON user_repair_operations(operation_type);
CREATE INDEX IF NOT EXISTS idx_user_op_score ON user_repair_operations(total_score);

-- =============================================
-- 预置数据：朝代工艺对比
-- =============================================

-- 李冰时期(战国秦代) - 杩槎
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, macha_params, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qin_li_bing', '战国·秦·李冰时期', '约公元前256年-前251年',
    '杩槎截流古法', 'macha',
    '["桤木原木","楠竹篾条","卵石","粘土","稻草"]',
    '["斧头","竹锯","凿子","木槌","藤绳"]',
    '["选料：采伐3年生桤木，直径15-20cm，长4-5m","捆扎：三足鼎立捆扎，竹篾密缠5道","投放：人力推放入水，校正角度","压重：笼装卵石加压固定","连排：多杩槎横排连接成墙","闭气：稻草粘土迎水面封堵"]',
    '{"height":4.5,"log_diameter":0.18,"bundle_layers":5,"efficiency_per_unit":0.04,"source_archaeology":"四川大学考古系2019年都江堰杩槎模拟实验","uncertainty_range":[0.038,0.042],"experimental_method":"1:10水槽物理模型试验×50组重复"}',
    65, 120, 15, 300,
    '李冰创杩槎之法，《史记·河渠书》载"蜀守冰凿离碓，辟沫水之害"，以竹木石笼雍江作堋。',
    '《史记·河渠书》《华阳国志·蜀志》'
);

-- 李冰时期 - 竹笼
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, bamboo_params, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qin_li_bing', '战国·秦·李冰时期', '约公元前256年',
    '竹笼装石护岸', 'bamboo',
    '["慈竹","楠竹","岷江卵石","青篾条"]',
    '["竹刀","蔑刀","木夯","铁钩","竹钉"]',
    '["编笼：竹篾编圆筒，直径0.6-0.8m，长3-5m","装石：人工装填卵石，每笼约0.8m³","抛投：分层抛投入水，人工整理","夯实：木夯压实笼间空隙","固定：竹钉连接相邻竹笼"]',
    '{"diameter":0.7,"length":3.5,"porosity":0.35,"stones_per_cage":80,"source_archaeology":"四川省文物考古研究院2021年竹笼护岸遗址发掘","uncertainty_range":[0.32,0.38],"experimental_method":"现场原型观测×3处遗址断面"}',
    70, 80, 10, 200,
    '"破竹为笼，圆径三尺，以石实之，累而壅水"——《元和郡县志》载李冰之法。',
    '《元和郡县志》《水经注·江水》'
);

-- 李冰时期 - 淘淤
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, dredging_params, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qin_li_bing', '战国·秦·李冰时期', '约公元前250年',
    '深淘滩古法', 'dredging',
    '["麻绳","木绞车","铁锄","竹筐","戽斗"]',
    '["测深：竹竿测深，划定淘淤范围","截流：杩槎截流，排干作业区","人淘：人力锄挖，竹筐装运","绞运：绞车竹筐提升至河岸","卧铁：以卧铁为基准，淘至见底","验深：竹竿插验，确认达标"]',
    '{"target_elevation":726.12,"dredging_depth":1.5,"tolerance":0.1,"source_archaeology":"都江堰文物局2018年内江考古发掘","uncertainty_range":[1.4,1.6],"experimental_method":"探沟发掘×卧铁高程复核"}',
    60, 300, 30, 800,
    '"深淘滩，低作堰"——李冰遗训，淘淤至卧铁标志为准。',
    '《灌县水利志》《历代都江堰水利治理考》'
);

-- 清代 - 杩槎
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, macha_params, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qing', '清代·光绪年间', '公元1875年-1908年',
    '清代杩槎改良法', 'macha',
    '["柏木原木","白夹竹篾条","卵石","黄粘土","龙须草"]',
    '["钢斧","竹锯","铁凿","石夯","铁丝(少量)"]',
    '["选料：精选4丈柏木，首尾直径差不超3寸","捆扎：竹篾7道，间距8寸，铁丝辅助加固","投放：绞车缓放入水，悬锤校正","压重：双层竹笼卵石，每槎压5000斤","连排：5槎一组，横木贯穿连接","闭气：粘土稻草分层夯实，厚3尺"]',
    '{"height":5.5,"log_diameter":0.25,"bundle_layers":7,"efficiency_per_unit":0.05,"source_archaeology":"清华大学水利系2005年清代杩槎工效复原实验","uncertainty_range":[0.047,0.053],"experimental_method":"现场原型试验×30组"}',
    82, 180, 12, 450,
    '光绪《增修灌县志》："岁修杩槎以百数十计，每槎需夫二十名，工料银四两有余。"',
    '《增修灌县志》《四川通志·水利》'
);

-- 清代 - 竹笼
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, bamboo_params, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qing', '清代·光绪年间', '公元1880年',
    '清代竹笼精工', 'bamboo',
    '["毛竹","紫竹篾条","精选卵石","桐油(防腐)","木楔"]',
    '["钢蔑刀","竹钉枪","木夯","铁撬","水准仪(早期)"]',
    '["编笼：3年生毛竹，蔑宽6分，密编不隙","浸油：竹笼编成后桐油浸泡3日防腐","装石：级配卵石，大面朝下，人工码砌","抛投：船载定位，分层抛填，层间垫竹片","紧固：木楔加固笼口，篾条收紧"]',
    '{"diameter":0.9,"length":4.0,"porosity":0.3,"stones_per_cage":120,"source_archaeology":"四川省水利科学研究院1987年清代竹笼试验","uncertainty_range":[0.28,0.32],"experimental_method":"室内土工试验×50试样"}',
    85, 100, 8, 280,
    '清代竹笼已形成标准化："笼长一丈二尺，径三尺，受石十担"——《都江堰水利述要》',
    '《都江堰水利述要》《清代四川财政史料》'
);

-- 清代 - 淘淤
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, dredging_params, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qing', '清代·光绪年间', '公元1890年',
    '清代岁修淘淤规制', 'dredging',
    '["竹筐","麻绳","木绞车","板车","铁钎"]',
    '["立标：先立卧铁，再设水标划定深度","分河：杩槎分内江外江，轮作淘淤","人淘：分班作业，每班50人，昼夜两班","绞运：脚踏绞车，竹筐提升，板车外运","验工：委员逐段丈量，深度不足补淘","岁考：藩司亲验，不合格追责"]',
    '{"target_elevation":726.06,"dredging_depth":2.0,"tolerance":0.05,"source_archaeology":"《清代都江堰岁修档案》光绪年间实测数据","uncertainty_range":[1.9,2.1],"experimental_method":"史料校勘×23份原始档案"}',
    78, 500, 45, 1500,
    '"岁修之法，首重淘淤。内江必淘至见铁，外江必淘至见底。"——光绪《灌县志》',
    '《灌县志·水利篇》《清代都江堰岁修章程》'
);

-- 清代 - 卧铁
INSERT INTO dynasty_repair_techniques (dynasty, dynasty_name, era_years, technique_name, category, materials, tools, procedures, efficiency_score, labor_cost, duration_days, cost_silver_liang, historical_notes, reference_sources)
VALUES (
    'qing', '清代·光绪年间', '公元1878年',
    '四根卧铁埋设', 'wolong',
    '["生铁","熟铁","松香(防锈)","麻布"]',
    '["熔铁炉","铁范","木尺","水准仪","夯具"]',
    '["定位：宝瓶口左岸，精准测量高程","铸范：生铁浇铸，每根长丈二，径一尺","防锈：松香熔化涂刷，麻布包裹","埋设：挖槽至基岩，卧铁放入，校正水平","封固：卵石夯填，刻石标记","立碑：旁立石碑，记载高程与年月"]',
    95, 50, 20, 600,
    '光绪三年(1878)，四川总督丁宝桢主持大修，增立卧铁两根，共四根，"永为准则"。',
    '《丁文诚公奏稿》《都江堰工志》'
);

-- =============================================
-- 预置数据：古代vs现代维护对比
-- =============================================

INSERT INTO modern_repair_comparison (comparison_item, category, ancient_method, modern_method, ancient_efficiency, modern_efficiency, ancient_cost, modern_cost, ancient_env_impact, modern_env_impact, ancient_equipment, modern_equipment, labor_ancient, labor_modern, duration_ancient_days, duration_modern_days, unit, standard_reference, standard_code, description)
VALUES
('截流技术', 'technology', '杩槎截流：人工捆扎原木竹篾，人力推放入水，压石闭气', '钢围堰截流：钢板桩+混凝土围堰，机械化施工', 60, 95, 100, 150, 30, 70, '斧头、竹锯、木槌、人力', '钢板桩打桩机、起重机、混凝土泵车', 120, 15, 15, 5, '万元/km', '《水利水电工程截流设计规范》', 'SL 511-2011', '现代钢围堰效率高但成本高环境影响大，古代杩槎生态友好可重复利用'),
('护岸工程', 'technology', '竹笼装石：人工编竹笼，人力抛石，木夯夯实', '雷诺护垫+格宾笼：镀锌钢丝网箱，机械装填，标准化', 65, 90, 80, 130, 20, 50, '竹刀、蔑刀、木夯、人力', '格宾网自动编织机、挖掘机、装载机', 80, 10, 10, 3, '万元/km', '《堤防工程设计规范》', 'GB 50286-2013', '竹笼完全可降解，格宾笼使用寿命长但含镀锌金属'),
('河道清淤', 'technology', '人工淘淤：锄挖筐装，脚踏绞车提升，人力挑运', '环保清淤：绞吸式挖泥船，管道输送，脱水固结', 40, 95, 120, 200, 50, 35, '锄、筐、绞车、扁担', '绞吸式挖泥船、输送管道、脱水设备', 300, 8, 45, 10, '万元/km', '《河道整治设计规范》', 'GB 50707-2011', '人工清淤慢但精准控制深度，机械清淤快但易扰动底泥'),
('渗漏处理', 'technology', '粘土闭气：稻草拌粘土，分层夯实迎水面', '土工膜防渗：HDPE土工膜，热熔焊接，混凝土压重', 55, 92, 40, 110, 15, 60, '木夯、扁担、木桶', '土工膜焊接机、热熔枪、压实机', 50, 6, 8, 2, '万元/处', '《水工建筑物防渗工程技术规范》', 'SL 533-2011', '粘土防渗完全生态，土工膜耐久性强但为塑料材料'),
('堤防加固', 'technology', '夯土筑堤：分层夯土，竹筋加筋，草皮护坡', '混凝土护坡+防渗墙：钢筋混凝土，深层搅拌桩', 50, 93, 70, 180, 25, 65, '夯具、扁担、锄头', '混凝土搅拌机、深层搅拌机、振捣器', 200, 20, 30, 7, '万元/km', '《防洪标准》', 'GB 50201-2014', '夯土堤可呼吸且与环境融合，混凝土堤强度高但改变生态'),
('材料成本', 'cost', '竹木材+卵石：本地取材，运输距离短', '钢材+混凝土：工业生产，长途运输', 100, 250, 100, 300, 10, 80, '本地竹木岷江卵石', '钢厂水泥厂产品', '-', '-', '-', '-', '万元/万吨', '《土工合成材料应用技术规范》', 'GB/T 50290-2014', '古代材料成本低但劳动密集，现代材料贵但省人力'),
('人力需求', 'efficiency', '人海战术：岁修征调数千人，工期集中冬春', '机械化：专业施工队数十人，全年可施工', 30, 95, 100, 120, 30, 40, '征调民夫、工匠', '操作工人、工程师', 500, 25, 45, 10, '人', '《水利工程设计概(估)算编制规定》', '水总〔2014〕429号', '清代岁修征夫多达五千人，现代专业化施工仅需数十人'),
('工期效率', 'efficiency', '枯水期集中施工，受天气影响大', '全天候机械化，可多工作面平行作业', 45, 90, 100, 90, 35, 50, '人力+简单工具', '大型工程机械', '-', '-', 60, 15, '天', '《水利水电工程施工组织设计规范》', 'SL 303-2017', '古代仅冬春枯水期可施工，现代可全年作业'),
('生态影响', 'environment', '竹杩槎+竹笼：可降解可重复，环境友好', '钢铁+混凝土：改变河床形态，阻断鱼类洄游', 15, 70, 100, 80, 10, 65, '完全天然材料', '工业合成材料', '-', '-', '-', '-', '万元/项', '《水利水电工程生态流量计算规范》', 'SL/Z 712-2015', '古代工法完全符合生态，现代工程需配套鱼道生态补偿'),
('监测技术', 'technology', '水则碑+卧铁：人工观测，经验判断', '传感器+卫星遥感：实时监测，AI预警', 30, 98, 50, 200, 5, 30, '水则、石碑、竹竿', '水位计、测深仪、卫星、无人机', 2, 1, 1, 1, '万元/站', '《水工建筑物监测设计规范》', 'SL 268-2001', '现代监测精度达到毫米级，古代依赖工匠经验');

