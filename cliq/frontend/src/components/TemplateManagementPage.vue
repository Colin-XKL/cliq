<template>
  <div class="p-4">
    <h2 class="text-2xl font-bold mb-4">模板管理</h2>
    <div v-if="favTemplates.length === 0" class="text-center text-gray-500">
      <p>暂无收藏模板。</p>
    </div>
    <div v-else>
      <DataTable :value="favTemplates" responsiveLayout="scroll">
        <Column field="name" header="模板名称"></Column>
        <Column field="description" header="描述"></Column>
        <Column field="author" header="作者"></Column>
        <Column field="version" header="版本"></Column>
        <Column header="操作">
          <template #body="slotProps">
            <Button icon="pi pi-eye" class="p-button-rounded p-button-info mr-2" @click="viewTemplate(slotProps.data)" />
            <Button icon="pi pi-pencil" class="p-button-rounded p-button-success mr-2" @click="editTemplate(slotProps.data)" />
            <Button icon="pi pi-trash" class="p-button-rounded p-button-warning" @click="confirmDeleteTemplate(slotProps.data)" />
          </template>
        </Column>
      </DataTable>
    </div>

    <Dialog v-model:visible="displayConfirmation" header="确认删除" :modal="true">
      <div class="confirmation-content">
        <i class="pi pi-exclamation-triangle mr-3" style="font-size: 2rem" />
        <span>您确定要删除模板 <b>{{ templateToDelete ? templateToDelete.name : '' }}</b> 吗？</span>
      </div>
      <template #footer>
        <Button label="取消" icon="pi pi-times" class="p-button-text" @click="closeConfirmation" />
        <Button label="删除" icon="pi pi-check" class="p-button-danger" @click="deleteTemplate" />
      </template>
    </Dialog>

    <Dialog v-model:visible="displayViewDialog" :header="`预览模板: ${templateToView ? templateToView.name : ''}`" :modal="true" :style="{ width: '50vw' }">
      <div class="p-fluid">
        <Textarea v-model="templateContentToView" rows="20" cols="30" readonly />
      </div>
      <template #footer>
        <Button label="关闭" icon="pi pi-times" class="p-button-text" @click="displayViewDialog = false" />
      </template>
    </Dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { models } from '../../wailsjs/go/models';
import { ListFavTemplates, DeleteFavTemplate, GetFavTemplate } from '../../wailsjs/go/main/App';
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import Button from 'primevue/button';
import Dialog from 'primevue/dialog';
import Textarea from 'primevue/textarea';
import { useToastNotifications } from '../composables/useToastNotifications';

const favTemplates = ref<models.TemplateFile[]>([]);
const displayConfirmation = ref(false);
const templateToDelete = ref<models.TemplateFile | null>(null);
const displayViewDialog = ref(false);
const templateToView = ref<models.TemplateFile | null>(null);
const templateContentToView = ref('');
const { showToast } = useToastNotifications();

const loadFavTemplates = async () => {
  try {
    const result = await ListFavTemplates();
    favTemplates.value = result || [];
  } catch (error) {
    console.error('Failed to list favorite templates:', error);
    showToast('错误', `加载收藏模板失败: ${error}`, 'error');
  }
};

const confirmDeleteTemplate = (template: models.TemplateFile) => {
  templateToDelete.value = template;
  displayConfirmation.value = true;
};

const deleteTemplate = async () => {
  if (templateToDelete.value) {
    try {
      await DeleteFavTemplate(templateToDelete.value.name);
      showToast('成功', `模板 ${templateToDelete.value.name} 已删除`, 'success');
      await loadFavTemplates(); // Reload the list
    } catch (error) {
      console.error('Failed to delete template:', error);
      showToast('错误', `删除模板失败: ${error}`, 'error');
    } finally {
      closeConfirmation();
    }
  }
};

const editTemplate = (template: models.TemplateFile) => {
  // TODO: Implement edit functionality
  showToast('信息', `编辑模板 ${template.name} (功能待实现)`, 'info');
};

const viewTemplate = async (template: models.TemplateFile) => {
  templateToView.value = template;
  try {
    const content = await GetFavTemplate(template.name);
    templateContentToView.value = JSON.stringify(content, null, 2);
    displayViewDialog.value = true;
  } catch (error) {
    console.error('Failed to get template content:', error);
    showToast('错误', `获取模板内容失败: ${error}`, 'error');
  }
};

const closeConfirmation = () => {
  displayConfirmation.value = false;
  templateToDelete.value = null;
};

onMounted(() => {
  loadFavTemplates();
});
</script>

<style scoped>
.confirmation-content {
  display: flex;
  align-items: center;
}
</style>