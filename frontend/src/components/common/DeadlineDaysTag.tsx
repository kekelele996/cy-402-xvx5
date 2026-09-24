import { Tag } from 'antd'
import { DeadlineStatus } from '@/constants/deadline'
import { dueDayInfo, dueToneColor } from '@/utils/deadlineDays'

// 期限剩余/逾期天数标签，已完成显示"已完成"
export default function DeadlineDaysTag({ dueAt, status }: { dueAt: string | Date; status: string }) {
  const info = dueDayInfo(dueAt, status)
  if (status === DeadlineStatus.COMPLETED) {
    return <Tag>{info.label}</Tag>
  }
  return <Tag color={dueToneColor[info.tone]}>{info.label}</Tag>
}
