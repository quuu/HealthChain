<template>
 <div class="is-centered">   

      <div class="column">
        <div class='container'>
            <img src='../../../g3884.png'/>
        </div>
      </div>  

    <div class="hero-body">
      <h1 class="title">Welcome to HealthChain!</h1>
      <h2 class="subtitle">A decentralized health record store with encrypted client records</h2>
    </div>

    <div class="hero-body">    
        <h1 class="title">What is HealthChain?</h1>
        <h2 class="subtitle">HealthChain is a proof of concept for a decentralized health record system.</h2>
        <h2 class="subtitle">  Our aim is to give health professionals rapid and easy access to their patients health records, regardless of location. </h2>
        <h2 class="subtitle">  This is done by assigning a private key to each person in the database, which when entered will unlock access to the patients medical records.</h2>
                          
    </div>

    <div class="tile is-ancestor">
            <div class= "tile is-1"></div>
            <div class="tile is-10 is-vertical is-child box">
                    <h2 class="title is-size-4 ">Consider the scenario:</h2>
                    <p class="box has-background-light">Tom from the US is vacationing in Greece. While there he encounters an injury and is sent to the hospital. 
                            Before a Doctor can provide him help, they may need to get his medical records mailed from home. 
                            Instead, the doctor can use Tom's private key to get access to his health information. 
                            With quick access to these records, the Doctor can more efficiently and rapidly treat Tom so that 
                            he can recover and enjoy his vacation.</p>
            </div>
            <div class="tile">
              <!-- This tile will take the rest: 2/3 -->
            </div>
    </div>
  </div>


</template>

<script>
export default {
  name: 'HealthChainInfo',
  props: {
  }
}
</script>

