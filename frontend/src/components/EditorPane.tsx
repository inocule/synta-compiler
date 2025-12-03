import React, { useEffect, useRef } from 'react'
import Editor, { OnMount } from '@monaco-editor/react'
import type { editor as MonacoEditor, IDisposable } from 'monaco-editor'
import { TokenDTO } from '../types'

type Props = {
  code: string
  setCode: (c: string) => void
  tokens?: TokenDTO[]
  onRun?: () => void
  theme?: 'light' | 'dark'
}

export default function EditorPane({ code, setCode, tokens = [], onRun, theme = 'dark' }: Props) {
  const editorRef = useRef<MonacoEditor.IStandaloneCodeEditor | null>(null)
  const monacoRef = useRef<any>(null)
  const decorationsRef = useRef<string[]>([])

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor
    monacoRef.current = monaco

    // Define Synta dark theme
    try {
      monaco.editor.defineTheme('synta-dark', {
        base: 'vs-dark',
        inherit: true,
        rules: [
          { token: 'keyword', foreground: 'ff6b9d', fontStyle: 'bold' },
          { token: 'identifier', foreground: '9db4c0' },
          { token: 'string', foreground: '7dd3c0' },
          { token: 'number', foreground: 'd4a7ff', fontStyle: 'bold' },
          { token: 'comment', foreground: '6d5a5e', fontStyle: 'italic' },
          { token: 'operator', foreground: 'ffd93d', fontStyle: 'bold' },
        ],
        colors: {
          'editor.background': '#2d1418',
          'editor.foreground': '#e8d4d6',
          'editorLineNumber.foreground': '#9d8589',
          'editorLineNumber.activeForeground': '#c4979b',
          'editorCursor.foreground': '#d32f2f',
          'editor.selectionBackground': '#4d2b32',
          'editor.inactiveSelectionBackground': '#3d1d24',
          'editor.lineHighlightBackground': '#3d1d24',
          'editor.lineHighlightBorder': '#4d2429',
          'editorWhitespace.foreground': '#4d2429',
          'editorIndentGuide.background': '#4d2429',
          'editorIndentGuide.activeBackground': '#6d343d',
          'editorBracketMatch.background': '#4d2b32',
          'editorBracketMatch.border': '#d32f2f',
        },
      })

      // Define Synta light theme
      monaco.editor.defineTheme('synta-light', {
        base: 'vs',
        inherit: true,
        rules: [
          { token: 'keyword', foreground: 'c41e3a', fontStyle: 'bold' },
          { token: 'identifier', foreground: '4a5568' },
          { token: 'string', foreground: '059669' },
          { token: 'number', foreground: '7c3aed', fontStyle: 'bold' },
          { token: 'comment', foreground: 'a0aec0', fontStyle: 'italic' },
          { token: 'operator', foreground: 'd97706', fontStyle: 'bold' },
        ],
        colors: {
          'editor.background': '#faf7f5',
          'editor.foreground': '#1a1a1a',
          'editorLineNumber.foreground': '#a0aec0',
          'editorLineNumber.activeForeground': '#c41e3a',
          'editorCursor.foreground': '#c41e3a',
          'editor.selectionBackground': '#fce7e9',
          'editor.inactiveSelectionBackground': '#f5e5e7',
          'editor.lineHighlightBackground': '#fef5f5',
          'editor.lineHighlightBorder': '#fce7e9',
          'editorWhitespace.foreground': '#e5e7eb',
          'editorIndentGuide.background': '#e5e7eb',
          'editorIndentGuide.activeBackground': '#cbd5e0',
          'editorBracketMatch.background': '#fce7e9',
          'editorBracketMatch.border': '#c41e3a',
        },
      })

      // Set initial theme
      monaco.editor.setTheme(theme === 'light' ? 'synta-light' : 'synta-dark')
    } catch (e) {
      // ignore theme errors
    }

    // Ctrl/Cmd+Enter to run
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      onRun?.()
    })
  }

  // Switch theme when theme prop changes
  useEffect(() => {
    const monaco = monacoRef.current
    if (monaco) {
      try {
        monaco.editor.setTheme(theme === 'light' ? 'synta-light' : 'synta-dark')
      } catch (e) {
        // ignore theme errors
      }
    }
  }, [theme])

  // Re-apply decorations when theme or tokens change
  useEffect(() => {
    const editor = editorRef.current
    if (!editor) return

    // Convert tokens to Monaco inline decorations
    const newDecorations: MonacoEditor.IModelDeltaDecoration[] = tokens
      .filter(t => {
        if (!t.lexeme || t.line <= 0) return false
        const tt = (t.type || '').toUpperCase()
        // Ignore standalone NEWLINE tokens
        if (tt === 'NEWLINE') return false
        if (t.lexeme === '\\n' || t.lexeme === '\n') return false
        return true
      })
      .map(t => {
        const startLine = t.line
        const startCol = Math.max(1, t.column)
        const endCol = startCol + Math.max(1, t.lexeme.length)
        const className = mapTokenTypeToClass(t.type, theme === 'light')
        return {
          range: { startLineNumber: startLine, startColumn: startCol, endLineNumber: startLine, endColumn: endCol } as any,
          options: { inlineClassName: className }
        }
      })

    // Apply decorations
    try {
      const monaco = monacoRef.current || (window as any).monaco
      if (monaco) {
        // @ts-ignore
        newDecorations.forEach(d => { d.range = new monaco.Range(d.range.startLineNumber, d.range.startColumn, d.range.endLineNumber, d.range.endColumn) })
      }
      decorationsRef.current = editor.deltaDecorations(decorationsRef.current, newDecorations)
    } catch (e) {
      // ignore if monaco not ready
    }
  }, [tokens, theme])

  return (
    <div className="editorContainer" style={{ height: '100%' }}>
      <Editor
        height="100%"
        defaultLanguage="plaintext"
        value={code}
        onChange={(value) => setCode(value ?? '')}
        options={{ 
          minimap: { enabled: false }, 
          fontSize: 14, 
          fontFamily: "Inter, ui-monospace, 'Cascadia Code', 'Fira Code', monospace",
          lineHeight: 24,
          padding: { top: 16, bottom: 16 },
          scrollBeyondLastLine: false,
          smoothScrolling: true,
          cursorBlinking: 'smooth',
          cursorSmoothCaretAnimation: 'on',
          fontLigatures: true,
        }}
        onMount={handleMount}
      />
    </div>
  )
}

function mapTokenTypeToClass(type: string, lightMode: boolean = false) {
  const keywords = new Set([
    'IF','ELIF','ELSE','WHILE','MATCH','RETURN','AWAIT','BREAK','CONTINUE','BIND','CONST','CRAFT','USE','AS','FROM','FN','STRUCT',
    'TRY','CATCH','RAISE','TYPE','CAST','ANY','NONE','TRAIT','INT_TYPE','FLOAT_TYPE','CHAR_TYPE','BOOL_TYPE','STR_TYPE','ASYNC'
  ])
  const t = type.toUpperCase()
  const prefix = lightMode ? 'tok-light-' : 'tok-'
  
  if (keywords.has(t) || t === 'AT_AGENT' || t === 'AT_TASK') return `${prefix}keyword`
  if (t === 'STRING') return `${prefix}string`
  if (t === 'INTEGER' || t === 'FLOAT') return `${prefix}number`
  if (t === 'COMMENT_LINE' || t === 'COMMENT_MULTI') return `${prefix}comment`
  if (t === 'ILLEGAL') return `${prefix}illegal`
  if (t === 'PLUS' || t === 'MINUS' || t === 'ARROW' || t.endsWith('_ASSIGN') || t === 'EQ' || t === 'NEQ') return `${prefix}operator`
  return `${prefix}identifier`
}