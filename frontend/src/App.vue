<template>
  <div class="app">
    <!-- 登录/注册弹窗 -->
    <div v-if="!currentUser" class="modal-overlay" @click.self="closeAuthModal">
      <div class="auth-modal">
        <div v-if="showLogin">
          <h2>登录</h2>
          <form @submit.prevent="login">
            <div class="form-group">
              <label>用户名/邮箱</label>
              <input v-model="loginForm.username" type="text" required placeholder="请输入用户名或邮箱" />
            </div>
            <div class="form-group">
              <label>密码</label>
              <input v-model="loginForm.password" type="password" required placeholder="请输入密码" />
            </div>
            <button type="submit" :disabled="loading">{{ loading ? '登录中...' : '登录' }}</button>
          </form>
          <p class="toggle-text">
            还没有账号？
            <a href="#" @click.prevent="showLogin = false">立即注册</a>
          </p>
        </div>
        <div v-else>
          <h2>注册</h2>
          <form @submit.prevent="register">
            <div class="form-group">
              <label>用户名</label>
              <input v-model="registerForm.username" type="text" required placeholder="请输入用户名" />
            </div>
            <div class="form-group">
              <label>邮箱（可选）</label>
              <input v-model="registerForm.email" type="email" placeholder="请输入邮箱" />
            </div>
            <div class="form-group">
              <label>密码（至少 6 位）</label>
              <input v-model="registerForm.password" type="password" required minlength="6" placeholder="请输入密码" />
            </div>
            <button type="submit" :disabled="loading">{{ loading ? '注册中...' : '注册' }}</button>
          </form>
          <p class="toggle-text">
            已有账号？
            <a href="#" @click.prevent="showLogin = true">立即登录</a>
          </p>
        </div>
        <div v-if="authError" class="error-msg">{{ authError }}</div>
      </div>
    </div>

    <!-- 侧边栏 -->
    <aside class="sidebar">
      <div class="logo" @click="showAuthModal = !showAuthModal">
        <div class="logo-icon">
          <svg viewBox="0 0 24 24" fill="currentColor">
            <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
          </svg>
        </div>
        <span class="logo-text">MiMo TTS</span>
      </div>

      <nav class="nav">
        <a href="#" :class="['nav-item', { active: currentView === 'tts' }]" @click.prevent="switchView('tts')">
          <svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/></svg>
          <span>语音合成</span>
        </a>
        <a href="#" :class="['nav-item', { active: currentView === 'history' }]" @click.prevent="switchView('history')">
          <svg viewBox="0 0 24 24" fill="currentColor"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>
          <span>历史记录</span>
        </a>
        <a href="#" class="nav-item" @click.prevent>
          <svg viewBox="0 0 24 24" fill="currentColor"><path d="M19.14 12.94c.04-.31.06-.63.06-.94 0-.31-.02-.63-.06-.94l2.03-1.58c.18-.14.23-.41.12-.61l-1.92-3.32c-.12-.22-.37-.29-.59-.22l-2.39.96c-.5-.38-1.03-.7-1.62-.94l-.36-2.54c-.04-.24-.24-.41-.48-.41h-3.84c-.24 0-.43.17-.47.41l-.36 2.54c-.59.24-1.13.57-1.62.94l-2.39-.96c-.22-.08-.47 0-.59.22L2.74 8.87c-.12.21-.08.47.12.61l2.03 1.58c-.04.31-.06.63-.06.94s.02.63.06.94l-2.03 1.58c-.18.14-.23.41-.12.61l1.92 3.32c.12.22.37.29.59.22l2.39-.96c.5.38 1.03.7 1.62.94l.36 2.54c.05.24.24.41.48.41h3.84c.24 0 .44-.17.47-.41l.36-2.54c.59-.24 1.13-.56 1.62-.94l2.39.96c.22.08.47 0 .59-.22l1.92-3.32c.12-.22.07-.47-.12-.61l-2.01-1.58zM12 15.6c-1.98 0-3.6-1.62-3.6-3.6s1.62-3.6 3.6-3.6 3.6 1.62 3.6 3.6-1.62 3.6-3.6 3.6z"/></svg>
          <span>设置</span>
        </a>
      </nav>

      <div class="sidebar-footer">
        <div v-if="currentUser" class="user-info" @click="showAuthModal = true">
          <div class="avatar">{{ getAvatar(currentUser.username) }}</div>
          <div class="user-details">
            <span class="username">{{ currentUser.username }}</span>
            <span class="api-key" @click.stop="copyAPIKey">点击复制 API Key</span>
          </div>
        </div>
        <div v-else class="user-info login-prompt" @click="showAuthModal = true">
          <div class="avatar">?</div>
          <span class="username">点击登录</span>
        </div>
      </div>
    </aside>

    <!-- 主内容区 -->
    <main class="main">
      <header class="header">
        <h1>语音合成</h1>
        <p class="subtitle">将文本转换为自然流畅的语音</p>
      </header>

      <div class="content">
        <!-- 左侧：输入区域 -->
        <div class="input-section">
          <div class="card">
            <div class="card-header">
              <h2>文本输入</h2>
              <span class="char-count">{{ text.length }}/1000</span>
            </div>
            <textarea
              v-model="text"
              placeholder="请输入要转换的文本，支持中英文混合输入..."
              maxlength="1000"
            ></textarea>
          </div>

          <div class="card">
            <div class="card-header">
              <h2>语音设置</h2>
            </div>
            <div class="settings-grid">
              <div class="setting-item">
                <label>风格</label>
                <select v-model="selectedStyle">
                  <option value="">默认</option>
                  <option value="温柔">温柔</option>
                  <option value="悄悄话">悄悄话</option>
                  <option value="东北话 开心">东北话 开心</option>
                  <option value="开心">开心</option>
                  <option value="严肃">严肃</option>
                </select>
              </div>
              <div class="setting-item">
                <label>语速</label>
                <div class="slider-container">
                  <input type="range" v-model="speed" min="0.5" max="2" step="0.1" />
                  <span class="slider-value">{{ speed }}x</span>
                </div>
              </div>
            </div>
          </div>

          <button class="generate-btn" :disabled="loading || !currentUser" @click="convert">
            <svg v-if="!loading" viewBox="0 0 24 24" fill="currentColor">
              <path d="M8 5v14l11-7z"/>
            </svg>
            <div v-else class="spinner"></div>
            <span>{{ loading ? '生成中...' : (!currentUser ? '请先登录' : '生成语音') }}</span>
          </button>

          <div class="error-msg" v-if="error">
            <svg viewBox="0 0 24 24" fill="currentColor">
              <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z"/>
            </svg>
            {{ error }}
          </div>
        </div>

        <!-- 右侧：输出区域 -->
        <div class="output-section">
          <div class="card output-card">
            <div class="card-header">
              <h2>音频输出</h2>
            </div>

            <div class="output-content" :class="{ 'has-audio': audioUrl }">
              <div v-if="!audioUrl" class="placeholder">
                <svg viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
                </svg>
                <p>生成的音频将在这里显示</p>
                <span>输入文本并点击"生成语音"</span>
              </div>

              <div v-else class="audio-player">
                <div class="waveform">
                  <div class="wave-bar" v-for="i in 40" :key="i" :style="{ animationDelay: `${i * 0.05}s` }"></div>
                </div>
                <audio ref="player" :src="audioUrl" controls></audio>
                <div class="audio-actions">
                  <button class="action-btn" @click="downloadAudio">
                    <svg viewBox="0 0 24 24" fill="currentColor">
                      <path d="M19 9h-4V3H9v6H5l7 7 7-7zM5 18v2h14v-2H5z"/>
                    </svg>
                    下载音频
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div class="card history-card">
            <div class="card-header">
              <h2>最近生成</h2>
              <div class="history-actions">
                <button class="clear-btn" @click="clearHistory" :disabled="history.length === 0">清空</button>
                <button class="refresh-btn" @click="loadHistory" :disabled="loading">刷新</button>
              </div>
            </div>
            <div class="history-list" v-if="history.length > 0">
              <div class="history-item" v-for="(item, index) in history" :key="index" @click="loadHistoryItem(item)">
                <div class="history-text">{{ item.text.slice(0, 30) }}{{ item.text.length > 30 ? '...' : '' }}</div>
                <div class="history-meta">{{ item.style || '默认' }} · {{ item.created_at }}</div>
              </div>
            </div>
            <div v-else class="empty-state">
              <svg viewBox="0 0 24 24" fill="currentColor">
                <path d="M13 3c-4.97 0-9 4.03-9 9H1l3.89 3.89.07.14L9 12H6c0-3.87 3.13-7 7-7s7 3.13 7 7-3.13 7-7 7c-1.93 0-3.68-.79-4.94-2.06l-1.42 1.42C8.27 19.99 10.51 21 13 21c4.97 0 9-4.03 9-9s-4.03-9-9-9zm-1 5v5l4.28 2.54.72-1.21-3.5-2.08V8H12z"/>
              </svg>
              <p>暂无历史记录</p>
              <span>首次生成后将在这里显示</span>
            </div>
          </div>
        </div>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const text = ref('你好，欢迎使用 MiMo 语音合成服务。')
const selectedStyle = ref('')
const speed = ref(1.0)
const loading = ref(false)
const audioUrl = ref('')
const error = ref('')
const player = ref(null)
const history = ref([])
const showAuthModal = ref(false)
const showLogin = ref(true)
const authError = ref('')
const currentUser = ref(null)
let currentUserAPIKey = ref('')
const currentView = ref("tts") // "tts" 或 "history"
// 登录表单
const loginForm = ref({
  username: '',
  password: ''
})

// 注册表单
const registerForm = ref({
  username: '',
  email: '',
  password: ''
})

const API_BASE = '/api/v1'

// 获取用户信息
async function loadUserInfo() {
  try {
    const res = await fetch(`${API_BASE}/user/info`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('api_key') || ''}`
      }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.code === 200 && data.data) {
        currentUser.value = data.data.user
        return
      }
    }
    // 尝试用 API key 获取信息
    const apiKeyRes = await fetch(`${API_BASE}/user/api-key`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('api_key') || ''}`
      }
    })
    if (apiKeyRes.ok) {
      const data = await apiKeyRes.json()
      if (data.code === 200 && data.data) {
        currentUserAPIKey.value = data.data.api_key
        localStorage.setItem('api_key', data.data.api_key)
      }
    }
    currentUser.value = null
  } catch (e) {
    currentUser.value = null
  }
}

// 加载历史记录
async function loadHistory() {
  if (!currentUser.value) return
  try {
    const res = await fetch(`${API_BASE}/history?page=1&page_size=10`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('api_key') || ''}`
      }
    })
    if (res.ok) {
      const data = await res.json()
      if (data.code === 200 && data.data) {
        history.value = data.data.records || []
      }
    }
  } catch (e) {
    console.error('加载历史失败:', e)
  }
}

// 从历史记录加载
function loadHistoryItem(item) {
  text.value = item.text
  selectedStyle.value = item.style || ''
}

// 清空历史记录
async function clearHistory() {
  if (!confirm('确定要清空所有历史记录吗？')) return
  try {
    const res = await fetch(`${API_BASE}/history`, {
      method: 'DELETE',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('api_key') || ''}`
      }
    })
    if (res.ok) {
      history.value = []
      error.value = '已清空历史记录'
    }
  } catch (e) {
    error.value = '清空失败：' + e.message
  }
}

// 登录
async function login() {
  authError.value = ''
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/auth/login`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(loginForm.value)
    })
    const data = await res.json()
    if (data.code === 200) {
      const apiKey = data.data.user.api_key
      localStorage.setItem('api_key', apiKey)
      currentUserAPIKey.value = apiKey
      await loadUserInfo()
      await loadHistory()
      showAuthModal.value = false
      authError.value = ''
    } else {
      authError.value = data.message || '登录失败'
    }
  } catch (e) {
    authError.value = '登录失败：' + e.message
  } finally {
    loading.value = false
  }
}

// 注册
async function register() {
  authError.value = ''
  loading.value = true
  try {
    const res = await fetch(`${API_BASE}/auth/register`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(registerForm.value)
    })
    const data = await res.json()
    if (data.code === 200) {
      const apiKey = data.data.user.api_key
      localStorage.setItem('api_key', apiKey)
      currentUserAPIKey.value = apiKey
      await loadUserInfo()
      showAuthModal.value = false
    } else {
      authError.value = data.message || '注册失败'
    }
  } catch (e) {
    authError.value = '注册失败：' + e.message
  } finally {
    loading.value = false
  }
}

// 关闭认证弹窗
function closeAuthModal() {
  showAuthModal.value = false
  authError.value = ''
}

// 复制 API Key
function copyAPIKey() {
  navigator.clipboard.writeText(currentUserAPIKey.value).then(() => {
    error.value = 'API Key 已复制'
  }).catch(() => {
    error.value = '复制失败'
  })
}

// 获取头像首字母
function getAvatar(username) {
  if (!username) return '?'
  return username.charAt(0).toUpperCase()
}

// 生成语音
async function convert() {
  const trimmed = text.value.trim()
  error.value = ''
  authError.value = ''

  if (!trimmed) {
    error.value = '请输入文本'
    return
  }
  if (trimmed.length > 1000) {
    error.value = '文本不能超过 1000 字'
    return
  }

  if (!currentUser.value) {
    error.value = '请先登录'
    showAuthModal.value = true
    return
  }

  if (audioUrl.value) {
    URL.revokeObjectURL(audioUrl.value)
    audioUrl.value = ''
  }

  loading.value = true

  try {
    const res = await fetch(`${API_BASE}/tts/convert`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${currentUserAPIKey.value}`
      },
      body: JSON.stringify({ text: trimmed, style: selectedStyle.value })
    })

    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.message || '请求失败')
    }

    const blob = await res.blob()
    audioUrl.value = URL.createObjectURL(blob)
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

// 下载音频
function downloadAudio() {
  if (!audioUrl.value) return
  const a = document.createElement('a')
  a.href = audioUrl.value
  a.download = `tts_${Date.now()}.mp3`
  a.click()
}

onMounted(async () => {
  const apiKey = localStorage.getItem('api_key')
  if (apiKey) {
    currentUserAPIKey.value = apiKey
    await loadUserInfo()
    await loadHistory()
  }
})
</script>

<style>
* {
  box-sizing: border-box;
  margin: 0;
  padding: 0;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  background: #0a0a0f;
  color: #e5e5e5;
}

.app {
  display: flex;
  min-height: 100vh;
  width: 100%;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.7);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.auth-modal {
  background: linear-gradient(180deg, #1a1a2e 0%, #16213e 100%);
  border: 1px solid rgba(99, 102, 241, 0.3);
  border-radius: 16px;
  padding: 32px;
  width: 100%;
  max-width: 400px;
}

.auth-modal h2 {
  font-size: 24px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 24px;
  text-align: center;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  font-size: 14px;
  color: #a1a1aa;
  margin-bottom: 8px;
}

.form-group input {
  width: 100%;
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  font-size: 15px;
  color: #e5e5e5;
  outline: none;
  transition: border-color 0.2s;
}

.form-group input:focus {
  border-color: rgba(99, 102, 241, 0.5);
}

.auth-modal button {
  width: 100%;
  padding: 14px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border: none;
  border-radius: 10px;
  font-size: 16px;
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
  margin-top: 8px;
}

.auth-modal button:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
}

.auth-modal button:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.toggle-text {
  text-align: center;
  margin-top: 20px;
  color: #71717a;
  font-size: 14px;
}

.toggle-text a {
  color: #818cf8;
  text-decoration: none;
  font-weight: 500;
}

.toggle-text a:hover {
  text-decoration: underline;
}

.error-msg {
  margin-top: 16px;
  padding: 12px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 8px;
  color: #f87171;
  font-size: 14px;
  text-align: center;
}

/* 侧边栏 */
.sidebar {
  width: 240px;
  background: linear-gradient(180deg, #12121a 0%, #0a0a0f 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.06);
  display: flex;
  flex-direction: column;
  padding: 20px 12px;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  margin-bottom: 32px;
  cursor: pointer;
}

.logo-icon {
  width: 36px;
  height: 36px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon svg {
  width: 20px;
  height: 20px;
  color: white;
}

.logo-text {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
}

.nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 10px;
  color: #71717a;
  text-decoration: none;
  transition: all 0.2s;
  cursor: pointer;
}

.nav-item svg {
  width: 20px;
  height: 20px;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.05);
  color: #e5e5e5;
}

.nav-item.active {
  background: rgba(99, 102, 241, 0.15);
  color: #818cf8;
}

.sidebar-footer {
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  cursor: pointer;
}

.login-prompt {
  color: #52525b;
}

.avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, #3b82f6 0%, #8b5cf6 100%);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: white;
}

.user-details {
  flex: 1;
  overflow: hidden;
}

.username {
  display: block;
  color: #e5e5e5;
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
}

.api-key {
  font-size: 11px;
  color: #818cf8;
  cursor: pointer;
  padding: 2px 6px;
  background: rgba(99, 102, 241, 0.1);
  border-radius: 4px;
  transition: background 0.2s;
}

.api-key:hover {
  background: rgba(99, 102, 241, 0.2);
}

/* 主内容区 */
.main {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.header {
  padding: 32px 40px 24px;
  background: rgba(255, 255, 255, 0.01);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.header h1 {
  font-size: 28px;
  font-weight: 600;
  color: #fff;
  margin-bottom: 8px;
}

.subtitle {
  color: #71717a;
  font-size: 15px;
}

.content {
  flex: 1;
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  padding: 24px 40px;
  overflow-y: auto;
}

/* 卡片 */
.card {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 16px;
  padding: 24px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.card-header h2 {
  font-size: 16px;
  font-weight: 500;
  color: #e5e5e5;
}

.char-count {
  font-size: 13px;
  color: #52525b;
}

.history-actions {
  display: flex;
  gap: 8px;
}

.clear-btn, .refresh-btn {
  padding: 6px 12px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 6px;
  color: #a1a1aa;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s;
}

.clear-btn:hover:not(:disabled), .refresh-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e5e5;
}

.clear-btn:disabled, .refresh-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* 输入区域 */
.input-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

textarea {
  width: 100%;
  height: 160px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  padding: 16px;
  font-size: 15px;
  color: #e5e5e5;
  resize: none;
  outline: none;
  transition: border-color 0.2s;
  font-family: inherit;
  line-height: 1.6;
}

textarea::placeholder {
  color: #52525b;
}

textarea:focus {
  border-color: rgba(99, 102, 241, 0.5);
}

.settings-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.setting-item label {
  display: block;
  font-size: 13px;
  color: #71717a;
  margin-bottom: 8px;
}

select {
  width: 100%;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 10px;
  padding: 12px 16px;
  font-size: 14px;
  color: #e5e5e5;
  outline: none;
  cursor: pointer;
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='%2371717a'%3E%3Cpath d='M7 10l5 5 5-5z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
}

.select:focus {
  border-color: rgba(99, 102, 241, 0.5);
}

.slider-container {
  display: flex;
  align-items: center;
  gap: 12px;
}

input[type="range"] {
  flex: 1;
  height: 4px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 2px;
  outline: none;
  -webkit-appearance: none;
}

input[type="range"]::-webkit-slider-thumb {
  -webkit-appearance: none;
  width: 16px;
  height: 16px;
  background: #6366f1;
  border-radius: 50%;
  cursor: pointer;
}

.slider-value {
  font-size: 14px;
  color: #a1a1aa;
  min-width: 40px;
}

.generate-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 16px 32px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  border: none;
  border-radius: 12px;
  font-size: 16px;
  font-weight: 500;
  color: white;
  cursor: pointer;
  transition: all 0.2s;
}

.generate-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.3);
}

.generate-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.generate-btn svg {
  width: 20px;
  height: 20px;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

.error-msg {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  border-radius: 10px;
  color: #f87171;
  font-size: 14px;
}

.error-msg svg {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

/* 输出区域 */
.output-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.output-card {
  flex: 1;
  min-height: 300px;
}

.output-content {
  height: calc(100% - 40px);
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: rgba(0, 0, 0, 0.2);
  border: 1px dashed rgba(255, 255, 255, 0.1);
}

.output-content.has-audio {
  border-style: solid;
  border-color: rgba(99, 102, 241, 0.3);
}

.placeholder {
  text-align: center;
  color: #52525b;
}

.placeholder svg {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.placeholder p {
  font-size: 15px;
  margin-bottom: 4px;
}

.placeholder span {
  font-size: 13px;
  color: #3f3f46;
}

.empty-state {
  text-align: center;
  padding: 40px 20px;
  color: #52525b;
}

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-state p {
  font-size: 15px;
  margin-bottom: 4px;
}

.empty-state span {
  font-size: 13px;
  color: #3f3f46;
}

.audio-player {
  width: 100%;
  padding: 24px;
  text-align: center;
}

.waveform {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 3px;
  height: 60px;
  margin-bottom: 20px;
}

.wave-bar {
  width: 4px;
  height: 20px;
  background: linear-gradient(180deg, #6366f1 0%, #8b5cf6 100%);
  border-radius: 2px;
  animation: wave 1s ease-in-out infinite;
}

@keyframes wave {
  0%, 100% { height: 20px; }
  50% { height: 50px; }
}

audio {
  width: 100%;
  margin-bottom: 16px;
  border-radius: 10px;
}

.audio-actions {
  display: flex;
  justify-content: center;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 20px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 8px;
  color: #a1a1aa;
  font-size: 14px;
  cursor: pointer;
  transition: all 0.2s;
}

.action-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  color: #e5e5e5;
}

.action-btn svg {
  width: 18px;
  height: 18px;
}

/* 历史记录 */
.history-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.history-item {
  padding: 12px 16px;
  background: rgba(0, 0, 0, 0.2);
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.history-text {
  font-size: 14px;
  color: #e5e5e5;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-meta {
  font-size: 12px;
  color: #52525b;
}

/* 响应式 */
@media (max-width: 1024px) {
  .sidebar {
    display: none;
  }

  .content {
    grid-template-columns: 1fr;
    padding: 20px;
  }
}

@media (max-width: 640px) {
  .header {
    padding: 20px;
  }

  .settings-grid {
    grid-template-columns: 1fr;
  }
}
</style>
