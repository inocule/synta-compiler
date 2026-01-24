// App.tsx

import React, { useState, useEffect } from 'react'
import HomePage from './HomePage'
import LexicalAnalyzer from './LexicalAnalyzer'
import SyntacticalAnalyzer from './SyntacticalAnalyzer'

type AnalyzerType = 'lexical' | 'syntactical' | null

function App() {
  const [selectedAnalyzer, setSelectedAnalyzer] = useState<AnalyzerType>(null)
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

  // Handle analyzer selection
  const handleSelectAnalyzer = (type: 'lexical' | 'syntactical') => {
    setSelectedAnalyzer(type)
  }

  // Handle back to home
  const handleBackToHome = () => {
    setSelectedAnalyzer(null)
  }

  // Show HomePage if no analyzer is selected
  if (!selectedAnalyzer) {
    return <HomePage onSelectAnalyzer={handleSelectAnalyzer} />
  }

  // Render the selected analyzer with shared UI components
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

      {/* Back to Home Button */}
      <button 
        className="back-home-btn" 
        onClick={handleBackToHome}
        aria-label="Back to home"
        title="Back to home"
      >
        ← HOME
      </button>

      {/* Analyzer Type Badge */}
      <div className="analyzer-badge">
        {selectedAnalyzer === 'lexical' ? '📊 Lexical Analyzer' : '🌲 Syntactical Analyzer'}
      </div>

      {/* Render Selected Analyzer */}
      {selectedAnalyzer === 'lexical' && <LexicalAnalyzer theme={theme} />}
      {selectedAnalyzer === 'syntactical' && <SyntacticalAnalyzer theme={theme} />}
    </>
  )
}

export default App