<template>
  <div class="login">
    <h1>Log in</h1>
    <form>
      <input type="text" placeholder="Username">
      <input type="text" placeholder="Password">
      <button
        class="btn waves-effect waves-light blue accent-4"
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
  methods: {
    login: function(e) {
      e.preventDefault();

      axios.post("/api/login", {}, { auth: { username: "a", password: "b" } });
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
</style>