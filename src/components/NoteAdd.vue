<script setup lang="ts">
import type { Options, Todo } from '@/types'
import { ref, inject } from 'vue'
import { useDraggable, useFocus } from '@vueuse/core'

import ColorPicker from '@/components/ColorPicker.vue'
import type { PropType } from 'vue'
import { onMounted } from 'vue'
import { onBeforeUnmount } from 'vue'

const $props = defineProps({
  todo: {
    type: Object as PropType<Todo>,
    required: false
  },
  modify: {
    type: Boolean as PropType<Boolean>,
    required: false
  }
})

const $emits = defineEmits<{
  (e: 'close'): void
  (e: 'addModifyNote', note: Todo): void
  (e: 'deleteTodo', todo: Todo): void
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

// If a todo is passed as a prop we insert the values in the refs
if ($props.todo) {
  _title.value = $props.todo.Title
  _description.value = $props.todo.Description
  _state.value = $props.todo.StateID
  _priority.value = $props.todo.PriorityID
  _color.value = $props.todo.ColorID
  _does_expire.value = $props.todo.Expiration.DoesExpire
  if (_does_expire.value) {
    // Ugly, I know
    let d = new Date($props.todo.Expiration.Date).toISOString().split('T')
    _date.value = d[0]
    _time.value = d[1].split('.')[0]
  }
}

function close() {
  $emits('close')
}

function deleteTodo() {
  if ($props.todo) {
    $emits('deleteTodo', $props.todo)
  }
}

function addModifyNote() {
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

  let d: Date = new Date()
  if (_does_expire.value) {
    d = new Date(_date.value + ' ' + _time.value)
  }

  let id = $props.todo?.ID
  let n: Todo = {
    ID: id ? id : 0,
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

  $emits('addModifyNote', n)
}

function keyboardLIstener(event: KeyboardEvent) {
  if (event.key == 'Escape') {
    close()
  }
}

onMounted(() => {
  window.addEventListener('keyup', keyboardLIstener)
})
onBeforeUnmount(() => {
  window.removeEventListener('keyup', keyboardLIstener)
})
</script>

<template>
  <div
    class="fixed w-[25rem] rounded border border-black bg-white drop-shadow-lg"
    :style="style"
    ref="__top_div"
  >
    <header class="flex justify-between rounded bg-gray-200" ref="__header">
      <h1 class="flex-auto p-5 text-center font-bold">
        {{ $props.modify ? 'Modify' : 'Add' }} a Note
      </h1>
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

        <div class="mt-5 flex justify-around">
          <button @click="close()" class="rounded bg-red-200 p-2 hover:bg-red-400">Close</button>
          <button @click="deleteTodo()" class="rounded bg-red-500 p-2 hover:bg-red-800">
            DELETE
          </button>
          <button @click="addModifyNote()" class="rounded bg-blue-200 p-2 hover:bg-blue-400">
            {{ $props.modify ? 'Modify' : 'Add' }}
          </button>
        </div>
      </form>
    </body>
  </div>
</template>

<style></style>
