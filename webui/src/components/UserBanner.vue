<script>
export default {
  props: {
    inputData: {
      type: Object,
      required: true,
    }
  },
  emits: ['bannerClicked'],
  data() {
    return {
      user: {
        usrId: '',
        userName: '',
        userPhoto: '',
      }
    }
  },
  mounted() {
    this.user['usrId']= this.inputData['usrId'];
    this.user['userName']= this.inputData['userName'];
    this.user['userPhoto']= this.inputData['userPhoto'];
  },
  methods: {
    loadUser(){
      this.$emit("bannerClicked", {usrId: this.user['usrId']});
    }
  }
};
</script>

<template>
  <div class="user-banner" @click="loadUser">
    <!-- Foto Profilo a destra -->
    <div v-if="user['userPhoto']" class="user-image-container">
      <img :src="'data:image/png;base64,'+ user['userPhoto']" alt="Profile Image">
    </div>

    <!-- Contenuto del banner -->
    <div class="username-container">
      <span>{{ user['userName'] }}</span>
    </div>
  </div>
</template>

<style scoped>

.user-banner {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 5px;
  margin-bottom: 5px;
  background-color: lightgray;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
  height: 70px;
  width: 100%;
}

.user-banner:hover {
  background-color: darkgray;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.4);
}

.user-image-container {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.user-image-container img {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.username-container {
  flex: 1;
  padding: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  margin-left: 10px;
  user-select: none;
}

.username-container span{
  width: 100%;
  font-size: 1.1em;
  font-weight: bold;
  color: #333;
}
</style>