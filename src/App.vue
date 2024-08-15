<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { computed } from 'vue'
import { LOGIN_ENDPOINT } from './consts'

const $route = useRoute()

const path = computed(() => {
  return $route.path
})
</script>

<template>
  <NavBar v-if="path != LOGIN_ENDPOINT" />
  <RouterView />
</template>

<style></style>
