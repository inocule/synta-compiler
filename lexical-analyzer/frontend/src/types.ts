export type TokenDTO = {
  lexeme: string
  type: string
  line: number
  column: number
  value?: string
  extra?: Record<string, any>
}

export type ParseError = {
  message: string
  line: number
  column: number
  expected?: string
  actual?: string
}

export type ParseResult = {
  success: boolean
  errors: ParseError[]
  message?: string
}

// Analysis mode types
export type AnalysisMode = 'syntax' | 'lexical'
export type LexicalView = 'table' | 'lineByLine' | 'singleLine' | 'codeBlock'