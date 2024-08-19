<script setup lang="ts">
import type { User } from '@/types'
import { useFocus } from '@vueuse/core'
import type { PropType } from 'vue'
import { computed } from 'vue'
import { ref } from 'vue'

const $props = defineProps({
  tried_submitted: {
    type: Boolean as PropType<Boolean>,
    required: false
  }
})

const $data_model = defineModel<User>('data', { required: true })
const $ok_model = defineModel<boolean>('ok', { required: true })

$data_model.value.Active = true

const _error = computed(() => {
  let v = validate()
  if ($props.tried_submitted && !v.isValid) {
    return v.err
  } else {
    return ''
  }
})

const __usernameRef = ref<HTMLElement | null>(null)
useFocus(__usernameRef, { initialValue: true })

function validate(): { isValid: Boolean; err: string } {
  if (!$data_model.value.Username) {
    return { isValid: false, err: 'Username is empty' }
  }
  if (!$data_model.value.Password) {
    return { isValid: false, err: 'Password missing or too short' }
  }
  if (!$data_model.value.Email) {
    return { isValid: false, err: 'Email is empty' }
  }
  $ok_model.value = true
  return { isValid: true, err: '' }
}
</script>

<template>
  <div class="w-96">
    <form @submit.prevent="">
      <div class="*:p-2">
        <div>
          <label for="username"> Username: </label>
          <input
            type="text"
            id="username"
            ref="__usernameRef"
            v-model="$data_model.Username"
            class="rounded p-2 outline-none invalid:border-red-400 invalid:bg-red-50 enabled:border"
            minlength="3"
          />
        </div>
        <div>
          <label for="password"> Password: </label>
          <input
            type="password"
            id="password"
            v-model="$data_model.Password"
            class="rounded p-2 outline-none invalid:border-red-400 invalid:bg-red-50 enabled:border"
            minlength="4"
            placeholder="********"
          />
        </div>
        <div>
          <label for="email"> Email: </label>
          <input
            type="email"
            id="email"
            v-model="$data_model.Email"
            class="rounded p-2 outline-none invalid:border-red-400 invalid:bg-red-50 enabled:border"
          />
        </div>
        <div>
          <label for="name"> Name: </label>
          <input
            type="text"
            id="name"
            v-model="$data_model.Name"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
        <div>
          <label for="surname"> Surname: </label>
          <input
            type="text"
            id="surname"
            v-model="$data_model.Surname"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
        <div>
          <label for="admin"> Admin: </label>
          <input
            type="checkbox"
            id="admin"
            v-model="$data_model.Admin"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>

        <div>
          <label for="active"> Active: </label>
          <input
            type="checkbox"
            id="active"
            v-model="$data_model.Active"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
      </div>
      <div v-if="_error.length > 0" class="p-2 text-red-600">Error: {{ _error }}</div>
    </form>
  </div>
</template>

<style></style>
