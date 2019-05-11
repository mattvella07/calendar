<template>
  <div class="hello">
    <div class="title">
      <a
        v-on:click="prevMonth"
        class="btn-floating btn-large waves-effect waves-light blue accent-4"
      >
        <i class="material-icons"><</i>
      </a>
      <a v-on:click="goToToday" class="btn waves-effect waves-light blue accent-4">Today</a>
      <a
        v-on:click="nextMonth"
        class="btn-floating btn-large waves-effect waves-light blue accent-4"
      >
        <i class="material-icons">></i>
      </a>

      <h4>{{ currMonthStr }} {{ currYear }}</h4>
    </div>
    <table>
      <tr>
        <th v-for="(d, dKey) in dow" v-bind:key="dKey">
          <span>{{ d }}</span>
        </th>
      </tr>
      <tr v-for="(week, weekKey) in weeks" v-bind:key="weekKey">
        <td v-for="(day, dayKey) in week" v-bind:key="dayKey">
          <span v-if="day" v-bind:class="{ hasEvent: day.hasEvent, today: day.today }">{{ day.day }}</span>
        </td>
      </tr>
    </table>
  </div>
</template>

<script>
import { formatDateForAPI } from "../utils";
import axios from "axios";
import {
  getDay,
  getMonth,
  getYear,
  getDaysInMonth,
  isToday,
  format
} from "date-fns";

export default {
  name: "Month",
  props: {
    msg: String
  },
  data: () => ({
    dow: [
      "Monday",
      "Tuesday",
      "Wednesday",
      "Thursday",
      "Friday",
      "Saturday",
      "Sunday"
    ],
    currYear: getYear(new Date()),
    currMonth: getMonth(new Date()),
    currMonthStr: format(new Date(), "MMMM"),
    weeks: [],
    events: []
  }),
  methods: {
    createMonth: function() {
      let firstDayOfMonth =
        getDay(new Date(this.currYear, this.currMonth, 1)) - 1;

      if (firstDayOfMonth == -1) {
        firstDayOfMonth = 6;
      }

      let numDaysInMonth = getDaysInMonth(
          new Date(this.currYear, this.currMonth)
        ),
        day = 1,
        allWeeks = [];

      for (let x = 0; x < 6; x++) {
        let week = [];

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
              week.push({
                day: day,
                hasEvent: hasEvent,
                today: isToday(new Date(this.currYear, this.currMonth, day))
              });
              day++;
            } else {
              week.push("");
            }
          } else {
            if (day <= numDaysInMonth) {
              week.push({
                day: day,
                hasEvent: hasEvent,
                today: isToday(new Date(this.currYear, this.currMonth, day))
              });
              day++;
            } else {
              week.push("");
            }
          }
        }

        // If all items in array are empty string, don't push array
        if (week.join("") != "") {
          allWeeks.push(week);
        }
      }

      this.weeks = allWeeks;
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
      this.currMonthStr = format(
        new Date(this.currYear, this.currMonth, 1),
        "MMMM"
      );

      let numDaysInMonth = getDaysInMonth(
        new Date(this.currYear, this.currMonth)
      );

      axios
        .get(
          `/api/getEvents?startDate=${formatDateForAPI(
            format(new Date(this.currYear, this.currMonth, 1), "MMM DD, YYYY")
          )}T00:00:00Z&endDate=${formatDateForAPI(
            format(
              new Date(this.currYear, this.currMonth, numDaysInMonth),
              "MMM DD, YYYY"
            )
          )}T11:59:00Z`
        )
        .then(response => {
          this.events = response.data || [];
          this.createMonth();
        })
        .catch(err => {
          // Unauthorized, send user back to log in page
          if (err && err.response && err.response.status === 401) {
            this.$emit("user");
          }

          this.createMonth();
        });
    },
    goToToday: function() {
      this.currYear = getYear(new Date());
      this.currMonth = getMonth(new Date());

      this.getEvents();
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
  color: white;
  margin: 0px 5px;
}
.title {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: start;
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
.today {
  color: red;
}
</style>
