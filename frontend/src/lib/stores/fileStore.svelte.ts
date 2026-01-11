import type { FileEntry } from '$lib/types/file';

type FileState = {
	currentFile: FileEntry | null;
	content: string;
	saveStatus: 'saved' | 'saving' | 'unsaved';
};

let fileState = $state<FileState>({
	currentFile: null,
	content: '',
	saveStatus: 'saved'
});

export const fileStore = {
	get currentFile() {
		return fileState.currentFile;
	},
	get content() {
		return fileState.content;
	},
	get saveStatus() {
		return fileState.saveStatus;
	},
	get isHurlFile() {
		if (!fileState.currentFile) return false;
		const ext = fileState.currentFile.name.split('.').pop()?.toLowerCase();
		return ext === 'hurl';
	},
	get isMarkdownFile() {
		if (!fileState.currentFile) return false;
		const ext = fileState.currentFile.name.split('.').pop()?.toLowerCase();
		return ext === 'md' || ext === 'markdown';
	},
	setCurrentFile(file: FileEntry | null) {
		fileState.currentFile = file;
	},
	setContent(content: string) {
		fileState.content = content;
	},
	setSaveStatus(status: 'saved' | 'saving' | 'unsaved') {
		fileState.saveStatus = status;
	},
	clear() {
		fileState.currentFile = null;
		fileState.content = '';
		fileState.saveStatus = 'saved';
	}
};
