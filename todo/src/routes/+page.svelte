<script lang="ts">
	import { onMount } from 'svelte';
	import { listTodos, addTodo, toggleTodo, deleteTodo, type Todo } from '$lib/api';

	let todos = $state<Todo[]>([]);
	let newTitle = $state('');
	let error = $state('');
	let loading = $state(true);

	const remaining = $derived(todos.filter((t) => !t.done).length);

	onMount(load);

	async function load() {
		loading = true;
		error = '';
		try {
			todos = await listTodos();
		} catch (e) {
			error = (e as Error).message;
		} finally {
			loading = false;
		}
	}

	async function handleAdd(event: SubmitEvent) {
		event.preventDefault();
		const title = newTitle.trim();
		if (!title) return;
		try {
			const created = await addTodo(title);
			todos = [...todos, created];
			newTitle = '';
		} catch (e) {
			error = (e as Error).message;
		}
	}

	async function handleToggle(todo: Todo) {
		try {
			const updated = await toggleTodo(todo.id, !todo.done);
			todos = todos.map((t) => (t.id === updated.id ? updated : t));
		} catch (e) {
			error = (e as Error).message;
		}
	}

	async function handleDelete(todo: Todo) {
		try {
			await deleteTodo(todo.id);
			todos = todos.filter((t) => t.id !== todo.id);
		} catch (e) {
			error = (e as Error).message;
		}
	}
</script>

<main class="mx-auto max-w-md p-6">
	<h1 class="mb-4 text-2xl font-bold">Ma ToDo</h1>

	<form class="mb-4 flex gap-2" onsubmit={handleAdd}>
		<input
			class="flex-1 rounded border border-gray-300 px-3 py-2"
			placeholder="Nouvelle tâche..."
			bind:value={newTitle}
		/>
		<button class="rounded bg-blue-600 px-4 py-2 font-medium text-white hover:bg-blue-700">
			Ajouter
		</button>
	</form>

	{#if error}
		<p class="mb-4 rounded bg-red-100 px-3 py-2 text-red-700">{error}</p>
	{/if}

	{#if loading}
		<p class="text-gray-500">Chargement...</p>
	{:else if todos.length === 0}
		<p class="text-gray-500">Aucune tâche pour l'instant.</p>
	{:else}
		<ul class="space-y-2">
			{#each todos as todo (todo.id)}
				<li class="flex items-center gap-3 rounded border border-gray-200 px-3 py-2">
					<input type="checkbox" checked={todo.done} onchange={() => handleToggle(todo)} />
					<span class="flex-1 {todo.done ? 'text-gray-400 line-through' : ''}">{todo.title}</span>
					<button class="text-sm text-red-600 hover:underline" onclick={() => handleDelete(todo)}>
						Supprimer
					</button>
				</li>
			{/each}
		</ul>
		<p class="mt-4 text-sm text-gray-500">{remaining} tâche(s) restante(s)</p>
	{/if}
</main>
