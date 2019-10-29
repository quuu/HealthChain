<template>
  <div>
    <b-navbar class="navbar is-fixed-top has-background-light">
      <template slot="brand">
        <b-navbar-item>
          <nuxt-link to="/index"><img
            src="../assets/logo.png"
            alt="HealthChain logo"
          /></nuxt-link>
        </b-navbar-item>
      </template>

        <template slot="start">
            <b-navbar-item>
                <nuxt-link to="/index"><b-button type="is-danger" outlined>Home</b-button></nuxt-link>
            </b-navbar-item>
            <b-navbar-item>
                <nuxt-link to="/records"><b-button type="is-danger" outlined>Get Records</b-button></nuxt-link>
            </b-navbar-item>
        </template>
        
    </b-navbar>




    <!-- <section class="main-content columns">
      <aside class="column is-2 section">
        <p class="menu-label is-hidden-touch">
          General
        </p>
        <ul class="menu-list">
          <li
            v-for="(item, key) of items"
            :key="key"
          >
            <nuxt-link
              :to="item.to"
              exact-active-class="is-active"
            >
              <b-icon :icon="item.icon" /> {{ item.title }}
            </nuxt-link>
          </li>
        </ul>
      </aside>

      <div class="container column is-10">
        <nuxt />
      </div>
    </section> -->
  </div>
</template>

<script>
export default {
  data () {
    return {
      items: [
        {
          title: 'Home',
          icon: 'home',
          to: { name: 'index' }
        }
        
      ]
    }
  }
}
</script>
