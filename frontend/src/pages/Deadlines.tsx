import { useEffect, useMemo, useState } from 'react'
import { Button, Card, message, Radio, Select, Space } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import DeadlineList from '@/components/common/DeadlineList'
import DeadlineFormModal, { DeadlineFormValues } from '@/components/common/DeadlineFormModal'
import { completeDeadline, createDeadline, updateDeadline } from '@/api/deadline'
import { useDeadlineStore } from '@/stores/deadlineStore'
import { useUserStore } from '@/stores/userStore'
import { DeadlineView, DeadlineViewOptions } from '@/constants/deadline'
import type { DeadlineItem } from '@/types'

export default function Deadlines() {
  const store = useDeadlineStore()
  const userStore = useUserStore()
  const [view, setView] = useState<string>(DeadlineView.UPCOMING)
  const [assigneeId, setAssigneeId] = useState<number | undefined>()
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<DeadlineItem | null>(null)

  useEffect(() => {
    userStore.fetchStaff()
  }, [userStore])

  const params = useMemo(
    () => ({ page, page_size: pageSize, view, assignee_id: assigneeId }),
    [page, pageSize, view, assigneeId],
  )

  useEffect(() => {
    store.fetchList(params)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params])

  function refresh() {
    store.fetchList(params)
  }

  function switchView(v: string) {
    setView(v)
    setPage(1)
  }

  async function handleSubmit(values: DeadlineFormValues) {
    const payload = {
      type: values.type,
      name: values.name.trim(),
      due_at: values.due_at.format('YYYY-MM-DD HH:mm'),
      assignee_id: values.assignee_id,
    }
    if (editing) {
      await updateDeadline(editing.id, payload)
      message.success('期限已修改')
    } else {
      await createDeadline({ case_id: Number(values.case_id), ...payload })
      message.success('期限登记成功')
    }
    setOpen(false)
    setEditing(null)
    refresh()
  }

  async function handleComplete(item: DeadlineItem) {
    await completeDeadline(item.id)
    message.success('期限已标记完成')
    refresh()
  }

  return (
    <Card
      title="期限中心"
      extra={
        <Button type="primary" icon={<PlusOutlined />} onClick={() => { setEditing(null); setOpen(true) }}>
          登记期限
        </Button>
      }
    >
      <Space style={{ marginBottom: 16, justifyContent: 'space-between', width: '100%' }}>
        <Radio.Group value={view} options={DeadlineViewOptions} optionType="button" buttonStyle="solid" onChange={(e) => switchView(e.target.value)} />
        <Select
          allowClear
          showSearch
          placeholder="按责任人筛选"
          style={{ width: 180 }}
          value={assigneeId}
          optionFilterProp="label"
          onChange={(v) => { setAssigneeId(v); setPage(1) }}
          options={userStore.staff.map((u) => ({ label: u.real_name || u.username, value: u.id }))}
        />
      </Space>
      <DeadlineList
        showCase
        data={store.list}
        pagination={{
          current: page,
          pageSize,
          total: store.total,
          onChange: (p, ps) => { setPage(p); setPageSize(ps) },
        }}
        onEdit={(item) => { setEditing(item); setOpen(true) }}
        onComplete={handleComplete}
      />
      <DeadlineFormModal
        open={open}
        editing={editing}
        onCancel={() => { setOpen(false); setEditing(null) }}
        onSubmit={handleSubmit}
      />
    </Card>
  )
}
