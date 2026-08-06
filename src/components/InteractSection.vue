<script setup>
import UrgeWall from './UrgeWall.vue'
import MessageBoard from './MessageBoard.vue'
import SectionTitle from './SectionTitle.vue'
</script>

<template>
  <section id="interact" class="interact">
    <div class="container">
      <SectionTitle title="互动区" sub="催更和留言,都是对猫猫酱的爱" />
      <div class="grid" v-reveal>
        <UrgeWall />
        <MessageBoard />
      </div>
    </div>
  </section>
</template>

<style scoped>
.interact {
  background: var(--bg-deep);
}

.grid {
  display: grid;
  grid-template-columns: 1fr 1.2fr;
  gap: 28px;
  max-width: 980px;
  margin: 0 auto;
  align-items: start;
}

@media (max-width: 720px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
