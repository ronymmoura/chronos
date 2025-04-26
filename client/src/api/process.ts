import { api } from ".";
import { Process } from "../entities";

export const ProcessApi = {
  list: () => api.get<Process[]>("process").json(),

  get: (id: string) => api.get<Process>(`process/${id}`).json(),

  create: (data: Process, file: File) => {
    const formData = new FormData();
    formData.append("file", file);
    formData.append("name", data.name);
    formData.append("description", data.description);
    formData.append("env", data.env);
    formData.append("execute_every_secs", data.execute_every_secs.toString());

    return api.post<Process>("process", { body: formData }).json();
  },
};
