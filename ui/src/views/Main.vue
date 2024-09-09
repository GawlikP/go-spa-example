<template>
  <div class="p-6 rounded border-red-900 shadow-xl shadow-rose-800 w-3/4 h-min-screen h-full md:w-full">
    <p class="text-center font-bold text-stone-100">
      Go Spa Example
    </p>
    <PostsForm @refresh="refreshList()" />
    <PostsList :fetchFirstPage="fetchFirstPage" @pageRefreshed="fetchFirstPage = false" />
  </div>
</template>
<script>
import PostsForm from '../components/PostsForm.vue';
import PostsList from '../components/PostsList.vue';

export default {
  name: 'Main',
  components: {
    PostsForm,
    PostsList
  },
  data() {
    return {
      fetchFirstPage: false
    }
  },
  methods: {
    refreshList() {
      console.log('Refreshing List');
      this.fetchFirstPage = true;
    }
  }
}
</script>
