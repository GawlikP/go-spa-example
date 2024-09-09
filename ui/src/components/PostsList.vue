<template>
  <div class="p-6 rounded border-red-900 shadow-xl shadow-rose-800 w-3/4 h-full md:w-full">
    <p class="text-center font-bold text-stone-100">
      Recently Added Posts
    </p>

    <div class="flex p-2 bg-rose-900 rounded-xl w-full flex-col py-4">
      <div v-for="post in posts" :key="post.id" class="p-4 bg-stone-900 rounded-xl mt-4 border-red-800">
        <p class="text-white text-xl font-bold">{{ post.title }}</p>
        <p class="text-white text-lg">{{ post.content }}</p>
        <p class="text-white text-sm">{{ post.created_at }}</p>
        <p class="text-white text-sm">{{ post.user.name }}</p>
        <p class="text-white text-sm">Author {{ post.user.email }}</p>
      </div>
      <div class="flex justify-center bg-rose-900 items-center mt-4">
        <button class="bg-rose-950 hover:bg-rose-700 text-white font-bold py-2 px-4 rounded text-xl" @click="fetchPreviousPage()">
          Previous
        </button>
        <button class="bg-rose-950 hover:bg-rose-700 text-white font-bold py-2 px-4 rounded text-xl" @click="fetchNextPage()">
          Next
        </button>
      </div>
    </div>
  </div>
</template>
<script setup>
import { ref, onMounted, watch } from 'vue';
const props = defineProps({
  fetchFirstPage: Boolean
});
const emit = defineEmits(['pageRefreshed']);

const posts = ref([]);
const currentPage = ref(1);
const pageSize = ref(10);
const totalPosts = ref(0);
const totalPages = ref(0);

onMounted(async () => {
  await fetchCurrentPage()
});

watch(() => props.fetchFirstPage, async (newVal, oldVal) => {
  console.log('Current Page Changed: ', currentPage.value);
  if (newVal === true) {
    await fetchCurrentPage();
    emit('pageRefreshed');
  }
});
const fetchCurrentPage = async () => {
  const response = await fetch(`http://${serverHost}:${serverPort}/api/v1/posts?pageSize=${pageSize.value}&page=${currentPage.value}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    },
    credentials: 'include',
  });

  if (response.status === 200 && response.ok) {
    const value = await response.json();
    console.log(value);
    posts.value = value.posts;
    totalPages.value = value.total_pages;
    currentPage.value = value.current_page;
  }
};

const fetchPreviousPage = async () => {
  currentPage.value -= 1;
  if (currentPage.value < 1) {
    currentPage.value = 1;
  }
  await fetchCurrentPage();
}

const fetchNextPage = async () => {
  currentPage.value += 1;
  if (currentPage.value > totalPages.value) {
    currentPage.value = totalPages.value;
  }
  await fetchCurrentPage();
}
</script>
