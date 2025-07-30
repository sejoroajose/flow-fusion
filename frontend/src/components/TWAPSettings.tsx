import React from 'react';
import { Clock, Sliders, TrendingDown, Info } from 'lucide-react';

interface TWAPConfig {
  timeWindow: number;
  intervalCount: number;
  maxSlippage: number;
}

interface TWAPSettingsProps {
  config: TWAPConfig;
  onChange: (config: TWAPConfig) => void;
  quote?: any;
}

const TWAPSettings: React.FC<TWAPSettingsProps> = ({ config, onChange, quote }) => {
  const intervalDuration = config.timeWindow / config.intervalCount;
  
  const handleTimeWindowChange = (value: number) => {
    onChange({ ...config, timeWindow: value });
  };

  const handleIntervalCountChange = (value: number) => {
    onChange({ ...config, intervalCount: Math.max(2, Math.min(20, value)) });
  };

  const handleSlippageChange = (value: number) => {
    onChange({ ...config, maxSlippage: Math.max(0.1, Math.min(5.0, value)) });
  };

  return (
    <div className="bg-white/5 rounded-xl p-4 border border-white/10 space-y-4">
      <div className="flex items-center space-x-2 text-sm font-medium text-gray-300">
        <Clock className="w-4 h-4" />
        <span>TWAP Configuration</span>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {/* Time Window */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-300">
            Time Window
          </label>
          <select
            value={config.timeWindow}
            onChange={(e) => handleTimeWindowChange(Number(e.target.value))}
            className="w-full bg-white/10 text-white rounded-lg px-3 py-2 border border-white/20 text-sm"
          >
            <option value={1800}>30 minutes</option>
            <option value={3600}>1 hour</option>
            <option value={7200}>2 hours</option>
            <option value={14400}>4 hours</option>
            <option value={28800}>8 hours</option>
          </select>
          <p className="text-xs text-gray-400">
            Total execution time
          </p>
        </div>

        {/* Interval Count */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-300">
            Intervals
          </label>
          <div className="flex items-center space-x-2">
            <button
              onClick={() => handleIntervalCountChange(config.intervalCount - 1)}
              className="bg-white/10 hover:bg-white/20 text-white rounded px-2 py-1 text-sm"
            >
              -
            </button>
            <input
              type="number"
              min="2"
              max="20"
              value={config.intervalCount}
              onChange={(e) => handleIntervalCountChange(Number(e.target.value))}
              className="flex-1 bg-white/10 text-white rounded px-2 py-1 text-center text-sm border border-white/20"
            />
            <button
              onClick={() => handleIntervalCountChange(config.intervalCount + 1)}
              className="bg-white/10 hover:bg-white/20 text-white rounded px-2 py-1 text-sm"
            >
              +
            </button>
          </div>
          <p className="text-xs text-gray-400">
            {Math.round(intervalDuration / 60)}min per interval
          </p>
        </div>

        {/* Max Slippage */}
        <div className="space-y-2">
          <label className="block text-sm font-medium text-gray-300">
            Max Slippage
          </label>
          <div className="flex items-center space-x-2">
            <input
              type="range"
              min="0.1"
              max="5.0"
              step="0.1"
              value={config.maxSlippage}
              onChange={(e) => handleSlippageChange(Number(e.target.value))}
              className="flex-1"
            />
            <span className="text-white text-sm font-medium w-12">
              {config.maxSlippage.toFixed(1)}%
            </span>
          </div>
          <p className="text-xs text-gray-400">
            Maximum slippage per interval
          </p>
        </div>
      </div>

      {/* Execution Schedule Preview */}
      {quote?.execution_schedule && (
        <div className="space-y-3">
          <div className="flex items-center space-x-2 text-sm font-medium text-gray-300">
            <Sliders className="w-4 h-4" />
            <span>Execution Schedule</span>
          </div>
          
          <div className="bg-white/5 rounded-lg p-3 max-h-32 overflow-y-auto">
            <div className="space-y-2">
              {quote.execution_schedule.slice(0, 5).map((schedule: any, index: number) => (
                <div key={index} className="flex justify-between items-center text-xs">
                  <div className="flex items-center space-x-2">
                    <div className="w-1.5 h-1.5 bg-blue-400 rounded-full"></div>
                    <span className="text-gray-300">
                      Interval {schedule.interval_index + 1}
                    </span>
                  </div>
                  <div className="text-right">
                    <div className="text-white">
                      {(parseFloat(schedule.amount) / 1000000).toFixed(2)} USDC
                    </div>
                    <div className="text-gray-400">
                      {new Date(schedule.execution_time).toLocaleTimeString()}
                    </div>
                  </div>
                </div>
              ))}
              {quote.execution_schedule.length > 5 && (
                <div className="text-center text-xs text-gray-400">
                  +{quote.execution_schedule.length - 5} more intervals
                </div>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Benefits Summary */}
      {quote?.impact_comparison && (
        <div className="bg-green-500/10 border border-green-500/20 rounded-lg p-3">
          <div className="flex items-start space-x-2">
            <TrendingDown className="w-4 h-4 text-green-400 mt-0.5" />
            <div className="text-sm">
              <div className="text-green-400 font-medium">TWAP Benefits</div>
              <div className="text-gray-300 mt-1">
                Reduces price impact by{' '}
                <span className="text-green-400 font-semibold">
                  {quote.impact_comparison.improvement_percent.toFixed(1)}%
                </span>
                {' '}compared to instant swap
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Recommendations */}
      {quote?.recommendations && (
        <div className="space-y-2">
          <div className="flex items-center space-x-2 text-sm font-medium text-gray-300">
            <Info className="w-4 h-4" />
            <span>Recommendations</span>
          </div>
          <div className="space-y-1">
            {quote.recommendations.slice(0, 2).map((rec: string, index: number) => (
              <div key={index} className="text-xs text-blue-400 bg-blue-500/10 rounded px-2 py-1">
                {rec}
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default TWAPSettings;