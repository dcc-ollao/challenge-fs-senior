export type Project = {
  id: string;
  name: string;
};

export type Task = {
  id: string;
  projectId?: string;
  title: string;
  description?: string | null;
  status?: string;
  assigneeId?: string | null;
};

export type CreateTaskInput = {
  title: string;
  description?: string;
  status?: string;
  assignee_id?: string;
};
