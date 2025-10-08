<template>
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
      <Dropdown v-model="selectedCommandInternal" :options="templateData.cmds" optionLabel="name" class="w-full"
        placeholder="选择要执行的命令">
        <template #value="slotProps">
          <div class="flex align-items-center">
            <div>{{ slotProps.value.name }}</div>
          </div>
        </template>
        <template #option="slotProps">
          <div class="flex flex-col text-left">
            <div>{{ slotProps.option.name }}</div>
          </div>
        </template>
      </Dropdown>
      <p class="mt-2 text-sm text-gray-500" v-if="selectedCommandInternal">{{ selectedCommandInternal.description }}</p>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, computed } from 'vue';
import { ImportTemplate } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import Dropdown from 'primevue/dropdown';
import { useToastNotifications } from '../composables/useToastNotifications';

const props = defineProps({
  templateData: { type: Object as () => main.TemplateFile, required: true },
  selectedCommand: { type: Object as () => any, default: null },
});

const emit = defineEmits(['update:templateData', 'update:selectedCommand', 'reset-template']);

const { showToast } = useToastNotifications();

const templateDataInternal = ref(props.templateData);
const selectedCommandInternal = ref(props.selectedCommand);

watch(() => props.templateData, (newValue) => {
  templateDataInternal.value = newValue;
});

watch(() => props.selectedCommand, (newValue) => {
  selectedCommandInternal.value = newValue;
});

watch(templateDataInternal, (newValue) => {
  emit('update:templateData', newValue);
}, { deep: true });

watch(selectedCommandInternal, (newValue) => {
  emit('update:selectedCommand', newValue);
});

const importTemplate = async () => {
  try {
    const result = await ImportTemplate();
    if (result) {
      templateDataInternal.value = result;
      selectedCommandInternal.value = null; // Reset selected command on new template import
      if (result.cmds && result.cmds.length > 0) {
        selectedCommandInternal.value = result.cmds[0];
      }
      emit('reset-template');
      showToast('成功', '模板导入成功', 'success');
    }
  } catch (error) {
    showToast('错误', `导入模板失败: ${error}`, 'error');
    console.error('导入模板失败:', error);
  }
};

</script>