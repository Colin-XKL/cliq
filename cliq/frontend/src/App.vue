<script lang="ts" setup>
import { ref } from 'vue';
import { OpenFileDialog, SaveFileDialog, ExecuteCommand } from '../wailsjs/go/main/App';
import { useToast } from 'primevue/usetoast';

const toast = useToast();
const inputFilePath = ref('');
const outputFilePath = ref('');
const isProcessing = ref(false);
const commandOutput = ref('');

// 打开文件选择器
const openFileDialog = async () => {
  try {
    const result = await OpenFileDialog();
    if (result) {
      inputFilePath.value = result;
    }
  } catch (error) {
    showToast('错误', '打开文件选择器失败', 'error');
    console.error('打开文件选择器失败:', error);
  }
};

// 打开保存文件对话框
const saveFileDialog = async () => {
  try {
    const result = await SaveFileDialog();
    if (result) {
      outputFilePath.value = result;
    }
  } catch (error) {
    showToast('错误', '打开保存文件对话框失败', 'error');
    console.error('打开保存文件对话框失败:', error);
  }
};

// 执行命令
const runCommand = async () => {
  if (!inputFilePath.value) {
    showToast('警告', '请选择输入文件', 'warn');
    return;
  }

  if (!outputFilePath.value) {
    showToast('警告', '请选择输出文件路径', 'warn');
    return;
  }

  isProcessing.value = true;
  commandOutput.value = '';

  try {
    const result = await ExecuteCommand(inputFilePath.value, outputFilePath.value);
    commandOutput.value = result;
    showToast('成功', '命令执行成功', 'success');
  } catch (error) {
    showToast('错误', `命令执行失败: ${error}`, 'error');
    console.error('命令执行失败:', error);
  } finally {
    isProcessing.value = false;
  }
};
type ToastSeverity = 'success' | 'info' | 'warn' | 'error' | 'secondary' | 'contrast'
// 显示提示消息
const showToast = (summary: string, detail: string, severity: ToastSeverity) => {
  toast.add({
    severity,
    summary,
    detail,
    life: 3000
  });
};
</script>

<template>
  <div class="homepage-bg h-full w-full">
    <div class="flex flex-col items-center justify-center h-full p-6">
      <div class="w-full max-w-2xl">
        <div class="text-center mb-8">
          <h1 class="text-4xl font-bold mt-4">cliQ</h1>
          <p class="text-xl mt-2">将复杂的 CLI 命令转化为直观、易用的图形用户界面</p>
        </div>

        <div class="bg-white p-6 rounded-lg shadow-md">
          <h2 class="text-2xl font-bold mb-4 text-black">PNG图片压缩工具</h2>
          <p class="mb-6 text-gray-600">使用pngquant命令压缩PNG图片，减小文件大小</p>

          <!-- 输入文件选择 -->
          <div class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">输入文件</label>
            <div class="flex">
              <input type="text" v-model="inputFilePath"
                class="flex-1 p-2 border rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-400 font-light "
                placeholder="选择PNG图片文件" readonly />
              <button @click="openFileDialog"
                class="bg-blue-500 text-white px-4 py-2 rounded-r-md hover:bg-blue-600 focus:outline-none ">
                浏览...
              </button>
            </div>
          </div>

          <!-- 输出文件选择 -->
          <div class="mb-6">
            <label class="block text-sm font-medium text-gray-700 mb-2">输出文件</label>
            <div class="flex">
              <input type="text" v-model="outputFilePath"
                class="flex-1 p-2 border rounded-l-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-gray-400 font-light"
                placeholder="选择保存位置" readonly />
              <button @click="saveFileDialog"
                class="bg-blue-500 text-white px-4 py-2 rounded-r-md hover:bg-blue-600 focus:outline-none">
                浏览...
              </button>
            </div>
          </div>

          <!-- 运行按钮 -->
          <div class="flex justify-center">
            <button @click="runCommand"
              class="bg-green-500 text-white px-6 py-3 rounded-md hover:bg-green-600 focus:outline-none disabled:bg-gray-400"
              :disabled="isProcessing">
              <span v-if="isProcessing">处理中...</span>
              <span v-else>运行命令</span>
            </button>
          </div>

          <!-- 命令输出 -->
          <div v-if="commandOutput" class="mt-6 p-4 bg-gray-100 rounded-md">
            <h3 class="font-medium mb-2">命令输出:</h3>
            <pre class="text-sm whitespace-pre-wrap">{{ commandOutput }}</pre>
          </div>
        </div>
      </div>
    </div>
  </div>
  <Toast />
</template>

<style>
@import './home-bg.css';

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
