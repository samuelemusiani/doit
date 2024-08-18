<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { onMounted, provide, computed } from 'vue'
import { LOGIN_ENDPOINT } from './consts'
import type { Options } from './types'
import { getColors, getPriorities, getStates } from './lib/api'

const $route = useRoute()

const path = computed(() => {
  return $route.path
})

const _options = {} as Options
provide('todoOptions', _options)

onMounted(() => {
  try {
    getStates().then((states) => (_options.States = states))
  } catch (err) {
    console.error(err)
  }

  try {
    getPriorities().then((priorities) => (_options.Priorities = priorities))
  } catch (err) {
    console.error(err)
  }

  try {
    getColors().then((colors) => (_options.Colors = colors))
  } catch (err) {
    console.error(err)
  }
})
</script>

<template>
  <NavBar v-if="path != LOGIN_ENDPOINT" />
  <RouterView />
</template>

<style></style>
