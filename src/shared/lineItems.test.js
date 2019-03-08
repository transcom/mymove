import * as lineItems from './lineItems';

function runTests(items) {
  for (let item of items) {
    expect(lineItems.displayBaseQuantityUnits(item.test)).toEqual(item.expected);
  }
}

describe('lineItems', () => {
  describe('displayBaseQuantityUnits', () => {
    describe('for full pack(205A), full unpack(105C), origin service charge(135A), destination service charge(135B)', () => {
      it('should display fee weight truncated to 0 decimal places', () => {
        const items = [
          { test: { tariff400ng_item: { code: '105A' }, quantity_1: 5000000 }, expected: '500 lbs' },
          {
            test: { tariff400ng_item: { code: '105C' }, quantity_1: Number.MAX_SAFE_INTEGER },
            expected: '900,719,925,474 lbs',
          },
          { test: { tariff400ng_item: { code: '135A' }, quantity_1: 50000 }, expected: '5 lbs' },
          { test: { tariff400ng_item: { code: '135B' }, quantity_1: 51111 }, expected: '5 lbs' },
          { test: { tariff400ng_item: { code: '105A' }, quantity_1: null }, expected: '0 lbs' },
          { test: null, expected: undefined },
          { test: { tariff400ng_item: { code: '105A' }, quantity_1: '5000000' }, expected: '0 lbs' }, // doesn't work w/ strings
          // negitives act funny due to floor() - but we shouldn't have negitive quantities so meh
          { test: { tariff400ng_item: { code: '105A' }, quantity_1: -55000 }, expected: '-6 lbs' },
        ];
        runTests(items);
      });
    });
    describe('for Pack Reg Crate(105B) and UnPack Reg Crate(105E)', () => {
      describe('for original accessorials', () => {
        it('should display value in quantity_1', () => {
          const items = [
            { test: { tariff400ng_item: { code: '105B' }, quantity_1: 5000000 }, expected: '500.0000' },
            {
              test: { tariff400ng_item: { code: '105E' }, quantity_1: Number.MAX_SAFE_INTEGER },
              expected: '900,719,925,474.0991',
            },
            { test: null, expected: undefined },
          ];

          runTests(items);
        });
      });

      describe('for robust accessorials', () => {
        it('should dispaly volume in cubic feet truncated to 2 decimal places', () => {
          const items = [
            {
              test: { tariff400ng_item: { code: '105B' }, quantity_1: 5000000, crate_dimensions: {} },
              expected: '500.00 cu ft',
            },
            {
              test: { tariff400ng_item: { code: '105E' }, quantity_1: Number.MAX_SAFE_INTEGER, crate_dimensions: {} },
              expected: '900,719,925,474.09 cu ft',
            },
            { test: null, expected: undefined },
          ];

          runTests(items);
        });
      });
    });
    describe('for Linehaul Transportation(LHS) and 105E Fule Surcharge-LHS(16A)', () => {
      it('should display weight and milage', () => {
        const items = [
          {
            test: { tariff400ng_item: { code: 'LHS' }, quantity_1: 5000000, quantity_2: 55550000 },
            expected: '500 lbs, 5,555 mi',
          },
          {
            test: {
              tariff400ng_item: { code: 'LHS' },
              quantity_1: Number.MAX_SAFE_INTEGER,
              quantity_2: Number.MAX_SAFE_INTEGER,
            },
            expected: '900,719,925,474 lbs, 900,719,925,474 mi',
          },
          { test: null, expected: undefined },
        ];
        runTests(items);
      });
    });
  });
});
