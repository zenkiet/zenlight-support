import { CheckUpdate, DoUpdate } from '../../../wailsjs/go/app/App';
import type { app } from '../../../wailsjs/go/models';

export interface Setting {
	runBackground: boolean;
	notifications: boolean;
	runOnStartup: boolean;
	autoCheckUpdate: boolean;
}

export class SystemStore {
	isDark = $state(true);
	setting = $state<Setting>({
		runBackground: false,
		notifications: false,
		runOnStartup: false,
		autoCheckUpdate: false
	});
	initialized = $state<boolean>(false);
	info = $state<app.UpdateInfo>({
		available: false,
		currentVersion: 'v0.0.0',
		build: '',
		latestVersion: '',
		releaseNotes: '',
		downloadUrl: ''
	});

	constructor() {
		const saved = localStorage.getItem('theme');
		if (saved) {
			this.isDark = saved === 'dark';
		} else {
			this.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
		}
		this.applyTheme();
	}

	async init() {
		if (this.initialized) return;
		try {
			this.info = await this.getVersion();
			this.initialized = true;
		} catch (err) {
			console.error('Error loading system info:', err);
		}
	}

	async getVersion() {
		const version = await CheckUpdate();
		return version;
	}

	async update() {
		const url = this.info?.downloadUrl ?? '';
		if (!url) return;

		await DoUpdate(url);
	}

	toggleTheme() {
		this.isDark = !this.isDark;
		this.applyTheme();
	}

	applyTheme() {
		if (this.isDark) {
			document.documentElement.setAttribute('data-theme', 'dark');
			localStorage.setItem('theme', 'dark');
		} else {
			document.documentElement.setAttribute('data-theme', 'light');
			localStorage.setItem('theme', 'light');
		}
	}
}

export const systemStore = new SystemStore();
