<template>
  <div class="container">
    <h1>MiMo 语音合成</h1>

    <label for="text">输入文本</label>
    <textarea
      id="text"
      v-model="text"
      placeholder="请输入要转换的文本..."
      maxlength="1000"
    ></textarea>

    <div class="row">
      <div>
        <label for="style">风格</label>
        <select id="style" v-model="selectedStyle">
          <option value="">默认</option>
          <option value="温柔">温柔</option>
          <option value="悄悄话">悄悄话</option>
          <option value="东北话 开心">东北话 开心</option>
          <option value="开心">开心</option>
          <option value="严肃">严肃</option>
        </select>
      </div>
    </div>

    <button class="primary" :disabled="loading" @click="convert">
      {{ loading ? '生成中...' : '生成语音' }}
    </button>

    <div class="result" v-if="audioUrl">
      <label>播放</label>
      <audio ref="player" :src="audioUrl" controls></audio>
    </div>

    <div class="error" v-if="error">{{ error }}</div>
    <div class="status" v-if="status">{{ status }}</div>
  </div>
</template>

<script setup>
import { ref, nextTick } from 'vue'

const text = ref('你好，欢迎使用 MiMo 语音合成服务。')
const selectedStyle = ref('')
const loading = ref(false)
const audioUrl = ref('')
const error = ref('')
const status = ref('')
const player = ref(null)

async function convert() {
  const trimmed = text.value.trim()
  error.value = ''
  status.value = ''

  if (!trimmed) {
    error.value = '请输入文本'
    return
  }
  if (trimmed.length > 1000) {
    error.value = '文本不能超过1000字'
    return
  }

  if (audioUrl.value) {
    URL.revokeObjectURL(audioUrl.value)
    audioUrl.value = ''
  }

  loading.value = true
  status.value = '正在调用 API，请稍候...'

  try {
    const res = await fetch('/api/v1/tts/convert', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ text: trimmed, style: selectedStyle.value })
    })

    if (!res.ok) {
      const data = await res.json()
      throw new Error(data.message || '请求失败')
    }

    const blob = await res.blob()
    audioUrl.value = URL.createObjectURL(blob)
    status.value = '生成成功，点击播放'

    await nextTick()
    player.value?.play()
  } catch (e) {
    error.value = e.message
    status.value = ''
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
* { box-sizing: border-box; margin: 0; padding: 0; }
.container {
  background: #fff;
  border-radius: 16px;
  padding: 32px;
  width: 480px;
  box-shadow: 0 4px 24px rgba(0,0,0,0.08);
}
h1 { font-size: 20px; margin-bottom: 24px; color: #1d1d1f; }
label { display: block; font-size: 14px; color: #6e6e73; margin-bottom: 6px; }
textarea {
  width: 100%; height: 100px;
  border: 1px solid #d2d2d7; border-radius: 10px;
  padding: 12px; font-size: 15px;
  resize: vertical; outline: none;
  transition: border-color 0.2s;
  font-family: inherit;
}
textarea:focus { border-color: #0071e3; }
.row { display: flex; gap: 12px; margin-top: 16px; }
.row > div { flex: 1; }
select {
  width: 100%;
  border: 1px solid #d2d2d7; border-radius: 10px;
  padding: 10px 12px; font-size: 14px;
  outline: none; background: #fff;
}
select:focus { border-color: #0071e3; }
button {
  margin-top: 20px; width: 100%;
  padding: 12px; border: none; border-radius: 10px;
  font-size: 16px; font-weight: 600;
  cursor: pointer; transition: background 0.2s;
}
button.primary { background: #0071e3; color: #fff; }
button.primary:hover { background: #0077ed; }
button.primary:disabled { background: #a1a1a6; cursor: not-allowed; }
.result { margin-top: 20px; }
audio { width: 100%; margin-top: 8px; }
.error { color: #ff3b30; font-size: 14px; margin-top: 12px; }
.status { color: #6e6e73; font-size: 13px; margin-top: 12px; text-align: center; }
</style>
