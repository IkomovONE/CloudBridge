import { writable } from 'svelte/store';

export type ToastType = 'info' | 'success' | 'error';

export interface Toast {
	id: string;
	message: string;
	type: ToastType;
}

export const toasts = writable<Toast[]>([]);

export function addToast(
	message: string,
	type: ToastType = 'info',
	duration = 3000
): void {
	const id = crypto.randomUUID();

	const toast: Toast = { id, message, type };

	toasts.update((all) => [...all, toast]);

	setTimeout(() => {
		toasts.update((all) => all.filter((t) => t.id !== id));
	}, duration);
}

export function removeToast(id: string): void {
	toasts.update((all) => all.filter((t) => t.id !== id));
}