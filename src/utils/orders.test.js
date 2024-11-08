import { matchesOrdersType } from 'utils/orders';
import { ORDERS_TYPE } from 'constants/orders';

describe('matchesOrdersType', () => {
  const PCSOrders = { orders_type: ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION };
  const retirementOrders = { orders_type: ORDERS_TYPE.RETIREMENT };
  const separationOrders = { orders_type: ORDERS_TYPE.SEPARATION };
  it.each([
    [PCSOrders, [ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION]],
    [PCSOrders, [ORDERS_TYPE.RETIREMENT, ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION]],
    [retirementOrders, [ORDERS_TYPE.RETIREMENT]],
    [retirementOrders, [ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION, ORDERS_TYPE.RETIREMENT]],
    [separationOrders, [ORDERS_TYPE.SEPARATION]],
    [separationOrders, [ORDERS_TYPE.RETIREMENT, ORDERS_TYPE.SEPARATION]],
  ])('returns true when orders matches at least one of the provided types', (orders, ordersTypes) => {
    expect(matchesOrdersType(orders, ...ordersTypes)).toEqual(true);
  });

  it.each([
    [PCSOrders, ORDERS_TYPE.RETIREMENT],
    [retirementOrders, ORDERS_TYPE.SEPARATION],
    [separationOrders, ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION],
  ])('returns false when the orders type does not match', (orders, ordersType) => {
    expect(matchesOrdersType(matchesOrdersType(orders, ordersType))).toEqual(false);
  });

  it.each([
    [undefined, ORDERS_TYPE.RETIREMENT],
    [null, ORDERS_TYPE.SEPARATION],
    [{}, ORDERS_TYPE.PERMANENT_CHANGE_OF_STATION],
  ])('returns false when the orders object does not contain an orders type', (orders, ordersType) => {
    expect(matchesOrdersType(matchesOrdersType(orders, ordersType))).toEqual(false);
  });

  it.each([
    [PCSOrders, undefined],
    [PCSOrders, null],
    [PCSOrders, ''],
  ])('returns false when the orders type value is falsey', (orders, ordersType) => {
    expect(matchesOrdersType(matchesOrdersType(orders, ordersType))).toEqual(false);
  });
});
