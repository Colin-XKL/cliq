<template>
  <div v-if="selectedCommand && commandVariables.length > 0" class="mb-6 p-4 bg-blue-50 rounded-md">
    <h3 class="font-medium mb-4">命令参数</h3>
    <div v-for="variable in commandVariables" :key="variable.name" class="mb-4">
      <label :for="variable.name" class="block text-sm font-medium text-gray-700 mb-2">
        {{ variable.label }}
        <span v-if="variable.required" class="text-red-500">*</span>
      </label>
      <!-- 文本输入 -->
      <InputText v-if="variable.type === 'text'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" class="w-full" :placeholder="variable.description" />
      <!-- 数字输入 -->
      <InputNumber v-else-if="variable.type === 'number'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" class="w-full" :placeholder="variable.description" />
      <!-- 布尔值 (Checkbox) -->
      <Checkbox v-else-if="variable.type === 'boolean'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" :binary="true" />
      <!-- 文件输入 -->
      <div v-else-if="variable.type === 'file_input' || variable.type === 'file_output'"
        class="flex items-center space-x-2">
        <InputText :id="variable.name" v-model="commandVariableValuesInternal[variable.name]" readonly
          :placeholder="variable.description" />
        <Button type="button" @click="openFileSelection(variable.name, variable.type)" size="small">
          选择文件
        </Button>
      </div>
      <!-- 下拉选择 -->
      <Select v-else-if="variable.type === 'select'" :id="variable.name"
        v-model="commandVariableValuesInternal[variable.name]" :options="Object.keys(variable.options || {})"
        class="w-full" :placeholder="variable.description" />
      <small v-if="variable.description" class="mt-1 text-sm text-gray-500">{{ variable.description }}</small>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue';
import { OpenFileDialog, SaveFileDialog } from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';
import InputText from 'primevue/inputtext';
import InputNumber from 'primevue/inputnumber';
import Checkbox from 'primevue/checkbox';
import Dropdown from 'primevue/dropdown';
import { useToastNotifications } from '../composables/useToastNotifications';

const props = defineProps({
  selectedCommand: { type: Object as () => any, default: null },
  commandVariableValues: { type: Object as () => { [key: string]: any }, required: true },
  inputFilePath: { type: String, default: '' },
  outputFilePath: { type: String, default: '' },
});

const emit = defineEmits(['update:commandVariableValues', 'update:inputFilePath', 'update:outputFilePath']);

const { showToast } = useToastNotifications();

const commandVariableValuesInternal = ref(props.commandVariableValues);
const inputFilePathInternal = ref(props.inputFilePath);
const outputFilePathInternal = ref(props.outputFilePath);

watch(() => props.commandVariableValues, (newValue) => {
  commandVariableValuesInternal.value = newValue;
}, { deep: true });


watch(commandVariableValuesInternal, (newValue) => {
  emit('update:commandVariableValues', newValue);
}, { deep: true });


const commandVariables = computed(() => {
  if (props.selectedCommand && props.selectedCommand.variables) {
    return Object.entries(props.selectedCommand.variables).map(([name, variable]) => ({
      name,
      ...(variable as models.Variable),
    }));
  }
  return [];
});

const openFileSelection = async (variableName: string, variableType: string) => {
  let filePath = '';
  try {
    if (variableType === 'file_input') {
      filePath = await OpenFileDialog();
      if (filePath) {
        inputFilePathInternal.value = filePath;
      }
    } else if (variableType === 'file_output') {
      filePath = await SaveFileDialog();
      if (filePath) {
        outputFilePathInternal.value = filePath;
      }
    }
    if (filePath) {
      commandVariableValuesInternal.value[variableName] = filePath;
    }
  } catch (error) {
    showToast('错误', `选择文件失败: ${error}`, 'error');
    console.error('选择文件失败:', error);
  }
};
</script>