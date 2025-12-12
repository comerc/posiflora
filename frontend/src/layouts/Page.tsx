import { Outlet } from 'react-router-dom'

function PageLayout() {
  return (
    <div className="flex min-h-full flex-col">
      <div className="flex min-h-0 flex-1 flex-col items-center overflow-y-auto">
        <Outlet />
      </div>
      <Footer className="flex-shrink-0" />
    </div>
  )
}

function Footer({ className }: { className?: string }) {
  return (
    <div className={`flex justify-between px-3 py-1 ${className || ''}`}>
      <div>Posiflora</div>
      <div>v0.1</div>
    </div>
  )
}

export default PageLayout
