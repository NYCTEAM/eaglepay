<template>
  <div class="banner_box">
    <div class="banner_title">
      <img :src="brandLogo" alt="EAGLE SWAP LLC" class="brand_mark" />
      {{ title }}
    </div>
    <div class="banner_img">
      <img :src="brandLogo" alt="EAGLE SWAP LLC" class="brand_hero" />
    </div>
  </div>
</template>
<script setup lang="ts">
import brandLogo from "@/assets/logo/eagle-swap.png";
// 全局title
const title = import.meta.env.VITE_GLOB_APP_TITLE;
</script>

<style lang="scss" scoped>
.banner_box {
  position: relative;
  box-sizing: border-box;
  display: flex;
  flex-direction: column;
  width: 650px;
  height: 100%;
  padding: 30px;
  .banner_title {
    position: absolute;
    display: flex;
    column-gap: $margin-text;
    align-items: center;
    font-size: $font-size-title-1;
    font-weight: bold;
    color: $color-primary;
  }
  .brand_mark {
    width: 28px;
    height: 28px;
    object-fit: contain;
    border-radius: 50%;
  }
  .banner_img {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: center;
    padding-top: 48px;
  }
  .brand_hero {
    width: min(78%, 420px);
    max-height: 360px;
    object-fit: contain;
    filter: drop-shadow(0 24px 36px rgba(22, 88, 201, 0.18));
  }
}
</style>
