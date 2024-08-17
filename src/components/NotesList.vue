<script setup lang="ts">
import type { PropType } from 'vue'
import type { Options, Todo } from '@/types'
import { inject } from 'vue'

const _todo_options = inject('todoOptions') as Options

const $props = defineProps({
  notes: {
    type: Array as PropType<Todo[]>,
    required: true
  }
})

const $emits = defineEmits<{
  (event: 'deleteNote', id: number): void
  (event: 'updateTodo', todo: Todo): void
}>()

function sortTodos(a: Todo, b: Todo) {
  let diff = b.PriorityID - a.PriorityID

  if (diff == 0) {
    let d1 = new Date(b.Expiration.Date)
    let d2 = new Date(a.Expiration.Date)

    diff = d1.getTime() - d2.getTime()
  }
  return b.PriorityID - a.PriorityID
}

function advanceStateTodo(todo: Todo) {
  todo.StateID = (todo.StateID + 1) % (_todo_options.States.length + 1)
  if (todo.StateID == 0) {
    todo.StateID = 1
  }

  $emits('updateTodo', todo)
}
</script>

<template>
  <div class="">
    <template v-for="note in $props.notes.sort(sortTodos)" :key="note.ID">
      <div class="m-2 flex justify-between rounded border border-black p-5">
        <div class="mr-5 grid items-center">
          <button
            class="w-24 rounded p-2 hover:opacity-90"
            @click="advanceStateTodo(note)"
            :style="{ 'background-color': _todo_options.Colors[note.ColorID - 1].Hex }"
          >
            {{ _todo_options.States[note.StateID - 1].State }}
          </button>
        </div>
        <div class="w-full">
          <h3 class="font-semibold">{{ note.Title }}</h3>
          <p>
            {{ note.Description }}
          </p>
          <div class="flex w-full justify-end">
            <span v-if="note.Expiration.DoesExpire">
              {{ new Date(note.Expiration.Date).toDateString() }}
              {{ new Date(note.Expiration.Date).toLocaleTimeString() }}
            </span>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style></style>
