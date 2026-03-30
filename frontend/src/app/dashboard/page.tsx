'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import { isAuthenticated } from '@/lib/auth';
import Navbar from '@/components/Navbar';
import ServerTile from '@/components/ServerTile';
import AddServerModal from '@/components/AddServerModal';

export default function DashboardPage() {
  const router = useRouter();
  const [showAdd, setShowAdd] = useState(false);

  useEffect(() => {
    if (!isAuthenticated()) {
      router.replace('/login');
    }
  }, [router]);

  const { data: servers, refetch } = useQuery({
    queryKey: ['servers'],
    queryFn: api.getServers,
    refetchInterval: 30_000,
  });

  return (
    <div className="min-h-screen">
      <Navbar />
      <div className="max-w-7xl mx-auto px-6 py-6">
        <div className="flex items-center justify-between mb-6">
          <h1 className="text-2xl font-bold">Servers</h1>
          <button
            onClick={() => setShowAdd(true)}
            className="px-4 py-2 bg-blue-600 hover:bg-blue-700 rounded-lg text-sm font-medium transition"
          >
            + Add Server
          </button>
        </div>

        {servers && servers.length === 0 && (
          <div className="text-center py-20 text-gray-500">
            <p className="text-lg mb-2">No servers yet</p>
            <p className="text-sm">Add a server to start monitoring</p>
          </div>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {servers?.map((server) => (
            <ServerTile key={server.id} server={server} />
          ))}
        </div>
      </div>

      <AddServerModal
        open={showAdd}
        onClose={() => setShowAdd(false)}
        onCreated={() => refetch()}
      />
    </div>
  );
}
