<script setup>
// 猪咪 — 猫耳 + 猪鼻子的 Q 版粉丝形象,手绘原创 SVG
defineProps({
  size: { type: [Number, String], default: 120 },
})
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 200 200" aria-hidden="true">
    <!-- 猫耳朵 -->
    <path d="M48 80 L36 34 L84 56 Z" fill="#ffd6e2" stroke="#f9718f" stroke-width="5" stroke-linejoin="round" />
    <path d="M152 80 L164 34 L116 56 Z" fill="#ffd6e2" stroke="#f9718f" stroke-width="5" stroke-linejoin="round" />
    <!-- 圆脑袋(猪咪粉) -->
    <circle cx="100" cy="120" r="60" fill="#ffd6e2" stroke="#f9718f" stroke-width="5" />
    <!-- 眼睛:开心眯眯眼 -->
    <path d="M66 112 Q76 102 86 112" fill="none" stroke="#5b3a47" stroke-width="5" stroke-linecap="round" />
    <path d="M114 112 Q124 102 134 112" fill="none" stroke="#5b3a47" stroke-width="5" stroke-linecap="round" />
    <!-- 猪鼻子 -->
    <ellipse cx="100" cy="132" rx="18" ry="13" fill="#ffb3c6" stroke="#f9718f" stroke-width="4" />
    <ellipse cx="93" cy="132" rx="3.5" ry="5" fill="#e85d7f" />
    <ellipse cx="107" cy="132" rx="3.5" ry="5" fill="#e85d7f" />
    <!-- 腮红 -->
    <ellipse cx="62" cy="134" rx="9" ry="5" fill="#ffb3c6" opacity="0.8" />
    <ellipse cx="138" cy="134" rx="9" ry="5" fill="#ffb3c6" opacity="0.8" />
    <!-- 小卷尾巴尖从头顶冒出来 -->
    <path d="M100 60 Q112 48 104 40 Q98 34 104 28" fill="none" stroke="#f9718f" stroke-width="4" stroke-linecap="round" />
  </svg>
</template>
