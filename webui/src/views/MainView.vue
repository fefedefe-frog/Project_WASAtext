<script>
export default {
  data: function () {
    return {
      token: '',
      errormsg: null,
      loading: false,
      userChats: [],
      users: []
    }
  },
  mounted() {
    this.token= sessionStorage.getItem('authToken');
    this.users()
  },
  methods: {
    async users() {
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
  }
}
</script>

<template>
  <div class="container">
    <div class="lists">
      <user-banner v-for="user in users" :key="user.usrId" :userData="user"/>
    </div>

    <ErrorMsg v-if="errormsg" :msg="errormsg" />
  </div>
</template>

<style scoped>

</style>
