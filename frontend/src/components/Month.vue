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
      <h4>{{ months[currMonth] }} {{ currYear }}</h4>
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
    jwt: localStorage.getItem("jwt"),
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
    getMonth: function() {
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
              endDate = new Date(e.end_time);

            if (startDate.getDate() == day || endDate.getDate() == day) {
              hasEvent = true;
              return;
            }
          });

          if (x == 0) {
            if (y >= firstDayOfMonth) {
              arr2.push({ day: day, hasEvent: hasEvent });
              day++;
            } else {
              arr2.push("");
            }
          } else {
            if (day <= numDaysInMonth) {
              arr2.push({ day: day, hasEvent: hasEvent });
              day++;
            } else {
              arr2.push("");
            }
          }
        }

        // If all items in array are empty string, don't push array
        if (arr2.join("") != "") {
          arr.push(arr2);
        }
      }

      this.weeks = arr;
    },
    nextMonth: function() {
      this.currMonth++;
      if (this.currMonth == 12) {
        this.currMonth = 0;
        this.currYear++;
      }

      this.getEvents();
    },
    prevMonth: function() {
      this.currMonth--;
      if (this.currMonth == -1) {
        this.currMonth = 11;
        this.currYear--;
      }

      this.getEvents();
    },
    getEvents: function() {
      let numDaysInMonth = new Date(
          this.currYear,
          this.currMonth + 1,
          0
        ).getDate(),
        currMonth = this.currMonth + 1;

      if (currMonth < 10) {
        currMonth = "0" + currMonth;
      }

      axios
        .get(
          `/api/getEvents?startDate=${
            this.currYear
          }-${currMonth}-01T00:00:00Z&endDate=${
            this.currYear
          }-${currMonth}-${numDaysInMonth}T11:59:00Z`,
          { headers: { jwt: this.jwt } }
        )
        .then(response => {
          this.events = response.data || [];
          this.getMonth();
        })
        .catch(err => {
          console.log("ERR: " + err);
        });
    }
  },
  created: function() {
    this.getEvents();
  }
};
</script>

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
.week {
  height: 125px;
  width: 100%;
}
.day {
  padding-right: 75px;
  padding-bottom: 75px;
  border: black solid 1px;
}
th {
  padding: 0px 15px;
}
th {
  text-align: center;
}
td {
  position: relative;
  height: 125px;
  width: 200px;
  border: 1px solid lightgray;
}
td span {
  position: absolute;
  left: 0;
  top: 0;
}
.hasEvent {
  color: blue;
}
</style>
