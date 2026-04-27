<template>
  <NavBar v-if="auth.isAuthenticated" />
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'
import NavBar from './components/NavBar.vue'
import { useAuthStore } from './stores/auth'

const auth = useAuthStore()

onMounted(() => {
  if (auth.isAuthenticated && !auth.userName) {
    auth.fetchMe()
  }
})
</script>
