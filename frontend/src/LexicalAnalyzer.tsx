// LexicalAnalyzer.tsx

import React, { useState } from 'react'
import EditorPane from './components/EditorPane'
import OutputTable from './components/OutputTable'
import { analyzeCode } from './api'
import { TokenDTO } from './types'

// Define the possible output modes
type ViewMode = 'table' | 'lineByLine' | 'singleLine' | 'codeBlock'

interface LexicalAnalyzerProps {
  theme: 'light' | 'dark'
}

const LexicalAnalyzer: React.FC<LexicalAnalyzerProps> = ({ theme }) => {
  const [code, setCode] = useState<string>('// type code here\n')
  const [tokens, setTokens] = useState<TokenDTO[]>([])
  const [loading, setLoading] = useState(false)
  const [err, setErr] = useState<string | null>(null)

  // State for view mode, default to 'table'
  const [viewMode, setViewMode] = useState<ViewMode>('table')
  
  // State for line navigation in singleLine mode
  const [currentLine, setCurrentLine] = useState(1)

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
          
          {/* View Switcher UI */}
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
  )
}

export default LexicalAnalyzer