<template>
  <div class="day">
    <HeaderTwo v-bind:title="currDate" v-bind:previous="prevDay" v-bind:next="nextDay" v-bind:goToToday="goToToday"></HeaderTwo>
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
import HeaderTwo from "./HeaderTwo.vue";
import axios from "axios";
import { addDays, subDays, format } from "date-fns";

export default {
  name: "Day",
  components: {
    HeaderTwo
  },
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

</style>
