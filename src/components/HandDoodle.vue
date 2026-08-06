<script setup>
// 手绘涂鸦小元素库:云朵 / 星星 / 爱心 / 卷卷箭头
defineProps({
  kind: { type: String, required: true }, // cloud | star | heart | arrow
  size: { type: [Number, String], default: 48 },
})
</script>

<template>
  <!-- 云朵 -->
  <svg v-if="kind === 'cloud'" :width="size" :height="size" viewBox="0 0 80 50" aria-hidden="true">
    <path
      d="M20 42 Q 8 42 8 33 Q 8 25 18 24 Q 20 13 32 15 Q 40 6 50 15 Q 61 13 62 24 Q 72 25 71 34 Q 70 42 58 42 Z"
      fill="#fff"
      stroke="#f4a9c0"
      stroke-width="3"
      stroke-linejoin="round"
    />
  </svg>

  <!-- 四角星 -->
  <svg v-else-if="kind === 'star'" :width="size" :height="size" viewBox="0 0 24 24" aria-hidden="true">
    <path
      d="M12 2 Q 13.5 9, 22 12 Q 13.5 15, 12 22 Q 10.5 15, 2 12 Q 10.5 9, 12 2 Z"
      fill="#ffe9a8"
      stroke="#f0b429"
      stroke-width="1.6"
      stroke-linejoin="round"
    />
  </svg>

  <!-- 描边爱心 -->
  <svg v-else-if="kind === 'heart'" :width="size" :height="size" viewBox="0 0 32 30" aria-hidden="true">
    <path
      d="M16 27 C 8 20, 2.5 14.5, 4 8.5 C 5.5 3.5, 12 3, 16 8.5 C 20 3, 26.5 3.5, 28 8.5 C 29.5 14.5, 24 20, 16 27 Z"
      fill="none"
      stroke="#f9718f"
      stroke-width="2.5"
      stroke-linejoin="round"
    />
  </svg>

  <!-- 卷卷箭头(指向右下) -->
  <svg v-else :width="size" :height="size" viewBox="0 0 60 60" aria-hidden="true">
    <path
      d="M6 8 Q 26 10, 34 26 Q 40 38, 50 46 M50 46 l -12 -1 M50 46 l -1 -12"
      fill="none"
      stroke="#e85d7f"
      stroke-width="3"
      stroke-linecap="round"
      stroke-linejoin="round"
    />
  </svg>
</template>
