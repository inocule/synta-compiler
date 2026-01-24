import { TokenDTO } from './types'

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
