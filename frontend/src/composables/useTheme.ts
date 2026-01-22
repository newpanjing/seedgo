import {ref} from 'vue'

export type ThemeColor = 'indigo' | 'green' | 'blue' | 'orange' | 'purple' | 'teal' | 'gray' | 'red' | 'yellow' | 'pink'

export const themeColors: Record<ThemeColor, { label: string; color: string; primary: string; primaryForeground: string; ring: string; accent: string; accentForeground: string }> = {
  indigo: {
    label: '靛蓝',
    color: '#6366f1',
    primary: '243 75% 59%',
    primaryForeground: '0 0% 100%',
    ring: '243 75% 59%',
    accent: '243 75% 96%',
    accentForeground: '243 75% 20%',
  },
  green: {
    label: '翡翠',
    color: '#10b981',
    primary: '150 99% 38%',
    primaryForeground: '0 0% 100%',
    ring: '150 99% 38%',
    accent: '150 60% 92%',
    accentForeground: '150 25% 15%',
  },
  blue: {
    label: '蔚蓝',
    color: '#3b82f6',
    primary: '217 91% 60%',
    primaryForeground: '0 0% 100%',
    ring: '217 91% 60%',
    accent: '217 80% 90%',
    accentForeground: '217 33% 20%',
  },
  orange: {
    label: '橙色',
    color: '#f97316',
    primary: '24 95% 56%',
    primaryForeground: '0 0% 100%',
    ring: '24 95% 56%',
    accent: '24 90% 92%',
    accentForeground: '22 30% 20%',
  },
  purple: {
    label: '紫色',
    color: '#a855f7',
    primary: '262 83% 58%',
    primaryForeground: '0 0% 100%',
    ring: '262 83% 58%',
    accent: '262 70% 90%',
    accentForeground: '262 40% 20%',
  },
  teal: {
    label: '青色',
    color: '#14b8a6',
    primary: '174 72% 45%',
    primaryForeground: '0 0% 100%',
    ring: '174 72% 45%',
    accent: '174 60% 90%',
    accentForeground: '174 30% 18%',
  },
  gray: {
    label: '灰色',
    color: '#6b7280',
    primary: '215 20% 40%',
    primaryForeground: '0 0% 100%',
    ring: '215 20% 40%',
    accent: '210 20% 92%',
    accentForeground: '215 25% 20%',
  },
  red: {
    label: '红色',
    color: '#ef4444',
    primary: '0 84% 60%',
    primaryForeground: '0 0% 100%',
    ring: '0 84% 60%',
    accent: '0 84% 96%',
    accentForeground: '0 84% 20%',
  },
  yellow: {
    label: '黄色',
    color: '#eab308',
    primary: '48 96% 53%',
    primaryForeground: '0 0% 100%',
    ring: '48 96% 53%',
    accent: '48 96% 96%',
    accentForeground: '48 96% 20%',
  },
  pink: {
    label: '粉色',
    color: '#ec4899',
    primary: '330 81% 60%',
    primaryForeground: '0 0% 100%',
    ring: '330 81% 60%',
    accent: '330 81% 96%',
    accentForeground: '330 81% 20%',
  }
}

const currentTheme = ref<ThemeColor>('indigo')

export function useTheme() {
  const applyThemeColor = (color: ThemeColor) => {
    const vars = themeColors[color]
    const root = document.documentElement
    root.style.setProperty('--primary', vars.primary)
    root.style.setProperty('--primary-foreground', vars.primaryForeground)
    root.style.setProperty('--ring', vars.ring)
    root.style.setProperty('--accent', vars.accent)
    root.style.setProperty('--accent-foreground', vars.accentForeground)
    currentTheme.value = color
    localStorage.setItem('sb-theme-color', color)
  }

  const initTheme = () => {
    const saved = localStorage.getItem('sb-theme-color') as ThemeColor | null
    if (saved && themeColors[saved]) {
      currentTheme.value = saved
      applyThemeColor(saved)
    } else {
      applyThemeColor('indigo')
    }
  }

  return {
    currentTheme,
    themeColors,
    applyThemeColor,
    initTheme
  }
}
