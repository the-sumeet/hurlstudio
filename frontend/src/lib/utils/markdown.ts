import { marked } from 'marked';


export function createMarkdownRenderer(baseDir: string) {
    const renderer = new marked.Renderer();

    // Custom image renderer to resolve relative paths
    renderer.image = ({ href, title, text }) => {
        let resolvedHref = href || '';

        if (
            typeof href === 'string' &&
            href &&
            !href.startsWith('http://') &&
            !href.startsWith('https://') &&
            !href.startsWith('data:')
        ) {
            // Handle relative paths
            if (href.startsWith('./') || href.startsWith('../') || !href.startsWith('/')) {
                resolvedHref = resolveRelativePath(baseDir, href);
            }
        }

        const titleAttr = title ? ` title="${escapeHtml(title)}"` : '';
        const altAttr = text ? escapeHtml(text) : '';
        return `<img src="${resolvedHref}" alt="${altAttr}"${titleAttr}>`;
    };

    return renderer;
}


function resolveRelativePath(baseDir: string, relativePath: string): string {
    const parts = relativePath.split('/');
    const dirParts = baseDir.split('/');

    for (const part of parts) {
        if (part === '..') {
            dirParts.pop();
        } else if (part !== '.' && part !== '') {
            dirParts.push(part);
        }
    }

    return dirParts.join('/');
}

function escapeHtml(text: string): string {
    const map: Record<string, string> = {
        '&': '&amp;',
        '<': '&lt;',
        '>': '&gt;',
        '"': '&quot;',
        "'": '&#039;'
    };
    return text.replace(/[&<>"']/g, (char) => map[char]);
}


export async function renderMarkdown(content: string, baseDir: string): Promise<string> {
    const renderer = createMarkdownRenderer(baseDir);
    return await marked(content, { renderer });
}
