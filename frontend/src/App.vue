<template>
  <div id="app">
    <div v-if="loggedIn">
      <nav>
        <div class="nav-wrapper blue accent-4">
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
        </div>
      </nav>
      <component v-bind:is="currView"></component>
    </div>
    <div v-else>
      <Login/>
    </div>
  </div>
</template>

<script>
import Day from "./components/Day.vue";
import Week from "./components/Week.vue";
import Month from "./components/Month.vue";
import Year from "./components/Year.vue";
import Login from "./components/Login.vue";

export default {
  name: "app",
  components: {
    Day,
    Week,
    Month,
    Year,
    Login
  },
  data: () => ({
    currView: "month",
    loggedIn: false
  }),
  methods: {
    setView: function(event) {
      this.currView = event.target.text.toLowerCase();
    }
  }
};
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
