<template>
  <section class="section">  
    <div class="card">
      <div class="card-content">
      <div class="media">
        <div class="media-left">
          <figure class="image is-128x128">
            <img src="@/assets/defaultprofile.jpeg" alt="Placeholder image">
          </figure>
        </div>
        <div class="media-content">
          <p class="title is-8">Summary: {{firstname}} {{lastname}}</p> 
          <ul>
            <li v-if="height"><b >Height: </b>{{height}}</li>
            <li v-if="weight"><b >Weight: </b>{{weight}}</li>
            <li v-if="check_vaccines(vaccines)"><b >Vaccintation(s): </b>{{list_vaccines(vaccines)}}</li>
          </ul>        
        </div>
      </div>
      </div>
      <footer class="card-footer">
        <div class="card-footer-item">Last Appointment: {{formatDate(latest_appt)}}</div>
      </footer>
    </div>
  </section>
</template>

<script>
import moment from 'moment';

export default {
  props:["firstname", "lastname", "latest_appt", "height", "weight", "vaccines"],
  methods:{
    formatDate(date){
      var myDate = new Date(date);
      return moment(myDate).format('MMMM Do YYYY, h:mm:ss a');
    },
    check_vaccines(vaccines){
      return (vaccines.length > 0);
    },
    list_vaccines(vaccines){
      var v_string = "";
      vaccines.forEach(function(vax, index){
        if(index != vaccines.length-1){
          v_string += vax + ", "
        }else{
          v_string += vax
        }
      });
      return v_string;
    }
  }

}
</script>

<style>

</style>