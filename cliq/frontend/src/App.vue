<script lang="ts" setup>
import { ref, computed } from 'vue';
import { OpenFileDialog, SaveFileDialog, ExecuteCommand, ImportTemplate } from '../wailsjs/go/main/App';
import { useToast } from 'primevue/usetoast';
import { main } from '../wailsjs/go/models';
import { Ref } from 'vue';
import Dropdown from 'primevue/dropdown';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Checkbox from 'primevue/checkbox';
import FileUpload from 'primevue/fileupload';

const toast = useToast();
const inputFilePath = ref('');
const outputFilePath = ref('');
const isProcessing = ref(false);
const commandOutput = ref('');
const templateData: Ref<main.TemplateFile> = ref({} as main.TemplateFile);
const selectedCommand = ref<any>(null); // 当前选中的命令
const commandVariableValues: Ref<{ [key: string]: any }> = ref({}); // 存储动态表单的变量值

// 计算属性，用于将 selectedCommand.variables 转换为数组，方便遍历
const commandVariables = computed(() => {
  if (selectedCommand.value && selectedCommand.value.variables) {
    return Object.entries(selectedCommand.value.variables).map(([name, variable]) => ({
      name,
      ...(variable as main.Variable),
    }));
  }
  return [];
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
    const result = await ExecuteCommand(inputFilePath.value, outputFilePath.value, selectedCommand.value.id, commandVariableValues.value);
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
          <div v-if="!templateData.name" class="text-center py-12">
            <h2 class="text-2xl font-bold text-black mb-4">欢迎使用 cliQ</h2>
            <p class="text-gray-600 mb-8">请导入模板文件以开始使用</p>
            <button @click="importTemplate"
              class="bg-purple-500 text-white px-6 py-3 rounded-md hover:bg-purple-600 focus:outline-none text-lg">
              导入模板
            </button>
          </div>

          <div v-else>
            <div class="flex justify-between items-center mb-4">
              <h2 class="text-2xl font-bold text-black">{{ templateData.name }}</h2>
              <button @click="importTemplate"
                class="bg-purple-500 text-white px-4 py-2 rounded-md hover:bg-purple-600 focus:outline-none">
                更换模板
              </button>
            </div>
            <p class="mb-6 text-gray-600">{{ templateData.description }}</p>

            <!-- 模板信息显示 -->
            <div class="mb-6 p-4 bg-blue-50 rounded-md">
              <p class="text-xs text-blue-500">作者: {{ templateData.author }} | 版本: {{ templateData.version }}</p>
            </div>

            <!-- 命令选择 -->
            <div class="mb-6" v-if="templateData.cmds && templateData.cmds.length > 0">
              <label class="block text-sm font-medium text-gray-700 mb-2">选择命令</label>
              <Dropdown v-model="selectedCommand" :options="templateData.cmds" optionLabel="name" class="w-full"
                placeholder="选择要执行的命令">
                <template #value="slotProps">
                  <div class="flex align-items-center">
                    <div>{{ slotProps.value.name }}</div>
                  </div>
                </template>
                <template #option="slotProps">
                  <div class="flex flex-col">
                    <div>{{ slotProps.option.name }}</div>
                    <small class="text-gray-500">{{ slotProps.option.description }}</small>
                  </div>
                </template>
              </Dropdown>
              <p class="mt-2 text-sm text-gray-500" v-if="selectedCommand">{{ selectedCommand.description }}</p>
            </div>

            <!-- 动态表单 -->
            <div v-if="selectedCommand && commandVariables.length > 0" class="mb-6 p-4 bg-blue-50 rounded-md">
              <h3 class="font-medium mb-4">命令参数</h3>
              <div v-for="variable in commandVariables" :key="variable.name" class="mb-4">
                <label :for="variable.name" class="block text-sm font-medium text-gray-700 mb-2">
                  {{ variable.label }}
                  <span v-if="variable.required" class="text-red-500">*</span>
                </label>
                <!-- 文本输入 -->
                <InputText v-if="variable.type === 'text'" :id="variable.name"
                  v-model="commandVariableValues[variable.name]" class="w-full" :placeholder="variable.description" />
                <!-- 数字输入 -->
                <InputNumber v-else-if="variable.type === 'number'" :id="variable.name"
                  v-model="commandVariableValues[variable.name]" class="w-full" :placeholder="variable.description" />
                <!-- 布尔值 (Checkbox) -->
                <Checkbox v-else-if="variable.type === 'boolean'" :id="variable.name"
                  v-model="commandVariableValues[variable.name]" :binary="true" />
                <!-- 文件输入 -->
                <div v-else-if="variable.type === 'file_input' || variable.type === 'file_output'"
                  class="flex items-center space-x-2">
                  <InputText :id="variable.name" v-model="commandVariableValues[variable.name]" class="w-full" readonly
                    :placeholder="variable.description" />
                  <button type="button" @click="openFileSelection(variable.name, variable.type)"
                    class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600 focus:outline-none">
                    选择文件
                  </button>
                </div>
                <!-- 下拉选择 -->
                <Dropdown v-else-if="variable.type === 'select'" :id="variable.name"
                  v-model="commandVariableValues[variable.name]" :options="Object.keys(variable.options || {})"
                  class="w-full" :placeholder="variable.description" />
                <small v-if="variable.description" class="mt-1 text-sm text-gray-500">{{ variable.description }}</small>
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
  </div>
  <Toast />
</template>

<style>
@import './home-bg.css';
</style>
