'use client';

import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { removeToken } from '@/lib/auth';

export default function Navbar() {
  const router = useRouter();

  function handleLogout() {
    removeToken();
    router.push('/login');
  }

  return (
    <nav className="bg-gray-900 border-b border-gray-800 px-6 py-3 flex items-center justify-between">
      <Link href="/dashboard" className="text-xl font-bold text-blue-400 hover:text-blue-300">
        ServerPoint
      </Link>
      <button
        onClick={handleLogout}
        className="text-sm text-gray-400 hover:text-gray-200 transition"
      >
        Logout
      </button>
    </nav>
  );
}
