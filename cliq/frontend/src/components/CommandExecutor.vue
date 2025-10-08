<template>
  <!-- 运行按钮 -->
  <div class="flex justify-center space-x-4">
    <button @click="runCommand"
      class="bg-green-500 text-white px-6 py-3 rounded-md hover:bg-green-600 focus:outline-none disabled:bg-gray-400"
      :disabled="isProcessingInternal">
      <span v-if="isProcessingInternal">处理中...</span>
      <span v-else>运行命令</span>
    </button>
    <button @click="showCommandInfo"
      class="bg-blue-500 text-white px-6 py-3 rounded-md hover:bg-blue-600 focus:outline-none">
      查看命令
    </button>
  </div>

  <!-- 命令输出 -->
  <div v-if="commandOutputInternal" class="mt-6 p-4 bg-gray-100 rounded-md">
    <h3 class="font-medium mb-2">命令输出:</h3>
    <pre class="text-sm whitespace-pre-wrap">{{ commandOutputInternal }}</pre>
  </div>

  <!-- 命令文本模态框 -->
  <div v-if="showCommandTextModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full"
    @click.self="showCommandTextModal = false">
    <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
      <div class="mt-3 text-center">
        <div class="text-gray-600 text-left">
          <h3 class="text-lg leading-6 font-medium text-gray-900">vars</h3>
          <pre class="text-sm whitespace-pre-wrap">{{ commandVariableValues }}</pre>
        </div>
        <h3 class="text-lg leading-6 font-medium text-gray-900">即将执行的命令</h3>
        <div class="mt-2 px-7 py-3">
          <p class="text-sm text-gray-500 break-all">{{ commandText }}</p>
        </div>
        <div class="items-center px-4 py-3">
          <button id="ok-btn"
            class="px-4 py-2 bg-blue-500 text-white text-base font-medium rounded-md w-full shadow-sm hover:bg-blue-600 focus:outline-none"
            @click="showCommandTextModal = false">
            关闭
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ExecuteCommand, GetCommandText } from '../../wailsjs/go/main/App';
import { useToastNotifications } from '../composables/useToastNotifications';

const props = defineProps({
  selectedCommand: { type: Object as () => any, default: null },
  commandVariableValues: { type: Object as () => { [key: string]: any }, required: true },
  isProcessing: { type: Boolean, default: false },
  commandOutput: { type: String, default: '' },
  inputFilePath: { type: String, default: '' },
  outputFilePath: { type: String, default: '' },
});

const emit = defineEmits(['update:isProcessing', 'update:commandOutput']);

const { showToast } = useToastNotifications();

const isProcessingInternal = ref(props.isProcessing);
const commandOutputInternal = ref(props.commandOutput);

watch(() => props.isProcessing, (newValue) => {
  isProcessingInternal.value = newValue;
});

watch(() => props.commandOutput, (newValue) => {
  commandOutputInternal.value = newValue;
});

watch(isProcessingInternal, (newValue) => {
  emit('update:isProcessing', newValue);
});

watch(commandOutputInternal, (newValue) => {
  emit('update:commandOutput', newValue);
});

const runCommand = async () => {
  if (!props.selectedCommand) {
    showToast('警告', '请选择要执行的命令', 'warn');
    return;
  }

  isProcessingInternal.value = true;
  commandOutputInternal.value = '';

  try {
    const result = await ExecuteCommand(props.selectedCommand.id, props.commandVariableValues);
    commandOutputInternal.value = result;
    showToast('成功', '命令执行成功', 'success');
  } catch (error) {
    showToast('错误', `命令执行失败: ${error}`, 'error');
    console.error('命令执行失败:', error);
  } finally {
    isProcessingInternal.value = false;
  }
};

const commandText = ref('');
const showCommandTextModal = ref(false);

const showCommandInfo = async () => {
  if (!props.selectedCommand) {
    showToast('警告', '请选择要执行的命令', 'warn');
    return;
  }

  try {
    const result = await GetCommandText(props.selectedCommand.id, props.commandVariableValues);
    commandText.value = result;
    showCommandTextModal.value = true;
  } catch (error) {
    showToast('错误', `获取命令文本失败: ${error}`, 'error');
    console.error('获取命令文本失败:', error);
  }
};
</script>