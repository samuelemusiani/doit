<script setup lang="ts">
import type { Options, Todo } from '@/types'
import type { PropType } from 'vue'
import { inject, ref } from 'vue'

const $props = defineProps({
  todos: {
    type: Array as PropType<Todo[]>,
    required: true
  }
})

const $emits = defineEmits<{
  (e: 'selected', s: number): void
}>()

const _todo_options = inject('todoOptions') as Options

function countState(state: number) {
  return $props.todos.filter((e) => e.StateID == state).length
}

const _selected = ref<number>(0)

function select(n: number) {
  if (_selected.value == n) {
    _selected.value = 0
  } else {
    _selected.value = n
  }
  $emits('selected', _selected.value)
}
</script>

<template>
  <div class="flex flex-col gap-1">
    <template v-for="state in _todo_options.States">
      <div
        class="flex w-52 items-center justify-between rounded border p-2 hover:bg-gray-100"
        :class="{ 'bg-gray-100': _selected == state.ID }"
        @click="select(state.ID)"
      >
        <span>
          {{ state.State }}
        </span>
        <span>
          {{ countState(state.ID) }}
        </span>
      </div>
    </template>
  </div>
</template>

<style></style>
