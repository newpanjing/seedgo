export interface Room {
  id: string;
  name: string;
  type: 'chess' | 'mahjong' | 'billiards';
  status: 'available' | 'occupied' | 'cleaning' | 'maintenance';
  currentSessionId?: string;
  hourlyRate: number;
}

export interface Order {
  id: string;
  productId: string;
  productName: string;
  quantity: number;
  price: number;
  total: number;
  timestamp: string;
}

export interface Session {
  id: string;
  roomId: string;
  startTime: string;
  endTime?: string;
  status: 'active' | 'completed';
  totalAmount: number;
  discount?: number;
  rounding?: number;
  orders: Order[];
  durationMinutes: number;
}
