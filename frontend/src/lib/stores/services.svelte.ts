import { SvelteMap } from 'svelte/reactivity';
import { domain } from '../../../wailsjs/go/models';
import { useWailsEvent } from '$lib/services/wails-event';
import { WailsBridge } from '$lib/services/wails-bridge';

export enum Status {
	STOPPED = 1,
	RUNNING = 4
}

export interface ServiceMetrics {
	pid: number;
	createTime: number;
	cpu: number;
	mem: number;
}

export interface ServiceStatus {
	id: string;
	status: Status;
	metrics?: ServiceMetrics;
}

export interface ServiceData {
	id: string;
	name: string;
	description: string;
	status: Status;
	installable: boolean;
	metrics?: ServiceMetrics;
}

export class Service {
	/** Props */
	readonly id: string;
	readonly name: string;
	readonly description: string;
	readonly installable: boolean = false;

	/** States */
	status = $state<Status>(Status.STOPPED);
	metrics = $state<ServiceMetrics>({ pid: 0, createTime: 0, cpu: 0, mem: 0 });
	loading = $state<boolean>(true);

	constructor(config: ServiceData) {
		this.id = config.id;
		this.name = config.name;
		this.description = config.description;
		this.installable = config.installable;

		if (config.metrics) {
			this.metrics = config.metrics;
		}
	}

	update(data: ServiceStatus) {
		if (this.status !== data.status) {
			this.status = data.status;
		}

		if (data.metrics) {
			this.metrics = data.metrics;
		}

		this.loading = false;
	}

	async start() {
		if (this.loading || this.status === Status.RUNNING) return;

		const prevStatus = this.status;
		this.loading = true;

		const result = await WailsBridge.service.start(this.id);
		if (!result.success) {
			this.status = prevStatus;
			this.loading = false;
			console.error(`Failed to start ${this.name}:`, result.error);
			throw new Error(result.error);
		}
	}

	async stop() {
		if (this.loading || this.status === Status.RUNNING) return;

		const prevStatus = this.status;
		const previousMetrics = { ...this.metrics };

		this.metrics = { pid: 0, createTime: 0, cpu: 0, mem: 0 };
		this.loading = true;

		const result = await WailsBridge.service.stop(this.id);
		if (!result.success) {
			this.status = prevStatus;
			this.metrics = previousMetrics;
			this.loading = false;
			console.error(`Failed to stop ${this.name}:`, result.error);
			throw new Error(result.error);
		}
	}

	async openExplorer() {
		await WailsBridge.service.openExplorer(this.id);
	}

	async install(files: domain.InstallFileDTO[]) {
		this.loading = true;
		const result = await WailsBridge.service.install(this.id, files);
		this.loading = false;

		if (!result.success) {
			console.error(`Failed to install ${this.name}:`, result.error);
			throw result.error;
		}
	}
}

export class ServiceStore {
	private _services = new SvelteMap<string, Service>();
	services = $derived(Array.from(this._services.values()));

	/** Derived */
	totalCount = $derived(this.services.length);
	runningCount = $derived(this.services.filter((s) => s.status === Status.RUNNING).length);
	stoppedCount = $derived(this.services.filter((s) => s.status === Status.STOPPED).length);

	/** States */
	initialized = $state<boolean>(false);

	constructor() {
		useWailsEvent<ServiceStatus[]>(`services-update`, (updates) => {
			updates.forEach((update) => {
				const service = this._services.get(update.id);
				if (service) {
					service.update(update);
				}
			});
		});
	}

	async init() {
		if (this.initialized) return;

		const result = await WailsBridge.service.get();
		if (result.success) {
			this._services.clear();
			result.data.forEach((service) => {
				this._services.set(
					service.id,
					new Service({
						id: service.id,
						name: service.name,
						description: service.description,
						installable: service.installable,
						status: Status.STOPPED
					})
				);
			});
			this.initialized = true;
		} else {
			console.error('Error fetching services:', result.error);
		}
	}
}

export const serviceStore = new ServiceStore();
