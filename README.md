# MiMo TTS Service

基于小米 MiMo API 的文本转语音（TTS）服务，后端使用 Go + Gin，前端使用 Vue 3 + Vite，部署在 Kubernetes 上。

## 功能

- 文本转语音：调用 MiMo `mimo-v2-tts` 模型，将文本转换为 MP3 音频
- 风格支持：温柔、悄悄话、东北话、开心、严肃等语音风格
- 批量转换：一次最多转换 10 段文本
- Vue 前端界面：文本输入、风格选择、在线播放

## 项目结构

```
├── main.go                  # 入口，Gin 路由配置
├── config/config.go         # 环境变量加载
├── handler/tts_handler.go   # HTTP 请求处理
├── service/tts_service.go   # MiMo API 调用逻辑
├── model/model.go           # 数据结构定义
├── utils/response.go        # 统一响应格式
├── static/index.html        # 静态前端（备用）
├── frontend/                # Vue 3 前端项目
│   ├── src/App.vue
│   ├── nginx.conf
│   └── Dockerfile
├── k8s/                     # Kubernetes 部署文件
│   ├── deployment.yaml      # 后端 Deployment
│   ├── service.yaml         # 后端 Service
│   ├── frontend-deployment.yaml
│   ├── frontend-service.yaml
│   ├── ingress.yaml         # Ingress 路由
│   └── secret.yaml.example  # Secret 模板
├── Dockerfile               # 后端多阶段构建
├── .env.example             # 环境变量模板
└── go.mod
```

## API

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/v1/tts/convert` | 单条文本转语音，返回 audio/mpeg |
| POST | `/api/v1/tts/batch` | 批量文本转语音 |
| GET  | `/api/v1/tts/voices` | 获取可用音色列表 |
| GET  | `/health` | 健康检查 |

### 请求示例

```bash
curl -X POST http://localhost:8080/api/v1/tts/convert \
  -H "Content-Type: application/json" \
  -d '{"text": "你好，欢迎使用 MiMo 语音合成服务。", "style": "温柔"}' \
  --output output.mp3
```

## 本地开发

### 后端

```bash
# 复制环境变量
cp .env.example .env
# 编辑 .env 填入你的 API Key

# 安装依赖
go mod download

# 启动
go run main.go
```

服务默认运行在 `http://localhost:8080`。

### 前端

```bash
cd frontend
npm install
npm run dev
```

前端开发服务器运行在 `http://localhost:5173`，API 请求会代理到后端 8080 端口。

## Docker 构建

```bash
# 后端
docker build -t your-registry/tts-server:v1 .

# 前端
docker build -t your-registry/tts-frontend:v1 ./frontend
```

## Kubernetes 部署

```bash
# 1. 创建 Secret（填入真实 API Key）
kubectl create secret generic tts-secret --from-literal=MIMO_API_KEY=your_key

# 2. 部署后端
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# 3. 部署前端
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml

# 4. 配置 Ingress
kubectl apply -f k8s/ingress.yaml
```

Ingress 路由规则：
- `/api/*` → 后端 `tts-service`
- `/` → 前端 `tts-frontend-service`

## 环境变量

| 变量 | 说明 | 默认值 |
|------|------|--------|
| `MIMO_API_KEY` | MiMo API 密钥 | 必填 |
| `MIMO_ENDPOINT` | MiMo API 地址 | `https://api.xiaomimimo.com/v1/chat/completions` |
| `SERVER_PORT` | 服务端口 | `8080` |

## 技术栈

- **后端**：Go 1.22, Gin
- **前端**：Vue 3, Vite
- **部署**：Docker, Kubernetes, Nginx Ingress
