export type StudyRecord = {
  id: number
  date: string
  studyMinutes: number
  focusScore: number
  note: string
}

export type RecordListResponse = {
  dailyAverageMinutes: number
  streakDays: number
  completionRate: number
  records: StudyRecord[]
}

const records: StudyRecord[] = [
  { id: 1, date: '2026-04-16', studyMinutes: 70, focusScore: 80, note: 'Go の handler を整理' },
  { id: 2, date: '2026-04-17', studyMinutes: 95, focusScore: 88, note: 'Vue 画面の再配置' },
  { id: 3, date: '2026-04-18', studyMinutes: 55, focusScore: 72, note: '認証周りの確認' },
  { id: 4, date: '2026-04-19', studyMinutes: 110, focusScore: 91, note: '一覧 UI の調整' },
  { id: 5, date: '2026-04-20', studyMinutes: 85, focusScore: 84, note: 'API レスポンスの確認' },
  { id: 6, date: '2026-04-21', studyMinutes: 100, focusScore: 89, note: 'frontend 構成の整理' },
]

export const fetchRecordList = async (): Promise<RecordListResponse> => {
  const totalMinutes = records.reduce((sum, record) => sum + record.studyMinutes, 0)
  const totalFocus = records.reduce((sum, record) => sum + record.focusScore, 0)

  return {
    dailyAverageMinutes: Math.round(totalMinutes / records.length),
    streakDays: records.length,
    completionRate: Math.round(totalFocus / records.length),
    records,
  }
}
