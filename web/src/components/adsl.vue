<template>
  <section class="frame">
    <b-loading
      :is-full-page="false"
      :active.sync="is_loading"
      :can-cancel="false"
    ></b-loading>
    <div
      class="message"
      v-bind:class="{ 'is-success': active, 'is-black': !active }"
    >
      <h4 class="message-header">ADSL</h4>
      <div class="message-body">
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-6">
              <p class="is-pulled-right">Ping:</p>
            </div>
            <div class="column is-6">
              <span
                class="tag is-pulled-left"
                v-bind:class="{ 'is-danger': !ping, 'is-success': ping }"
                >{{ _pingText }}</span
              >
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-6">
              <p class="is-pulled-right">Active (@Shack):</p>
            </div>
            <div class="column is-6">
              <span
                class="tag is-pulled-left"
                v-bind:class="{ 'is-danger': !active, 'is-success': active }"
                >{{ _activeText }}</span
              >
            </div>
          </div>
        </div>
        <hr />
        <nav class="level">
          <div class="level-item has-text-centered">
            <div>
              <p>Change to ADSL (Shack)</p>
              <b-button
                class="button is-info"
                :loading="activating_adsl"
                :disabled="active || is_loading"
                v-on:click="activateAdsl"
                >Activate</b-button
              >
            </div>
          </div>
        </nav>
      </div>
    </div>
  </section>
</template>

<script lang="ts">
import { Component, Prop, Emit, Watch, Vue } from "vue-property-decorator";

@Component({})
export default class Adsl extends Vue {
  @Prop() is_loading!: boolean;
  @Prop() active!: boolean;
  @Prop() ping!: boolean;

  activating_adsl: boolean = false;

  @Emit("activateadsl")
  activateAdsl() {
    this.activating_adsl = true;
  }

  get _activeText(): string {
    if (this.active) {
      return "Active";
    }
    return "Inactive";
  }

  get _pingText(): string {
    if (this.ping) {
      return "Received";
    } else {
      return "Not received";
    }
  }

  @Watch("active") onActiveChanged(): void {
    this.activating_adsl = false;
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.frame {
  position: relative;
}
</style>
