<template>
  <div
    class="columns is-vcentered service is-mobile is-multiline has-background-grey-lighter"
  >
    <div
      class="column is-5 is-hidden-tablet has-text-left has-text-weight-bold"
    >
      Name:
    </div>
    <div
      class="column is-7-mobile is-2-tablet has-text-left has-text-weight-bold"
    >
      {{ serviceName }}
    </div>
    <div
      class="column is-5 is-hidden-tablet has-text-left has-text-weight-bold"
    >
      Description:
    </div>
    <div class="column is-7-mobile is-4-tablet has-text-left">
      {{ description }}
    </div>
    <div
      class="column is-5 is-hidden-tablet has-text-left has-text-weight-bold"
    >
      State:
    </div>
    <div
      class="column is-7-mobile is-1-tablet is-2-desktop has-text-left-mobile"
    >
      <b-tag v-bind:class="stateTag" rounded size="is-medium">{{
        activeState
      }}</b-tag>
    </div>
    <div
      class="column is-5 is-hidden-tablet has-text-left has-text-weight-bold"
    >
      Actions:
    </div>
    <div
      class="column is-7-mobile is-5-tablet is-4-desktop has-text-left-mobile buttons"
    >
      <b-button
        type="is-info"
        :disabled="startBtnState"
        :loading="startBtnLoading"
        icon-left="play"
        @click="startService(name)"
        >Start</b-button
      >
      <b-button
        type="is-info"
        :disabled="stopBtnState"
        :loading="stopBtnLoading"
        icon-left="stop"
        @click="stopService(name)"
        >Stop</b-button
      >
      <b-button
        type="is-info"
        :loading="restartBtnLoading"
        icon-left="restart"
        @click="restartService(name)"
        >Restart</b-button
      >
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Prop, Emit, Watch, Vue } from "vue-property-decorator";

@Component({})
export default class Service extends Vue {
  @Prop() name!: string;
  @Prop() description!: string;
  @Prop() activeState!: string;
  @Prop() loadState!: string;
  @Prop() subState!: string;

  get serviceName(): string {
    var s = this.name.split(".");
    return s[0];
  }

  public startBtnLoading: boolean = false;
  public stopBtnLoading: boolean = false;
  public restartBtnLoading: boolean = false;

  get startBtnState(): boolean {
    if (this.activeState == "active") {
      return true;
    } else {
      return false;
    }
  }

  get stopBtnState(): boolean {
    if (this.activeState != "active") {
      return true;
    } else {
      return false;
    }
  }

  get stateTag(): string {
    if (this.activeState == "active") {
      return "is-success";
    } else {
      return "is-danger";
    }
  }

  get activeStateColor(): boolean {
    console.log("activestate: ", this.activeState);
    if (this.activeState == "active") {
      return true;
    }
    return false;
  }

  set active(newState: boolean) {
    if (newState == true) {
      this.$emit("startService", this.name);
    } else {
      this.$emit("stopService", this.name);
    }
  }

  //   public changeState(newState: boolean): boolean {
  //     var self = this;
  //   }

  //   public get changeState(): boolean {
  //     console.log("reqested state");
  //     return this.active;
  //   }

  @Watch("activeState")
  onActiveStateChanged(newValue: string, oldValue: string) {
    if (newValue == "active") {
      this.startBtnLoading = false;
      this.restartBtnLoading = false;
    } else if (newValue == "inactive") {
      this.stopBtnLoading = false;
    }
  }

  @Emit("startService")
  startService(serviceName: string): void {
    this.startBtnLoading = true;
  }
  @Emit("stopService")
  stopService(serviceName: string): void {
    this.stopBtnLoading = true;
  }

  @Emit("restartService")
  restartService(serviceName: string): void {
    this.restartBtnLoading = true;
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.service {
  margin-top: 15px;
}
.service-btn {
  min-width: 120px;
}
</style>
