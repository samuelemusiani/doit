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
  <div class="p-5">
    <template v-for="note in $props.notes" :key="note.ID">
      <div class="border-black border m-2 p-5 rounded">
        <div class="flex justify-between">
          <h3 class="font-semibold">{{ note.Title }}</h3>
          <span> ID: {{ note.ID }} </span>
        </div>
        <p>
          {{ note.Description }}
        </p>
      </div>
    </template>
  </div>
</template>

<style></style>
