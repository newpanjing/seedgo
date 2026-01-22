export type CategoryType = 'chess' | 'mahjong' | 'billiards' | 'rental' | 'other';

export interface Category {
  id: string;
  name: string;
  type: CategoryType;
}

export interface Product {
  id: string;
  name: string;
  categoryId: string;
  price: number;
  stock: number;
  imageUrl?: string;
  description?: string;
  pricingType: 'fixed' | 'hourly' | 'package';
  status: 'active' | 'inactive';
  createdAt: string;
}
