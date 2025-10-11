<template>
  <div v-if="visible" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-full max-w-6xl h-5/6 flex flex-col">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-xl font-semibold">模板编辑器</h3>
        <Button @click="closeModal" label="关闭" class="p-button-secondary" />
      </div>

      <div class="flex flex-col h-full">
        <div class="flex gap-4 mb-4">
          <Button @click="validateTemplate" label="校验模板" />
          <Button @click="previewForm" label="预览表单" />
          <Button @click="applyChanges" label="应用更改" class="p-button-success" />
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 h-full">
          <!-- 编辑器区域 -->
          <div class="flex flex-col">
            <h4 class="text-lg font-medium mb-2">模板编辑</h4>
            <div class="flex-grow border border-gray-300 rounded-md">
              <MonacoEditor :value="templateYaml" :key="editorKey" language="yaml" :height="400" theme="vs" :options="{
                minimap: { enabled: false },
                automaticLayout: true,
                fontSize: 14,
                scrollBeyondLastLine: false,
                wordWrap: 'on',
                formatOnType: true,
                formatOnPaste: true,
                readOnly: false
              }" @change="onEditorChange" />
            </div>
          </div>

          <!-- 预览区域 -->
          <div class="flex flex-col">
            <h4 class="text-lg font-medium mb-2">表单预览</h4>
            <div class="flex-grow p-4 bg-gray-50 rounded-md min-h-[400px]">
              <div v-if="previewCommand" class="w-full">
                <DynamicCommandForm :selectedCommand="previewCommand" :commandVariableValues="commandVariableValues"
                  @update:commandVariableValues="updateCommandVariableValues" />
              </div>
              <div v-else-if="hasValidationError" class="flex flex-col items-center justify-center h-64 text-red-500">
                <i class="pi pi-exclamation-triangle text-4xl mb-3"></i>
                <p>模板格式无效，请检查YAML语法</p>
              </div>
              <div v-else class="flex items-center justify-center h-64 text-gray-500">
                <p>校验模板后将显示表单预览</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue';
import MonacoEditor from 'monaco-editor-vue3';
import DynamicCommandForm from './DynamicCommandForm.vue';
import { ValidateYAMLTemplate, ParseYAMLToTemplate } from '../../wailsjs/go/main/App';
import { useToastNotifications } from '../composables/useToastNotifications';

const { showToast } = useToastNotifications();

interface Props {
  visible: boolean;
  initialYaml: string;
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  initialYaml: ''
});

const emit = defineEmits(['close', 'save']);

const templateYaml = ref(props.initialYaml);
const editorKey = ref(0);
const previewCommand = ref<any>(null);
const hasValidationError = ref(false);
const commandVariableValues = reactive<{ [key: string]: any }>({});

// Track if the change is coming from the editor to avoid infinite loops
let isUpdatingFromEditor = false;

const onEditorChange = (value: string) => {
  isUpdatingFromEditor = true;
  templateYaml.value = value;
  // The watch handler will take care of the preview update
};

// Watch for changes in templateYaml and update preview accordingly
watch(templateYaml, async (newYaml) => {
  if (!isUpdatingFromEditor) {
    // This means the value was updated from outside the editor
    // Force refresh the editor by updating the key
    editorKey.value += 1;
  } else {
    // Reset the flag after delay to ensure the editor state is stable
    setTimeout(() => {
      isUpdatingFromEditor = false;
    }, 0);
  }

  if (newYaml) {
    await updatePreviewFromYaml();
  } else {
    previewCommand.value = null;
    hasValidationError.value = false;
  }
});

watch(() => props.visible, (newVisible) => {
  if (newVisible && props.initialYaml) {
    templateYaml.value = props.initialYaml;
  }
});

const updatePreview = async (templateObj: any) => {
  if (templateObj && templateObj.cmds && templateObj.cmds.length > 0) {
    // Use the first command for preview
    previewCommand.value = templateObj.cmds[0];

    // Initialize command variable values
    if (previewCommand.value.variables) {
      Object.keys(previewCommand.value.variables).forEach(key => {
        commandVariableValues[key] = undefined; // Initialize with undefined
      });
    }
  }
};

const updatePreviewFromYaml = async () => {
  if (!templateYaml.value) {
    previewCommand.value = null;
    hasValidationError.value = false;
    return;
  }

  try {
    // Parse the YAML back to a template object for preview
    const templateObj = await ParseYAMLToTemplate(templateYaml.value);
    updatePreview(templateObj);
    hasValidationError.value = false; // Clear any validation errors
  } catch (error) {
    console.error('解析YAML模板失败:', error);
    previewCommand.value = null;
    hasValidationError.value = true; // Set validation error state
  }
};

const updateCommandVariableValues = (newValues: { [key: string]: any }) => {
  Object.assign(commandVariableValues, newValues);
};

const validateTemplate = async () => {
  if (!templateYaml.value) {
    showToast('错误', '没有可校验的模板', 'error');
    return;
  }

  try {
    await ValidateYAMLTemplate(templateYaml.value);
    showToast('成功', '模板格式有效', 'success');
    hasValidationError.value = false; // Clear validation error state

    // Update preview after validation since it's a valid template
    await updatePreviewFromYaml();
  } catch (error) {
    showToast('错误', '模板格式无效: ' + error, 'error');
    hasValidationError.value = true; // Set validation error state
  }
};

const previewForm = async () => {
  if (!templateYaml.value) {
    showToast('错误', '没有可预览的模板', 'error');
    previewCommand.value = null;
    hasValidationError.value = false;
    return;
  }

  try {
    // Validate and parse the YAML to update preview
    await ValidateYAMLTemplate(templateYaml.value);

    // If validation passes, update the preview
    const templateObj = await ParseYAMLToTemplate(templateYaml.value);
    updatePreview(templateObj);
    hasValidationError.value = false; // Clear any validation errors

    showToast('成功', '模板预览已更新', 'success');
  } catch (error) {
    // If validation fails, set preview to null to show error message
    console.error('YAML validation failed:', error);
    previewCommand.value = null;
    hasValidationError.value = true; // Set validation error state
    showToast('错误', '模板格式无效: ' + error, 'error');
  }
};

const applyChanges = () => {
  emit('save', templateYaml.value);
  closeModal();
};

const closeModal = () => {
  emit('close');
};

</script>