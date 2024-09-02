<script setup lang="ts">
import UserAdd from '@/components/User/UserAdd.vue'

import { onMounted, ref, inject, toRaw } from 'vue'
import type { User } from '@/types'
import { addUser, getUsers } from '@/lib/api'

const _users = ref<User[]>([])

const $modals = inject('$modals') as any

function newUser() {
  $modals.show('addUserModal').then(
    (data: any) => {
      let user: User = toRaw(data)
      addUser(user)
        .then(() => {
          fetchUsers()
        })
        .catch((err) => {
          console.error(err)
        })
    },
    () => {
      // Modal rejected
    }
  )
}

function fetchUsers() {
  getUsers().then((users) => (_users.value = users))
}

onMounted(() => {
  fetchUsers()
})
</script>

<template>
  <div class="flex justify-center p-2">
    <div class="w-full max-w-[30rem] rounded border p-5">
      <div class="mt-2 flex flex-col justify-between md:flex-row">
        <h1 class="text-center font-bold">Users list</h1>
        <div>Total: {{ _users.length }}</div>
      </div>
      <div class="mt-5 flex flex-col gap-2">
        <template v-for="user in _users" :key="user.ID">
          <RouterLink :to="{ name: 'user_details', params: { id: user.ID } }">
            <div class="rounded border p-2 hover:bg-gray-100">
              {{ user.Username }}
            </div>
          </RouterLink>
        </template>
      </div>
    </div>
  </div>
  <button
    class="fixed bottom-5 right-5 rounded bg-orange-200 p-5 hover:bg-orange-400"
    @click="newUser"
  >
    Add User
  </button>

  <Transition>
    <Modal
      name="addUserModal"
      title="Add user"
      accept_button="Create"
      :component="UserAdd"
      class="absolute left-0 top-0 h-full w-full"
    >
    </Modal>
  </Transition>
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
