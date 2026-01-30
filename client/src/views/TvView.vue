<script setup lang="ts">
import {onMounted, ref} from "vue"
import { useRoute } from "vue-router";

const route = useRoute()

const tvId = parseInt(route.params.id[0])
const API_BASE = 'http://localhost:8080';
const isBlue = ref(true)

async function toggleBlue() {
  const res = await fetch(`${API_BASE}/tv/${tvId}/flip`, { method: "POST" })
  isBlue.value = (await res.json()).data
}

async function getBlue() {
  const res = await fetch(`${API_BASE}/tv/${tvId}`)
  isBlue.value = (await res.json()).data
}

onMounted(() => {
  getBlue();
})
</script>

<template>
  <div class="tv" :class="{ blue: isBlue }" @click="toggleBlue">
    <img v-show="isBlue" src="@/assets/tv-blue.svg" />
    <img v-show="!isBlue" src="@/assets/tv-orange.svg" />
  </div>
</template>

<style scoped>
.tv {
  height: calc(100vh - 4rem);
  width: 100%;
  background-color: orange;
}

.tv img {
  width: 100%;
  height: 100%;
}

.blue {
  background-color: blue;
}
</style>
