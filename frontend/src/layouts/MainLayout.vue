<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import Sidebar from './components/Sidebar.vue'
import Header from './components/Header.vue'

const isSidebarOpen = ref(false)
const isSidebarCollapsed = ref(false)
</script>

<template>
  <div class="flex h-screen overflow-hidden bg-background text-foreground transition-colors duration-300">
    <Sidebar 
      v-model:isSidebarOpen="isSidebarOpen"
      v-model:isSidebarCollapsed="isSidebarCollapsed"
    />

    <!-- Main Content Wrapper -->
    <div 
      class="flex-1 flex flex-col min-w-0 overflow-hidden transition-all duration-300 ease-in-out"
      :class="[
        isSidebarCollapsed ? 'lg:pl-20' : 'lg:pl-64'
      ]"
    >
      <Header @toggleSidebar="isSidebarOpen = !isSidebarOpen" />
      
      <!-- Page Content -->
      <main class="flex-1 p-4 sm:p-6 overflow-auto bg-muted/20">
        <RouterView />
      </main>
    </div>
  </div>
</template>