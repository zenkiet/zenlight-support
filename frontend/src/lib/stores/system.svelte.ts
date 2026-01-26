import { CheckUpdate, DoUpdate } from '../../../wailsjs/go/app/App';
import type { app } from '../../../wailsjs/go/models';

export interface Setting {
	runBackground: boolean;
	notifications: boolean;
	runOnStartup: boolean;
	autoCheckUpdate: boolean;
}

export class SystemStore {
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
}

export const systemStore = new SystemStore();
