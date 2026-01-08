import { useEffect, useState } from "react";
import { listProjects, createProject } from "./api";
import type { Project } from "./types";
import { useSnackbar } from "../../components/snackbar/SnackbarContext";

export default function ProjectsPage() {
  const { showError, showSuccess } = useSnackbar();

  const [projects, setProjects] = useState<Project[]>([]);
  const [loading, setLoading] = useState(true);

  const [name, setName] = useState("");
  const [creating, setCreating] = useState(false);

  async function loadProjects() {
    setLoading(true);
    try {
      const data = await listProjects();
      setProjects(data);
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ??
        err?.response?.data?.error ??
        err?.message ??
        "Failed to load projects.";

      showError(msg);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreate(e: React.FormEvent) {
    e.preventDefault();
    if (!name.trim()) return;

    setCreating(true);
    try {
      await createProject(name.trim());
      setName("");
      showSuccess("Project created");
      await loadProjects();
    } catch (err: any) {
      const msg =
        err?.response?.data?.message ??
        err?.response?.data?.error ??
        err?.message ??
        "Failed to create project.";

      showError(msg);
    } finally {
      setCreating(false);
    }
  }

  useEffect(() => {
    void loadProjects();
  }, []);

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="space-y-1">
        <h1 className="text-2xl font-semibold tracking-tight">Projects</h1>
        <p className="text-slate-600">
          Organize your work and collaborate by project.
        </p>
      </div>

      {/* Create project */}
      <div className="rounded-lg border bg-white p-4 space-y-3">
        <div className="text-sm font-medium">Create project</div>

        <form onSubmit={handleCreate} className="flex gap-2">
          <input
            className="flex-1 rounded-md border px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-slate-900"
            placeholder="Project name"
            value={name}
            onChange={(e) => setName(e.target.value)}
            disabled={creating}
          />
          <button
            type="submit"
            className="rounded-md bg-slate-900 px-4 py-2 text-sm font-medium text-white hover:bg-slate-800 disabled:opacity-50"
            disabled={creating || !name.trim()}
          >
            Create
          </button>
        </form>
      </div>

      {/* Projects list */}
      <div className="space-y-3">
        <div className="text-sm font-medium">Your projects</div>

        {loading && (
          <div className="text-sm text-slate-500">Loading projects…</div>
        )}

        {!loading && projects.length === 0 && (
          <div className="rounded-lg border border-dashed bg-white p-8 text-center text-sm text-slate-600">
            No projects yet.<br />
            Create one to start organizing your work.
          </div>
        )}

        {!loading && projects.length > 0 && (
          <div className="space-y-2">
            {projects.map((project) => (
              <div
                key={project.id}
                className="flex items-center justify-between rounded-lg border bg-white px-4 py-3"
              >
                <div className="font-medium text-sm">{project.name}</div>
                <div className="text-xs text-slate-500">Project</div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
