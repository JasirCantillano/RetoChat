new Vue({
    el: '#app',
    data () {
      return {
        channels: []
      }
    },
    mounted () {
      axios
        .get('channels.json')
        .then(response => (this.channels = response.data))
    }
  })