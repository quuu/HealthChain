<template>
  <div class="column is-four-fifths">
     <b-select class="column" placeholder="Select Field" v-model="field" @input="addField">
       <option value="height">Height</option>
      <option value="weight">Weight</option>
      <option value="vaccination">Vaccination</option>
      <option value="sickness">Sickness</option>
      <option value="eyesight">Eyesight</option>
      <option value="blood-pressure">Blood Pressure</option>
    </b-select>
    <b-input class ="column" placeholder="Value" v-model="value" @input="addField()"></b-input>
  </div>
</template>


<script>
export default {
  props: ["id"],
  data(){
    
    return{
      field: null,
      value: null,
      show: true
    }
  },
  methods:{
    addField(){
      this.$emit('inputData', {
        id:this.id,
        field: this.field,
        value: this.value
      });
    }
  }
}
</script>