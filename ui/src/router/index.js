import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/Login.vue';
import Main from '../views/Main.vue';
import Register from '../views/Register.vue';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'login',
      component: Login
    },
    {
      path: '/main',
      name: 'main',
      component: Main
    },
    {
      path: '/register',
      name: 'register',
      component: Register
    }
  ]
});

async function authorize() {
  let route = null;
  try {
    const response = await fetch(`http://${serverHost}:${serverPort}/api/v1/session/authorize`, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json'
      },
      credentials: 'include',
      cors: true,
      cache: 'no-cache',
    });
    if (response.status != 200) {
      throw new Error('Session Expired');
    }
  } catch (error) {
    console.log(error);
    route = { name: 'login', query: { error: error.message } };
  } finally {
    return route;
  }
}

router.beforeEach(async (to, from, next) => {
  if (to.name !== 'login' && to.name !== 'register') {
    const route = await authorize();
    if (route) {
      next(route);
      return;
    } else {
      next();
      return;
    }
  }
  next();
});

export default router;
