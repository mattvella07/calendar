<template>
  <div class="week">
    <div class="title">
      <a
        v-on:click="prevWeek"
        class="btn-floating btn-large waves-effect waves-light blue accent-4"
      >
        <i class="material-icons"><</i>
      </a>
      <h4>{{ startDate }} - {{ endDate }}</h4>
      <a
        v-on:click="nextWeek"
        class="btn-floating btn-large waves-effect waves-light blue accent-4"
      >
        <i class="material-icons">></i>
      </a>
    </div>
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
            <span>{{ day.date }}</span>&nbsp;&nbsp;&nbsp;
            <span>{{ day.dow }}</span>
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
import axios from "axios";
import moment from "moment";

export default {
  name: "Week",
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
    startDate: new Date(),
    endDate: new Date(),
    days: [],
    events: []
  }),
  methods: {
    getWeek: function() {
      let arr = [],
        d = new Date(this.startDate);

      for (let x = 0; x < 7; x++) {
        arr.push({
          dow: this.dow[x],
          date: moment()
            .set({
              year: d.getFullYear(),
              month: d.getMonth(),
              date: d.getDate()
            })
            .add(x, "days")
            .get("date")
        });
      }

      this.days = arr;
    },
    nextWeek: function() {
      this.startDate = moment(this.startDate).add(7, "days");

      this.getEvents();
    },
    prevWeek: function() {
      this.startDate = moment(this.startDate).subtract(7, "days");

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
      // Get week start and end date
      if (new Date(this.startDate).getDay() === 0) {
        this.startDate = moment(this.startDate)
          .subtract(6, "days")
          .format("MMM DD, YYYY");
      } else {
        this.startDate = moment(this.startDate)
          .startOf("week")
          .add(1, "days")
          .format("MMM DD, YYYY");
      }

      this.endDate = moment(this.startDate)
        .add(6, "days")
        .format("MMM DD, YYYY");

      axios
        .get(
          `/api/getEvents?startDate=${this.formatDateForAPI(
            this.startDate
          )}T00:00:00Z&endDate=${this.formatDateForAPI(this.endDate)}T11:59:00Z`
        )
        .then(response => {
          this.events = response.data || [];
          this.getWeek();
        })
        .catch(err => {
          // Unauthorized, send user back to log in page
          if (err.response.status === 401) {
            this.$emit("user");
          }

          this.getWeek();
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
