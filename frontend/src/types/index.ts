export interface User {
  id: number
  username: string
  real_name: string
  role: string
  license_no: string
  email: string
  phone: string
  avatar: string
  created_at: string
}

export interface Client {
  id: number
  name: string
  id_number: string
  contact: string
  address: string
  remark: string
  created_at: string
}

export interface CaseItem {
  id: number
  case_no: string
  title: string
  case_type: string
  status: string
  accept_date: string | null
  close_date: string | null
  summary: string
  client_id: number
  lead_lawyer_id: number
  co_lawyer_ids: number[]
  created_at: string
}

export interface DocumentItem {
  id: number
  title: string
  file_type: string
  file_url: string
  upload_time: string
  case_id: number
  uploader_id: number
  created_at: string
}

export interface Billing {
  id: number
  bill_no: string
  billing_type: string
  amount: number | string
  status: string
  case_id: number
  client_id: number
  invoice_info: string
  created_at: string
}

export interface Deadline {
  id: number
  case_id: number
  type: string
  name: string
  due_at: string
  due_date: string
  assignee_id: number
  status: string
  completed_by_id: number | null
  completed_at: string | null
  created_at: string
  updated_at: string
}

// 期限中心/案件期限列表项（后端 JOIN 用户与案件后的视图）
export interface DeadlineItem extends Deadline {
  case_no: string
  case_title: string
  case_status: string
  assignee_name: string
  completer_name: string
}

export interface AuditLog {
  id: number
  operator_id: number
  operator_name: string
  action: string
  entity_type: string
  entity_id: string
  detail: string
  ip: string
  created_at: string
}

export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
}
