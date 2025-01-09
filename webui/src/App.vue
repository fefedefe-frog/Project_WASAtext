<script setup>
import { RouterLink, RouterView } from 'vue-router'
</script>
<script>
export default {
  data (){
    return {
      isAuthenticated: false, //Stato che verifica se l'utente è autenticato
      }
  },
  mounted(){
    const token = localStorage.getItem('authToken');
    if (token){
      this.isAuthenticated = true;
    }
  }
}
</script>

<template>

	<header class="navbar navbar-dark sticky-top bg-dark flex-md-nowrap p-0 shadow">
		<RouterLink class="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" to="/chats">WASAtext</RouterLink>
		<button class="navbar-toggler position-absolute d-md-none collapsed" type="button" data-bs-toggle="collapse" data-bs-target="#sidebarMenu" aria-controls="sidebarMenu" aria-expanded="false" aria-label="Toggle navigation">
			<span class="navbar-toggler-icon"></span>
		</button>
	</header>

	<div class="container-fluid">
		<div class="row">
			<nav id="sidebarMenu" class="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse" v-if="isAuthenticated">
				<div class="position-sticky pt-3 sidebar-sticky">
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/chats" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-square"/></svg>
                Chats
							</RouterLink>
						</li>
						<li class="nav-item">
							<RouterLink to="/users" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users"/></svg>
								Users
							</RouterLink>
						</li>
					</ul>

					<h6 class="sidebar-heading d-flex justify-content-between align-items-center px-3 mt-4 mb-1 text-muted text-uppercase">
						<span>Preferences</span>
					</h6>
					<ul class="nav flex-column">
						<li class="nav-item">
							<RouterLink to="/profile" class="nav-link">
								<svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user"/></svg>
								Your Profile
							</RouterLink>
						</li>
            <li class="nav-item">
              <RouterLink to="/session" class="nav-link">
                <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out"/></svg>
                Logout
              </RouterLink>
            </li>
					</ul>
				</div>
			</nav>

			<main class="col-md-9 ms-sm-auto col-lg-10 px-md-4">
				<RouterView />
			</main>
		</div>
	</div>
</template>

<style>
</style>
