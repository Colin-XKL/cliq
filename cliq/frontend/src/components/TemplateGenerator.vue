<template>
  <div class="p-6 max-w-6xl mx-auto">
    <div class="">
      <!-- 输入区域 -->
      <div>
        <InputText v-model="commandInput" type="text"
          placeholder="输入命令例如: pngquant {{input_file}} --output {{output_file}}"
          class="w-full p-3 border border-gray-300 rounded-md" />
        <div class="flex gap-3 my-6">
          <Button @click="generateTemplate" label="开始转换" class="bg-purple-500 hover:bg-purple-600 text-white" />
        </div>
      </div>

      <!-- 输出区域 -->
      <div>
        <Card>
          <template #content>
            <div class="flex justify-between items-center mb-4">
              <div class="flex gap-2">
                <Button @click="copyToClipboard" label="复制" icon="pi pi-copy" size="small"
                  class="p-button-outlined p-button-secondary" />
                <Button @click="exportTemplate" label="导出" icon="pi pi-download" size="small"
                  class="p-button-outlined p-button-secondary" />
                <Button @click="validateTemplate" label="校验模板" class="p-button-outlined p-button-secondary"
                  size="small" />
              </div>
            </div>

            <div class="min-h-96">
              <textarea v-model="generatedYaml"
                class="w-full h-full p-4 font-mono text-sm border border-gray-300 rounded-md resize-none text-gray-400"
                readonly placeholder="生成的模板将显示在这里"></textarea>
            </div>
          </template>
        </Card>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { ParseCommandToTemplate, GenerateYAMLFromTemplate, ValidateYAMLTemplate, SaveYAMLToFile } from '../../wailsjs/go/main/App';
import { models } from '../../wailsjs/go/models';
import InputText from 'primevue/inputtext';
import Button from 'primevue/button';
import { useToastNotifications } from '../composables/useToastNotifications';
import Card from 'primevue/card';
const { showToast } = useToastNotifications();

const commandInput = ref('');
const generatedYaml = ref('');

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