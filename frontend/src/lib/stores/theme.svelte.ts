class ThemeStore {
  isDark = $state(true);

  constructor() {
    const saved = localStorage.getItem('theme');
    if (saved) {
      this.isDark = saved === 'dark';
    } else {
      this.isDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    }
    this.apply();
  }

  toggle() {
    this.isDark = !this.isDark;
    this.apply();
  }

  apply() {
    if (this.isDark) {
      document.documentElement.setAttribute('data-theme', 'dark');
			localStorage.setItem('theme', 'dark');
		} else {
      document.documentElement.setAttribute('data-theme', 'light');
      localStorage.setItem('theme', 'light');
    }
  }
}

export const theme = new ThemeStore();