<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import * as monaco from 'monaco-editor';
	import editorWorker from 'monaco-editor/esm/vs/editor/editor.worker?worker';
	import { GetFlattenedVariables, GetActiveEnvironment } from '$lib/wailsjs/go/main/App';
	import { EDITOR_CONFIG } from '$lib/constants';

	interface Props {
		value?: string;
		language?: string;
		theme?: string;
		readonly?: boolean;
		disableFind?: boolean;
		onchange?: (newValue: string) => void;
		onRunEntry?: (entryIndex: number) => void;
	}

	let {
		value = '',
		language = 'plaintext',
		theme = 'vs-dark',
		readonly = false,
		disableFind = false,
		onchange = undefined,
		onRunEntry = undefined
	}: Props = $props();

	let editorContainer: HTMLDivElement;
	let editor: monaco.editor.IStandaloneCodeEditor;
	let codeLensProvider: monaco.IDisposable | null = null;
	let hoverProvider: monaco.IDisposable | null = null;
	let commandDisposables: monaco.IDisposable[] = [];
	let commandIdCounter = 0;

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
		const httpMethods = [
			'GET',
			'POST',
			'PUT',
			'DELETE',
			'PATCH',
			'HEAD',
			'OPTIONS',
			'CONNECT',
			'TRACE'
		];

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
			fontSize: EDITOR_CONFIG.FONT_SIZE,
			minimap: { enabled: false },
			scrollBeyondLastLine: false,
			wordWrap: 'on',
			codeLens: true,
			links: false,
			folding: true, // Enable code folding (collapse/expand blocks)
			glyphMargin: false, // Disable document symbols
			hover: {
				enabled: true
			}
		});

		// Disable Cmd+F / Ctrl+F find widget if requested
		if (disableFind) {
			editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.KeyF, () => {
				window.dispatchEvent(
					new KeyboardEvent('keydown', { key: 'f', metaKey: true, bubbles: true })
				);
			});
		}

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
			// Register for all languages - we'll filter in the provider
			codeLensProvider = monaco.languages.registerCodeLensProvider('*', {
				provideCodeLenses: (model) => {
					const text = model.getValue();
					const entryLines = findHurlEntries(text);

					// If no entries found, return empty
					if (entryLines.length === 0) {
						return { lenses: [], dispose: () => {} };
					}

					// Clear previous command disposables
					commandDisposables.forEach((d) => d.dispose());
					commandDisposables = [];

					const lenses = entryLines.map((lineNumber, index) => {
						const commandId = `run-hurl-entry-${commandIdCounter++}`;

						// Register command and store disposable
						const disposable = monaco.editor.registerCommand(commandId, () => {
							if (onRunEntry) {
								onRunEntry(index + 1);
							}
						});
						commandDisposables.push(disposable);

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
	});

	onDestroy(() => {
		if (editor) {
			editor.dispose();
		}
		if (codeLensProvider) {
			codeLensProvider.dispose();
		}
		if (hoverProvider) {
			hoverProvider.dispose();
		}
		// Dispose all command registrations
		commandDisposables.forEach((d) => d.dispose());
		commandDisposables = [];
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

	// Manage hover provider based on language
	$effect(() => {
		if (!editor) return;

		// Dispose existing hover provider when language changes
		if (hoverProvider) {
			hoverProvider.dispose();
			hoverProvider = null;
		}

		// Only register hover provider for plaintext (Hurl files)
		if (language === 'plaintext') {
			hoverProvider = monaco.languages.registerHoverProvider('plaintext', {
				provideHover: async (model, position) => {
					const line = model.getLineContent(position.lineNumber);

					// Find if we're hovering over a variable {{...}}
					const lineUpToPosition = line.substring(0, position.column - 1);
					const lineFromPosition = line.substring(position.column - 1);

					// Check if there's {{ before and }} after the current position
					const beforeMatch = lineUpToPosition.match(/\{\{([^}]*)$/);
					const afterMatch = lineFromPosition.match(/^([^}]*)\}\}/);

					if (!beforeMatch || !afterMatch) return null;

					// Extract the full variable name
					const variableName = (beforeMatch[1] + afterMatch[1]).trim();

					if (!variableName) return null;

					try {
						// Get the active environment and flattened variables
						const activeEnv = await GetActiveEnvironment();
						const variables = await GetFlattenedVariables(activeEnv);

						if (variables && variables[variableName]) {
							return {
								range: new monaco.Range(
									position.lineNumber,
									position.column - beforeMatch[1].length - 2,
									position.lineNumber,
									position.column + afterMatch[1].length + 2
								),
								contents: [{ value: `**${variableName}**: \`${variables[variableName]}\`` }]
							};
						} else {
							return {
								range: new monaco.Range(
									position.lineNumber,
									position.column - beforeMatch[1].length - 2,
									position.lineNumber,
									position.column + afterMatch[1].length + 2
								),
								contents: [{ value: `**${variableName}**: _Variable not defined_` }]
							};
						}
					} catch (error) {
						console.error('Error fetching environment variables:', error);
						return null;
					}
				}
			});
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
