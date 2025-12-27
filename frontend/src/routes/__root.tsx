import { createRootRoute, Outlet } from '@tanstack/react-router'

export const Route = createRootRoute({
  component: RootComponent,
  notFoundComponent: () => <p>Page not found</p>,
})

function RootComponent() {
  return (
    <main className="container">
      <Outlet />
    </main>
  )
}
