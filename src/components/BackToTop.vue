<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import PawPrint from './PawPrint.vue'

const show = ref(false)

function onScroll() {
  show.value = window.scrollY > 600
}

function toTop() {
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <button v-show="show" class="to-top" aria-label="回到顶部" @click="toTop">
    <PawPrint :size="24" color="#fff" />
  </button>
</template>

<style scoped>
.to-top {
  position: fixed;
  right: 24px;
  bottom: 24px;
  z-index: 60;
  width: 52px;
  height: 52px;
  border-radius: 50%;
  border: none;
  background: var(--pink);
  box-shadow: 0 6px 18px rgba(249, 113, 143, 0.45);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.2s, background 0.2s;
}

.to-top:hover {
  background: var(--pink-deep);
  transform: translateY(-3px);
}
</style>
