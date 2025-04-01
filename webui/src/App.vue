<script setup>
import {RouterLink, RouterView} from 'vue-router'
</script>
<script>
export default {
  data() {
    return {
      isAuthenticated: false
    };
  },
  created() {
    this.checkAuth();
    // Ascolto i cambiamenti nel localStorage (da altre schede)
    window.addEventListener("storage", this.checkAuth);
  },
  beforeUnmount() {
    window.removeEventListener("storage", this.checkAuth);
  },
  methods: {
    checkAuth() {
      const token = localStorage.getItem('authToken');
      if (token && token !== "") {
        this.isAuthenticated = true;
      }
    },
    logout() {
      localStorage.removeItem("authToken");
      this.checkAuth();
    }
  }
};
</script>
<template>
  <header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
    <RouterLink class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" to="/chats">WASAtext</RouterLink>
    <button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon" />
    </button>
  </header>

  <div class="container-fluid">
    <main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
      <RouterView />
    </main>
  </div>
</template>

<style>

main {
  overflow: hidden;
  width: 100%;
  height: 90vh;
}

</style>
