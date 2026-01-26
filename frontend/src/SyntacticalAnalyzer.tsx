// SyntacticalAnalyzer.tsx

import React, { useState, useEffect } from 'react'
import './SyntacticalAnalyzer.css'
import EditorPane from './components/EditorPane'

// Type definitions for parse results
interface ParseResult {
  success: boolean
  errors: ParseError[]
  warnings: ParseWarning[]
  parseTree?: string
}

interface ParseError {
  line: number
  message: string
  type: 'syntax' | 'semantic'
}

interface ParseWarning {
  line: number
  message: string
}

type ViewMode = 'parse' | 'errors'

interface SyntacticalAnalyzerProps {
  theme: 'light' | 'dark'
}

const SyntacticalAnalyzer: React.FC<SyntacticalAnalyzerProps> = ({ theme }) => {
  const [code, setCode] = useState<string>('!> type code here\n')
  const [parseResult, setParseResult] = useState<ParseResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [viewMode, setViewMode] = useState<ViewMode>('parse')

  // Calculate stats
  const lineCount = code.split('\n').length
  const tokenCount = 0 // TODO: Get from tokenizer
  const parseStatus = !parseResult ? 'Not analyzed' : parseResult.success ? 'Success' : `Failed (${parseResult.errors?.length || 0} errors)`

  const handleParse = async () => {
    setLoading(true)
    try {
      const response = await fetch('http://localhost:8080/api/analyze', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ code }),
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()

      setParseResult({
        success: data.success,
        errors: data.errors || [],
        warnings: data.warnings || [],
        parseTree: data.parseTree,
      })

      // Auto-switch to errors tab if there are errors
      if (data.errors && data.errors.length > 0) {
        setViewMode('errors')
      } else {
        setViewMode('parse')
      }
    } catch (error) {
      console.error('Parse error:', error)
      setParseResult({
        success: false,
        errors: [
          {
            line: 1,
            message: `Failed to parse code: ${error instanceof Error ? error.message : String(error)}`,
            type: 'syntax',
          },
        ],
        warnings: [],
      })
      setViewMode('errors')
    } finally {
      setLoading(false)
    }
  }

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target?.files?.[0]
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

  const handleCopyOutput = () => {
    if (!parseResult) {
      alert('No output to copy yet')
      return
    }

    let output = ''

    if (viewMode === 'errors') {
      if (parseResult.errors.length === 0 && parseResult.warnings.length === 0) {
        output = 'No errors found!'
      } else {
        if (parseResult.errors.length > 0) {
          output += `ERRORS (${parseResult.errors.length}):\n`
          parseResult.errors.forEach((error) => {
            output += `  Line ${error.line} [${error.type}]: ${error.message}\n`
          })
          output += '\n'
        }
        if (parseResult.warnings.length > 0) {
          output += `WARNINGS (${parseResult.warnings.length}):\n`
          parseResult.warnings.forEach((warning) => {
            output += `  Line ${warning.line}: ${warning.message}\n`
          })
        }
      }
    } else {
      output = parseResult.parseTree || 'No parse tree available'
    }

    navigator.clipboard.writeText(output).then(() => {
      alert('Output copied to clipboard!')
    }).catch(() => {
      alert('Failed to copy to clipboard')
    })
  }

  const renderParseView = () => {
    if (!parseResult) {
      return (
        <div className="placeholder">
          <div className="placeholder-icon">🔍</div>
          <h3>No Parse Results Yet</h3>
          <p>Run the analyzer to see syntax parsing results</p>
        </div>
      )
    }

    return (
      <div className="parse-results">
        <div className="parse-header">
          <div className={`status-badge ${parseResult.success ? 'success' : 'failed'}`}>
            {parseResult.success ? '✓ Parse Successful' : '✕ Parse Failed'}
          </div>
        </div>
        
        <div className="parse-content">
          <h3>Parse Tree</h3>
          <pre className="parse-tree">
            {parseResult.parseTree || 'No parse tree available'}
          </pre>
        </div>
      </div>
    )
  }

  const renderErrorsView = () => {
    if (!parseResult) {
      return (
        <div className="placeholder">
          <div className="placeholder-icon">✓</div>
          <h3>No Errors Yet</h3>
          <p>Run the analyzer to check for errors</p>
        </div>
      )
    }

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
                      <span className="error-location">Line {error.line}</span>
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
                      <span className="warning-location">Line {warning.line}</span>
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
    <div className="syntactical-analyzer">
      <div className="app-grid">
        <div className="pane left">
          <div className="toolbar">
            <button onClick={handleParse} disabled={loading}>
              {loading ? 'Running...' : '▶ RUN'}
            </button>
            <button className="file-btn" title="Save output">
              📄 PSI
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
          <div className="result-tabs">
            <button
              className={`result-tab ${viewMode === 'errors' ? 'active' : ''}`}
              onClick={() => setViewMode('errors')}
            >
              Errors
            </button>
            <button
              className={`result-tab ${viewMode === 'parse' ? 'active' : ''}`}
              onClick={() => setViewMode('parse')}
            >
              Parse
            </button>
            <button
              className="result-tab copy-btn"
              onClick={handleCopyOutput}
              title="Copy output to clipboard"
            >
              📋 Copy
            </button>
          </div>
          
          <div className="outputContainer">
            {viewMode === 'parse' && renderParseView()}
            {viewMode === 'errors' && renderErrorsView()}
          </div>
        </div>
      </div>

      {/* Status Bar */}
      <div className="status-bar">
        <span className="status-item">Lines: {lineCount}</span>
        <span className="status-separator">|</span>
        <span className="status-item">Tokens: {tokenCount}</span>
        <span className="status-separator">|</span>
        <span className="status-item">Parse: {parseStatus}</span>
      </div>
    </div>
  )
}

export default SyntacticalAnalyzer