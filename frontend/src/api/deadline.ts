import request from '@/utils/request'

export function listDeadlines(params: { page?: number; page_size?: number; view?: string; owner_id?: number }) {
  return request.get('/deadlines', { params })
}

export function listDeadlinesByCase(caseId: number) {
  return request.get(`/deadlines/by-case/${caseId}`)
}

export function createDeadline(data: { case_id: number; deadline_type: string; name: string; due_date: string; owner_id: number }) {
  return request.post('/deadlines', data)
}

export function updateDeadline(id: number, data: { deadline_type?: string; name?: string; due_date?: string; owner_id?: number }) {
  return request.put(`/deadlines/${id}`, data)
}

export function completeDeadline(id: number) {
  return request.post(`/deadlines/${id}/complete`)
}
