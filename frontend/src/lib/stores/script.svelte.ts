import { bridge } from '$lib/services/wails-bridge';
import type { domain, sql } from '../../../wailsjs/go/models';

class ScriptStore {
	initialized = $state<boolean>(false);
	config = $state<domain.SQLConfig>();
	loading = $state<boolean>(false);

	result = $state<{
		state: 'idle' | 'running' | 'completed' | 'error';
		data?: sql.Result;
		error?: string;
	}>({ state: 'idle' });

	async execute(script: string) {
		if (!script.trim()) return;

		const id = this.config?.id;
		if (!id) return;

		this.result = { state: 'running' };
		this.loading = true;

		try {
			const res = await bridge.script.executeSQL(id, script);

			if (res.success) {
				this.result = { state: 'completed', data: Object.freeze(res.data) };
			} else {
				this.result = { state: 'error', error: res.error };
			}
		} catch (err) {
			this.result = { state: 'error', error: String(err) };
		} finally {
			this.loading = false;
		}
	}

	async init() {
		try {
			const res = await bridge.script.get();

			if (res.success) {
				this.config = res.data;
			}

			this.initialized = true;
		} catch (error) {
			console.error('Failed to initialize ScriptStore:', error);
			return;
		}
	}

	reset() {
		this.result = { state: 'idle', data: undefined, error: undefined };
	}
}

export const scriptStore = new ScriptStore();
