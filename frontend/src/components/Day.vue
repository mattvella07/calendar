<template>
  <div class="day">
    <div class="title">
      <a v-on:click="prevDay" class="btn-floating btn-large waves-effect waves-light blue accent-4">
        <i class="material-icons"><</i>
      </a>
      <h4>{{ currDate }}</h4>
      <a v-on:click="nextDay" class="btn-floating btn-large waves-effect waves-light blue accent-4">
        <i class="material-icons">></i>
      </a>
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
import axios from "axios";
import moment from "moment";

export default {
  name: "Day",
  data: () => ({
    timeSlots: [
      "12am",
      "",
      "1am",
      "",
      "2am",
      "",
      "3am",
      "",
      "4am",
      "",
      "5am",
      "",
      "6am",
      "",
      "7am",
      "",
      "8am",
      "",
      "9am",
      "",
      "10am",
      "",
      "11am",
      "",
      "12pm",
      "",
      "1pm",
      "",
      "2pm",
      "",
      "3pm",
      "",
      "4pm",
      "",
      "5pm",
      "",
      "6pm",
      "",
      "7pm",
      "",
      "8pm",
      "",
      "9pm",
      "",
      "10pm",
      "",
      "11pm",
      ""
    ],
    currDate: moment().format("dddd, MMMM DD, YYYY"),
    events: []
  }),
  methods: {
    getDay: function() {},
    nextDay: function() {
      this.currDate = moment(this.currDate)
        .add(1, "days")
        .format("dddd, MMMM DD, YYYY");

      this.getEvents();
    },
    prevDay: function() {
      this.currDate = moment(this.currDate)
        .subtract(1, "days")
        .format("dddd, MMMM DD, YYYY");

      this.getEvents();
    },
    formatDateForAPI: function(date) {
      let d = new Date(date),
        yr = d.getFullYear(),
        mo = d.getMonth() + 1,
        day = d.getDate();

      if (mo < 10) {
        mo = "0" + mo;
      }

      if (day < 10) {
        day = "0" + day;
      }

      return `${yr}-${mo}-${day}`;
    },
    getEvents: function() {
      axios
        .get(
          `/api/getEvents?startDate=${this.formatDateForAPI(
            this.currDate
          )}T00:00:00Z&endDate=${this.formatDateForAPI(
            this.currDate
          )}T11:59:00Z`
        )
        .then(response => {
          this.events = response.data || [];
          this.getDay();
        })
        .catch(err => {
          // Unauthorized, send user back to log in page
          if (err.response.status === 401) {
            this.$emit("user");
          }

          this.getDay();
        });
    }
  },
  created: function() {
    this.getEvents();
  }
};
</script>

<style scoped>
#fullCalendar {
  display: flex;
  flex-direction: row;
}
#time {
  width: 5%;
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
</style>
