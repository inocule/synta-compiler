import React, { useState, useEffect } from 'react'

type LoadingStage = 'tokenizing' | 'parsing' | 'complete'

type Props = {
  stage: LoadingStage
  onComplete?: () => void
}

export default function LoadingIndicator({ stage, onComplete }: Props) {
  const [progress, setProgress] = useState(0)

  useEffect(() => {
    if (stage === 'tokenizing') {
      setProgress(0)
      const timer = setInterval(() => {
        setProgress(prev => Math.min(prev + 10, 50))
      }, 50)
      return () => clearInterval(timer)
    } else if (stage === 'parsing') {
      setProgress(50)
      const timer = setInterval(() => {
        setProgress(prev => Math.min(prev + 10, 90))
      }, 50)
      return () => clearInterval(timer)
    } else if (stage === 'complete') {
      setProgress(100)
      setTimeout(() => onComplete?.(), 300)
    }
  }, [stage, onComplete])

  const stageLabels: Record<LoadingStage, string> = {
    tokenizing: '🔤 Tokenizing code...',
    parsing: '🔍 Parsing syntax...',
    complete: '✓ Analysis complete!'
  }

  return (
    <div style={{
      position: 'fixed',
      top: 0,
      left: 0,
      right: 0,
      height: '3px',
      background: 'var(--surface-2)',
      zIndex: 1000,
      overflow: 'hidden'
    }}>
      <div style={{
        height: '100%',
        width: `${progress}%`,
        background: stage === 'complete' 
          ? 'var(--success)' 
          : 'linear-gradient(90deg, var(--red-gradient-start), var(--red-gradient-end))',
        transition: 'width 0.3s ease-out',
        boxShadow: '0 0 10px currentColor'
      }} />
      
      {/* Loading Label */}
      <div style={{
        position: 'absolute',
        top: '10px',
        right: 'var(--spacing-xl)',
        background: 'var(--panel)',
        padding: '6px 12px',
        borderRadius: 'var(--radius-md)',
        border: '1px solid var(--border)',
        fontSize: '0.75rem',
        fontWeight: 600,
        color: 'var(--text)',
        boxShadow: 'var(--shadow-md)',
        display: 'flex',
        alignItems: 'center',
        gap: '8px'
      }}>
        <span>{stageLabels[stage]}</span>
        {stage !== 'complete' && (
          <div className="loading-spinner" style={{
            width: '12px',
            height: '12px',
            border: '2px solid var(--border)',
            borderTopColor: 'var(--accent)',
            borderRadius: '50%',
            animation: 'spin 0.6s linear infinite'
          }} />
        )}
      </div>

      <style>{`
        @keyframes spin {
          to { transform: rotate(360deg); }
        }
      `}</style>
    </div>
  )
}