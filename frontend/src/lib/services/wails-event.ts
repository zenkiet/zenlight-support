import { EventsOn, EventsOff } from '../../../wailsjs/runtime';

export function useWailsEvent<T>(eventName: string, callback: (data: T) => void) {
	EventsOn(eventName, callback);

	return () => {
		EventsOff(eventName);
	};
}
