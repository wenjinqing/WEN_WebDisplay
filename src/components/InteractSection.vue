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
  grid-template-columns: 2fr 3fr; /* 留言板更宽,催更墙收窄 */
  gap: 26px;
  max-width: 1020px;
  margin: 0 auto;
  align-items: stretch; /* 两卡等高 */
}

@media (max-width: 720px) {
  .grid {
    grid-template-columns: 1fr;
  }
}
</style>
