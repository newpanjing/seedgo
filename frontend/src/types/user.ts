export interface User {
  id: string;
  name: string;
  phone: string;
  level: 'normal' | 'silver' | 'gold' | 'diamond';
  discount: number; // 0-1, e.g. 0.9 for 10% off
  balance: number;
  points: number;
  totalConsumed: number;
  tags: string[];
  createdAt: string;
  lastVisit: string;
}
