<script setup>
// Q版猫猫头 — 手绘原创 SVG,眨眼表情
defineProps({
  size: { type: [Number, String], default: 200 },
})
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 200 200" aria-hidden="true">
    <!-- 耳朵 -->
    <path d="M45 78 L30 28 L82 52 Z" fill="#fff" stroke="#f9718f" stroke-width="5" stroke-linejoin="round" />
    <path d="M155 78 L170 28 L118 52 Z" fill="#fff" stroke="#f9718f" stroke-width="5" stroke-linejoin="round" />
    <path d="M50 66 L42 40 L70 54 Z" fill="#ffc2d4" />
    <path d="M150 66 L158 40 L130 54 Z" fill="#ffc2d4" />
    <!-- 头 -->
    <circle cx="100" cy="118" r="62" fill="#fff" stroke="#f9718f" stroke-width="5" />
    <!-- 左眼(圆) -->
    <circle cx="76" cy="112" r="7" fill="#5b3a47" />
    <circle cx="78" cy="110" r="2.5" fill="#fff" />
    <!-- 右眼(wink) -->
    <path d="M114 112 Q124 104 134 112" fill="none" stroke="#5b3a47" stroke-width="5" stroke-linecap="round" />
    <!-- 腮红 -->
    <ellipse cx="64" cy="132" rx="10" ry="6" fill="#ffc2d4" opacity="0.9" />
    <ellipse cx="136" cy="132" rx="10" ry="6" fill="#ffc2d4" opacity="0.9" />
    <!-- 鼻子+嘴 ω -->
    <path d="M97 124 L103 124 L100 128 Z" fill="#f9718f" />
    <path d="M100 128 Q100 134 92 134 M100 128 Q100 134 108 134" fill="none" stroke="#5b3a47" stroke-width="3.5" stroke-linecap="round" />
    <!-- 胡须 -->
    <g stroke="#e8a7ba" stroke-width="3" stroke-linecap="round">
      <path d="M30 112 L52 116" />
      <path d="M30 128 L52 126" />
      <path d="M170 112 L148 116" />
      <path d="M170 128 L148 126" />
    </g>
  </svg>
</template>
