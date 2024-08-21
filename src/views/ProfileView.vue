<script setup lang="ts">
import UserPorfile from '@/components/User/UserProfile.vue'
import { onMounted, ref } from 'vue'
import type { User } from '@/types'
import { getCurrentUser } from '@/lib/api'

const _user = ref<User>({} as User)

onMounted(() => {
  getCurrentUser().then((u) => (_user.value = u))
})
</script>

<template>
  <div class="flex justify-center">
    <div class="mt-5">
      <UserPorfile :user="_user" v-if="_user.ID" />
    </div>
  </div>
</template>

<style></style>
