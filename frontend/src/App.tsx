import React, { useState, useEffect } from 'react'
import EditorPane from './components/EditorPane'
import OutputTable from './components/OutputTable'
import { analyzeCode } from './api'
import { TokenDTO } from './types'

function App() {
  const [code, setCode] = useState<string>('// type code here\n')
  const [tokens, setTokens] = useState<TokenDTO[]>([])
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState<string | null>(null)
  
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

  async function run() {
    setLoading(true)
    setErr(null)
    try {
      const tok = await analyzeCode(code)
      setTokens(tok)
    } catch (e: any) {
      setErr(e.message || 'Analysis error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <>
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
                {loading ? 'Running...' : 'Run'}
              </button>
            </div>
            <div className="grow" />
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
            <OutputTable tokens={tokens} />
          </div>
        </div>
      </div>
    </>
  )
}

export default App