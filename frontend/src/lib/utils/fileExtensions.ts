/**
 * Utility functions for file extension handling
 */

export type FileExtension = 'md' | 'markdown' | 'hurl' | 'unknown';


export function getFileExtension(filename: string): string {
    return filename.split('.').pop()?.toLowerCase() || '';
}

export function isMarkdownFile(filename: string): boolean {
    const ext = getFileExtension(filename);
    return ext === 'md' || ext === 'markdown';
}


export function isHurlFile(filename: string): boolean {
    return getFileExtension(filename) === 'hurl';
}


export function getEditorLanguage(filename: string): string {
    const ext = getFileExtension(filename);

    if (ext === 'md' || ext === 'markdown') return 'markdown';
    if (ext === 'hurl') return 'plaintext'; // Can add custom Hurl syntax later
    if (ext === 'json') return 'json';
    if (ext === 'html') return 'html';
    if (ext === 'css') return 'css';
    if (ext === 'js' || ext === 'javascript') return 'javascript';
    if (ext === 'ts' || ext === 'typescript') return 'typescript';
    if (ext === 'xml') return 'xml';
    if (ext === 'yaml' || ext === 'yml') return 'yaml';

    return 'plaintext';
}


export function hasValidExtension(filename: string): boolean {
    const validExtensions = ['.md', '.markdown', '.hurl'];
    return validExtensions.some((ext) => filename.toLowerCase().endsWith(ext));
}


export function getExtensionValidationError(): string {
    return 'Invalid file extension. File must end with .md, .markdown, or .hurl';
}
