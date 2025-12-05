// App.tsx

import React, { useState, useEffect } from 'react'
import EditorPane from './components/EditorPane'
import OutputTable from './components/OutputTable'
import { analyzeCode } from './api'
import { TokenDTO } from './types'

// Define the possible output modes (added 'codeBlock')
type ViewMode = 'table' | 'lineByLine' | 'singleLine' | 'codeBlock'

function App() {
  const [code, setCode] = useState<string>('// type code here\n')
  const [tokens, setTokens] = useState<TokenDTO[]>([])
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState<string | null>(null)

  // New state for view mode, default to 'singleLine' based on the required feature
  const [viewMode, setViewMode] = useState<ViewMode>('singleLine')
  
  // State for line navigation in singleLine mode
  const [currentLine, setCurrentLine] = useState(1)

  // Theme state management 
  const [theme, setTheme] = useState<'light' | 'dark'>('light')

  // Load saved theme on mount 
  useEffect(() => {
    const savedTheme = (localStorage.getItem('theme') as 'light' | 'dark') || 'light'
    setTheme(savedTheme)
    document.body.setAttribute('data-theme', savedTheme)
  }, [])

  // Toggle theme function 
  const toggleTheme = () => {
    const newTheme = theme === 'light' ? 'dark' : 'light'
    setTheme(newTheme)
    document.body.setAttribute('data-theme', newTheme)
    localStorage.setItem('theme', newTheme)
  }

  // Handle line change for singleLine mode
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
    setErr(null)
    try {
      const tok = await analyzeCode(code)
      setTokens(tok)
      
      // Reset currentLine to 1 after a successful run
      setCurrentLine(1)
    } catch (e: any) {
      setErr(e.message || 'Analysis error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
      {/* Theme Toggle Button  */}
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
                {loading ? 'Running...' : 'Run'}
              </button>
            </div>
            
            <div className="grow" /> 
            
            {/* START: Updated View Switcher UI with SingleLine */}
            <div className="view-switch-container">
                <button
                    className={`view-switch-btn ${viewMode === 'singleLine' ? 'active' : ''}`}
                    onClick={() => setViewMode('singleLine')}
                    title="Single Line Navigation"
                >
                    LINE
                </button>
                <button
                    className={`view-switch-btn ${viewMode === 'lineByLine' ? 'active' : ''}`}
                    onClick={() => setViewMode('lineByLine')}
                    title="All Lines View"
                >
                    ALL
                </button>
                <button
                    className={`view-switch-btn ${viewMode === 'table' ? 'active' : ''}`}
                    onClick={() => setViewMode('table')}
                    title="Classic Token Table"
                >
                    TABLE
                </button>
                <button
                  className={`view-switch-btn ${viewMode === 'codeBlock' ? 'active' : ''}`}
                  onClick={() => setViewMode('codeBlock')}
                  title="Code Blocks View"
                >
                  BLOCKS
                </button>
            </div>
            {/* END: Updated View Switcher UI */}
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
            {/* Pass the new line state and handler */}
            <OutputTable 
                tokens={tokens} 
                code={code} 
                viewMode={viewMode} 
                currentLine={currentLine} 
                onLineChange={handleLineChange} 
            />
          </div>
        </div>
      </div>
    </>
  )
}

export default App