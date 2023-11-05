<template>
  <div id="app">
<!--    <LoginRegister />-->
<!--    <OrderForm />-->
    <nav>
      <router-link to="/login">登录</router-link> |
      <router-link to="/order">order</router-link>
    </nav>
    <router-view />
  </div>
</template>

<script>
export default {
  name: "App",
  components: {
    // LoginRegister: () => import("./components/LoginRegister.vue"),
    // OrderForm: () => import("./components/OrderForm.vue"),
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
