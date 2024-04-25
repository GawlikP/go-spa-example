import { createRouter, createWebHistory } from 'vue-router';
import HelloComponent from '../components/HelloComponent.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HelloComponent
    },
  ]
});

export default router;
