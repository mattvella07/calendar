<template>
    <div class="title">
      <button
        v-on:click="prevMonth"
        class="nextPrevBtn"
      >
      <i class="small material-icons">chevron_left</i>
      </button>
      <button
        v-on:click="nextMonth"
        class="nextPrevBtn"
      >
        <i class="small material-icons">chevron_right</i>
      </button>
      <a v-on:click="goToToday" class="todayBtn btn btn-flat waves-effect waves-light blue accent-4"><strong>Today</strong></a>

      <h4>{{ currMonthStr }} {{ currYear }}</h4>
    </div>
</template>

<script>

</script>

<style>
.title {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: start;
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
  background-color: #82b1ff;
}
.todayBtn {
  margin-right: 10px;
}
.todayBtn:hover {
  background-color: white;
}
</style>