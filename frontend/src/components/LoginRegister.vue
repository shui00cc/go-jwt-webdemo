<template>
  <div>
    <input v-model="username" placeholder="用户名" class="input"/><br />
    <input v-model="password" type="password" placeholder="密码" class="input"/><br />
    <button @click="register" class="button">注册</button>
    <button @click="login" class="button">登录</button>
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
        const response = await fetch("http://47.108.72.107:9099/register", {
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
        if (response.status === 200) {
          alert('用户注册成功');
          location.reload();
        } else {
          alert(data)
        }
        alert(data);
      } catch (error) {
        console.error("Error:", error);
      }
    },
    async login() {
      try {
        const response = await fetch("http://47.108.72.107:9099/login", {
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
        if (response.status === 200) {
          document.cookie = `token=${data.token}`;
          this.$router.push("/order");
          alert('用户登录成功');
        } else {
          alert(data)
        }

      } catch (error) {
        console.error("Error:", error);
      }
    },
  },

};
</script>

<style scoped>
.login-register {
  text-align: center;
}

.input {
  width: 200px;
  padding: 10px;
  margin-bottom: 10px;
}

.button {
  padding: 10px 20px;
  margin-right: 10px;
}
</style>
