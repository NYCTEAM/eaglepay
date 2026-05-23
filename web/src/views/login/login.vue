<template>
  <div class="container">
    <div class="login">
      <LoginBanner v-if="isPc" />
      <div class="login_box">
        <img :src="brandLogo" alt="EAGLE SWAP LLC" class="login_brand_logo" />
        <div class="login_title">EAGLE SWAP LLC</div>
        <div class="login_title_desc">Secure crypto payment operations for USDT and EFUN settlement.</div>
        <LoginForm />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import brandLogo from "@/assets/logo/eagle-swap.png";
import { useDevicesSize } from "@/hooks/useDevicesSize";
import LoginBanner from "@/views/login/components/login-banner.vue";
import LoginForm from "@/views/login/components/login-form.vue";

const { isPc } = useDevicesSize();
</script>

<style lang="scss" scoped>
.container {
  position: relative;
  height: 100vh;
  overflow: hidden;

  .login {
    position: absolute;
    top: 50%;
    left: 50%;
    display: flex;
    align-items: center;
    max-width: 1000px;
    height: 500px;
    box-shadow: 0 0 8px 1px $color-fill-2;
    transform: translate(-50%, -50%);

    .login_box {
      position: relative;
      box-sizing: border-box;
      width: 350px;
      height: 100%;
      padding: 40px 30px 30px;

      .login_brand_logo {
        width: 84px;
        height: 84px;
        object-fit: contain;
        margin-bottom: 18px;
        filter: drop-shadow(0 12px 24px rgba(18, 99, 221, 0.18));
      }

      .login_title {
        margin-bottom: $margin-text;
        font-size: $font-size-title-2;
        font-weight: 700;
        color: $color-text-1;
      }

      .login_title_desc {
        margin-bottom: 8px;
        line-height: 1.6;
        font-size: $font-size-body-1;
        color: $color-text-3;
      }
    }
  }
}
</style>
