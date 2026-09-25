import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { Card, Select, Space, Table, Tabs, Tag } from 'antd'
import { useDeadlineStore } from '@/stores/deadlineStore'
import { useUserStore } from '@/stores/userStore'
import { DeadlineView, DeadlineViewText, DeadlineTypeText } from '@/constants/deadline'
import { userNameOf } from '@/components/common/DeadlineList'
import { formatDate, formatDateTime } from '@/utils/dateFormat'
import type { DeadlineItem } from '@/types'

function daysTag(view: string, days?: number) {
  if (days === undefined) return '-'
  if (view === DeadlineView.OVERDUE) {
    return <Tag color="red">逾期 {-days} 天</Tag>
  }
  if (days === 0) {
    return <Tag color="red">今天到期</Tag>
  }
  return <Tag color={days <= 3 ? 'orange' : 'blue'}>剩余 {days} 天</Tag>
}

export default function Deadlines() {
  const { list, total, fetchCenter } = useDeadlineStore()
  const { staff, fetchStaff } = useUserStore()
  const [view, setView] = useState<string>(DeadlineView.UPCOMING)
  const [ownerId, setOwnerId] = useState<number>()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)

  useEffect(() => {
    fetchStaff()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  useEffect(() => {
    fetchCenter({ view, owner_id: ownerId, page, page_size: pageSize })
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [view, ownerId, page, pageSize])

  const baseColumns = [
    {
      title: '案号',
      dataIndex: 'case_no',
      render: (v: string, row: DeadlineItem) => <Link to={`/cases/${row.case_id}`}>{v}</Link>,
    },
    { title: '案件标题', dataIndex: 'case_title' },
    { title: '期限类型', dataIndex: 'deadline_type', render: (v: string) => DeadlineTypeText[v] || v },
    { title: '期限名称', dataIndex: 'name' },
    { title: '截止时间', dataIndex: 'due_date', render: (v: string) => formatDate(v) },
    { title: '责任人', dataIndex: 'owner_id', render: (v: number) => userNameOf(staff, v) },
  ]
  const columns =
    view === DeadlineView.COMPLETED
      ? [
          ...baseColumns,
          { title: '处理人', dataIndex: 'completed_by', render: (v: number) => userNameOf(staff, v) },
          { title: '完成时间', dataIndex: 'completed_at', render: (v: string) => (v ? formatDateTime(v) : '-') },
        ]
      : [...baseColumns, { title: '剩余/逾期', dataIndex: 'days', render: (v: number) => daysTag(view, v) }]

  return (
    <Card>
      <Tabs
        activeKey={view}
        onChange={(k) => {
          setView(k)
          setPage(1)
        }}
        items={Object.entries(DeadlineViewText).map(([key, label]) => ({ key, label }))}
      />
      <Space style={{ marginBottom: 16 }}>
        <Select
          placeholder="按责任人筛选"
          allowClear
          style={{ width: 160 }}
          value={ownerId}
          onChange={(v) => {
            setOwnerId(v)
            setPage(1)
          }}
          options={staff.map((u) => ({ label: u.real_name || u.username, value: u.id }))}
        />
      </Space>
      <Table<DeadlineItem>
        rowKey="id"
        dataSource={list}
        columns={columns}
        pagination={{ current: page, pageSize, total, onChange: (p, ps) => { setPage(p); setPageSize(ps) } }}
      />
    </Card>
  )
}
