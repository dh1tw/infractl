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
      <h4 class="message-header">LTE / 4G</h4>
      <div class="message-body">
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Status:</p>
            </div>
            <div class="column is-5">
              <span
                class="tag is-pulled-left"
                v-bind:class="{
                  'is-danger': !connected,
                  'is-success': connected,
                }"
                >{{ _connectedText }}</span
              >
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Ping:</p>
            </div>
            <div class="column is-5">
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
            <div class="column is-7">
              <p class="is-pulled-right">Active (@Shack):</p>
            </div>
            <div class="column is-5">
              <span
                class="tag is-pulled-left"
                v-bind:class="{ 'is-danger': !active, 'is-success': active }"
                >{{ _activeText }}</span
              >
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Provider:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">{{ _provider }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Network Type:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">{{ _network_type }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Signal:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">
                <b-icon :icon="_signalbars"></b-icon>
                {{ _signal }}
              </p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Uptime:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">{{ _uptime }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Upload:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">{{ _upload }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Download:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">{{ _download }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right has-text-weight-bold">Monthly quota:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left has-text-weight-bold">{{ quota }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right has-text-weight-bold">Monthly usage:</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left has-text-weight-bold">
                {{ _consumption }}
              </p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7">
              <p class="is-pulled-right">Monthly usage (Up):</p>
            </div>
            <div class="column is-5">
              <p class="is-pulled-left">{{ _consumption_upload }}</p>
            </div>
          </div>
        </div>
        <div class="container">
          <div class="columns is-mobile">
            <div class="column is-7 is-two-thirds-mobile">
              <p class="is-pulled-right">Monthly usage (Down):</p>
            </div>
            <div class="column is-5 is-one-third-mobile">
              <p class="is-pulled-left">{{ _consumption_download }}</p>
            </div>
          </div>
        </div>
        <hr />
        <nav class="level">
          <div class="level-item has-text-centered">
            <div>
              <p>Reset 4G Modem</p>
              <b-button
                class="button is-info"
                :loading="resetting_4g"
                :disabled="is_loading"
                v-on:click="reset4g"
                >Reset</b-button
              >
            </div>
          </div>
          <div class="level-item has-text-centered">
            <div>
              <p>Change to 4G (Shack)</p>
              <b-button
                class="button is-info"
                :loading="activating_4g"
                :disabled="active || is_loading"
                v-on:click="activate4g"
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
export default class Lte extends Vue {
  @Prop() connected!: boolean;
  @Prop() active!: boolean;
  @Prop() is_loading!: boolean;
  @Prop() ping!: boolean;
  @Prop() provider!: string;
  @Prop() signal!: number;
  @Prop() signalbars!: number;
  @Prop() uptime!: number;
  @Prop() network_type!: string;
  @Prop() upload_realtime!: number;
  @Prop() download_realtime!: number;
  @Prop() quota!: string;
  @Prop() consumption!: number;
  @Prop() consumption_upload!: number;
  @Prop() consumption_download!: number;

  activating_4g: boolean = false;
  resetting_4g: boolean = false;

  formatData(data: number): string {
    if (data >= 1000000000) {
      //Gbyte
      return `${(data / 1000000000).toFixed(2)}GB`;
    }
    if (data >= 1000000) {
      //Mbyte
      return `${(data / 1000000).toFixed(2)}MB`;
    }
    if (data >= 1000) {
      //KiloByte
      return `${(data / 1000).toFixed(2)}kB`;
    }
    if (data >= 0) {
      //Bytes
      return `${data.toFixed(0)}Bytes`;
    }
    return "n/a";
  }

  @Emit("reset4g")
  reset4g(): void {
    this.resetting_4g = true;
  }

  @Emit("activate4g")
  activate4g(): void {
    this.activating_4g = true;
  }

  get _connectedText(): string {
    if (this.connected) {
      return "Connected";
    }
    return "Disconnected";
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
      return "Not Received";
    }
  }

  get _provider(): string {
    if (this.provider == undefined) {
      return "n/a";
    }
    return this.provider;
  }

  get _network_type(): string {
    if (this.network_type == undefined) {
      return "n/a";
    }
    return this.network_type;
  }

  get _signal(): string {
    if (this.connected && this.signal != undefined) {
      return `(${this.signal}dBm)`;
    }
    return "";
  }

  get _signalbars(): string {
    switch (this.signalbars) {
      case 0:
        return "wifi-strength-outline";
      case 1:
        return "wifi-strength-1";
      case 2:
        return "wifi-strength-2";
      case 3:
        return "wifi-strength-3";
      case 4:
        return "wifi-strength-3";
      case 5:
        return "wifi-strength-4";
      default:
        return "wifi-strength-off";
    }
  }

  get _uptime(): string {
    if (this.uptime == undefined) {
      return "n/a";
    }
    var date = new Date(this.uptime * 1000);
    var dateString =
      ("0" + date.getUTCHours()).slice(-2) +
      ":" +
      ("0" + date.getUTCMinutes()).slice(-2) +
      ":" +
      ("0" + date.getUTCSeconds()).slice(-2);
    return dateString;
  }

  get _upload(): string {
    var rate = this.formatData(this.upload_realtime);
    if (rate == "n/a") {
      return rate;
    }
    return `${rate}/s`;
  }

  get _download(): string {
    var rate = this.formatData(this.download_realtime);
    if (rate == "n/a") {
      return rate;
    }
    return `${rate}/s`;
  }

  get _consumption(): string {
    return this.formatData(this.consumption);
  }

  get _consumption_upload(): string {
    return this.formatData(this.consumption_upload);
  }

  get _consumption_download(): string {
    return this.formatData(this.consumption_download);
  }

  @Watch("active") onActiveChanged(): void {
    this.activating_4g = false;
  }

  @Watch("ping") onPingChanged(): void {
    if (this.ping && this.resetting_4g) {
      this.resetting_4g = false;
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.frame {
  position: relative;
}
</style>
