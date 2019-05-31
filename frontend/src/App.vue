<template>
  <div id="app">
    <div v-if="loggedIn">
      <nav>
        <div class="nav-wrapper indigo darken-2">
          <ul class="left hide-on-med-and-down">
            <li v-on:click="setView" v-bind:class="{active: currView == 'day'}">
              <a href="#">Day</a>
            </li>
            <li v-on:click="setView" v-bind:class="{active: currView == 'week'}">
              <a href="#">Week</a>
            </li>
            <li v-on:click="setView" v-bind:class="{active: currView == 'month'}">
              <a href="#">Month</a>
            </li>
            <li v-on:click="setView" v-bind:class="{active: currView == 'year'}">
              <a href="#">Year</a>
            </li>
          </ul>
          <ul class="right hide-on-med-and-down">
            <li v-on:click="logout">
              <a href="#">Logout</a>
            </li>
          </ul>
        </div>
      </nav>
      <component @user="isValidUser" v-bind:is="currView"></component>
    </div>
    <div v-else>
      <div v-if="!showSignup">
        <Login @user="isValidUser"/>
        <a href="#" v-on:click="signup">New user? Sign up</a>
      </div>
      <Signup v-else @user="isValidUser"/>
    </div>
  </div>
</template>

<script>
import axios from "axios";
import Day from "./components/Day.vue";
import Week from "./components/Week.vue";
import Month from "./components/Month.vue";
import Year from "./components/Year.vue";
import Login from "./components/Login.vue";
import Signup from "./components/Signup.vue";

export default {
  name: "app",
  components: {
    Day,
    Week,
    Month,
    Year,
    Login,
    Signup
  },
  data: () => ({
    currView: "month",
    loggedIn: false,
    showSignup: false
  }),
  methods: {
    setView: function(event) {
      this.currView = event.target.text.toLowerCase();
    },
    logout: function(event) {
      axios
        .post("/api/logout")
        .then(response => {
          this.loggedIn = false;
        })
        .catch(err => {
          console.log(err);
        });
    },
    signup: function(event) {
      event.preventDefault();

      this.showSignup = true;
    },
    isValidUser: function() {
      axios
        .get("/api/isValidUser")
        .then(response => {
          this.loggedIn = true;
        })
        .catch(err => {
          this.loggedIn = false;
        });
    }
  },
  created: function() {
    this.isValidUser();
  }
};
</script>

<style>
#app {
  font-family: 'Nunito', sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}

</style>
