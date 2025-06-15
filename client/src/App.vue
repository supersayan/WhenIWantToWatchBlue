<script setup lang="ts">
import {ref} from "vue"

const isBlue = ref(true)

async function toggleBlue() {
  const res = await fetch("http://localhost:8080/tv/0/flip")
  updateBlue()
}

async function updateBlue() {
  const res = await fetch("http://localhost:8080/tv/0")
  isBlue.value = (await res.json()).blue
}
</script>

<template>
  <div class="page" :class="{ blue: isBlue }" @click="toggleBlue">

  </div>
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
