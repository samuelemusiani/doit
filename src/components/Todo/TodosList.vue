<script setup lang="ts">
import TodoViewer from './TodoViewer.vue'
import type { PropType } from 'vue'
import type { Options, Todo } from '@/types'
import { inject, ref } from 'vue'
import { onClickOutside } from '@vueuse/core'

const _todo_options = inject('todoOptions') as Options

const $props = defineProps({
  notes: {
    type: Array as PropType<Todo[]>,
    required: true
  }
})

const $emits = defineEmits<{
  (event: 'deleteTodo', id: number): void
  (event: 'updateTodo', todo: Todo): void
}>()

const _todo_to_modify = ref<Todo>()
const _modify_todo = ref<boolean>()

function sortTodos(a: Todo, b: Todo) {
  let diff = b.PriorityID - a.PriorityID

  if (diff == 0) {
    let d1 = new Date(b.Expiration.Date)
    let d2 = new Date(a.Expiration.Date)

    diff = d1.getTime() - d2.getTime()
  }
  return b.PriorityID - a.PriorityID
}

function advanceStateTodo(todo: Todo, negative: boolean) {
  todo.StateID = todo.StateID + (negative ? -1 : 1)
  todo.StateID = Math.min(todo.StateID, _todo_options.States.length)
  todo.StateID = Math.max(todo.StateID, 1)

  $emits('updateTodo', todo)
}

function callModify(todo: Todo) {
  // deep copy
  _todo_to_modify.value = JSON.parse(JSON.stringify(todo))
  _modify_todo.value = true
}

function modifyTodo(todo: Todo) {
  $emits('updateTodo', todo)
  _modify_todo.value = false
}

function deleteTodo(todo: Todo) {
  $emits('deleteTodo', todo.ID)
  _modify_todo.value = false
}

function closeModify() {
  _modify_todo.value = false
}
</script>

<template>
  <div>
    <!--
    <TodoAdd :todo="_todo_to_modify" :modify="true" v-if="_modify_todo" @close="closeModify" @addModifyTodo="modifyTodo"
      @deleteTodo="deleteTodo" ref="__modify_ref" />
-->
    <TodoViewer
      :todo="_todo_to_modify"
      v-if="_modify_todo && _todo_to_modify"
      @close="closeModify"
      @modify="modifyTodo"
      @delete="deleteTodo"
      class="absolute right-0 top-0 h-full w-full"
    />
    <template v-for="note in $props.notes.sort(sortTodos)" :key="note.ID">
      <div
        class="my-2 flex flex-col gap-5 rounded border border-black p-5 hover:bg-gray-100 md:max-h-28 md:flex-row md:justify-between"
        @click="callModify(note)"
      >
        <div class="flex overflow-hidden">
          <div class="mr-5 grid items-center">
            <button
              class="w-24 rounded p-2 hover:saturate-150"
              @click.stop.exact="advanceStateTodo(note, false)"
              @click.stop.shift="advanceStateTodo(note, true)"
              :style="{ 'background-color': _todo_options.Colors[note.ColorID - 1].Hex }"
            >
              {{ _todo_options.States[note.StateID - 1].State }}
            </button>
          </div>
          <div class="overflow-auto">
            <h3 class="truncate font-semibold">{{ note.Title }}</h3>
            <p class="truncate">
              {{ note.Description }}
            </p>
          </div>
        </div>
        <div class="flex justify-between md:mt-0 md:flex-col">
          <div class="min-w-20 md:text-right">
            {{ _todo_options.Priorities[note.PriorityID - 1].Priority }}
          </div>
          <div v-if="note.Expiration.DoesExpire" class="min-w-40 text-right">
            {{ new Date(note.Expiration.Date).toDateString() }}
            {{ new Date(note.Expiration.Date).toLocaleTimeString() }}
          </div>
        </div>
      </div>
    </template>
    <div v-if="notes.length == 0" class="text-gray-800">
      Looks like theres is nothing to do... :)
    </div>
  </div>
</template>

<style></style>
