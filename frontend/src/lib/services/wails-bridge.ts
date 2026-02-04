import { Quit, BrowserOpenURL, WindowMinimise } from '../../../wailsjs/runtime/runtime';
import {
	StartService,
	StopService,
	GetServices,
	OpenExplorer,
	GetDirectories,
	Install,
	GetResourceMetrics,
	ExecuteSQLScript,
	GetSqlConfig
} from '../../../wailsjs/go/app/App';
import type { domain } from '../../../wailsjs/go/models';
import { safeCall } from '$lib/utils/result';

export const WailsBridge = {
	quit: () => Quit(),
	openURL: (url: string) => BrowserOpenURL(url),
	minimiseWindow: () => WindowMinimise(),
	GetResourceMetrics: (resourceName: string) => safeCall(GetResourceMetrics(resourceName)),
	openExplorer: (id: string) => safeCall(OpenExplorer(id)),
	startService: (id: string) => safeCall(StartService(id)),
	stopService: (id: string) => safeCall(StopService(id)),
	installService: (id: string, files: domain.InstallFileDTO[]) => safeCall(Install(id, files)),
	fetchServices: () => safeCall(GetServices()),
	fetchDirectories: () => safeCall(GetDirectories()),
	fetchSQLConfig: () => safeCall(GetSqlConfig()),
	executeSQLScript: (id: string, script: string) => safeCall(ExecuteSQLScript(id, script))
};

export const bridge = {
	quit: WailsBridge.quit,
	openURL: WailsBridge.openURL,
	minimiseWindow: WailsBridge.minimiseWindow,
	service: {
		get: WailsBridge.fetchServices,
		start: WailsBridge.startService,
		stop: WailsBridge.stopService,
	},
	directory: {
		get: WailsBridge.fetchDirectories,
	},
	openExplorer: WailsBridge.openExplorer,
	install: WailsBridge.installService,
	metrics: WailsBridge.GetResourceMetrics,
	script: {
		get: WailsBridge.fetchSQLConfig,
		executeSQL: WailsBridge.executeSQLScript,
	}
}


