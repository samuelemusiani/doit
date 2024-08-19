<script setup lang="ts">
import { deleteUser, modifyUser } from '@/lib/api'
import type { User } from '@/types'
import type { PropType } from 'vue'
import { inject } from 'vue'
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const $props = defineProps({
  user: {
    type: Object as PropType<User>,
    required: true
  }
})

const $router = useRouter()

const _modify = ref<boolean>(false)
const _error = ref<string>('')

const _username = ref<string>('')
const _password = ref<string>('')
const _email = ref<string>('')
const _name = ref<string>('')
const _surname = ref<string>('')
const _admin = ref<boolean>(false)
const _active = ref<boolean>(false)

function userToRefs(u: User) {
  _username.value = u.Username
  _password.value = ''
  _email.value = u.Email
  _name.value = u.Name
  _surname.value = u.Surname
  _admin.value = u.Admin
  _active.value = u.Active
}

function close() {
  // Reset properties
  userToRefs($props.user)

  _modify.value = false
}

function done() {
  let u: User = {
    ID: $props.user.ID,
    Username: _username.value,
    Email: _email.value,
    Name: _name.value,
    Surname: _surname.value,
    Admin: _admin.value,
    Active: _active.value
  }

  if (_password.value.length > 0) {
    u.Password = _password.value
  }

  modifyUser(u)
    .then(() => (_modify.value = false))
    .catch(() => (_error.value = 'Could not update user'))
}

userToRefs($props.user)

const $modals = inject('$modals') as any
function _deleteUser() {
  $modals.show('deletePrompt').then(
    () => {
      deleteUser($props.user.ID)
        .then(() => {
          $router.go(-1)
        })
        .catch((err) => {
          console.error(err)
        })
    },
    () => {
      // Rejected
    }
  )
}
</script>

<template>
  <div class="rounded border p-5">
    <div class="flex justify-between border-b pb-5">
      <div class="flex items-center">
        <h1 class="font-bol text-2xl">Profile</h1>
      </div>
      <button class="rounded border p-2 hover:bg-gray-100" @click="_modify = true" v-if="!_modify">
        Modify
      </button>
      <button
        class="rounded border bg-red-500 p-2 hover:bg-red-600"
        @click="_deleteUser"
        v-if="_modify"
      >
        DELETE
      </button>
    </div>
    <form @submit.prevent="">
      <div class="*:p-2">
        <div>
          <label for="username"> Username: </label>
          <input
            type="text"
            id="username"
            v-model="_username"
            disabled
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
        <div v-if="_modify">
          <label for="password"> Password: </label>
          <input
            type="password"
            id="password"
            v-model="_password"
            class="rounded p-2 outline-none enabled:border"
            placeholder="********"
          />
        </div>
        <div>
          <label for="email"> Email: </label>
          <input
            type="text"
            id="email"
            v-model="_email"
            :disabled="!_modify"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
        <div>
          <label for="name"> Name: </label>
          <input
            type="text"
            id="name"
            v-model="_name"
            :disabled="!_modify"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
        <div>
          <label for="surname"> Surname: </label>
          <input
            type="text"
            id="surname"
            v-model="_surname"
            :disabled="!_modify"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
        <div v-if="$props.user.Admin">
          <label for="admin"> Admin: </label>
          <input
            type="checkbox"
            id="admin"
            v-model="_admin"
            :disabled="!_modify"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>

        <div>
          <label for="active"> Active: </label>
          <input
            type="checkbox"
            id="active"
            v-model="_active"
            :disabled="!_modify"
            class="rounded p-2 outline-none enabled:border"
          />
        </div>
      </div>

      <div class="mt-5 flex justify-evenly" v-if="_modify">
        <button class="w-20 rounded bg-red-300 p-2 hover:bg-red-500" @click="close">Cancel</button>

        <button class="w-20 rounded border bg-blue-300 p-2 hover:bg-blue-500" @click="done">
          Done
        </button>
      </div>
      <div v-if="_modify && _error.length > 0" class="p-2 text-red-600">Error: {{ _error }}</div>
    </form>

    <Transition>
      <Modal
        name="deletePrompt"
        title="Delete user"
        accept_button="DELETE"
        class="absolute left-0 top-0 h-full w-full"
      >
        Do you really want to delete the user?
      </Modal>
    </Transition>
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
