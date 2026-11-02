import { BrowserRouter, Route, Routes } from 'react-router-dom'
import { ThemeToggleButton } from './components/ThemeToggleButton'
import { BoilerDetailPage } from './pages/BoilerDetailPage'
import { BoilerGridPage } from './pages/BoilerGridPage'
import { NotFoundPage } from './pages/NotFoundPage'
import { WarehouseSelectPage } from './pages/WarehouseSelectPage'
import { ThemeProvider } from './theme/ThemeContext'

export default function App() {
  return (
    <ThemeProvider>
      <BrowserRouter>
        <ThemeToggleButton />
        <main className="app-shell">
          <Routes>
            <Route path="/" element={<WarehouseSelectPage />} />
            <Route path="/warehouse/:warehouseId" element={<BoilerGridPage />} />
            <Route path="/warehouse/:warehouseId/boiler/:boilerId" element={<BoilerDetailPage />} />
            <Route path="*" element={<NotFoundPage />} />
          </Routes>
        </main>
      </BrowserRouter>
    </ThemeProvider>
  )
}
