import { useEffect, useState } from 'react'
import { Button, DatePicker, Form, Input, message, Modal, Select, Space, Table, Tag, Typography } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import { createDeadline, updateDeadline, completeDeadline } from '@/api/deadline'
import { useDeadlineStore } from '@/stores/deadlineStore'
import { useUserStore } from '@/stores/userStore'
import { DeadlineStatus, DeadlineStatusText, DeadlineTypeOptions, DeadlineTypeText } from '@/constants/deadline'
import { CaseStatus } from '@/constants/case'
import { formatDate, formatDateTime } from '@/utils/dateFormat'
import type { DeadlineItem, User } from '@/types'

export function userNameOf(users: User[], id?: number | null): string {
  if (!id) return '-'
  const u = users.find((x) => x.id === id)
  return u ? u.real_name || u.username : `#${id}`
}

export default function DeadlineList({ caseId, caseStatus }: { caseId: number; caseStatus: string }) {
  const { byCase, fetchByCase } = useDeadlineStore()
  const { staff, fetchStaff } = useUserStore()
  const [open, setOpen] = useState(false)
  const [editing, setEditing] = useState<DeadlineItem | null>(null)
  const [form] = Form.useForm()
  // 已结案或归档案件只能补录过去日期
  const closedCase = caseStatus === CaseStatus.CLOSED || caseStatus === CaseStatus.ARCHIVED

  useEffect(() => {
    fetchByCase(caseId)
    fetchStaff()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [caseId])

  // Modal 打开后再回填/重置表单，避免表单未挂载时操作告警
  useEffect(() => {
    if (!open) return
    if (editing) {
      form.setFieldsValue({
        deadline_type: editing.deadline_type,
        name: editing.name,
        due_date: dayjs(editing.due_date),
        owner_id: editing.owner_id,
      })
    } else {
      form.resetFields()
    }
  }, [open, editing, form])

  function openCreate() {
    setEditing(null)
    setOpen(true)
  }

  function openEdit(row: DeadlineItem) {
    setEditing(row)
    setOpen(true)
  }

  async function onSubmit() {
    const values = await form.validateFields()
    const payload = {
      deadline_type: values.deadline_type,
      name: values.name,
      due_date: values.due_date.format('YYYY-MM-DD'),
      owner_id: values.owner_id,
    }
    if (editing) {
      await updateDeadline(editing.id, payload)
      message.success('期限已更新')
    } else {
      await createDeadline({ case_id: caseId, ...payload })
      message.success('期限登记成功')
    }
    setOpen(false)
    fetchByCase(caseId)
  }

  async function onComplete(row: DeadlineItem) {
    await completeDeadline(row.id)
    message.success('期限已标记完成')
    fetchByCase(caseId)
  }

  return (
    <>
      <Space style={{ marginBottom: 16 }}>
        <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>登记期限</Button>
        {closedCase && <Typography.Text type="warning">案件已结案/归档，只能补录过去日期的期限</Typography.Text>}
      </Space>
      <Table<DeadlineItem>
        rowKey="id"
        dataSource={byCase}
        pagination={false}
        size="small"
        columns={[
          { title: '类型', dataIndex: 'deadline_type', render: (v) => DeadlineTypeText[v] || v },
          { title: '名称', dataIndex: 'name' },
          { title: '截止时间', dataIndex: 'due_date', render: (v) => formatDate(v) },
          { title: '责任人', dataIndex: 'owner_id', render: (v) => userNameOf(staff, v) },
          {
            title: '状态',
            dataIndex: 'status',
            render: (v) => (
              <Tag color={v === DeadlineStatus.COMPLETED ? 'green' : 'orange'}>{DeadlineStatusText[v] || v}</Tag>
            ),
          },
          { title: '处理人', dataIndex: 'completed_by', render: (v) => userNameOf(staff, v) },
          { title: '完成时间', dataIndex: 'completed_at', render: (v) => (v ? formatDateTime(v) : '-') },
          {
            title: '操作',
            render: (_, row) =>
              row.status === DeadlineStatus.PENDING ? (
                <Space>
                  <Button size="small" onClick={() => openEdit(row)}>修改</Button>
                  <Button size="small" type="primary" onClick={() => onComplete(row)}>标记完成</Button>
                </Space>
              ) : null,
          },
        ]}
      />
      <Modal
        title={editing ? '修改期限' : '登记期限'}
        open={open}
        onOk={onSubmit}
        onCancel={() => setOpen(false)}
        destroyOnClose
      >
        <Form form={form} layout="vertical">
          <Form.Item name="deadline_type" label="期限类型" rules={[{ required: true, message: '请选择期限类型' }]}>
            <Select options={DeadlineTypeOptions} placeholder="请选择期限类型" />
          </Form.Item>
          <Form.Item name="name" label="期限名称" rules={[{ required: true, message: '请输入期限名称' }, { max: 200 }]}>
            <Input placeholder="如：一审开庭、上诉截止" />
          </Form.Item>
          <Form.Item name="due_date" label="截止时间" rules={[{ required: true, message: '请选择截止时间' }]}>
            <DatePicker
              style={{ width: '100%' }}
              disabledDate={closedCase ? (d) => !d.isBefore(dayjs().startOf('day')) : undefined}
            />
          </Form.Item>
          <Form.Item name="owner_id" label="责任人" rules={[{ required: true, message: '请选择责任人' }]}>
            <Select
              placeholder="请选择责任人"
              options={staff.map((u) => ({ label: u.real_name || u.username, value: u.id }))}
            />
          </Form.Item>
        </Form>
      </Modal>
    </>
  )
}
