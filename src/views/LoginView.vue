<script setup lang="ts">
import { LOGIN_URL } from '@/consts'
import router from '@/router'
import { ref } from 'vue'

const username = ref('')
const password = ref('')

const errorText = ref('')

async function login() {
  try {
    const response = await fetch(LOGIN_URL, {
      method: 'POST',
      body: JSON.stringify({
        username: username.value,
        password: password.value
      }),
      credentials: 'include' // Used to set coockies; DOTO Check if this should be in production
    })
    password.value = '' // Reset password

    if (response.ok) {
      router.push('/') // Redirect to home after successful login
    } else {
      errorText.value = await response.text()
      switch (response.status) {
        case 400: {
          throw new Error(errorText.value)
        }
        case 404: {
          throw new Error(errorText.value)
        }
        case 500: {
          throw new Error(errorText.value)
        }
        default:
          break
      }
    }
  } catch (error) {
    console.error(error)
  }
}
</script>

<template>
  <h1>LOGIN PAGE</h1>
  <form @submit.prevent="login()">
    <div>
      <label for="username">Username:</label>
      <input id="username" type="text" v-model="username" />
    </div>

    <div>
      <label for="password">Password:</label>
      <input id="password" type="password" v-model="password" />
    </div>
    <button>Login</button>
    <div>
      {{ errorText }}
    </div>
  </form>
</template>

<style></style>
