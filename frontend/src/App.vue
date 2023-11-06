<template>
  <div id="app">
<!--    <LoginRegister />-->
<!--    <OrderForm />-->
    <nav class="navigation">
      <router-link to="/login" class="nav-link">登录</router-link> |
      <router-link to="/order" class="nav-link">AI绘画</router-link>
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
.navigation {
  text-align: center;
  margin-bottom: 20px;
}

.nav-link {
  margin: 0 10px;
  font-size: 18px;
  text-decoration: none;
  color: #333;
  cursor: pointer;
}
</style>
