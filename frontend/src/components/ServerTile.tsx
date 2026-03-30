'use client';

import Link from 'next/link';
import { useQuery } from '@tanstack/react-query';
import { api, Server } from '@/lib/api';
import { Area, AreaChart, ResponsiveContainer } from 'recharts';

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

export default function ServerTile({ server }: { server: Server }) {
  const { data: metrics } = useQuery({
    queryKey: ['metrics', server.id],
    queryFn: () => api.getMetrics(server.id, '1h'),
    refetchInterval: 30_000,
  });

  const { data: latest } = useQuery({
    queryKey: ['latest-metric', server.id],
    queryFn: () => api.getLatestMetric(server.id),
    refetchInterval: 30_000,
  });

  const cpuData = (metrics || []).map((m) => ({ v: m.cpu_percent }));
  const ramData = (metrics || []).map((m) => ({ v: m.ram_percent }));

  return (
    <Link href={`/servers/${server.id}`}>
      <div className="bg-gray-900 border border-gray-800 rounded-xl p-5 hover:border-gray-600 transition cursor-pointer">
        <div className="flex items-center justify-between mb-3">
          <h3 className="font-semibold text-lg truncate">{server.name}</h3>
          <span
            className={`w-2.5 h-2.5 rounded-full ${server.online ? 'bg-green-500' : 'bg-red-500'}`}
          />
        </div>

        <p className="text-xs text-gray-500 mb-3 truncate">{server.host}</p>

        {latest && server.online ? (
          <>
            <div className="grid grid-cols-2 gap-3 mb-3">
              <div>
                <p className="text-xs text-gray-500">CPU</p>
                <p className="text-sm font-mono">{latest.cpu_percent.toFixed(1)}%</p>
              </div>
              <div>
                <p className="text-xs text-gray-500">RAM</p>
                <p className="text-sm font-mono">{latest.ram_percent.toFixed(1)}%</p>
              </div>
              <div>
                <p className="text-xs text-gray-500">Disk</p>
                <p className="text-sm font-mono">{latest.disk_percent.toFixed(1)}%</p>
              </div>
              <div>
                <p className="text-xs text-gray-500">RAM Used</p>
                <p className="text-sm font-mono">
                  {formatBytes(latest.ram_used)}
                </p>
              </div>
            </div>

            <div className="grid grid-cols-2 gap-2 h-12">
              <div>
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={cpuData}>
                    <Area
                      type="monotone"
                      dataKey="v"
                      stroke="#3b82f6"
                      fill="#3b82f6"
                      fillOpacity={0.1}
                      strokeWidth={1.5}
                      dot={false}
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
              <div>
                <ResponsiveContainer width="100%" height="100%">
                  <AreaChart data={ramData}>
                    <Area
                      type="monotone"
                      dataKey="v"
                      stroke="#8b5cf6"
                      fill="#8b5cf6"
                      fillOpacity={0.1}
                      strokeWidth={1.5}
                      dot={false}
                    />
                  </AreaChart>
                </ResponsiveContainer>
              </div>
            </div>

            {server.uptime && (
              <p className="text-xs text-gray-600 mt-2 truncate">{server.uptime}</p>
            )}
          </>
        ) : (
          <p className="text-sm text-gray-600">{server.online ? 'Waiting for data...' : 'Offline'}</p>
        )}
      </div>
    </Link>
  );
}
