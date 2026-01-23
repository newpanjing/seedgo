export const formatDate = (val: string | number | Date | undefined | null, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!val) return '-'
  const date = new Date(val)
  if (isNaN(date.getTime())) return '-'

  const o: Record<string, number> = {
    'M+': date.getMonth() + 1, // month
    'D+': date.getDate(), // day
    'H+': date.getHours(), // hour
    'm+': date.getMinutes(), // minute
    's+': date.getSeconds(), // second
    'q+': Math.floor((date.getMonth() + 3) / 3), // quarter
    'S': date.getMilliseconds() // millisecond
  }

  if (/(Y+)/.test(format)) {
    format = format.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
  }

  for (const k in o) {
    if (new RegExp('(' + k + ')').test(format)) {
      format = format.replace(RegExp.$1, RegExp.$1.length === 1 ? String(o[k]) : ('00' + o[k]).substr(('' + o[k]).length))
    }
  }
  return format
}
