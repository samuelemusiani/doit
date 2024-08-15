<script setup lang="ts">
import type { PropType } from 'vue'
import type { Note } from '@/types'

const $props = defineProps({
  notes: {
    type: Array as PropType<Note[]>,
    required: true
  }
})
</script>

<template>
  <h1>Notes:</h1>
  <template v-for="note in $props.notes" :key="note.ID">
    <div>
      <h3>
        {{ note.Title }}
      </h3>
      <span>
        {{ note.ID }}
      </span>
      <p>
        {{ note.Description }}
      </p>
    </div>
  </template>
</template>

<style></style>
