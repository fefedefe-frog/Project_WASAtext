<script>
import {RouterView} from 'vue-router'

export default {
  components: {RouterView},
  data: function () {
    return {
      token: '',
      usrId: '',
      errormsg: null,
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.usrId= sessionStorage.getItem('usrId');
  },
  methods: {
    componentsErrorHandler(error){
      this.errormsg= error;
    },
  }
}
</script>

<template>
  <div class="main-container">
    <div class="lists bobby">
      <SidebarList items="chats" :bannerComponent="'ChatBanner'" @error="componentsErrorHandler"/>
    </div>
    <div class="chat-container bobby">
      <RouterView />
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
.main-container {
  height: 100%;
  width: 100%;

  display: flex;
  flex-direction: row;

  position: relative;
}


.lists {
  display: flex;
  flex-direction: column;
  align-items: center;

  height: 100%;
  width: 25%;
  padding: 5px;

  margin-right: 5px;
}

.chat-container {
  height: 100%;
  width: 70%;

  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}
</style>