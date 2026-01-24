import React from 'react'
import { TokenDTO, ParseResult } from '../types'

type Props = {
  code: string
  tokens: TokenDTO[]
  parseResult: ParseResult | null
  analysisTime?: number
}

export default function StatusBar({ code, tokens, parseResult, analysisTime }: Props) {
  const lines = code.split('\n').length
  const totalTokens = tokens.filter(t => {
    const tt = (t.type || '').toUpperCase()
    return tt !== 'NEWLINE' && t.lexeme !== '\\n' && t.lexeme !== '\n'
  }).length

  const parseStatus = !parseResult 
    ? { icon: '⏸', text: 'Not analyzed', color: 'var(--muted)' }
    : parseResult.success
    ? { icon: '✓', text: 'Valid', color: 'var(--success)' }
    : { icon: '✗', text: `${parseResult.errors.length} error${parseResult.errors.length !== 1 ? 's' : ''}`, color: 'var(--error)' }

  return (
    <div style={{
      position: 'fixed',
      bottom: 0,
      left: 0,
      right: 0,
      height: '32px',
      background: 'var(--table-header)',
      borderTop: '1px solid var(--border)',
      display: 'flex',
      alignItems: 'center',
      padding: '0 var(--spacing-lg)',
      gap: 'var(--spacing-lg)',
      fontSize: '0.75rem',
      fontWeight: 500,
      color: 'var(--text-secondary)',
      zIndex: 100,
      boxShadow: '0 -2px 8px rgba(0, 0, 0, 0.05)',
      transition: 'all var(--transition-base)'
    }}>
      {/* Lines Count */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
        <span style={{ opacity: 0.7 }}>📄</span>
        <span>Lines: <strong style={{ color: 'var(--text)' }}>{lines}</strong></span>
      </div>

      <div style={{ width: '1px', height: '16px', background: 'var(--border)', opacity: 0.5 }} />

      {/* Tokens Count */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
        <span style={{ opacity: 0.7 }}>🔤</span>
        <span>Tokens: <strong style={{ color: 'var(--text)' }}>{totalTokens}</strong></span>
      </div>

      <div style={{ width: '1px', height: '16px', background: 'var(--border)', opacity: 0.5 }} />

      {/* Parse Status */}
      <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
        <span style={{ color: parseStatus.color }}>{parseStatus.icon}</span>
        <span>Parse: <strong style={{ color: parseStatus.color }}>{parseStatus.text}</strong></span>
      </div>

      {/* Analysis Time */}
      {analysisTime !== undefined && (
        <>
          <div style={{ width: '1px', height: '16px', background: 'var(--border)', opacity: 0.5 }} />
          <div style={{ display: 'flex', alignItems: 'center', gap: '6px' }}>
            <span style={{ opacity: 0.7 }}>⚡</span>
            <span><strong style={{ color: 'var(--text)' }}>{analysisTime}ms</strong></span>
          </div>
        </>
      )}

      {/* Spacer */}
      <div style={{ flex: 1 }} />

      {/* Theme Indicator */}
      <div style={{ 
        display: 'flex', 
        alignItems: 'center', 
        gap: '6px',
        opacity: 0.6,
        fontSize: '0.7rem'
      }}>
        <span>SYNTA Compiler v1.0</span>
      </div>
    </div>
  )
}