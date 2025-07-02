<script>
export default {
  data: function () {
    return {
      errormsg: null,
      loading: false,
      token: "",
      myUsrId: "",
      user: {
        usrId: "",
        userName: "",
        userPhoto: "",
      },
    }
  },
  async mounted () {
    this.user['usrId']= this.$route.params.usr_id;

    this.token= sessionStorage.getItem('authToken');
    this.myUsrId= sessionStorage.getItem('usrId');

    this.loading= true;
    await this.getUserInfo();
  },
  methods: {
    sendMessage(rawMessage){
      console.log(rawMessage)
    },
    async getUserInfo(){
      this.errormsg= null;
      try {
        let response= await this.$axios.get(`/users/${this.user['usrId']}`, {
          headers: {Authorization: this.token},
        });

        if (response.data) {
          this.user= response.data;
          this.user['userPhoto']= `data:image/png;base64,${this.user['userPhoto']}`;
        }
      }catch(e) {
        this.errormsg= e.toString();
      }finally {
        this.loading= false
      }
    },
    imageUpload() {
      const input= document.createElement('input');
      input.type= "file";
      input.accept= "image/*";

      input.addEventListener("change", async (event) => await this.changeUserPhoto(event));
      input.click();
    },
    async changeUserPhoto(event){
      let newProfileImage= null;
      let oldProfileImage= this.user['userPhoto'];

      const file= event.target.files[0];

      if (file && file.type.startsWith("image/")) {
        const reader= new FileReader();

        reader.onload= (e) => {
          newProfileImage= file;
          this.user['userPhoto']= e.target.result;
        };
        reader.readAsDataURL(file);

        // Faccio la richiesta per modificare l'immagine al backend
        // Preparo il formData per la richiesta
        const requestFormData= new FormData();
        requestFormData.append('newUserPhoto', newProfileImage);

        this.errormsg= null;
        try{
          let response= await this.$axios.put('/profile/propic', requestFormData, {
            headers: {Authorization: this.token},
          });
        }catch (e){
          this.errormsg= e.toString();
          this.user['userPhoto']= oldProfileImage;
        }
      }else {
        this.user['userPhoto'] = oldProfileImage;
      }
    }
  },
}
</script>

<template>
  <LoadingSpinner :loading="loading" loading-text="Caricando le info dell'utente..." />
  <div v-if="!loading" class="info-container">
    <div class="btn-toolbar mb-2 mb-md-0 w-100">
      <div class="btn-group me-2">
        <button type="button" class="btn btn-sm btn-outline-danger" @click="$router.replace('/users')">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg> Chiudi
        </button>
        <button type="button" class="btn btn-sm btn-outline-primary shadow-none" @click="getUserInfo">
          <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#rotate-cw" /></svg> Ricarica info
        </button>
      </div>
    </div>
    <div class="user-info">
      <div v-if="user['usrId'] === myUsrId" class="user-image-container">
        <button class="user-image-button" type="button" @click="imageUpload()">
          <img v-if="user['userPhoto']" :src="`${user['userPhoto']}` || '/images/def_single.png'" alt="Anteprima" draggable="false">
          <span>
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
            </span>
        </button>
      </div>
      <div v-else-if="user['userPhoto']" class="user-image-container">
        <img :src="'data:image/png;base64,'+ user['userPhoto']" alt="Profile Image">
      </div>
      <span>{{ user['userName'].substring(0,12) }}{{ user['userName'].length > 12 ? "..." : "" }}</span>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" @close="this.errormsg= null" />
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
  aspect-ratio: 1/1;

  display: flex;
  justify-content: center;
  align-items: center;
}

.user-image-button {
  height: 100%;
  aspect-ratio: 1/1;
  border-radius: 50%;


  padding: 0;
  border: none;
  background: lightgray;

  position: relative;
}

.user-image-button img {
  display: block;

  width: 100%;
  height: 100%;


  object-fit: cover;
  object-position: center;
  pointer-events: none;
}

.user-image-button:hover {
  filter: brightness(0.6);
  transition: filter 0.2s ease-in-out;
}

.user-image-button span {
  position: absolute;
  width: 40%;
  height: 40%;

  border-radius: 50%;
  padding: 10px;

  display: flex;

  top: 50%;
  left:50%;
  transform: translate(-50%, -50%);

  background: rgba(0, 0, 0, 0.5);
}

.user-image-button svg {
  width: 100%;
  height: 100%;
  color: white;
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
