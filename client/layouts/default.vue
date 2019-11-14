<template>
  <div>
    <b-navbar class="navbar is-fixed-top has-background-light">
      <template slot="brand">
        <b-navbar-item>
          <nuxt-link to="/"><img
            src="@/assets/logo.png"
            alt="HealthChain logo"
          /></nuxt-link>
        </b-navbar-item>
      </template>

        <template slot="start">
            <b-navbar-item>
                <nuxt-link to="/"><b-button type="is-danger" outlined>Home</b-button></nuxt-link>
            </b-navbar-item>
            <b-navbar-item>
                <nuxt-link to="/records"><b-button type="is-danger" outlined>Get Records</b-button></nuxt-link>
            </b-navbar-item>
            <b-navbar-item>
                <nuxt-link to="/uploads"><b-button type="is-danger" outlined>Add Records</b-button></nuxt-link>
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
  
    <section class="section"><nuxt/></section>

     <footer class="footer ">
        <div class="content has-text-centered">
          <p>
            <strong>HealthChain</strong>, a project for Software Documentation and Design at Rensselaer Polytechnic Institute. The source code is licensed
            <a href="http://opensource.org/licenses/mit-license.php">MIT</a>.
          </p>
        </div>
      </footer>

  </div>
</template>

<style>
  .navbar-link.is-active, .navbar-link:focus, .navbar-link:focus-within, .navbar-link:hover, a.navbar-item.is-active, a.navbar-item:focus, a.navbar-item:focus-within, a.navbar-item:hover{
    color: #E80000;
}

a {
    color: #E80000;
}

.navbar-item img {
    max-height: 2rem;
}

.color_1{
    background-color: #62C2FE;
    color: whitesmoke;
}

.color_2{
    color: #E80000;
}

.is-horizontal-center {
    justify-content: center;
  }

.select:not(.is-multiple):not(.is-loading)::after {
    border-color: rgb(232,0,0,.50);
    border-top-color: rgb(232,0,0,.50);
    border-right-color: rgb(232,0,0,.50);
    border-bottom-color: rgb(232,0,0,.50);
    border-left-color: rgb(232,0,0,.50);
}

.input:hover, .is-hovered.input, .is-hovered.textarea, .select select.is-hovered, .select select:hover, .taginput .is-hovered.taginput-container.is-focusable, .taginput .taginput-container.is-focusable:hover, .textarea:hover{
    border-color: #E80000;
    border-top-color: rgb(232,0,0,.50);
    border-right-color: rgb(232,0,0,.50);
    border-bottom-color: rgb(232,0,0,.50);
    border-left-color: rgb(232,0,0,.50);    
}

.input:active, .input:focus, .is-active.input, .is-active.textarea, .is-focused.input, .is-focused.textarea, .select select.is-active, .select select.is-focused, .select select:active, .select select:focus, .taginput .is-active.taginput-container.is-focusable, .taginput .is-focused.taginput-container.is-focusable, .taginput .taginput-container.is-focusable:active, .taginput .taginput-container.is-focusable:focus, .textarea:active, .textarea:focus {
    border-color:  #E80000;  
    box-shadow: 0 0 0 0.125em  rgb(232,0,0,.25);
}
</style>

