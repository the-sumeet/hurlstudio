import { toast } from 'svelte-sonner';

export class AppError extends Error {
    constructor(
        message: string,
        public readonly title: string = 'Error',
        public readonly originalError?: unknown
    ) {
        super(message);
        this.name = 'AppError';
    }
}


export function handleError(error: unknown, context?: string): void {
    console.error(context ? `[${context}]` : '', error);

    let title = 'Error';
    let message = 'An unexpected error occurred';

    if (error instanceof AppError) {
        title = error.title;
        message = error.message;
    } else if (error instanceof Error) {
        message = error.message;
    } else if (typeof error === 'string') {
        message = error;
    }

    if (context) {
        title = context;
    }

    toast.error(title, {
        description: message
    });
}


export function handleSuccess(message: string, description?: string): void {
    toast.success(message, {
        description
    });
}

/**
 * Create a typed error for file operations
 */
export function createFileError(operation: string, filename: string, error: unknown): AppError {
    const message = error instanceof Error ? error.message : String(error);
    return new AppError(message, `Failed to ${operation} "${filename}"`);
}
