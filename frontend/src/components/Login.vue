<template>
  <div class="login">
    <h1>Log in</h1>
    <form>
      <transition name="bounce">
        <p class="errorMsg" v-if="invalidCreds">Invalid username and/or password!</p>
      </transition>
      <input type="text" v-model="username" placeholder="Username" v-on:keyup="inputChange">
      <input type="password" v-model="password" placeholder="Password" v-on:keyup="inputChange">
      <button
        class="btn waves-effect waves-light blue accent-4"
        v-bind:class="{disabled: submitDisabled}"
        type="submit"
        v-on:click="login"
      >Log in</button>
    </form>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Login",
  data: () => ({
    username: '',
    password: '',
    submitDisabled: true,
    invalidCreds: false
  }),
  methods: {
    login: function(e) {
      e.preventDefault();

      if (this.username.trim() !== "" && this.password.trim() !== "") {
        axios.post("/api/login", {}, { auth: { username: this.username.trim(), password: this.password.trim() } })
        .then(response => {
          console.log('Success')
          localStorage.setItem("jwt", response.data)
          location.reload();
        })
        .catch(err => {
          console.log("ERR: " + err);
          this.invalidCreds = true;
        })
      }
    },
    inputChange: function(e) {
      if (this.invalidCreds) {
        this.invalidCreds = false;
      }

      if (this.username.trim() !== "" && this.password.trim() !== "") {
        this.submitDisabled = false;
      } else {
        this.submitDisabled = true;
      }
    }
  }
};
</script>

<style scoped>
form {
  width: 300px;
  height: 300px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin: auto;
}
.errorMsg {
  color: white;
  background-color: red;
  border-radius: 4px;
  padding: 4px;
}
.bounce-enter-active {
  animation: bounce-in .5s;
}
@keyframes bounce-in {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.5);
  }
  100% {
    transform: scale(1);
  }
}
</style>