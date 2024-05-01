<template>
  <div class="p-8 rounded-xl shadow-2xl shadow-rose-500  w-1/2 h-1/2 text-2xl">
    <p class="text-center font-bold text-blue-100">
      Login Page
    </p>
    <input type="text" class="w-full mt-4 p-2 rounded border bg-rose-950 text-white" v-model="email">
    <input type="password" class="w-full mt-4 p-2 rounded border bg-rose-950 text-white" v-model="password">
    <button class="w-full mt-4 bg-rose-950 hover:bg-rose-700 text-white font-bold py-2 px-4 rounded" @click="Login">
      Login!
    </button>
    <p class="text-red-500 mt-4" v-if="error">{{ error }}</p>
  </div>
</template>

<script>
export default {
  name: 'HelloComponent',
  data() {
    return {
      email: '',
      password: '',
      error: null,
    }
  },
  beforeMount() {
    if (this.$route.query.error) {
      this.error = this.$route.query.error;
    }
  },
  watch: {
    '$route' (to, from) {
      if(to !== from ) {
        if (to.query.error) {
          this.error = to.query.error;
        }
      }
    }
  },
  methods: {
    async Login() {
      this.error = null;
      try {
        const response = await fetch(`http://${serverHost}:${serverPort}/api/v1/login`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'Accept': 'application/json',
          },
          body: JSON.stringify({
            email: this.email,
            password: this.password,
          }),
        });

        if (response.status === 201) {
          const payload = await response.json();
          if (payload.nickname) {
            console.log("LOG IN SUCCESSFUL");
            this.$router.push({ name: 'main' });
          } else {
            throw new Error(payload.message);
          }
        } else {
          const payload = await response.json();
          if (payload.message) {
            throw new Error(payload.message);
          } else {
            throw new Error('Something went wrong with the request');
          }
        }
      } catch (error) {
        console.error(error)
        this.error = error.message;
      }
    }
  }
}
</script>
