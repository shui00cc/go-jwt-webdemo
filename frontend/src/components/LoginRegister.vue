<template>
  <div>
    <input v-model="username" placeholder="用户名" /><br />
    <input v-model="password" type="password" placeholder="密码" /><br />
    <button @click="register">注册</button>
    <button @click="login">登录</button>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      password: "",
    };
  },
  methods: {
    async register() {
      try {
        const response = await fetch("http://localhost:9099/register", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: this.username,
            password: this.password,
          }),
        });

        const data = await response.json();
        alert(data);
        location.reload();
      } catch (error) {
        console.error("Error:", error);
      }
    },
    async login() {
      try {
        const response = await fetch("http://localhost:9099/login", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            username: this.username,
            password: this.password,
          }),
        });

        const data = await response.json();
        alert('user login successfully');
        document.cookie = `token=${data.token}`;
        this.$router.push("/order");
      } catch (error) {
        console.error("Error:", error);
      }
    },
  },

};
</script>
