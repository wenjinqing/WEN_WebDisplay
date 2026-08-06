<script setup>
// 区块标题 + 手绘波浪下划线
defineProps({
  title: { type: String, required: true },
  sub: { type: String, default: '' },
})
</script>

<template>
  <div class="sec-head" v-reveal>
    <h2 class="section-title">{{ title }}</h2>
    <svg class="squiggle" viewBox="0 0 140 12" width="130" height="12" aria-hidden="true">
      <path
        d="M4 7 Q 12 2, 20 7 T 36 7 T 52 7 T 68 7 T 84 7 T 100 7 T 116 7 T 132 7"
        fill="none"
        stroke="#f9718f"
        stroke-width="3.5"
        stroke-linecap="round"
      />
    </svg>
    <p v-if="sub" class="section-sub">{{ sub }}</p>
  </div>
</template>

<style scoped>
.sec-head {
  text-align: center;
  margin-bottom: 40px;
}

.section-title {
  margin-bottom: 0;
}

.squiggle {
  display: block;
  margin: 6px auto 8px;
}

.section-sub {
  margin-bottom: 0;
}
</style>
