'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '@/lib/api';
import { isAuthenticated } from '@/lib/auth';

export default function Home() {
  const router = useRouter();

  useEffect(() => {
    async function check() {
      try {
        const { setup_required } = await api.checkSetup();
        if (setup_required) {
          router.replace('/setup');
        } else if (isAuthenticated()) {
          router.replace('/dashboard');
        } else {
          router.replace('/login');
        }
      } catch {
        router.replace('/login');
      }
    }
    check();
  }, [router]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-gray-400">Loading...</div>
    </div>
  );
}
