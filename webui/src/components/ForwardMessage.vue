<script>
import SidebarList from "./SidebarList.vue";

export default {
  components: {SidebarList},
  props: {
    chatId: {
      type: Number,
      required: true
    },
    msgId: {
      type: String,
      required: true
    },
  },
  emits: ['close', 'error'],
  data: function () {
    return{
      status: "minus",
      date: "",
      textContent: "",
      photoContent: {},
      usrId: "",

      message: {
        msgId: -1,
        senderId: "",
        respondTo: -1,
        textContent: "",
        photoContent: [],
        timestamp: "",
        comments: [],
      },

      respondToData: {}
    }
  },
  mounted(){
    this.usrId= sessionStorage.getItem('usrId')
    this.token= sessionStorage.getItem('authToken');
  },
  methods: {
    async forwardToChat(chatId){
      try {
        await this.$axios.post(`/chats/${this.chatId}/messages/${this.msgId}`, {
          chatToForward: chatId
        },{
          headers: {Authorization: `${this.token}`}
        });
      } catch(e) {
        let error_string = ""
        if (e.response.status === 400 ||  //Bad request
            e.response.status === 401 ||  //Unauthorized
            e.response.status === 403 ||  //Forbidden
            e.response.status === 404 ||  //Not found
            e.response.status === 500) {   //Internal server error
          error_string = `Error: ${e.response.status}. ${e.response.data}`;
        } else {  //Axios error
          error_string = `Internal axios error: ${e}`;
          console.log(e);
        }
        this.$emit('error' ,error_string);
      }
      this.$emit('close');
    }
  }
}
</script>

<template>
  <div class="forward-main">
    <div class="close-button">
      <button type="button" class="btn btn-sm btn-danger shadow-none" @click="$emit('close')">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#x" /></svg>
      </button>
    </div>
    <div class="list-to-forward">
      <sidebar-list banner-component="ChatBanner" items="chats" @banner-data="forwardToChat" @error="$emit('error')" />
    </div>
  </div>
</template>

<style scoped>
.forward-main{
  width: 100%;
  height: 100%;

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.close-button{
  border: 1px solid black;
  width: 100%;
  height: 15%;

  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}

.list-to-forward{
  border: 1px solid green;
  width: 80%;
  height: 85%;
  overflow: hidden;
}

.btn{
  height: 60%;
  aspect-ratio: 1/1;

  display: flex;
  align-items: center;
  justify-content: center;
}

</style>