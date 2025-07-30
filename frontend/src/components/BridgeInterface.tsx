import React, { useState, useEffect } from 'react';
import { ArrowUpDown, Settings, TrendingUp, Zap, Shield } from 'lucide-react';
import { useWallet } from '../context/WalletContext';
import { useBridge } from '../hooks/useBridge';
import { useToast } from '../context/ToastContext';
import TWAPSettings from './TWAPSettings';
import TransactionPreview from './TransactionPreview';
import { PriceChart } from './PlaceholderComponent';

interface BridgeState {
  sourceChain: string;
  destChain: string;
  sourceToken: string;
  destToken: string;
  amount: string;
  useTWAP: boolean;
  twapConfig: {
    timeWindow: number;
    intervalCount: number;
    maxSlippage: number;
  };
}

const BridgeInterface: React.FC = () => {
  const { isConnected, address, connectWallet } = useWallet();
  const { createOrder, getQuote, getTWAPQuote } = useBridge();
  const { showToast } = useToast();

  const [bridgeState, setBridgeState] = useState<BridgeState>({
    sourceChain: 'ethereum',
    destChain: 'cosmos',
    sourceToken: 'USDC',
    destToken: 'ATOM',
    amount: '',
    useTWAP: false,
    twapConfig: {
      timeWindow: 3600, // 1 hour
      intervalCount: 12, // 5-minute intervals
      maxSlippage: 0.5,  // 0.5%
    },
  });

  const [quote, setQuote] = useState<any>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [showAdvanced, setShowAdvanced] = useState(false);

  // Get quote when parameters change
  useEffect(() => {
    if (bridgeState.amount && parseFloat(bridgeState.amount) > 0) {
      handleGetQuote();
    }
  }, [bridgeState.amount, bridgeState.sourceToken, bridgeState.destToken, bridgeState.useTWAP]);

  const handleGetQuote = async () => {
    if (!bridgeState.amount || parseFloat(bridgeState.amount) <= 0) return;

    setIsLoading(true);
    try {
      let quoteData;
      if (bridgeState.useTWAP) {
        quoteData = await getTWAPQuote({
          sourceChain: bridgeState.sourceChain,
          destChain: bridgeState.destChain,
          sourceToken: bridgeState.sourceToken,
          destToken: bridgeState.destToken,
          totalAmount: bridgeState.amount,
          timeWindow: bridgeState.twapConfig.timeWindow,
          intervalCount: bridgeState.twapConfig.intervalCount,
        });
      } else {
        quoteData = await getQuote({
          sourceChain: bridgeState.sourceChain,
          destChain: bridgeState.destChain,
          sourceToken: bridgeState.sourceToken,
          destToken: bridgeState.destToken,
          amount: bridgeState.amount,
        });
      }
      setQuote(quoteData);
    } catch (error) {
      showToast('Failed to get quote', 'error');
    }
    setIsLoading(false);
  };

  const handleSwapChains = () => {
    setBridgeState(prev => ({
      ...prev,
      sourceChain: prev.destChain,
      destChain: prev.sourceChain,
      sourceToken: prev.destToken,
      destToken: prev.sourceToken,
    }));
  };

  const handleCreateOrder = async () => {
    if (!isConnected) {
      await connectWallet();
      return;
    }

    if (!bridgeState.amount || parseFloat(bridgeState.amount) <= 0) {
      showToast('Please enter a valid amount', 'error');
      return;
    }

    setIsLoading(true);
    try {
      const orderData = {
        userAddress: address!,
        sourceChain: bridgeState.sourceChain,
        destChain: bridgeState.destChain,
        sourceToken: bridgeState.sourceToken,
        destToken: bridgeState.destToken,
        amount: bridgeState.amount,
        useTWAP: bridgeState.useTWAP,
        twapConfig: bridgeState.useTWAP ? bridgeState.twapConfig : undefined,
      };

      const order = await createOrder(orderData);
      showToast(`Order created successfully! ID: ${order.id}`, 'success');
      
      // Reset form
      setBridgeState(prev => ({ ...prev, amount: '' }));
      setQuote(null);
    } catch (error) {
      showToast('Failed to create order', 'error');
    }
    setIsLoading(false);
  };

  return (
    <div className="max-w-4xl mx-auto space-y-8">
      {/* Hero Section */}
      <div className="text-center space-y-4">
        <h1 className="text-4xl font-bold text-white bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
          Flow Fusion
        </h1>
        <p className="text-xl text-gray-300">
          TWAP-Enabled Cross-Chain Bridge: Ethereum ↔ Cosmos
        </p>
        <div className="flex justify-center space-x-6 text-sm text-gray-400">
          <div className="flex items-center space-x-2">
            <Zap className="w-4 h-4" />
            <span>Fast & Secure</span>
          </div>
          <div className="flex items-center space-x-2">
            <TrendingUp className="w-4 h-4" />
            <span>TWAP Optimization</span>
          </div>
          <div className="flex items-center space-x-2">
            <Shield className="w-4 h-4" />
            <span>Atomic Swaps</span>
          </div>
        </div>
      </div>

      <div className="grid lg:grid-cols-3 gap-8">
        {/* Main Bridge Interface */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20">
            <div className="flex justify-between items-center mb-6">
              <h2 className="text-2xl font-semibold text-white">Bridge Assets</h2>
              <button
                onClick={() => setShowAdvanced(!showAdvanced)}
                className="flex items-center space-x-2 text-blue-400 hover:text-blue-300 transition-colors"
              >
                <Settings className="w-4 h-4" />
                <span>Advanced</span>
              </button>
            </div>

            {/* From Section */}
            <div className="space-y-4">
              <div className="bg-white/5 rounded-xl p-4 border border-white/10">
                <div className="flex justify-between items-center mb-3">
                  <label className="text-sm font-medium text-gray-300">From</label>
                  <select
                    value={bridgeState.sourceChain}
                    onChange={(e) => setBridgeState(prev => ({ ...prev, sourceChain: e.target.value }))}
                    className="bg-white/10 text-white rounded-lg px-3 py-1 border border-white/20 text-sm"
                  >
                    <option value="ethereum">Ethereum</option>
                    <option value="cosmos">Cosmos</option>
                  </select>
                </div>
                <div className="flex space-x-3">
                  <input
                    type="number"
                    value={bridgeState.amount}
                    onChange={(e) => setBridgeState(prev => ({ ...prev, amount: e.target.value }))}
                    placeholder="0.0"
                    className="flex-1 bg-transparent text-white text-2xl font-semibold placeholder-gray-500 outline-none"
                  />
                  <select
                    value={bridgeState.sourceToken}
                    onChange={(e) => setBridgeState(prev => ({ ...prev, sourceToken: e.target.value }))}
                    className="bg-white/10 text-white rounded-lg px-3 py-2 border border-white/20"
                  >
                    <option value="USDC">USDC</option>
                    <option value="ETH">ETH</option>
                    <option value="WETH">WETH</option>
                  </select>
                </div>
                <div className="flex justify-between text-sm text-gray-400 mt-2">
                  <span>Balance: 0.00</span>
                  <button className="text-blue-400 hover:text-blue-300">MAX</button>
                </div>
              </div>

              {/* Swap Button */}
              <div className="flex justify-center">
                <button
                  onClick={handleSwapChains}
                  className="bg-white/10 hover:bg-white/20 rounded-full p-3 border border-white/20 transition-all duration-200 hover:scale-105"
                >
                  <ArrowUpDown className="w-5 h-5 text-white" />
                </button>
              </div>

              {/* To Section */}
              <div className="bg-white/5 rounded-xl p-4 border border-white/10">
                <div className="flex justify-between items-center mb-3">
                  <label className="text-sm font-medium text-gray-300">To</label>
                  <select
                    value={bridgeState.destChain}
                    onChange={(e) => setBridgeState(prev => ({ ...prev, destChain: e.target.value }))}
                    className="bg-white/10 text-white rounded-lg px-3 py-1 border border-white/20 text-sm"
                  >
                    <option value="cosmos">Cosmos</option>
                    <option value="ethereum">Ethereum</option>
                  </select>
                </div>
                <div className="flex space-x-3">
                  <input
                    type="text"
                    value={quote?.estimated_output || '0.0'}
                    readOnly
                    placeholder="0.0"
                    className="flex-1 bg-transparent text-white text-2xl font-semibold placeholder-gray-500 outline-none"
                  />
                  <select
                    value={bridgeState.destToken}
                    onChange={(e) => setBridgeState(prev => ({ ...prev, destToken: e.target.value }))}
                    className="bg-white/10 text-white rounded-lg px-3 py-2 border border-white/20"
                  >
                    <option value="ATOM">ATOM</option>
                    <option value="USDC">USDC</option>
                    <option value="OSMO">OSMO</option>
                  </select>
                </div>
                <div className="flex justify-between text-sm text-gray-400 mt-2">
                  <span>Balance: 0.00</span>
                  {quote && (
                    <span className="text-green-400">
                      Price Impact: {(quote.price_impact * 100).toFixed(2)}%
                    </span>
                  )}
                </div>
              </div>
            </div>

            {/* TWAP Settings */}
            <div className="mt-6">
              <div className="flex items-center space-x-3 mb-4">
                <input
                  type="checkbox"
                  id="useTWAP"
                  checked={bridgeState.useTWAP}
                  onChange={(e) => setBridgeState(prev => ({ ...prev, useTWAP: e.target.checked }))}
                  className="w-4 h-4 text-blue-500 rounded"
                />
                <label htmlFor="useTWAP" className="text-white font-medium">
                  Enable TWAP Execution
                </label>
                <div className="flex items-center space-x-1 text-xs text-green-400">
                  <TrendingUp className="w-3 h-3" />
                  <span>Reduces Price Impact</span>
                </div>
              </div>

              {bridgeState.useTWAP && (
                <TWAPSettings
                  config={bridgeState.twapConfig}
                  onChange={(config) => setBridgeState(prev => ({ ...prev, twapConfig: config }))}
                  quote={quote}
                />
              )}
            </div>

            {/* Submit Button */}
            <button
              onClick={handleCreateOrder}
              disabled={isLoading || !bridgeState.amount}
              className="w-full bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 disabled:from-gray-600 disabled:to-gray-600 text-white font-semibold py-4 rounded-xl transition-all duration-200 disabled:cursor-not-allowed mt-6"
            >
              {isLoading ? (
                <div className="flex items-center justify-center space-x-2">
                  <div className="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></div>
                  <span>Processing...</span>
                </div>
              ) : !isConnected ? (
                'Connect Wallet'
              ) : (
                bridgeState.useTWAP ? 'Create TWAP Order' : 'Bridge Assets'
              )}
            </button>
          </div>

          {/* Transaction Preview */}
          {quote && (
            <TransactionPreview
              quote={quote}
              useTWAP={bridgeState.useTWAP}
              twapConfig={bridgeState.twapConfig}
            />
          )}
        </div>

        {/* Side Panel */}
        <div className="space-y-6">
          {/* Price Chart */}
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20">
            <h3 className="text-lg font-semibold text-white mb-4">Price Chart</h3>
            <PriceChart 
              sourceToken={bridgeState.sourceToken}
              destToken={bridgeState.destToken}
            />
          </div>

          {/* TWAP Benefits */}
          {bridgeState.useTWAP && quote && (
            <div className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20">
              <h3 className="text-lg font-semibold text-white mb-4 flex items-center space-x-2">
                <TrendingUp className="w-5 h-5 text-green-400" />
                <span>TWAP Benefits</span>
              </h3>
              <div className="space-y-3">
                <div className="flex justify-between text-sm">
                  <span className="text-gray-300">Instant Swap Impact:</span>
                  <span className="text-red-400">{(quote.impact_comparison?.instant_swap * 100).toFixed(2)}%</span>
                </div>
                <div className="flex justify-between text-sm">
                  <span className="text-gray-300">TWAP Impact:</span>
                  <span className="text-green-400">{(quote.impact_comparison?.twap_execution * 100).toFixed(2)}%</span>
                </div>
                <div className="flex justify-between text-sm font-semibold">
                  <span className="text-gray-300">Improvement:</span>
                  <span className="text-green-400">
                    +{quote.impact_comparison?.improvement_percent.toFixed(1)}%
                  </span>
                </div>
              </div>
            </div>
          )}

          {/* Recent Activity */}
          <div className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20">
            <h3 className="text-lg font-semibold text-white mb-4">Recent Activity</h3>
            <div className="space-y-3 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-300">24h Volume:</span>
                <span className="text-white">$1.2M</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">Total Bridges:</span>
                <span className="text-white">1,234</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-300">TWAP Orders:</span>
                <span className="text-green-400">456</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default BridgeInterface;