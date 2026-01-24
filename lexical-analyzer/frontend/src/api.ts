import { TokenDTO, ParseResult } from './types'

export async function analyzeCode(code: string): Promise<TokenDTO[]> {
  const res = await fetch('http://localhost:8080/api/analyze', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code }),
  })
  const data = await res.json()
  if (!res.ok || data.error) throw new Error(data.error || 'analysis failed')
  return data.tokens as TokenDTO[]
}

export async function parseCode(code: string): Promise<ParseResult> {
  try {
    const res = await fetch('http://localhost:8080/api/parse', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ code }),
    })
    const data = await res.json()
    if (!res.ok) throw new Error(data.error || 'parse failed')
    return data as ParseResult
  } catch (error: any) {
    // Return mock error for testing when backend isn't ready
    return {
      success: false,
      errors: [{
        message: error.message || 'Parse endpoint not available yet',
        line: 0,
        column: 0
      }],
      message: 'Backend not ready - this is expected during development'
    }
  }
}