import dayjs from 'dayjs'
import { DeadlineStatus } from '@/constants/deadline'

export interface DueDayInfo {
  days: number
  label: string
  tone: 'ok' | 'warn' | 'danger' | 'done'
}

// 按自然日计算剩余（正数）/逾期（负数）天数，今天为 0
export function dueDays(dueAt: string | Date, now: dayjs.Dayjs = dayjs()): number {
  const due = dayjs(dueAt).startOf('day')
  return due.diff(now.startOf('day'), 'day')
}

// 期限剩余/逾期文案与色彩基调，已完成直接返回"已完成"
export function dueDayInfo(dueAt: string | Date, status: string, now: dayjs.Dayjs = dayjs()): DueDayInfo {
  if (status === DeadlineStatus.COMPLETED) {
    return { days: 0, label: '已完成', tone: 'done' }
  }
  const days = dueDays(dueAt, now)
  if (days === 0) return { days, label: '今天到期', tone: 'danger' }
  if (days > 0) return { days, label: `剩余 ${days} 天`, tone: days <= 3 ? 'warn' : 'ok' }
  return { days, label: `已逾期 ${-days} 天`, tone: 'danger' }
}

export const dueToneColor: Record<DueDayInfo['tone'], string> = {
  ok: 'green',
  warn: 'orange',
  danger: 'red',
  done: 'default',
}
