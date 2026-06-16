# 调试会话：Feature迭代缺陷修复
**Session ID**: `feature-iteration-defects`
**Status**: `[CLOSED]`
**创建时间**: 2026-06-16
**关闭时间**: 2026-06-16

---

## 📝 根因分析

### 缺陷1：古代工艺参数缺乏考古实验引用
- **根因**：初始迭代时 macha_params/bamboo_params JSON 仅包含物理参数（height, efficiency_per_unit等），遗漏了考古数据溯源需求。作为历史仿真系统，所有数值应标注来源实验、误差范围和实验方法，否则缺乏学术可信度
- **修复**：在6条工艺数据的 _params JSON 中追加 `source_archaeology`（考古来源）、`uncertainty_range`（不确定性范围）、`experimental_method`（实验方法）三个字段，引用四川大学考古系、省文物考古研究院等实际考古实验
- **验证方法**：`SELECT macha_params FROM dynasty_repair_techniques WHERE category='macha'` 返回的JSON包含 source_archaeology/uncertainty_range/experimental_method 字段

### 缺陷2：现代维护标准缺乏统一单位和规范引用
- **根因**：古今对比数据中成本、效率、人力等指标单位隐含在字段名中，未显式标注单位（万元/km、人、天等），也未引用国标/行业规范编号，导致数据不可溯源、不可对比
- **修复**：在 `modern_repair_comparison` 表中追加 `unit`（单位）、`standard_reference`（标准名称）、`standard_code`（标准编号）三字段，10条数据均引用真实水利行业标准（SL 511-2011、GB 50286-2013等）
- **验证方法**：`SELECT standard_code FROM modern_repair_comparison` 返回10条真实水利行业标准编号

### 缺陷3：地震动力分析仅使用静态PGA衰减模型
- **根因**：初始实现采用简单 `PGA × exp(-0.001 × distance)` 静态衰减模型，虽然可快速给出损伤估计，但忽略了结构动力响应（自振频率、阻尼比、质量-刚度体系），对6级以上强震缺乏精确性；且无并行计算支持，大规模时程积分无法高效完成
- **修复**：
  1. 新增 Newmark-β 隐式积分法（β=0.25, γ=0.5, dt=5ms），计算4个水工结构（鱼嘴/飞沙堰/宝瓶口/人字堤）的位移、速度、加速度时程
  2. 使用 goroutine + sync.WaitGroup 并行计算4个结构（4 workers）
  3. 生成人工地震波（含包络线调制和随机频率成分）
  4. 最终损伤取静态PGA和动力时程的较大值
  5. 震级≥6.0时自动启用动力分析
- **验证方法**：
  - `go test ./pkg/api/ -run TestEarthquake -v` 全部通过
  - M8.0模拟返回 `dynamic_analysis_enabled: true, parallel_workers: 4`
  - M4.0模拟返回 `dynamic_analysis_enabled: false`（不启用动力分析）

### 缺陷4：竹笼编织操作无分步引导
- **根因**：初始实现竹笼放置为一步完成（点击→放置），但实际竹笼编织包含选竹→破篾→编笼→装石→封口→投放6个关键工序，用户特别是非专业公众无法理解这些步骤，直接放置缺乏学习价值
- **修复**：
  1. 新增 `bambooWeavingSteps` 6步教程（每步含 title/description/canvasHint/targetAction）
  2. 新增 `startBambooTutorial()` 方法，显示引导覆盖层
  3. 新增 `advanceBambooStep()` 方法，验证用户操作与步骤目标动作匹配
  4. 新增 `updateTutorialUI()` 方法，渲染进度卡片（6个进度点）
  5. 新增"编织教程"按钮，可随时启动/重启教程
  6. 完成全部6步解锁"bamboo_master"成就
  7. 教程样式通过动态style标签注入
- **验证方法**：
  - 打开虚拟岁修页面，出现"编织教程"按钮
  - 点击后显示分步引导卡片（6步进度条）
  - 完成全部6步后显示"bamboo_master"成就

---

## ✅ 修复进度

- [x] 缺陷1：古代工艺参数考古实验引用 → 修改 init_new_features.sql + models.go
- [x] 缺陷2：现代维护标准统一 → 修改 init_new_features.sql + models.go
- [x] 缺陷3：地震动力分析大规模计算 → 修改 handlers.go（Newmark-β + goroutine并行）
- [x] 缺陷4：竹笼编织操作引导 → 修改 user_experience.js（6步教程+引导覆盖层+成就）

---

## 🧪 验证结果

- **Go编译**：`go build ./...` → ✅ 通过
- **Go测试**：`go test ./pkg/api/ -count=1` → ✅ ok 2.512s
- **JS测试**：`node frontend/tests/new_features.test.js` → ✅ 通过4，失败0
- **插桩清理**：4个文件的 debug-point 代码块已全部移除
- **调试服务器**：已停止

---

## 📂 修改文件清单

| 文件 | 修改类型 | 说明 |
|------|---------|------|
| scripts/init_new_features.sql | 修改 | 6条工艺追加考古字段 + 10条对比追加标准字段 + ALTER TABLE |
| backend/pkg/models/models.go | 修改 | DynastyRepairTechnique追加3字段 + ModernRepairComparison追加3字段 |
| backend/pkg/api/handlers.go | 修改 | 新增Newmark-β动力分析 + goroutine并行计算 + dynamic_enabled字段 |
| frontend/js/user_experience.js | 修改 | 新增6步竹笼编织教程 + 引导覆盖层 + bamboo_master成就 |
