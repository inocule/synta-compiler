// EditorPane.tsx
import React, { useEffect, useRef } from 'react'
import Editor, { OnMount } from '@monaco-editor/react'
import type { editor as MonacoEditor } from 'monaco-editor'
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
          { token: 'statement-end', foreground: 'ff6b9d', fontStyle: 'bold' },
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
          { token: 'statement-end', foreground: 'c41e3a', fontStyle: 'bold' },
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

      monaco.editor.setTheme(theme === 'light' ? 'synta-light' : 'synta-dark')
    } catch (e) {}

    editor.addCommand(monaco.KeyMod.CtrlCmd | monaco.KeyCode.Enter, () => onRun?.())
  }

  useEffect(() => {
    if (monacoRef.current) {
      try {
        monacoRef.current.editor.setTheme(theme === 'light' ? 'synta-light' : 'synta-dark')
      } catch (e) {}
    }
  }, [theme])

  useEffect(() => {
    const editor = editorRef.current
    const model = editor?.getModel()
    if (!editor || !model) return

    const newDecorations: MonacoEditor.IModelDeltaDecoration[] = tokens
      .filter(t => {
        if (!t.lexeme || t.line <= 0) return false
        const tt = (t.type || '').toUpperCase()
        return tt !== 'NEWLINE' && t.lexeme !== '\\n' && t.lexeme !== '\n'
      })
      .map(t => {
        const startLine = t.line
        const startCol = Math.max(1, t.column)
        const tt = (t.type || '').toUpperCase()
        let endLine = startLine
        let endCol = startCol + Math.max(1, t.lexeme.length)

        // Handle multi-line comments: <! ... !>
        if (tt === 'COMMENT_MULTI') {
          const lineContent = model.getLineContent(startLine)
          const restOfLine = lineContent.substring(startCol - 1)
          
          if (restOfLine.includes('!>')) {
            endCol = startCol + restOfLine.indexOf('!>') + 2
          } else {
            const maxLine = model.getLineCount()
            for (let line = startLine + 1; line <= maxLine; line++) {
              const content = model.getLineContent(line)
              if (content.includes('!>')) {
                endLine = line
                endCol = content.indexOf('!>') + 3
                break
              }
            }
          }
        }

        // Handle single-line comments: !> ...
        if (tt === 'COMMENT_LINE') {
          endCol = model.getLineContent(startLine).length + 1
        }

        // Handle strings
        if (tt === 'STRING') {
          const lineContent = model.getLineContent(startLine)
          const beforeToken = lineContent.substring(0, startCol - 1)
          const quote = (beforeToken.match(/["']$/) || ['"'])[0]
          
          let i = startCol
          while (i < lineContent.length) {
            if (lineContent[i] === quote && lineContent[i - 1] !== '\\') {
              endCol = i + 2
              break
            }
            i++
          }
        }

        return {
          range: { startLineNumber: startLine, startColumn: startCol, endLineNumber: endLine, endColumn: endCol } as any,
          options: { inlineClassName: mapTokenTypeToClass(t.type) }
        }
      })

    try {
      const monaco = monacoRef.current || (window as any).monaco
      if (monaco) {
        // @ts-ignore
        newDecorations.forEach(d => { 
          d.range = new monaco.Range(d.range.startLineNumber, d.range.startColumn, d.range.endLineNumber, d.range.endColumn) 
        })
      }
      decorationsRef.current = editor.deltaDecorations(decorationsRef.current, newDecorations)
    } catch (e) {}
  }, [tokens])

  return (
    <div className="editorContainer" style={{ height: '100%' }}>
      <Editor
        height="100%"
        defaultLanguage="plaintext"
        value={code}
        onChange={(value) => setCode(value ?? '')}
        options={{ 
          minimap: { enabled: false }, 
          fontSize: 13, 
          fontFamily: "Inter, ui-sans-serif",
          lineNumbers: 'on',
          renderWhitespace: 'selection',
          scrollBeyondLastLine: false,
          wordWrap: 'on',
        }}
        onMount={handleMount}
      />
    </div>
  )
}

function mapTokenTypeToClass(type: string) {
  const keywords = new Set([
    'IF', 'ELIF', 'ELSE', 'WHILE', 'FOR', 'MATCH', 'CASE', 'DEFAULT', 'RETURN', 'AWAIT', 'BREAK', 'CONTINUE',
    'BIND', 'CONST', 'CRAFT', 'USE', 'AS', 'FROM', 'FN', 'STRUCT',
    'TRY', 'CATCH', 'RAISE',
    'TYPE', 'CAST', 'ANY', 'NONE', 'TRAIT', 'INT_TYPE', 'FLOAT_TYPE', 'CHAR_TYPE', 'BOOL_TYPE', 'STR_TYPE', 'MAP_TYPE', 'ARRAY_TYPE',
    'ASYNC', 'EMIT', 'LISTEN', 'DISPATCH', 'MERGE', 'TASK', 'CONCURRENT', 'STAGE',
    'WITH', 'THEN', 'DEFER', 'PIPE', 'PASS', 'THROUGH', 'RANGE', 'ALLOW', 'PSEUDO', 'STRATEGY', 'TIMEOUT', 'WINDOW', 'ALERT_THRESHOLD',
    'THINK', 'ASK', 'PROMPT', 'ADAPT', 'CALL_API', 'TRAIN', 'EVALUATE', 'REASON', 'OBSERVE',
    'READ', 'WRITE', 'PRINT', 'LOG', 'SAVE', 'FLOW', 'CONTEXT', 'MEMORY',
    'DEBUG', 'CHECKPOINT', 'TRACE', 'ASSERT', 'CONFIGURE', 'GENERATE_REPORT',
    'AGENT', 'CORE', 'MODEL', 'TOOLS', 'ROLE', 'MODE', 'SYS_PROMPT', 'MAX_CONCURRENT_REQUESTS', 'RETRY_POLICY',
    'CREATE_POOL', 'MAX_WORKERS', 'SUBMIT', 'SUBMIT_DELAYED', 'JOIN', 'NOW', 'EXECUTION_TIME', 'REPORT',
  ])
  
  const t = (type || '').toUpperCase()
  
  if (keywords.has(t)) return 'tok-keyword'
  if (t.startsWith('AT_') || t === 'DECORATOR') return 'tok-keyword'
  if (t === 'STRING') return 'tok-string'
  if (t === 'INTEGER' || t === 'FLOAT') return 'tok-number'
  if (t === 'COMMENT_LINE' || t === 'COMMENT_MULTI') return 'tok-comment'
  if (t === 'STATEMENT_END') return 'tok-statement-end'
  if (t === 'ILLEGAL') return 'tok-illegal'
  if (['PLUS', 'MINUS', 'MULTIPLY', 'DIVIDE', 'MODULO', 'ARROW', 'FAT_ARROW'].includes(t)) return 'tok-operator'
  if (t.endsWith('_ASSIGN') || t === 'ASSIGN' || t === 'BIND_ASSIGN') return 'tok-operator'
  if (['EQ', 'NEQ', 'LT', 'GT', 'LTE', 'GTE', 'AND', 'OR', 'NOT', 'INCREMENT', 'DECREMENT'].includes(t)) return 'tok-operator'
  
  return 'tok-identifier'
}