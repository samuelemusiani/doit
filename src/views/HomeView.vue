<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'
import NoteAdd from '@/components/NoteAdd.vue'

import { onMounted, ref } from 'vue'
import { NOTES_URL } from '@/consts'
import type { Note } from '@/types.ts'

const _notes = ref<Note[]>([])
const __addNote = ref<boolean>(false)

async function fetchNotes() {
  try {
    const response = await fetch(NOTES_URL, {
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('Could not fetch notes')
    }
    _notes.value = (await response.json()) as Note[]
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

function addNote(note: Note) {
  fetch(NOTES_URL, {
    method: 'POST',
    credentials: 'include',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(note)
  })
    .then(() => {
      fetchNotes()
    })
    .catch((error) => {
      console.error(error)
    })

  __addNote.value = false
}

onMounted(async () => {
  fetchNotes()
})
</script>

<template>
  <NotesList :notes="_notes" @deleteNote="deleteNote" />
  <NoteAdd @addNote="addNote" @close="__addNote = false" v-if="__addNote" />
  <button
    class="absolute bottom-5 right-5 rounded bg-blue-200 p-5 hover:bg-blue-400"
    @click="__addNote = true"
  >
    Add note
  </button>
</template>

<style></style>
