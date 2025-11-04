type FileEntry = {
	name: string;
	path: string;
	isDirectory: boolean;
	size: number;
	modTime: string;
};

type FileState = {
	currentFile: FileEntry | null;
	content: string;
};

let fileState = $state<FileState>({
	currentFile: null,
	content: ''
});

export const fileStore = {
	get currentFile() {
		return fileState.currentFile;
	},
	get content() {
		return fileState.content;
	},
	setCurrentFile(file: FileEntry | null) {
		fileState.currentFile = file;
	},
	setContent(content: string) {
		fileState.content = content;
	},
	clear() {
		fileState.currentFile = null;
		fileState.content = '';
	}
};
