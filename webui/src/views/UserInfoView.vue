<script>
export default {
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

      myUsrId: "",
      updateUserName: false,
      newUserName: "",
    }
  },
  computed: {
    newNameIsValid(){
      return (this.newUserName.length >= 3 && this.newUserName.length <= 16 && ((/^\S.*\S$/).test(this.newUserName)));
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
      let oldProfileImage= this.user['userPhoto'];

      const file= event.target.files[0];

      if (file && file.type.startsWith("image/")) {

        // Faccio la richiesta per modificare l'immagine al backend
        // Preparo il formData per la richiesta
        const requestFormData= new FormData();
        requestFormData.append('newUserPhoto', file);

        this.errormsg= null;
        try{
          let response= await this.$axios.put('/profile/propic', requestFormData, {
            headers: {Authorization: this.token},
          });

          if (response.data){
            this.user['userPhoto']= response.data['userPhoto']
          }
        }catch (e){
          this.errormsg= e.toString();
          this.user['userPhoto']= oldProfileImage;
        }
      }else {
        this.user['userPhoto'] = oldProfileImage;
      }
    },
    enableUsernameUpdate(){
      this.updateUserName= !this.updateUserName;
      let button= this.$refs.enableUsernameUpdate;

      if (this.updateUserName){
        button.classList.add('cancel-edit');
        button.children.item(0).children.item(0).setAttribute("href", "/feather-sprite-v4.29.0.svg#x");

      }else {
        button.classList.remove('cancel-edit');
        button.children.item(0).children.item(0).setAttribute("href", "/feather-sprite-v4.29.0.svg#edit-2");

      }
    },
    async changeUserName(){
      console.log("new name "+ this.newUserName);
      this.errormsg= null;

      try{
        let response= await this.$axios.put(`/profile`, {
          newUserName: this.newUserName,
        },{
          headers: {Authorization: this.token},
        });

        if (response.data){
          this.user['userName']= response.data['userName'];
        }

      }catch (e){
        this.errormsg= e.toString();
      }finally {
        this.updateUserName= false;
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
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#edit" /></svg>
          </span>
        </button>
      </div>
      <div v-else-if="user['userPhoto']" class="user-image-container">
        <img :src="user['userPhoto']" alt="Profile Image" draggable="false">
      </div>
      <div v-if="user['usrId'] === myUsrId" class="username-container">
        <input v-model="newUserName" class="update-username" type="text" :placeholder=" updateUserName ? 'Inserisci nome utente' : user['userName']" :disabled="!updateUserName">

        <button v-if="updateUserName" class="edit-username-button" type="button" :disabled="!newNameIsValid" @click="changeUserName">
          <svg class="feather"> <use href="/feather-sprite-v4.29.0.svg#navigation-2" /></svg>
        </button>

        <button ref="enableUsernameUpdate" class="edit-username-button" type="button" @click="enableUsernameUpdate">
          <svg class="feather"> <use href="/feather-sprite-v4.29.0.svg#edit-2" /></svg>
        </button>
      </div>
      <div v-else class="username-container">
        <span>{{ user['userName'] }}</span>
      </div>
    </div>
    <ErrorMsg v-if="errormsg" :msg="errormsg" @close="errormsg= null" />
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

.username-container{
  flex: 1;
  width: 100%;

  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}

.username-container input {
  width: 70%;
  margin-top: 2px;
  margin-bottom: 2px;
  border-radius: 5%, 5%, 5%, 5%;
  border: none;
  align-items: center;
}

.username-container input:disabled {
  background: none;
}

.username-container span {
  user-select: none;
  font-size: 7vh;
  color: #333;

  max-width: 80%;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.update-username{
  flex: 1;
  user-select: none;
  font-size: 7vh;
  color: #333;

  max-width: 80%;
}

.edit-username-button {
  width: 10%;
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

.edit-username-button:not(:disabled):hover {
  transition: .4s;
  border: 2px dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.edit-username-button:active {
  background-color: lightseagreen;
}

.edit-username-button:disabled{
  background-color: lightgray;
  border: 2px solid lightgray;
  cursor: default;
}

.cancel-edit {
  border: 2px dashed red;
  color: black;
  background-color: red;
}

.cancel-edit:not(:disabled):hover {
  transition: .4s;
  border: 2px dashed red;
  background-color: white;
  color: red;
}

.cancel-edit:active {
  background-color: red;
}
</style>