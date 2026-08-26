<script setup>
import { ref } from 'vue'
import NavBar from './components/NavBar.vue'
import HeroSection from './components/HeroSection.vue'
import NoticeBoard from './components/NoticeBoard.vue'
import AboutAuthor from './components/AboutAuthor.vue'
import NovelMenu from './components/NovelMenu.vue'
import NovelReader from './components/NovelReader.vue'
import CommentBox from './components/CommentBox.vue'
import GalleryWall from './components/GalleryWall.vue'
import InteractSection from './components/InteractSection.vue'
import PigmiWall from './components/PigmiWall.vue'
import PigFarm from './components/PigFarm.vue'
import PigmiClub from './components/PigmiClub.vue'
import Achievements from './components/Achievements.vue'
import RecapCard from './components/RecapCard.vue'
import SiteFooter from './components/SiteFooter.vue'
import BackToTop from './components/BackToTop.vue'
import DanmakuLayer from './components/DanmakuLayer.vue'
import PawRain from './components/PawRain.vue'
import PixelPet from './components/PixelPet.vue'
import WhackGame from './components/WhackGame.vue'
import WeatherLayer from './components/WeatherLayer.vue'
import FunMenu from './components/FunMenu.vue'
import ShareCard from './components/ShareCard.vue'
import NightMode from './components/NightMode.vue'
import AppTabBar from './components/AppTabBar.vue'
import ProfilePage from './components/ProfilePage.vue'
import CheckinCalendar from './components/CheckinCalendar.vue'
import ShakeBox from './components/ShakeBox.vue'
import HubApp from './hub/HubApp.vue'
import { isApp } from './app.js'

const readingNovel = ref(null)
const commentingNovel = ref(null)
const tab = ref('home')
const showHub = ref(false)

// APP 模式下,组件里的页内锚点(#novels 等)改为切换 Tab
if (isApp) {
  const hashToTab = {
    '#home': 'home',
    '#notice': 'home',
    '#novels': 'novels',
    '#gallery': 'gallery',
    '#interact': 'interact',
    '#pigmi': 'interact',
  }
  document.addEventListener(
    'click',
    (e) => {
      // 「进入猪咪聚集地」→ APP 内嵌打开,不跳线上网页
      const hubLink = e.target.closest('a[href="/hub"]')
      if (hubLink) {
        e.preventDefault()
        e.stopPropagation()
        showHub.value = true
        return
      }
      // 聚集地里的「回猫咖」→ 关闭内嵌页
      if (showHub.value && e.target.closest('a[href="/"]')) {
        e.preventDefault()
        e.stopPropagation()
        showHub.value = false
        return
      }
      const a = e.target.closest('a[href^="#"]')
      if (!a) return
      const t = hashToTab[a.getAttribute('href')]
      if (t) {
        e.preventDefault()
        e.stopPropagation()
        tab.value = t
        window.scrollTo({ top: 0 })
      }
    },
    true
  )
}
</script>

<template>
  <!-- ============ 网页模式:顶部导航 + 长滚动页 ============ -->
  <template v-if="!isApp">
    <NavBar />
    <main>
      <HeroSection />
      <NoticeBoard />
      <NovelMenu @read="readingNovel = $event" @comments="commentingNovel = $event" />
      <GalleryWall />
      <AboutAuthor />
      <InteractSection />
      <ShakeBox />
      <PigmiWall />
      <PigFarm />
      <PigmiClub />
      <CheckinCalendar />
      <Achievements />
      <RecapCard />
    </main>
    <SiteFooter />
  </template>

  <!-- ============ APP 模式:底部 Tab 切换 ============ -->
  <main v-else class="app-main">
    <div v-show="tab === 'home'">
      <HeroSection />
      <NoticeBoard />
      <RecapCard />
    </div>
    <div v-show="tab === 'novels'">
      <NovelMenu @read="readingNovel = $event" @comments="commentingNovel = $event" />
    </div>
    <div v-show="tab === 'gallery'">
      <GalleryWall />
    </div>
    <div v-show="tab === 'interact'">
      <InteractSection />
      <ShakeBox />
      <PigmiWall />
      <PigFarm />
      <PigmiClub />
    </div>
    <div v-show="tab === 'me'">
      <ProfilePage />
      <CheckinCalendar />
      <Achievements />
      <AboutAuthor />
      <SiteFooter />
    </div>
    <AppTabBar v-model="tab" />
    <!-- 猪咪聚集地:APP 内嵌全屏页,不再跳线上 -->
    <div v-if="showHub" class="hub-overlay">
      <button class="hub-close" @click="showHub = false">✕ 收起聚集地</button>
      <HubApp />
    </div>
  </main>

  <!-- 悬浮层两种模式共用 -->
  <BackToTop />
  <DanmakuLayer />
  <PawRain />
  <PixelPet />
  <WhackGame />
  <WeatherLayer />
  <FunMenu />
  <ShareCard />
  <NightMode />
  <NovelReader v-if="readingNovel" :novel="readingNovel" @close="readingNovel = null" />
  <CommentBox v-if="commentingNovel" :novel="commentingNovel" @close="commentingNovel = null" />
</template>

<style scoped>
.app-main {
  padding-bottom: calc(76px + env(safe-area-inset-bottom));
}

/* APP 内嵌聚集地:全屏覆盖,盖在 Tab 栏之上 */
.hub-overlay {
  position: fixed;
  inset: 0;
  z-index: 70;
  overflow-y: auto;
  background: var(--bg, #fff6f8);
  -webkit-overflow-scrolling: touch;
}

.hub-close {
  position: sticky;
  top: 10px;
  left: 10px;
  z-index: 72;
  margin: 10px 0 0 10px;
  border: none;
  border-radius: 999px;
  padding: 8px 16px;
  background: rgba(91, 58, 71, 0.85);
  color: #fff;
  font-size: 0.85rem;
  font-family: inherit;
  cursor: pointer;
  backdrop-filter: blur(4px);
}
</style>
