<script>
export default {
  props: {chatId: String, isVisible: Boolean},
  emits: ["error"],
  data: function () {
    return{
      chatName: '',
      chatPhoto: '',
      participants: [],
      token: '',
    }
  },
  created () {
    this.token= localStorage.getItem('authToken').split(' ')[1];
  },
  mounted () {
    this.token= localStorage.getItem('authToken').split(' ')[1];
  },
  methods: {
    async refreshChatInfo() {
      this.loading= true
      this.errormsg= null
      this.messages=[]

      try {
        let response = await this.$axios.get('/chats/'+ this.chatId, {headers: {Authorization: `${this.token}`}});
        if (response.data) {
          this.chatName= response.data.chatName
          this.chatPhoto= response.data.chatPhoto
          this.participants= response.data.participants
        }
      } catch (e) {
          this.$emit('error', e)
      }finally {
        this.loading = false;
      }
    }
  }
};
</script>

<template>
  <div v-if="isVisible" class="info-sidebar">
    <div class="info-sidebar-header">
      <h3>Chat Info</h3>
      <p>Foto Chat modificabile</p>
    </div>
    <div class="info-sidebar-content">
      <p>Partecipanti, con nome</p>
      <p>Aggiungi partecipanti</p>
    </div>
  </div>
</template>

<style scoped>
.info-sidebar {
  position: absolute;
  top: 0;
  right: 0;
  width: 300px;
  height: 400px;
  background-color: #f8f9fa;
  box-shadow: -2px 0 5px rgba(0, 0, 0, 0.5);
  border-radius: 10px;
  z-index: 1000;
  padding: 20px;
  display: flex;
  flex-direction: column;
}

.info-sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ddd;
  padding-bottom: 10px;
  margin-bottom: 20px;
}

.close-info-sidebar-btn {
  background-color: #dc3545;
  color: white;
  border: none;
  padding: 5px 10px;
  cursor: pointer;
  border-radius: 4px;
}

.close-info-sidebar-btn:hover {
  background-color: #c82333;
}

.info-sidebar-content {
  flex-grow: 1;
  overflow-y: auto;
  color: #333;
  font-size: 14px;
}
</style>