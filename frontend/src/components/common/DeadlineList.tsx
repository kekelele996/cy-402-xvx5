import { Button, Popconfirm, Space, Table, Tag } from 'antd'
import type { ColumnsType } from 'antd/es/table'
import { Link } from 'react-router-dom'
import DeadlineDaysTag from './DeadlineDaysTag'
import { DeadlineStatusText, DeadlineTypeText } from '@/constants/deadline'
import { formatDateTime } from '@/utils/dateFormat'
import type { DeadlineItem } from '@/types'

interface Props {
  loading?: boolean
  data: DeadlineItem[]
  // 期限中心显示案号/标题列与分页；案件详情不显示
  showCase?: boolean
  pagination?: false | { current: number; pageSize: number; total: number; onChange: (page: number, pageSize: number) => void }
  onComplete?: (item: DeadlineItem) => void
  onEdit?: (item: DeadlineItem) => void
}

// 期限列表：案件详情与期限中心共用
export default function DeadlineList({ loading, data, showCase, pagination, onComplete, onEdit }: Props) {
  const columns: ColumnsType<DeadlineItem> = []
  if (showCase) {
    columns.push({
      title: '案号 / 标题',
      render: (_, r) => (
        <Space direction="vertical" size={0}>
          <Link to={`/cases/${r.case_id}`}>{r.case_no || `#${r.case_id}`}</Link>
          <span style={{ color: 'rgba(0,0,0,0.65)' }}>{r.case_title}</span>
        </Space>
      ),
    })
  }
  columns.push(
    { title: '类型', dataIndex: 'type', width: 100, render: (v: string) => DeadlineTypeText[v] || v },
    {
      title: '名称',
      dataIndex: 'name',
      render: (v: string) => <strong>{v}</strong>,
    },
    {
      title: '截止时间',
      dataIndex: 'due_at',
      width: 160,
      render: (v: string) => formatDateTime(v),
    },
    {
      title: '剩余 / 逾期',
      dataIndex: 'due_at',
      width: 110,
      render: (v: string, r) => <DeadlineDaysTag dueAt={v} status={r.status} />,
    },
    { title: '责任人', dataIndex: 'assignee_name', width: 100, render: (v: string, r) => v || `#${r.assignee_id}` },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      render: (v: string) => <Tag color={v === 'completed' ? 'green' : 'orange'}>{DeadlineStatusText[v] || v}</Tag>,
    },
    {
      title: '完成人 / 时间',
      width: 180,
      render: (_, r) =>
        r.status === 'completed' ? (
          <Space direction="vertical" size={0}>
            <span>{r.completer_name || (r.completed_by_id ? `#${r.completed_by_id}` : '-')}</span>
            <span style={{ color: 'rgba(0,0,0,0.45)' }}>{formatDateTime(r.completed_at)}</span>
          </Space>
        ) : (
          '-'
        ),
    },
  )
  if (onComplete || onEdit) {
    columns.push({
      title: '操作',
      width: 140,
      render: (_, r) =>
        r.status === 'pending' ? (
          <Space>
            {onEdit && <Button size="small" onClick={() => onEdit(r)}>修改</Button>}
            {onComplete && (
              <Popconfirm title="确认该期限已处理完成？" onConfirm={() => onComplete(r)}>
                <Button size="small" type="primary">标记完成</Button>
              </Popconfirm>
            )}
          </Space>
        ) : null,
    })
  }

  return (
    <Table<DeadlineItem>
      rowKey="id"
      loading={loading}
      dataSource={data}
      columns={columns}
      pagination={pagination}
      size="middle"
    />
  )
}
