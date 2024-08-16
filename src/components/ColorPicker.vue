<script setup lang="ts">
import type { TodoColor } from '@/types'
import type { PropType } from 'vue'
import { ref } from 'vue'
const $props = defineProps({
  colors: {
    type: Array as PropType<TodoColor[]>,
    required: true
  }
})

const $model = defineModel()

function update(n: number) {
  $model.value = n
}
</script>

<template>
  <div>
    <div class="grid grid-cols-6 gap-4">
      <div v-for="color in $props.colors" :key="color.ID" class="rounded-md border border-black">
        <div
          :style="{ 'background-color': color.Hex }"
          class="grid min-h-7 rounded hover:border hover:border-white"
          @click="update(color.ID)"
        >
          <span v-show="$model == color.ID" class="self-center text-center"> X </span>
        </div>
      </div>
    </div>
  </div>
</template>

<style></style>
