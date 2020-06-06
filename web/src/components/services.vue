<template>
  <section class="frame">
    <div class="message is-dark">
      <h4 class="message-header">Services</h4>
      <div class="message-body">
        <div class="container">
          <div class="columns is-hidden-mobile">
            <div class="column is-2 has-text-left has-text-weight-bold">Service Name</div>
            <div class="column is-4 has-text-left has-text-weight-bold">Description</div>
            <div class="column is-2 has-text-weight-bold">Status</div>
            <div class="column is-4 has-text-weight-bold">Actions</div>
          </div>
          <div class="container service" v-for="service in servicesSorted" :key="service.name">
            <Service
              :name="service.Name"
              :description="service.Description"
              :activeState="service.ActiveState"
              :loadState="service.LoadState"
              :subState="service.SubState"
              v-on:startService="startService"
              v-on:stopService="stopService"
              v-on:restartService="restartService"
            ></Service>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { Component, Prop, Emit, Watch, Vue } from "vue-property-decorator";
import Service from "./service.vue";

@Component({
  components: {
    Service
  }
})
export default class Services extends Vue {
  @Prop() services: Array<object> = [];

  get servicesSorted(): Array<object> {
    var self = this;
    var sortedServices = self.services;
    sortedServices.sort((a: any, b: any) => (a.Name > b.Name ? 1 : -1));
    return sortedServices;
  }

  public isActive(active: string): boolean {
    if (active == "active") {
      return true;
    }
    return false;
  }

  @Emit("startService")
  startService(serviceName: string): void {}

  @Emit("stopService")
  stopService(serviceName: string): void {}

  @Emit("restartService")
  restartService(serviceName: string): void {}
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
