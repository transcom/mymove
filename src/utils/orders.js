export const matchesOrdersType = (orders, ...ordersTypes) => {
  return ordersTypes?.includes(orders?.orders_type);
};

export default matchesOrdersType;
