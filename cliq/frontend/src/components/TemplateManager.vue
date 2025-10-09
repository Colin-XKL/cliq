<template>
  <div v-if="!templateData.name" class="text-center py-12">
    <h2 class="text-2xl font-bold text-black mb-4">欢迎使用 cliQ</h2>
    <p class="text-gray-600 mb-8">请导入模板文件以开始使用</p>
    <div class="flex justify-center gap-4">
      <button @click="importTemplate"
        class="bg-purple-500 text-white px-6 py-3 rounded-md hover:bg-purple-600 focus:outline-none text-lg">
        导入模板
      </button>
      <button @click="showUrlImportDialog = true"
        class="bg-blue-500 text-white px-6 py-3 rounded-md hover:bg-blue-600 focus:outline-none text-lg">
        URL导入
      </button>
    </div>

    <div v-if="favTemplates && favTemplates.length > 0" class="mt-8">
      <h3 class="text-xl font-bold text-black mb-4">或从收藏夹选择</h3>
      <Listbox v-model="selectedFavTemplate" :options="favTemplates" optionLabel="name" 
        class="w-full md:w-56 mx-auto" @change="loadFavTemplate" />
    </div>
  </div>

  <div v-else>
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-2xl font-bold text-black">{{ templateData.name }}</h2>
      <div class="flex gap-2">
        <button @click="importTemplate"
          class="bg-purple-500 text-white px-4 py-2 rounded-md hover:bg-purple-600 focus:outline-none">
          更换模板
        </button>
        <button @click="showUrlImportDialog = true"
          class="bg-blue-500 text-white px-4 py-2 rounded-md hover:bg-blue-600 focus:outline-none">
          URL导入
        </button>
        <button @click="addTemplateToFavorites"
          class="bg-yellow-500 text-white px-4 py-2 rounded-md hover:bg-yellow-600 focus:outline-none">
          收藏
        </button>
      </div>
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

  <!-- URL导入对话框 -->
  <div v-if="showUrlImportDialog" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-6 rounded-lg w-full max-w-md">
      <h3 class="text-lg font-semibold mb-4">从URL导入模板</h3>
      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">模板URL</label>
        <input v-model="templateUrl" type="text" 
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          placeholder="https://example.com/template.cliqfile.yaml">
      </div>
      <div class="flex justify-end gap-2">
        <button @click="cancelUrlImport" 
          class="px-4 py-2 text-gray-700 hover:bg-gray-100 rounded-md">
          取消
        </button>
        <button @click="importTemplateFromUrl" 
          class="px-4 py-2 bg-blue-500 text-white rounded-md hover:bg-blue-600">
          导入
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import { ImportTemplate, ImportTemplateFromURL } from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';
import Dropdown from 'primevue/dropdown';
import Listbox from 'primevue/listbox';
import { useToastNotifications } from '../composables/useToastNotifications';
import { SaveFavTemplate } from '../../wailsjs/go/main/App';

const props = defineProps({
  templateData: { type: Object as () => models.TemplateFile, required: true },
  selectedCommand: { type: Object as () => any, default: null },
  favTemplates: { type: Array as () => models.TemplateFile[], default: () => [] },
});

const emit = defineEmits(['update:templateData', 'update:selectedCommand', 'reset-template', 'fav-template-updated']);

const { showToast } = useToastNotifications();

const templateDataInternal = ref(props.templateData);
const selectedCommandInternal = ref(props.selectedCommand);
const showUrlImportDialog = ref(false);
const templateUrl = ref('');
const selectedFavTemplate = ref<models.TemplateFile | null>(null);

watch(() => props.templateData, (newValue) => {
  templateDataInternal.value = newValue;
});

watch(() => props.selectedCommand, (newValue) => {
  selectedCommandInternal.value = newValue;
});

watch(() => props.favTemplates, (newValue) => {
  if (newValue.length > 0 && !selectedFavTemplate.value) {
    selectedFavTemplate.value = newValue[0];
  }
});

watch(templateDataInternal, (newValue) => {
  emit('update:templateData', newValue);
}, { deep: true });

watch(selectedCommandInternal, (newValue) => {
  emit('update:selectedCommand', newValue);
});

const importTemplate = async () => {
  try {
    emit('reset-template');
    const result = await ImportTemplate();
    if (result) {
      templateDataInternal.value = result;
      selectedCommandInternal.value = null; // Reset selected command on new template import
      if (result.cmds && result.cmds.length > 0) {
        selectedCommandInternal.value = result.cmds[0];
      }
      showToast('成功', '模板导入成功', 'success');
    }
  } catch (error) {
    showToast('错误', `导入模板失败: ${error}`, 'error');
    console.error('导入模板失败:', error);
  }
};

const importTemplateFromUrl = async () => {
  if (!templateUrl.value.trim()) {
    showToast('错误', '请输入模板URL', 'error');
    return;
  }

  try {
    emit('reset-template');
    const result = await ImportTemplateFromURL(templateUrl.value);
    if (result) {
      templateDataInternal.value = result;
      selectedCommandInternal.value = null; // Reset selected command on new template import
      if (result.cmds && result.cmds.length > 0) {
        selectedCommandInternal.value = result.cmds[0];
      }
      showToast('成功', '模板导入成功', 'success');
      cancelUrlImport(); // Close the dialog
    }
  } catch (error) {
    showToast('错误', `从URL导入模板失败: ${error}`, 'error');
    console.error('从URL导入模板失败:', error);
  }
};

const loadFavTemplate = () => {
  if (selectedFavTemplate.value) {
    emit('reset-template');
    templateDataInternal.value = selectedFavTemplate.value;
    selectedCommandInternal.value = null;
    if (selectedFavTemplate.value.cmds && selectedFavTemplate.value.cmds.length > 0) {
      selectedCommandInternal.value = selectedFavTemplate.value.cmds[0];
    }
    showToast('成功', `已加载收藏模板: ${selectedFavTemplate.value.name}`, 'success');
  }
};

const cancelUrlImport = () => {
  showUrlImportDialog.value = false;
  templateUrl.value = '';
};

const addTemplateToFavorites = async () => {
  if (!templateDataInternal.value || !templateDataInternal.value.name) {
    showToast('错误', '没有可收藏的模板', 'error');
    return;
  }

  try {
    await SaveFavTemplate(templateDataInternal.value);
    showToast('成功', `模板 ${templateDataInternal.value.name} 已收藏`, 'success');
    emit('fav-template-updated'); // Notify parent to refresh favorite templates
  } catch (error) {
    showToast('错误', `收藏模板失败: ${error}`, 'error');
    console.error('收藏模板失败:', error);
  }
};

</script>