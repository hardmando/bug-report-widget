'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function ApiKeysPage() {
    const [apiKey, setApiKey] = useState<string | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const router = useRouter();

    useEffect(() => {
        fetchApiKey();
    }, []);

    const fetchApiKey = async () => {
        const token = localStorage.getItem('token');
        if (!token) {
            router.push('/login');
            return;
        }

        try {
            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/api-keys`, {
                headers: { 'Authorization': `Bearer ${token}` }
            });

            if (res.status === 401) {
                localStorage.removeItem('token');
                router.push('/login');
                return;
            }

            const data = await res.json();
            setApiKey(data.api_key);
        } catch (err) {
            console.error(err);
            setError('Failed to fetch API key');
        } finally {
            setLoading(false);
        }
    };

    const generateApiKey = async () => {
        const token = localStorage.getItem('token');
        if (!token) return;

        try {
            const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081'}/api-keys`, {
                method: 'POST',
                headers: { 'Authorization': `Bearer ${token}` }
            });

            const data = await res.json();
            setApiKey(data.api_key);
        } catch (err) {
            setError('Failed to generate API Key');
        }
    };

    if (loading) return <div style={{ padding: '2rem' }}>Loading...</div>;

    return (
        <div style={{ padding: '2rem', maxWidth: '800px', margin: '0 auto' }}>
            <h1>API Keys</h1>
            <p style={{ marginBottom: '2rem' }}>Manage your API keys for the bug report widget.</p>

            {error && <p style={{ color: 'red' }}>{error}</p>}

            <div style={{
                padding: '1.5rem',
                border: '1px solid #ddd',
                borderRadius: '8px',
                background: 'white',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
            }}>
                <div>
                    <h3 style={{ margin: '0 0 0.5rem 0' }}>Current API Key</h3>
                    <code style={{
                        background: '#f0f0f0',
                        padding: '0.25rem 0.5rem',
                        borderRadius: '4px',
                        fontFamily: 'monospace',
                        display: 'block'
                    }}>
                        {apiKey || 'No API Key found'}
                    </code>
                </div>
                <button
                    onClick={generateApiKey}
                    style={{
                        padding: '0.5rem 1rem',
                        background: '#0070f3',
                        color: 'white',
                        border: 'none',
                        borderRadius: '4px',
                        cursor: 'pointer'
                    }}
                >
                    {apiKey ? 'Regenerate Key' : 'Generate Key'}
                </button>
            </div>

            <div style={{ marginTop: '2rem' }}>
                <a href="/dashboard" style={{ color: '#0070f3', textDecoration: 'underline' }}>Back to Dashboard</a>
            </div>
        </div>
    );
}
