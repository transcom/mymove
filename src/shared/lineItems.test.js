import * as lineItems from './lineItems';

describe('lineItems', () => {
  describe('displayBaseQuantityUnits', () => {
    it('display full pack, full unpack, origin and dest fee weight truncated to 0 decimal places', () => {
      const item105A = { tariff400ng_item: { code: '105A' }, quantity_1: 5000000 };
      const item105C = { tariff400ng_item: { code: '105C' }, quantity_1: Number.MAX_SAFE_INTEGER };

      const item135A = { tariff400ng_item: { code: '135A' }, quantity_1: 50000 };
      const item135B = { tariff400ng_item: { code: '135B' }, quantity_1: 51111 };
      const itemQuantityNull = { tariff400ng_item: { code: '105A' }, quantity_1: null };
      const itemNegitive = { tariff400ng_item: { code: '105A' }, quantity_1: -55000 };
      const item105AString = { tariff400ng_item: { code: '105A' }, quantity_1: '5000000' };
      const itemNull = null;

      expect(lineItems.displayBaseQuantityUnits(item105A)).toEqual('500 lbs');
      expect(lineItems.displayBaseQuantityUnits(item105C)).toEqual('900,719,925,474 lbs');
      expect(lineItems.displayBaseQuantityUnits(item135A)).toEqual('5 lbs');
      expect(lineItems.displayBaseQuantityUnits(item135B)).toEqual('5 lbs');
      expect(lineItems.displayBaseQuantityUnits(itemQuantityNull)).toEqual('0 lbs');
      expect(lineItems.displayBaseQuantityUnits(itemNull)).toEqual(undefined);
      expect(lineItems.displayBaseQuantityUnits(item105AString)).toEqual('0 lbs'); // doesn't work w/ strings
      // negitives act funny due to floor() - but we shouldn't have negitive quantities so meh
      expect(lineItems.displayBaseQuantityUnits(itemNegitive)).toEqual('-6 lbs');
    });
    it('display 105B/E for original accessorials', () => {
      const item105B = { tariff400ng_item: { code: '105B' }, quantity_1: 5000000 };
      const item105E = { tariff400ng_item: { code: '105E' }, quantity_1: Number.MAX_SAFE_INTEGER };
      const itemNull = null;

      expect(lineItems.displayBaseQuantityUnits(item105B)).toEqual('500.0000');
      expect(lineItems.displayBaseQuantityUnits(item105E)).toEqual('900,719,925,474.0991');
      expect(lineItems.displayBaseQuantityUnits(itemNull)).toEqual(undefined);
    });

    it('display 105B/E cu ft for robust accessorials', () => {
      const item105B = { tariff400ng_item: { code: '105B' }, quantity_1: 5000000, crate_dimensions: {} };
      const item105E = {
        tariff400ng_item: { code: '105E' },
        quantity_1: Number.MAX_SAFE_INTEGER,
        crate_dimensions: {},
      };
      const itemNull = null;

      expect(lineItems.displayBaseQuantityUnits(item105B)).toEqual('500.00 cu ft');
      expect(lineItems.displayBaseQuantityUnits(item105E)).toEqual('900,719,925,474.09 cu ft');
      expect(lineItems.displayBaseQuantityUnits(itemNull)).toEqual(undefined);
    });
  });
});
