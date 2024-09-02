<script setup lang="ts">
import NotesList from '@/components/NotesList.vue'
import NoteAdd from '@/components/NoteAdd.vue'
import TodoStats from '@/components/TodoStats.vue'
import type { Options, Todo } from '@/types'
import { onMounted, ref, computed, inject } from 'vue'
import { addTodo, deleteTodo, fetchNotes, updateTodo } from '@/lib/api'

const _notes = ref<Todo[]>([])
const __addNote = ref<boolean>(false)

function newNote() {
  __addNote.value = true
}

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

onMounted(() => {
  fetchNotes().then((notes) => (_notes.value = notes))
})

// Responsive
const __show_mobile = ref<boolean>(false)
</script>

<template>
  <div class="p-2">
    <div class="flex flex-col-reverse md:flex-row md:justify-center">
      <NotesList
        :notes="_actual_todos"
        @updateTodo="_updateTodo"
        @deleteTodo="_deleteTodo"
        class="md:min-w-[40rem] md:max-w-[60rem]"
      />
      <div class="flex gap-5 md:ml-5 md:flex-col">
        <TodoStats class="hidden md:block" :todos="_notes" @selected="filterTodos" />
        <button
          class="rounded bg-orange-200 p-2 shadow md:hidden"
          @click="__show_mobile = !__show_mobile"
        >
          Filter
        </button>
        <TodoStats
          :class="{ hidden: !__show_mobile }"
          class="w-full md:hidden"
          :todos="_notes"
          @selected="filterTodos"
        />

        <input
          :class="{ hidden: __show_mobile }"
          type="text"
          v-model="_filter_search"
          class="w-full rounded border p-2 outline-none md:block"
          placeholder="search..."
        />
        <button :class="{ hidden: __show_mobile }" class="rounded border p-2 md:hidden">
          Search
        </button>
      </div>
    </div>
    <Transition>
      <NoteAdd
        @addTodo="_addNote"
        @close="__addNote = false"
        v-if="__addNote"
        class="absolute left-0 top-0 z-10 h-full w-full"
      />
    </Transition>
    <button
      class="fixed bottom-5 right-5 rounded bg-sky-200 p-5 shadow hover:bg-sky-400"
      @click="newNote"
    >
      Add Todo
    </button>
  </div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.1s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
