<script setup>
// Q版猫猫头 v2 — 圆脸圆耳幼态比例,手绘风
defineProps({
  size: { type: [Number, String], default: 200 },
})
</script>

<template>
  <svg :width="size" :height="size" viewBox="0 0 200 200" aria-hidden="true">
    <g stroke-linecap="round" stroke-linejoin="round">
      <!-- 小圆耳带一点尖:短而微翘(Pusheen 比例 + 小尖角) -->
      <path d="M46 102 Q 50 66 66 58 Q 82 64 84 90 Z" fill="#fff" stroke="#f9718f" stroke-width="5" />
      <path d="M154 102 Q 150 66 134 58 Q 118 64 116 90 Z" fill="#fff" stroke="#f9718f" stroke-width="5" />
      <path d="M56 94 Q 59 74 67 68 Q 76 74 76 86 Z" fill="#ffc2d4" />
      <path d="M144 94 Q 141 74 133 68 Q 124 74 124 86 Z" fill="#ffc2d4" />

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
