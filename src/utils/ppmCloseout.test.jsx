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
