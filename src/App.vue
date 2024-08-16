<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { onMounted, ref, provide, computed } from 'vue'
import { COLORS_URL, LOGIN_ENDPOINT, PRIORITIES_URL, STATES_URL } from './consts'
import type { Options } from './types'

const $route = useRoute()

const path = computed(() => {
  return $route.path
})

const _options = {} as Options
provide('todoOptions', _options)

onMounted(() => {
  fetch(STATES_URL)
    .then(async (res) => {
      _options.States = await res.json()
    })
    .catch((error) => {
      console.error(error)
    })

  fetch(PRIORITIES_URL)
    .then(async (res) => {
      _options.Priorities = await res.json()
    })
    .catch((error) => {
      console.error(error)
    })

  fetch(COLORS_URL)
    .then(async (res) => {
      _options.Colors = await res.json()
    })
    .catch((error) => {
      console.error(error)
    })
})
</script>

<template>
  <NavBar v-if="path != LOGIN_ENDPOINT" />
  <RouterView />
</template>

<style></style>
