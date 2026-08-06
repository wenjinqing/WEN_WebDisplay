<script setup>
defineProps({
  size: { type: [Number, String], default: 32 },
  color: { type: String, default: '#ffc2d4' },
  rotate: { type: [Number, String], default: 0 },
})
</script>

<template>
  <svg
    :width="size"
    :height="size"
    viewBox="0 0 64 64"
    :style="{ transform: `rotate(${rotate}deg)` }"
    aria-hidden="true"
  >
    <g :fill="color">
      <ellipse cx="32" cy="42" rx="13" ry="10" />
      <ellipse cx="15" cy="26" rx="6" ry="8" transform="rotate(-20 15 26)" />
      <ellipse cx="28" cy="18" rx="6" ry="8" />
      <ellipse cx="41" cy="18" rx="6" ry="8" />
      <ellipse cx="52" cy="26" rx="6" ry="8" transform="rotate(20 52 26)" />
    </g>
  </svg>
</template>
