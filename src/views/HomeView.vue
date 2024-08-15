<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'

import { onMounted, ref } from 'vue'
import { NOTES_URL } from '@/consts'
import type { Note } from '@/types.ts'

const notes = ref<Note[]>([])

async function fetchNotes() {
  try {
    const response = await fetch(NOTES_URL, {
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('Could not fetch notes')
    }
    notes.value = (await response.json()) as Note[]
    console.log(notes.value)
  } catch (error) {
    console.error(error)
  }
}

function deleteNote(id: number) {
  fetch(NOTES_URL + '/' + id, {
    method: 'DELETE',
    credentials: 'include'
  })
    .then(() => {
      fetchNotes()
    })
    .catch((error) => {
      console.error(error)
    })
}

onMounted(async () => {
  fetchNotes()
})
</script>

<template>
  <NotesList :notes="notes" @deleteNote="deleteNote" />
</template>

<style></style>
