<script setup lang="ts">
import type { Todo, Options } from '@/types'
import type { PropType } from 'vue'
import { inject, ref, onMounted, onBeforeUnmount } from 'vue'

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
const $modals = inject('$modals') as any

const _deleting = ref<boolean>(false)

function close() {
  if (!_deleting.value) {
    $emits('close')
  }
}

function deleteTodo() {
  _deleting.value = true
  $modals
    .show('deleteTodo')
    .then(
      () => {
        $emits('delete', $props.todo)
      },
      () => {}
    )
    .finally(() => {
      _deleting.value = false
    })
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
      class="fixed place-self-center rounded-lg bg-white shadow-2xl md:w-[35rem] md:min-w-[25rem] lg:w-[50rem]"
      @click.stop=""
    >
      <header class="flex justify-end border-b border-b-gray-400 p-1">
        <button class="h-7 w-7 rounded p-1 hover:bg-gray-200" @click="deleteTodo()">
          <span class="icon-[mdi--bin] h-5 w-5 text-red-500"></span>
        </button>

        <button class="h-7 w-7 rounded p-1 hover:bg-gray-200" @click="">
          <span class="icon-[mdi--pencil] h-5 w-5 text-black"></span>
        </button>
        <button class="h-7 w-7 rounded p-1 hover:bg-gray-200" @click="close()">
          <span class="icon-[material-symbols--close] h-5 w-5 text-black"></span>
        </button>
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

      <Transition>
        <Modal
          name="deleteTodo"
          title="Delete todo?"
          accept_button="DELETE"
          class="absolute left-0 top-0 h-full w-full"
        >
          Do you really want to delete the TODO?
        </Modal>
      </Transition>
    </div>
  </div>
</template>

<style scoped>
.v-enter-active,
.v-leave-active {
  transition: opacity 0.1s ease;
}

.v-enter-from,
.v-leave-to {
  opacity: 0;
}
</style>
