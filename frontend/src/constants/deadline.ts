// 期限类型/状态/视图枚举（与后端 backend/internal/constants/deadline.go 保持一致）
export const DeadlineType = {
  HEARING: 'hearing',
  APPEAL: 'appeal',
  EVIDENCE: 'evidence',
  TRIAL: 'trial',
  EXECUTION: 'execution',
  LIMITATION: 'limitation',
  REPLY: 'reply',
  REGISTRATION: 'registration',
  OTHER: 'other',
} as const

export const DeadlineTypeText: Record<string, string> = {
  [DeadlineType.HEARING]: '开庭',
  [DeadlineType.APPEAL]: '上诉',
  [DeadlineType.EVIDENCE]: '举证',
  [DeadlineType.TRIAL]: '审理',
  [DeadlineType.EXECUTION]: '执行',
  [DeadlineType.LIMITATION]: '诉讼时效',
  [DeadlineType.REPLY]: '答辩',
  [DeadlineType.REGISTRATION]: '立案登记',
  [DeadlineType.OTHER]: '其他',
}

export const DeadlineStatus = {
  PENDING: 'pending',
  COMPLETED: 'completed',
} as const

export const DeadlineStatusText: Record<string, string> = {
  [DeadlineStatus.PENDING]: '未完成',
  [DeadlineStatus.COMPLETED]: '已完成',
}

// 期限中心视图
export const DeadlineView = {
  UPCOMING: 'upcoming',
  OVERDUE: 'overdue',
  COMPLETED: 'completed',
} as const

export const DeadlineViewText: Record<string, string> = {
  [DeadlineView.UPCOMING]: '即将到期',
  [DeadlineView.OVERDUE]: '已逾期',
  [DeadlineView.COMPLETED]: '已完成',
}

export const DeadlineTypeOptions = Object.entries(DeadlineTypeText).map(([value, label]) => ({ label, value }))
export const DeadlineViewOptions = Object.entries(DeadlineViewText).map(([value, label]) => ({ label, value }))
