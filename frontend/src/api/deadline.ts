import request from '@/utils/request'

export interface DeadlinePayload {
  type: string
  name: string
  due_at: string
  assignee_id: number
}

export function listDeadlines(params: {
  page?: number
  page_size?: number
  view?: string
  assignee_id?: number
  case_id?: number
}) {
  return request.get('/deadlines', { params })
}

export function listDeadlinesByCase(caseId: number) {
  return request.get(`/cases/${caseId}/deadlines`)
}

export function createDeadline(data: { case_id: number } & DeadlinePayload) {
  return request.post('/deadlines', data)
}

export function updateDeadline(id: number, data: DeadlinePayload) {
  return request.put(`/deadlines/${id}`, data)
}

export function completeDeadline(id: number) {
  return request.post(`/deadlines/${id}/complete`)
}
