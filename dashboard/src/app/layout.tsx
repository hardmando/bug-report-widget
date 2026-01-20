import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
    title: "Bug Report Dashboard",
    description: "Manage bug reports and API keys",
};

export default function RootLayout({
    children,
}: Readonly<{
    children: React.ReactNode;
}>) {
    return (
        <html lang="en">
            <body>{children}</body>
        </html>
    );
}
