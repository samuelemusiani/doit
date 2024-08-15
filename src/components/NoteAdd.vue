<script setup lang="ts">
import type { Note } from '@/types'
import { ref } from 'vue'

const $emits = defineEmits<{
  (e: 'close'): void
  (e: 'addNote', note: Note): void
}>()

const _title = ref<string>('')
const _description = ref<string>('')

function close() {
  $emits('close')
}

function addNote() {
  let n: Note = { ID: 0, Title: _title.value, Description: _description.value }

  $emits('addNote', n)
}
</script>

<template>
  <div class="w-[20rem] rounded border">
    <header class="flex justify-between bg-gray-200">
      <h1 class="flex-auto p-5 text-center font-bold">Add a Note</h1>
    </header>

    <body class="">
      <form class="flex flex-col p-5" @submit.prevent="">
        <label>Title</label>
        <input type="text" class="rounded border p-2" required v-model="_title" />

        <label class="mt-5">Description</label>
        <textarea type="text" class="rounded border p-2" v-model="_description"></textarea>

        <div class="mt-5 flex justify-around">
          <button @click="close()" class="rounded bg-red-200 p-2 hover:bg-red-400">Close</button>
          <button @click="addNote()" class="rounded bg-blue-200 p-2 hover:bg-blue-400">Add</button>
        </div>
      </form>
    </body>
  </div>
</template>

<style></style>
