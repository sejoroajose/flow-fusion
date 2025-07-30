import { useState, useEffect } from 'react';
import axios from 'axios';

const API_BASE_URL = import.meta.env.REACT_APP_RELAYER_API || 'http://localhost:8080/api/v1';

interface BridgeQuoteRequest {
  sourceChain: string;
  destChain: string;
  sourceToken: string;
  destToken: string;
  amount: string;
}

interface TWAPQuoteRequest {
  sourceChain: string;
  destChain: string;
  sourceToken: string;
  destToken: string;
  totalAmount: string;
  timeWindow: number;
  intervalCount: number;
}

interface CreateOrderRequest {
  userAddress: string;
  sourceChain: string;
  destChain: string;
  sourceToken: string;
  destToken: string;
  amount: string;
  useTWAP?: boolean;
  twapConfig?: {
    timeWindow: number;
    intervalCount: number;
    maxSlippage: number;
  };
}

interface BridgeOrder {
  id: string;
  userAddress: string;
  sourceChain: string;
  destChain: string;
  sourceToken: string;
  destToken: string;
  amount: string;
  status: string;
  createdAt: string;
  ethereumTxHash?: string;
  cosmosTxHash?: string;
  twapOrderId?: string;
}

export const useBridge = () => {
  const [orders, setOrders] = useState<BridgeOrder[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [wsConnection, setWsConnection] = useState<WebSocket | null>(null);

  // Initialize WebSocket connection for real-time updates
  useEffect(() => {
    const wsUrl = API_BASE_URL.replace('http', 'ws').replace('/api/v1', '/ws');
    const ws = new WebSocket(wsUrl);

    ws.onopen = () => {
      console.log('WebSocket connected');
      setWsConnection(ws);
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === 'order_update') {
          handleOrderUpdate(data.order);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    };

    ws.onclose = () => {
      console.log('WebSocket disconnected');
      setWsConnection(null);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    return () => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.close();
      }
    };
  }, []);

  const handleOrderUpdate = (updatedOrder: BridgeOrder) => {
    setOrders(prevOrders => {
      const existingIndex = prevOrders.findIndex(order => order.id === updatedOrder.id);
      if (existingIndex >= 0) {
        const newOrders = [...prevOrders];
        newOrders[existingIndex] = updatedOrder;
        return newOrders;
      } else {
        return [updatedOrder, ...prevOrders];
      }
    });
  };

  const getQuote = async (request: BridgeQuoteRequest) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/bridge/quote`, request);
      return response.data;
    } catch (error) {
      console.error('Error getting bridge quote:', error);
      throw error;
    }
  };

  const getTWAPQuote = async (request: TWAPQuoteRequest) => {
    try {
      const response = await axios.post(`${API_BASE_URL}/twap/quote`, request);
      return response.data;
    } catch (error) {
      console.error('Error getting TWAP quote:', error);
      throw error;
    }
  };

  const createOrder = async (request: CreateOrderRequest): Promise<BridgeOrder> => {
    setIsLoading(true);
    try {
      const response = await axios.post(`${API_BASE_URL}/bridge/order`, request);
      const newOrder = response.data;
      setOrders(prev => [newOrder, ...prev]);
      return newOrder;
    } catch (error) {
      console.error('Error creating order:', error);
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  const getOrder = async (id: string): Promise<BridgeOrder> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/bridge/order/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error getting order:', error);
      throw error;
    }
  };

  const getOrders = async (): Promise<BridgeOrder[]> => {
    try {
      const response = await axios.get(`${API_BASE_URL}/bridge/orders`);
      const fetchedOrders = response.data.orders || [];
      setOrders(fetchedOrders);
      return fetchedOrders;
    } catch (error) {
      console.error('Error getting orders:', error);
      throw error;
    }
  };

  const getTWAPOrder = async (id: string) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/twap/order/${id}`);
      return response.data;
    } catch (error) {
      console.error('Error getting TWAP order:', error);
      throw error;
    }
  };

  const getTWAPSchedule = async (id: string) => {
    try {
      const response = await axios.get(`${API_BASE_URL}/twap/order/${id}/schedule`);
      return response.data;
    } catch (error) {
      console.error('Error getting TWAP schedule:', error);
      throw error;
    }
  };

  const cancelOrder = async (id: string) => {
    try {
      await axios.post(`${API_BASE_URL}/bridge/order/${id}/cancel`);
      // Update local state
      setOrders(prev => 
        prev.map(order => 
          order.id === id ? { ...order, status: 'cancelled' } : order
        )
      );
    } catch (error) {
      console.error('Error cancelling order:', error);
      throw error;
    }
  };

  const getMetrics = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/metrics`);
      return response.data;
    } catch (error) {
      console.error('Error getting metrics:', error);
      throw error;
    }
  };

  const getStatus = async () => {
    try {
      const response = await axios.get(`${API_BASE_URL}/status`);
      return response.data;
    } catch (error) {
      console.error('Error getting status:', error);
      throw error;
    }
  };

  return {
    orders,
    isLoading,
    wsConnection,
    getQuote,
    getTWAPQuote,
    createOrder,
    getOrder,
    getOrders,
    getTWAPOrder,
    getTWAPSchedule,
    cancelOrder,
    getMetrics,
    getStatus,
  };
};