// Petit client pour l'API du backend Go.
const BASE = 'http://localhost:8080/api';

export interface Todo {
	id: number;
	title: string;
	done: boolean;
}

export async function listTodos(): Promise<Todo[]> {
	const res = await fetch(`${BASE}/todos`);
	if (!res.ok) throw new Error('Impossible de charger les tâches');
	return res.json();
}

export async function addTodo(title: string): Promise<Todo> {
	const res = await fetch(`${BASE}/todos`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ title })
	});
	if (!res.ok) throw new Error("Impossible d'ajouter la tâche");
	return res.json();
}

export async function toggleTodo(id: number, done: boolean): Promise<Todo> {
	const res = await fetch(`${BASE}/todos/${id}`, {
		method: 'PUT',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({ done })
	});
	if (!res.ok) throw new Error('Impossible de mettre à jour la tâche');
	return res.json();
}

export async function deleteTodo(id: number): Promise<void> {
	const res = await fetch(`${BASE}/todos/${id}`, { method: 'DELETE' });
	if (!res.ok) throw new Error('Impossible de supprimer la tâche');
}
