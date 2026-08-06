<script setup>
// 猪咪 v2 — 圆脸圆耳,猫耳 + 猪鼻的 Q 版粉丝形象
defineProps({
  size: { type: [Number, String], default: 120 },
})
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 200 200" aria-hidden="true">
    <g stroke-linecap="round" stroke-linejoin="round">
      <!-- 圆猫耳 -->
      <path d="M62 88 Q 44 42 70 47 Q 92 51 89 78 Z" fill="#ffd6e2" stroke="#f9718f" stroke-width="5" />
      <path d="M138 88 Q 156 42 130 47 Q 108 51 111 78 Z" fill="#ffd6e2" stroke="#f9718f" stroke-width="5" />

      <!-- 圆脑袋 -->
      <ellipse cx="100" cy="126" rx="64" ry="56" fill="#ffd6e2" stroke="#f9718f" stroke-width="5" />

      <!-- 开心眯眯眼 -->
      <path d="M62 122 Q 73 110 84 122" fill="none" stroke="#5b3a47" stroke-width="5.5" />
      <path d="M116 122 Q 127 110 138 122" fill="none" stroke="#5b3a47" stroke-width="5.5" />

      <!-- 大猪鼻子 -->
      <ellipse cx="100" cy="142" rx="19" ry="14" fill="#ffb3c6" stroke="#f9718f" stroke-width="4" />
      <ellipse cx="93" cy="142" rx="3.5" ry="5.5" fill="#e85d7f" />
      <ellipse cx="107" cy="142" rx="3.5" ry="5.5" fill="#e85d7f" />

      <!-- 腮红 -->
      <ellipse cx="58" cy="142" rx="10" ry="6" fill="#ffb3c6" opacity="0.85" />
      <ellipse cx="142" cy="142" rx="10" ry="6" fill="#ffb3c6" opacity="0.85" />

      <!-- 头顶卷卷呆毛 -->
      <path d="M100 70 Q 114 58 105 49 Q 98 43 105 35" fill="none" stroke="#f9718f" stroke-width="4" />
    </g>
  </svg>
</template>
