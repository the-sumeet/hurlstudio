/**
 * Monaco Editor utilities and helpers
 */

import * as monaco from 'monaco-editor';
import { HTTP_METHODS } from '$lib/constants';


export function findHurlEntries(text: string): number[] {
    const lines = text.split('\n');
    const entryLines: number[] = [];

    lines.forEach((line, index) => {
        const trimmed = line.trim();

        // Skip empty lines and comments
        if (!trimmed || trimmed.startsWith('#')) {
            return;
        }

        // Check if line starts with an HTTP method followed by space
        for (const method of HTTP_METHODS) {
            if (trimmed.startsWith(method + ' ') || trimmed === method) {
                entryLines.push(index + 1); // Monaco uses 1-based line numbers
                break;
            }
        }
    });

    return entryLines;
}


export class DisposableManager {
    private disposables: monaco.IDisposable[] = [];

    /**
     * Add a disposable to be managed
     */
    add(disposable: monaco.IDisposable): void {
        this.disposables.push(disposable);
    }

    /**
     * Clear all disposables
     */
    clear(): void {
        this.disposables.forEach((d) => d.dispose());
        this.disposables = [];
    }

    /**
     * Dispose all resources
     */
    dispose(): void {
        this.clear();
    }

    /**
     * Get the number of managed disposables
     */
    get count(): number {
        return this.disposables.length;
    }
}


export function createHurlCodeLensProvider(
    onRunEntry: (entryIndex: number) => void,
    disposableManager: DisposableManager
): monaco.languages.CodeLensProvider {
    let commandIdCounter = 0;

    return {
        provideCodeLenses: (model) => {
            const text = model.getValue();
            const entryLines = findHurlEntries(text);

            // Clear previous command disposables
            disposableManager.clear();

            const lenses = entryLines.map((lineNumber, index) => {
                const commandId = `run-hurl-entry-${commandIdCounter++}`;

                // Register command and store disposable
                const disposable = monaco.editor.registerCommand(commandId, () => {
                    onRunEntry(index + 1);
                });
                disposableManager.add(disposable);

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

            return { lenses, dispose: () => { } };
        },
        resolveCodeLens: (model, codeLens) => codeLens
    };
}


export function createHurlHoverProvider(
    getVariables: () => Promise<Record<string, string>>
): monaco.languages.HoverProvider {
    return {
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
                const variables = await getVariables();

                const range = new monaco.Range(
                    position.lineNumber,
                    position.column - beforeMatch[1].length - 2,
                    position.lineNumber,
                    position.column + afterMatch[1].length + 2
                );

                if (variables && variables[variableName]) {
                    return {
                        range,
                        contents: [{ value: `**${variableName}**: \`${variables[variableName]}\`` }]
                    };
                } else {
                    return {
                        range,
                        contents: [{ value: `**${variableName}**: _Variable not defined_` }]
                    };
                }
            } catch (error) {
                console.error('Error fetching environment variables:', error);
                return null;
            }
        }
    };
}


export function getLanguageFromContentType(contentType: string): string {
    const lowerContentType = contentType.toLowerCase();

    if (lowerContentType.includes('json')) return 'json';
    if (lowerContentType.includes('html')) return 'html';
    if (lowerContentType.includes('xml')) return 'xml';
    if (lowerContentType.includes('javascript') || lowerContentType.includes('ecmascript'))
        return 'javascript';
    if (lowerContentType.includes('css')) return 'css';
    if (lowerContentType.includes('yaml')) return 'yaml';

    return 'plaintext';
}
