import { useEffect, useMemo, useState } from "react";
import { createTask, listProjects, listTasksByProject } from "./api";
import type { Project, Task } from "./types";

export default function TasksPage() {
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<string>("");

  const [tasks, setTasks] = useState<Task[]>([]);
  const [loadingProjects, setLoadingProjects] = useState(true);
  const [loadingTasks, setLoadingTasks] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const [title, setTitle] = useState("");
  const [creating, setCreating] = useState(false);

  const selectedProject = useMemo(
    () => projects.find((p) => p.id === selectedProjectId) ?? null,
    [projects, selectedProjectId]
  );

  useEffect(() => {
    let mounted = true;

    async function load() {
      setError(null);
      setLoadingProjects(true);
      try {
        console.log("selectedProjectId", selectedProjectId);
        const data = await listProjects();
        if (!mounted) return;

        setProjects(data);
        if (data.length > 0) setSelectedProjectId(data[0].id);
      } catch {
        setError("Failed to load projects.");
      } finally {
        if (mounted) setLoadingProjects(false);
      }
    }

    load();
    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    if (!selectedProjectId) return;

    let mounted = true;

    async function loadTasks() {
      setError(null);
      setLoadingTasks(true);
      try {
        const data = await listTasksByProject(selectedProjectId);
        if (!mounted) return;
        setTasks(data ?? []);
      } catch {
        setError("Failed to load tasks.");
      } finally {
        if (mounted) setLoadingTasks(false);
      }
    }

    loadTasks();
    return () => {
      mounted = false;
    };
  }, [selectedProjectId]);

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!selectedProjectId) return;

    const trimmed = title.trim();
    if (!trimmed) return;

    setCreating(true);
    setError(null);
    try {
      await createTask(selectedProjectId, trimmed);
      setTitle("");
      const data = await listTasksByProject(selectedProjectId);
      setTasks(data ?? []);
    } catch {
      setError("Failed to create task.");
    } finally {
      setCreating(false);
    }
  }

  return (
    <div className="space-y-6">
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold">Tasks</h1>
        <p className="text-slate-600">Project-scoped tasks (MVP).</p>
      </div>

      {error ? (
        <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          {error}
        </div>
      ) : null}

      <div className="flex items-center gap-3">
        <label className="text-sm text-slate-600">Project</label>
        <select
          className="rounded-md border px-3 py-2 text-sm"
          value={selectedProjectId}
          onChange={(e) => setSelectedProjectId(e.target.value)}
          disabled={loadingProjects || projects.length === 0}
        >
          {projects.map((p) => (
            <option key={p.id} value={p.id}>
              {p.name}
            </option>
          ))}
        </select>

        {loadingProjects ? <span className="text-sm text-slate-500">Loading…</span> : null}
      </div>

      {!loadingProjects && projects.length === 0 ? (
        <div className="rounded-md border p-4 text-sm text-slate-600">
          No projects found. Create a project first.
        </div>
      ) : null}

      {selectedProject ? (
        <form onSubmit={onCreate} className="flex items-center gap-2">
          <input
            className="w-full rounded-md border px-3 py-2 text-sm"
            placeholder="New task title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            disabled={creating}
          />
          <button
            type="submit"
            className="rounded-md border px-3 py-2 text-sm hover:bg-slate-50 disabled:opacity-50"
            disabled={creating || !title.trim()}
          >
            {creating ? "Adding…" : "Add"}
          </button>
        </form>
      ) : null}

      <div className="space-y-2">
        <div className="text-sm font-medium">
          {selectedProject ? `Tasks for: ${selectedProject.name}` : "Tasks"}
        </div>

        {loadingTasks ? (
          <div className="text-sm text-slate-500">Loading tasks…</div>
        ) : tasks.length === 0 && selectedProjectId ? (
          <div className="rounded-md border p-4 text-sm text-slate-600">No tasks yet.</div>
        ) : (
          <ul className="divide-y rounded-md border">
            {tasks.map((t) => (
              <li key={t.id} className="flex items-center justify-between px-4 py-3">
                <div className="text-sm">{t.title}</div>
                <div className="text-xs text-slate-500">{t.status ?? ""}</div>
              </li>
            ))}
          </ul>
        )}
      </div>
    </div>
  );
}
