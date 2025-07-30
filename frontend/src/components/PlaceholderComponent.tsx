import React from 'react';
import { TrendingUp, History, BarChart3 } from 'lucide-react';

// PriceChart Component
interface PriceChartProps {
  sourceToken: string;
  destToken: string;
}

export const PriceChart: React.FC<PriceChartProps> = ({ sourceToken, destToken }) => {
  return (
    <div className="h-48 flex items-center justify-center bg-white/5 rounded-lg">
      <div className="text-center">
        <TrendingUp className="w-8 h-8 text-blue-400 mx-auto mb-2" />
        <p className="text-gray-300 text-sm">
          {sourceToken}/{destToken} Price Chart
        </p>
        <p className="text-gray-500 text-xs mt-1">Coming Soon</p>
      </div>
    </div>
  );
};

// TransactionHistory Component
export const TransactionHistory: React.FC = () => {
  return (
    <div className="max-w-4xl mx-auto">
      <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 border border-white/20">
        <div className="text-center">
          <History className="w-12 h-12 text-blue-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-white mb-2">Transaction History</h2>
          <p className="text-gray-300 mb-6">
            View your cross-chain bridge transactions and TWAP executions
          </p>
          <div className="bg-white/5 rounded-lg p-6">
            <p className="text-gray-400">No transactions yet</p>
            <p className="text-sm text-gray-500 mt-2">
              Start bridging assets to see your transaction history here
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

// Analytics Component
export const Analytics: React.FC = () => {
  return (
    <div className="max-w-6xl mx-auto">
      <div className="bg-white/10 backdrop-blur-md rounded-2xl p-8 border border-white/20">
        <div className="text-center">
          <BarChart3 className="w-12 h-12 text-purple-400 mx-auto mb-4" />
          <h2 className="text-2xl font-bold text-white mb-2">Analytics Dashboard</h2>
          <p className="text-gray-300 mb-6">
            Monitor TWAP performance, price impact savings, and bridge statistics
          </p>
          
          <div className="grid md:grid-cols-3 gap-6 mt-8">
            <div className="bg-white/5 rounded-lg p-6">
              <h3 className="text-lg font-semibold text-white mb-2">TWAP Savings</h3>
              <div className="text-3xl font-bold text-green-400 mb-1">2.3%</div>
              <p className="text-sm text-gray-400">Average price impact reduction</p>
            </div>
            
            <div className="bg-white/5 rounded-lg p-6">
              <h3 className="text-lg font-semibold text-white mb-2">Total Volume</h3>
              <div className="text-3xl font-bold text-blue-400 mb-1">$1.2M</div>
              <p className="text-sm text-gray-400">24h bridge volume</p>
            </div>
            
            <div className="bg-white/5 rounded-lg p-6">
              <h3 className="text-lg font-semibold text-white mb-2">Success Rate</h3>
              <div className="text-3xl font-bold text-purple-400 mb-1">99.8%</div>
              <p className="text-sm text-gray-400">Transaction success rate</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};