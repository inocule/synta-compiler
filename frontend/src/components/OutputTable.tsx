// OutputTable.tsx
import React, { useMemo } from 'react'
import { TokenDTO } from '../types'

type ViewMode = 'table' | 'lineByLine' | 'singleLine' | 'codeBlock'
type Props = {
    tokens: TokenDTO[]
    code: string
    viewMode: ViewMode
    currentLine: number
    onLineChange: (direction: 'up' | 'down') => void
    onViewModeChange: (mode: ViewMode) => void
}


const filterNewlines = (tokens: TokenDTO[]) => tokens.filter(t => {
    if (!t) return false
    const tt = (t.type || '').toUpperCase()
    const isNewline = tt === 'NEWLINE' || t.lexeme === "\\n" || t.lexeme === "\n"
    // Also filter out tokens without valid line numbers
    return !isNewline && t.line !== undefined && t.line !== null && t.line > 0
})

const calculateWhitespaceCount = (code: string) => 
    code.split('\n').map(line => (line.match(/\s/g) || []).length)

const calculateTokenCounts = (tokens: TokenDTO[]) => {
    const counts: Record<number, number> = {}
    filterNewlines(tokens).forEach(t => {
        const line = Math.max(1, t.line || 1)
        counts[line] = (counts[line] || 0) + 1
    })
    return counts
}

type CodeBlock = {
    id: number
    type: 'curly_brace' | 'statement' | 'single_comment' | 'multi_comment'
    startLine: number
    endLine: number
    code: string
    tokens: TokenDTO[]
    tokenCount: number
    wsCount: number
}

const parseCodeBlocks = (code: string, tokens: TokenDTO[]): CodeBlock[] => {
    const lines = code.split('\n')
    const blocks: CodeBlock[] = []
    let blockId = 0
    let i = 0

    const tokensByLine = filterNewlines(tokens).reduce((acc, t) => {
        if (!acc[t.line]) acc[t.line] = []
        acc[t.line].push(t)
        return acc
    }, {} as Record<number, TokenDTO[]>)

    const countWS = (str: string) => str.split('').filter(c => /\s/.test(c)).length

    while (i < lines.length) {
        const line = lines[i]
        const trimmed = line.trim()
        const lineNum = i + 1

        if (!trimmed) { i++; continue }

        // Multi-line comment: <! ... !>
        if (trimmed.startsWith('<!') || trimmed.startsWith('/*') || trimmed.startsWith('/*')) {
            const startLine = lineNum
            let endLine = lineNum
            let blockCode = line
            const closers = ['!>', '*/', '*/']
            
            if (!closers.some(c => line.includes(c))) {
                for (let j = i + 1; j < lines.length; j++) {
                    blockCode += '\n' + lines[j]
                    if (closers.some(c => lines[j].includes(c))) {
                        endLine = j + 1
                        i = j
                        break
                    }
                }
            }

            const blockTokens = []
            for (let l = startLine; l <= endLine; l++) {
                if (tokensByLine[l]) blockTokens.push(...tokensByLine[l])
            }

            blocks.push({
                id: blockId++,
                type: 'multi_comment',
                startLine, endLine, code: blockCode,
                tokens: blockTokens,
                tokenCount: blockTokens.length,
                wsCount: countWS(blockCode)
            })
            i++
            continue
        }

        // Single-line comment: !> or //
        if (trimmed.startsWith('!>') || trimmed.startsWith('//')) {
            blocks.push({
                id: blockId++,
                type: 'single_comment',
                startLine: lineNum, endLine: lineNum, code: line,
                tokens: tokensByLine[lineNum] || [],
                tokenCount: (tokensByLine[lineNum] || []).length,
                wsCount: countWS(line)
            })
            i++
            continue
        }

        // Curly brace block
        if (trimmed.includes('{')) {
            const startLine = lineNum
            let endLine = lineNum
            let blockCode = line
            let braceCount = (line.match(/{/g) || []).length - (line.match(/}/g) || []).length

            if (braceCount > 0) {
                for (let j = i + 1; j < lines.length; j++) {
                    blockCode += '\n' + lines[j]
                    braceCount += (lines[j].match(/{/g) || []).length - (lines[j].match(/}/g) || []).length
                    if (braceCount === 0) {
                        endLine = j + 1
                        i = j
                        break
                    }
                }
            }

            const blockTokens = []
            for (let l = startLine; l <= endLine; l++) {
                if (tokensByLine[l]) blockTokens.push(...tokensByLine[l])
            }

            blocks.push({
                id: blockId++,
                type: 'curly_brace',
                startLine, endLine, code: blockCode,
                tokens: blockTokens,
                tokenCount: blockTokens.length,
                wsCount: countWS(blockCode)
            })
            i++
            continue
        }

        // Statement or default
        const blockTokens = tokensByLine[lineNum] || []
        blocks.push({
            id: blockId++,
            type: 'statement',
            startLine: lineNum, endLine: lineNum, code: line,
            tokens: blockTokens,
            tokenCount: blockTokens.length,
            wsCount: countWS(line)
        })
        i++
    }

    return blocks
}

function CodeBlockViewer({ tokens, code }: { tokens: TokenDTO[], code: string }) {
    const blocks = useMemo(() => parseCodeBlocks(code, tokens), [code, tokens])

    const typeLabels: Record<string, string> = {
        curly_brace: '{ Block }',
        statement: 'Statement',
        single_comment: '!> Comment',
        multi_comment: '<! Comment !>'
    }

    const typeColors: Record<string, string> = {
        curly_brace: '#7c3aed',
        statement: '#059669',
        single_comment: '#9ca3af',
        multi_comment: '#6b7280'
    }

    return (
        <div style={{ flex: 1, overflow: 'auto', padding: 'var(--spacing-lg)' }}>
            <div style={{ padding: '10px 0', color: 'var(--text)', fontWeight: 700, fontSize: '0.95rem', marginBottom: 'var(--spacing-md)' }}>
                Total Blocks: {blocks.length}
            </div>

            {blocks.map((block) => (
                <div key={block.id} style={{
                    marginBottom: 'var(--spacing-lg)', border: '2px solid var(--border)',
                    borderRadius: 'var(--radius-md)', overflow: 'hidden',
                    background: 'var(--panel-secondary)', transition: 'all var(--transition-base)'
                }}
                onMouseEnter={(e) => {
                    e.currentTarget.style.borderColor = 'var(--accent)'
                    e.currentTarget.style.boxShadow = 'var(--shadow-md)'
                }}
                onMouseLeave={(e) => {
                    e.currentTarget.style.borderColor = 'var(--border)'
                    e.currentTarget.style.boxShadow = 'none'
                }}>
                    <div style={{
                        display: 'flex', justifyContent: 'space-between', alignItems: 'center',
                        padding: 'var(--spacing-sm) var(--spacing-md)',
                        background: 'var(--table-header)', borderBottom: '1px solid var(--border)'
                    }}>
                        <div style={{ display: 'flex', gap: 'var(--spacing-md)', alignItems: 'center' }}>
                            <span style={{
                                background: typeColors[block.type], color: '#ffffff',
                                padding: '4px 10px', borderRadius: 'var(--radius-sm)',
                                fontSize: '0.7rem', fontWeight: 700, textTransform: 'uppercase', letterSpacing: '0.05em'
                            }}>
                                {typeLabels[block.type]}
                            </span>
                            <span style={{ fontSize: '0.8rem', color: 'var(--text-secondary)' }}>
                                Lines {block.startLine}{block.endLine !== block.startLine ? `-${block.endLine}` : ''}
                            </span>
                        </div>
                        <div style={{ display: 'flex', gap: 'var(--spacing-md)', fontSize: '0.75rem', color: 'var(--muted)' }}>
                            <span>Tokens: <strong style={{ color: 'var(--accent-bright)' }}>{block.tokenCount}</strong></span>
                            <span>WS: <strong style={{ color: 'var(--accent-bright)' }}>{block.wsCount}</strong></span>
                        </div>
                    </div>

                    <pre style={{
                        margin: 0, padding: 'var(--spacing-md)', background: 'var(--accent-light)',
                        color: 'var(--text)', fontSize: '0.9rem',
                        fontFamily: 'ui-monospace, "Cascadia Code", "Fira Code", monospace',
                        whiteSpace: 'pre-wrap', wordBreak: 'break-word',
                        borderBottom: block.tokens.length > 0 ? '1px solid var(--border)' : 'none'
                    }}>
                        {block.code}
                    </pre>

                    {block.tokens.length > 0 && (
                        <div style={{ padding: 'var(--spacing-sm) var(--spacing-md)' }}>
                            <div style={{ display: 'flex', flexWrap: 'wrap', gap: 'var(--spacing-sm)' }}>
                                <div className="token-item token-count">
                                    <span className="token-type">Token Count</span>
                                    <span className="token-lexeme">{block.tokenCount}</span>
                                </div>
                                <div className="token-item ws-count">
                                    <span className="token-type">WS Count</span>
                                    <span className="token-lexeme">{block.wsCount}</span>
                                </div>
                                {block.tokens.map((t, idx) => (
                                    <div key={idx} className="token-item" style={{ fontSize: '0.7rem' }}>
                                        <span className="token-type">{t.semanticGroup}</span>
                                        <span className="token-lexeme">"{t.lexeme}"</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    )}
                </div>
            ))}

            {blocks.length === 0 && (
                <div style={{ padding: '20px', color: 'var(--muted)', textAlign: 'center' }}>
                    No code blocks found.
                </div>
            )}
        </div>
    )
}

function ClassicTokenTable({ tokens, whitespaceCounts }: { tokens: TokenDTO[], whitespaceCounts: number[] }) {
    const visible = useMemo(() => filterNewlines(tokens), [tokens])
    const tokenCounts = useMemo(() => calculateTokenCounts(tokens), [tokens])
    const totalTokens = useMemo(() => visible.length, [visible])

    return (
        <div style={{ flex: 1, overflow: 'auto' }}>
            <div style={{ padding: '10px 12px', color: 'var(--text)', fontWeight: 700, fontSize: '0.95rem' }}>
                Total Tokens: {totalTokens}
            </div>
            <table>
                <thead>
                    <tr><th>Lexeme</th><th>Semantic Group</th><th>Line</th><th>WS Count</th></tr>
                </thead>
                <tbody>
                    {visible.map((t, i) => (
                        <tr key={i}>
                            <td>{t.lexeme}</td>
                            <td>{t.semanticGroup}</td>
                            <td>{t.line}</td>
                            <td>{whitespaceCounts[t.line - 1] ?? 0}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    )
}

function LineByLineViewer({ tokens, code, whitespaceCounts }: { tokens: TokenDTO[], code: string, whitespaceCounts: number[] }) {
    const lines = useMemo(() => code.split('\n'), [code])
    const tokensByLine = useMemo(() => 
        filterNewlines(tokens).reduce((acc, t) => {
            const lineNum = t.line
            if (!acc[lineNum]) acc[lineNum] = []
            acc[lineNum].push(t)
            return acc
        }, {} as Record<number, TokenDTO[]>), [tokens])

    return (
        <div className="line-by-line-viewer">
            {lines.map((line, lineIndex) => {
                const lineNumber = lineIndex + 1
                const lineTokens = tokensByLine[lineNumber] || []
                const wsCount = whitespaceCounts[lineIndex] ?? 0
                const isBlank = line.trim().length === 0 && lineTokens.length === 0
                if (isBlank) return null

                return (
                    <div key={lineIndex} className="line-row">
                        <div className="line-num">{lineNumber}</div>
                        <div className="line-content">
                            <span className="source-code">{line}</span>
                            <div className="tokens-list">
                                <div className="token-item token-count">
                                    <span className="token-type">Token Count</span>
                                    <span className="token-lexeme">{lineTokens.length}</span>
                                </div>
                                <div className="token-item ws-count">
                                    <span className="token-type">WS Count</span>
                                    <span className="token-lexeme">{wsCount}</span>
                                </div>
                                {lineTokens.map((t, i) => (
                                    <div key={i} className="token-item">
                                        <span className="token-type">{t.type}</span>
                                        <span className="token-lexeme">"{t.lexeme}"</span>
                                        <span className="token-col">Col: {t.column}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                )
            })}
        </div>
    )
}

function SingleLineViewer({ tokens, code, currentLine, onLineChange, whitespaceCounts }: Props & { whitespaceCounts: number[] }) {
    const lines = useMemo(() => code.split('\n'), [code])
    const currentLineCode = lines[currentLine - 1] || ""
    const currentLineTokens = useMemo(() => 
        filterNewlines(tokens).filter(t => t.line === currentLine), [tokens, currentLine])
    const wsCount = whitespaceCounts[currentLine - 1] ?? 0
    const tokenCountForLine = currentLineTokens.length

    return (
        <div className="single-line-viewer">
            <div className="single-line-header">
                <div className="line-label">LINE {currentLine}</div>
                <div className="nav-buttons">
                    <button onClick={() => onLineChange('down')} disabled={currentLine === 1}
                        className="nav-btn up" title="Previous Line">🔺</button>
                    <button onClick={() => onLineChange('up')} disabled={currentLine === lines.length}
                        className="nav-btn down" title="Next Line">🔻</button>
                </div>
            </div>
            <div className="line-code-preview">
                {currentLine} | {currentLineCode}
                <span style={{float: 'right', opacity: 0.7, fontSize: '0.8em'}}>WS: {wsCount}</span>
            </div>
            <table>
                <thead>
                    <tr><th>LEXEME</th><th>TOKEN</th><th>Token Count</th><th>WS COUNT</th></tr>
                </thead>
                <tbody>
                    {currentLineTokens.map((t, i) => (
                        <tr key={i}>
                            <td>{t.lexeme}</td>
                            <td>{t.type}</td>
                            <td>{tokenCountForLine}</td>
                            <td>{wsCount}</td>
                        </tr>
                    ))}
                    {currentLineTokens.length === 0 && (
                        <tr>
                            <td colSpan={4} style={{ textAlign: 'center', color: 'var(--muted)' }}>
                                No tokens found on Line {currentLine}. WS Count: {wsCount}
                            </td>
                        </tr>
                    )}
                </tbody>
            </table>
        </div>
    )
}

export default function OutputTable(props: Props) {
    const whitespaceCounts = useMemo(() => calculateWhitespaceCount(props.code), [props.code])

    return (
        <div style={{ height: '100%', display: 'flex', flexDirection: 'column' }}>
            {/* Result Tabs */}
            <div className="result-tabs">
                <button
                    className={`result-tab ${props.viewMode === 'table' ? 'active' : ''}`}
                    onClick={() => props.onViewModeChange('table')}
                >
                    Table
                </button>
                <button
                    className={`result-tab ${props.viewMode === 'singleLine' ? 'active' : ''}`}
                    onClick={() => props.onViewModeChange('singleLine')}
                >
                    Line
                </button>
                <button
                    className={`result-tab ${props.viewMode === 'lineByLine' ? 'active' : ''}`}
                    onClick={() => props.onViewModeChange('lineByLine')}
                >
                    All
                </button>
                <button
                    className={`result-tab ${props.viewMode === 'codeBlock' ? 'active' : ''}`}
                    onClick={() => props.onViewModeChange('codeBlock')}
                >
                    Blocks
                </button>
            </div>
            
            {/* Content */}
            <div className="outputContainer">
                {props.viewMode === 'table' && <ClassicTokenTable tokens={props.tokens} whitespaceCounts={whitespaceCounts} />}
                {props.viewMode === 'lineByLine' && <LineByLineViewer tokens={props.tokens} code={props.code} whitespaceCounts={whitespaceCounts} />}
                {props.viewMode === 'singleLine' && <SingleLineViewer {...props} whitespaceCounts={whitespaceCounts} />}
                {props.viewMode === 'codeBlock' && <CodeBlockViewer tokens={props.tokens} code={props.code} />}

                {props.tokens.length === 0 && props.code.trim().length === 0 && (
                    <div style={{ padding: '20px', color: 'var(--muted)', textAlign: 'center' }}>
                        Run the code analyzer to see the lexical output here.
                    </div>
                )}
            </div>
        </div>
    )
}