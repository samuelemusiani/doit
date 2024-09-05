<script setup lang="ts">
import type { Options, Todo } from '@/types'
import { ref, inject } from 'vue'
import { useFocus } from '@vueuse/core'

import ColorPicker from '@/components/ColorPicker.vue'
import { onMounted, onBeforeUnmount } from 'vue'

const $emits = defineEmits<{
  (e: 'close'): void
  (e: 'addTodo', note: Todo): void
}>()

const __title_input = ref<HTMLElement | null>(null)
useFocus(__title_input, { initialValue: true })

const _todo = ref<Todo>({} as Todo)
_todo.value.StateID = 1
_todo.value.PriorityID = 1
_todo.value.ColorID = 1
_todo.value.Expiration = {
  DoesExpire: false,
  Date: new Date().toISOString()
}

const _date = ref<string>('')
const _time = ref<string>('')

const _todo_options = inject('todoOptions') as Options

function addTodo() {
  if (_todo.value.Title === '') {
    return
  }
  if (_todo.value.Expiration.DoesExpire) {
    if (_date.value == null) {
      return
    }
    if (_time.value == null) {
      return
    }
  }

  let d: Date = new Date()
  if (_todo.value.Expiration.DoesExpire) {
    d = new Date(_date.value + ' ' + _time.value)
  }

  _todo.value.Expiration.Date = d.toISOString()

  $emits('addTodo', _todo.value)
}

function close() {
  $emits('close')
}

function keyboardListener(event: KeyboardEvent) {
  if (event.key == 'Escape') {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keyup', keyboardListener)
})
onBeforeUnmount(() => {
  window.removeEventListener('keyup', keyboardListener)
})
</script>

<template>
  <div class="grid bg-gray-200 bg-opacity-30" @click="close()">
    <div
      class="fixed place-self-center rounded-lg bg-white shadow-2xl md:min-w-[25rem]"
      @click.stop=""
    >
      <header class="rounded-t-lg bg-gray-300">
        <h1 class="p-4 text-center text-xl font-bold">Add note</h1>
      </header>

      <body class="">
        <form class="flex flex-col p-5" @submit.prevent="">
          <label>Title</label>
          <input
            type="text"
            class="rounded border p-2 focus:border-black focus:outline-none"
            required
            v-model="_todo.Title"
            ref="__title_input"
          />

          <label class="mt-5">Description</label>
          <textarea
            type="text"
            class="rounded border p-2 focus:border-black focus:outline-none"
            v-model="_todo.Description"
          ></textarea>

          <div class="mt-5 flex flex-col justify-between md:flex-row md:gap-10">
            <label class="">State:</label>
            <div class="w-full">
              <select
                v-model="_todo.StateID"
                class="w-full rounded border border-black bg-white p-2"
              >
                <option v-for="state in _todo_options.States" :key="state.ID" :value="state.ID">
                  {{ state.State }}
                </option>
              </select>
            </div>
          </div>

          <div class="mt-5 flex flex-col justify-between md:flex-row md:gap-10">
            <label class="">Priority:</label>
            <div class="w-full">
              <select
                v-model="_todo.PriorityID"
                class="w-full rounded border border-black bg-white p-2"
              >
                <option
                  v-for="priority in _todo_options.Priorities"
                  :key="priority.ID"
                  :value="priority.ID"
                >
                  {{ priority.Priority }}
                </option>
              </select>
            </div>
          </div>

          <label class="mt-5">Color:</label>
          <ColorPicker :colors="_todo_options.Colors" v-model="_todo.ColorID" />

          <label class="mt-5">Expiration:</label>
          <div>
            <input
              id="expire_checkbox"
              type="checkbox"
              class="m-2 h-4 w-4 rounded"
              v-model="_todo.ColorID"
            />
            <label for="expire_checkbox">Does Expire</label>
          </div>
          <div class="flex flex-col justify-between gap-2 p-2 md:flex-row">
            <input
              type="date"
              :disabled="!_todo.Expiration.DoesExpire"
              class="disabled:opacity-20"
              :required="_todo.Expiration.DoesExpire"
              v-model="_date"
              placeholder="mm-dd-yyyy"
            />
            <input
              type="time"
              :disabled="!_todo.Expiration.DoesExpire"
              class="disabled:opacity-20"
              :required="_todo.Expiration.DoesExpire"
              v-model="_time"
            />
          </div>

          <div class="mt-5 flex justify-around">
            <button @click="close()" class="rounded bg-red-200 p-2 hover:bg-red-400">Close</button>
            <button @click="addTodo()" class="rounded bg-blue-200 p-2 hover:bg-blue-400">
              Add
            </button>
          </div>
        </form>
      </body>
    </div>
  </div>
</template>

<style></style>
