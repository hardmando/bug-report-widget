'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

interface Bug {
    id: string;
    description: string;
    status: string;
    created_at: string;
}

export default function DashboardPage() {
    const [bugs, setBugs] = useState<Bug[]>([]); // Initialize as empty array
    const [loading, setLoading] = useState(true);
    const router = useRouter();

    useEffect(() => {
        // Check auth
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }

        // Fetch bugs (Mock for now until bug-service is updated)
        // In real implementation: invoke fetch to bug-service
        setLoading(false);
    }, [router]);

    return (
        <div style={{ padding: '2rem' }}>
            <header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '2rem' }}>
                <h1>Dashboard</h1>
                <div style={{ display: 'flex', gap: '1rem' }}>
                    <Link href="/api-keys" style={{ color: '#0070f3', fontWeight: 'bold' }}>Manage API Keys</Link>
                    <button onClick={() => { localStorage.removeItem('token'); router.push('/login'); }} style={{ background: 'none', border: 'none', color: 'red', cursor: 'pointer' }}>Logout</button>
                </div>
            </header>

            <div style={{ background: 'white', padding: '1.5rem', borderRadius: '8px', boxShadow: '0 2px 4px rgba(0,0,0,0.1)' }}>
                <h2>Recent Bug Reports</h2>
                {bugs.length === 0 ? (
                    <p style={{ color: '#666' }}>No bug reports found (or functionality not yet implemented).</p>
                ) : (
                    <ul>
                        {bugs.map((bug) => (
                            <li key={bug.id}>{bug.description}</li>
                        ))}
                    </ul>
                )}
            </div>
        </div>
    );
}
