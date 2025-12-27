/// <reference types="vite/client" />
import type { ReactNode } from 'react'
import {
  Outlet,
  createRootRoute,
  HeadContent,
  Scripts,
} from '@tanstack/react-router'

export const Route = createRootRoute({
  head: () => ({
    meta: [
      { charSet: 'utf-8' },
      { name: 'viewport', content: 'width=device-width, initial-scale=1' },
      { title: 'Lofam - Task Manager' },
    ],
    links: [
      {
        rel: 'stylesheet',
        href: 'https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css',
      },
    ],
  }),
  component: RootComponent,
  notFoundComponent: () => <p>Page not found</p>,
})

function RootComponent() {
  return (
    <RootDocument>
      <Outlet />
    </RootDocument>
  )
}

function RootDocument({ children }: Readonly<{ children: ReactNode }>) {
  return (
    <html lang="en">
      <head>
        <HeadContent />
      </head>
      <body>
        <main className="container">
          {children}
        </main>
        <Scripts />
      </body>
    </html>
  )
}
