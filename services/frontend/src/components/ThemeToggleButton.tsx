import { useTheme } from '../theme/useTheme'

function SunIcon() {
  return (
    <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" aria-hidden="true">
      <circle cx="12" cy="12" r="4.5" />
      <path d="M12 2v2.5M12 19.5V22M2 12h2.5M19.5 12H22M4.93 4.93l1.77 1.77M17.3 17.3l1.77 1.77M4.93 19.07l1.77-1.77M17.3 6.7l1.77-1.77" />
    </svg>
  )
}

function MoonIcon() {
  return (
    <svg viewBox="0 0 24 24" width="22" height="22" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
      <path d="M20.5 14.1A8.5 8.5 0 1 1 9.9 3.5a6.8 6.8 0 0 0 10.6 10.6Z" />
    </svg>
  )
}

export function ThemeToggleButton() {
  const { theme, toggleTheme } = useTheme()
  const isDark = theme === 'dark'
  const label = isDark ? 'Включить светлую тему' : 'Включить тёмную тему'

  return (
    <button type="button" className="theme-toggle" onClick={toggleTheme} aria-label={label} title={label}>
      {isDark ? <SunIcon /> : <MoonIcon />}
    </button>
  )
}
