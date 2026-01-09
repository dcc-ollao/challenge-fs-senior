import { useEffect, useMemo, useState } from "react";
import {
  createTask,
  listProjects,
  listTasksByProject,
  updateTask,
  listAssignableUsers,
} from "./api";
import type { Project, Task } from "./types";
import { useSnackbar } from "../../components/snackbar/SnackbarContext";
import { useAuthContext } from "../auth/AuthContext";

const STATUS_LABELS: Record<string, string> = {
  todo: "To do",
  in_progress: "In progress",
  done: "Done",
};

type AssignableUser = {
  id: string;
  email: string;
};

function StatusSelect({
  value,
  onChange,
  disabled,
}: {
  value?: string;
  onChange: (v: string) => void;
  disabled?: boolean;
}) {
  return (
    <select
      value={value ?? "todo"}
      onChange={(e) => onChange(e.target.value)}
      disabled={disabled}
      className="rounded-md border px-2 py-1 text-xs bg-white hover:bg-slate-50 disabled:opacity-60 disabled:cursor-not-allowed"
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

function getErrorMessage(err: any, fallback: string) {
  return (
    err?.response?.data?.message ??
    err?.response?.data?.error ??
    err?.message ??
    fallback
  );
}

export default function TasksPage() {
  const { showError } = useSnackbar();
  const { user } = useAuthContext();

  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<string>("");

  const [tasks, setTasks] = useState<Task[]>([]);
  const [users, setUsers] = useState<AssignableUser[]>([]);

  const [loadingProjects, setLoadingProjects] = useState(true);
  const [loadingTasks, setLoadingTasks] = useState(false);

  const [title, setTitle] = useState("");
  const [creating, setCreating] = useState(false);

  const selectedProject = useMemo(
    () => projects.find((p) => p.id === selectedProjectId) ?? null,
    [projects, selectedProjectId]
  );

  const projectMap = useMemo(
    () => Object.fromEntries(projects.map((p) => [p.id, p])),
    [projects]
  );

  // Load projects
  useEffect(() => {
    let mounted = true;

    async function loadProjects() {
      setLoadingProjects(true);
      try {
        const data = await listProjects();
        if (!mounted) return;
        setProjects(data);
        setSelectedProjectId(data.length > 0 ? data[0].id : "");
      } catch (err: any) {
        showError(getErrorMessage(err, "Failed to load projects."));
      } finally {
        if (mounted) setLoadingProjects(false);
      }
    }

    void loadProjects();
    return () => {
      mounted = false;
    };
  }, [showError]);

  // Load tasks
  useEffect(() => {
    if (!selectedProjectId) {
      setTasks([]);
      return;
    }

    let mounted = true;

    async function loadTasks() {
      setLoadingTasks(true);
      try {
        const data = await listTasksByProject(selectedProjectId);
        if (!mounted) return;
        setTasks(data ?? []);
      } catch (err: any) {
        showError(getErrorMessage(err, "Failed to load tasks."));
      } finally {
        if (mounted) setLoadingTasks(false);
      }
    }

    void loadTasks();
    return () => {
      mounted = false;
    };
  }, [selectedProjectId, showError]);

  useEffect(() => {
    let mounted = true;

    async function loadUsers() {
      try {
        const data = await listAssignableUsers();
        if (mounted) setUsers(data);
      } catch {
        showError("Failed to load users");
      }
    }

    void loadUsers();
    return () => {
      mounted = false;
    };
  }, [showError]);

  async function onCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!selectedProjectId) return;

    const trimmed = title.trim();
    if (!trimmed) return;

    setCreating(true);
    try {
      await createTask(selectedProjectId, trimmed);
      setTitle("");

      const data = await listTasksByProject(selectedProjectId);
      setTasks(data ?? []);
    } catch (err: any) {
      showError(getErrorMessage(err, "Failed to create task."));
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
        assigneeId: task.assigneeId ?? null,
      });
    } catch (err: any) {
      setTasks((prev) =>
        prev.map((t) => (t.id === task.id ? { ...t, status: prevStatus } : t))
      );
      showError(getErrorMessage(err, "Failed to update task status."));
    }
  }

  async function onAssigneeChange(task: Task, nextAssignee: string | null) {
    const prev = task.assigneeId;

    setTasks((prevTasks) =>
      prevTasks.map((t) =>
        t.id === task.id ? { ...t, assigneeId: nextAssignee } : t
      )
    );

    try {
      await updateTask(task.id, {
        title: task.title,
        description: task.description ?? "",
        status: task.status,
        assigneeId: nextAssignee,
      });
    } catch (err: any) {
      setTasks((prevTasks) =>
        prevTasks.map((t) =>
          t.id === task.id ? { ...t, assigneeId: prev } : t
        )
      );
      showError(getErrorMessage(err, "Failed to update assignee."));
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
          <div className="text-sm font-medium">
            Tasks for {selectedProject.name}
          </div>

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
          {tasks.map((t) => {
            const canChangeStatus =
              user?.role === "admin" || t.assigneeId === user?.id;

            return (
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

                <div className="flex items-center gap-3">
                  <select
                    value={t.assigneeId ?? ""}
                    disabled={user?.role !== "admin"}
                    onChange={(e) =>
                      onAssigneeChange(t, e.target.value || null)
                    }
                    className="rounded-md border px-2 py-1 text-xs bg-white disabled:opacity-60 disabled:cursor-not-allowed"
                  >
                    <option value="">Unassigned</option>
                    {users.map((u) => (
                      <option key={u.id} value={u.id}>
                        {u.email}
                      </option>
                    ))}
                  </select>

                  <StatusSelect
                    value={t.status}
                    onChange={(v) => onStatusChange(t, v)}
                    disabled={!canChangeStatus}
                  />
                </div>
              </li>
            );
          })}
        </ul>
        )
      ) : null}
    </div>
  );
}
