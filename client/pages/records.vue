<template>
  <div class="hero-body">
          <div class="container ">
            <div class="columns is-5-tablet is-4-desktop is-3-widescreen">
                <div class="column">
                    <form class="box has-background-light" >
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
                                  <input class="input" type="password" v-model="ssn" placeholder="SSN">
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
                              <button class="button is-danger" @click=" () => {
                                getHealthData()
                                }" >
                                  Request Records
                              </button>
                          </div>
                        </div>
                  
                    </form>
                </div>
            </div>
          </div>

          <p v-if="firstname">First Name: {{firstname}}</p>
      <p v-if="lastname">Last Name: {{lastname}}</p>
      <p v-if="country">Country: {{country}} </p>
      <p v-if="ssn">SSN: {{ssn}}</p>
      </div>
  
</template>

<script>
import axios from 'axios';

export default {
  props:["firstname", "lastname", "country", "ssn"],

  methods: {

    async getHealthData(){
      alert("HELLO")
      await axios({
        method: 'post',
        url: 'localhost:3000/new_record',
        data: {
          firstname: firstname,
          lastname: lastname,
          country: country,
          ssn: ssn
        }
      }).then(function(response){

        console.log(response);
      });
    }
  }


  
}
</script>