export type TokenDTO = {
  lexeme: string
  type: string
  line: number
  column: number
  value?: string
  extra?: Record<string, any>
}
