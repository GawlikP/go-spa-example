<template>
  <div class="p-6 rounded border-red-900 shadow-2xl shadow-rose-800 w-full">
    <p class="text-center text-xl font-bold text-stone-100">
      Add Post
    </p>
    <div class="p-4 bg-rose-950 rounded-xl">
      <input type="text" class=" w-full mt-4 p-2 rounded border bg-rose-950 text-white" v-model="title" placeholder="Title">
      <textarea class="w-full mt-4 p-2 rounded border bg-rose-950 text-white" v-model="content" placeholder="Content"></textarea>
      <button class="w-full mt-4 bg-rose-950 hover:bg-rose-700 text-white font-bold py-2 px-4 rounded" @click="AddPost">
        Add Post
      </button>
    <div class="p-4 rounded-xl mt-4" v-show="postAdded">
      <p class="text-white text-xl">Post added successfully!</p>
    </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
const emit = defineEmits(['refresh']);

const title = ref('');
const content = ref('');
const postAdded = ref(false);

const AddPost = async () => {
  const response = await fetch(`http://${serverHost}:${serverPort}/api/v1/posts`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    body: JSON.stringify({
      title: title.value,
      content: content.value,
    }),
    credentials: 'include',
  });

  console.log(response);
  if (response.status === 201 && response.ok) {
    console.log('Post added successfully!');
    postAdded.value = true;
    emit('refresh');
  }
};
</script>
