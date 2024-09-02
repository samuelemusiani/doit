<script setup lang="ts">
import type { Todo, Options } from '@/types'
import type { PropType } from 'vue'
import { inject } from 'vue'

const $props = defineProps({
  todo: {
    type: Object as PropType<Todo>,
    required: true
  }
})

const $emits = defineEmits<{
  (e: 'close'): void
  (e: 'modify', t: Todo): void
  (e: 'delete', t: Todo): void
}>()

const _todo_options = inject('todoOptions') as Options

function close() {
  $emits('close')
}
</script>

<template>
  <div class="grid bg-black bg-opacity-30" @click="close()">
    <div
      class="fixed place-self-center rounded-lg bg-white md:w-[35rem] md:min-w-[25rem] lg:w-[50rem]"
      @click.stop=""
    >
      <header class="flex justify-end border-b border-b-gray-400 p-1">
        <button class="mr-3 h-7 w-7 rounded p-1 hover:bg-gray-200">X</button>
      </header>

      <body class="flex justify-between gap-5 p-5">
        <div>
          <h1 class="mb-5 pb-2 text-xl font-bold">
            {{ $props.todo.Title }}
          </h1>

          <p>
            {{ $props.todo.Description }}
          </p>
        </div>

        <div class="flex min-w-40 flex-col gap-4">
          <div>
            <label>State:</label>
            <div class="font-bold">
              {{ _todo_options.States[$props.todo.StateID - 1].State }}
            </div>
          </div>

          <div>
            <label class="">Priority:</label>
            <div class="font-bold">
              {{ _todo_options.Priorities[$props.todo.PriorityID - 1].Priority }}
            </div>
          </div>

          <div>
            <label>Color:</label>
            <div
              :style="{ 'background-color': _todo_options.Colors[$props.todo.ColorID - 1].Hex }"
              class="h-8 rounded"
            ></div>
          </div>
        </div>
      </body>
    </div>
  </div>
</template>

<style></style>
