export type Command = {
	id: string;
	label: string;
	icon?: string;
	action: () => void;
	group: 'Service' | 'System' | 'Navigation';
};

class CommandRegistry {
	commands = $state<Command[]>([]);

	register(cmd: Command) {
		this.commands.push(cmd);
	}

	search(query: string) {
		return this.commands.filter((c) => c.label.toLowerCase().includes(query.toLowerCase()));
	}
}
export const cmdStore = new CommandRegistry();
