<script setup>
import CatFace from './CatFace.vue'
import PigmiFace from './PigmiFace.vue'
import PawPrint from './PawPrint.vue'
import { site } from '../data.js'
</script>

<template>
  <section id="home" class="hero">
    <!-- 漂浮的猫爪背景 -->
    <div class="paws" aria-hidden="true">
      <PawPrint :size="44" class="paw p1" :rotate="-24" />
      <PawPrint :size="30" class="paw p2" :rotate="18" />
      <PawPrint :size="56" class="paw p3" :rotate="40" />
      <PawPrint :size="26" class="paw p4" :rotate="-10" />
    </div>

    <div class="container hero-inner">
      <div class="hero-text">
        <span class="badge font-cute">☕ 营业中 · OPEN</span>
        <h1 class="font-cute">{{ site.title }}</h1>
        <p class="slogan">{{ site.slogan }}</p>
        <p class="sub">
          这里是 <b>{{ site.author }}</b> 的非官方粉丝后援会
        </p>
        <div class="actions">
          <a href="#novels" class="btn btn-primary">📖 去看小说菜单</a>
          <a :href="site.authorPixiv" target="_blank" rel="noopener" class="btn btn-ghost">
            作者主页 ↗
          </a>
        </div>
      </div>

      <div class="hero-mascots">
        <CatFace :size="220" class="cat" />
        <PigmiFace :size="110" class="pigmi" />
      </div>
    </div>
  </section>
</template>

<style scoped>
.hero {
  position: relative;
  padding: 96px 0 80px;
  background: linear-gradient(160deg, var(--bg) 0%, var(--bg-deep) 100%);
  overflow: hidden;
}

.hero-inner {
  display: flex;
  align-items: center;
  gap: 48px;
}

.hero-text {
  flex: 1.2;
}

.badge {
  display: inline-block;
  padding: 4px 16px;
  border-radius: 999px;
  background: #fff;
  border: 2px solid var(--pink-soft);
  color: var(--pink-deep);
  font-size: 0.9rem;
  margin-bottom: 20px;
}

h1 {
  font-size: clamp(2rem, 6vw, 3.4rem);
  line-height: 1.25;
  color: var(--ink);
  margin-bottom: 16px;
}

.slogan {
  font-size: 1.15rem;
  color: var(--pink-deep);
  margin-bottom: 8px;
}

.sub {
  color: var(--muted);
  margin-bottom: 32px;
}

.sub b {
  color: var(--pink-deep);
}

.actions {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}

.hero-mascots {
  position: relative;
  flex-shrink: 0;
}

.cat {
  display: block;
  animation: bob 3.5s ease-in-out infinite;
}

.pigmi {
  position: absolute;
  right: -24px;
  bottom: -28px;
  animation: bob 3.5s ease-in-out 0.8s infinite;
}

@keyframes bob {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.paws .paw {
  position: absolute;
  opacity: 0.5;
}

.p1 { top: 12%; left: 4%; }
.p2 { top: 64%; left: 46%; }
.p3 { bottom: 8%; right: 6%; }
.p4 { top: 8%; right: 28%; }

@media (max-width: 720px) {
  .hero-inner {
    flex-direction: column-reverse;
    text-align: center;
    gap: 32px;
  }
  .actions {
    justify-content: center;
  }
  .hero-mascots {
    margin: 0 auto;
  }
  .paws .paw {
    display: none;
  }
}

@media (prefers-reduced-motion: reduce) {
  .cat, .pigmi { animation: none; }
}
</style>
