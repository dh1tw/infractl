<template>
  <div class="message is-dark">
    <h4 class="message-header">Servicios</h4>
    <div class="message-body">
      <nav class="level">
        <div class="level-item has-text-centered">
          <div>
            <p>Reiniciar NATS Broker</p>
            <button class="button is-info" v-on:click="restartNats">Restart</button>
          </div>
        </div>
      </nav>
    </div>
  </div>
</template>

<script>
export default {
  name: "Services",
  methods: {
    restartNats: function() {
      this.$emit("restart-nats");
      console.log("ok");
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
