<script setup lang="ts">
import type { User } from '@/types'
import type { PropType } from 'vue'

const $props = defineProps({
  user: {
    type: Object as PropType<User>,
    required: true
  }
})
</script>

<template>
  <div class="rounded border p-5">
    <div class="p-2">
      <span> Username: </span>
      <span> {{ $props.user.Username }} </span>
    </div>
    <div class="p-2">
      <span> Email: </span>
      <span> {{ $props.user.Email }} </span>
    </div>
    <div class="p-2">
      <span> Name: </span>
      <span> {{ $props.user.Name }} </span>
    </div>
    <div class="p-2">
      <span> Surname: </span>
      <span> {{ $props.user.Surname }} </span>
    </div>
    <div class="p-2">
      <span> Admin: </span>
      <span> {{ $props.user.Admin }} </span>
    </div>
    <div class="p-2">
      <span> External: </span>
      <span> {{ $props.user.External }} </span>
    </div>
  </div>
</template>

<style></style>
