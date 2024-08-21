<script setup lang="ts">
import UserPorfile from '@/components/User/UserProfile.vue'
import { onMounted, ref } from 'vue'
import type { User } from '@/types'
import { useRoute, useRouter } from 'vue-router'
import { getUser } from '@/lib/api'
const $route = useRoute()
const $router = useRouter()

const _user = ref<User>({} as User)

onMounted(() => {
  let id: number = parseInt($route.params.id as string, 10)
  getUser(id).then((u) => (_user.value = u))
})
</script>

<template>
  <div class="flex justify-center p-2">
    <UserPorfile :user="_user" v-if="_user.ID" />
  </div>
</template>

<style></style>
