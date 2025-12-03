import React, { useEffect, useRef } from 'react'
import Editor, { OnMount } from '@monaco-editor/react'
import type { editor as MonacoEditor, IDisposable } from 'monaco-editor'
import { TokenDTO } from '../types'

type Props = {
  code: string
  setCode: (c: string) => void
  tokens?: TokenDTO[]
  onRun?: () => void
}

export default function EditorPane({ code, setCode, tokens = [], onRun }: Props) {
  const editorRef = useRef<MonacoEditor.IStandaloneCodeEditor | null>(null)
  const monacoRef = useRef<any>(null)
  const decorationsRef = useRef<string[]>([])

  const handleMount: OnMount = (editor, monaco) => {
    editorRef.current = editor
    monacoRef.current = monaco

    // Define Synta dark theme matching the visual identity
    try {
      monaco.editor.defineTheme('synta-dark', {
        base: 'vs-dark',
        inherit: true,
        rules: [
          { token: 'keyword', foreground: 'ff6b9d', fontStyle: 'bold' },      // bright pink
          { token: 'identifier', foreground: '9db4c0' },                      // light blue-gray
          { token: 'string', foreground: '7dd3c0' },                          // teal
          { token: 'number', foreground: 'd4a7ff', fontStyle: 'bold' },       // light purple
          { token: 'comment', foreground: '6d5a5e', fontStyle: 'italic' },    // muted brown
          { token: 'operator', foreground: 'ffd93d', fontStyle: 'bold' },     // yellow
        ],
        colors: {
          'editor.background': '#2d1418',             // dark maroon
          'editor.foreground': '#e8d4d6',             // light silver
          'editorLineNumber.foreground': '#9d8589',   // muted text
          'editorLineNumber.activeForeground': '#c4979b', // accent
          'editorCursor.foreground': '#d32f2f',       // bright red cursor
          'editor.selectionBackground': '#4d2b32',    // dark selection
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
      monaco.editor.setTheme('synta-dark')
    } catch (e) {
      // ignore theme errors
    }

    // Ctrl/Cmd+Enter to run
    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => {
      onRun?.()
    })
  }

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
        const className = mapTokenTypeToClass(t.type)
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
  }, [tokens])

  useEffect(() => {
    return () => {}
  }, [])

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

function mapTokenTypeToClass(type: string) {
  const keywords = new Set([
    'IF','ELIF','ELSE','WHILE','MATCH','RETURN','AWAIT','BREAK','CONTINUE','BIND','CONST','CRAFT','USE','AS','FROM','FN','STRUCT',
    'TRY','CATCH','RAISE','TYPE','CAST','ANY','NONE','TRAIT','INT_TYPE','FLOAT_TYPE','CHAR_TYPE','BOOL_TYPE','STR_TYPE','ASYNC'
  ])
  const t = type.toUpperCase()
  if (keywords.has(t) || t === 'AT_AGENT' || t === 'AT_TASK') return 'tok-keyword'
  if (t === 'STRING') return 'tok-string'
  if (t === 'INTEGER' || t === 'FLOAT') return 'tok-number'
  if (t === 'COMMENT_LINE' || t === 'COMMENT_MULTI') return 'tok-comment'
  if (t === 'ILLEGAL') return 'tok-illegal'
  if (t === 'PLUS' || t === 'MINUS' || t === 'ARROW' || t.endsWith('_ASSIGN') || t === 'EQ' || t === 'NEQ') return 'tok-operator'
  return 'tok-identifier'
}