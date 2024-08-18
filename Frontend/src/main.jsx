import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import { AuthProvider } from './Providers/AuthProvider.jsx';
import { Routes } from './Routers/routes.jsx';

export const API_ADRESS = "http://127.0.0.1:8080" // Enter here your backend url (without / at the end)

ReactDOM.createRoot(document.getElementById('root')).render( // Render all pages
  <React.StrictMode>
    <AuthProvider>
      <Routes />
    </AuthProvider>
  </React.StrictMode>,
)
