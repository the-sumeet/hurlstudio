<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import * as monaco from 'monaco-editor';
  import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';

  let editorContainer: HTMLDivElement;
  let editor: monaco.editor.IStandaloneCodeEditor;

  export let value = '';
  export let language = 'plaintext';
  export let theme = 'vs-dark';
  export let readonly = false;

  // Setup Monaco Editor workers
  self.MonacoEnvironment = {
    getWorker(_: any, label: string) {
      return new editorWorker();
    }
  };

  onMount(() => {
    // Create the editor
    editor = monaco.editor.create(editorContainer, {
      value: value,
      language: language,
      theme: theme,
      readOnly: readonly,
      automaticLayout: true,
      fontSize: 14,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      wordWrap: 'on',
    });

    // Listen to content changes
    editor.onDidChangeModelContent(() => {
      value = editor.getValue();
    });

    return () => {
      editor.dispose();
    };
  });

  onDestroy(() => {
    if (editor) {
      editor.dispose();
    }
  });

  // Update editor when value changes externally
  $: if (editor && value !== editor.getValue()) {
    editor.setValue(value);
  }
</script>

<div bind:this={editorContainer} class="editor-container"></div>

<style>
  .editor-container {
    width: 100%;
    height: 100%;
    min-height: 400px;
  }
</style>
