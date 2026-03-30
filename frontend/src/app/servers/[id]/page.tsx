'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { isAuthenticated } from '@/lib/auth';
import Navbar from '@/components/Navbar';
import CpuChart from '@/components/CpuChart';
import RamChart from '@/components/RamChart';
import DiskUsage from '@/components/DiskUsage';
import ProcessTable from '@/components/ProcessTable';

export default function ServerDetailPage() {
  const params = useParams();
  const router = useRouter();
  const id = Number(params.id);
  const [duration, setDuration] = useState('1h');

  useEffect(() => {
    if (!isAuthenticated()) router.replace('/login');
  }, [router]);

  const { data: server } = useQuery({
    queryKey: ['server', id],
    queryFn: () => api.getServer(id),
    refetchInterval: 30_000,
  });

  const { data: metrics } = useQuery({
    queryKey: ['metrics', id, duration],
    queryFn: () => api.getMetrics(id, duration),
    refetchInterval: 30_000,
  });

  const { data: latest } = useQuery({
    queryKey: ['latest-metric', id],
    queryFn: () => api.getLatestMetric(id),
    refetchInterval: 30_000,
  });

  const { data: processes } = useQuery({
    queryKey: ['processes', id],
    queryFn: () => api.getProcesses(id),
    refetchInterval: 30_000,
  });

  async function handleDelete() {
    if (!confirm('Are you sure you want to delete this server?')) return;
    await api.deleteServer(id);
    router.push('/dashboard');
  }

  return (
    <div className="min-h-screen">
      <Navbar />
      <div className="max-w-6xl mx-auto px-6 py-6">
        <div className="flex items-center justify-between mb-6">
          <div>
            <button
              onClick={() => router.push('/dashboard')}
              className="text-sm text-gray-500 hover:text-gray-300 mb-1"
            >
              &larr; Back to Dashboard
            </button>
            <div className="flex items-center gap-3">
              <h1 className="text-2xl font-bold">{server?.name || 'Loading...'}</h1>
              {server && (
                <span
                  className={`px-2 py-0.5 text-xs rounded-full ${
                    server.online
                      ? 'bg-green-500/20 text-green-400'
                      : 'bg-red-500/20 text-red-400'
                  }`}
                >
                  {server.online ? 'Online' : 'Offline'}
                </span>
              )}
            </div>
          </div>
          <button
            onClick={handleDelete}
            className="px-4 py-2 bg-red-600/20 text-red-400 hover:bg-red-600/30 rounded-lg text-sm transition"
          >
            Delete Server
          </button>
        </div>

        {server && (
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <p className="text-xs text-gray-500">Host</p>
              <p className="font-mono text-sm">{server.host}:{server.port}</p>
            </div>
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <p className="text-xs text-gray-500">OS</p>
              <p className="text-sm truncate">{server.os || 'N/A'}</p>
            </div>
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <p className="text-xs text-gray-500">Hostname</p>
              <p className="text-sm truncate">{server.hostname || 'N/A'}</p>
            </div>
            <div className="bg-gray-900 border border-gray-800 rounded-xl p-4">
              <p className="text-xs text-gray-500">Uptime</p>
              <p className="text-sm truncate">{server.uptime || 'N/A'}</p>
            </div>
          </div>
        )}

        <div className="flex gap-2 mb-4">
          {['1h', '6h', '24h'].map((d) => (
            <button
              key={d}
              onClick={() => setDuration(d)}
              className={`px-3 py-1 text-sm rounded-lg transition ${
                duration === d
                  ? 'bg-blue-600 text-white'
                  : 'bg-gray-800 text-gray-400 hover:text-white'
              }`}
            >
              {d}
            </button>
          ))}
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-4">
          <CpuChart metrics={metrics || []} />
          <RamChart metrics={metrics || []} />
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <DiskUsage metric={latest || null} />
          <ProcessTable processes={processes || []} />
        </div>
      </div>
    </div>
  );
}
