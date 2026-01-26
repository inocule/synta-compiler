// App.tsx

import React, { useState, useEffect } from 'react'
import LexicalAnalyzer from './LexicalAnalyzer'
import SyntacticalAnalyzer from './SyntacticalAnalyzer'
import './App.css'

type AnalyzerType = 'lexical' | 'syntactical'

function App() {
  const [activeTab, setActiveTab] = useState<AnalyzerType>('syntactical')
  const [theme, setTheme] = useState<'light' | 'dark'>('light')

  // Force light theme on mount 
  useEffect(() => {
    setTheme('light')
    document.body.setAttribute('data-theme', 'light')
    localStorage.setItem('theme', 'light')
  }, [])

  // Toggle theme function 
  const toggleTheme = () => {
    const newTheme = theme === 'light' ? 'dark' : 'light'
    setTheme(newTheme)
    document.body.setAttribute('data-theme', newTheme)
    localStorage.setItem('theme', newTheme)
  }

  return (
    <div className="app-container">
      {/* Header with Logo and Theme Toggle */}
      <header className="app-header">
        <div className="logo"></div>
        <button 
          className="theme-toggle" 
          onClick={toggleTheme}
          aria-label="Toggle theme"
        >
          <div className="theme-toggle-slider">
            {theme === 'light' ? '☼' : '☾'}
          </div>
        </button>
      </header>

      {/* Tabs */}
      <div className="tabs">
        <button 
          className={`tab ${activeTab === 'syntactical' ? 'active' : ''}`}
          onClick={() => setActiveTab('syntactical')}
        >
          Syntax
        </button>
        <button 
          className={`tab ${activeTab === 'lexical' ? 'active' : ''}`}
          onClick={() => setActiveTab('lexical')}
        >
          Lexical
        </button>
      </div>

      {/* Render Active Analyzer */}
      <div className="analyzer-content">
        {activeTab === 'lexical' && <LexicalAnalyzer theme={theme} />}
        {activeTab === 'syntactical' && <SyntacticalAnalyzer theme={theme} />}
      </div>
    </div>
  )
}

export default App