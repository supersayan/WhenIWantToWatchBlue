<script setup lang="ts">
import { RouterView } from 'vue-router'
</script>

<template>
  <RouterView />
</template>

<style scoped>
.page {
  height: calc(100vh - 4rem);
  width: 100%;
  background-color: orange;
}

.blue {
  background-color: blue;
}
</style>
