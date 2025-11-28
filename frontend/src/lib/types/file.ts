/**
 * Represents a file or directory entry in the file system
 */
export type FileEntry = {
	name: string;
	path: string;
	isDirectory: boolean;
	size: number;
	modTime: string;
};

/**
 * Represents the current files state including directory and file
 */
export type CurrentFilesState = {
	currentDir?: FileEntry;
	currentFile?: FileEntry;
};
