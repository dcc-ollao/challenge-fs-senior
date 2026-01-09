export type Project = {
  id: string;
  name: string;
};

export type Task = {
  id: string;
  projectId: string;
  title: string;
  description?: string;
  status: string;
  assigneeId: string | null;
  createdAt?: string;
  updatedAt?: string;
};


export type CreateTaskInput = {
  title: string;
  description?: string;
  status?: string;
  assignee_id?: string;
};
