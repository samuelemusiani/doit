<script setup lang="ts">
import { LOGIN_URL } from '@/consts'
import { useRouter } from 'vue-router'

const $router = useRouter()

function logout() {
  fetch(LOGIN_URL, {
    method: 'DELETE',
    credentials: 'include'
  }).catch((error) => console.log(error))

  $router.push('login')
}
</script>

<template>
  <div class="flex justify-between bg-gray-800">
    <h1 class="p-2 text-4xl font-bold text-gray-200">DOIT</h1>
    <button class="p-2 text-gray-200 hover:bg-gray-700" @click="logout()">Logout</button>
  </div>
</template>

<style></style>
