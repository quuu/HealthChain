
<template>
  <div id="app">
    <b-navbar class="navbar is-fixed-top has-background-light">
      <template slot="brand">
        <b-navbar-item href="index.html">
          <img
            src="../../g3947.png"
            alt="HealthChain logo"
          >
        </b-navbar-item>
        </template>
        <template slot="start">
            <b-navbar-item href="index.html">
                <b-button type="is-danger" outlined>Home</b-button>
            </b-navbar-item>
            <b-navbar-item href="record_request.html">
                <b-button type="is-danger" outlined>Get Records</b-button>
            </b-navbar-item>
        </template>
        <template slot="end">
            
            <b-navbar-item href="https://github.com/quuu/HealthChain">
              <b-button size="is-medium"
                icon-left="github-circle">
                
            </b-button>
            </b-navbar-item>
             
           
        </template>
        
    </b-navbar>
    <HealthChainInfo/>
  </div>
</template>

<script>
import HealthChainInfo from './components/HealthChainInfo.vue'

export default {
  name: 'app',
  components: {
    HealthChainInfo
  }
}
</script>



<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
