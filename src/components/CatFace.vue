<script setup>
// Q版猫猫头 v2 — 圆脸圆耳幼态比例,手绘风
defineProps({
  size: { type: [Number, String], default: 200 },
})
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 200 200" aria-hidden="true">
    <g stroke-linecap="round" stroke-linejoin="round">
      <!-- 小圆耳:贴着头顶的短短小三角(Pusheen 比例) -->
      <path d="M48 100 Q 44 68 62 64 Q 80 60 84 86 Z" fill="#fff" stroke="#f9718f" stroke-width="5" />
      <path d="M152 100 Q 156 68 138 64 Q 120 60 116 86 Z" fill="#fff" stroke="#f9718f" stroke-width="5" />
      <path d="M56 92 Q 55 74 64 72 Q 74 70 76 84 Z" fill="#ffc2d4" />
      <path d="M144 92 Q 145 74 136 72 Q 126 70 124 84 Z" fill="#ffc2d4" />

      <!-- 圆脸蛋:横宽椭圆更幼态 -->
      <ellipse cx="100" cy="124" rx="66" ry="58" fill="#fff" stroke="#f9718f" stroke-width="5" />

      <!-- 左眼:大而圆,双层高光 -->
      <circle cx="72" cy="122" r="11" fill="#5b3a47" />
      <circle cx="76" cy="118" r="3.8" fill="#fff" />
      <circle cx="68.5" cy="126" r="2" fill="#fff" opacity="0.8" />

      <!-- 右眼:wink -->
      <path d="M116 122 Q 127 111 138 122" fill="none" stroke="#5b3a47" stroke-width="5.5" />

      <!-- 大腮红 -->
      <ellipse cx="58" cy="142" rx="12" ry="7" fill="#ffc2d4" opacity="0.9" />
      <ellipse cx="142" cy="142" rx="12" ry="7" fill="#ffc2d4" opacity="0.9" />

      <!-- 小鼻子 + ω 嘴 -->
      <path d="M96 134 Q 100 131 104 134 Q 100 139 96 134 Z" fill="#f9718f" />
      <path d="M100 138 Q 100 144 92 144 M100 138 Q 100 144 108 144" fill="none" stroke="#5b3a47" stroke-width="3.5" />

      <!-- 手绘弯胡须 -->
      <g stroke="#eda9be" stroke-width="3" fill="none">
        <path d="M28 118 Q 40 120 50 123" />
        <path d="M28 136 Q 40 134 50 132" />
        <path d="M172 118 Q 160 120 150 123" />
        <path d="M172 136 Q 160 134 150 132" />
      </g>
    </g>
  </svg>
</template>
