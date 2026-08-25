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

const readingNovel = ref(null)
const commentingNovel = ref(null)
</script>

<template>
  <NavBar />
  <main>
    <HeroSection />
    <NoticeBoard />
    <NovelMenu @read="readingNovel = $event" @comments="commentingNovel = $event" />
    <GalleryWall />
    <AboutAuthor />
    <InteractSection />
    <PigmiWall />
    <PigFarm />
    <PigmiClub />
    <Achievements />
    <RecapCard />
  </main>
  <SiteFooter />
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
