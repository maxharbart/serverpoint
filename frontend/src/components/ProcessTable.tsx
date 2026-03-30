'use client';

import { useState } from 'react';
import { Process } from '@/lib/api';

export default function ProcessTable({ processes }: { processes: Process[] }) {
  const [sortBy, setSortBy] = useState<'cpu' | 'ram'>('cpu');

  const sorted = [...processes].sort((a, b) =>
    sortBy === 'cpu' ? b.cpu - a.cpu : b.ram - a.ram
  );

  return (
    <div className="bg-gray-900 border border-gray-800 rounded-xl p-5">
      <div className="flex items-center justify-between mb-3">
        <h3 className="text-sm font-semibold text-gray-400">Top Processes</h3>
        <div className="flex gap-1">
          <button
            onClick={() => setSortBy('cpu')}
            className={`px-3 py-1 text-xs rounded-lg transition ${
              sortBy === 'cpu' ? 'bg-blue-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            CPU
          </button>
          <button
            onClick={() => setSortBy('ram')}
            className={`px-3 py-1 text-xs rounded-lg transition ${
              sortBy === 'ram' ? 'bg-purple-600 text-white' : 'bg-gray-800 text-gray-400 hover:text-white'
            }`}
          >
            RAM
          </button>
        </div>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full text-sm">
          <thead>
            <tr className="text-gray-500 text-xs border-b border-gray-800">
              <th className="text-left py-2 font-medium">PID</th>
              <th className="text-left py-2 font-medium">User</th>
              <th className="text-right py-2 font-medium">CPU %</th>
              <th className="text-right py-2 font-medium">RAM %</th>
              <th className="text-left py-2 font-medium pl-4">Command</th>
            </tr>
          </thead>
          <tbody>
            {sorted.map((p) => (
              <tr key={p.id} className="border-b border-gray-800/50 hover:bg-gray-800/30">
                <td className="py-1.5 font-mono text-gray-400">{p.pid}</td>
                <td className="py-1.5 text-gray-400">{p.user}</td>
                <td className="py-1.5 text-right font-mono">{p.cpu.toFixed(1)}</td>
                <td className="py-1.5 text-right font-mono">{p.ram.toFixed(1)}</td>
                <td className="py-1.5 pl-4 font-mono text-gray-400 truncate max-w-xs">{p.command}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
