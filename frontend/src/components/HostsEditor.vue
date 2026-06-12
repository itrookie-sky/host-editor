<template>
  <div ref="containerRef" class="editor-container"></div>
</template>

<script lang="ts" setup>
import { ref, onMounted, onBeforeUnmount, watch } from 'vue'
import * as monaco from 'monaco-editor'

const props = defineProps<{
  modelValue: string
  readOnly?: boolean
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

const containerRef = ref<HTMLElement>()
let editor: monaco.editor.IStandaloneCodeEditor | null = null
let ignoreNextChange = false

onMounted(() => {
  if (!containerRef.value) return

  editor = monaco.editor.create(containerRef.value, {
    value: props.modelValue,
    language: 'plaintext',
    theme: 'vs-dark',
    readOnly: props.readOnly ?? false,
    minimap: { enabled: false },
    fontSize: 13,
    lineNumbers: 'on',
    scrollBeyondLastLine: false,
    wordWrap: 'on',
    automaticLayout: true,
    padding: { top: 8, bottom: 8 },
    renderLineHighlight: 'line',
    overviewRulerBorder: false,
  })

  editor.onDidChangeModelContent(() => {
    if (ignoreNextChange) {
      ignoreNextChange = false
      return
    }
    emit('update:modelValue', editor!.getValue())
  })
})

watch(() => props.modelValue, (val) => {
  if (!editor) return
  if (editor.getValue() === val) return
  ignoreNextChange = true
  editor.setValue(val)
})

watch(() => props.readOnly, (ro) => {
  editor?.updateOptions({ readOnly: ro ?? false })
})

onBeforeUnmount(() => {
  editor?.dispose()
})
</script>

<style scoped>
.editor-container {
  width: 100%;
  height: 100%;
}
</style>
