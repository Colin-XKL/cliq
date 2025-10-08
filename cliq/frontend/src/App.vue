<script lang="ts" setup>
import { ref } from 'vue';
import { main } from '../wailsjs/go/models';
import TemplateManager from './components/TemplateManager.vue';
import DynamicCommandForm from './components/DynamicCommandForm.vue';
import CommandExecutor from './components/CommandExecutor.vue';
import TemplateGenerator from './components/TemplateGenerator.vue';
import Button from 'primevue/button';

const templateData = ref<main.TemplateFile>({} as main.TemplateFile);
const selectedCommand = ref<any>(null);
const commandVariableValues = ref<{ [key: string]: any }>({});
const isProcessing = ref(false);
const commandOutput = ref('');
const currentView = ref<'main' | 'generator'>('main'); // Add view state

const resetTemplate = () => {
  templateData.value = {} as main.TemplateFile;
  selectedCommand.value = null;
  commandVariableValues.value = {};
  isProcessing.value = false;
  commandOutput.value = '';
};
</script>

<template>
  <div class="homepage-bg h-full w-full">
    <div class="flex flex-col items-center justify-center h-[100]vh p-6">
      <div class="w-full max-w-4xl">
        <div class="text-center mb-8">
          <h1 class="text-4xl font-bold mt-4">cliQ</h1>
          <p class="text-xl mt-2">将复杂的 CLI 命令转化为直观、易用的图形用户界面</p>
        </div>
        
        <!-- Navigation tabs -->
        <div class="flex justify-center mb-6">
          <div class="inline-flex bg-gray-100 rounded-lg p-1">
            <button 
              @click="currentView = 'main'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'main' ? 'bg-white shadow text-gray-900' : 'text-gray-600']"
            >
              主界面
            </button>
            <button 
              @click="currentView = 'generator'"
              :class="['px-4 py-2 rounded-md text-sm font-medium', currentView === 'generator' ? 'bg-white shadow text-gray-900' : 'text-gray-600']"
            >
              模板生成器
            </button>
          </div>
        </div>

        <div class="bg-white p-6 rounded-lg shadow-md overflow-y-auto max-h-[70vh]">
          <!-- Main View -->
          <div v-if="currentView === 'main'">
            <Button @click="resetTemplate">Reset</Button>

            <TemplateManager v-model:templateData="templateData" v-model:selectedCommand="selectedCommand"
              @reset-template="resetTemplate" />

            <DynamicCommandForm v-if="templateData.name" :selectedCommand="selectedCommand"
              v-model:commandVariableValues="commandVariableValues"/>

            <CommandExecutor v-if="templateData.name" :selectedCommand="selectedCommand"
              :commandVariableValues="commandVariableValues" v-model:isProcessing="isProcessing"
              v-model:commandOutput="commandOutput" />
          </div>

          <!-- Template Generator View -->
          <div v-if="currentView === 'generator'">
            <TemplateGenerator />
          </div>
        </div>
      </div>
    </div>
  </div>
  <Toast />
</template>

<style>
@import './styles/home-bg.css';
</style>
