import { Outlet } from 'react-router-dom'

function AppLayout() {
  return (
    <div className="h-screen w-screen">
      <Outlet />
    </div>
  )
}

export default AppLayout
