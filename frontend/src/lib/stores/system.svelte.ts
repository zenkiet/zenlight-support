import { CheckUpdate, DoUpdate } from '../../../wailsjs/go/app/App';
import type { app } from '../../../wailsjs/go/models';

export interface Setting {
	runBackground: boolean;
	notifications: boolean;
	runOnStartup: boolean;
	autoCheckUpdate: boolean;
}

export type ThemeScheme = 'light' | 'dark' | 'system';

export class SystemStore {
	scheme = $state<ThemeScheme>('system');

	isDark = $derived.by(() => {
		if (this.scheme === 'system') {
			return typeof window !== 'undefined'
				? window.matchMedia('(prefers-color-scheme: dark)').matches
				: false;
		}
		return this.scheme === 'dark';
	});

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

	private schemeMediaQuery: MediaQueryList | null = null;

	constructor() {
		const saved = localStorage.getItem('theme') as ThemeScheme | null;
		if (saved === 'light' || saved === 'dark' || saved === 'system') {
			this.scheme = saved;
		}

		if (typeof window !== 'undefined') {
			window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', () => {
				if (this.scheme === 'system') {
					this.applyTheme();
				}
			});
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
		const order: ThemeScheme[] = ['light', 'dark', 'system'];
		const current = order.indexOf(this.scheme);
		this.scheme = order[(current + 1) % 3];
		this.applyTheme();
	}

	setScheme(scheme: ThemeScheme) {
		this.scheme = scheme;
		this.applyTheme();
	}

	private applyTheme() {
		document.documentElement.setAttribute('data-theme', this.isDark ? 'dark' : 'light');
		localStorage.setItem('theme', this.scheme);
	}
}

export const systemStore = new SystemStore();
