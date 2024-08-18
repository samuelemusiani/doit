<script setup lang="ts">
import { onMounted, ref } from 'vue'
import type { User } from '@/types'
import { getUsers } from '@/lib/api'

const _users = ref<User[]>([])

onMounted(() => {
  getUsers().then((users) => (_users.value = users))
})
</script>

<template>
  <div class="ml-5 mt-5">
    <button class="rounded border p-2 hover:bg-gray-100" @click="$router.go(-1)">Go back</button>
  </div>
  <div class="flex justify-center">
    <div class="w-1/3 rounded border p-5">
      <div class="mt-2 flex justify-between">
        <h1 class="text-center font-bold">Users list</h1>
        <div>Total: {{ _users.length }}</div>
      </div>
      <div class="mt-5 flex flex-col gap-2">
        <template v-for="user in _users" :key="user.ID">
          <RouterLink :to="{ name: 'user_details', params: { id: user.ID } }">
            <div class="rounded border p-2 hover:bg-gray-100">
              {{ user.Username }}
            </div>
          </RouterLink>
        </template>
      </div>
    </div>
  </div>
</template>

<style></style>
