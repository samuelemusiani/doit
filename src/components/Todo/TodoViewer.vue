<script setup lang="ts">
import type { Todo, Options } from '@/types'
import type { PropType } from 'vue'
import { inject, ref, onMounted, onBeforeUnmount } from 'vue'

import ColorPicker from '@/components/ColorPicker.vue'
import VueMarkdown from 'vue-markdown-render'

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
const _modifying = ref<boolean>(false)

function close() {
  if (!_deleting.value) {
    cancelModify()
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

function startModify() {
  _modifying.value = true
}

function cancelModify() {
  _modifying.value = false
}

function modify() {
  $emits('modify', $props.todo)
  _modifying.value = false
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
      class="fixed h-full place-self-center overflow-auto rounded-lg bg-white shadow-2xl md:w-[35rem] md:min-w-[25rem] lg:w-[50rem]"
      :class="{ 'max-h-[30rem]': !_modifying, 'max-h-[44rem]': _modifying }"
      @click.stop=""
    >
      <div class="flex h-full flex-col justify-between">
        <div>
          <header class="flex justify-end border-b border-b-gray-400 p-1">
            <button
              class="h-7 w-7 rounded p-1 hover:bg-gray-200"
              @click="deleteTodo()"
              v-show="!_modifying"
            >
              <span class="icon-[mdi--bin] h-5 w-5 text-red-500"></span>
            </button>

            <button
              class="h-7 w-7 rounded p-1 hover:bg-gray-200"
              @click="startModify()"
              v-show="!_modifying"
            >
              <span class="icon-[mdi--pencil] h-5 w-5 text-black"></span>
            </button>

            <button class="h-7 w-7 rounded p-1 hover:bg-gray-200" @click="close()">
              <span class="icon-[material-symbols--close] h-5 w-5 text-black"></span>
            </button>
          </header>

          <main>
            <form class="flex flex-col justify-between gap-5 p-5 lg:flex-row">
              <div>
                <div class="mb-5 pb-2">
                  <h1 class="text-xl font-bold" v-if="!_modifying">
                    {{ $props.todo.Title }}
                  </h1>
                  <input
                    v-else
                    class="rounded border p-1 text-xl font-bold"
                    v-model="$props.todo.Title"
                  />
                </div>

                <div class="prose">
                  <vue-markdown v-if="!_modifying" :source="$props.todo.Description" />
                  <textarea
                    v-else
                    class="h-52 rounded border p-1 lg:min-w-96"
                    v-model="$props.todo.Description"
                  >
                  </textarea>
                </div>
              </div>

              <div class="flex min-w-40 flex-col gap-4">
                <div>
                  <label>State:</label>
                  <div class="font-bold" v-if="!_modifying">
                    {{ _todo_options.States[$props.todo.StateID - 1].State }}
                  </div>
                  <div class="w-full" v-else>
                    <select
                      v-model="$props.todo.StateID"
                      class="w-full rounded border border-black bg-white p-2 font-bold"
                    >
                      <option
                        v-for="state in _todo_options.States"
                        :key="state.ID"
                        :value="state.ID"
                      >
                        {{ state.State }}
                      </option>
                    </select>
                  </div>
                </div>
                <div>
                  <label class="">Priority:</label>
                  <div class="font-bold" v-if="!_modifying">
                    {{ _todo_options.Priorities[$props.todo.PriorityID - 1].Priority }}
                  </div>
                  <div class="w-full" v-else>
                    <select
                      v-model="$props.todo.PriorityID"
                      class="w-full rounded border border-black bg-white p-2 font-bold"
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

                <div>
                  <label>Color:</label>
                  <div
                    :style="{
                      'background-color': _todo_options.Colors[$props.todo.ColorID - 1].Hex
                    }"
                    class="h-8 rounded"
                    v-if="!_modifying"
                  ></div>
                  <ColorPicker
                    :colors="_todo_options.Colors"
                    v-model="$props.todo.ColorID"
                    v-else
                  />
                </div>
              </div>
            </form>
          </main>
        </div>

        <footer v-show="_modifying" class="p-2">
          <div class="flex justify-evenly">
            <button class="rounded border bg-red-300 p-2 hover:bg-red-500" @click="cancelModify()">
              Cancel
            </button>

            <button class="rounded border bg-blue-300 p-2 hover:bg-blue-500" @click="modify()">
              Modify
            </button>
          </div>
        </footer>
      </div>

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
