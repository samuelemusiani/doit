<script setup lang="ts"></script>

<template>
  <div class="flex justify-center">
    <div>
      <div>
        <h1 class="m-5 text-2xl font-bold">Admin page</h1>
      </div>
      <div class="w-full">
        <RouterLink :to="{ name: 'users' }">
          <div class="w-full rounded border p-2 hover:bg-gray-100">Users</div>
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<style></style>
