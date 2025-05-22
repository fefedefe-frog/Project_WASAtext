<script>
export default {
  props: ['usr_id'],
  data: function () {
    return {
      errormsg: null,
      loading: false,
      token: "",
      user: {
        usrId: "",
        userName: "",
        userPhoto: "",
      },
    }
  },
  created() {
    this.$watch(
        () => this.$route.params.usr_id,
        (newUsrId, oldUsrId) => {
          this.user['usrId']= newUsrId;
          this.getUserInfo();
        }
    )
  },
  mounted () {
    this.user['usrId']= this.$route.params.usr_id;
    this.token= sessionStorage.getItem('authToken');
    this.loading= true;
    this.getUserInfo();
  },
  methods: {
    sendMessage(rawMessage){
      this.$emit('reqNewChat', {sendTo: this.user['usrId'], messageData: rawMessage});
    },
    async getUserInfo(){
      this.errormsg= null;
      try {
        let response= await this.$axios.get(`/users/${this.user['usrId']}`, {
          headers: {Authorization: this.token},
        });

        if (response.data) {
          this.user= response.data;
        }
      }catch(e) {
        this.errormsg= e;
      }finally {
        this.loading= false
      }
    }
  },
}
</script>

<template>
  <div v-if="!loading" class="info-container">
    <div class="btn-toolbar mb-2 mb-md-0 w-100">
      <div class="btn-group me-2">
        <button type="button" class="btn btn-sm btn-outline-danger" @click="this.$router.replace('/users')">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg> Chiudi
        </button>
        <button type="button" class="btn btn-sm btn-outline-primary shadow-none" @click="getUserInfo">
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

    <span>Invia un messaggio</span>
    <div class="send-message">
      <messageForm @prepMessage="sendMessage"/>
    </div>
    <LoadingSpinner :loading="loading" loadingText="Caricando le info dell'utente..." />
    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>
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
</style>
