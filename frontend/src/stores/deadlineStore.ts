import { create } from 'zustand'
import { listDeadlines, listDeadlinesByCase } from '@/api/deadline'
import type { DeadlineItem } from '@/types'

interface DeadlineState {
  list: DeadlineItem[]
  total: number
  byCase: DeadlineItem[]
  fetchList: (params?: Record<string, unknown>) => Promise<void>
  fetchByCase: (caseId: number) => Promise<void>
}

export const useDeadlineStore = create<DeadlineState>((set) => ({
  list: [],
  total: 0,
  byCase: [],
  async fetchList(params = {}) {
    const res: any = await listDeadlines(params)
    set({ list: res.data.list, total: res.data.total })
  },
  async fetchByCase(caseId: number) {
    const res: any = await listDeadlinesByCase(caseId)
    set({ byCase: res.data })
  },
}))
