# 古代都江堰岁修工艺仿真与河床演变分析系统

## 项目概述

本系统是为水利史研究团队开发的全栈应用，用于模拟和分析古代都江堰的岁修工艺与河床演变过程。系统结合了三维可视化、离散元法仿真、水沙模型预测和实时数据监控等功能。

## 系统架构

```
┌─────────────────────────────────────────────────────────────────────────────┐
│                           前端 (Nginx + Gzip)                                │
│  ┌──────────┐  ┌──────────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ 3D模型   │  │ 河床演变面板  │  │ 岁修仿真面板  │  │ 告警中心        │   │
│  └──────────┘  └──────────────┘  └──────────────┘  └─────────────────┘   │
└─────────────────────────────┬───────────────────────────────────────────────┘
                              │ WebSocket + REST API
┌─────────────────────────────▼───────────────────────────────────────────────┐
│                         Go 后端服务 (:8080)                                  │
│  ┌───────────────────────────────────────────────────────────────────────┐  │
│  │                        消息总线 (Channel)                              │  │
│  └─────┬──────────┬──────────────┬───────────────────┬──────────────────┘  │
│  ┌─────▼──┐  ┌───▼───────┐  ┌───▼──────────┐  ┌─────▼──────────┐        │
│  │DTU接收  │  │岁修仿真   │  │河床演变分析  │  │ 告警MQTT推送   │        │
│  │模块     │  │模块       │  │模块          │  │ 模块           │        │
│  └────┬───┘  └───────────┘  └──────────────┘  └─────┬──────────┘        │
│       │                                               │                     │
│       │  ┌───────────────┐  ┌──────────────────┐   │                     │
│       └──┤ Prometheus    │  │ pprof (:6060)    │   │                     │
│          │ /metrics      │  │ 性能分析         │   │                     │
│          └───────────────┘  └──────────────────┘   │                     │
└─────────────┬───────────────────────────────────────┼─────────────────────┘
              │                                       │
    ┌─────────▼─────────┐                   ┌───────▼──────────┐
    │  TimescaleDB      │                   │  MQTT Broker     │
    │  (时序数据库)     │                   │  (Mosquitto)     │
    │  - 超表存储       │                   │  - 告警推送       │
    │  - 连续聚合降采样 │                   │  - QoS 1         │
    │  - 保留策略       │                   └──────────────────┘
    │  - 数据压缩       │
    └───────────────────┘

┌─────────────────────────────────────────────────────────────────────────────┐
│                         都江堰水文模拟器                                      │
│  - 6种场景预设 (平水/洪水/枯水/岁修/高沙/冲刷)                               │
│  - 可配置水位、含沙量、流量倍率                                              │
│  - 8个监测站点模拟数据                                                      │
└─────────────────────────────────────────────────────────────────────────────┘
```

## 技术栈

### 后端技术栈
- **语言**: Go 1.21+
- **Web框架**: Gin v1.9.1
- **数据库**: PostgreSQL + TimescaleDB 2.13+
- **数据库驱动**: pgx/v5 v5.5.0
- **消息队列**: MQTT (Eclipse Paho)
- **实时通信**: WebSocket (Gorilla)
- **科学计算**: Gonum v0.14.0
- **监控指标**: Prometheus client_golang v1.18.0
- **性能分析**: pprof (Go内置)

### 前端技术栈
- **3D渲染**: Three.js v0.158.0
- **图表库**: Chart.js v4.4.0
- **UI**: 原生HTML5 + CSS3 + JavaScript
- **2D绘图**: Canvas API
- **并发计算**: Web Worker

### 基础设施
- **容器化**: Docker + Docker Compose
- **Web服务器**: Nginx (Gzip压缩)
- **监控**: Prometheus + Grafana (可选)
- **MQTT Broker**: Eclipse Mosquitto

## 核心功能模块

### 1. 数据采集与存储
- 每1小时模拟传感器上报水位、流量、含沙量、河床高程
- 8个监测站点（内江3个、外江2个、飞沙堰2个、人字堤1个）
- TimescaleDB时序数据库优化存储
- 连续聚合视图提供快速统计查询

### 2. 岁修工艺仿真
- **杩槎截流仿真**: 基于离散元法模拟杩槎结构的截流过程
- **竹笼装石仿真**: 离散元法模拟石块填充竹笼的物理过程
- 空间哈希网格优化碰撞检测
- 多goroutine并行计算

### 3. 河床演变分析
- 基于水沙模型的河床冲淤计算
- 堰坝边界条件修正（飞沙堰、宝瓶口、人字堤）
- 未来10年河床高程预测
- 动态高程图可视化

### 4. 三维可视化
- 都江堰渠首三维地形模型
- 水流粒子动画系统
- 水利工程结构3D模型
- 卧铁标记与监测站点可视化

### 5. 告警系统
- 河床淤积超过卧铁高程自动触发预警
- 三级告警级别（严重/警告/通知）
- MQTT消息推送
- 数据库触发器自动创建告警

## 目录结构

```
AI_solo_coder_task_A_130/
├── backend/                          # Go后端服务
│   ├── cmd/
│   │   └── server/
│   │       └── main.go               # 主服务入口 (Application编排器)
│   ├── pkg/
│   │   ├── msg/
│   │   │   └── message.go            # 消息总线协议定义
│   │   ├── metrics/
│   │   │   ├── metrics.go            # Prometheus指标定义
│   │   │   └── middleware.go         # Gin Prometheus中间件
│   │   ├── dtu_receiver/
│   │   │   └── receiver.go           # DTU数据采集模块
│   │   ├── maintenance_simulator/
│   │   │   └── simulator.go          # 岁修工艺仿真模块
│   │   ├── riverbed_analyzer/
│   │   │   └── analyzer.go           # 河床演变分析模块
│   │   ├── alarm_mqtt/
│   │   │   └── alarm.go              # 告警评估和MQTT推送
│   │   ├── models/
│   │   │   ├── models.go             # 数据模型定义
│   │   │   └── database.go           # 数据库操作
│   │   └── simulation/               # 旧模块(保留兼容)
│   ├── Dockerfile                    # 多阶段构建Dockerfile
│   ├── .env                          # 环境变量
│   ├── .env.example                  # 环境变量示例
│   └── go.mod                        # Go依赖
├── frontend/                         # 前端应用
│   ├── index.html                    # 主页面
│   ├── css/
│   └── js/
│       ├── main.js                   # 主应用逻辑
│       ├── dujiangyan_3d.js          # 都江堰3D模块
│       ├── riverbed_panel.js         # 河床演变面板模块
│       └── simulation/
│           └── dem-worker.js         # Web Worker高程计算
├── scripts/                          # 脚本文件
│   ├── init_timescaledb.sql          # 数据库初始化脚本
│   ├── timescaledb_retention.sql     # 降采样和保留策略脚本
│   └── hydrology_simulator.go        # 水文数据模拟器
├── config/                           # 配置文件
│   ├── craft_params.json             # 工艺参数配置
│   └── sediment_params.json          # 水沙参数配置
├── docker/                           # Docker配置
│   ├── mosquitto/
│   │   └── mosquitto.conf            # Mosquitto配置
│   ├── nginx/
│   │   └── nginx.conf                # Nginx配置(Gzip)
│   └── prometheus/
│       └── prometheus.yml            # Prometheus配置
├── docker-compose.yml                # Docker Compose编排
├── Dockerfile.simulator              # 模拟器Dockerfile
└── README.md                         # 项目文档
```

## 快速开始

### 方式一：Docker Compose 部署（推荐）

#### 1. 环境要求
- Docker 20.10+
- Docker Compose v2+
- 至少 2GB 可用内存

#### 2. 启动基础服务

```bash
# 克隆项目
cd AI_solo_coder_task_A_130

# 启动核心服务 (TimescaleDB + MQTT + Go后端)
docker-compose up -d

# 查看服务状态
docker-compose ps
```

#### 3. 启动水文模拟器

```bash
# 启动默认模拟器(平水期场景)
docker-compose --profile simulator up -d

# 或启动洪水场景
docker-compose --profile flood up -d

# 或启动枯水期场景
docker-compose --profile drought up -d
```

#### 4. 启动Nginx（前端静态资源）

```bash
docker-compose --profile nginx up -d
```

#### 5. 启动监控栈（Prometheus + Grafana）

```bash
docker-compose --profile monitoring up -d
```

#### 6. 访问服务

| 服务 | 地址 | 说明 |
|------|------|------|
| 前端 | http://localhost/ | 通过Nginx访问(Gzip压缩) |
| API | http://localhost:8080/api/v1 | Go后端API |
| Metrics | http://localhost:8080/metrics | Prometheus指标 |
| pprof | http://localhost:6060/debug/pprof/ | Go性能分析 |
| Prometheus | http://localhost:9090 | Prometheus UI |
| Grafana | http://localhost:3000 | Grafana仪表板 |
| MQTT | localhost:1883 | MQTT Broker端口 |

#### 7. 停止服务

```bash
# 停止所有服务
docker-compose down

# 停止并删除数据卷(清空数据)
docker-compose down -v
```

### 方式二：本地开发部署

#### 1. 环境要求
- Go 1.21+
- PostgreSQL 14+ / TimescaleDB 2.13+
- MQTT Broker (如 Mosquitto)
- 现代浏览器（支持WebGL）

#### 2. 数据库初始化

```sql
-- 创建数据库
CREATE DATABASE dujiangyan;

-- 连接数据库
\c dujiangyan

-- 执行初始化脚本
\i scripts/init_timescaledb.sql

-- 执行保留策略和降采样脚本
\i scripts/timescaledb_retention.sql
```

#### 3. 后端配置

复制 `.env.example` 为 `.env` 并修改配置：

```bash
cd backend
cp .env.example .env
```

#### 4. 编译运行后端

```bash
cd backend
go mod download
go build -o dujiangyan-server cmd/server/main.go
./dujiangyan-server
```

#### 5. 运行水文模拟器

```bash
cd scripts
go run hydrology_simulator.go \
  -api http://localhost:8080/api/v1 \
  -interval 1s \
  -speed 3600 \
  -scenario normal
```

#### 6. 启动前端

```bash
# 使用Python启动简单服务器
cd frontend
python -m http.server 8081

# 或使用Nginx (推荐, 支持Gzip)
# 参考 docker/nginx/nginx.conf 配置
```

## 都江堰水文模拟器使用说明

### 查看可用场景

```bash
go run scripts/hydrology_simulator.go -list-scenarios
```

### 场景预设

| 场景 | 说明 | 水位 | 流量 | 含沙量 | 淤积速率 |
|------|------|------|------|--------|----------|
| `normal` | 平水期 - 正常水文条件 | 基准 | 1.0x | 基准 | 1.0x |
| `flood` | 丰水期/洪水 | +3.0m | 2.0x | +1.5, 2.5x | 3.0x |
| `drought` | 枯水期 | -2.5m | 0.4x | -0.3, 0.5x | 0.3x |
| `maintenance` | 岁修期 - 截流后冲刷 | -4.0m | 0.15x | -0.5, 0.3x | -0.5x |
| `high_sediment` | 高含沙期 | +1.0m | 1.2x | +3.0, 4.0x | 5.0x |
| `erosion` | 冲刷期 - 河床下切 | +2.0m | 1.8x | -0.2, 0.8x | -2.0x |

### 命令行参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `-api` | string | `http://localhost:8080/api/v1` | API服务器地址 |
| `-interval` | duration | `1s` | 数据上报间隔(真实时间) |
| `-speed` | float | `1.0` | 时间加速倍数(小时/tick) |
| `-historical` | int | `0` | 预生成历史数据天数 |
| `-scenario` | string | `normal` | 场景预设名称 |
| `-water-level-bonus` | float | `0` | 水位偏移量(米), 覆盖场景设置 |
| `-sediment-bonus` | float | `0` | 含沙量偏移量(kg/m³), 覆盖场景设置 |
| `-flow-multiplier` | float | `1.0` | 流量倍率, 覆盖场景设置 |
| `-list-scenarios` | bool | `false` | 列出所有可用场景并退出 |

### 使用示例

```bash
# 1. 正常模式, 1秒=1小时
go run hydrology_simulator.go -speed 3600

# 2. 洪水场景, 快速模拟(1秒=24小时)
go run hydrology_simulator.go -scenario flood -speed 86400

# 3. 自定义高水位高含沙
go run hydrology_simulator.go \
  -water-level-bonus 5.0 \
  -sediment-bonus 2.0 \
  -flow-multiplier 1.5

# 4. 先生成30天历史数据, 再实时模拟
go run hydrology_simulator.go -historical 30 -speed 3600

# 5. Docker方式启动洪水场景
docker-compose --profile flood up -d
```

## 监控与性能分析

### Prometheus 指标

访问 `http://localhost:8080/metrics` 查看所有指标。

#### 核心指标分类

| 分类 | 指标名称 | 说明 |
|------|----------|------|
| HTTP | `http_requests_total` | HTTP请求总数 |
| HTTP | `http_request_duration_seconds` | HTTP请求耗时直方图 |
| HTTP | `http_requests_in_flight` | 正在处理的请求数 |
| 水文 | `hydrology_data_received_total` | 接收的水文数据点数 |
| 水文 | `hydrology_water_level_meters` | 最新水位(米) |
| 水文 | `hydrology_sediment_concentration_kg_per_m3` | 最新含沙量 |
| 水文 | `hydrology_flow_rate_m3_per_s` | 最新流量 |
| 水文 | `hydrology_bed_elevation_meters` | 最新河床高程 |
| 仿真 | `simulations_started_total` | 仿真启动数 |
| 仿真 | `simulations_completed_total` | 仿真完成数 |
| 仿真 | `simulation_duration_seconds` | 仿真耗时直方图 |
| 预测 | `predictions_started_total` | 预测启动数 |
| 预测 | `prediction_duration_seconds` | 预测耗时直方图 |
| 告警 | `alerts_triggered_total` | 告警触发数 |
| 告警 | `alerts_acknowledged_total` | 告警确认数 |
| 总线 | `bus_messages_total` | 总线消息数 |
| 总线 | `bus_messages_dropped_total` | 丢弃的消息数 |
| WebSocket | `websocket_connections` | 活跃连接数 |
| WebSocket | `websocket_messages_sent_total` | 发送消息数 |
| MQTT | `mqtt_messages_published_total` | 发布消息数 |
| MQTT | `mqtt_reconnects_total` | 重连次数 |
| Go运行时 | `go_goroutines` | Goroutine数量 |
| Go运行时 | `go_memory_alloc_bytes` | 已分配内存字节数 |

#### Prometheus配置

Prometheus配置文件位于 `docker/prometheus/prometheus.yml`：

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'dujiangyan-server'
    static_configs:
      - targets: ['server:8080']
    metrics_path: /metrics
```

### pprof 性能分析

Go pprof端点默认监听在 `:6060` 端口。

#### 常用命令

```bash
# 查看goroutine概况
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 查看CPU使用情况(采样30秒)
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# 查看内存使用情况
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看阻塞分析
go tool pprof http://localhost:6060/debug/pprof/block

# 查看锁竞争
go tool pprof http://localhost:6060/debug/pprof/mutex
```

#### pprof Web界面

```bash
# 启动Web界面查看CPU profile
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile?seconds=30
```

然后访问 `http://localhost:8081` 查看火焰图、调用图等。

## TimescaleDB 降采样与保留策略

### 数据保留策略

| 数据表 | 保留时长 | 说明 |
|--------|----------|------|
| 原始水文数据 | 2年 | 高频原始数据 |
| 小时级聚合 | 1年 | 每小时统计数据 |
| 日级聚合 | 5年 | 每日统计数据 |
| 月级聚合 | 永久 | 长期趋势分析 |
| 告警数据 | 5年 | 告警记录 |
| 杩槎仿真数据 | 1年 | 仿真时序数据 |
| 河床演变预测 | 10年 | 预测结果 |

### 数据压缩策略

| 数据表 | 压缩延迟 | 分段字段 | 排序字段 |
|--------|----------|----------|----------|
| 水文数据 | 7天 | station_id | time DESC |
| 告警数据 | 30天 | station_id, alert_level | alert_time DESC |
| 杩槎仿真 | 1天 | simulation_id | time DESC |

### 连续聚合视图

- `hydrology_hourly` - 小时级聚合, 每30分钟刷新
- `hydrology_daily` - 日级聚合, 每2小时刷新
- `hydrology_monthly` - 月级聚合, 每天刷新
- `alerts_daily_stats` - 告警日统计, 每小时刷新

### 查询保留策略状态

```sql
-- 查看所有保留策略
SELECT * FROM get_retention_info();

-- 查看压缩状态
SELECT * FROM get_compression_info();

-- 查看所有后台任务
SELECT * FROM timescaledb_information.jobs;

-- 查看连续聚合
SELECT * FROM timescaledb_information.continuous_aggregates;
```

## 前端静态资源优化

### Gzip 压缩

Nginx已配置Gzip压缩，支持以下类型：
- 文本文件 (HTML, CSS, JS, JSON, XML)
- 字体文件 (TTF, OTF, WOFF, WOFF2)
- 图片 (SVG, ICO)
- WebAssembly

压缩级别：6（平衡压缩率与CPU消耗）
最小压缩文件大小：1KB

### 缓存策略

静态资源缓存 7 天，带 `immutable` 标志。

### Nginx 配置

配置文件位于 `docker/nginx/nginx.conf`，包含：
- Gzip 压缩
- 静态资源缓存
- API 反向代理
- WebSocket 反向代理

## Docker 镜像说明

### Go后端服务镜像

使用多阶段构建，最终镜像基于 scratch，仅包含静态二进制：

```dockerfile
# Builder阶段: Go编译环境
FROM golang:1.21-alpine AS builder
# 编译为静态二进制 (CGO_ENABLED=0)

# Runtime阶段: 空镜像
FROM scratch
# 只包含二进制文件和时区信息
```

镜像特性：
- 体积小 (~20MB)
- 无多余依赖
- 静态编译，可移植性强
- 暴露端口: 8080 (API) / 6060 (pprof)

### 模拟器镜像

同样使用多阶段构建，基于 scratch，体积最小化。

## API 接口文档

### 水文数据接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/hydrology/data` | 上报水文数据 |
| GET | `/api/v1/hydrology/:station_id/history` | 查询历史数据 |
| GET | `/api/v1/hydrology/:station_id/latest` | 获取最新数据 |
| GET | `/api/v1/hydrology/latest` | 获取所有站点最新数据 |
| GET | `/api/v1/hydrology/:station_id/daily-stats` | 获取日统计数据 |
| GET | `/api/v1/hydrology/:station_id/evolution-rate` | 获取演变率 |
| GET | `/api/v1/stations` | 获取监测站点列表 |

### 预测与仿真接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/prediction/bed-evolution` | 运行河床演变预测 |
| POST | `/api/v1/simulation/bamboo-cage` | 运行竹笼装石仿真 |
| POST | `/api/v1/simulation/macha-interception` | 运行杩槎截流仿真 |
| GET | `/api/v1/dem-grid` | 获取DEM高程网格 |
| GET | `/api/v1/annual-repair/records` | 获取岁修记录 |

### 告警接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/alerts` | 获取告警列表 |
| PUT | `/api/v1/alerts/:id/acknowledge` | 确认告警 |

### 监控接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/metrics` | Prometheus指标 |
| GET | `/debug/pprof/` | pprof性能分析 |

### WebSocket接口

- `ws://host:port/api/v1/ws` - 实时数据推送

推送消息类型：
- `hydrology_data` - 水文数据更新
- `simulation_result` - 仿真结果
- `prediction_result` - 预测结果

## 配置文件

### 工艺参数配置 (config/craft_params.json)

```json
{
  "dem": {
    "gravity": 9.81,
    "restitution": 0.3,
    "friction": 0.98
  },
  "bamboo_cage": { ... },
  "macha": { ... },
  "terrain": { ... }
}
```

### 水沙参数配置 (config/sediment_params.json)

```json
{
  "sediment_transport": { ... },
  "weirs": [
    { "name": "飞沙堰", "type": "overflow", ... },
    { "name": "宝瓶口", "type": "orifice", ... },
    { "name": "人字堤", "type": "overflow", ... }
  ],
  "prediction": { ... },
  "evolution_rate": { ... }
}
```

## 故障排除

### 后端无法连接数据库
- 检查PostgreSQL和TimescaleDB是否正常运行
- 确认 `.env` 中的数据库配置正确
- 确认数据库已执行初始化脚本
- Docker部署时确认服务名正确 (timescaledb)

### MQTT连接失败
- 检查MQTT Broker是否运行
- 确认端口1883是否开放
- 系统支持无MQTT运行（告警仅存储不推送）
- Docker部署时确认服务名正确 (mosquitto)

### 前端3D场景不显示
- 确认浏览器支持WebGL
- 检查控制台是否有错误信息
- 尝试刷新页面重新加载

### 粒子动画不流畅
- 降低粒子数量（2000以下）
- 关闭其他浏览器标签释放资源
- 使用性能更好的显卡

### Docker容器无法启动
```bash
# 查看容器日志
docker-compose logs server
docker-compose logs timescaledb

# 检查端口占用
netstat -an | findstr "8080 5432 1883"
```

### Prometheus无数据
- 确认target是否可达
- 检查 `/metrics` 端点是否正常返回
- 查看Prometheus UI的Targets页面

## 开发说明

### 后端开发
```bash
cd backend
go run cmd/server/main.go
```

### 前端开发
前端使用纯HTML/JS/CSS，无需构建工具，直接修改文件即可。

### 添加新的监测指标
1. 在 `init_timescaledb.sql` 中添加对应字段
2. 更新 `models.go` 中的数据结构
3. 在 `database.go` 添加相应的查询方法
4. 在 `metrics.go` 中添加对应指标

## 许可证

本项目仅供水利史研究使用。

## 参考文献

1. 《都江堰水利工程史》
2. 《离散元法及其在岩土工程中的应用》
3. 《泥沙运动力学》
4. TimescaleDB官方文档
5. Three.js官方文档
6. Prometheus官方文档
