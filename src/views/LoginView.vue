<script setup lang="ts">
import { LOGIN_URL } from '@/consts'
import router from '@/router'
import { ref } from 'vue'
import { useFocus } from '@vueuse/core'

const _username = ref('')
const _password = ref('')

const _userinput = ref<HTMLElement | null>(null)

const _errorText = ref('')

useFocus(_userinput, { initialValue: true })

async function login() {
  try {
    const response = await fetch(LOGIN_URL, {
      method: 'POST',
      body: JSON.stringify({
        username: _username.value,
        password: _password.value
      }),
      credentials: 'include' // Used to set coockies; DOTO Check if this should be in production
    })
    _password.value = '' // Reset password

    if (response.ok) {
      router.push('/') // Redirect to home after successful login
    } else {
      _errorText.value = await response.text()
      switch (response.status) {
        case 400: {
          throw new Error(_errorText.value)
        }
        case 404: {
          throw new Error(_errorText.value)
        }
        case 500: {
          throw new Error(_errorText.value)
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
  <div
    class="absolute top-0 -z-10 hidden h-full w-full items-center font-mono text-[20rem] text-gray-100 opacity-60 lg:flex lg:text-[24rem] xl:text-[32rem] 2xl:text-[40rem]"
  >
    <div class="flex w-full justify-between">
      <div>DO</div>
      <div>IT</div>
    </div>
  </div>
  <div class="grid h-full">
    <div class="w-96 place-self-center rounded-lg border bg-white p-5 shadow-lg">
      <form class="mt-5" @submit.prevent="login()">
        <div class="mt-5 w-full">
          <input
            id="username"
            type="text"
            v-model="_username"
            class="w-full rounded-lg border p-2 outline-none"
            required
            placeholder="user"
            ref="_userinput"
          />
        </div>

        <div class="mt-5 w-full">
          <input
            id="password"
            type="password"
            v-model="_password"
            class="w-full rounded-lg border p-2 outline-none"
            required
            placeholder="**********"
          />
        </div>
        <div class="mt-5 flex justify-center">
          <button class="rounded-lg border p-2 hover:bg-gray-200">Login</button>
        </div>
        <div class="mt-5 text-red-500" v-show="_errorText.length > 0">
          {{ _errorText }}
        </div>
      </form>
    </div>
  </div>
</template>

<style></style>
