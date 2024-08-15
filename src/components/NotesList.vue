<script setup lang="ts">
import type { PropType } from 'vue'
import type { Todo } from '@/types'
import { NOTES_URL } from '@/consts'

const $props = defineProps({
  notes: {
    type: Array as PropType<Todo[]>,
    required: true
  }
})

const $emits = defineEmits<{
  (event: 'deleteNote', id: number): void
}>()

function deleteNote(id: number) {
  $emits('deleteNote', id)
}
</script>

<template>
  <div class="">
    <template v-for="note in $props.notes" :key="note.ID">
      <div class="m-2 flex justify-between rounded border border-black p-5">
        <div class="">
          <h3 class="font-semibold">{{ note.Title }}</h3>
          <p>
            {{ note.Description }}
          </p>
        </div>
        <div class="flex">
          <button class="rounded bg-red-200 p-2 hover:bg-red-400" @click="deleteNote(note.ID)">
            Delete
          </button>
        </div>
      </div>
    </template>
  </div>
</template>

<style></style>
