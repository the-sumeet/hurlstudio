/**
 * Get the color class for an HTTP status code
 * @param status - HTTP status code
 * @param prefix - CSS class prefix (e.g., 'text-', 'bg-', 'border-'). Use empty string for no prefix.
 * @returns CSS color class
 */
export function getHttpStatusColorClass(status: number, prefix: string = 'text-'): string {
	if (status >= 200 && status < 300) return `${prefix}green-600`;
	if (status >= 300 && status < 400) return `${prefix}blue-600`;
	if (status >= 400 && status < 500) return `${prefix}orange-600`;
	if (status >= 500) return `${prefix}red-600`;
	return `${prefix}gray-600`;
}
