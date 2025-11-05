<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as monaco from 'monaco-editor';
	import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';

	let {
		value = '',
		language = 'plaintext',
		theme = 'vs-dark',
		readonly = false,
		onchange = undefined
	}: {
		value?: string;
		language?: string;
		theme?: string;
		readonly?: boolean;
		onchange?: (newValue: string) => void;
	} = $props();

	let editorContainer: HTMLDivElement;
	let editor: monaco.editor.IStandaloneCodeEditor;

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
			wordWrap: 'on'
		});

		// Listen to content changes
		editor.onDidChangeModelContent(() => {
			const newValue = editor.getValue();
			value = newValue;
			if (onchange) {
				onchange(newValue);
			}
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
	$effect(() => {
		if (editor && value !== editor.getValue()) {
			editor.setValue(value);
		}
	});

	// Update language when it changes
	$effect(() => {
		if (editor) {
			const model = editor.getModel();
			if (model) {
				monaco.editor.setModelLanguage(model, language);
			}
		}
	});

	// Update theme when it changes
	$effect(() => {
		if (editor) {
			monaco.editor.setTheme(theme);
		}
	});
</script>

<div bind:this={editorContainer} class="editor-container"></div>

<style>
	.editor-container {
		width: 100%;
		height: 100%;
		min-height: 400px;
	}
</style>
