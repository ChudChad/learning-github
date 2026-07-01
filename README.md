# ToDo — Go + SvelteKit

Une simple application ToDo. Backend en Go (API REST, stockage en mémoire), frontend en SvelteKit.

## Structure

- `backend/` — API REST en Go
- `todo/` — frontend SvelteKit

## Lancer le backend

```bash
cd backend
go run .
```

L'API écoute sur http://localhost:8080

### Endpoints

| Méthode | Route             | Description        |
| ------- | ----------------- | ------------------ |
| GET     | `/api/todos`      | Lister les tâches  |
| POST    | `/api/todos`      | Créer une tâche    |
| PUT     | `/api/todos/{id}` | Modifier une tâche |
| DELETE  | `/api/todos/{id}` | Supprimer une tâche |

## Lancer le frontend

```bash
cd todo
npm install
npm run dev
```

Le frontend démarre sur http://localhost:5173 et appelle le backend sur le port 8080.
Lance les deux en parallèle.
