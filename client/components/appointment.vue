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
                    Appointment: {{date}}
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
export default {
  props:['date', 'message']
}
</script>