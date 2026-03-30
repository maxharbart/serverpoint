'use client';

import { Metric } from '@/lib/api';

function formatBytes(bytes: number): string {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

export default function DiskUsage({ metric }: { metric: Metric | null }) {
  if (!metric) return null;

  const percent = metric.disk_percent;
  const color = percent > 90 ? 'bg-red-500' : percent > 70 ? 'bg-yellow-500' : 'bg-green-500';

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
      <h3 className="text-sm font-semibold text-gray-400 mb-3">Disk Usage</h3>
      <div className="flex items-end gap-4">
        <div className="flex-1">
          <div className="w-full bg-gray-800 rounded-full h-4 mb-2">
            <div
              className={`h-4 rounded-full ${color} transition-all`}
              style={{ width: `${Math.min(percent, 100)}%` }}
            />
          </div>
          <div className="flex justify-between text-xs text-gray-500">
            <span>{formatBytes(metric.disk_used)} used</span>
            <span>{formatBytes(metric.disk_total)} total</span>
          </div>
        </div>
        <div className="text-2xl font-bold font-mono">{percent.toFixed(1)}%</div>
      </div>
    </div>
  );
}
