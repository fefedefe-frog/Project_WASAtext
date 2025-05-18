<script>
export default {
	props: {
    userData: {
      type: Object,
      required: true
    }
  },
  emits: ['closeUserInfo', 'reqNewChat'],
  data: function () {
    return {
      errormsg: null,
      usrId: "",
      user: {
        usrId: "",
        userName: "",
        userPhoto: "",
      },
    }
  },
  methods: {
    sendMessage(rawMessage){
      this.$emit('reqNewChat', {sendTo: this.user['usrId'], messageData: rawMessage});
    },
  },
  mounted () {
    this.user['usrId']= this.userData['usrId'];
    this.user['userName']= this.userData['userName'];
    this.user['userPhoto']= this.userData['userPhoto'];

    this.usrId= sessionStorage.getItem('usrId');
  }
}
</script>

<template>
  <div class="info-background">
    <div class="info-container">
      <div class="btn-toolbar mb-2 mb-md-0 w-100">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-danger shadow-none" @click="this.$emit('closeUserInfo')">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg> Chiudi
          </button>
          <button type="button" class="btn btn-sm btn-outline-primary shadow-none" @click="this.$emit('closeUserInfo')">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#rotate-cw" /></svg> Ricarica info
          </button>
        </div>
      </div>
      <div class="user-info">
        <div v-if="user['userPhoto']" class="user-image-container">
          <img :src="'data:image/png;base64,'+ user['userPhoto']" alt="Profile Image">
        </div>
        <span>{{ user['userName'].substring(0,12) }}{{ user['userName'].length > 12 ? "..." : "" }}</span>
      </div>

      <div v-if="usrId !== user['usrId']" class="send-message">
        <span>Invia un messaggio</span>
        <div class="message-fomr-container">
          <messageForm @prepMessage="sendMessage"/>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.info-background {
  width: 100%;
  height: 100%;

  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.6);
}

.info-container {
  height: 70%;
  width: 40%;

  display: flex;
  flex-direction: column;

  align-items: center;

  padding: 5px;
  margin-bottom: 5px;

  background-color: lightgray;
  border-radius: 8px;
  box-shadow: 0 2px 5px rgba(0, 0, 0, 0.1);
}

.user-info {
  width: 80%;
  height: fit-content;

  margin-top: 5px;
  border-bottom: 1px grey solid;

  display: flex;
  flex-direction: column;
  align-items: center;


}

.user-image-container {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
}

.user-image-container img {
  width: 20vh;
  height: 20vh;
  border-radius: 50%;
  object-fit: cover;
  user-select: none;
}

.user-info span{
  flex: 1;
  user-select: none;
  font-size: 7vh;
  color: #333;
}


.send-message {
  display: flex;
  flex-direction: column;

  align-items: center;

  width: 80%;
  height: fit-content;
}

.send-message span {
  font-size: 1rem;
  text-align: center;
  margin-bottom: 5px;

}

.message-fomr-container {
  height: 100%;
}
</style>
