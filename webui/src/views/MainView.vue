<script>
export default {
  data: function () {
    return {
      token: '',
      search: '',
      errormsg: null,
      loading: false,
      userChats: [],
      users: []
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.getUsers()
    this.getUserChats()
  },
  methods: {
    async getUsers() {
      this.loading= true
      this.errormsg= null

      try {
        let response= await this.$axios.get(`/users`, {headers: {Authorization: this.token}});
        if (response.data) {
          this.users= []
          response.data['users'].forEach(user => {
            this.users.push(user)
          })
        }
      }catch(e) {
        this.errormsg = e.toString();
      }finally {
        this.loading = false
      }
    },
    async getUserChats() {
      this.loading= true
      this.errormsg= null

      try {
        let response= await this.$axios.get(`/chats`, {headers: {Authorization: this.token}});
        if (response.data) {
          this.userChats= []
          response.data['chats'].forEach(chat => {
            this.userChats.push(chat)
          })
        }
      }catch(e) {
        this.errormsg = e.toString();
      }finally {
        this.loading = false
      }

    }
  }
}
</script>

<template>
  <div class="container">
    <div class="lists">
      <div class="search-box">
        <svg class="feather"><use href="/feather-sprite-v4.29.0.svg#search" /></svg>
        <input v-model="search" type="text" placeholder="Inserisci nome utente o chat" required>
      </div>
      <div v-if="false" class="users-list">
        <userBanner v-for="user in users" :key="user.usrId" :userData="user"/>
      </div>
      <div v-if="true" class="chat-list">
        <chatBanner v-for="chat in userChats" :key="chat.chatId" :chatData="chat"/>
      </div>
    </div>
    <div class="chat-container">

    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>

.container {
  border: 2px cyan solid;

  display: flex;
  flex-direction: row;

  position: relative;
  height: 100%;
  width: 100%;
}

.lists {
  display: flex;
  flex-direction: column;

  height: 100%;
  width: fit-content;
  padding: 5px;

  border: 2px cornflowerblue solid;
  border-radius: 10px;
}

.search-box {
  position: relative;
  width: 100%;
  max-width: 300px;
}

.search-box input {
  width: 100%;
  height: 30px;
  padding: 10px 40px 10px 20px; /* Spazio per l'icona */
  border: 1px solid #ccc;
  border-radius: 50px; /* Arrotonda i bordi */
  outline: none;
}

.search-box svg {
  position: absolute;
  right: 15px;
  top: 50%;
  transform: translateY(-50%);
  color: #888;
}

.users-list {
  border: 2px darkmagenta solid;
  border-radius: 10px;

  height: 100%;
  width: 100%;
  padding: 5px;

  overflow: hidden;
  overflow-y: scroll;
}

.chat-list {
  border: 2px magenta solid;
  border-radius: 10px;

  height: 100%;
  width: 100%;
  padding: 5px;

  overflow: hidden;
  overflow-y: scroll;
}


.chat-container {
  border: 2px peru solid;
  border-radius: 10px;

  height: 100%;
  width: 100%;

}

</style>
