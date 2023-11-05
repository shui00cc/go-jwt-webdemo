import {createRouter, createWebHistory} from 'vue-router';

// createApp(App).use(VueRouter);

const routes = [
    { path: '/login', component: () => import('./components/LoginRegister.vue')},
    { path: '/order', component: () => import('./components/OrderForm.vue')},
];

const router = createRouter({
    history: createWebHistory(), // 不带#
    routes,
});

export default router;
