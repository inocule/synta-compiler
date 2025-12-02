import React from 'react'
import { TokenDTO } from '../types'

export default function OutputTable({ tokens }: { tokens: TokenDTO[] }) {
  const visible = tokens.filter(t => {
    const tt = (t.type || '').toUpperCase()
    if (tt === 'NEWLINE') return false
    if (t.lexeme === "\\n" || t.lexeme === "\n") return false
    return true
  })

  return (
    <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
      <div style={{ padding: '8px 12px', borderBottom: '1px solid rgba(255,255,255,0.02)', fontWeight: 600 }}>Editor Reader</div>
      <div style={{ flex: 1, overflow: 'auto' }}>
        <table>
          <thead>
            <tr>
              <th>Lexeme</th>
              <th>Token</th>
              <th>Line</th>
              <th>Value (chars)</th>
            </tr>
          </thead>
          <tbody>
            {visible.map((t, i) => (
              <tr key={i}>
                <td>{t.lexeme}</td>
                <td>{t.type}</td>
                <td>{t.line}</td>
                <td>{(t.lexeme ?? '').length}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
