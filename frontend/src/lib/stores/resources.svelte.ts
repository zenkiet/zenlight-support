export interface ResourceMetrics {
	pid: number;
	createTime: number;
	cpu: number;
	mem: number;
	// --- DIRECTORY METRICS ---
	totalSize?: number;
	lastModified?: number;
}

import { ServiceStatus } from '../enums/service-status.enum';
import { domain } from '../../../wailsjs/go/models';
import { ResourceType } from '../enums/resource-type.enum';
import { bridge } from '$lib/services/wails-bridge';
import { useWailsEvent } from '$lib/services/wails-event';
import { SvelteMap } from 'svelte/reactivity';

export class Resource {
	/** States */
	config = $state<domain.ResourceConfig>();
	status = $state<ServiceStatus>(ServiceStatus.STOPPED);
	metrics = $state<ResourceMetrics>({
		pid: 0,
		createTime: 0,
		cpu: 0,
		mem: 0
	});
	loading = $state<boolean>(true);

	constructor(config: domain.ResourceConfig) {
		this.config = config;

		if (this.isDirectory) {
			this.getMetrics();
		}
	}

	/** Getter */
	get isService() {
		return this.config?.type === ResourceType.SERVICE;
	}

	get isDirectory() {
		return this.config?.type === ResourceType.DIRECTORY;
	}

	get isRunning() {
		return this.status === ServiceStatus.RUNNING;
	}

	get isStopped() {
		return this.status === ServiceStatus.STOPPED;
	}

	async start() {
		if (!this.isService) throw new Error('Cannot start a directory');
		if (this.loading || this.isRunning) return;

		const prevStatus = this.status;
		this.loading = true;

		const res = await bridge.service.start(this.config!.id);
		this.loading = !res.success;
		if (!res.success) {
			this.status = prevStatus;
			console.error(`Failed to start ${this.config!.name}:`, res.error);
			throw new Error(res.error);
		}
	}

	async stop() {
		if (!this.isService) throw new Error('Cannot stop a directory');
		if (this.loading || this.isStopped) return;

		const prevStatus = this.status;
		this.loading = true;

		const res = await bridge.service.stop(this.config!.id);
		this.loading = !res.success;
		if (!res.success) {
			this.status = prevStatus!;
			console.error(`Failed to stop ${this.config!.name}:`, res.error);
			throw new Error(res.error);
		}
	}

	async openExplorer() {
		const res = await bridge.openExplorer(this.config!.id);
		if (!res.success) throw new Error(res.error);
	}

	async install(files: domain.InstallFileDTO[]) {
		this.loading = true;
		const result = await bridge.install(this.config!.id, files);
		this.loading = false;

		if (!result.success) {
			console.error(`Failed to install ${this.config!.name}:`, result.error);
			throw result.error;
		} else {
			await this.getMetrics();
		}
	}

	async getMetrics() {
		const res = await bridge.metrics(this.config!.id);
		if (!res.success) {
			throw new Error(res.error);
		}
		this.metrics = res.data;
	}

	async save(data: Partial<domain.ResourceConfig>) {
		this.loading = true;
		const res = await bridge.createOrUpdate({ ...this.config!, ...data });
		this.loading = false;

		if (!res.success) {
			console.error(`Failed to save ${this.config!.name}:`, res.error);
			throw new Error(res.error);
		} else {
			this.config = res.data;
		}
	}

	async delete() {
		this.loading = true;
		const res = await bridge.delete(this.config!.id);
		this.loading = false;

		if (!res.success) {
			console.error(`Failed to delete ${this.config!.name}:`, res.error);
			throw new Error(res.error);
		}
	}

	update(payload: { status?: ServiceStatus; metrics?: Partial<ResourceMetrics> }) {
		if (payload.status) {
			this.status = payload.status;
			this.metrics = {
				pid: 0,
				createTime: 0,
				cpu: 0,
				mem: 0
			};
		}
		if (payload.metrics && payload.status === ServiceStatus.RUNNING) {
			this.metrics = { ...this.metrics, ...payload.metrics };
		}
		this.loading = false;
	}
}

export class ResourceStore {
	initialized = $state<boolean>(false);
	private _items = new SvelteMap<string, Resource>();
	items = $derived(Array.from(this._items.values()));

	services = $derived(this.items.filter((r) => r.isService));
	directories = $derived(this.items.filter((r) => r.isDirectory));

	constructor() {
		useWailsEvent<
			{
				id: string;
				status: ServiceStatus;
				metrics?: ResourceMetrics;
			}[]
		>(`services-update`, (updates) => {
			updates.forEach((u) => {
				const item = this._items.get(u.id);
				if (item) item.update({ status: u.status, metrics: u.metrics });
			});
		});
	}

	async init() {
		try {
			const res = await Promise.all([bridge.service.get(), bridge.directory.get()]);
			const [svcRes, dirRes] = res;

			if (svcRes.success) {
				svcRes.data.forEach((cfg) => this._items.set(cfg.id, new Resource(cfg)));
			}

			if (dirRes.success) {
				dirRes.data.forEach((cfg) => this._items.set(cfg.id, new Resource(cfg)));
			}
			this.initialized = true;
		} catch (error) {
			console.error('Failed to initialize ResourceStore:', error);
			return;
		}
	}

	get(id: string) {
		return this._items.get(id);
	}
}

export const resourceStore = new ResourceStore();
