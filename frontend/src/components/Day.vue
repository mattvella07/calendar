<template>
  <div class="day">
    <div class="title">
      <button
        v-on:click="prevDay"
        class="nextPrevBtn"
      >
      <i class="small material-icons">chevron_left</i>
      </button>
      <button
        v-on:click="nextDay"
        class="nextPrevBtn"
      >
        <i class="small material-icons">chevron_right</i>
      </button>
      <button v-on:click="goToToday" class="todayBtn">Today</button>

      <h4>{{ currDate }}</h4>
    </div>
    <div id="fullCalendar">
      <table id="time">
        <tr v-for="(time, timeKey) in timeSlots" v-bind:key="timeKey">
          <td>
            <span>{{ time }}</span>
          </td>
        </tr>
      </table>
      <table id="calendar">
        <tr v-for="(time, timeKey) in timeSlots" v-bind:key="timeKey">
          <td>
            <span></span>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import { formatDateForAPI, createTimeSlots } from "../utils";
import axios from "axios";
import { addDays, subDays, format } from "date-fns";

export default {
  name: "Day",
  data: () => ({
    timeSlots: [],
    currDate: format(new Date(), "dddd, MMMM DD, YYYY"),
    events: []
  }),
  methods: {
    createDay: function() {},
    nextDay: function() {
      this.currDate = format(
        addDays(new Date(this.currDate), 1),
        "dddd, MMMM DD, YYYY"
      );

      this.getEvents();
    },
    prevDay: function() {
      this.currDate = format(
        subDays(new Date(this.currDate), 1),
        "dddd, MMMM DD, YYYY"
      );

      this.getEvents();
    },
    getEvents: function() {
      axios
        .get(
          `/api/getEvents?startDate=${formatDateForAPI(
            this.currDate
          )}T00:00:00Z&endDate=${formatDateForAPI(this.currDate)}T11:59:00Z`
        )
        .then(response => {
          this.events = response.data || [];
          this.createDay();
        })
        .catch(err => {
          // Unauthorized, send user back to log in page
          if (err && err.response && err.response.status === 401) {
            this.$emit("user");
          }

          this.createDay();
        });
    },
    goToToday: function() {
      this.currDate = format(new Date(), "dddd, MMMM DD, YYYY");

      this.getEvents();
    }
  },
  created: function() {
    this.timeSlots = createTimeSlots(this.currDate);
    this.getEvents();
  }
};
</script>

<style scoped>
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
#fullCalendar {
  display: flex;
  flex-direction: row;
}
#time {
  width: 5%;
}
#calendar {
  width: 95%;
  margin-top: 10px;
  margin-left: 2.5px;
}
th,
td {
  position: relative;
  height: 50px;
  width: 200px;
}
#calendar td {
  border: 1px solid lightgray;
}
#time tr {
  border: none;
}
#time td {
  position: relative;
}
#time td span {
  position: absolute;
  top: 0;
  right: 0;
}
.nextPrevBtn {
  color: black;
  border-radius: 100px;
  width: 40px;
  height: 40px;
  background-color: white;
  border: none;
  margin-left: 2.5px;
  margin-right: 2.5px;
  display: flex;
  justify-content: center;
}
.nextPrevBtn:hover {
  background-color: #e8eaf6;
}
.todayBtn {
  background-color: #3949ab;
  color: white;
  margin-left: 5px;
  margin-right: 15px;
  width: 80px;
  height: 40px;
  border-radius: 4px;
  border: none;
}
.todayBtn:hover {
  background-color: #1a237e;
}
</style>
