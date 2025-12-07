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

  // New state for view mode, default to 'table'
  const [viewMode, setViewMode] = useState<ViewMode>('table')
  
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

  const handleCreateNewFile = () => {
  // Pseudocode content for AI Agents explanation
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

## 🔄 Execution Flow Diagram

START
  │
  ▼
┌────────────────────────────────────┐
│ Initialize System                  │
│ • Load agents (AICoder, ClaudeOpus)│
│ • Configure debug settings         │
│ • Create task_pool (4 workers)     │
└─────────────┬──────────────────────┘
              │
              ▼
        ┌─────────────┐
        │ Task Queue  │
        └──────┬──────┘
               │
        ┌──────▼──────────────────────┐
        │  FOR each task in queue:    │
        └──────┬──────────────────────┘
               │
        ┌──────▼─────────────────────────────┐
        │ IF task.type == "code_fix":        │
        └──────┬─────────────────────────────┘
               │
        ┌──────▼───────────────────────┐
    YES │ Route to AICoder Agent       │
        │  ┌────────────────────────┐  │
        │  │ 1. Parse input         │  │
        │  │ 2. Analyze syntax      │  │
        │  │ 3. Apply fixes         │  │
        │  │ 4. Validate output     │  │
        │  └────────────────────────┘  │
        │  result = AICoder.process()  │
        │  log(result)                 │
        └──────┬───────────────────────┘
               │
        ┌──────▼─────────────────────────────┐
    NO  │ ELSE IF task.type == "summary":    │
        └──────┬─────────────────────────────┘
               │
        ┌──────▼────────────────────────┐
    YES │ Route to ClaudeOpus Agent     │
        │  ┌────────────────────────┐   │
        │  │ 1. Extract key points  │   │
        │  │ 2. Generate summary    │   │
        │  │ 3. Format output       │   │
        │  └────────────────────────┘   │
        │  result = ClaudeOpus.process()│
        │  log(result)                  │
        └──────┬────────────────────────┘
               │
        ┌──────▼──────────────────────┐
        │ Store in intent_log         │
        │ Analyze in ai_insights      │
        └──────┬──────────────────────┘
               │
        ┌──────▼──────────────────────┐
        │ END task iteration          │
        └──────┬──────────────────────┘
               │
               ▼
┌──────────────────────────────────┐
│ All tasks complete               │
│ Generate execution report        │
└──────────────────────────────────┘
  │
  ▼
END

## Agents Detailed Explanation

### STEP 1: Define AICoder agent
- Role: coding assistance
- Tools:
  - GitHub MCP (fetch code, track changes, integrate version control)
  - slm_chatbot (natural language code interaction)
  - pdf_scanner (parse code from PDFs)
- Model: local Llama 3.1 8B (~16–20GB VRAM for FP16, 2k–4k token window)
- Mode: hybrid (local primary, cloud fallback)
- Debug Notes: Monitor VRAM usage, task size, and malformed code edge cases

### STEP 2: Define ClaudeOpus agent
- Role: reasoning & text generation
- Tools: text_summarizer, code_explainer, idea_generator
- Model: cloud Claude Opus 4.5 (large context windows, scalable)
- Mode: cloud
- Settings: max 3 concurrent requests, timeout 60s, linear backoff retry
- Debug Notes: Monitor queue, latency, type errors, hallucinations

### STEP 3: Create Example Task
- Intent: Demonstrate agent interaction
- Actions:
  1. Use AICoder to fix syntax errors
  2. Use ClaudeOpus to summarize code
  3. Print results
- Edge Cases: cloud timeout, large/malformed code, task pool limits
- Debug Notes: Outputs logged in intent_log, analyzed in ai_insights

## Execution Flow (Pseudocode)
START
  Initialize agents (AICoder, ClaudeOpus)
  Configure debug settings
  Create task_pool with 4 workers
  FOR each task:
      IF task.type == "code_fix":
          result = AICoder.process(task.input)
          log(result)
      ELSE IF task.type == "summary":
          result = ClaudeOpus.process(task.input)
          log(result)
      ENDIF
  ENDFOR
END

## Revision History
- 2025-12-06: Initial generation
- AI Insight: Agents designed for hybrid local/cloud execution; task example demonstrates intent and concurrency tracking
`;

    // Trigger download of the pseudocode file
    const element = document.createElement('a');
    const file = new Blob([content], { type: 'text/plain' });
    element.href = URL.createObjectURL(file);
    element.download = 'ai_agents_psi.md';
    document.body.appendChild(element);
    element.click();
    document.body.removeChild(element);
    
    // Also clear the editor
    setCode('')
    setTokens([])
    setCurrentLine(1)
    setErr(null)
  }

  const handleFileUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (!file) return

    // Validate file extension
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