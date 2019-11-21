<template>
  <section>

        <b-collapse class="card" aria-id="contentIdForA11y3" :open="false">
            <div
                slot="trigger" 
                slot-scope="props"
                class="card-header"
                role="button"
                aria-controls="contentIdForA11y3">
                <p class="card-header-title">
                  {{formatDate(date)}}
                </p>
                <a class="card-header-icon">
                    <b-icon
                        :icon="props.open ? 'menu-down' : 'menu-up'">
                    </b-icon>
                </a>
            </div>
            <div class="card-content">
                <div class="content">
                    <b-message><b>Summary:</b> {{appt_info.summary}}</b-message>
                    <ul>
                      <li v-if="appt_info.height"><b >Height: </b>{{appt_info.height}}</li>
                      <li v-if="appt_info.weight"><b >Weight: </b>{{appt_info.weight}}</li>
                      <li v-if="appt_info.vaccination"><b >Vaccination: </b>{{appt_info.vaccination}}</li>
                      <li v-if="appt_info.sickness"><b >Sickness: </b>{{appt_info.sickness}}</li>
                      <li v-if="appt_info.eyesight"><b >Eyesight: </b>{{appt_info.eyesight}}</li>
                      <li v-if="appt_info['blood-pressure']"><b >Blood Pressure: </b>{{appt_info["blood-pressure"]}}</li>
                    </ul>
                </div>
            </div>
        </b-collapse>
    </section>
</template>


<script>
import moment from 'moment';
export default {
  

  props:['date', 'appt_info'],
  methods:{
    formatDate(date){
      var myDate = new Date(date);
      return moment(myDate).format('MMMM Do YYYY, h:mm:ss a');
    }
  }
}
</script>