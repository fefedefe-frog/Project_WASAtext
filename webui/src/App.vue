<script setup>
import {RouterView} from 'vue-router'
</script>

<template>
  <div class="app-div">
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style>
.app-div{
  width: 100vw;
  height: 100vh;
  /* schema di sfondo preso da uiverse */
  background-color: #313131;
  background-image: radial-gradient(rgba(255, 255, 255, 0.171) 2px, transparent 0);
  background-size: 30px 30px;
  background-position: -5px -5px;
}

.app-main {
  width: 100%;
  height: 100%;
  margin: auto;
  padding: 20px;

  display: flex;
  justify-content: center;
  align-items: center;
}
</style>
