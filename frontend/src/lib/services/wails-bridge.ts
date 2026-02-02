import { Quit, BrowserOpenURL, WindowMinimise } from '../../../wailsjs/runtime/runtime';
import {
	StartService,
	StopService,
	GetServices,
	InstallService,
	OpenExplorer
} from '../../../wailsjs/go/app/App';
import type { domain } from '../../../wailsjs/go/models';
import { safeCall } from '$lib/utils/result';

export const WailsBridge = {
	openUrl: (url: string) => BrowserOpenURL(url),
	quit: () => Quit(),
	minimize: () => WindowMinimise(),
	service: {
		get: () => safeCall(GetServices()),
		start: (id: string) => safeCall(StartService(id)),
		stop: (id: string) => safeCall(StopService(id)),
		install: (id: string, files: domain.InstallFileDTO[]) => safeCall(InstallService(id, files)),
		openExplorer: (id: string) => safeCall(OpenExplorer(id))
	}
};
