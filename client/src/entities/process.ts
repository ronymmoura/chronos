export type Process = {
  id: string;
  name: string;
  description: string;
  path: string;
  env: string;
  execute_every_secs: number;
  status: string;
  running: boolean;
  created_at: Date;
};
