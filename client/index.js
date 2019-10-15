var app = new Vue({
    el: '#app',
    data: {
        showForm: true,
        firstname: null,
        lastname: null,
        country: null,
        ssn: null
    },
    methods:{
        getRecords: function(event){
           console.log("CLICKED");
            event.preventDefault();
            let formData = {
               firstname: this.firstname,
                lastname: this.lastname,
                country: this.country,
                ssn: this.ssn
            }
            console.log(formData)
        }
    }
})



function getHealthData(data){
    fetch('http://localhost:5000?/data?firstname='+data.firstname + "&lastname=" + data.lastname+ "&country="+ data.country + "&ssn=" + data.ssn)
        .then(function(response){
            return response.json();
        })
        .then(function(hData){
            console.log*hData();
        });
}