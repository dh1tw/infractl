<template>
  <div id="app">
    <section class="hero">
      <div class="hero-body">
        <div class="container">
          <h1 class="title">Connectivity & Services</h1>
          <h2 class="subtitle">@ED1R</h2>
        </div>
      </div>
    </section>
    <section class="section">
      <div class="container">
        <div class="columns">
          <div class="column is-half">
            <Adsl
              :active="adsl_active"
              :ping="adsl_ping"
              :is_loading="is_loading"
              v-on:activateadsl="activateAdsl"
            ></Adsl>
          </div>
          <div class="column is-half">
            <Lte
              :connected="lte_connected"
              :active="lte_active"
              :ping="lte_ping"
              :provider="lte_provider"
              :signal="lte_signal"
              :signalbars="lte_signalbars"
              :uptime="lte_uptime"
              :network_type="lte_network_type"
              :upload_realtime="lte_upload_realtime"
              :download_realtime="lte_download_realtime"
              :quota="lte_quota"
              :consumption="lte_consumption"
              :consumption_upload="lte_consumption_upload"
              :consumption_download="lte_consumption_download"
              :is_loading="is_loading"
              v-on:reset4g="reset4g"
              v-on:activate4g="activate4g"
            ></Lte>
          </div>
        </div>
      </div>
    </section>
    <div class="section">
      <div class="container">
        <Services
          :services="services"
          v-on:startService="startService"
          v-on:stopService="stopService"
          v-on:restartService="restartService"
        ></Services>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { Component, Vue } from "vue-property-decorator";
import Lte from "./components/lte.vue";
import Adsl from "./components/adsl.vue";
import Services from "./components/services.vue";
import axios, { AxiosError } from "axios";

// set base URL if a remote server is used instead of the local golang app
// this is helpful if the golang app runs on a different machine
if (process.env.NODE_ENV === "development") {
  axios.defaults.baseURL = "http://localhost:7070";
}

@Component({
  components: {
    Adsl,
    Lte,
    Services
  }
})
export default class App extends Vue {
  private ajax_timeout: number = 2500; //ms
  private adsl_active: boolean = false;
  private adsl_ping: boolean = false;
  private loaded_status4g: boolean = false;
  private loaded_routes: boolean = false;
  private lte_restarting: boolean = false;
  private lte_signal: number = -1;
  private lte_signalbars: number = 0;
  private lte_provider: string = "";
  private lte_connected: boolean = false;
  private lte_ping: boolean = false;
  private lte_active: boolean = false;
  private lte_uptime: number = 0;
  private lte_network_type: string = "";
  private lte_upload_realtime: number = 0;
  private lte_download_realtime: number = 0;
  private lte_quota: string = "unlimited";
  private lte_consumption: number = 0;
  private lte_consumption_upload: number = 0;
  private lte_consumption_download: number = 0;
  private services: Array<object> = [];

  beforeCreated(): void {
    this.getServices();
  }

  mounted(): void {
    var self = this;
    setInterval(function() {
      self.getPingADSL();
      self.getPing4G();
      self.getRouteStatus();
    }, 3000);
    setInterval(function() {
      self.getStatus4g();
      self.getServices();
    }, 3000);
  }

  getPingADSL(): void {
    var self = this;
    axios
      .get("/api/ping/google.com", {
        timeout: this.ajax_timeout
      })
      .then(function(response) {
        // console.log(response);
        var pingRes = response.data;
        if (!pingRes.failed && Number(pingRes.rtt) > 0) {
          self.adsl_ping = true;
        } else {
          self.adsl_ping = false;
        }
      })
      .catch(function(error) {
        if (error == "Error: Network Error") {
          return;
        }
        self.adsl_ping = false;
        // self.notify(`unable to get ping over ADSL (${error})`, "is-danger");
      });
  }

  getPing4G(): void {
    var self = this;

    // return if LTE modem is resetting
    if (self.lte_restarting) {
      return;
    }

    axios
      .get("/api/ping/nats.ddns.net", {
        timeout: this.ajax_timeout
      })
      .then(function(response) {
        // console.log(response);
        var pingRes = response.data;
        if (!pingRes.failed && Number(pingRes.rtt) > 0) {
          self.lte_ping = true;
        } else {
          self.lte_ping = false;
        }
      })
      .catch(function(error) {
        if (error == "Error: Network Error") {
          return;
        }
        self.lte_ping = false;
        // self.notify(`unable to get ping over 4G (${error})`, "is-danger");
      });
  }

  getRouteStatus(): void {
    var self = this;
    axios
      .get("/api/routes", {
        timeout: this.ajax_timeout
      })
      .then(function(response) {
        // console.log(response);
        var data = response.data;
        var lte = data["4g"];
        if (lte.active) {
          self.lte_active = true;
        } else {
          self.lte_active = false;
        }
        var adsl = data["adsl"];
        if (adsl.active) {
          self.adsl_active = true;
        } else {
          self.adsl_active = false;
        }
        self.loaded_routes = true;
      })
      .catch(function(error) {
        if (error == "Error: Network Error") {
          return;
        }
        self.notify(`unable to get route information (${error})`, "is-danger");
      });
  }

  getStatus4g(): void {
    var self = this;
    axios
      .get("/api/status4g", {
        timeout: this.ajax_timeout
      })
      .then(function(response) {
        // console.log(response);
        var data = response.data;
        self.lte_signal = Number(data.lte_rsrp);
        self.lte_signalbars = Number(data.signalbar);
        self.lte_provider = data.network_provider;
        if (data.ppp_status == "ppp_connected") {
          self.lte_connected = true;
        } else {
          self.lte_connected = false;
        }
        self.lte_uptime = Number(data.realtime_time);
        self.lte_network_type = data.network_type;
        // for unknown reasons realtime_tx_thrpt is the sum of tx_thrpt and rx_thrpt
        self.lte_upload_realtime =
          Number(data.realtime_tx_thrpt * 10) -
          Number(data.realtime_rx_thrpt * 10);
        self.lte_download_realtime = Number(data.realtime_rx_thrpt * 10);
        self.lte_consumption =
          Number(data.monthly_rx_bytes) + Number(data.monthly_tx_bytes);
        self.lte_consumption_download = Number(data.monthly_rx_bytes);
        self.lte_consumption_upload = Number(data.monthly_tx_bytes);
        self.loaded_status4g = true;
        self.lte_restarting = false;
      })
      .catch(function(error) {
        // omit error if LTE modem is resetting
        if (self.lte_restarting) {
          return;
        }
        if (error == "Error: Network Error") {
          return;
        }
        self.notify(`unable to update 4G status (${error})`, "is-danger");
        self.printAxiosError(error);
      });
  }

  reset4g(): void {
    var self = this;
    this.loaded_status4g = false;
    axios
      .get("/api/reset4g", {
        timeout: this.ajax_timeout
      })
      .then(function() {
        // disable error messages until lte has restarted
        self.lte_restarting = true;
        // reset fields
        self.lte_signal = 0;
        self.lte_signalbars = -1;
        self.lte_provider = "";
        self.lte_connected = false;
        self.lte_ping = false;
        self.lte_active = false;
        self.lte_uptime = 0;
        self.lte_network_type = "";
        self.lte_upload_realtime = 0;
        self.lte_download_realtime = 0;
        self.lte_consumption = 0;
        self.lte_consumption_upload = 0;
        self.lte_consumption_download = 0;
      })
      .catch(function(error) {
        self.notify(`unable to reset the 4G modem (${error})`, "is-danger");
        self.printAxiosError(error);
      });
  }

  activate4g(): void {
    var self = this;
    axios
      .get("/api/route/adsl/disable", {
        timeout: this.ajax_timeout
      })
      .then(function() {
        // console.log("4g enabled ok");
      })
      .catch(function(error) {
        self.notify(`unable to activate 4G (${error})`, "is-danger");
        self.printAxiosError(error);
      });
  }

  activateAdsl(): void {
    var self = this;
    axios
      .get("/api/route/adsl/enable", {
        timeout: this.ajax_timeout
      })
      .then(function() {
        // nothing to do
      })
      .catch(function(error) {
        self.notify(`unable to activate ADSL (${error})`, "is-danger");
        self.printAxiosError(error);
      });
  }

  startService(serviceName: string): void {
    var self = this;
    axios
      .get("/api/service/" + serviceName + "/start")
      .then(function() {
        // nothing to do
      })
      .catch(function(error) {
        self.notify(
          `unable to start service ${serviceName} (${error})`,
          "is-danger"
        );
        self.printAxiosError(error);
      });
  }

  stopService(serviceName: string): void {
    var self = this;
    axios
      .get("/api/service/" + serviceName + "/stop")
      .then(function() {
        // nothing to do
      })
      .catch(function(error) {
        self.notify(
          `unable to start service ${serviceName} (${error})`,
          "is-danger"
        );
        self.printAxiosError(error);
      });
  }

  restartService(serviceName: string): void {
    var self = this;
    axios
      .get("/api/service/" + serviceName + "/restart")
      .then(function() {
        // nothing to do
      })
      .catch(function(error) {
        self.notify(`unable to restart service NATS (${error})`, "is-danger");
        self.printAxiosError(error);
      });
  }

  getServices(): void {
    var self = this;
    axios
      .get("/api/services")
      .then(function(services) {
        self.services = services.data;
        // console.log(services.data);
      })
      .catch(function(error) {
        self.notify(
          `unable get the list of systemd services (${error})`,
          "is-danger"
        );
        self.printAxiosError(error);
      });
  }

  notify(msg: string, type: string): void {
    this.$notification.open({
      duration: 3000,
      message: msg,
      position: "is-top-right",
      type: "is-danger",
      queue: false
      // hasIcon: true,
    });
  }

  printAxiosError(error: AxiosError): void {
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      // console.log(error.response.data);
      // console.log(error.response.status);
      // console.log(error.response.headers);
    } else if (error.request) {
      // The request was made but no response was received
      // `error.request` is an instance of XMLHttpRequest in the browser and an instance of
      // http.ClientRequest in node.js
      // console.log(error.request);
    } else {
      // Something happened in setting up the request that triggered an Error
      // console.log("Error", error.message);
    }
    // console.log(error.config);
  }

  get is_loading(): boolean {
    if (this.loaded_status4g && this.loaded_routes) {
      return false;
    }
    return true;
  }
}
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
