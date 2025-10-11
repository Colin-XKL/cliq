<template>
  <div class="p-6 max-w-6xl mx-auto">
    <!-- 输入区域 -->
    <div class="mb-6">
      <InputText v-model="commandInput" type="text"
        placeholder="输入命令例如: pngquant {{input_file}} --output {{output_file}}"
        class="w-full p-3 border border-gray-300 rounded-md" />
      <div class="flex gap-3 my-6">
        <Button @click="generateTemplate" label="开始转换" class="bg-purple-500 hover:bg-purple-600 text-white" />
      </div>
    </div>

    <!-- 模板输出区域 -->
    <div>
      <Card>
        <template #title>
          <h2>生成的模板</h2>
          <div class="flex items-center justify-between">
            <div class="flex gap-2">
              <Button @click="copyToClipboard" label="复制" icon="pi pi-copy" size="small" />
              <Button @click="exportTemplate" label="导出" icon="pi pi-download" />
              <Button @click="saveToLocal" label="收藏保存" icon="pi pi-bookmark" />
              <Button @click="openTemplateEditor" label="编辑模板" icon="pi pi-pencil" />
            </div>
          </div>
        </template>
        <template #content>
          <div class="w-full">
            <textarea
              v-model="generatedYaml"
              readonly
              class="w-full h-96 p-3 border border-gray-300 rounded-md font-mono text-sm resize-none"
              placeholder="生成的模板将显示在这里..."></textarea>
          </div>
        </template>
      </Card>
    </div>

    <!-- Template Editor Modal -->
    <TemplateEditorModal 
      :visible="showEditorModal" 
      :initialYaml="generatedYaml"
      @close="showEditorModal = false"
      @save="onTemplateSaved" />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ParseCommandToTemplate, GenerateYAMLFromTemplate, SaveYAMLToFile, ParseYAMLToTemplate, SaveFavTemplate } from '../../wailsjs/go/main/App';
import { useToastNotifications } from '../composables/useToastNotifications';
import TemplateEditorModal from './TemplateEditorModal.vue';

const { showToast } = useToastNotifications();

const commandInput = ref('');
const generatedYaml = ref('');
const showEditorModal = ref(false);

const generateTemplate = async () => {
  if (!commandInput.value.trim()) {
    showToast('错误', '请输入CLI命令', 'error');
    return;
  }

  try {
    // 首先解析命令为模板对象
    const templateObj = await ParseCommandToTemplate(commandInput.value);
    if (templateObj) {
      // 然后将模板对象转换为YAML字符串
      const yamlStr = await GenerateYAMLFromTemplate(templateObj);
      generatedYaml.value = yamlStr;

      showToast('成功', '模板生成成功', 'success');
    }
  } catch (error) {
    showToast('错误', `生成模板失败: ${error}`, 'error');
    console.error('生成模板失败:', error);
  }
};

const copyToClipboard = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可复制的模板', 'error');
    return;
  }

  try {
    await navigator.clipboard.writeText(generatedYaml.value);
    showToast('成功', '模板已复制到剪贴板', 'success');
  } catch (error) {
    showToast('错误', '复制失败: ' + error, 'error');
  }
};

const exportTemplate = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可导出的模板', 'error');
    return;
  }

  try {
    await SaveYAMLToFile(generatedYaml.value);
    showToast('成功', '模板已导出', 'success');
  } catch (error) {
    showToast('错误', '导出失败: ' + error, 'error');
  }
};

const saveToLocal = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可保存的模板', 'error');
    return;
  }

  try {
    // Parse the YAML back to a template object
    const templateObj = await ParseYAMLToTemplate(generatedYaml.value);
    
    // Ensure the template has a name for saving
    if (!templateObj.name) {
      // Generate a name based on the command if not provided
      const commandName = commandInput.value.trim().split(' ')[0] || 'unnamed-template';
      templateObj.name = `${commandName}-${Date.now()}`;
    }
    
    // Save the template to local storage via backend
    await SaveFavTemplate(templateObj);
    showToast('成功', `模板 "${templateObj.name}" 已收藏保存`, 'success');
  } catch (error) {
    showToast('错误', '保存模板失败: ' + error, 'error');
    console.error('保存模板失败:', error);
  }
};

const openTemplateEditor = () => {
  if (!generatedYaml.value) {
    showToast('提示', '请先生成模板', 'warn');
    return;
  }
  showEditorModal.value = true;
};

const onTemplateSaved = (updatedYaml: string) => {
  generatedYaml.value = updatedYaml;
  showToast('成功', '模板已更新', 'success');
};
</script>