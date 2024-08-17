<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'
import NoteAdd from '@/components/NoteAdd.vue'

import { onClickOutside } from '@vueuse/core'

import { onMounted, ref } from 'vue'
import { NOTES_URL } from '@/consts'
import type { Todo } from '@/types'
import { onBeforeUnmount } from 'vue'

const _notes = ref<Todo[]>([])
const __addNote = ref<boolean>(false)

const __addNote_ref = ref<HTMLElement | null>(null)
onClickOutside(__addNote_ref, () => (__addNote.value = false))

async function fetchNotes() {
  try {
    const response = await fetch(NOTES_URL, {
      credentials: 'include'
    })
    if (!response.ok) {
      throw new Error('Could not fetch notes')
    }
    _notes.value = (await response.json()) as Todo[]
  } catch (error) {
    console.error(error)
  }
}

function updateTodo(todo: Todo) {
  fetch(NOTES_URL + '/' + todo.ID, {
    method: 'PUT',
    credentials: 'include',
    body: JSON.stringify(todo)
  })
    .then(() => {
      fetchNotes()
    })
    .catch((error) => {
      console.error(error)
    })
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

function addNote(note: Todo) {
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

// For now this does not work because if we use need to type 'a' in a nested
// component the function will fire anyway
//function keyboardListener(event: KeyboardEvent) {
//  if (event.key == 'a') {
//    __addNote.value = true
//  }
//}
//
onMounted(async () => {
  //  window.addEventListener('keypress', keyboardListener)
  fetchNotes()
})
//
//onBeforeUnmount(() => {
//  window.removeEventListener('keypress', keyboardListener)
//})
</script>

<template>
  <div class="flex justify-center">
    <div class="flex w-[70%] flex-col">
      <h1 class="p-2 text-xl font-bold">Todos:</h1>
      <NotesList :notes="_notes" @updateTodo="updateTodo" @deleteNote="deleteNote" class="" />
    </div>
  </div>
  <NoteAdd
    @addModifyNote="addNote"
    @close="__addNote = false"
    v-if="__addNote"
    ref="__addNote_ref"
  />
  <button
    class="absolute bottom-5 right-5 rounded bg-blue-200 p-5 hover:bg-blue-400"
    @click="__addNote = true"
  >
    Add Todo
  </button>
</template>

<style></style>
