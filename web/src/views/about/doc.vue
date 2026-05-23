<template>
  <div class="doc-container">
    <div class="card">
      <h2>EAGLE SWAP LLC Documentation</h2>
      <div class="divider"></div>
      <p>
        Repository:
        <a href="https://github.com/NYCTEAM/eaglepay" target="_blank" class="link">https://github.com/NYCTEAM/eaglepay</a>
      </p>
      <p style="margin-top: 12px">
        API guide:
        <a href="https://github.com/NYCTEAM/eaglepay/blob/main/docs/api/api.md" target="_blank" class="link">
          https://github.com/NYCTEAM/eaglepay/blob/main/docs/api/api.md
        </a>
      </p>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.doc-container {
  padding: 20px;
}

.card {
  min-width: 0;
  background: $color-bg-2;
  border: 1px solid $color-border-2;
  border-radius: 4px;
  padding: 20px;
  box-shadow: $shadow2-center;
}

h2 {
  margin: 0 0 10px 0;
  font-size: 18px;
  color: $color-text-1;
}

.divider {
  height: 1px;
  background-color: $color-border-2;
  margin: 16px 0;
}

p {
  margin: 0;
  color: $color-text-2;
  font-size: 14px;
  line-height: 1.6;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.link {
  display: inline;
  color: $color-link;
  text-decoration: none;
  overflow-wrap: anywhere;
  word-break: break-word;
}

.link:hover {
  color: rgb(var(--link-5));
}
</style>
