<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as monaco from 'monaco-editor';
	import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';

	let {
		value = '',
		language = 'plaintext',
		theme = 'vs-dark',
		readonly = false,
		onchange = undefined,
		onRunEntry = undefined
	}: {
		value?: string;
		language?: string;
		theme?: string;
		readonly?: boolean;
		onchange?: (newValue: string) => void;
		onRunEntry?: (entryIndex: number) => void;
	} = $props();

	let editorContainer: HTMLDivElement;
	let editor: monaco.editor.IStandaloneCodeEditor;
	let codeLensProvider: monaco.IDisposable | null = null;

	// Setup Monaco Editor workers
	self.MonacoEnvironment = {
		getWorker(_: any, label: string) {
			return new editorWorker();
		}
	};

	// Function to find Hurl entry start lines
	// Based on Hurl grammar: entry starts with method + SP + url
	function findHurlEntries(text: string): number[] {
		const lines = text.split('\n');
		const entryLines: number[] = [];
		const httpMethods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS', 'CONNECT', 'TRACE'];

		lines.forEach((line, index) => {
			const trimmed = line.trim();

			// Skip empty lines and comments
			if (!trimmed || trimmed.startsWith('#')) {
				return;
			}

			// Check if line starts with an HTTP method followed by space
			// According to Hurl grammar: method SP url
			for (const method of httpMethods) {
				if (trimmed.startsWith(method + ' ') || trimmed === method) {
					entryLines.push(index + 1); // Monaco uses 1-based line numbers
					break;
				}
			}
		});

		return entryLines;
	}

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
			codeLens: true, // Enable CodeLens
			// Disable features that aren't supported for plaintext to avoid console warnings
			links: false,
			folding: false,
			hover: {
				enabled: false
			}
		});

		// Listen to content changes
		editor.onDidChangeModelContent(() => {
			const newValue = editor.getValue();
			value = newValue;
			if (onchange) {
				onchange(newValue);
			}
		});

		// Register CodeLens provider for Hurl files
		if (onRunEntry) {
			codeLensProvider = monaco.languages.registerCodeLensProvider(language, {
				provideCodeLenses: (model) => {
					const text = model.getValue();
					const entryLines = findHurlEntries(text);

					const lenses = entryLines.map((lineNumber, index) => {
						const commandId = `run-hurl-entry-${index}`;

						// Register command for this specific entry
						monaco.editor.registerCommand(commandId, () => {
							if (onRunEntry) {
								onRunEntry(index + 1);
							}
						});

						return {
							range: {
								startLineNumber: lineNumber,
								startColumn: 1,
								endLineNumber: lineNumber,
								endColumn: 1
							},
							id: `run-entry-${index}`,
							command: {
								id: commandId,
								title: `▶ Run Entry ${index + 1}`,
								arguments: [index + 1]
							}
						};
					});

					return { lenses, dispose: () => {} };
				},
				resolveCodeLens: (model, codeLens) => codeLens
			});
		}

		return () => {
			editor.dispose();
			if (codeLensProvider) {
				codeLensProvider.dispose();
			}
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
