import { SvelteMap } from 'svelte/reactivity';
import {
	GetConfig,
	StartService,
	StopService,
	OpenExplorer,
	InstallService
} from '../../../wailsjs/go/app/App';
import { domain } from '../../../wailsjs/go/models';
import { EventsOn } from '../../../wailsjs/runtime';

export enum Status {
	ERROR = 0,
	STOPPED = 1,
	START_PENDING = 2,
	STOP_PENDING = 3,
	RUNNING = 4,
	CONTINUE_PENDING = 5,
	PAUSE_PENDING = 6,
	PAUSED = 7
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
	id: string;
	name: string;
	description: string;
	installable: boolean = false;

	status = $state<Status>(Status.STOPPED);
	metrics = $state<ServiceMetrics>({
		pid: 0,
		createTime: 0,
		cpu: 0,
		mem: 0
	});
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
		this.status = data.status;
		if (data.metrics) {
			this.metrics = data.metrics;
		}

		if (this.loading && [Status.RUNNING, Status.STOPPED].includes(data.status)) {
			this.loading = false;
		}
	}

	async start() {
		this.loading = true;
		await StartService(this.id);
	}

	async stop() {
		this.loading = true;
		await StopService(this.id).then(() => {
			this.metrics = {
				pid: 0,
				createTime: 0,
				cpu: 0,
				mem: 0
			};
		});
	}

	async openExplorer() {
		await OpenExplorer(this.id);
	}

	async install(files: domain.InstallFileDTO[]) {
		await InstallService(this.id, files);
	}
}

export class ServiceStore {
	private _services = new SvelteMap<string, Service>();
	services = $derived(Array.from(this._services.values()));

	totalCount = $derived(this.services.length);
	runningCount = $derived(
		this.services.filter((service) => service.status === Status.RUNNING).length
	);
	stoppedCount = $derived(
		this.services.filter((service) => service.status === Status.STOPPED).length
	);

	initialized = $state<boolean>(false);

	async init() {
		if (this.initialized) return;
		try {
			const configs = await GetConfig();

			const sMap = new SvelteMap<string, Service>();
			configs.services.forEach((config) => {
				sMap.set(
					config.id,
					new Service({
						id: config.id,
						name: config.name,
						description: config.description,
						installable: config.installable,
						status: Status.STOPPED
					})
				);
			});
			this._services = sMap;
			this.initialized = true;

			EventsOn(`services-update`, (updates: ServiceStatus[]) => {
				updates.forEach((update) => {
					const service = this._services.get(update.id);
					if (service) {
						service.update(update);
					}
				});
			});
		} catch (err) {
			console.error('Error loading config in ServiceStore:', err);
		}
	}
}

export const serviceStore = new ServiceStore();
