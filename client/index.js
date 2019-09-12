new Vue({
    el: '#app'
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