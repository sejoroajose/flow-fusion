import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { Wallet, Activity, BarChart3, Github } from 'lucide-react';
import { useWallet } from '../context/WalletContext';

const Header: React.FC = () => {
  const { isConnected, address, connectWallet, disconnectWallet, balance, chainId } = useWallet();
  const location = useLocation();

  const navigation = [
    { name: 'Bridge', href: '/', icon: Activity },
    { name: 'History', href: '/history', icon: BarChart3 },
    { name: 'Analytics', href: '/analytics', icon: BarChart3 },
  ];

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`;
  };

  const getChainName = (chainId: number | null) => {
    switch (chainId) {
      case 1:
        return 'Ethereum';
      case 11155111:
        return 'Sepolia';
      default:
        return 'Unknown';
    }
  };

  return (
    <header className="bg-black/20 backdrop-blur-md border-b border-white/10 sticky top-0 z-40">
      <div className="container mx-auto px-4">
        <div className="flex items-center justify-between h-16">
          {/* Logo */}
          <Link to="/" className="flex items-center space-x-2">
            <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
              <span className="text-white font-bold text-sm">FF</span>
            </div>
            <span className="text-xl font-bold text-white">Flow Fusion</span>
          </Link>

          {/* Navigation */}
          <nav className="hidden md:flex items-center space-x-8">
            {navigation.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.href;
              
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`
                    flex items-center space-x-2 px-3 py-2 rounded-lg transition-all duration-200
                    ${isActive 
                      ? 'bg-white/10 text-white' 
                      : 'text-gray-300 hover:text-white hover:bg-white/5'
                    }
                  `}
                >
                  <Icon className="w-4 h-4" />
                  <span className="font-medium">{item.name}</span>
                </Link>
              );
            })}
          </nav>

          {/* Wallet Connection */}
          <div className="flex items-center space-x-4">
            {/* GitHub Link */}
            <a
              href="https://github.com/your-username/flow-fusion"
              target="_blank"
              rel="noopener noreferrer"
              className="text-gray-400 hover:text-white transition-colors"
            >
              <Github className="w-5 h-5" />
            </a>

            {/* Chain Indicator */}
            {isConnected && chainId && (
              <div className="hidden sm:flex items-center space-x-2 bg-white/10 rounded-lg px-3 py-1.5 border border-white/20">
                <div className="w-2 h-2 bg-green-400 rounded-full"></div>
                <span className="text-sm text-white font-medium">
                  {getChainName(chainId)}
                </span>
              </div>
            )}

            {/* Wallet Button */}
            {isConnected ? (
              <div className="flex items-center space-x-3">
                {/* Balance */}
                {balance && (
                  <div className="hidden sm:block text-right">
                    <div className="text-sm text-white font-medium">
                      {parseFloat(balance).toFixed(4)} ETH
                    </div>
                    <div className="text-xs text-gray-400">Balance</div>
                  </div>
                )}

                {/* Address Button */}
                <button
                  onClick={disconnectWallet}
                  className="flex items-center space-x-2 bg-white/10 hover:bg-white/20 rounded-lg px-4 py-2 border border-white/20 transition-all duration-200"
                >
                  <Wallet className="w-4 h-4 text-white" />
                  <span className="text-white font-medium">
                    {formatAddress(address!)}
                  </span>
                </button>
              </div>
            ) : (
              <button
                onClick={connectWallet}
                className="flex items-center space-x-2 bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 rounded-lg px-4 py-2 transition-all duration-200"
              >
                <Wallet className="w-4 h-4 text-white" />
                <span className="text-white font-medium">Connect Wallet</span>
              </button>
            )}
          </div>
        </div>

        {/* Mobile Navigation */}
        <div className="md:hidden py-4 border-t border-white/10">
          <div className="flex justify-center space-x-6">
            {navigation.map((item) => {
              const Icon = item.icon;
              const isActive = location.pathname === item.href;
              
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`
                    flex flex-col items-center space-y-1 p-2 rounded-lg transition-all duration-200
                    ${isActive 
                      ? 'text-white' 
                      : 'text-gray-400 hover:text-white'
                    }
                  `}
                >
                  <Icon className="w-5 h-5" />
                  <span className="text-xs font-medium">{item.name}</span>
                </Link>
              );
            })}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;