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
                            <b-field label="First Name"
                              label-position="on-border">
                            <div class="control">
                              <input class="input" type="text" v-model="firstname" placeholder="John">
                            </div></b-field>
                          </div>
                          <div class="field column">
                            <b-field label="Last Name"
                              label-position="on-border">
                            <div class="control">
                              <input class="input" type="text" v-model="lastname" placeholder="Doe">
                            </div></b-field>
                          </div>
                        </div>
                        <div class = "columns">
                          <div class="field column">
                            <p class="control has-icons-left">
                              <span class="select">
                                <b-field label="Country"
                                  label-position="on-border">
                                <b-select v-model="country" placeholder="Country">
                                  <option value ="USA">United States (USA)</option>
                                  <option value="GBR">United Kingdom (GBR)</option>
                                  <option value="DEU">Germany (DEU)</option>
                                </b-select></b-field>
                              </span>
                              <span class="icon is-small is-left">
                                <font-awesome-icon :icon="['fas', 'globe']"/>
                              </span>
                            </p>
                          </div>
                          <div class="field column">
                              <p class="control has-icons-left">
                                  <b-field label="Code"
                                    label-position="on-border">
                                  <input class="input" type="password" v-model="code" placeholder="Code">
                                  <span class="icon is-small is-left">
                                    <font-awesome-icon :icon="['fas', 'lock']"/>
                                  </span></b-field>
                                </p>
                          </div>
                          <div class="field column">
                              <button class="button is-danger" @click.prevent=" () => {
                                getHealthData()
                                }" >
                                  Request Records
                              </button>
                          </div>
                          <div class="field column">
                            <button class="button is-pulled-right" >
                              Clear Form
                            </button>
                          </div>
                        </div>
                      <b-loading :is-full-page="false" :active.sync="isLoading" :can-cancel="true">
                      </b-loading>
                    </form>
                </div>
            </div>
          </div>
  
    
    <section class="container" v-if="!showForm">
      <patientCard
        :firstname="firstname"
        :lastname="lastname"
        :latest_appt="latest_appt"
        :height="latest_height"
        :weight="latest_weight"
        :vaccines="vaccine_list"
      />
      
      <b-field class ="field columns">
        <div class="column"><b-button type="is-danger" @click="resetForm"><font-awesome-icon :icon="['fas', 'arrow-left']"/> Go Back</b-button></div>
        
        <b-field class="column is-7" grouped>
          
        <b-field label="Search..." label-position="on-border">
            <b-input placeholder="Search..." type="search" v-model="searchTerm"></b-input>
            <p class="control">
                <b-button class="button is-danger" @click="() => { search() }">Search</b-button>
            </p>
        </b-field>
        <b-field label="Sort By:"
            label-position="on-border">
            <b-select v-model="sortMethod" placeholder="Sort Method..." @input="() => {sortData()}">
                <option value="1">Date (Newest First)</option>
                <option value="2">Date (Oldest First)</option>
            </b-select>
        </b-field>
        </b-field>

      </b-field>
        
      
      <appointment 
        v-for="data in displayData" 
        :key="data.ID"
        :date="data.Date"
        :appt_info="data"
      />
    </section>

      
  
      </div>
  
</template>

<script>
import axios from 'axios';
import { NotificationProgrammatic as Notification } from 'buefy'
import appointment from '@/components/appointment'
import patientCard from '@/components/patientCard'
function newestSort(a, b) {
    return new Date(b.Date).getTime() - new Date(a.Date).getTime();
}
function oldestSort(a, b) {
    return new Date(a.Date).getTime() - new Date(b.Date).getTime();
}
export default {
  components:{
    appointment,
    patientCard
  },

  data() {
    return {
      showForm: true,
      firstname: null,
      lastname: null,
      country: null,
      code: null,
      showForm: true,
      isLoading: false,
      healthData: [],
      displayData: [],
      searchTerm: "",
      sortMethod:1,
      latest_appt: null,
      latest_height: null,
      latest_weight: null,
      vaccine_list: []
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
      let self =this;
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
        if(response.data == null){
          Notification.open({
                    duration: 5000,
                    message: `Error: No records available for this person`,
                    position: 'is-bottom-right',
                    type: 'is-danger',
                    hasIcon: true
                })
           self.isLoading = false;     
        }else{
          self.healthData = JSON.parse(JSON.stringify(response.data));
          self.healthData = self.healthData.map(JSON.parse);
          self.healthData.sort(newestSort);
          self.latest_appt = self.healthData[0].Date;
          self.findLatest();
          self.isLoading = false;
          self.displayData = self.healthData
          self.showForm = false;
        }
      });
    },
    openLoading() {
        this.isLoading = true
        
    },
    resetForm(){
      this.showForm = true;
      this.healthData = [];
      this.firstname = null;
      this.lastname = null;
      this.country = null;
      this.code = null;
    },
    search(){
      this.displayData = []
      let re = new RegExp(this.searchTerm)
      for(let i=0;i<this.healthData.length;i++){
        let record = JSON.stringify(this.healthData[i])
        let results = re.exec(record)
        if(results){
          this.displayData.push(JSON.parse(results.input))
        }
      }
    },
    sortData(){
      console.log("change");
        if(this.sortMethod == 1){
          this.displayData.sort(newestSort);
        }
        if(this.sortMethod == 2){
          this.displayData.sort(oldestSort);
        }

    },
    findLatest(){
      self = this;
      this.healthData.forEach(function(appt){
        if(appt.height != "" && self.latest_height == null){
          self.latest_height= appt.height;
        }
        if(appt.weight != "" && self.latest_weight == null){
          self.latest_weight= appt.weight;  
        }
        if(appt.vaccination != "" ){
          self.vaccine_list.push(appt.vaccination);
        }
      });
    
    }
  }

  
}
</script>