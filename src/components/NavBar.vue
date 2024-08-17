<script setup lang="ts">
import { LOGIN_URL } from '@/consts'
import { useRouter } from 'vue-router'

const $router = useRouter()

function logout() {
  fetch(LOGIN_URL, {
    method: 'DELETE',
    credentials: 'include'
  }).catch((error) => console.log(error))
}
</script>

<template>
  <div class="flex justify-between bg-gray-800">
    <RouterLink :to="{ name: 'home' }" class="grid p-2 text-gray-200 hover:bg-gray-700">
      <h1 class="text-4xl font-bold text-gray-200">DOIT</h1>
    </RouterLink>
    <RouterLink :to="{ name: 'profile' }" class="grid p-2 text-gray-200 hover:bg-gray-700">
      <span class="place-self-center"> Profile </span>
    </RouterLink>
    <RouterLink
      :to="{ name: 'login' }"
      class="grid p-2 text-gray-200 hover:bg-gray-700"
      @click="logout()"
    >
      <span class="place-self-center"> Logout </span>
    </RouterLink>
  </div>
</template>

<style></style>
