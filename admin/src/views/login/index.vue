<template>
  <div class="login-page">
    <div class="login-card">
      <div class="login-header">
        <h1 class="login-title">T-LINE</h1>
        <p class="login-subtitle">Sports Club 管理后台</p>
      </div>

      <a-form :model="form" layout="vertical" @submit="handleLogin">
        <a-form-item field="phone" label="手机号" :rules="[{ required: true, message: '请输入手机号' }]">
          <a-input v-model="form.phone" placeholder="请输入管理员手机号" size="large" />
        </a-form-item>

        <a-form-item field="password" label="密码" :rules="[{ required: true, message: '请输入密码' }]">
          <a-input-password v-model="form.password" placeholder="请输入密码" size="large" />
        </a-form-item>

        <a-form-item>
          <a-button type="primary" html-type="submit" long size="large" :loading="loading"
            style="background: #2255CC; height: 48px; font-size: 16px;">
            登 录
          </a-button>
        </a-form-item>
      </a-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { post } from '@/api/request'
import { Message } from '@arco-design/web-vue'

const router = useRouter()
const loading = ref(false)

const form = reactive({
  phone: '',
  password: '',
})

async function handleLogin() {
  if (!form.phone || !form.password) return

  loading.value = true
  try {
    const data: any = await post('/auth/password-login', {
      phone: form.phone,
      password: form.password,
    })
    localStorage.setItem('access_token', data.access_token)
    localStorage.setItem('refresh_token', data.refresh_token)
    Message.success('登录成功')
    router.push('/dashboard')
  } catch (err: any) {
    Message.error(err?.message || '登录失败')
  } finally {
    loading.value = false
  }
}
</script>

<style lang="scss" scoped>
.login-page {
  min-height: 100vh; display: flex; align-items: center; justify-content: center;
  background: linear-gradient(135deg, #1a1f36 0%, #2255CC 100%);
}

.login-card {
  width: 400px; background: #fff; border-radius: 12px;
  padding: 48px 40px; box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2);
}

.login-header { text-align: center; margin-bottom: 40px; }

.login-title {
  font-size: 36px; font-weight: 800; color: #2255CC;
  letter-spacing: 4px; margin: 0;
}

.login-subtitle {
  font-size: 14px; color: #999; margin-top: 8px;
}
</style>
