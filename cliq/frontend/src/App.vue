<script lang="ts" setup>
import { ref, watch } from 'vue';
import { OpenFileDialog, SaveFileDialog, ExecuteCommand, ImportTemplate } from '../wailsjs/go/main/App';
import { main } from '../wailsjs/go/models';
import TemplateManager from './components/TemplateManager.vue';
import DynamicCommandForm from './components/DynamicCommandForm.vue';
import CommandExecutor from './components/CommandExecutor.vue';
import { useToastNotifications } from './composables/useToastNotifications';

const { showToast } = useToastNotifications();

const templateData = ref<main.TemplateFile>({} as main.TemplateFile);
const selectedCommand = ref<any>(null);
const commandVariableValues = ref<{ [key: string]: any }>({});
const isProcessing = ref(false);
const commandOutput = ref('');
const inputFilePath = ref('');
const outputFilePath = ref('');

watch(templateData, () => {
  selectedCommand.value = null;
  commandVariableValues.value = {};
  inputFilePath.value = '';
  outputFilePath.value = '';
});

watch(selectedCommand, () => {
  commandVariableValues.value = {};
  inputFilePath.value = '';
  outputFilePath.value = '';
});

// 导入模板
const importTemplate = async () => {
  try {
    const result = await ImportTemplate();
    if (result) {
      templateData.value = result;
      // 显式重置 selectedCommand 为 null 以触发更新
      selectedCommand.value = null;
      // 设置为第一个命令对象
      if (result.cmds && result.cmds.length > 0) {
        selectedCommand.value = result.cmds[0];
      }
      // 重置输入和输出路径
      inputFilePath.value = '';
      outputFilePath.value = '';
      showToast('成功', '模板导入成功', 'success');
      console.log('导入的模板数据:', result);
    }
  } catch (error) {
    showToast('错误', `导入模板失败: ${error}`, 'error');
    console.error('导入模板失败:', error);
  }
};

// 处理文件选择
const openFileSelection = async (variableName: string, variableType: string) => {
  let filePath = '';
  try {
    if (variableType === 'file_input') {
      filePath = await OpenFileDialog();
    } else if (variableType === 'file_output') {
      filePath = await SaveFileDialog();
    }
    if (filePath) {
      commandVariableValues.value[variableName] = filePath;
    }
  } catch (error) {
    showToast('错误', `选择文件失败: ${error}`, 'error');
    console.error('选择文件失败:', error);
  }
};

// 执行命令
const runCommand = async () => {
  if (!selectedCommand.value) {
    showToast('警告', '请选择要执行的命令', 'warn');
    return;
  }

  isProcessing.value = true;
  commandOutput.value = '';

  try {
    // 使用选中的命令ID执行命令，并传递变量
    const result = await ExecuteCommand(selectedCommand.value.id, commandVariableValues.value);
    commandOutput.value = result;
    showToast('成功', '命令执行成功', 'success');
  } catch (error) {
    showToast('错误', `命令执行失败: ${error}`, 'error');
    console.error('命令执行失败:', error);
  } finally {
    isProcessing.value = false;
  }
};
</script>

<template>
  <div class="homepage-bg h-full w-full">
    <div class="flex flex-col items-center justify-center h-[100]vh p-6">
      <div class="w-full max-w-2xl">
        <div class="text-center mb-8">
          <h1 class="text-4xl font-bold mt-4">cliQ</h1>
          <p class="text-xl mt-2">将复杂的 CLI 命令转化为直观、易用的图形用户界面</p>
        </div>

        <div class="bg-white p-6 rounded-lg shadow-md overflow-y-auto max-h-4/5">
          <TemplateManager v-model:templateData="templateData" v-model:selectedCommand="selectedCommand" />

          <DynamicCommandForm v-if="templateData.name" :selectedCommand="selectedCommand"
            v-model:commandVariableValues="commandVariableValues" v-model:inputFilePath="inputFilePath"
            v-model:outputFilePath="outputFilePath" />

          <CommandExecutor v-if="templateData.name" :selectedCommand="selectedCommand"
            :commandVariableValues="commandVariableValues" v-model:isProcessing="isProcessing"
            v-model:commandOutput="commandOutput" :inputFilePath="inputFilePath" :outputFilePath="outputFilePath" />
        </div>
      </div>
    </div>
  </div>
  <Toast />
</template>

<style>
@import './home-bg.css';
</style>
