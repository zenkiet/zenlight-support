<script lang="ts">
	import { onMount, onDestroy } from 'svelte';
	import { EditorState, Compartment } from '@codemirror/state';
	import {
		EditorView,
		keymap,
		lineNumbers,
		highlightActiveLineGutter,
		highlightSpecialChars,
		drawSelection,
		dropCursor,
		rectangularSelection,
		crosshairCursor,
		highlightActiveLine
	} from '@codemirror/view';
	import { defaultKeymap, history, historyKeymap } from '@codemirror/commands';
	import { sql } from '@codemirror/lang-sql';
	import { HighlightStyle, syntaxHighlighting } from '@codemirror/language';
	import { tags as t } from '@lezer/highlight';
	import { systemStore } from '$lib/stores/system.svelte';

	let {
		value,
		onchange,
		readonly = false
	} = $props<{
		value: string;
		onchange?: (val: string) => void;
		readonly?: boolean;
	}>();

	let container = $state<Element>();
	let view: EditorView;

	const themeConf = new Compartment();
	const syntaxConf = new Compartment();
	const readOnlyConf = new Compartment();

	const getStyles = (isDark: boolean) =>
		HighlightStyle.define([
			{ tag: t.keyword, color: `${isDark ? '#c084fc' : '#7e22ce'}`, fontWeight: 'bold' },
			{ tag: t.comment, color: `${isDark ? '#94a3b8' : '#64748b'}`, fontStyle: 'italic' },
			{ tag: t.string, color: `${isDark ? '#4ade80' : '#16a34a'}` },
			{ tag: t.number, color: `${isDark ? '#facc15' : '#ca8a04'}` },
			{ tag: t.function(t.variableName), color: `${isDark ? '#60a5fa' : '#2563eb'}` },
			{ tag: t.operator, color: `${isDark ? '#94a3b8' : '#64748b'}` },
			{ tag: t.punctuation, color: `${isDark ? '#cbd5e1' : '#94a3b8'}` }
		]);

	const getTheme = (isDark: boolean) => {
		return EditorView.theme(
			{
				'&': {
					color: isDark ? '#f8fafc' : '#0f172a',
					backgroundColor: 'transparent',
					height: '100%'
				},
				'&.cm-focused': {
					outline: 'none',
					boxShadow: 'none'
				},
				'.cm-content': {
					caretColor: isDark ? '#f8fafc' : '#0f172a',
					fontFamily: "'JetBrains Mono', 'Fira Code', 'Consolas', monospace",
					fontSize: '14px',
					lineHeight: '1.6',
					padding: '8px 0'
				},
				'.cm-cursor, .cm-dropCursor': {
					borderLeftColor: isDark ? '#f8fafc' : '#0f172a',
					borderLeftWidth: '2px'
				},
				'&.cm-focused .cm-selectionBackground, ::selection': {
					backgroundColor: isDark ? 'rgba(59, 130, 246, 0.3)' : 'rgba(59, 130, 246, 0.2)'
				},
				'.cm-gutters': {
					backgroundColor: 'transparent',
					color: isDark ? '#64748b' : '#94a3b8',
					border: 'none',
					paddingRight: '8px',
					minWidth: '40px'
				},
				'.cm-activeLineGutter': {
					backgroundColor: 'transparent',
					color: isDark ? '#f8fafc' : '#0f172a'
				},
				'.cm-activeLine': {
					backgroundColor: isDark ? 'rgba(148, 163, 184, 0.1)' : 'rgba(148, 163, 184, 0.08)'
				},
				'.cm-line': {
					padding: '0 4px'
				},
				'.cm-placeholder': {
					color: isDark ? '#64748b' : '#94a3b8',
					fontStyle: 'italic'
				}
			},
			{ dark: isDark }
		);
	};

	onMount(() => {
		view = new EditorView({
			parent: container,
			state: EditorState.create({
				doc: value,
				extensions: [
					lineNumbers(),
					highlightActiveLineGutter(),
					highlightSpecialChars(),
					history(),
					drawSelection(),
					dropCursor(),
					rectangularSelection(),
					crosshairCursor(),
					highlightActiveLine(),
					keymap.of([...defaultKeymap, ...historyKeymap]),
					sql(),
					themeConf.of(getTheme(systemStore.isDark)),
					syntaxConf.of(syntaxHighlighting(getStyles(systemStore.isDark))),
					readOnlyConf.of(EditorState.readOnly.of(readonly)),
					EditorView.lineWrapping,
					EditorView.updateListener.of((u) => {
						if (u.docChanged) {
							const newValue = u.state.doc.toString();
							value = newValue;
							onchange?.(newValue);
						}
					})
				]
			})
		});
	});

	onDestroy(() => view?.destroy());

	$effect(() => {
		if (!view) return;
		view.dispatch({
			effects: [
				themeConf.reconfigure(getTheme(systemStore.isDark)),
				syntaxConf.reconfigure(syntaxHighlighting(getStyles(systemStore.isDark)))
			]
		});
	});

	$effect(() => {
		view?.dispatch({ effects: readOnlyConf.reconfigure(EditorState.readOnly.of(readonly)) });
	});

	$effect(() => {
		if (view && value !== view.state.doc.toString()) {
			view.dispatch({
				changes: { from: 0, to: view.state.doc.length, insert: value }
			});
		}
	});

	export const focus = () => view?.focus();
</script>

<div bind:this={container} class="h-full w-full bg-surface/60 border-b border-muted/15"></div>
