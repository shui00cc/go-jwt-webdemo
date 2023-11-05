<template>
  <div>
    <input v-model="textInput" placeholder="请输入你想AI绘画的内容" /><br />
    <button @click="submit">提交</button>
    <img :src="cdnImageUrl" v-if="cdnImageUrl" alt="AI绘画结果" />
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

        const response = await fetch("http://localhost:9099/api/order", {
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
        if (data.success === true) {
          const uid = data.uid;
          const statusResponse = await fetch(
              `https://open.nolibox.com/prod-open-aigc/engine/status/${uid}`,
              {
                headers: {
                  'Authorization': 'Basic aXZsdEFRRXVmVEJ1Omg1SjkyWGtrOVF2RUdES3RUN1VMRVFrazBraW56ck9a'
                }
              }
          );
          const statusData = await statusResponse.json();

          // 使用cdnUrl展示图像
          this.cdnImageUrl = statusData.data.cdn;
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
