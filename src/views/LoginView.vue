<script setup lang="ts">
import router from '@/router'
import { ref } from 'vue'
import { useFocus } from '@vueuse/core'
import { login } from '@/lib/api'

const _username = ref('')
const _password = ref('')

const _userinput = ref<HTMLElement | null>(null)

const _errorText = ref('')

useFocus(_userinput, { initialValue: true })

async function _login() {
  login(_username.value, _password.value)
    .then(() => {
      router.push({ name: 'home' })
    })
    .catch((err) => {
      _errorText.value = "Something went wrong. Can't login"
      console.error(err)
    })

  _password.value = ''
}
</script>

<template>
  <!-- Desktop -->
  <div
    class="absolute top-0 -z-10 hidden h-full w-full items-center justify-between font-mono text-[20rem] text-gray-100 opacity-60 md:flex lg:flex lg:text-[24rem] xl:text-[32rem] 2xl:text-[40rem]"
  >
    <div>DO</div>
    <div>IT</div>
  </div>

  <!-- Mobile -->
  <div
    class="absolute top-0 -z-10 flex h-screen w-screen flex-col items-center justify-between font-mono text-[16rem] text-gray-100 opacity-60 md:hidden"
  >
    <div>DO</div>
    <div>IT</div>
  </div>
  <div class="grid h-full">
    <div class="place-self-center rounded-lg border bg-white p-5 shadow-lg md:min-w-96">
      <form class="mt-5" @submit.prevent="_login()">
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
