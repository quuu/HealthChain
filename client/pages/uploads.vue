<template>

  <div class="hero-body">
          <div class="container ">
            <div class="columns is-5-tablet is-4-desktop is-3-widescreen">
                <div class="column">
                    <form class="box has-background-light" >
                        <div class="field has-text-centered">
                            <img src="@/assets/logo.png" width="400">
                        </div>
                        <p class="title is-3 has-text-centered">Health Record Submission</p>
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
                            <label class="label">Country</label>
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
                              <label class="label">Unique Code</label>
                              <p class="control has-icons-left">
                                  <input class="input" type="password" v-model="code" placeholder="Code">
                                  <span class="icon is-small is-left">
                                    <i class="fas fa-lock"></i>
                                  </span>
                                </p>
                          </div>

                         
                          
                          <div class="field column"></div>
                          
                          
                        </div>
                        <div class="columns">
                          <div class="field column">
                            Appointment Information:
                            
                            <b-input class="control" type="textarea" placeholder="Summary" v-model="summary"></b-input>

                            <div class="control">
                              Fields:
                                <b-field v-for="(field, index) in fields" :key="field.id">
                                 <div class="columns"> 
                                  <formField 
                                   :id="field.id"
                                   @inputData= "inputData"/>
                                  <b-button type="is-danger" icon-right="delete" @click="fields.splice(index, 1)" />
                                  </div>
                               </b-field>

                            </div>
                          </div>
                          
                        </div>
                        <b-button rounded @click="addField(numFields++)">
                            Add Field
                          </b-button>

                      <div class="field column is-pulled-right">
                              <button class="button is-danger" @click.prevent=" () => {
                                postHealthData()
                                }" >
                                  Upload Records 
                              </button>
                          </div>
                      <b-loading :is-full-page="false" :active.sync="isLoading" :can-cancel="true">
                      </b-loading>
                    </form>
                </div>
            </div>
          </div>

  </div>
</template>

<script>

import axios from 'axios';
import formField from '@/components/formField'
import { NotificationProgrammatic as Notification } from 'buefy'
export default {
  components:{
    formField
  },
  data() {
    return {
      firstname: null,
      lastname: null,
      country: null,
      code: null,
      message: {},
      isLoading: false,
      numFields:0,
      fields: [],
      summary: null
    
    };
  },
  methods:{
    openLoading() {
        this.isLoading = true
    },
    loadData(){
      let form = new FormData()
      form.append("first", this.firstname)
      form.append("last", this.lastname)
      form.append("country", this.country)
      form.append("code", this.code)
      this.message["summary"]=this.summary;
      var self=this;
      this.fields.forEach(function(data, index){
        self.message[data.field]= data.value;
      });
      form.append("appointment_info", this.message);
      return form
    },
    loadData2(){
      this.message["summary"]=this.summary;
      var self=this;
      this.fields.forEach(function(data, index){
        self.message[data.field]= data.value;
      });
      return  {
        "first": this.firstname,
        "last": this.lastname,
        "country": this.country,
        "code": this.code,
        "appointment_info": this.message

      }
    },
    
    async postHealthData(){
      this.openLoading();
      let self =this;
      console.log(this.loadData2());
      await axios({
        method: 'post',
        url: '/api/new_record',

        headers: { 'Access-Control-Allow-Origin': '*',
        'Access-Control-Allow-Methods': 'GET,PUT,POST,DELETE,OPTIONS',
        'Content-Type': 'multipart/form-data'},
        data: this.loadData2()
        // {
        //   first: this.firstname,
        //   last: this.lastname,
        //   country: this.country,
        //   code: this.code
        // }
      }).then(function(response){
        //this.showForm =false;        
        console.log(response);
        if(response.data == "Saved!"){
          Notification.open({
                    duration: 5000,
                    message: `Record for `+ self.firstname + ` uploaded successfully!`,
                    position: 'is-bottom-right',
                    type: 'is-success',
                    hasIcon: true
                })
        }
        self.isLoading = false;
        self.firstname = null;
        self.lastname = null;
        self.country = null; 
        self.code = null;
        self.message = null;

      
      });
    },
    addField(ID){
      this.fields.push({id: ID});
      
    },
    inputData(data){
      var i;
      this.fields.forEach(function(field, index){
        if(data.id == field.id){
          i = index;
        }
      });
      if(data.field != null){
        this.fields[i].field = data.field;
        this.fields[i].value = data.value;
      }
      
    }
  }

  
  
}
</script>