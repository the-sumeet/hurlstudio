/**
 * HTTP utility functions
 */

export type HttpStatusRange = '1xx' | '2xx' | '3xx' | '4xx' | '5xx';

export function getStatusRange(status: number): HttpStatusRange {
    if (status >= 100 && status < 200) return '1xx';
    if (status >= 200 && status < 300) return '2xx';
    if (status >= 300 && status < 400) return '3xx';
    if (status >= 400 && status < 500) return '4xx';
    return '5xx';
}

export function isSuccessStatus(status: number): boolean {
    return status >= 200 && status < 300;
}

export function isRedirectStatus(status: number): boolean {
    return status >= 300 && status < 400;
}

export function isClientErrorStatus(status: number): boolean {
    return status >= 400 && status < 500;
}

export function isServerErrorStatus(status: number): boolean {
    return status >= 500;
}
