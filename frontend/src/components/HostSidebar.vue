<template>
  <aside class="sidebar">
    <div class="sidebar-header">
      <span class="sidebar-title">Versions</span>
      <button class="btn-icon" @click="$emit('create')" title="New version">+</button>
    </div>
    <div class="sidebar-list">
      <div
        v-for="file in files"
        :key="file.name"
        class="sidebar-item"
        :class="{ active: file.name === activeName }"
        @click="$emit('select', file.name)"
      >
        <span class="item-name">{{ file.name }}</span>
        <span v-if="file.name === activeName && isDirty" class="dirty-dot" title="Unsaved changes"></span>
        <button
          v-if="file.name !== 'hosts'"
          class="btn-delete"
          @click.stop="$emit('delete', file.name)"
          title="Delete"
        >×</button>
      </div>
    </div>
  </aside>
</template>

<script lang="ts" setup>
import type { HostFileInfo } from '../types/host'

defineProps<{
  files: HostFileInfo[]
  activeName: string
  isDirty: boolean
}>()

defineEmits<{
  (e: 'select', name: string): void
  (e: 'create'): void
  (e: 'delete', name: string): void
}>()
</script>

<style scoped>
.sidebar {
  width: 200px;
  min-width: 160px;
  background: #1a1f2b;
  display: flex;
  flex-direction: column;
  border-right: 1px solid #2a2f3b;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 10px 12px;
  border-bottom: 1px solid #2a2f3b;
}

.sidebar-title {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.5px;
  color: #888;
}

.btn-icon {
  width: 22px;
  height: 22px;
  border: none;
  border-radius: 4px;
  background: #2a2f3b;
  color: #aaa;
  font-size: 14px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-icon:hover {
  background: #3a3f4b;
  color: #fff;
}

.sidebar-list {
  flex: 1;
  overflow-y: auto;
}

.sidebar-item {
  display: flex;
  align-items: center;
  padding: 8px 12px;
  cursor: pointer;
  font-size: 13px;
  color: #ccc;
  gap: 6px;
}

.sidebar-item:hover {
  background: #252a36;
}

.sidebar-item.active {
  background: #2a3a2a;
  color: #6fbf73;
}

.item-name {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.dirty-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: #e6a817;
  flex-shrink: 0;
}

.btn-delete {
  border: none;
  background: none;
  color: #666;
  font-size: 16px;
  cursor: pointer;
  padding: 0 2px;
  line-height: 1;
  display: none;
}

.sidebar-item:hover .btn-delete {
  display: block;
}

.btn-delete:hover {
  color: #e55;
}
</style>
