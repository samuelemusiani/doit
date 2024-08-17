<script setup lang="ts">
import type { Options, Todo } from '@/types'
import { ref, inject } from 'vue'
import { useDraggable, useFocus } from '@vueuse/core'

import ColorPicker from '@/components/ColorPicker.vue'

const $emits = defineEmits<{
  (e: 'close'): void
  (e: 'addNote', note: Todo): void
}>()

const __top_div = ref<HTMLElement | null>(null)
const __header = ref<HTMLElement | null>(null)
const { style } = useDraggable(__top_div, {
  initialValue: {
    x: document.documentElement.clientWidth / 4,
    y: document.documentElement.clientHeight / 4
  },
  handle: __header
})

const __title_input = ref<HTMLElement | null>(null)
useFocus(__title_input, { initialValue: true })

const _title = ref<string>('')
const _description = ref<string>('')
const _state = ref<number>(1)
const _priority = ref<number>(1)
const _color = ref<number>(1)
const _does_expire = ref<boolean>(false)

const _date = ref<string>('')
const _time = ref<string>('')

const _todo_options = inject('todoOptions') as Options

function close() {
  $emits('close')
}

function addNote() {
  if (_title.value === '') {
    return
  }
  if (_does_expire.value) {
    if (_date.value == null) {
      return
    }
    if (_time.value == null) {
      return
    }
  }

  let d: Date
  if (_does_expire.value) {
    d = new Date(_date.value)
    d.setHours(parseInt(_time.value.split(':')[0]))
    d.setMinutes(parseInt(_time.value.split(':')[1]))
  } else {
    d = new Date()
  }

  let n: Todo = {
    ID: 0,
    Title: _title.value,
    Description: _description.value,
    StateID: _state.value,
    PriorityID: _priority.value,
    ColorID: _color.value,
    Expiration: {
      DoesExpire: _does_expire.value,
      Date: d.toISOString()
    }
  }

  $emits('addNote', n)
}
</script>

<template>
  <div
    class="fixed w-[25rem] rounded border border-black bg-white drop-shadow-lg"
    :style="style"
    ref="__top_div"
  >
    <header class="flex justify-between rounded bg-gray-200" ref="__header">
      <h1 class="flex-auto p-5 text-center font-bold">Add a Note</h1>
    </header>

    <body class="">
      <form class="flex flex-col p-5" @submit.prevent="">
        <label>Title</label>
        <input
          type="text"
          class="rounded border p-2 focus:border-black focus:outline-none"
          required
          v-model="_title"
          ref="__title_input"
        />

        <label class="mt-5">Description</label>
        <textarea
          type="text"
          class="rounded border p-2 focus:border-black focus:outline-none"
          v-model="_description"
        ></textarea>

        <div class="mt-5 flex justify-between gap-10">
          <label class="">State:</label>
          <div class="w-full">
            <select v-model="_state" class="w-full rounded border border-black bg-white p-2">
              <option v-for="state in _todo_options.States" :key="state.ID" :value="state.ID">
                {{ state.State }}
              </option>
            </select>
          </div>
        </div>

        <div class="mt-5 flex justify-between gap-10">
          <label class="">Priority:</label>
          <div class="w-full">
            <select v-model="_priority" class="w-full rounded border border-black bg-white p-2">
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
        <ColorPicker :colors="_todo_options.Colors" v-model="_color" />

        <label class="mt-5">Expiration:</label>
        <div>
          <input
            id="expire_checkbox"
            type="checkbox"
            class="m-2 h-4 w-4 rounded"
            v-model="_does_expire"
          />
          <label for="expire_checkbox">Does Expire</label>
        </div>
        <div class="flex justify-between p-2">
          <input
            type="date"
            :disabled="!_does_expire"
            class="disabled:opacity-20"
            :required="_does_expire"
            v-model="_date"
            placeholder="mm-dd-yyyy"
          />
          <input
            type="time"
            :disabled="!_does_expire"
            class="disabled:opacity-20"
            :required="_does_expire"
            v-model="_time"
          />
        </div>
        {{ _date }}
        {{ _time }}

        <div class="mt-5 flex justify-around">
          <button @click="close()" class="rounded bg-red-200 p-2 hover:bg-red-400">Close</button>
          <button @click="addNote()" class="rounded bg-blue-200 p-2 hover:bg-blue-400">Add</button>
        </div>
      </form>
    </body>
  </div>
</template>

<style></style>
