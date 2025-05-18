<script>
export default {
	props: {
    participantIds: {
      type: Array,
      required: true
    }
  },
  emits: ['closeCreateGroup', 'newGroupData', 'error'],
  data: function () {
    return {
      usrId: "",
      groupData: {
        groupName: "",
        groupPhoto: null,
      },
      firstMessage: {
        textContent: "",
        photoContent: null
      },
      imageName: "",
      maxLength: 1024
    }
  },
  methods: {
    async createGroup(){
      let requestFormData= FormData.new()

      // Assegno le informazioni sul messaggio
      requestFormData.append('textContent', this.firstMessage['textContent']);
      if (this.firstMessage['photoContent'] === null){
        this.firstMessage['photoContent']= new Blob([], {type: 'image/png'});
      }
      requestFormData.append('photoContent', this.firstMessage['photoContent']);
      requestFormData.append('respondTo', -1);

      // Assegno le informazioni sulla chat
      requestFormData.append('chatName', this.groupData['groupName']);
      if (this.groupData['groupPhoto'] === null){
        this.groupData['groupPhoto']= new Blob([], {type: 'image/png'});
      }
      requestFormData.append('chatPhoto', this.groupData['groupPhoto']);
      requestFormData.append('isGroup', true);
      requestFormData.append('participants', this.participantIds);

      try{
        let response= await this.$axios.post(`/chats`, requestFormData, {headers: {Authorization: this.token}});

        if (response.data) {
          if(response.data){
            let groupData= response.data
            this.$emit('newGroupData', groupData);
          }
        }
      }catch (e){
        this.$emit('error', e);
      }
    },
    prepMessage() {
      let emptyPhoto= new Blob([], {type: 'image/png'});
      if (this.firstMessage['textContent'] || this.firstMessage['photoContent']){
        if (this.firstMessage['textContent'].trim()){
          this.firstMessage['textContent']= "";
        }
        if (this.firstMessage['photoContent'].trim()){
          this.firstMessage['photoContent']= emptyPhoto;
        }


        this.firstMessage['textContent']= "";
        this.firstMessage['photoContent']= null;
        this.imageName= "";
      }
    },
    imageUpload() {
      const input= document.createElement('input');
      input.type= "file";
      input.accept= "image/*";
      input.addEventListener("change", this.handleFileChange);
      input.click();
    },
    handleFileChange(event) {
      const file= event.target.files[0];
      if(file && file.type.startsWith("image/")) {
        this.firstMessage['photoContent']= file;
        this.imageName= file.name;
      }else {
        this.firstMessage['photoContent']= null;
        this.imageName= "";
      }
    }
  },
  mounted () {
    this.token= sessionStorage.getItem('authToken')
    this.usrId= sessionStorage.getItem('usrId');
  }
}
</script>

<template>
  <div class="group-background">
    <div class="group-container">
      <div class="btn-toolbar mb-2 mb-md-0 w-100">
        <div class="btn-group me-2">
          <button type="button" class="btn btn-sm btn-outline-danger shadow-none" @click="this.$emit('closeCreateGroup')">
            <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg> Chiudi
          </button>
        </div>
      </div>
      <!-- TODO  caperi come strutturare qua -->

      <form class="create-group-form" @submit.prevent="prepMessage">
        <div class="groupName-section">
          <label for="groupName">Nome Gruppo:  <!--<span v-if="!isUsernameValid && groupName" class="error">Nome non valido</span>--></label>
          <input id="groupName" v-model="groupData['groupName']" type="text" placeholder="Inserisci il nome del gruppo" required>
        </div>

        <div class="message-section">
          <textarea v-if="!firstMessage['photoContent']" class="textarea-content" v-model="firstMessage['photoContent']" ref="textareaMessage" placeholder="Scrivi un messaggio..." rows="2" :maxlength="maxLength"></textarea>
          <div v-if="firstMessage['photoContent']" class="image-name">
            <button class="form-buttons delete-button" type="button" @click="firstMessage['photoContent']=null; imageName= '';">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
            </button>
            <span>{{ imageName }}</span>
          </div>

          <div class="buttom-column">
            <button class="form-buttons" type="button" @click="imageUpload">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#image" /></svg>
            </button>

            <button class="form-buttons" type="submit">
              <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#navigation" /></svg>
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.group-background {
  width: 100%;
  height: 100%;

  display: flex;
  justify-content: center;
  align-items: center;
  background: rgba(0, 0, 0, 0.6);
}

.group-container {
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

.create-group-form{
  display: flex;
  flex-direction: column;
  align-items: center;
}

.groupName-section label{
  margin: 0 5px 5px 0;
}

.message-section {
  display: flex;
  flex-direction: row;
}

.buttom-column{
  display: flex;
  flex-direction: column;
}
</style>
