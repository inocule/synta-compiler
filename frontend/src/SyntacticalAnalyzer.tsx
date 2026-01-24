// SyntacticalAnalyzer.tsx

import React, { useState, useEffect } from 'react'
import './SyntacticalAnalyzer.css'
import EditorPane from './components/EditorPane'

// Type definitions for parse results
interface ParseResult {
  success: boolean
  errors: ParseError[]
  warnings: ParseWarning[]
}

interface ParseError {
  line: number
  column: number
  message: string
  type: 'syntax' | 'semantic'
}

interface ParseWarning {
  line: number
  column: number
  message: string
}

type ViewMode = 'summary' | 'errors'

interface SyntacticalAnalyzerProps {
  theme: 'light' | 'dark'
}

const SyntacticalAnalyzer: React.FC<SyntacticalAnalyzerProps> = ({ theme }) => {
  const [code, setCode] = useState<string>('// Enter your .synta code here\n')
  const [parseResult, setParseResult] = useState<ParseResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [viewMode, setViewMode] = useState<ViewMode>('summary')

  // Theme state management removed - now controlled by parent App.tsx

  const handleParse = async () => {
    setLoading(true)
    try {
      // TODO: Replace with actual backend API call
      // const result = await fetch('/api/parse', {
      //   method: 'POST',
      //   headers: { 'Content-Type': 'application/json' },
      //   body: JSON.stringify({ code })
      // }).then(res => res.json())
      
      // Mock result for now
      await new Promise(resolve => setTimeout(resolve, 800))
      
      const mockResult: ParseResult = {
        success: true,
        errors: [],
        warnings: []
      }
      
      setParseResult(mockResult)
    } catch (error) {
      console.error('Parse error:', error)
      setParseResult({
        success: false,
        errors: [{
          line: 1,
          column: 1,
          message: 'Failed to parse code',
          type: 'syntax'
        }],
        warnings: []
      })
    } finally {
      setLoading(false)
    }
  }

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    if (!file.name.endsWith('.synta')) {
      alert('Error: Only .synta files are accepted')
      return
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      const content = e.target?.result as string
      setCode(content)
      setParseResult(null)
    }
    reader.readAsText(file)
  }

  const renderSummaryView = () => {
    if (!parseResult) return null

    const totalErrors = parseResult.errors.length
    const totalWarnings = parseResult.warnings.length
    const syntaxErrors = parseResult.errors.filter(e => e.type === 'syntax').length
    const semanticErrors = parseResult.errors.filter(e => e.type === 'semantic').length

    return (
      <div className="summary-container">
        <div className="summary-header">
          <h2>Parse Analysis Summary</h2>
          <div className={`status-badge ${parseResult.success ? 'success' : 'failed'}`}>
            {parseResult.success ? '✓ Parse Successful' : '✕ Parse Failed'}
          </div>
        </div>

        <div className="summary-stats">
          <div className="stat-card errors">
            <div className="stat-number">{totalErrors}</div>
            <div className="stat-label">Total Errors</div>
            {totalErrors > 0 && (
              <div className="stat-breakdown">
                <span>{syntaxErrors} syntax</span>
                <span>{semanticErrors} semantic</span>
              </div>
            )}
          </div>

          <div className="stat-card warnings">
            <div className="stat-number">{totalWarnings}</div>
            <div className="stat-label">Warnings</div>
          </div>
        </div>

        {totalErrors > 0 && (
          <div className="quick-errors">
            <h3>Recent Errors</h3>
            {parseResult.errors.slice(0, 3).map((error, idx) => (
              <div key={idx} className="quick-error-item">
                <span className="quick-error-location">[{error.line}:{error.column}]</span>
                <span className="quick-error-message">{error.message}</span>
              </div>
            ))}
            {totalErrors > 3 && (
              <div className="more-errors">
                +{totalErrors - 3} more errors (switch to ERRORS view)
              </div>
            )}
          </div>
        )}
      </div>
    )
  }

  const renderErrorsView = () => {
    if (!parseResult) return null

    return (
      <div className="errors-container">
        {parseResult.errors.length === 0 && parseResult.warnings.length === 0 ? (
          <div className="no-errors">
            <div className="success-icon">✓</div>
            <h3>No Errors Found</h3>
            <p>Your code parsed successfully!</p>
          </div>
        ) : (
          <>
            {parseResult.errors.length > 0 && (
              <div className="error-section">
                <h3 className="error-heading">
                  <span className="error-icon">✕</span>
                  Errors ({parseResult.errors.length})
                </h3>
                {parseResult.errors.map((error, idx) => (
                  <div key={idx} className="error-item">
                    <div className="error-header">
                      <span className="error-type">{error.type}</span>
                      <span className="error-location">[{error.line}:{error.column}]</span>
                    </div>
                    <div className="error-message">{error.message}</div>
                  </div>
                ))}
              </div>
            )}
            
            {parseResult.warnings.length > 0 && (
              <div className="warning-section">
                <h3 className="warning-heading">
                  <span className="warning-icon">⚠</span>
                  Warnings ({parseResult.warnings.length})
                </h3>
                {parseResult.warnings.map((warning, idx) => (
                  <div key={idx} className="warning-item">
                    <div className="warning-header">
                      <span className="warning-location">[{warning.line}:{warning.column}]</span>
                    </div>
                    <div className="warning-message">{warning.message}</div>
                  </div>
                ))}
              </div>
            )}
          </>
        )}
      </div>
    )
  }

  return (
    <div className="app-grid">
      <div className="pane left">
        <div className="toolbar">
          <div className="flex">
            <button onClick={handleParse} disabled={loading}>
              {loading ? 'Parsing...' : '▶ Parse'}
            </button>
            <label className="file-btn-label" title="Upload a .synta file">
              📂 OPEN
              <input 
                type="file" 
                accept=".synta" 
                onChange={handleFileUpload}
                style={{ display: 'none' }}
              />
            </label>
          </div>
          
          <div className="grow" />
          
          <div className="view-switch-container">
            <button
              className={`view-switch-btn ${viewMode === 'summary' ? 'active' : ''}`}
              onClick={() => setViewMode('summary')}
              title="Summary View"
            >
              SUMMARY
            </button>
            <button
              className={`view-switch-btn ${viewMode === 'errors' ? 'active' : ''}`}
              onClick={() => setViewMode('errors')}
              title="Errors & Warnings"
            >
              ERRORS
            </button>
          </div>
        </div>
        
        <div className="editor">
          <EditorPane 
            code={code} 
            setCode={setCode} 
            tokens={[]} 
            onRun={handleParse}
            theme={theme}
          />
        </div>
      </div>

      <div className="pane right">
        <div className="outputContainer">
          {!parseResult ? (
            <div className="placeholder">
              <div className="placeholder-icon">🌲</div>
              <h3>Ready to Parse</h3>
              <p>Enter code and click "Parse" to analyze syntax</p>
            </div>
          ) : (
            <div className="output-content">
              {viewMode === 'summary' && renderSummaryView()}
              {viewMode === 'errors' && renderErrorsView()}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default SyntacticalAnalyzer