<template>
  <div class="week">
    <HeaderTwo v-bind:title="title" v-bind:previous="prevWeek" v-bind:next="nextWeek" v-bind:goToToday="goToToday"></HeaderTwo>
    <div id="fullCalendar">
      <table id="time">
        <tr>
          <th></th>
        </tr>
        <tr v-for="(time, timeKey) in timeSlots" v-bind:key="timeKey">
          <td>
            <span>{{ time }}</span>
          </td>
        </tr>
      </table>
      <table id="calendar">
        <tr>
          <th v-for="(day, dayKey) in days" v-bind:key="dayKey">
            <div class="dayHdr">
              <button class="dateBtn" v-bind:class="{ today: day.today }">{{ day.date }}</button>
              <br/>
              <button class="dayBtn" v-bind:class="{ today: day.today }">{{ day.dow.toUpperCase() }}</button>
            </div>
          </th>
        </tr>
        <tr v-for="(time, timeKey) in timeSlots" v-bind:key="timeKey">
          <td v-for="(day, dayKey) in days" v-bind:key="dayKey">
            <span></span>
          </td>
        </tr>
      </table>
    </div>
  </div>
</template>

<script>
import { formatDateForAPI, createTimeSlots } from "../utils";
import HeaderTwo from "./HeaderTwo.vue";
import axios from "axios";
import {
  getDate,
  getMonth,
  getYear,
  addHours,
  addDays,
  subDays,
  startOfWeek,
  endOfWeek,
  isToday,
  format
} from "date-fns";

export default {
  name: "Week",
  components: {
    HeaderTwo
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
    timeSlots: [],
    startDate: new Date(),
    endDate: new Date(),
    title: "",
    days: [],
    events: []
  }),
  methods: {
    createWeek: function() {
      let daysInWeek = [];

      for (let x = 0; x < 7; x++) {
        daysInWeek.push({
          dow: this.dow[x],
          date: getDate(addDays(new Date(this.startDate), x)),
          today: isToday(addDays(new Date(this.startDate), x))
        });
      }

      this.days = daysInWeek;
    },
    nextWeek: function() {
      this.startDate = addDays(new Date(this.startDate), 7);

      this.getEvents();
    },
    prevWeek: function() {
      this.startDate = subDays(new Date(this.startDate), 7);

      this.getEvents();
    },
    getEvents: function() {
      this.startDate = format(
        startOfWeek(new Date(this.startDate), { weekStartsOn: 1 }),
        "MMM DD, YYYY"
      );

      this.endDate = format(
        endOfWeek(new Date(this.startDate), { weekStartsOn: 1 }),
        "MMM DD, YYYY"
      );

      this.title = this.startDate + " - " + this.endDate;

      axios
        .get(
          `/api/getEvents?startDate=${formatDateForAPI(
            this.startDate
          )}T00:00:00Z&endDate=${formatDateForAPI(this.endDate)}T11:59:00Z`
        )
        .then(response => {
          this.events = response.data || [];
          this.createWeek();
        })
        .catch(err => {
          // Unauthorized, send user back to log in page
          if (err && err.response && err.response.status === 401) {
            this.$emit("user");
          }

          this.createWeek();
        });
    },
    goToToday: function() {
      this.startDate = new Date();
      this.endDate = new Date();

      this.getEvents();
    }
  },
  created: function() {
    this.timeSlots = createTimeSlots(this.startDate);
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
  margin-top: 22px;
  margin-right: 2.5px;
}
#calendar {
  width: 95%;
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
.today {
  color: #3949ab;
}
.dayHdr {
  text-align: center;
}
.dayHdr button {
  background-color: white;
  border:none;
}
.dateBtn {
  height: 30px;
  width: 30px;
  border-radius: 50px;
}
.dateBtn:hover {
  background-color: #e8eaf6;
}
.dateBtn.today {
  background-color: #3949ab;
  height: 30px;
  width: 30px;
  border-radius: 50px;
  color: white;
}
.dateBtn.today:hover {
    background-color: #1a237e;
}
</style>
