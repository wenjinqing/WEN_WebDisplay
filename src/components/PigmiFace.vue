<script setup>
// 猪咪吉祥物 — 梦幻可爱小猪(いらすとや 免费素材,页脚已标注)
import pigUrl from '../assets/pig-yumeko.png'

defineProps({
  size: { type: [Number, String], default: 120 },
})
</script>

<template>
  <img :src="pigUrl" :width="size" alt="猪咪" draggable="false" />
</template>
