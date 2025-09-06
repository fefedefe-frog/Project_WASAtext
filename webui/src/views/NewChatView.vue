<script>
export default {
  data: function () {
    return {
      token: '',
      usrId: '',
      userName: '',
      errormsg: null,
      users: [],

      chatInfo: {
        chatName: "",
        chatPhoto: null,
        chatPhotoPreview: "",
        isGroup: false,
        participantsId: [],
      },


      selectedParticipantInfo: [],

      initialMessage: {
        textContent: "",
        photoContent: null,
        photoName: ""
      },

      setIntervalId: null,
    }
  },
  computed: {
    chatNameIsValid() {
      const name=  this.chatInfo['chatName'];
      return (name.length >= 3 && name.length <= 16 && ((/^\S.*\S$/).test(name)));
    },
    userFilteredResult(){
      let searchPool= this.users;
      let query= this.searchQuery.trim();

      if (query === ''){
        return searchPool;
      }
      return searchPool.filter(user =>
        user['userName'].toLowerCase().includes(this.searchQuery.toLowerCase())
      );
    },
    isFormSendable() {
      let messageValid= this.initialMessage['textContent'].trim() || this.initialMessage['photoContent'];
      let chatDataValid= this.chatNameIsValid && (this.chatInfo['participantsId'].length >= 1);
      // return messageValid && chatDataValid
      return chatDataValid
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.usrId= sessionStorage.getItem('usrId');
    this.userName= sessionStorage.getItem('userName');

    this.getUsers();
    this.setIntervalId= setInterval(async () => {
      await this.getUsers();
    }, 5000);
  },
  beforeUnmount() {
    clearInterval(this.setIntervalId);
  },
  methods: {
    async getUsers() {
      this.errormsg= null;

      try {
        let response= await this.$axios.get(`/users`, {
          headers: {Authorization: this.token}
        });

        if (response.data) {
          if(response.data['users']){
            this.users= [];
            response.data['users'].forEach(user => {
              this.users.push(user);
            });
          }
        }
      }catch(e) {
        if (e.response.status === 404){
          this.users= [];
        }else {
          let error_string= ""
          if (e.response.status === 400 ||  //Bad request
              e.response.status === 401 ||  //Unauthorized
              e.response.status === 403 ||  //Forbidden
              e.response.status === 500){   //Internal server error
            error_string= `Error: ${e.response.status}. ${e.response.data}`;
          }else{  //Axios error
            error_string= `Internal axios error: ${e}`;
            console.log(e);
          }
          this.errormsg= error_string;
        }
      }
    },
    async startNewChat(){
      this.errormsg= null;

      if (this.chatInfo['participantsId'].length < 1){
        this.errormsg= new Error("numero di partecipanti non valido");
        return
      }

      // Preparo il formData per la richiesta
      const requestFormData= new FormData();

      // Assegno le informazioni sulla chat
      let chatName= "";
      let chatPhoto= new Blob([], {type: 'image/png'});
      if (this.chatInfo['isGroup']){
        chatName= this.chatInfo['chatName'];
        if (this.chatInfo['chatPhoto']){
          chatPhoto= this.chatInfo['chatPhoto'];
        }
      }
      requestFormData.append('chatName', chatName);
      requestFormData.append('chatPhoto', chatPhoto);
      requestFormData.append('isGroup', this.chatInfo['isGroup']);
      requestFormData.append('participants', this.chatInfo['participantsId']);

      // Assegno le informazioni sul messaggio
      let textContent= this.initialMessage['textContent'];
      let photoContent= new Blob([], {type: 'image/png'});

      // Come avevo pensato di fare
      // if (this.initialMessage['photoContent']){
      //   textContent= "";
      //   photoContent= this.initialMessage['photoContent'];
      // }
      //
      // requestFormData.append('messageTextContent', textContent);
      // requestFormData.append('messagePhotoContent', photoContent);

      // Test per vedere se la chat può essere creata senza messaggio dell'utente
      textContent= `Chat creata da ${this.userName}`
      requestFormData.append('messageTextContent', textContent);
      requestFormData.append('messagePhotoContent', photoContent);

      try{
        let response= await this.$axios.post(`/chats`, requestFormData, {
          headers: {
            Authorization: this.token,
          }
        });

        if(response.data){
          let newChat= response.data
          setTimeout(()=> {
            this.$router.push(`/chats/${newChat['chatId']}`);
          }, 500);
        }
      }catch(e) {
        let error_string= ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 500){   //Internal server error
          error_string= `Error: ${e.response.status}. ${e.response.data}`;
        }else{  //Axios error
          error_string= `Internal axios error: ${e}`;
          console.log(e);
        }
        this.errormsg= error_string;
      }
    },
    imageUpload(target) {
      const input= document.createElement('input');
      input.type= "file";
      input.accept= "image/*";

      input.addEventListener("change", (event) => {

        const file= event.target.files[0];
        if (file && file.type.startsWith("image/")) {
          const reader= new FileReader();

          reader.onload= (e) => {
            if (target === 'chatPhoto'){
              this.chatInfo['chatPhoto']= file;
              this.chatInfo['chatPhotoPreview']= e.target.result;
            }else if (target === 'messagePhoto'){
              this.initialMessage['textContent']= "";
              this.initialMessage['photoContent']= file;
              this.initialMessage['photoName']= file.name;
            }
          };
          reader.readAsDataURL(file);
        }else {
          if (target === 'chatPhoto') {
            this.chatInfo['chatPhoto'] = null;
            this.chatInfo['chatPhotoPreview'] = "";
          } else if (target === 'messagePhoto') {
            this.initialMessage['photoContent'] = null;
            this.initialMessage['photoName'] = "";
          }
        }

      });
      input.click();
    },
    addParticipant(bannerData){
      if(!this.chatInfo['isGroup']){
        this.chatInfo['participantsId']= [];
        this.selectedParticipantInfo= [];

        this.chatInfo['chatPhotoPreview']= `data:image/png;base64,${bannerData['userPhoto']}`;
        this.chatInfo['chatName']= bannerData['userName'];
      }

      // Aggiungo l'utente alla lista di utenti, solo se non è già presente
      if (!this.chatInfo['participantsId'].some(usrId => usrId === bannerData['usrId'])){
        this.selectedParticipantInfo.push(bannerData);
        this.chatInfo['participantsId'].push(bannerData['usrId']);
      }
    },
    removeParticipant(id){
      this.selectedParticipantInfo= this.selectedParticipantInfo.filter(item =>
        item['usrId'] !== id
      );
      this.chatInfo['participantsId']= this.chatInfo['participantsId'].filter(pId =>
          pId !== id
      );
    },
    clearFormData(){
      this.chatInfo['chatName']= "";
      this.chatInfo['chatPhoto']= null;
      this.chatInfo['chatPhotoPreview']= "";
      this.chatInfo['participantsId']= [];
      this.selectedParticipantInfo= [];

      this.initialMessage['textContent']= "";
      this.initialMessage['photoContent']= null;
      this.initialMessage['photoName']= "";
    },
    componentsErrorHandler(e){
      this.errormsg= e;
    },
  }
}
</script>

<template>
  <div class="main-container bobby">
    <div class="select-participant">
      <sidebarList :banner-component="'userBanner'" items="users" @error="componentsErrorHandler" @banner-data="addParticipant" />
    </div>
    <form class="new-chat-form" @submit.prevent="startNewChat">
      <!-- Sezione per le info della chat -->
      <div class="new-chat-info">
        <div class="new-chat-image">
          <button class="chat-image-button" type="button" :disabled="!chatInfo['isGroup']" @click="imageUpload('chatPhoto')">
            <img v-if="chatInfo['chatPhotoPreview'] || chatInfo['isGroup']" :src=" chatInfo['chatPhotoPreview'] || '/images/def_group.png'" alt="Anteprima" draggable="false">
            <span>
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
            </span>
          </button>
        </div>
        <div class="new-chat-text">
          <div class="new-chat-name">
            <span>
              <label for="chatName">Nome della chat:  <span v-if="!chatNameIsValid && chatInfo['chatName']" style="color: red; margin-top: 10px;">Nome non valido</span></label>
              <input id="chatName" v-model="chatInfo['chatName']" type="text" :placeholder="chatInfo['isGroup'] ? 'Inserisci il nome' : ''" :disabled="!chatInfo['isGroup']" required>
            </span>
            <span>
              <label for="isGroup">Nuovo gruppo</label>
              <input id="isGroup" v-model="chatInfo['isGroup']" type="checkbox" @click="clearFormData">
            </span>
          </div>
          <!-- Sezione per la visualizzazione dei partecipanti di un gruppo -->
          <div v-if="chatInfo['isGroup']" class="participants">
            <span v-for="p in selectedParticipantInfo" :key="p['usrId']" class="participant">
              <img :src="'data:image/png;base64,'+ p['userPhoto']" alt="Profile Image" draggable="false">
              <span>{{ p['userName'] }}</span>
              <button @click="removeParticipant(p['usrId'])">
                <svg class="feather"> <use href="/feather-sprite-v4.29.0.svg#x" /></svg>
              </button>
            </span>
          </div>
        </div>
      </div>
      <!-- Sezione per il messaggio iniziale (COME DOVREBBE ESSERE) -->
      <!--
      <div class="initial-message">
        <div class="mess-form">
          <textarea v-if="!initialMessage['photoContent']" v-model="initialMessage['textContent']" class="textarea-content" placeholder="Scrivi un messaggio..." rows="2" :maxlength="1024" />
          <div v-if="initialMessage['photoContent']" class="image-name">
            <button class="form-buttons delete-button" type="button" @click="initialMessage['photoContent']=null; initialMessage['photoName']= ''">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
            </button>
            <span>{{ initialMessage['photoName'] }}</span>
          </div>
          <div class="buttom-column">
            <button class="form-buttons" type="button" @click="imageUpload('messagePhoto')">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
            </button>

            <button class="form-buttons" type="submit" :disabled="!isFormSendable">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigation" /></svg>
            </button>
          </div>
        </div>
      </div>
      -->
      <!-- TEST -->
      <div class="initial-message">
        <div class="mess-form">
          <button class="form-buttons-test" type="submit" :disabled="!isFormSendable">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#arrow-up" /></svg>
          </button>
        </div>
      </div>
    </form>
    <ErrorMsg v-if="errormsg" :msg="errormsg" @close="errormsg= null" />
  </div>
</template>

<style scoped>
.main-container {
  height: 100%;
  width: 100%;

  padding: 10px;
  display: flex;
  flex-direction: row;
}

@media (min-width: 2000px) {
  .main-container {
    max-width: 1400px;
    max-height: calc(1400px * 9 / 16);
  }
}

.select-participant {
  display: flex;
  flex-direction: column;
  align-items: center;

  width: 25%;
  padding: 4px;
}

.new-chat-form{
  width: 75%;
  height: 100%;
}

.new-chat-info {
  width: 75%;
  height: 50%;

  padding: 10px;

  display: flex;
  flex-direction: row;
  align-items: center;
}

.new-chat-image {
  height: 70%;
  aspect-ratio: 1/1;

  display: flex;
  justify-content: center;
  align-items: center;
}

.chat-image-button {
  height: 90%;
  aspect-ratio: 1/1;
  border-radius: 50%;

  padding: 0;
  border: none;
  background: lightgray;
  position: relative;
  overflow: hidden;
}

.chat-image-button img {
  display: block;

  width: 100%;
  height: 100%;


  object-fit: cover;
  object-position: center;
  pointer-events: none;
}

.chat-image-button:not(:disabled):hover {
  filter: brightness(0.6);
  transition: filter 0.2s ease-in-out;
}

.chat-image-button span {
  position: absolute;
  width: 40%;
  height: 40%;

  border-radius: 50%;
  padding: 10px;

  top: 50%;
  left:50%;
  transform: translate(-50%, -50%);

  color: white;
  background: rgba(0, 0, 0, 0.5);
}

.chat-image-button svg {
  width: 100%;
  height: 100%;
}


.new-chat-text {
  height: 70%;
  width: 70%;

  margin-left: 5px;


  display: flex;
  flex-direction: column;
}

.new-chat-name {
  width: fit-content;
  max-width: 80%;
  height: 60%;

  display: flex;
  flex-direction: column;

}

#chatName, #isGroup {
  margin-left: 5px;
}


.participants {
  width: 100%;
  height: 40%;

  padding-left: 5px;
  padding-right: 5px;
  gap: 5px;

  display: flex;
  flex-direction: row;
  align-items: center;

  overflow: hidden;
  overflow-x: scroll;
}

.participant {
  width: fit-content;
  height: fit-content;

  max-width: 100px;

  padding: 3px;

  background: skyblue;
  border-radius: 10px;

  display: flex;
  flex-direction: row;
  align-items: center;
}

.participant:hover {
  background: dodgerblue;
}

.participant img {
  height: 20px;
  width: 20px;

  border-radius: 10px;
  user-select: none;
}

.participant span{
  font-size: 14px;
  text-align: center;

  margin-left: 2px;

  user-select: none;

  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.participant button {
  background: none;
  border: none;
  width: 20px;
  height: 20px;
  border-radius: 20px;

  cursor: pointer;

  display: inline-flex;
  align-items: center;
  justify-content: center;

  color: red;
}

.participant svg {
  stroke-width: 3;
}

.participant button:hover {
  background: rgba(0, 0, 0, 0.5);
}

/* message form */
.initial-message {
  width: 75%;
  height: 50%;

  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.mess-form {
  width: 50%;
  height: 30%;

  display: flex;
  flex-direction: row;
  justify-content: right;

  /* TEST */
  justify-content: center;
}

.textarea-content {
  width: 100%;
  max-height: 100px;
  margin: 2px;

  border-radius: 10px;
  padding: 2px 4px 0 4px;

  resize: none;
  overflow-x: hidden;
  font-size: 0.8rem;
  line-height: 1rem;
  box-sizing: border-box;
  border: 1px solid white;
}

.buttom-column {
  width: fit-content;
  height: fit-content;
  margin-top: auto;
  margin-left: 0.3vh;
  display: flex;
  flex-direction: column;
}

/* COME DOVREBBE ESSERE
.form-buttons {
  width: 25px;
  height: 25px;

  margin: 1px 5px 2px 0;

  border-radius: 25%;
  border: 2px dashed lightseagreen;


  color: white;
  background-color: lightseagreen;
  cursor: pointer;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
  transition: .4s;

  display: flex;
  justify-content: center;
  align-items: center;
}

.form-buttons:not(:disabled):hover {
  transition: .4s;
  border: 2px dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.form-buttons:active {
  background-color: lightseagreen;
}

.form-buttons:disabled{
  background-color: lightgray;
  border: 2px solid lightgray;
  cursor: default;
}
*/

/* TEST */
.form-buttons-test {
  width: 30%;
  aspect-ratio: 1/1;


  border-radius: 25%;
  border: 2px dashed lightseagreen;


  color: white;
  background-color: lightseagreen;
  cursor: pointer;

  box-shadow: rgba(50, 50, 93, 0.25) 0px 2px 5px -1px, rgba(0, 0, 0, 0.3) 0px 1px 3px -1px;
  transition: .4s;

  display: flex;
  justify-content: center;
  align-items: center;
}

.form-buttons-test:not(:disabled):hover {
  transition: .4s;
  border: 2px dashed lightseagreen;
  background-color: white;
  color: lightseagreen;
}

.form-buttons-test:active {
  background-color: lightseagreen;
}

.form-buttons-test:disabled{
  background-color: lightgray;
  border: 2px solid lightgray;
  cursor: default;
}

.delete-button {
  border: 2px dashed red;
  background: red;
}
.delete-button:hover {
  border: 2px dashed darkred;
  background-color: white;
  color: red;
}
.delete-button:active {
  border: 2px dashed darkred;
}

.image-name {
  height: fit-content;
  object-fit: contain;

  margin-right: 2px;

  display: flex;
  flex-direction: row;

  align-items: center;
  justify-content: center;

}
</style>