<template>
  <div class="order-form">
    <input v-model="textInput" placeholder="请输入你想AI绘画的内容(英文)" class="input" /><br />
    <button @click="submit" class="button">提交</button>
    <img :src="cdnImageUrl" v-if="cdnImageUrl" alt="AI绘画结果" class="result-image" />
  </div>
</template>

<script>
export default {
  data() {
    return {
      textInput: "",
      cdnImageUrl: "",
    };
  },
  methods: {
    async submit() {
      try {
        const token = document.cookie.replace(
            /(?:(?:^|.*;\s*)token\s*=\s*([^;]*).*$)|^.*$/,
            "$1"
        );

        const response = await fetch("http://47.108.72.107:9099/api/order", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            token: `${token}`,
          },
          body: JSON.stringify({
            text: this.textInput,
          }),
        });

        const data = await response.json();
        if (response.status === 200) {
          if (data.success === true) {
            const uid = data.uid;

            const response = await fetch('http://47.108.72.107:9099/api/config',
                {
                  method: "POST",
                  headers: {
                    token: `${token}`,
                  },
                });
            const res = await response.json();
            const authorization = res.authorization;

            // 循环查询任务提交结果
            let status = '';
            let statusData;
            while (status !== 'finished') {
              const statusResponse = await fetch(
                  `https://open.nolibox.com/prod-open-aigc/engine/status/${uid}`,
                  {
                    headers: {
                      'Authorization': authorization
                    }
                  }
              );
              statusData = await statusResponse.json();

              status = statusData.status;

              if (status !== 'finished') {
                // 如果状态不是 "finished"，等待一段时间再重试
                await new Promise(resolve => setTimeout(resolve, 1000)); // 1秒钟
              }
            }

            this.cdnImageUrl = statusData.data.cdn;
          } else {
            alert('出现了一些错误，请联系管理员')
          }
        } else {
          alert(data);
        }
      } catch (error) {
        console.error("Error:", error);
      }
    },
  },
};
</script>

<style scoped>
.order-form {
  text-align: center;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.input {
  width: 300px;
  padding: 10px;
  margin-bottom: 10px;
}

.button {
  padding: 10px 20px;
}

.result-image {
  max-width: 512px;
  max-height: 512px;
  margin-top: 20px;
}
</style>

