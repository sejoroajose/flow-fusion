import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import { ToastProvider } from './context/ToastContext';
import { WalletProvider } from './context/WalletContext';
import Header from './components/Header';
import BridgeInterface from './components/BridgeInterface';
import { TransactionHistory, Analytics } from './components/PlaceholderComponent';
import './App.css';

function App() {
  return (
    <WalletProvider>
      <ToastProvider>
        <Router>
          <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
            <Header />
            
            <main className="container mx-auto px-4 py-8">
              <Routes>
                <Route path="/" element={<BridgeInterface />} />
                <Route path="/history" element={<TransactionHistory />} />
                <Route path="/analytics" element={<Analytics />} />
              </Routes>
            </main>

            {/* Background Effects */}
            <div className="fixed inset-0 pointer-events-none">
              <div className="absolute top-0 left-1/4 w-72 h-72 bg-purple-500/10 rounded-full blur-3xl animate-pulse-slow"></div>
              <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-blue-500/10 rounded-full blur-3xl animate-pulse-slow delay-1000"></div>
            </div>
          </div>
        </Router>
      </ToastProvider>
    </WalletProvider>
  );
}

export default App;