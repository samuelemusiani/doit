<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NOTES_URL } from '@/consts'
const notes = ref()

onMounted(async () => {
  console.log('mounted')
  try {
    const response = await fetch(NOTES_URL, {
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('Could not fetch notes')
    }
    notes.value = await response.json()
  } catch (error) {
    console.error(error)
  }
})
</script>

<template>
  <h1>Notes:</h1>
  <div>
    {{ notes }}
  </div>
</template>
