import {ref} from 'vue'

export type ThemeColor = string

// Helper to convert hex to HSL
function hexToHSL(hex: string) {
  let r = 0, g = 0, b = 0
  if (hex.length === 4) {
    r = parseInt("0x" + hex[1] + hex[1])
    g = parseInt("0x" + hex[2] + hex[2])
    b = parseInt("0x" + hex[3] + hex[3])
  } else if (hex.length === 7) {
    r = parseInt("0x" + hex[1] + hex[2])
    g = parseInt("0x" + hex[3] + hex[4])
    b = parseInt("0x" + hex[5] + hex[6])
  }
  r /= 255
  g /= 255
  b /= 255
  let cmin = Math.min(r, g, b),
      cmax = Math.max(r, g, b),
      delta = cmax - cmin,
      h = 0,
      s = 0,
      l = 0

  if (delta === 0)
    h = 0
  else if (cmax === r)
    h = ((g - b) / delta) % 6
  else if (cmax === g)
    h = (b - r) / delta + 2
  else
    h = (r - g) / delta + 4

  h = Math.round(h * 60)
  if (h < 0) h += 360

  l = (cmax + cmin) / 2
  s = delta === 0 ? 0 : delta / (1 - Math.abs(2 * l - 1))
  s = +(s * 100).toFixed(1)
  l = +(l * 100).toFixed(1)

  return { h, s, l }
}

export const themeColors: Record<string, { label: string; color: string; primary: string; primaryForeground: string; ring: string; accent: string; accentForeground: string }> = {
  indigo: {
    label: '靛蓝',
    color: '#6366f1',
    primary: '243 75% 59%',
    primaryForeground: '0 0% 100%',
    ring: '243 75% 59%',
    accent: '243 75% 96%',
    accentForeground: '243 75% 20%',
  },
  // ... other existing colors can remain if needed, but we will support dynamic hex
}

const currentTheme = ref<string>('indigo')

export function useTheme() {
  const applyThemeColor = (color: string) => {
    let vars
    if (themeColors[color]) {
      vars = themeColors[color]
    } else {
      // Dynamic hex color
      const { h, s, l } = hexToHSL(color)
      vars = {
        label: 'Custom',
        color: color,
        primary: `${h} ${s}% ${l}%`,
        primaryForeground: '0 0% 100%', // Assuming white text for now
        ring: `${h} ${s}% ${l}%`,
        accent: `${h} ${s}% 96%`, // Very light version
        accentForeground: `${h} ${s}% 20%`, // Dark version
      }
    }

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
    const saved = localStorage.getItem('sb-theme-color')
    if (saved) {
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
