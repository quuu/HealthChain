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
                    {{message}}
                </div>
            </div>
        </b-collapse>
    </section>
</template>


<script>
import moment from 'moment';
export default {
  

  props:['date', 'message'],
  methods:{
    formatDate(date){
      var myDate = new Date(date);
      return moment(myDate).format('MMMM Do YYYY, h:mm:ss a');
    }
  }
}
</script>