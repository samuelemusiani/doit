<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'
import NoteAdd from '@/components/NoteAdd.vue'
import TodoStats from '@/components/TodoStats.vue'
import type { Options, Todo } from '@/types'
import { onClickOutside } from '@vueuse/core'
import { onMounted, ref, computed, inject } from 'vue'
import { addTodo, deleteTodo, fetchNotes, updateTodo } from '@/lib/api'

const _notes = ref<Todo[]>([])
const __addNote = ref<boolean>(false)

const __addNote_ref = ref<HTMLElement | null>(null)
onClickOutside(__addNote_ref, () => (__addNote.value = false))

function _updateTodo(todo: Todo) {
  updateTodo(todo)
    .then(() => {
      fetchNotes().then((notes) => (_notes.value = notes))
    })
    .catch((error) => {
      console.error(error)
    })
}

function _deleteTodo(id: number) {
  deleteTodo(id)
    .then(() => {
      fetchNotes().then((notes) => (_notes.value = notes))
    })
    .catch((error) => {
      console.error(error)
    })
}

function _addNote(todo: Todo) {
  addTodo(todo)
    .then(() => {
      fetchNotes().then((notes) => (_notes.value = notes))
    })
    .catch((error) => {
      console.error(error)
    })

  __addNote.value = false
}

const _filter_state = ref<number>(0)
const _filter_search = ref<string>('')

const _todos_options = inject('todoOptions') as Options

const _actual_todos = computed(() => {
  return _notes.value.filter((e) => {
    let c1: boolean
    let c2: boolean

    if (_filter_state.value === 0) {
      c1 = e.StateID != _todos_options.States.length
    } else {
      c1 = e.StateID == _filter_state.value
    }

    if (_filter_search) {
      c2 = e.Title.toLowerCase().includes(_filter_search.value)
      c2 = c2 || e.Description.toLowerCase().includes(_filter_search.value)
    } else {
      c2 = true
    }

    return c1 && c2
  })
})

function filterTodos(s: number) {
  _filter_state.value = s
}

// For now this does not work because if we use need to type 'a' in a nested
// component the function will fire anyway
//function keyboardListener(event: KeyboardEvent) {
//  if (event.key == 'a') {
//    __addNote.value = true
//  }
//}
//
onMounted(() => {
  //  window.addEventListener('keypress', keyboardListener)
  fetchNotes().then((notes) => (_notes.value = notes))
})
//
//onBeforeUnmount(() => {
//  window.removeEventListener('keypress', keyboardListener)
//})
</script>

<template>
  <div class="mt-2 flex justify-center">
    <div class="flex w-1/2 flex-col">
      <NotesList
        :notes="_actual_todos"
        @updateTodo="_updateTodo"
        @deleteTodo="_deleteTodo"
        class=""
      />
    </div>
    <div class="ml-5 mt-2">
      <div class="">
        <TodoStats :todos="_notes" @selected="filterTodos" />
      </div>
      <div class="mt-5">
        <input
          type="text"
          v-model="_filter_search"
          class="rounded border p-2 outline-none"
          placeholder="search..."
        />
      </div>
    </div>
  </div>
  <NoteAdd
    @addModifyNote="_addNote"
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
