<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'

import { onMounted, ref } from 'vue'
import { NOTES_URL } from '@/consts'
import type { Note } from '@/types.ts'

const notes = ref<Note[]>([])

onMounted(async () => {
  try {
    const response = await fetch(NOTES_URL, {
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('Could not fetch notes')
    }
    notes.value = (await response.json()) as Note[]
  } catch (error) {
    console.error(error)
  }
})
</script>

<template>
  <NotesList :notes="notes" />
</template>

<style></style>
