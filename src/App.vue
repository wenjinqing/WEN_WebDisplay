<script setup>
import { ref } from 'vue'
import NavBar from './components/NavBar.vue'
import HeroSection from './components/HeroSection.vue'
import NoticeBoard from './components/NoticeBoard.vue'
import AboutAuthor from './components/AboutAuthor.vue'
import NovelMenu from './components/NovelMenu.vue'
import NovelReader from './components/NovelReader.vue'
import GalleryWall from './components/GalleryWall.vue'
import PigmiClub from './components/PigmiClub.vue'
import SiteFooter from './components/SiteFooter.vue'
import BackToTop from './components/BackToTop.vue'

const readingNovel = ref(null)
</script>

<template>
  <NavBar />
  <main>
    <HeroSection />
    <NoticeBoard />
    <AboutAuthor />
    <NovelMenu @read="readingNovel = $event" />
    <GalleryWall />
    <PigmiClub />
  </main>
  <SiteFooter />
  <BackToTop />
  <NovelReader v-if="readingNovel" :novel="readingNovel" @close="readingNovel = null" />
</template>
