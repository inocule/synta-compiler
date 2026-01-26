// SyntacticalAnalyzer.tsx

import React, { useState, useEffect, useRef } from 'react'
import './SyntacticalAnalyzer.css'
import EditorPane from './components/EditorPane'
import type { TokenDTO } from './types'

// Type definitions for parse results
interface ParseResult {
  success: boolean
  errors: ParseError[]
  warnings: ParseWarning[]
  parseTree?: string
  astTree?: ASTNode
}

interface ASTNode {
  type: string
  name?: string
  value?: string
  children?: ASTNode[]
  line?: number
  column?: number
  metadata?: Record<string, any>
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

type ViewMode = 'ast' | 'parse' | 'errors'

interface SyntacticalAnalyzerProps {
  theme: 'light' | 'dark'
}

const SyntacticalAnalyzer: React.FC<SyntacticalAnalyzerProps> = ({ theme }) => {
  const [code, setCode] = useState<string>('!> type code here\n')
  const [parseResult, setParseResult] = useState<ParseResult | null>(null)
  const [tokens, setTokens] = useState<TokenDTO[]>([])
  const [loading, setLoading] = useState(false)
  const [viewMode, setViewMode] = useState<ViewMode>('ast')
  const [selectedNode, setSelectedNode] = useState<ASTNode | null>(null)
  const [expandedNodes, setExpandedNodes] = useState<Set<string>>(new Set())
  const [searchTerm, setSearchTerm] = useState('')

  // Calculate stats
  const lineCount = code.split('\n').length
  const tokenCount = tokens.length
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

      setTokens(data.tokens || [])
      setParseResult({
        success: data.success,
        errors: data.errors || [],
        warnings: data.warnings || [],
        parseTree: data.parseTree,
        astTree: data.astTree || parseTreeToAST(data.parseTree),
      })

      // Auto-switch to errors tab if there are errors, otherwise show AST
      if (data.errors && data.errors.length > 0) {
        setViewMode('errors')
      } else {
        setViewMode('ast')
      }

      // Expand all nodes by default
      if (data.astTree) {
        const allIds = getAllNodeIds(data.astTree)
        setExpandedNodes(new Set(allIds))
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

  // Convert parse tree string to AST structure (fallback if backend doesn't provide AST)
  const parseTreeToAST = (parseTree: string | undefined): ASTNode | undefined => {
    if (!parseTree) return undefined
    
    // Simple heuristic parser - in production, this should come from backend
    const lines = parseTree.split('\n').filter(line => line.trim())
    const root: ASTNode = { type: 'Program', children: [] }
    
    return root
  }

  const getAllNodeIds = (node: ASTNode, prefix = '0'): string[] => {
    const ids = [prefix]
    node.children?.forEach((child, idx) => {
      ids.push(...getAllNodeIds(child, `${prefix}-${idx}`))
    })
    return ids
  }

  const toggleNode = (nodeId: string) => {
    const newExpanded = new Set(expandedNodes)
    if (newExpanded.has(nodeId)) {
      newExpanded.delete(nodeId)
    } else {
      newExpanded.add(nodeId)
    }
    setExpandedNodes(newExpanded)
  }

  const expandAll = () => {
    if (parseResult?.astTree) {
      const allIds = getAllNodeIds(parseResult.astTree)
      setExpandedNodes(new Set(allIds))
    }
  }

  const collapseAll = () => {
    setExpandedNodes(new Set(['0']))
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
      setTokens([])
      setParseResult(null)
      setSelectedNode(null)
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
    } else if (viewMode === 'ast') {
      output = JSON.stringify(parseResult.astTree, null, 2)
    } else {
      output = parseResult.parseTree || 'No parse tree available'
    }

    navigator.clipboard.writeText(output).then(() => {
      alert('Output copied to clipboard!')
    }).catch(() => {
      alert('Failed to copy to clipboard')
    })
  }

  const getNodeIcon = (type: string): string => {
    const iconMap: Record<string, string> = {
      'Program': '📄',
      'Declaration': '🔷',
      'Function': '⚡',
      'Variable': '📦',
      'Statement': '📝',
      'Expression': '🔢',
      'BinaryOp': '⚙️',
      'UnaryOp': '🔧',
      'Literal': '💎',
      'Identifier': '🏷️',
      'IfStatement': '🔀',
      'WhileStatement': '🔄',
      'ForStatement': '➰',
      'ReturnStatement': '↩️',
      'CallExpression': '📞',
      'ArrayLiteral': '📚',
      'MapLiteral': '🗺️',
      'AsyncStatement': '⏳',
      'AwaitExpression': '⏰',
    }
    return iconMap[type] || '🔹'
  }

  const getNodeColor = (type: string): string => {
    // Unique color palette for different node types
    const colorMap: Record<string, string> = {
      'Program': 'var(--node-program)',
      'Declaration': 'var(--node-declaration)',
      'Function': 'var(--node-function)',
      'Variable': 'var(--node-variable)',
      'Statement': 'var(--node-statement)',
      'Expression': 'var(--node-expression)',
      'BinaryOp': 'var(--node-operator)',
      'UnaryOp': 'var(--node-operator)',
      'Literal': 'var(--node-literal)',
      'Identifier': 'var(--node-identifier)',
      'IfStatement': 'var(--node-control)',
      'WhileStatement': 'var(--node-control)',
      'ForStatement': 'var(--node-control)',
      'SwitchStatement': 'var(--node-control)',
      'ReturnStatement': 'var(--node-return)',
      'CallExpression': 'var(--node-call)',
      'ArrayLiteral': 'var(--node-collection)',
      'MapLiteral': 'var(--node-collection)',
      'AsyncStatement': 'var(--node-async)',
      'AwaitExpression': 'var(--node-async)',
    }
    return colorMap[type] || 'var(--node-default)'
  }

  const renderASTNode = (node: ASTNode, nodeId: string, depth: number = 0): JSX.Element => {
    const isExpanded = expandedNodes.has(nodeId)
    const hasChildren = node.children && node.children.length > 0
    const isSelected = selectedNode === node
    const matchesSearch = !searchTerm || 
      node.type.toLowerCase().includes(searchTerm.toLowerCase()) ||
      node.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
      node.value?.toLowerCase().includes(searchTerm.toLowerCase())

    if (!matchesSearch && searchTerm) {
      return <></>
    }

    return (
      <div key={nodeId} className="ast-node-container" style={{ marginLeft: depth > 0 ? '24px' : '0' }}>
        <div
          className={`ast-node ${isSelected ? 'selected' : ''}`}
          onClick={() => setSelectedNode(node)}
          style={{ 
            '--node-color': getNodeColor(node.type),
            animationDelay: `${depth * 0.05}s`
          } as React.CSSProperties}
        >
          <div className="ast-node-header">
            {hasChildren && (
              <button
                className={`expand-btn ${isExpanded ? 'expanded' : ''}`}
                onClick={(e) => {
                  e.stopPropagation()
                  toggleNode(nodeId)
                }}
              >
                <svg width="12" height="12" viewBox="0 0 12 12">
                  <path d="M4 2 L8 6 L4 10" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" />
                </svg>
              </button>
            )}
            
            <span className="node-icon">{getNodeIcon(node.type)}</span>
            
            <div className="node-info">
              <span className="node-type">{node.type}</span>
              {node.name && <span className="node-name">{node.name}</span>}
              {node.value && <span className="node-value">= {node.value}</span>}
            </div>
            
            {(node.line || node.column) && (
              <span className="node-location">
                {node.line}:{node.column}
              </span>
            )}
          </div>
        </div>

        {hasChildren && isExpanded && (
          <div className="ast-children">
            {node.children!.map((child, idx) => 
              renderASTNode(child, `${nodeId}-${idx}`, depth + 1)
            )}
          </div>
        )}
      </div>
    )
  }

  const renderASTView = () => {
    if (!parseResult?.astTree) {
      return (
        <div className="placeholder">
          <div className="placeholder-icon">🌳</div>
          <h3>No AST Available</h3>
          <p>Run the analyzer to see the abstract syntax tree</p>
        </div>
      )
    }

    return (
      <div className="ast-view">
        <div className="ast-controls">
          <div className="ast-search">
            <svg width="16" height="16" viewBox="0 0 16 16" fill="none">
              <circle cx="7" cy="7" r="5" stroke="currentColor" strokeWidth="2"/>
              <path d="M11 11 L15 15" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
            </svg>
            <input
              type="text"
              placeholder="Search nodes..."
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              className="search-input"
            />
          </div>
          <div className="ast-actions">
            <button onClick={expandAll} className="ast-action-btn" title="Expand All">
              <svg width="16" height="16" viewBox="0 0 16 16">
                <path d="M2 6 L8 12 L14 6" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
              </svg>
            </button>
            <button onClick={collapseAll} className="ast-action-btn" title="Collapse All">
              <svg width="16" height="16" viewBox="0 0 16 16">
                <path d="M2 10 L8 4 L14 10" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
              </svg>
            </button>
          </div>
        </div>

        <div className="ast-content-wrapper">
          <div className="ast-tree-container">
            <div className="ast-tree">
              {renderASTNode(parseResult.astTree, '0')}
            </div>
          </div>

          <div className="ast-details-panel">
            <div className="details-header">
              <svg width="20" height="20" viewBox="0 0 20 20">
                <rect x="2" y="2" width="16" height="16" rx="2" fill="none" stroke="currentColor" strokeWidth="2"/>
                <path d="M6 10 L14 10 M10 6 L10 14" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
              </svg>
              <h3>Node Details</h3>
            </div>

            {selectedNode ? (
              <div className="details-content">
                <div className="detail-section">
                  <div className="detail-label">TYPE</div>
                  <div className="detail-value type-badge" style={{ '--node-color': getNodeColor(selectedNode.type) } as React.CSSProperties}>
                    {getNodeIcon(selectedNode.type)} {selectedNode.type}
                  </div>
                </div>

                {selectedNode.name && (
                  <div className="detail-section">
                    <div className="detail-label">NAME</div>
                    <div className="detail-value code">{selectedNode.name}</div>
                  </div>
                )}

                {selectedNode.value && (
                  <div className="detail-section">
                    <div className="detail-label">VALUE</div>
                    <div className="detail-value code">{selectedNode.value}</div>
                  </div>
                )}

                {(selectedNode.line || selectedNode.column) && (
                  <div className="detail-section">
                    <div className="detail-label">LOCATION</div>
                    <div className="detail-value">
                      Line {selectedNode.line}, Column {selectedNode.column}
                    </div>
                  </div>
                )}

                {selectedNode.metadata && Object.keys(selectedNode.metadata).length > 0 && (
                  <div className="detail-section">
                    <div className="detail-label">METADATA</div>
                    <div className="detail-value metadata">
                      {Object.entries(selectedNode.metadata).map(([key, value]) => (
                        <div key={key} className="metadata-item">
                          <span className="metadata-key">{key}:</span>
                          <span className="metadata-value">{JSON.stringify(value)}</span>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {selectedNode.children && selectedNode.children.length > 0 && (
                  <div className="detail-section">
                    <div className="detail-label">CHILDREN</div>
                    <div className="detail-value">{selectedNode.children.length} child node(s)</div>
                  </div>
                )}
              </div>
            ) : (
              <div className="details-placeholder">
                <div className="placeholder-icon-small">👆</div>
                <p>Select a node to view details</p>
              </div>
            )}
          </div>
        </div>
      </div>
    )
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
            <button onClick={handleParse} disabled={loading} className="run-btn">
              {loading ? (
                <>
                  <span className="spinner"></span>
                  Running...
                </>
              ) : (
                <>
                  <svg width="14" height="14" viewBox="0 0 14 14">
                    <path d="M3 2 L11 7 L3 12 Z" fill="currentColor" />
                  </svg>
                  RUN
                </>
              )}
            </button>
            <button className="toolbar-btn" title="Save output">
              <svg width="16" height="16" viewBox="0 0 16 16">
                <path d="M2 2 L14 2 L14 14 L2 14 Z M5 2 L5 6 L11 6 L11 2 M8 9 L8 14" fill="none" stroke="currentColor" strokeWidth="1.5"/>
              </svg>
              PSI
            </button>
            <label className="toolbar-btn" title="Upload a .synta file">
              <svg width="16" height="16" viewBox="0 0 16 16">
                <path d="M3 5 L3 13 L13 13 L13 5 M6 1 L6 8 M6 1 L3 4 M6 1 L9 4" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
              </svg>
              OPEN
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
              tokens={tokens} 
              onRun={handleParse}
              theme={theme}
            />
          </div>
        </div>

        <div className="pane right">
          <div className="result-tabs">
            <button
              className={`result-tab ${viewMode === 'ast' ? 'active' : ''}`}
              onClick={() => setViewMode('ast')}
            >
              <svg width="16" height="16" viewBox="0 0 16 16">
                <circle cx="8" cy="3" r="2" fill="currentColor"/>
                <circle cx="4" cy="13" r="2" fill="currentColor"/>
                <circle cx="12" cy="13" r="2" fill="currentColor"/>
                <path d="M8 5 L8 8 M8 8 L4 11 M8 8 L12 11" stroke="currentColor" strokeWidth="1.5"/>
              </svg>
              Abstract Syntax Tree
            </button>
            <button
              className={`result-tab ${viewMode === 'parse' ? 'active' : ''}`}
              onClick={() => setViewMode('parse')}
            >
              <svg width="16" height="16" viewBox="0 0 16 16">
                <rect x="2" y="2" width="12" height="2" fill="currentColor"/>
                <rect x="2" y="7" width="8" height="2" fill="currentColor"/>
                <rect x="2" y="12" width="10" height="2" fill="currentColor"/>
              </svg>
              Parse
            </button>
            <button
              className={`result-tab ${viewMode === 'errors' ? 'active' : ''}`}
              onClick={() => setViewMode('errors')}
            >
              <svg width="16" height="16" viewBox="0 0 16 16">
                <circle cx="8" cy="8" r="6" fill="none" stroke="currentColor" strokeWidth="1.5"/>
                <path d="M8 4 L8 9 M8 11 L8 12" stroke="currentColor" strokeWidth="2" strokeLinecap="round"/>
              </svg>
              Errors
            </button>
            <button
              className="result-tab copy-btn"
              onClick={handleCopyOutput}
              title="Copy output to clipboard"
            >
              <svg width="16" height="16" viewBox="0 0 16 16">
                <rect x="4" y="4" width="8" height="10" rx="1" fill="none" stroke="currentColor" strokeWidth="1.5"/>
                <path d="M10 4 L10 2 L4 2 L4 10" fill="none" stroke="currentColor" strokeWidth="1.5"/>
              </svg>
              Copy
            </button>
          </div>
          
          <div className="outputContainer">
            {viewMode === 'ast' && renderASTView()}
            {viewMode === 'parse' && renderParseView()}
            {viewMode === 'errors' && renderErrorsView()}
          </div>
        </div>
      </div>

      {/* Status Bar */}
      <div className="status-bar">
        <span className="status-item">
          <svg width="12" height="12" viewBox="0 0 12 12">
            <path d="M2 3 L10 3 M2 6 L10 6 M2 9 L10 9" stroke="currentColor" strokeWidth="1.5"/>
          </svg>
          Lines: {lineCount}
        </span>
        <span className="status-separator">|</span>
        <span className="status-item">
          <svg width="12" height="12" viewBox="0 0 12 12">
            <circle cx="6" cy="6" r="4" fill="currentColor"/>
          </svg>
          Tokens: {tokenCount}
        </span>
        <span className="status-separator">|</span>
        <span className="status-item">
          <svg width="12" height="12" viewBox="0 0 12 12">
            <path d="M6 2 L10 6 L6 10 L2 6 Z" fill="currentColor"/>
          </svg>
          Parse: {parseStatus}
        </span>
      </div>
    </div>
  )
}

export default SyntacticalAnalyzer