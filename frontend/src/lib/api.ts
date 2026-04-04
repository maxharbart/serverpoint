const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081/api';

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = typeof window !== 'undefined' ? localStorage.getItem('token') : null;

  const res = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...options.headers,
    },
  });

  if (res.status === 401) {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
      window.location.href = '/login';
    }
    throw new Error('Unauthorized');
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `Request failed: ${res.status}`);
  }

  return res.json();
}

export const api = {
  checkSetup: () => request<{ setup_required: boolean }>('/auth/check'),

  register: (username: string, password: string) =>
    request<{ token: string }>('/auth/register', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),

  login: (username: string, password: string) =>
    request<{ token: string }>('/auth/login', {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    }),

  getServers: () => request<Server[]>('/servers'),

  getServer: (id: number) => request<Server>(`/servers/${id}`),

  createServer: (data: CreateServerRequest) =>
    request<Server>('/servers', {
      method: 'POST',
      body: JSON.stringify(data),
    }),

  deleteServer: (id: number) =>
    request<{ message: string }>(`/servers/${id}`, { method: 'DELETE' }),

  getMetrics: (id: number, duration = '1h') =>
    request<Metric[]>(`/servers/${id}/metrics?duration=${duration}`),

  getLatestMetric: (id: number) =>
    request<Metric | null>(`/servers/${id}/metrics/latest`),

  getProcesses: (id: number) => request<Process[]>(`/servers/${id}/processes`),
};

export interface Server {
  id: number;
  name: string;
  host: string;
  port: number;
  username: string;
  auth_type: string;
  os: string;
  hostname: string;
  uptime: string;
  online: boolean;
  created_at: string;
}

export interface Metric {
  id: number;
  server_id: number;
  cpu_percent: number;
  ram_total: number;
  ram_used: number;
  ram_percent: number;
  disk_total: number;
  disk_used: number;
  disk_percent: number;
  created_at: string;
}

export interface Process {
  id: number;
  server_id: number;
  pid: number;
  user: string;
  cpu: number;
  ram: number;
  command: string;
}

export interface CreateServerRequest {
  name: string;
  host: string;
  port: number;
  username: string;
  auth_type: 'password' | 'key';
  password?: string;
  private_key?: string;
}
