<script lang="ts" setup>
import { ref } from 'vue';

const githubUrl = "https://github.com/wailsapp/wails";
const docsUrl = "https://wails.io/docs/introduction";
const templateUrl = "https://github.com/wailsapp/wails/tree/master/v2/pkg/templates";

const openURL = (url: string) => {
  // 简单使用a标签打开链接
  window.open(url, '_blank');
};
</script>

<template>
  <div class="flex flex-col items-center justify-center h-full">
    <div class="text-center">
      <img id="logo" alt="cliQ logo" src="./assets/images/logo-universal.png"/>
      <h1 class="text-4xl font-bold mt-4">cliQ</h1>
      <p class="text-xl mt-2">将复杂的 CLI 命令转化为直观、易用的图形用户界面</p>
      
      <div class="flex justify-center gap-4 mt-8">
        <PButton icon="pi pi-github" label="GitHub" class="p-button-outlined" @click="() => openURL(githubUrl)" />
        <PButton icon="pi pi-file" label="文档" class="p-button-outlined" @click="() => openURL(docsUrl)" />
        <PButton icon="pi pi-code" label="模板" class="p-button-outlined" @click="() => openURL(templateUrl)" />
      </div>
    </div>
  </div>
</template>

<style>
#logo {
  display: block;
  width: 150px;
  height: auto;
  margin: auto;
  background-position: center;
  background-repeat: no-repeat;
  background-size: 100% 100%;
  background-origin: content-box;
}
</style>
