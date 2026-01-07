import api from "../../lib/api";
import type { Project } from "./types";

export async function listProjects(): Promise<Project[]> {
  const res = await api.get<Project[]>("/api/projects");
  return res.data;
}

export async function createProject(name: string): Promise<Project> {
  const res = await api.post<Project>("/api/projects", { name });
  return res.data;
}
