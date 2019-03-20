<template>
  <div id="spinner" class="preloader-wrapper small active">
    <div class="spinner-layer spinner-blue-only">
      <div class="circle-clipper left">
        <div class="circle"></div>
      </div>
      <div class="gap-patch">
        <div class="circle"></div>
      </div>
      <div class="circle-clipper right">
        <div class="circle"></div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: "LoadingSpinner",
  props: {
    active: Boolean
  }
};
</script>

<style scoped>
#spinner {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  margin-left: auto;
  margin-right: auto;
}
</style>