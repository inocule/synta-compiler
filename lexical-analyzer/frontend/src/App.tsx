// App.tsx

import React, { useState, useEffect } from 'react'
import EditorPane from './components/EditorPane'
import OutputTable from './components/OutputTable'
import ParseResults from './components/ParseResults'
import StatusBar from './components/StatusBar'
import LoadingIndicator from './components/LoadingIndicator'
import { analyzeCode, parseCode } from './api'
import { TokenDTO, ParseResult, AnalysisMode, LexicalView } from './types'

type LoadingStage = 'tokenizing' | 'parsing' | 'complete' | null

function App() {
  const [code, setCode] = useState<string>('// type code here\n')
  const [tokens, setTokens] = useState<TokenDTO[]>([])
  const [parseResult, setParseResult] = useState<ParseResult | null>(null)
  const [loading, setLoading] = useState(false)
  const [loadingStage, setLoadingStage] = useState<LoadingStage>(null)
  const [analysisTime, setAnalysisTime] = useState<number | undefined>(undefined)
  const [err, setErr] = useState<string | null>(null)

  // HYBRID MODE: Top-level analysis mode
  const [analysisMode, setAnalysisMode] = useState<AnalysisMode>('syntax')
  
  // Lexical sub-view (only used when analysisMode is 'lexical')
  const [lexicalView, setLexicalView] = useState<LexicalView>('table')
  
  const [currentLine, setCurrentLine] = useState(1)
  const [theme, setTheme] = useState<'light' | 'dark'>('light')

  // Load saved theme on mount 
  useEffect(() => {
    const savedTheme = (localStorage.getItem('theme') as 'light' | 'dark') || 'light'
    setTheme(savedTheme)
    document.body.setAttribute('data-theme', savedTheme)
  }, [])

  const toggleTheme = () => {
    const newTheme = theme === 'light' ? 'dark' : 'light'
    setTheme(newTheme)
    document.body.setAttribute('data-theme', newTheme)
    localStorage.setItem('theme', newTheme)
  }

  const handleLineChange = (direction: 'up' | 'down') => {
    const maxLine = code.split('\n').length
    setCurrentLine(prev => {
        let newLine = prev
        if (direction === 'up') {
            newLine = Math.min(maxLine, prev + 1)
        } else if (direction === 'down') {
            newLine = Math.max(1, prev - 1)
        }
        return newLine
    })
  }

  async function run() {
    setLoading(true)
    setLoadingStage('tokenizing')
    setErr(null)
    const startTime = performance.now()
    
    try {
      // Tokenizing stage
      const tok = await analyzeCode(code)
      setTokens(tok)
      
      // Parsing stage
      setLoadingStage('parsing')
      const parseRes = await parseCode(code)
      setParseResult(parseRes)
      
      // Complete
      setLoadingStage('complete')
      const endTime = performance.now()
      setAnalysisTime(Math.round(endTime - startTime))
      
      setCurrentLine(1)
    } catch (e: any) {
      setErr(e.message || 'Analysis error')
      setLoadingStage(null)
    } finally {
      setTimeout(() => {
        setLoading(false)
        setLoadingStage(null)
      }, 500)
    }
  }

  const handleCreateNewFile = () => {
    const content = String.raw`# AI Agents - Detailed Markdown Explanation

## Metadata
- **File**: ai_agents.synta
- **Author**: inocule on 2025-12-06
- **Last Modified**: 2025-12-06 02:15:00

## Purpose
Define AI agents, their tools, execution settings, and demonstrate example tasks.
Provide detailed context for AI reasoning, debugging, and concurrency tracking.

## 🏗️ System Architecture


┌─────────────────────────────────────────────────────────────┐
│                     AI AGENT SYSTEM                         │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────────┐              ┌──────────────────┐     │
│  │   AICoder Agent  │              │ ClaudeOpus Agent │     │
│  │   (LOCAL/HYBRID) │              │     (CLOUD)      │     │
│  └────────┬─────────┘              └────────┬─────────┘     │
│           │                                 │               │
│           │         ┌──────────────┐        │               │
│           └────────►│  TASK POOL   │◄───────┘               │
│                     │  (4 Workers) │                        │
│                     └──────┬───────┘                        │
│                            │                                │
│                     ┌──────▼───────┐                        │
│                     │ INTENT LOG & │                        │
│                     │  AI INSIGHTS │                        │
│                     └──────────────┘                        │
└─────────────────────────────────────────────────────────────┘

## 🤖 Agent Configurations

### STEP 1: AICoder Agent Definition

╔════════════════════════════════════════════════════════╗
║                   AICODER AGENT                        ║
╠════════════════════════════════════════════════════════╣
║ Role: Coding Assistance                                ║
║ Model: Llama 3.1 8B                                    ║
║ VRAM: ~16-20GB (FP16)                                  ║
║ Context: 2k-4k tokens                                  ║
║ Mode: HYBRID (local primary, cloud fallback)           ║
╠════════════════════════════════════════════════════════╣
║ TOOLS:                                                 ║
║  ├─ 📂 GitHub MCP                                     ║
║  │   └─ Fetch code, track changes, version control     ║
║  ├─ 💬 slm_chatbot                                    ║
║  │   └─ Natural language code interaction              ║
║  └─ 📄 pdf_scanner                                     ║
║      └─ Parse code from PDFs                           ║
╠════════════════════════════════════════════════════════╣
║ DEBUG MONITORING:                                      ║
║  • VRAM usage tracking                                 ║
║  • Task size validation                                ║
║  • Malformed code edge cases                           ║
╚════════════════════════════════════════════════════════╝

### STEP 2: ClaudeOpus Agent Definition

╔════════════════════════════════════════════════════════╗
║                CLAUDEOPUS AGENT                        ║
╠════════════════════════════════════════════════════════╣
║ Role: Reasoning & Text Generation                      ║
║ Model: Claude Opus 4.5                                 ║
║ Context: Large context windows                         ║
║ Mode: CLOUD (scalable)                                 ║
╠════════════════════════════════════════════════════════╣
║ TOOLS:                                                 ║
║  ├─ 📝 text_summarizer                                ║
║  ├─ 🔍 code_explainer                                  ║
║  └─ 💡 idea_generator                                  ║
╠════════════════════════════════════════════════════════╣
║ SETTINGS:                                              ║
║  • Max concurrent: 3 requests                          ║
║  • Timeout: 60 seconds                                 ║
║  • Retry: Linear backoff                               ║
╠════════════════════════════════════════════════════════╣
║ DEBUG MONITORING:                                      ║
║  • Queue depth tracking                                ║
║  • Latency measurements                                ║
║  • Type error detection                                ║
║  • Hallucination analysis                              ║
╚════════════════════════════════════════════════════════╝
`;

    const element = document.createElement('a');
    const file = new Blob([content], { type: 'text/plain' });
    element.href = URL.createObjectURL(file);
    element.download = 'ai_agents_psi.md';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
    
    setCode('')
    setTokens([])
    setParseResult(null)
    setCurrentLine(1)
    setErr(null)
  }

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    if (!file.name.endsWith('.synta')) {
      setErr('Error: Only .synta files are accepted')
      return
    }

    const reader = new FileReader()
    reader.onload = (e) => {
      try {
        const content = e.target?.result as string
        setCode(content)
        setTokens([])
        setParseResult(null)
        setCurrentLine(1)
        setErr(null)
      } catch (error) {
        setErr('Error reading file')
      }
    }
    reader.onerror = () => {
      setErr('Error reading file')
    }
    reader.readAsText(file)
  }

  return (
    <>
      {/* Loading Indicator */}
      {loadingStage && loadingStage !== 'complete' && (
        <LoadingIndicator 
          stage={loadingStage as 'tokenizing' | 'parsing' | 'complete'} 
          onComplete={() => setLoadingStage(null)} 
        />
      )}

      {/* Theme Toggle Button */}
      <button 
        className="theme-toggle" 
        onClick={toggleTheme}
        aria-label="Toggle theme"
      >
        <div className="theme-toggle-slider">
          {theme === 'light' ? '☼' : '☾'}
        </div>
      </button>

      <div className="app-grid">
        <div className="pane left">
          <div className="toolbar">
            <div className="flex">
              <button onClick={run} disabled={loading}>
                {loading ? 'Analyzing...' : 'Run'}
              </button>
              <button 
                onClick={handleCreateNewFile}
                title="Create a new .synta file"
                className="file-btn"
              >
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
            
            <div className="grow" /> 
            
            {/* HYBRID MODE SELECTOR */}
            <div className="mode-selector-container">
              {/* Top-Level Mode Buttons */}
              <div className="mode-selector">
                <button
                  className={`mode-btn ${analysisMode === 'syntax' ? 'active' : ''}`}
                  onClick={() => setAnalysisMode('syntax')}
                  title="View Syntax Analysis Results"
                >
                  <span className="mode-icon">🔍</span>
                  <span>PARSE</span>
                </button>
                <button
                  className={`mode-btn ${analysisMode === 'lexical' ? 'active' : ''}`}
                  onClick={() => setAnalysisMode('lexical')}
                  title="View Lexical Analysis Results"
                >
                  <span className="mode-icon">🔤</span>
                  <span>TOKENS</span>
                  <span className="dropdown-arrow">▼</span>
                </button>
              </div>

              {/* Lexical Sub-Tabs (only shown when lexical mode is active) */}
              {analysisMode === 'lexical' && (
                <div className="lexical-subtabs">
                  <button
                    className={`subtab-btn ${lexicalView === 'singleLine' ? 'active' : ''}`}
                    onClick={() => setLexicalView('singleLine')}
                    title="Single Line Navigation"
                  >
                    LINE
                  </button>
                  <button
                    className={`subtab-btn ${lexicalView === 'lineByLine' ? 'active' : ''}`}
                    onClick={() => setLexicalView('lineByLine')}
                    title="All Lines View"
                  >
                    ALL
                  </button>
                  <button
                    className={`subtab-btn ${lexicalView === 'table' ? 'active' : ''}`}
                    onClick={() => setLexicalView('table')}
                    title="Classic Token Table"
                  >
                    TABLE
                  </button>
                  <button
                    className={`subtab-btn ${lexicalView === 'codeBlock' ? 'active' : ''}`}
                    onClick={() => setLexicalView('codeBlock')}
                    title="Code Blocks View"
                  >
                    BLOCKS
                  </button>
                </div>
              )}
            </div>
            
            {err && <div className="err">{err}</div>}
          </div>
          <div className="editor">
            <EditorPane 
              code={code} 
              setCode={setCode} 
              tokens={tokens} 
              onRun={run}
              theme={theme}
            />
          </div>
        </div>
        <div className="pane right">
          <div className="outputContainer">
            {analysisMode === 'syntax' ? (
              <ParseResults parseResult={parseResult} code={code} />
            ) : (
              <OutputTable 
                tokens={tokens} 
                code={code} 
                viewMode={lexicalView} 
                currentLine={currentLine} 
                onLineChange={handleLineChange} 
              />
            )}
          </div>
        </div>
      </div>

      {/* Status Bar */}
      <StatusBar 
        code={code}
        tokens={tokens}
        parseResult={parseResult}
        analysisTime={analysisTime}
      />
    </>
  )
}

export default App