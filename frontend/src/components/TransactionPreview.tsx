import React from 'react';
import { Clock, TrendingUp, DollarSign, Zap } from 'lucide-react';

interface TransactionPreviewProps {
  quote: any;
  useTWAP: boolean;
  twapConfig: {
    timeWindow: number;
    intervalCount: number;
    maxSlippage: number;
  };
}

const TransactionPreview: React.FC<TransactionPreviewProps> = ({ quote, useTWAP, twapConfig }) => {
  if (!quote) return null;

  return (
    <div className="bg-white/10 backdrop-blur-md rounded-2xl p-6 border border-white/20">
      <h3 className="text-lg font-semibold text-white mb-4">Transaction Preview</h3>
      
      <div className="space-y-4">
        {/* Basic Quote Info */}
        <div className="grid grid-cols-2 gap-4">
          <div className="bg-white/5 rounded-lg p-3">
            <div className="flex items-center space-x-2 mb-1">
              <DollarSign className="w-4 h-4 text-blue-400" />
              <span className="text-sm text-gray-300">You'll Receive</span>
            </div>
            <div className="text-lg font-semibold text-white">
              {(parseFloat(quote.estimated_output || '0') / 1000000).toFixed(6)} ATOM
            </div>
          </div>
          
          <div className="bg-white/5 rounded-lg p-3">
            <div className="flex items-center space-x-2 mb-1">
              <TrendingUp className="w-4 h-4 text-green-400" />
              <span className="text-sm text-gray-300">Price Impact</span>
            </div>
            <div className="text-lg font-semibold text-green-400">
              {(quote.price_impact * 100).toFixed(2)}%
            </div>
          </div>
        </div>

        {/* TWAP Specific Info */}
        {useTWAP && quote.execution_schedule && (
          <div className="bg-blue-500/10 border border-blue-500/20 rounded-lg p-4">
            <div className="flex items-center space-x-2 mb-3">
              <Clock className="w-4 h-4 text-blue-400" />
              <span className="text-sm font-medium text-blue-400">TWAP Execution Plan</span>
            </div>
            
            <div className="grid grid-cols-3 gap-4 text-sm">
              <div>
                <div className="text-gray-300">Total Time</div>
                <div className="text-white font-medium">
                  {Math.round(twapConfig.timeWindow / 60)} minutes
                </div>
              </div>
              <div>
                <div className="text-gray-300">Intervals</div>
                <div className="text-white font-medium">
                  {twapConfig.intervalCount}
                </div>
              </div>
              <div>
                <div className="text-gray-300">Per Interval</div>
                <div className="text-white font-medium">
                  {Math.round(twapConfig.timeWindow / twapConfig.intervalCount / 60)}min
                </div>
              </div>
            </div>

            {quote.impact_comparison && (
              <div className="mt-3 pt-3 border-t border-blue-500/20">
                <div className="text-xs text-gray-300 mb-1">Price Impact Comparison</div>
                <div className="flex justify-between text-sm">
                  <span>Instant: <span className="text-red-400">{(quote.impact_comparison.instant_swap * 100).toFixed(2)}%</span></span>
                  <span>TWAP: <span className="text-green-400">{(quote.impact_comparison.twap_execution * 100).toFixed(2)}%</span></span>
                  <span>Saved: <span className="text-blue-400">{quote.impact_comparison.improvement_percent.toFixed(1)}%</span></span>
                </div>
              </div>
            )}
          </div>
        )}

        {/* Fee Breakdown */}
        {quote.fees && (
          <div className="bg-white/5 rounded-lg p-4">
            <h4 className="text-sm font-medium text-gray-300 mb-3">Fee Breakdown</h4>
            <div className="space-y-2 text-sm">
              <div className="flex justify-between">
                <span className="text-gray-400">Ethereum Gas</span>
                <span className="text-white">{quote.fees.ethereum_gas} ETH</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">Cosmos Gas</span>
                <span className="text-white">{quote.fees.cosmos_gas} ATOM</span>
              </div>
              <div className="flex justify-between">
                <span className="text-gray-400">Relayer Fee</span>
                <span className="text-white">{quote.fees.relayer_fee} USDC</span>
              </div>
              <div className="flex justify-between border-t border-white/10 pt-2 font-medium">
                <span className="text-gray-300">Total Fees</span>
                <span className="text-white">${quote.fees.total}</span>
              </div>
            </div>
          </div>
        )}

        {/* Execution Time */}
        <div className="flex items-center justify-between bg-white/5 rounded-lg p-3">
          <div className="flex items-center space-x-2">
            <Zap className="w-4 h-4 text-yellow-400" />
            <span className="text-sm text-gray-300">Estimated Time</span>
          </div>
          <span className="text-white font-medium">
            {useTWAP ? `${Math.round(twapConfig.timeWindow / 60)} minutes` : '~5 minutes'}
          </span>
        </div>

        {/* Recommendations */}
        {quote.recommendations && quote.recommendations.length > 0 && (
          <div className="bg-purple-500/10 border border-purple-500/20 rounded-lg p-4">
            <h4 className="text-sm font-medium text-purple-400 mb-2">Recommendations</h4>
            <div className="space-y-1">
              {quote.recommendations.slice(0, 2).map((rec: string, index: number) => (
                <div key={index} className="text-xs text-gray-300 bg-white/5 rounded px-2 py-1">
                  {rec}
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default TransactionPreview;