import api from "../../lib/api";
import type { Project, Task } from "./types";

export async function listProjects(): Promise<Project[]> {
  const res = await api.get<Project[]>("/api/projects");
  return res.data;
}

export async function listTasksByProject(projectId: string): Promise<Task[]> {
  const res = await api.get<Task[]>(`/api/projects/${projectId}/tasks`);
  return res.data;
}

export async function createTask(projectId: string, title: string): Promise<Task> {
  const res = await api.post<Task>(`/api/projects/${projectId}/tasks`, { title });
  return res.data;
}
