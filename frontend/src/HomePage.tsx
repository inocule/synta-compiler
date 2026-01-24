// HomePage.tsx

import React, { useState, useEffect } from 'react'
import './HomePage.css'
import syntaLogo from './components/title.png'

interface HomePageProps {
  onSelectAnalyzer: (type: 'lexical' | 'syntactical') => void
}

const HomePage: React.FC<HomePageProps> = ({ onSelectAnalyzer }) => {
  const [hoveredCard, setHoveredCard] = useState<'lexical' | 'syntactical' | null>(null)
  const [mounted, setMounted] = useState(false)

  useEffect(() => {
    setMounted(true)
  }, [])

  return (
    <div className={`homepage ${mounted ? 'mounted' : ''}`}>
      {/* Animated background gradient */}
      <div className="bg-gradient-overlay"></div>

      <div className="homepage-content">
        {/* Header with Logo */}
        <header className="homepage-header">
          <img 
            src={syntaLogo} 
            alt="SYNTA Logo" 
            className="synta-logo"
          />
          <p className="subtitle">Choose your analyzer and start exploring code syntax</p>
        </header>

        {/* Analyzer Cards */}
        <div className="analyzer-cards">
          {/* Lexical Analyzer Card */}
          <div 
            className={`analyzer-card lexical-card ${hoveredCard === 'lexical' ? 'hovered' : ''} ${hoveredCard && hoveredCard !== 'lexical' ? 'dimmed' : ''}`}
            onMouseEnter={() => setHoveredCard('lexical')}
            onMouseLeave={() => setHoveredCard(null)}
            onClick={() => onSelectAnalyzer('lexical')}
          >
            <div className="card-glow"></div>
            <div className="card-content">
              <div className="card-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M4 7h16M4 12h16M4 17h16" />
                  <circle cx="8" cy="7" r="1.5" fill="currentColor" />
                  <circle cx="8" cy="12" r="1.5" fill="currentColor" />
                  <circle cx="8" cy="17" r="1.5" fill="currentColor" />
                </svg>
              </div>
              <h2 className="card-title">Lexical Analyzer</h2>
              <p className="card-description">
                Break down source code into tokens and understand its fundamental building blocks.
              </p>
              <ul className="card-features">
                <li>
                  <span className="feature-icon">→</span>
                  <span>Token identification</span>
                </li>
                <li>
                  <span className="feature-icon">→</span>
                  <span>Line-by-line analysis</span>
                </li>
                <li>
                  <span className="feature-icon">→</span>
                  <span>Syntax highlighting</span>
                </li>
              </ul>
              <div className="card-cta">
                <span className="cta-text">Start Tokenizing</span>
                <span className="cta-arrow">→</span>
              </div>
            </div>
            <div className="card-accent-line"></div>
          </div>

          {/* Syntactical Analyzer Card */}
          <div 
            className={`analyzer-card syntactical-card ${hoveredCard === 'syntactical' ? 'hovered' : ''} ${hoveredCard && hoveredCard !== 'syntactical' ? 'dimmed' : ''}`}
            onMouseEnter={() => setHoveredCard('syntactical')}
            onMouseLeave={() => setHoveredCard(null)}
            onClick={() => onSelectAnalyzer('syntactical')}
          >
            <div className="card-glow"></div>
            <div className="card-content">
              <div className="card-icon">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M12 2L2 7l10 5 10-5-10-5z" />
                  <path d="M2 17l10 5 10-5" />
                  <path d="M2 12l10 5 10-5" />
                </svg>
              </div>
              <h2 className="card-title">Syntactical Analyzer</h2>
              <p className="card-description">
                Parse and validate code structure to detect syntax errors and ensure grammatical correctness.
              </p>
              <ul className="card-features">
                <li>
                  <span className="feature-icon">→</span>
                  <span>Grammar validation</span>
                </li>
                <li>
                  <span className="feature-icon">→</span>
                  <span>Error location tracking</span>
                </li>
                <li>
                  <span className="feature-icon">→</span>
                  <span>Syntax error detection</span>
                </li>
              </ul>
              <div className="card-cta">
                <span className="cta-text">Start Parsing</span>
                <span className="cta-arrow">→</span>
              </div>
            </div>
            <div className="card-accent-line"></div>
          </div>
        </div>

        {/* Footer Info */}
        <footer className="homepage-footer">
          <p className="footer-text">
            Built for developers and students exploring compiler design
          </p>
        </footer>
      </div>
    </div>
  )
}

export default HomePage