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
  <div class="ml-5 mt-5">
    <button class="rounded border p-2 hover:bg-gray-100" @click="$router.go(-1)">Go back</button>
  </div>
  <div class="flex justify-center">
    <div class="w-1/3">
      <UserPorfile :user="_user" v-if="_user.ID" />
    </div>
  </div>
</template>

<style></style>
