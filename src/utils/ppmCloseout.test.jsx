import {
  calculateNetWeightForProGearWeightTicket,
  calculateNetWeightForWeightTicket,
  calculateTotalMovingExpensesAmount,
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForWeightTickets,
} from 'utils/ppmCloseout';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { createCompleteProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';

describe('calculateNetWeightForWeightTicket', () => {
  it.each([
    [0, 400, 400],
    [15000, 18000, 3000],
    [null, 1500, 0],
    [0, null, 0],
    [null, null, 0],
    [undefined, 1500, 0],
    [0, undefined, 0],
    [undefined, undefined, 0],
    ['not a number', 1500, 0],
    [0, 'not a number', 0],
    ['not a number', 'not a number', 0],
  ])(
    `calculates net weight properly | emptyWeight: %s | fullWeight: %s | expectedNetWeight: %s`,
    (emptyWeight, fullWeight, expectedNetWeight) => {
      const weightTicket = createCompleteWeightTicket(
        {},
        {
          emptyWeight,
          fullWeight,
        },
      );

      expect(calculateNetWeightForWeightTicket(weightTicket)).toEqual(expectedNetWeight);
    },
  );
});

describe('calculateTotalNetWeightForWeightTickets', () => {
  it.each([
    [[{ emptyWeight: 0, fullWeight: 400 }], 400],
    [
      [
        { emptyWeight: 0, fullWeight: 400 },
        { emptyWeight: 15000, fullWeight: 18000 },
      ],
      3400,
    ],
    [
      [
        { emptyWeight: null, fullWeight: 400 },
        { emptyWeight: 14000, fullWeight: 17000 },
      ],
      3000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: null },
        { emptyWeight: 14000, fullWeight: 19000 },
      ],
      5000,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null },
        { emptyWeight: 14000, fullWeight: 18000 },
      ],
      4000,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null },
        { emptyWeight: null, fullWeight: null },
      ],
      0,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: 400 },
        { emptyWeight: 14000, fullWeight: 17000 },
      ],
      3000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: undefined },
        { emptyWeight: 14000, fullWeight: 19000 },
      ],
      5000,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: undefined },
        { emptyWeight: 14000, fullWeight: 18000 },
      ],
      4000,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: undefined },
        { emptyWeight: undefined, fullWeight: undefined },
      ],
      0,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 400 },
        { emptyWeight: 14000, fullWeight: 17000 },
      ],
      3000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 19000 },
      ],
      5000,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 18000 },
      ],
      4000,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
        { emptyWeight: 'not a number', fullWeight: 'not a number' },
      ],
      0,
    ],
    [[], 0],
  ])(`calculates total net weight properly`, (weightTicketsFields, expectedNetWeight) => {
    const weightTickets = [];

    weightTicketsFields.forEach((fieldOverrides) => {
      weightTickets.push(createCompleteWeightTicket({}, fieldOverrides));
    });

    expect(calculateTotalNetWeightForWeightTickets(weightTickets)).toEqual(expectedNetWeight);
  });
});

describe('calculateNetWeightForProGearWeightTicket', () => {
  it.each([
    [0, 400, null, 400],
    [15000, 17000, null, 2000],
    [null, null, 1200, 1200],
    [null, 1500, null, 0],
    [0, null, null, 0],
    [null, null, null, 0],
    [undefined, 1500, undefined, 0],
    [0, undefined, undefined, 0],
    [undefined, undefined, undefined, 0],
    ['not a number', 1500, 'not a number', 0],
    [0, 'not a number', 'not a number', 0],
    ['not a number', 'not a number', 'not a number', 0],
  ])(
    `calculates net weight properly | emptyWeight: %s | fullWeight: %s | constructedWeight: %s | expectedNetWeight: %s`,
    (emptyWeight, fullWeight, constructedWeight, expectedNetWeight) => {
      const proGearWeightTicket = createCompleteProGearWeightTicket(
        {},
        {
          emptyWeight,
          fullWeight,
          constructedWeight,
        },
      );

      expect(calculateNetWeightForProGearWeightTicket(proGearWeightTicket)).toEqual(expectedNetWeight);
    },
  );
});

describe('calculateTotalNetWeightForProGearWeightTickets', () => {
  it.each([
    [[{ emptyWeight: 0, fullWeight: 400, constructedWeight: null }], 400],
    [
      [
        { emptyWeight: 0, fullWeight: 400, constructedWeight: null },
        { emptyWeight: 15000, fullWeight: 16000, constructedWeight: null },
      ],
      1400,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null, constructedWeight: 250 },
        { emptyWeight: 15000, fullWeight: 16500, constructedWeight: null },
      ],
      1750,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null, constructedWeight: 250 },
        { emptyWeight: null, fullWeight: null, constructedWeight: 500 },
      ],
      750,
    ],
    [
      [
        { emptyWeight: null, fullWeight: 400, constructedWeight: null },
        { emptyWeight: 14000, fullWeight: 16000, constructedWeight: null },
      ],
      2000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: null, constructedWeight: null },
        { emptyWeight: 14000, fullWeight: 15500, constructedWeight: null },
      ],
      1500,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null, constructedWeight: null },
        { emptyWeight: 14000, fullWeight: 14500, constructedWeight: null },
      ],
      500,
    ],
    [
      [
        { emptyWeight: null, fullWeight: null, constructedWeight: null },
        { emptyWeight: null, fullWeight: null, constructedWeight: null },
      ],
      0,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: 400, constructedWeight: undefined },
        { emptyWeight: 14000, fullWeight: 16000, constructedWeight: undefined },
      ],
      2000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: undefined, constructedWeight: undefined },
        { emptyWeight: 14000, fullWeight: 15500, constructedWeight: undefined },
      ],
      1500,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: undefined, constructedWeight: undefined },
        { emptyWeight: 14000, fullWeight: 14500, constructedWeight: undefined },
      ],
      500,
    ],
    [
      [
        { emptyWeight: undefined, fullWeight: undefined, constructedWeight: undefined },
        { emptyWeight: undefined, fullWeight: undefined, constructedWeight: undefined },
      ],
      0,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 400, constructedWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 16000, constructedWeight: 'not a number' },
      ],
      2000,
    ],
    [
      [
        { emptyWeight: 0, fullWeight: 'not a number', constructedWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 15500, constructedWeight: 'not a number' },
      ],
      1500,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number', constructedWeight: 'not a number' },
        { emptyWeight: 14000, fullWeight: 14500, constructedWeight: 'not a number' },
      ],
      500,
    ],
    [
      [
        { emptyWeight: 'not a number', fullWeight: 'not a number', constructedWeight: 'not a number' },
        { emptyWeight: 'not a number', fullWeight: 'not a number', constructedWeight: 'not a number' },
      ],
      0,
    ],
    [[], 0],
  ])(`calculates total net weight properly`, (proGearWeightTicketsFields, expectedNetWeight) => {
    const proGearWeightTickets = [];

    proGearWeightTicketsFields.forEach((fieldOverrides) => {
      proGearWeightTickets.push(createCompleteProGearWeightTicket({}, fieldOverrides));
    });

    expect(calculateTotalNetWeightForProGearWeightTickets(proGearWeightTickets)).toEqual(expectedNetWeight);
  });
});

describe('calculateTotalMovingExpensesAmount', () => {
  it.each([
    [[{ amount: 400 }], 400],
    [[{ amount: 300 }, { amount: 200 }], 500],
    [[{ amount: 300 }, { amount: 200 }, { amount: 600 }], 1100],
    [[{ amount: null }, { amount: 350 }], 350],
    [[{ amount: undefined }, { amount: 250 }], 250],
    [[{ amount: 'not a number' }, { amount: 600 }], 600],
    [[{ amount: 750 }, { amount: null }], 750],
    [[{ amount: null }, { amount: null }], 0],
    [[], 0],
  ])(`calculates total net weight properly`, (movingExpensesFields, expectedTotal) => {
    const expenses = [];

    movingExpensesFields.forEach((fieldOverrides) => {
      expenses.push(createCompleteMovingExpense({}, fieldOverrides));
    });

    expect(calculateTotalMovingExpensesAmount(expenses)).toEqual(expectedTotal);
  });
});
