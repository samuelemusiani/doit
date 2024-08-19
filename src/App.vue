<script setup lang="ts">
import { RouterView, useRoute } from 'vue-router'
import NavBar from '@/components/NavBar.vue'
import { onMounted, provide } from 'vue'
import type { Options } from './types'
import { getColors, getPriorities, getStates } from './lib/api'

const $route = useRoute()

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
  <NavBar v-if="!$route.meta.hide_navbar" />
  <RouterView />
</template>

<style></style>
