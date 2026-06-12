export interface HostFileInfo {
  name: string
  isDirty: boolean
}

export interface SaveHostFileRequest {
  name: string
  content: string
}
