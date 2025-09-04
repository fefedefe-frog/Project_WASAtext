<script setup>
import {RouterLink, RouterView} from 'vue-router'
</script>
<script>
export default {
  methods: {
    logout(){
      sessionStorage.removeItem('authToken')
      this.$router.push('/session')
    },
    myProfile(){
      let usrId= sessionStorage.getItem('usrId');
      this.$router.push(`/users/${usrId}`);
    }
  }
}
</script>

<template>
  <div class="app-div">
    <main v-if="$route.name !== 'login'" class="app-main">
      <div class="navmenu bobby">
        <div class="navmenu-submenu">
          <RouterLink to="/users" class="navmenu-button">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#users" /></svg>
          </RouterLink>
          <RouterLink to="/chats" class="navmenu-button">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#message-circle" /></svg>
          </RouterLink>
          <RouterLink to="/newChat" class="navmenu-button">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#plus-circle" /></svg>
          </RouterLink>
          <RouterLink to="/home" class="navmenu-button">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#home" /></svg>
          </RouterLink>
        </div>
        <div class="spacer" />
        <div class="navmenu-submenu">
          <button class="navmenu-button" @click="myProfile">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#user" /></svg>
          </button>
          <button class="navmenu-button logout-button" @click="logout">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#log-out" /></svg>
          </button>
        </div>
      </div>
      <div class="router-view">
        <RouterView />
      </div>
    </main>
    <main v-else>
      <RouterView />
    </main>
  </div>
</template>

<style>
.app-div{
  width: 100vw;
  height: 100vh;

  padding: 15px;

  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;


  /* schema di sfondo preso da uiverse */
  background-color: #313131;
  background-image: radial-gradient(rgba(255, 255, 255, 0.171) 2px, transparent 0);
  background-size: 30px 30px;
  background-position: -5px -5px;
}

.app-main {
  width: 100%;
  height: 100%;

  display: flex;
  justify-content: center;
  align-items: center;
}

.navmenu {
  display: flex;
  flex-direction: column;
  align-items: center;

  height: 100%;
  width: 5%;
  padding: 2px;

  margin-right: 5px;
  overflow: hidden;
  overflow-y: scroll;
}

.navmenu-submenu {
  width: 100%;
  height: fit-content;

  display: flex;
  flex-direction: column;
  align-items: center;
}

.spacer {
  width: 100%;
  height: 100%;

  border-top: 1px solid dimgray;
  border-bottom: 1px solid dimgray;
}

.navmenu-button {
  width: 70%;
  aspect-ratio: 1/1;

  margin: 5px;

  border-radius: 10px;
  border: 2px dashed lightseagreen;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;

  text-decoration: none;
  color: white;
  background-color: lightseagreen;
  cursor: pointer;
  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
  transition: .4s;
}

.logout-button {
  background-color: red;
  border: 2px dashed red;
}

.navmenu-button:hover {
  transition: .4s;
  border: 2px dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.logout-button:hover {
  border: 2px dashed red;
  color: red;
}

.navmenu-button:active {
  background-color: lightseagreen;
}

.logout-button:active{
  background-color: red;

}

.router-view {
  height: 100%;
  width: 100%;

  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}

.bobby{
  border-radius: 10px;
  border: groove 2px solid;

  background: linear-gradient(45deg, #666666, #999999);
}
</style>