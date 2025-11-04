import { browser } from '$app/environment';

type Theme = 'light' | 'dark';

// Get initial theme from localStorage or default to dark
function getInitialTheme(): Theme {
	if (!browser) return 'dark';

	const stored = localStorage.getItem('theme') as Theme | null;
	if (stored) return stored;

	// Check system preference
	if (window.matchMedia('(prefers-color-scheme: dark)').matches) {
		return 'dark';
	}

	return 'light';
}

let theme = $state<Theme>(getInitialTheme());

// Apply theme to document
function applyTheme(newTheme: Theme) {
	if (!browser) return;

	const root = document.documentElement;
	if (newTheme === 'dark') {
		root.classList.add('dark');
	} else {
		root.classList.remove('dark');
	}

	localStorage.setItem('theme', newTheme);
}

// Apply initial theme
if (browser) {
	applyTheme(theme);
}

export const themeStore = {
	get current() {
		return theme;
	},
	toggle() {
		theme = theme === 'dark' ? 'light' : 'dark';
		applyTheme(theme);
	},
	set(newTheme: Theme) {
		theme = newTheme;
		applyTheme(theme);
	}
};
