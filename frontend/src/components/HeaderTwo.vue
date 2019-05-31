<template>
    <div class="title">
      <button
        v-on:click="previous"
        class="nextPrevBtn"
      >
      <i class="small material-icons">chevron_left</i>
      </button>
      <button
        v-on:click="next"
        class="nextPrevBtn"
      >
        <i class="small material-icons">chevron_right</i>
      </button>
      <a v-on:click="goToToday" class="todayBtn btn btn-flat waves-effect waves-light"><strong>Today</strong></a>

      <h4>{{ title }}</h4>
    </div>
</template>

<script>
  export default {
    name: "HeaderTwo",
    props: {
      title: String,
      previous: Function,
      next: Function,
      goToToday: Function
    }
  }
</script>

<style>

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
.nextPrevBtn:focus {
  background-color: white;
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