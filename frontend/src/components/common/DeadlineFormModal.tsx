import { useEffect } from 'react'
import { DatePicker, Form, Input, InputNumber, Modal, Select } from 'antd'
import dayjs, { Dayjs } from 'dayjs'
import { DeadlineTypeOptions } from '@/constants/deadline'
import { useUserStore } from '@/stores/userStore'
import type { DeadlineItem } from '@/types'

export interface DeadlineFormValues {
  case_id?: number
  type: string
  name: string
  due_at: Dayjs
  assignee_id: number
}

interface Props {
  open: boolean
  // 传入 caseId 为案件详情场景（锁定所属案件）；不传则期限中心场景需手填案件 ID
  caseId?: number
  // 已结案/归档：只能补录过去日期
  caseLocked?: boolean
  editing?: DeadlineItem | null
  onCancel: () => void
  onSubmit: (values: DeadlineFormValues) => Promise<void> | void
}

export default function DeadlineFormModal({ open, caseId, caseLocked, editing, onCancel, onSubmit }: Props) {
  const [form] = Form.useForm<DeadlineFormValues>()
  const userStore = useUserStore()

  useEffect(() => {
    if (open) {
      userStore.fetchStaff()
      if (editing) {
        form.setFieldsValue({
          case_id: editing.case_id,
          type: editing.type,
          name: editing.name,
          due_at: dayjs(editing.due_at),
          assignee_id: editing.assignee_id,
        })
      } else {
        form.resetFields()
        form.setFieldsValue({ case_id: caseId, type: 'hearing', due_at: dayjs() })
      }
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, editing])

  async function handleOk() {
    const values = await form.validateFields()
    await onSubmit(values)
  }

  return (
    <Modal
      title={editing ? '修改期限' : '登记期限'}
      open={open}
      onOk={handleOk}
      onCancel={onCancel}
      okText={editing ? '保存' : '登记'}
      destroyOnClose
    >
      <Form form={form} layout="vertical">
        {caseId === undefined && (
          <Form.Item name="case_id" label="所属案件ID" rules={[{ required: true, message: '请输入案件ID' }]}>
            <InputNumber style={{ width: '100%' }} min={1} precision={0} placeholder="输入案件 ID" disabled={!!editing} />
          </Form.Item>
        )}
        <Form.Item name="type" label="类型" rules={[{ required: true, message: '请选择期限类型' }]}>
          <Select options={DeadlineTypeOptions} placeholder="选择期限类型" />
        </Form.Item>
        <Form.Item name="name" label="名称" rules={[{ required: true, whitespace: true, message: '请输入期限名称' }, { max: 200 }]}>
          <Input placeholder="如：一审开庭、上诉期限届满" />
        </Form.Item>
        <Form.Item
          name="due_at"
          label={caseLocked ? '截止时间（已结案/归档仅可补录过去日期）' : '截止时间'}
          rules={[{ required: true, message: '请选择截止时间' }]}
        >
          <DatePicker
            showTime={{ format: 'HH:mm' }}
            format="YYYY-MM-DD HH:mm"
            style={{ width: '100%' }}
            disabledDate={(d) => (caseLocked ? d.startOf('day').isAfter(dayjs().startOf('day')) : false)}
          />
        </Form.Item>
        <Form.Item name="assignee_id" label="责任人" rules={[{ required: true, message: '请选择责任人' }]}>
          <Select
            placeholder="选择责任人"
            options={userStore.staff.map((u) => ({ label: `${u.real_name || u.username}（${u.role === 'lawyer' ? '律师' : '助理'}）`, value: u.id }))}
          />
        </Form.Item>
      </Form>
    </Modal>
  )
}
