<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { LOGIN_URL } from '@/consts'
import type { User } from '@/types'

const _user = ref<User>()

onMounted(() => {
  fetch(LOGIN_URL, {
    credentials: 'include'
  })
    .then(async (res) => {
      _user.value = (await res.json()) as User
    })
    .catch((error) => {
      console.error(error)
    })
})
</script>

<template>
  <div>
    {{ _user }}
  </div>
</template>

<style></style>
