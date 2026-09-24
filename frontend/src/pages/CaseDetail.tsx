import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { Button, Card, Descriptions, Tabs, Select, Space, message, Tag } from 'antd'
import { PlusOutlined } from '@ant-design/icons'
import { getCase, changeCaseStatus, assignLawyer } from '@/api/case'
import { getClient } from '@/api/client'
import { completeDeadline, createDeadline, updateDeadline } from '@/api/deadline'
import DocumentList from '@/components/common/DocumentList'
import BillingCard from '@/components/common/BillingCard'
import StatusBadge from '@/components/common/StatusBadge'
import PermissionGuard from '@/components/common/PermissionGuard'
import TimelineItem from '@/components/common/TimelineItem'
import DeadlineList from '@/components/common/DeadlineList'
import DeadlineFormModal, { DeadlineFormValues } from '@/components/common/DeadlineFormModal'
import { useDocumentStore } from '@/stores/documentStore'
import { useBillingStore } from '@/stores/billingStore'
import { useDeadlineStore } from '@/stores/deadlineStore'
import { useUserStore } from '@/stores/userStore'
import { CaseStatusOptions, CaseTypeOptions, CaseStatus } from '@/constants/case'
import type { CaseItem, Client, DeadlineItem } from '@/types'

export default function CaseDetail() {
  const { id } = useParams()
  const caseId = Number(id)
  const [item, setItem] = useState<CaseItem | null>(null)
  const [client, setClient] = useState<Client | null>(null)
  const [status, setStatus] = useState('')
  const [lawyer, setLawyer] = useState<number>()
  const docStore = useDocumentStore()
  const billingStore = useBillingStore()
  const deadlineStore = useDeadlineStore()
  const userStore = useUserStore()
  const [deadlineOpen, setDeadlineOpen] = useState(false)
  const [editingDeadline, setEditingDeadline] = useState<DeadlineItem | null>(null)

  useEffect(() => {
    userStore.fetchLawyers()
    load()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [caseId])

  async function load() {
    const res: any = await getCase(caseId)
    setItem(res.data)
    setStatus(res.data.status)
    if (res.data.client_id) {
      const cr: any = await getClient(res.data.client_id)
      setClient(cr.data.client)
    }
    docStore.fetchByCase(caseId)
    billingStore.fetchByCase(caseId)
    deadlineStore.fetchByCase(caseId)
  }

  async function onStatusChange() {
    await changeCaseStatus(caseId, status)
    message.success('状态已更新')
    load()
  }

  async function onAssign() {
    if (!lawyer) return
    await assignLawyer(caseId, { lead_lawyer_id: lawyer })
    message.success('律师已分配')
    load()
  }

  // 已结案/归档案件只能补录过去日期
  const caseLocked = !!item && (item.status === CaseStatus.CLOSED || item.status === CaseStatus.ARCHIVED)

  async function onDeadlineSubmit(values: DeadlineFormValues) {
    const payload = {
      type: values.type,
      name: values.name.trim(),
      due_at: values.due_at.format('YYYY-MM-DD HH:mm'),
      assignee_id: values.assignee_id,
    }
    if (editingDeadline) {
      await updateDeadline(editingDeadline.id, payload)
      message.success('期限已修改')
    } else {
      await createDeadline({ case_id: caseId, ...payload })
      message.success('期限登记成功')
    }
    setDeadlineOpen(false)
    setEditingDeadline(null)
    deadlineStore.fetchByCase(caseId)
  }

  async function onDeadlineComplete(d: DeadlineItem) {
    await completeDeadline(d.id)
    message.success('期限已标记完成')
    deadlineStore.fetchByCase(caseId)
  }

  if (!item) return null

  return (
    <Card>
      <Space style={{ marginBottom: 16 }}>
        <h2 style={{ margin: 0 }}>{item.case_no} {item.title}</h2>
        <StatusBadge status={item.status} />
        <Tag>{CaseTypeOptions.find((o) => o.value === item.case_type)?.label || item.case_type}</Tag>
      </Space>
      <Tabs
        items={[
          {
            key: 'info',
            label: '案件信息',
            children: (
              <>
                <Descriptions bordered column={2} size="small">
                  <Descriptions.Item label="案号">{item.case_no}</Descriptions.Item>
                  <Descriptions.Item label="状态">{item.status}</Descriptions.Item>
                  <Descriptions.Item label="类型">{item.case_type}</Descriptions.Item>
                  <Descriptions.Item label="主办律师">#{item.lead_lawyer_id}</Descriptions.Item>
                  <Descriptions.Item label="受理日期">{item.accept_date || '-'}</Descriptions.Item>
                  <Descriptions.Item label="结案日期">{item.close_date || '-'}</Descriptions.Item>
                  <Descriptions.Item label="摘要" span={2}>{item.summary || '-'}</Descriptions.Item>
                </Descriptions>
                <PermissionGuard roles={['admin', 'lawyer']}>
                  <Space style={{ marginTop: 16 }}>
                    <Select value={status} style={{ width: 150 }} options={CaseStatusOptions} onChange={setStatus} />
                    <Button type="primary" onClick={onStatusChange}>更新状态</Button>
                  </Space>
                  <Space style={{ marginTop: 8 }}>
                    <Select
                      placeholder="分配主办律师"
                      style={{ width: 180 }}
                      value={lawyer}
                      onChange={setLawyer}
                      options={userStore.lawyers.map((l) => ({ label: l.real_name || l.username, value: l.id }))}
                    />
                    <Button onClick={onAssign}>分配</Button>
                  </Space>
                </PermissionGuard>
              </>
            ),
          },
          {
            key: 'client',
            label: '关联客户',
            children: client ? (
              <Descriptions bordered column={1} size="small">
                <Descriptions.Item label="姓名">{client.name}</Descriptions.Item>
                <Descriptions.Item label="证件号">{client.id_number}</Descriptions.Item>
                <Descriptions.Item label="联系方式">{client.contact}</Descriptions.Item>
                <Descriptions.Item label="地址">{client.address}</Descriptions.Item>
              </Descriptions>
            ) : null,
          },
          {
            key: 'docs',
            label: '文档',
            children: <DocumentList documents={docStore.byCase} />,
          },
          {
            key: 'billings',
            label: '账单',
            children: billingStore.byCase.map((b) => <BillingCard key={b.id} item={b} />),
          },
          {
            key: 'deadlines',
            label: `期限（${deadlineStore.byCase.length}）`,
            children: (
              <Space direction="vertical" style={{ width: '100%' }} size="middle">
                <PermissionGuard roles={['admin', 'lawyer', 'assistant']}>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => { setEditingDeadline(null); setDeadlineOpen(true) }}
                  >
                    登记期限
                  </Button>
                </PermissionGuard>
                <DeadlineList
                  data={deadlineStore.byCase}
                  pagination={false}
                  onEdit={(d) => { setEditingDeadline(d); setDeadlineOpen(true) }}
                  onComplete={onDeadlineComplete}
                />
              </Space>
            ),
          },
          {
            key: 'timeline',
            label: '时间线',
            children: (
              <TimelineItem
                items={[
                  { id: 1, time: item.created_at, text: `案件创建（${item.case_no}）` },
                  { id: 2, time: item.accept_date || item.created_at, text: '案件受理' },
                  { id: 3, time: item.close_date || item.created_at, text: item.close_date ? '案件结案' : '案件进行中' },
                ]}
              />
            ),
          },
        ]}
      />
      <DeadlineFormModal
        open={deadlineOpen}
        caseId={caseId}
        caseLocked={caseLocked}
        editing={editingDeadline}
        onCancel={() => { setDeadlineOpen(false); setEditingDeadline(null) }}
        onSubmit={onDeadlineSubmit}
      />
    </Card>
  )
}
