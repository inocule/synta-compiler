# Synta-Lexical Frontend

This is a minimal React + Vite frontend for the `synta-lexical` analyzer.

Quick start (PowerShell):

```powershell
cd D:\Programming\PPL\synta-lexical\frontend
npm install
npm run dev
```

The frontend expects the analyzer HTTP server at `http://localhost:8080/api/analyze`.
If you haven't added the server, run the Go analyzer locally and expose `/api/analyze`.

Editor: uses `@monaco-editor/react`. Run button sends the code to the analyzer and displays tokens.

Next steps:
- Add token highlighting in the editor using Monaco decorations.
- Add abortable requests and debounce for live analysis.
- Optionally compile the analyzer to WASM for in-browser analysis.
