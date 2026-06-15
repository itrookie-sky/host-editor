<template>
  <div class="app-layout">
    <header class="titlebar" data-wails-draggable>
      <span class="titlebar-title">Host Editor</span>
      <div class="titlebar-actions">
        <button class="btn-save" :disabled="!isDirty" @click="handleSave">
          {{ isDirty ? "Save" : "Saved" }}
        </button>
      </div>
    </header>
    <div class="main-area">
      <HostSidebar :files="files" :activeName="activeName" :isDirty="isDirty" @select="handleSelect" @create="handleCreate" @delete="handleDelete" />
      <div class="editor-wrap">
        <HostsEditor v-if="activeName" v-model="content" :readOnly="loading" />
        <div v-else class="empty-hint">Select or create a version</div>
      </div>
    </div>
    <div v-if="error" class="toast" @click="error = ''">{{ error }}</div>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch } from "vue";
import HostSidebar from "./components/HostSidebar.vue";
import HostsEditor from "./components/HostsEditor.vue";
import type { HostFileInfo } from "./types/host";
import { ListHostFiles, ReadHostFile, SaveHostFile, CreateHostFile, DeleteHostFile } from "../wailsjs/go/view/App";

const files = ref<HostFileInfo[]>([]);
const activeName = ref("");
const content = ref("");
const savedContent = ref("");
const loading = ref(false);
const error = ref("");

const isDirty = ref(false);

watch(content, val => {
  isDirty.value = val !== savedContent.value;
});

function showError(msg: string) {
  error.value = msg;
  setTimeout(() => {
    error.value = "";
  }, 3000);
}

async function loadFiles() {
  try {
    files.value = await ListHostFiles();
    if (files.value.length > 0 && !activeName.value) {
      await handleSelect(files.value[0].name);
    }
  } catch (e: any) {
    showError(e.message || String(e));
  }
}

async function handleSelect(name: string) {
  if (isDirty.value && activeName.value) {
    const ok = confirm(`You have unsaved changes in "${activeName.value}". Discard?`);
    if (!ok) return;
  }
  loading.value = true;
  try {
    content.value = await ReadHostFile(name);
    savedContent.value = content.value;
    activeName.value = name;
    isDirty.value = false;
  } catch (e: any) {
    showError(e.message || String(e));
  } finally {
    loading.value = false;
  }
}

async function handleSave() {
  if (!activeName.value || !isDirty.value) return;
  try {
    await SaveHostFile({ name: activeName.value, content: content.value });
    savedContent.value = content.value;
    isDirty.value = false;
  } catch (e: any) {
    showError(e.message || String(e));
  }
}

async function handleCreate() {
  const name = prompt("New version name:")?.trim();
  if (!name) return;
  try {
    const file = await CreateHostFile(name);
    files.value = [...files.value.filter(item => item.name !== file.name), file];
    await handleSelect(name);
  } catch (e: any) {
    showError(e.message || String(e));
  }
}

async function handleDelete(name: string) {
  if (!confirm(`Delete version "${name}"?`)) return;
  try {
    await DeleteHostFile(name);
    if (activeName.value === name) {
      activeName.value = "";
      content.value = "";
      savedContent.value = "";
      isDirty.value = false;
    }
    await loadFiles();
  } catch (e: any) {
    showError(e.message || String(e));
  }
}

function handleKeydown(e: KeyboardEvent) {
  if ((e.metaKey || e.ctrlKey) && e.key === "s") {
    e.preventDefault();
    handleSave();
  }
}

onMounted(() => {
  loadFiles();
  window.addEventListener("keydown", handleKeydown);
});
</script>

<style>
* {
  box-sizing: border-box;
}

html,
body {
  margin: 0;
  padding: 0;
  height: 100%;
  background: #1e222a;
  color: #ccc;
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", sans-serif;
}

#app {
  height: 100vh;
  overflow: hidden;
}
</style>

<style scoped>
.app-layout {
  display: flex;
  flex-direction: column;
  height: 100vh;
}

.titlebar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 66px;
  padding: 28px 78px 0 14px;
  background: #181c24;
  border-bottom: 1px solid #2a2f3b;
  --wails-draggable: drag;
  -webkit-app-region: drag;
  user-select: none;
  flex-shrink: 0;
}

.titlebar-title {
  font-size: 13px;
  font-weight: 600;
  color: #999;
}

.titlebar-actions {
  display: flex;
  gap: 8px;
  --wails-draggable: no-drag;
  -webkit-app-region: no-drag;
}

.btn-save {
  padding: 3px 12px;
  border: 1px solid #3a3f4b;
  border-radius: 4px;
  background: transparent;
  color: #aaa;
  font-size: 12px;
  cursor: pointer;
}

.btn-save:not(:disabled):hover {
  background: #2a3a2a;
  border-color: #4a5a4a;
  color: #6fbf73;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: default;
}

.main-area {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.editor-wrap {
  flex: 1;
  position: relative;
}

.empty-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #555;
  font-size: 14px;
}

.toast {
  position: fixed;
  bottom: 20px;
  left: 50%;
  transform: translateX(-50%);
  background: #c44;
  color: #fff;
  padding: 8px 16px;
  border-radius: 6px;
  font-size: 13px;
  cursor: pointer;
  z-index: 100;
}
</style>
