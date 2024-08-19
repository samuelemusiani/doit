<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { User } from '@/types'
import { getCurrentUser, logout } from '@/lib/api'

const _user = ref<User>({} as User)

onMounted(() => {
  getCurrentUser()
    .then((u) => (_user.value = u))
    .catch((err: Error) => {
      console.error(err)
    })
})
</script>

<template>
  <div class="flex justify-between bg-gray-800">
    <RouterLink :to="{ name: 'home' }" class="grid p-2 text-gray-200 hover:bg-gray-700">
      <h1 class="text-4xl font-bold text-gray-200">DOIT</h1>
    </RouterLink>
    <div class="flex">
      <RouterLink
        :to="{ name: 'admin' }"
        class="grid p-2 text-gray-200 hover:bg-gray-700"
        v-if="_user.Admin"
      >
        <span class="place-self-center"> Admin </span>
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
  </div>
</template>

<style></style>
