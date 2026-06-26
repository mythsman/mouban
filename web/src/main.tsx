import React from 'react'
import ReactDOM from 'react-dom/client'
import { BrowserRouter } from 'react-router-dom'
import { App } from 'antd'
import 'antd/dist/reset.css'
import './styles.css'
import RootApp from './App'

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <App>
        <RootApp />
      </App>
    </BrowserRouter>
  </React.StrictMode>,
)
