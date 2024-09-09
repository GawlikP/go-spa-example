<template>
  <div class="p-8 rounded-xl shadow-2xl shadow-rose-500  w-1/2 h-1/2 text-2xl">
    <p class="text-center font-bold text-blue-100">
      Registration
    </p>
    <p class="text-center font-bold text-blue-100">
      To create an account, please enter your email and password and nickname.
    </p>
    <input
      type="text"
      class="w-full mt-4 p-2 rounded border bg-rose-950 text-white"
      placeholder="Nickname"
      v-model="nickname">
    <input
      type="text"
      class="w-full mt-4 p-2 rounded border bg-rose-950 text-white"
      placeholder="Email"
      v-model="email">
    <input
      type="password"
      class="w-full mt-4 p-2 rounded border bg-rose-950 text-white"
      placeholder="Password"
      v-model="password">
    <button class="w-full mt-4 bg-rose-950 hover:bg-rose-700 text-white font-bold py-2 px-4 rounded" @click="Register">
      Create Account
    </button>
    <p class="text-red-500 mt-4" v-if="error">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter, useRoute } from 'vue-router'

const nickname = ref('');
const email = ref('');
const password = ref('');
const error = ref(null);
const router = useRouter();

const Register = async () => {
  error.value = null;
  try {
    const response = await fetch(`http://${serverHost}:${serverPort}/api/v1/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      body: JSON.stringify({
        nickname: nickname.value,
        email: email.value,
        password: password.value,
      }),
      credentials: 'include',
    });

    if (response.status === 201 && response.ok) {
      console.log('Account created successfully!');
      router.push({ name: 'login' });
    } else {
      const payload = await response.json();
      error.value = payload.message;
    }
  } catch (error) {
    console.error(error);
    error.value = 'An error occurred while creating your account. Please try again later.';
  }
};
</script>
