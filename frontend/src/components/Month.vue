<template>
  <div class="hello">
    <!-- <h1>{{ msg }}</h1> -->
    <div class="title">
      <a
        v-on:click="prevMonth"
        class="btn-floating btn-large waves-effect waves-light blue accent-4"
      >
        <i class="material-icons"><</i>
      </a>
      <h1>{{ months[currMonth] }} {{ currYear }}</h1>
      <a
        v-on:click="nextMonth"
        class="btn-floating btn-large waves-effect waves-light blue accent-4"
      >
        <i class="material-icons">></i>
      </a>
    </div>
    <table>
      <tr>
        <th v-for="(d, dKey) in dow" v-bind:key="dKey">
          <span>{{ d }}</span>
        </th>
      </tr>
      <tr v-for="(week, weekKey) in weeks" v-bind:key="weekKey">
        <td v-for="(day, dayKey) in week" v-bind:key="dayKey">
          <span v-if="day" v-bind:class="{ hasEvent: day.hasEvent }">{{ day.day }}</span>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import axios from "axios";

export default {
  name: "Month",
  props: {
    msg: String
  },
  data: () => ({
    months: [
      "January",
      "February",
      "March",
      "April",
      "May",
      "June",
      "July",
      "August",
      "September",
      "October",
      "November",
      "December"
    ],
    dow: [
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday",
      "Sunday"
    ],
    currYear: new Date().getFullYear(),
    currMonth: new Date().getMonth(),
    weeks: [],
    events: []
  }),
  methods: {
    getMonth: function(monthAdd) {
      if (monthAdd == 1) {
        this.currMonth++;
        if (this.currMonth == 12) {
          this.currMonth = 0;
          this.currYear++;
        }
      } else if (monthAdd == -1) {
        this.currMonth--;
        if (this.currMonth == -1) {
          this.currMonth = 11;
          this.currYear--;
        }
      }

      let firstDayOfMonth =
        new Date(this.currYear, this.currMonth, 1).getDay() - 1;

      if (firstDayOfMonth == -1) {
        firstDayOfMonth = 6;
      }

      let numDaysInMonth = new Date(
        this.currYear,
        this.currMonth + 1,
        0
      ).getDate();

      let day = 1;
      let arr = [];

      for (let x = 0; x < 6; x++) {
        let arr2 = [];

        for (let y = 0; y < 7; y++) {
          let hasEvent = false;

          this.events.forEach(e => {
            let startDate = new Date(e.start_time),
                endDate = new Date(e.end_time)
              
            if (startDate.getDate() == day || endDate.getDate() == day) {
              hasEvent = true;
              return;
            }
          })

          if (x == 0) {
            if (y >= firstDayOfMonth) {
              arr2.push({day: day, hasEvent: hasEvent})
              day++;
            } else {
              arr2.push("")
            }
          } else {
            if (day <= numDaysInMonth) {
              arr2.push({day: day, hasEvent: hasEvent})
              day++;
            } else {
              arr2.push("")
            }
          }
        }

        // If all items in array are empty string, don't push array
        if (arr2.join("") != "") {
          arr.push(arr2)
        }
      }

      this.weeks = arr;
    },
    nextMonth: function() {
      this.getMonth(1);
    },
    prevMonth: function() {
      this.getMonth(-1);
    }
  },
  created: function() {
    let jwt = localStorage.getItem("jwt"),
        numDaysInMonth = new Date(this.currYear, this.currMonth + 1, 0).getDate();

    //this.getMonth(0);

    axios.get(`/api/getEvents?startDate=${this.currYear}-0${this.currMonth + 1}-01T00:00:00Z&endDate=${this.currYear}-0${this.currMonth + 1}-${numDaysInMonth}T11:59:00Z`, {headers: {jwt: jwt}})
      .then(response => {
        this.events = response.data || [];
        this.getMonth(0);
      })
      .catch(function(err) {
        console.log("ERR: " + err);
      });
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.title {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: center;
}
.title > h1 {
  margin-left: 20px;
  margin-right: 20px;
}
.week {
  height: 125px;
  width: 100%;
}
.day {
  padding-right: 75px;
  padding-bottom: 75px;
  border: black solid 1px;
}
table {
  border: black solid 1px;
  height: 450px;
  width: 450px;
  margin-left: auto;
  margin-right: auto;
}
th {
  padding: 0px 15px;
}
td, th {
  text-align: center;
}
.hasEvent {
  color: blue;
}
</style>
