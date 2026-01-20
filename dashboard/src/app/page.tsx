'use client';
import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function Home() {
    const router = useRouter();

    useEffect(() => {
        const token = localStorage.getItem('token');
        if (token) {
            router.push('/dashboard');
        } else {
            router.push('/login');
        }
    }, [router]);

    return (
        <main style={{ padding: '2rem', textAlign: 'center' }}>
            <h1>Bug Report Widget</h1>
            <p>Redirecting...</p>
        </main>
    );
}
