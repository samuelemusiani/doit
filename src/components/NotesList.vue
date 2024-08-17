<script setup lang="ts">
import NoteAdd from './NoteAdd.vue'
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
const __modify_ref = ref<HTMLElement | null>(null)
onClickOutside(__modify_ref, () => {
  closeModify()
})

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

function callModify(todo: Todo) {
  _todo_to_modify.value = todo
  _modify_todo.value = true
}

function modifyNote(todo: Todo) {
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
    <NoteAdd
      :todo="_todo_to_modify"
      :modify="true"
      v-if="_modify_todo"
      @close="closeModify"
      @addModifyNote="modifyNote"
      @deleteTodo="deleteTodo"
      ref="__modify_ref"
    />
    <template v-for="note in $props.notes.sort(sortTodos)" :key="note.ID">
      <div
        class="m-2 flex justify-between rounded border border-black p-5 hover:bg-gray-100"
        @click="callModify(note)"
      >
        <div class="mr-5 grid items-center">
          <button
            class="w-24 rounded p-2 hover:saturate-150"
            @click.stop="advanceStateTodo(note)"
            :style="{ 'background-color': _todo_options.Colors[note.ColorID - 1].Hex }"
          >
            {{ _todo_options.States[note.StateID - 1].State }}
          </button>
        </div>
        <div class="flex w-full justify-between">
          <div>
            <h3 class="font-semibold">{{ note.Title }}</h3>
            <p>
              {{ note.Description }}
            </p>
          </div>
          <div class="flex flex-col justify-between">
            <div class="text-right">
              {{ _todo_options.Priorities[note.PriorityID - 1].Priority }}
            </div>
            <div v-if="note.Expiration.DoesExpire">
              {{ new Date(note.Expiration.Date).toDateString() }}
              {{ new Date(note.Expiration.Date).toLocaleTimeString() }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<style></style>
