<template>
  <div class="signup">
    <h1>Sign up</h1>
    <form>
      <transition name="bounce">
        <p class="errorMsg" v-if="userAlreadyExists">User already exists!</p>
      </transition>
      <input type="text" v-model="firstName" placeholder="First Name" v-on:keyup="inputChange">
      <input type="text" v-model="lastName" placeholder="Last Name" v-on:keyup="inputChange">
      <input type="text" v-model="username" placeholder="Username" v-on:keyup="inputChange">
      <input type="password" v-model="password" placeholder="Password" v-on:keyup="inputChange">
      <button
        class="btn waves-effect waves-light blue accent-4"
        v-bind:class="{disabled: submitDisabled}"
        type="submit"
        v-on:click="signup"
      >Sign Up
        <LoadingSpinner v-if="signingUp"/>
      </button>
    </form>
  </div>
</template>

<script>
import axios from "axios";
import LoadingSpinner from "./LoadingSpinner.vue";

export default {
  name: "Signup",
  components: {
    LoadingSpinner
  },
  data: () => ({
    firstName: "",
    lastName: "",
    username: "",
    password: "",
    submitDisabled: true,
    signingUp: false,
    userAlreadyExists: false
  }),
  methods: {
    signup: function(e) {
      e.preventDefault();

      this.submitDisabled = true;
      this.signingUp = true;

      if (
        this.firstName.trim() !== "" &&
        this.lastName.trim() !== "" &&
        this.username !== "" &&
        this.password !== ""
      ) {
        axios
          .post(
            "/api/signup",
            { data: { firstName: this.firstName, lastName: this.lastName } },
            {
              auth: {
                username: this.username.trim(),
                password: this.password.trim()
              }
            }
          )
          .then(response => {
            localStorage.setItem("jwt", response.data);
            this.signingUp = false;
            this.$emit("user");
          })
          .catch(() => {
            this.signingUp = false;
            this.userAlreadyExists = true;
          });
      }
    },
    inputChange: function() {
      if (this.userAlreadyExists) {
        this.userAlreadyExists = false;
      }

      if (
        this.firstName.trim() !== "" &&
        this.lastName.trim() !== "" &&
        this.username !== "" &&
        this.password !== ""
      ) {
        this.submitDisabled = false;
      } else {
        this.submitDisabled = true;
      }
    }
  }
};
</script>
