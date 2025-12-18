import { createBrowserRouter } from 'react-router-dom'

import AppLayout from '@/layouts/App'
import PageLayout from '@/layouts/Page'
import ShopPage from '@/pages/Shop'
import NotFoundPage from '@/pages/NotFound'

export default createBrowserRouter([
  {
    Component: AppLayout,
    children: [
      {
        Component: PageLayout,
        children: [
          {
            path: 'shops/:shopId/growth/telegram',
            Component: ShopPage,
          },
        ],
      },
      {
        path: '*',
        Component: NotFoundPage,
      },
    ],
  },
])
