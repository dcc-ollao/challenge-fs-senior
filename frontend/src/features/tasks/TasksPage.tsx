import { useEffect, useMemo, useState } from "react";
import { createTask, listProjects, listTasksByProject, updateTask } from "./api";
import type { Project, Task } from "./types";

const STATUS_LABELS: Record<string, string> = {
  todo: "To do",
  in_progress: "In progress",
  done: "Done",
};

function StatusSelect({
  value,
  onChange,
}: {
  value?: string;
  onChange: (v: string) => void;
}) {
  return (
    <select
      value={value ?? "todo"}
      onChange={(e) => onChange(e.target.value)}
      className="rounded-md border px-2 py-1 text-xs bg-white hover:bg-slate-50"
    >
      {Object.entries(STATUS_LABELS).map(([key, label]) => (
        <option key={key} value={key}>
          {label}
        </option>
      ))}
    </select>
  );
}

function ProjectBadge({ name }: { name: string }) {
  return (
    <span className="rounded-full bg-slate-100 px-2 py-0.5 text-xs text-slate-700">
      {name}
    </span>
  );
}

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

  const projectMap = useMemo(() => {
    return Object.fromEntries(projects.map((p) => [p.id, p]));
  }, [projects]);

  useEffect(() => {
    let mounted = true;

    async function loadProjects() {
      setLoadingProjects(true);
      setError(null);
      try {
        const data = await listProjects();
        if (!mounted) return;
        setProjects(data);
        setSelectedProjectId(data.length > 0 ? data[0].id : "");
      } catch {
        setError("Failed to load projects.");
      } finally {
        if (mounted) setLoadingProjects(false);
      }
    }

    loadProjects();
    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    if (!selectedProjectId) {
      setTasks([]);
      return;
    }

    let mounted = true;

    async function loadTasks() {
      setLoadingTasks(true);
      setError(null);
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

  async function onStatusChange(task: Task, nextStatus: string) {
    const prevStatus = task.status;

    setTasks((prev) =>
      prev.map((t) => (t.id === task.id ? { ...t, status: nextStatus } : t))
    );

    try {
      await updateTask(task.id, {
        title: task.title,
        description: task.description ?? "",
        status: nextStatus,
        assignee_id: task.assigneeId ?? null,
      });
    } catch {
      setTasks((prev) =>
        prev.map((t) => (t.id === task.id ? { ...t, status: prevStatus } : t))
      );
      setError("Failed to update task status.");
    }
  }

  const showNoProjects = !loadingProjects && projects.length === 0;

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">Tasks</h1>
        <p className="text-slate-600">Manage tasks within each project.</p>
      </div>

      {error && (
        <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
          {error}
        </div>
      )}

      {/* Project selector */}
      <div className="rounded-lg border bg-white p-4 flex items-center gap-4">
        <label className="text-sm font-medium text-slate-600">Project</label>
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

        {loadingProjects && (
          <span className="text-sm text-slate-500">Loading…</span>
        )}
      </div>

      {/* No projects */}
      {showNoProjects && (
        <div className="rounded-lg border border-dashed bg-white p-8 text-center text-sm text-slate-600">
          No projects found. Create a project first to start adding tasks.
        </div>
      )}

      {/* Create task */}
      {selectedProject && (
        <div className="rounded-lg border bg-white p-4 space-y-3">
          <div className="text-sm font-medium">Tasks for {selectedProject.name}</div>

          <form onSubmit={onCreate} className="flex gap-2">
            <input
              className="flex-1 rounded-md border px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-900"
              placeholder="New task title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              disabled={creating}
            />
            <button
              type="submit"
              className="rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800 disabled:opacity-50"
              disabled={creating || !title.trim()}
            >
              {creating ? "Adding…" : "Add"}
            </button>
          </form>
        </div>
      )}

      {/* Tasks list */}
      {selectedProjectId && !showNoProjects ? (
        loadingTasks ? (
          <div className="text-sm text-slate-500">Loading tasks…</div>
        ) : tasks.length === 0 ? (
          <div className="rounded-lg border border-dashed bg-white p-8 text-center text-sm text-slate-600">
            No tasks yet.
          </div>
        ) : (
          <ul className="divide-y rounded-lg border bg-white">
            {tasks.map((t) => (
              <li
                key={t.id}
                className="flex items-center justify-between gap-4 px-4 py-3 hover:bg-slate-50 transition"
              >
                <div className="flex flex-col gap-1">
                  <div className="text-sm font-medium">{t.title}</div>
                  {t.projectId && projectMap[t.projectId] && (
                    <ProjectBadge name={projectMap[t.projectId].name} />
                  )}
                </div>

                <StatusSelect value={t.status} onChange={(v) => onStatusChange(t, v)} />
              </li>
            ))}
          </ul>
        )
      ) : null}
    </div>
  );
}
