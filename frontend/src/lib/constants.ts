/**
 * Application-wide constants
 */

/**
 * Supported file extensions
 */
export const SUPPORTED_EXTENSIONS = {
    MARKDOWN: ['.md', '.markdown'],
    HURL: ['.hurl'],
    ALL: ['.md', '.markdown', '.hurl']
} as const;

/**
 * HTTP methods supported by Hurl
 */
export const HTTP_METHODS = [
    'GET',
    'POST',
    'PUT',
    'DELETE',
    'PATCH',
    'HEAD',
    'OPTIONS',
    'CONNECT',
    'TRACE'
] as const;

/**
 * Monaco Editor configuration
 */
export const EDITOR_CONFIG = {
    FONT_SIZE: 14,
    TAB_SIZE: 2,
    MIN_HEIGHT: 400,
    THEMES: {
        LIGHT: 'vs',
        DARK: 'vs-dark'
    }
} as const;

/**
 * Debounce delays (in milliseconds)
 */
export const DEBOUNCE_DELAYS = {
    SEARCH: 300,
    SAVE: 500,
    RESIZE: 150
} as const;

/**
 * Toast notification durations (in milliseconds)
 */
export const TOAST_DURATION = {
    SHORT: 2000,
    MEDIUM: 4000,
    LONG: 6000
} as const;

/**
 * Local storage keys
 */
export const STORAGE_KEYS = {
    THEME: 'theme',
    LAST_OPENED_FILE: 'lastOpenedFile',
    LAST_OPENED_DIR: 'lastOpenedDir',
    ACTIVE_ENVIRONMENT: 'activeEnvironment'
} as const;


export const SIDEBAR_CONFIG = {
    WIDTH: 300,
    ICON_WIDTH: 'calc(var(--sidebar-width-icon) + 1px)'
} as const;


export const VALIDATION_MESSAGES = {
    INVALID_EXTENSION: 'Invalid file extension. File must end with .md, .markdown, or .hurl',
    EMPTY_FILENAME: 'Filename cannot be empty',
    FILE_EXISTS: 'A file with this name already exists',
    INVALID_CHARACTERS: 'Filename contains invalid characters'
} as const;

export const ERROR_MESSAGES = {
    FILE_LOAD_FAILED: 'Failed to load file',
    FILE_SAVE_FAILED: 'Failed to save file',
    FILE_DELETE_FAILED: 'Failed to delete file',
    FILE_RENAME_FAILED: 'Failed to rename file',
    FILE_CREATE_FAILED: 'Failed to create file',
    DIR_CREATE_FAILED: 'Failed to create directory',
    HURL_RUN_FAILED: 'Failed to run Hurl',
    HURL_ENTRY_RUN_FAILED: 'Failed to run Hurl entry',
    CLIPBOARD_COPY_FAILED: 'Failed to copy to clipboard',
    REPORT_LOAD_FAILED: 'Failed to load report'
} as const;


export const SUCCESS_MESSAGES = {
    FILE_SAVED: 'File saved successfully',
    FILE_DELETED: 'File deleted successfully',
    FILE_RENAMED: 'File renamed successfully',
    FILE_CREATED: 'File created successfully',
    DIR_CREATED: 'Directory created successfully',
    HURL_EXECUTED: 'Hurl executed successfully',
    COPIED_TO_CLIPBOARD: 'Copied to clipboard'
} as const;


export const CONTENT_TYPE_TO_LANGUAGE: Record<string, string> = {
    'application/json': 'json',
    'text/html': 'html',
    'application/xml': 'xml',
    'text/xml': 'xml',
    'text/javascript': 'javascript',
    'application/javascript': 'javascript',
    'text/css': 'css',
    'application/x-yaml': 'yaml',
    'text/yaml': 'yaml',
    'text/plain': 'plaintext'
};


export const TOOLTIP_DELAY = 500;
export const ANIMATION_DURATION = {
    FAST: 150,
    NORMAL: 300,
    SLOW: 500
} as const;
