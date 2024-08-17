<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'
import NoteAdd from '@/components/NoteAdd.vue'
import TodoStats from '@/components/TodoStats.vue'
import { onClickOutside } from '@vueuse/core'
import { onMounted, ref } from 'vue'
import { NOTES_URL } from '@/consts'
import type { Options, Todo } from '@/types'
import { computed } from 'vue'
import { inject } from 'vue'

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

const _filter = ref<number>(0)

const _todos_options = inject('todoOptions') as Options

const _actual_todos = computed(() => {
  if (_filter.value === 0) {
    // We assume that the 'done' state is the last one
    return _notes.value.filter((e) => e.StateID != _todos_options.States.length)
  } else {
    return _notes.value.filter((e) => e.StateID == _filter.value)
  }
})

function filterTodos(s: number) {
  _filter.value = s
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
    <div class="flex w-1/2 flex-col">
      <NotesList
        :notes="_actual_todos"
        @updateTodo="updateTodo"
        @deleteNote="deleteNote"
        class=""
      />
    </div>
    <div class="ml-5 mt-9">
      <TodoStats :todos="_notes" @selected="filterTodos" />
    </div>
  </div>
  <NoteAdd
    @addModifyNote="addNote"
    @close="__addNote = false"
    v-if="__addNote"
    ref="__addNote_ref"
  />
  <button
    class="fixed bottom-5 right-5 rounded bg-blue-200 p-5 hover:bg-blue-400"
    @click="__addNote = true"
  >
    Add Todo
  </button>
</template>

<style></style>
