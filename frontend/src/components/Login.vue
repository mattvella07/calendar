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
      >
        Log in
        <LoadingSpinner v-if="loggingIn"/>
      </button>
    </form>
  </div>
</template>

<script>
import axios from "axios";
import LoadingSpinner from "./LoadingSpinner.vue";

export default {
  name: "Login",
  components: {
    LoadingSpinner
  },
  data: () => ({
    username: "",
    password: "",
    submitDisabled: true,
    loggingIn: false,
    invalidCreds: false
  }),
  methods: {
    login: function(e) {
      e.preventDefault();

      this.submitDisabled = true;
      this.loggingIn = true;

      if (this.username.trim() !== "" && this.password.trim() !== "") {
        axios
          .post(
            "/api/login",
            {},
            {
              auth: {
                username: this.username.trim(),
                password: this.password.trim()
              }
            }
          )
          .then(response => {
            this.loggingIn = false;
            this.$emit("user");
          })
          .catch(err => {
            console.log(err);
            this.loggingIn = false;
            this.invalidCreds = true;
          });
      }
    },
    inputChange: function() {
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
