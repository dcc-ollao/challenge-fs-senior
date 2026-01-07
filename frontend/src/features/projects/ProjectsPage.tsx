import { useEffect, useState } from "react";
import { listProjects, createProject } from "./api";
import type { Project } from "./types";

export default function ProjectsPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [name, setName] = useState("");
  const [creating, setCreating] = useState(false);

  async function loadProjects() {
    try {
      setLoading(true);
      setError(null);
      const data = await listProjects();
      setProjects(data);
    } catch (err) {
      setError("Failed to load projects");
    } finally {
      setLoading(false);
    }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!name.trim()) return;

    try {
      setCreating(true);
      await createProject(name.trim());
      setName("");
      await loadProjects();
    } catch (err) {
      setError("Failed to create project");
    } finally {
      setCreating(false);
    }
  }

  useEffect(() => {
    loadProjects();
  }, []);

  return (
    <div className="max-w-3xl mx-auto px-4 py-6">
      <h1 className="text-2xl font-semibold">Projects</h1>
      <p className="text-sm text-gray-500 mb-6">Your projects</p>

      <form onSubmit={handleCreate} className="mb-6 flex gap-2">
        <input
          type="text"
          placeholder="Project name"
          className="flex-1 border rounded px-3 py-2 text-sm"
          value={name}
          onChange={(e) => setName(e.target.value)}
          disabled={creating}
        />
        <button
          type="submit"
          className="bg-black text-white px-4 py-2 rounded text-sm disabled:opacity-50"
          disabled={creating}
        >
          Create
        </button>
      </form>

      {loading && <p className="text-sm text-gray-500">Loading projects…</p>}

      {error && <p className="text-sm text-red-600">{error}</p>}

      {!loading && !error && projects.length === 0 && (
        <p className="text-sm text-gray-500">
          No projects yet. Create one below.
        </p>
      )}

      {!loading && projects.length > 0 && (
        <ul className="space-y-2">
          {projects.map((project) => (
            <li
              key={project.id}
              className="border rounded px-3 py-2 text-sm"
            >
              {project.name}
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}
