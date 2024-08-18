<script setup lang="ts"></script>

<template>
  <div>
    <h1 class="m-5 text-2xl font-bold">Admin page</h1>
  </div>
  <div>
    <RouterLink :to="{ name: 'users' }" class="rounded border p-2 hover:bg-gray-100">
      Users
    </RouterLink>
  </div>
</template>

<style></style>
