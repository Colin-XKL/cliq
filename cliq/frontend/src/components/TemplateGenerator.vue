<template>
  <div class="p-6 max-w-6xl mx-auto">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">模板生成器</h1>
    
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
      <!-- 左侧输入区域 -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <h2 class="text-xl font-semibold mb-4">输入命令</h2>
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-2">CLI命令</label>
          <InputText 
            v-model="commandInput" 
            type="text" 
            placeholder="例如: pngquant {{input_file}} --output {{output_file}}"
            class="w-full p-3 border border-gray-300 rounded-md"
          />
        </div>
        
        <div class="flex gap-3 mt-6">
          <Button 
            @click="generateTemplate" 
            label="开始转换" 
            class="bg-purple-500 hover:bg-purple-600 text-white"
          />
          <Button 
            @click="validateTemplate" 
            label="校验模板" 
            class="bg-blue-500 hover:bg-blue-600 text-white"
          />
        </div>
      </div>
      
      <!-- 右侧输出区域 -->
      <div class="bg-white p-6 rounded-lg shadow-md">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">生成的模板</h2>
          <div class="flex gap-2">
            <Button 
              @click="copyToClipboard" 
              label="复制" 
              icon="pi pi-copy"
              class="p-button-outlined p-button-secondary"
            />
            <Button 
              @click="exportTemplate" 
              label="导出" 
              icon="pi pi-download"
              class="p-button-outlined p-button-secondary"
            />
          </div>
        </div>
        
        <div class="h-96">
          <textarea 
            v-model="generatedYaml" 
            class="w-full h-full p-4 font-mono text-sm border border-gray-300 rounded-md resize-none text-gray-400"
            readonly
            placeholder="生成的模板将显示在这里"
          ></textarea>
        </div>
      </div>
    </div>
    
    <!-- 错误提示 -->
    <div v-if="errorMessage" class="mt-6 p-4 bg-red-50 border border-red-200 rounded-md">
      <p class="text-red-600">{{ errorMessage }}</p>
    </div>
    
    <!-- 成功提示 -->
    <div v-if="successMessage" class="mt-6 p-4 bg-green-50 border border-green-200 rounded-md">
      <p class="text-green-600">{{ successMessage }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ParseCommandToTemplate, GenerateYAMLFromTemplate, ValidateYAMLTemplate, ExportTemplateToFile } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import { useToastNotifications } from '../composables/useToastNotifications';

const { showToast } = useToastNotifications();

const commandInput = ref('');
const generatedYaml = ref('');
const errorMessage = ref('');
const successMessage = ref('');

const generateTemplate = async () => {
  if (!commandInput.value.trim()) {
    errorMessage.value = '请输入CLI命令';
    successMessage.value = '';
    return;
  }
  
  try {
    errorMessage.value = '';
    successMessage.value = '';
    
    // 首先解析命令为模板对象
    const templateObj = await ParseCommandToTemplate(commandInput.value);
    if (templateObj) {
      // 然后将模板对象转换为YAML字符串
      const yamlStr = await GenerateYAMLFromTemplate(templateObj);
      generatedYaml.value = yamlStr;
      
      successMessage.value = '模板生成成功';
    }
  } catch (error) {
    errorMessage.value = `生成模板失败: ${error}`;
    successMessage.value = '';
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
    // In a real implementation, we would call the backend to save the file
    // For now, create and download a file directly from frontend
    const blob = new Blob([generatedYaml.value], { type: 'text/yaml' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'generated_template.cliqfile.yaml';
    document.body.appendChild(a);
    a.click();
    window.URL.revokeObjectURL(url);
    document.body.removeChild(a);
    
    showToast('成功', '模板已导出', 'success');
  } catch (error) {
    showToast('错误', '导出失败: ' + error, 'error');
  }
};

const validateTemplate = async () => {
  if (!generatedYaml.value) {
    showToast('错误', '没有可校验的模板', 'error');
    return;
  }
  
  try {
    await ValidateYAMLTemplate(generatedYaml.value);
    showToast('成功', '模板格式有效', 'success');
  } catch (error) {
    showToast('错误', '模板格式无效: ' + error, 'error');
  }
};
</script>