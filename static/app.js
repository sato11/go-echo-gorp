const app = new Vue({
  el: '#app',
  data: {
    comments: [],
    name: '',
    text: '',
  },
  created() { this.update() },
  methods: {
    add: () => {
      const payload = { name: app.name, text: app.text };
      fetch('/api/comments', {
        body: JSON.stringify(payload),
        headers: {
          'Content-Type': 'application/json',
        },
        method: 'POST',
      }).then(res => res.json())
        .then(() => {
          app.name = '';
          app.text = '';
          app.update();
        })
        .catch(err => {
          alert(err)
        })
    },
    update: () => {
      fetch('/api/comments').then(res => res.json())
        .then(response => app.comments = response || [])
        .catch(err => console.log(err))
    }
  }
})
