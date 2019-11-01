<template>
  <div class="hero-body">
          <div class="container ">
            <div class="columns is-5-tablet is-4-desktop is-3-widescreen">
                <div class="column">
                    <form class="box has-background-light" v-if="showForm" >
                        <div class="field has-text-centered">
                            <img src="@/assets/logo.png" width="400">
                        </div>
                        <p class="title is-3 has-text-centered">Health Records Request</p>
                        <div class="columns"> 
                          <div class="field column">
                            <label class="label">First Name</label>
                            <div class="control">
                              <input class="input" type="text" v-model="firstname" placeholder="John">
                            </div>
                          </div>
                          <div class="field column">
                            <label class="label">Last Name</label>
                            <div class="control">
                              <input class="input" type="text" v-model="lastname" placeholder="Doe">
                            </div>
                          </div>
                        </div>
                        <div class = "columns">
                          <div class="field column">
                            <p class="control has-icons-left">
                              <span class="select">
                                <select v-model="country">
                                  <option value="" selected>Country</option>
                                  <option value ="USA">United States (USA)</option>
                                  <option value="GBR">United Kingdom (GBR)</option>
                                  <option value="DEU">Germany (DEU)</option>
                                </select>
                              </span>
                              <span class="icon is-small is-left">
                                <i class="fas fa-globe"></i>
                              </span>
                            </p>
                          </div>
                          <div class="field column">
                              <p class="control has-icons-left">
                                  <input class="input" type="password" v-model="code" placeholder="Code">
                                  <span class="icon is-small is-left">
                                    <i class="fas fa-lock"></i>
                                  </span>
                                </p>
                          </div>
                          <div class="field column">
                            <button class="button is-pulled-right" >
                              Clear Form
                            </button>
                          </div>
                          <div class="field column">
                              <button class="button is-danger" @click.prevent=" () => {
                                getHealthData()
                                }" >
                                  Request Records
                              </button>
                          </div>
                        </div>
                      <b-loading :is-full-page="false" :active.sync="isLoading" :can-cancel="true">
                      </b-loading>
                    </form>
                </div>
            </div>
          </div>
  
    <h1 v-if="firstname" class="title is-1 text-justify-center">Health Info for  {{firstname + " " + lastname}}</h1>

      <appointment
        date="Friday, November 1st, 2019"
        Message="You have aids, we're gonna have to kill you"
      />
  
      </div>
  
</template>

<script>
import axios from 'axios';

import appointment from '@/components/appointment'

export default {
  components:{
    appointment
  },

  data() {
    return {
      showForm: true,
      firstname: null,
      lastname: null,
      country: null,
      code: null,
      showForm: true,
      isLoading: false
    };
  },
  
  methods: {
    loadData(){
      let form = new FormData()
      form.append("first", this.firstname)
      form.append("last", this.lastname)
      form.append("country", this.country)
      form.append("code", this.code)

      return form
    },
    async getHealthData(){
      this.openLoading();
      await axios({
        method: 'post',
        url: '/api/get_records',

        headers: { 'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET,PUT,POST,DELETE,OPTIONS',
        'Content-Type': 'multipart/form-data'},
        data: this.loadData()
        // {
        //   first: this.firstname,
        //   last: this.lastname,
        //   country: this.country,
        //   code: this.code
        // }
      }).then(function(response){
        //this.showForm =false;        
        console.log(response);
      });
    },
    openLoading() {
        this.isLoading = true
        setTimeout(() => {
          this.isLoading = false
        }, 10 * 1000)
      }
  }

  
}
</script>