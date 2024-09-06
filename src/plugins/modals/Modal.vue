<script setup lang="ts">
import type { PropType } from 'vue'
import { computed, inject, ref } from 'vue'

const $props = defineProps({
  title: {
    type: String as PropType<String>,
    required: true
  },
  name: {
    type: String as PropType<String>,
    required: true
  },
  accept_button: {
    type: String as PropType<String>,
    required: true
  },
  component: {
    required: false
  }
})

const _data = ref<any>({} as any)
const _can_submit = ref<boolean>(false)
const _tried_submitted = ref<boolean>(false)

const $modals = inject('$modals') as any

const _show = computed(() => {
  return $modals.active() == $props.name
})

function reset() {
  _data.value = {} as any
  _can_submit.value = false
  _tried_submitted.value = false
}

function closeModal(accept: boolean = false) {
  if (!$props.component) {
    accept ? $modals.accept() : $modals.reject()
    return
  }

  if (accept) {
    _tried_submitted.value = true

    if (_can_submit.value) {
      $modals.accept(_data.value)
      reset()
    }
  } else {
    $modals.reject()
    reset()
  }
}
</script>

<template>
  <div v-if="_show" class="bg-gray-200 bg-opacity-30">
    <div class="flex h-full items-center justify-center">
      <div
        class="flex h-full max-h-[40rem] flex-col justify-between overflow-auto rounded border bg-white p-5 shadow-xl"
      >
        <div>
          <header class="text-xl font-bold">
            {{ $props.title }}
          </header>

          <main class="my-5">
            <component
              v-if="$props.component"
              :is="$props.component"
              v-model:data="_data"
              v-model:ok="_can_submit"
              :tried_submitted="_tried_submitted"
            />
            <slot> </slot>
          </main>
        </div>

        <footer class="flex justify-evenly">
          <button @click="closeModal(false)" class="rounded border bg-red-300 p-2 hover:bg-red-500">
            Cancel
          </button>

          <button
            @click="closeModal(true)"
            class="rounded border bg-blue-300 p-2 hover:bg-blue-500"
          >
            {{ $props.accept_button }}
          </button>
        </footer>
      </div>
    </div>
  </div>
</template>

<style></style>
