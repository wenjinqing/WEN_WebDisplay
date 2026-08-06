<script setup>
// 猫猫吉祥物 — 梦幻可爱猫(いらすとや 免费素材,页脚已标注)
import catUrl from '../assets/cat-yumeko.png'

defineProps({
  size: { type: [Number, String], default: 200 },
})
</script>

<template>
  <img :src="catUrl" :width="size" alt="梦幻可爱猫" draggable="false" />
</template>
