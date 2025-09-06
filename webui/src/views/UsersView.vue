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
    showUserInfo(bannerData){
      this.$router.replace(`/users/${bannerData['usrId']}`);
    },
    componentsErrorHandler(e){
      this.errormsg= e;
    },
  }
}
</script>

<template>
  <div class="main-container">
    <div class="lists bobby">
      <SidebarList items="users" :banner-component="'UserBanner'" @error="componentsErrorHandler" @banner-data="showUserInfo" />
    </div>
    <div class="user-info-container bobby">
      <RouterView :key="$route.fullPath" />
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" @close="errormsg= null" />
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

.user-info-container {
  height: 100%;
  width: 75%;

  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
}
</style>