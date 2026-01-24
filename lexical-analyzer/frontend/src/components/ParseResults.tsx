import React from 'react'
import { ParseResult } from '../types'

type Props = {
  parseResult: ParseResult | null
  code: string
}

export default function ParseResults({ parseResult, code }: Props) {
  if (!parseResult) {
    return (
      <div style={{ 
        flex: 1, 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        color: 'var(--muted)',
        padding: 'var(--spacing-xl)'
      }}>
        <div style={{ textAlign: 'center' }}>
          <div style={{ fontSize: '3rem', marginBottom: 'var(--spacing-md)' }}>🔍</div>
          <div style={{ fontSize: '1rem', fontWeight: 600 }}>No Parse Results Yet</div>
          <div style={{ fontSize: '0.875rem', marginTop: 'var(--spacing-sm)' }}>
            Run the analyzer to see syntax parsing results
          </div>
        </div>
      </div>
    )
  }

  const lines = code.split('\n')

  return (
    <div style={{ flex: 1, overflow: 'auto', padding: 'var(--spacing-lg)' }}>
      {/* Parse Status Header */}
      <div style={{
        padding: 'var(--spacing-md)',
        background: parseResult.success ? 'var(--success)' : 'var(--error)',
        color: '#ffffff',
        borderRadius: 'var(--radius-md)',
        marginBottom: 'var(--spacing-lg)',
        display: 'flex',
        alignItems: 'center',
        gap: 'var(--spacing-md)',
        fontWeight: 700,
        boxShadow: 'var(--shadow-md)'
      }}>
        <span style={{ fontSize: '1.5rem' }}>
          {parseResult.success ? '✓' : '✗'}
        </span>
        <div>
          <div style={{ fontSize: '1rem', letterSpacing: '0.05em' }}>
            {parseResult.success ? 'PARSE SUCCESSFUL' : 'PARSE FAILED'}
          </div>
          {parseResult.message && (
            <div style={{ fontSize: '0.75rem', opacity: 0.9, marginTop: '4px' }}>
              {parseResult.message}
            </div>
          )}
        </div>
        {!parseResult.success && (
          <div style={{ marginLeft: 'auto', fontSize: '0.875rem' }}>
            {parseResult.errors.length} error{parseResult.errors.length !== 1 ? 's' : ''} found
          </div>
        )}
      </div>

      {/* Parse Errors List */}
      {!parseResult.success && parseResult.errors.length > 0 && (
        <div style={{ display: 'flex', flexDirection: 'column', gap: 'var(--spacing-md)' }}>
          {parseResult.errors.map((error, idx) => (
            <div
              key={idx}
              style={{
                background: 'var(--panel-secondary)',
                border: '2px solid var(--error)',
                borderRadius: 'var(--radius-md)',
                overflow: 'hidden',
                transition: 'all var(--transition-base)',
                boxShadow: 'var(--shadow-sm)'
              }}
              onMouseEnter={(e) => {
                e.currentTarget.style.transform = 'translateX(4px)'
                e.currentTarget.style.boxShadow = 'var(--shadow-md)'
              }}
              onMouseLeave={(e) => {
                e.currentTarget.style.transform = 'translateX(0)'
                e.currentTarget.style.boxShadow = 'var(--shadow-sm)'
              }}
            >
              {/* Error Header */}
              <div style={{
                padding: 'var(--spacing-sm) var(--spacing-md)',
                background: 'var(--error)',
                color: '#ffffff',
                display: 'flex',
                alignItems: 'center',
                gap: 'var(--spacing-md)',
                fontSize: '0.75rem',
                fontWeight: 700,
                textTransform: 'uppercase',
                letterSpacing: '0.05em'
              }}>
                <span>Error {idx + 1}</span>
                {error.line > 0 && (
                  <span style={{ opacity: 0.9 }}>
                    Line {error.line}
                    {error.column > 0 && `, Column ${error.column}`}
                  </span>
                )}
              </div>

              {/* Error Message */}
              <div style={{ padding: 'var(--spacing-md)' }}>
                <div style={{
                  color: 'var(--text)',
                  fontSize: '0.9rem',
                  marginBottom: error.expected || error.actual || (error.line > 0 && error.line <= lines.length) ? 'var(--spacing-md)' : 0,
                  fontWeight: 500
                }}>
                  {error.message}
                </div>

                {/* Expected vs Actual */}
                {(error.expected || error.actual) && (
                  <div style={{
                    display: 'flex',
                    gap: 'var(--spacing-md)',
                    marginBottom: 'var(--spacing-md)',
                    fontSize: '0.8rem'
                  }}>
                    {error.expected && (
                      <div style={{
                        background: 'var(--surface-2)',
                        padding: 'var(--spacing-sm)',
                        borderRadius: 'var(--radius-sm)',
                        border: '1px solid var(--border)',
                        flex: 1
                      }}>
                        <div style={{ color: 'var(--muted)', fontSize: '0.7rem', marginBottom: '4px' }}>
                          Expected:
                        </div>
                        <div style={{ color: 'var(--success)', fontFamily: 'monospace' }}>
                          {error.expected}
                        </div>
                      </div>
                    )}
                    {error.actual && (
                      <div style={{
                        background: 'var(--surface-2)',
                        padding: 'var(--spacing-sm)',
                        borderRadius: 'var(--radius-sm)',
                        border: '1px solid var(--border)',
                        flex: 1
                      }}>
                        <div style={{ color: 'var(--muted)', fontSize: '0.7rem', marginBottom: '4px' }}>
                          Actual:
                        </div>
                        <div style={{ color: 'var(--error)', fontFamily: 'monospace' }}>
                          {error.actual}
                        </div>
                      </div>
                    )}
                  </div>
                )}

                {/* Code Context */}
                {error.line > 0 && error.line <= lines.length && (
                  <div style={{
                    background: 'var(--accent-light)',
                    padding: 'var(--spacing-sm) var(--spacing-md)',
                    borderRadius: 'var(--radius-sm)',
                    borderLeft: '4px solid var(--error)',
                    fontFamily: 'monospace',
                    fontSize: '0.85rem'
                  }}>
                    <div style={{
                      color: 'var(--muted)',
                      fontSize: '0.7rem',
                      marginBottom: '4px'
                    }}>
                      Code at Line {error.line}:
                    </div>
                    <div style={{ color: 'var(--text)' }}>
                      {lines[error.line - 1]}
                    </div>
                    {error.column > 0 && (
                      <div style={{ color: 'var(--error)', marginTop: '2px' }}>
                        {' '.repeat(error.column - 1)}^
                      </div>
                    )}
                  </div>
                )}
              </div>
            </div>
          ))}
        </div>
      )}

      {/* Success Message with Details */}
      {parseResult.success && (
        <div style={{
          background: 'var(--panel-secondary)',
          border: '2px solid var(--success)',
          borderRadius: 'var(--radius-md)',
          padding: 'var(--spacing-lg)',
          textAlign: 'center'
        }}>
          <div style={{ fontSize: '3rem', marginBottom: 'var(--spacing-md)' }}>🎉</div>
          <div style={{ fontSize: '1.1rem', fontWeight: 600, color: 'var(--text)', marginBottom: 'var(--spacing-sm)' }}>
            No Syntax Errors Found
          </div>
          <div style={{ fontSize: '0.875rem', color: 'var(--text-secondary)' }}>
            Your code is syntactically correct and ready for execution
          </div>
        </div>
      )}
    </div>
  )
}