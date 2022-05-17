Vue.use(HighchartsVue.default)
Vue.component('v-select', VueSelect.VueSelect);
 
var app = new Vue({
    el: "#app2",
    data() {
        return {
            title: '',
            options: ['line', 'bar', 'pie'],
            modo: 'bar',
            series: [
              {
                  name: "Channels",
                  colorByPoint: true,
                  data: []
              }
          ],
        }
    },
    mounted () {
      axios
        .get('statistic.json')
        .then(response => (this.series = response))
    },
    computed: {
        chartOptions() { 
            return {
                    chart: {  type: this.modo},
                    title: {  text: this.title  },
                    series: this.series,
              }
        },
    },
    template:
    `
    <div>
    <div class="title-row">
            <span>Style:</span>
            <v-select :options="options" v-model="modo" crearable=False></v-select>
        </div>
         
        <highcharts :options="chartOptions" ></highcharts>
    </div>
    `
})