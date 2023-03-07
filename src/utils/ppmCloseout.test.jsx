import {
  calculateNetWeightForProGearWeightTicket,
  calculateNetWeightForWeightTicket,
  calculateTotalMovingExpensesAmount,
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForWeightTickets,
  formatExpenseItems,
} from 'utils/ppmCloseout';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { createCompleteProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteMovingExpense, createCompleteSITMovingExpense } from 'utils/test/factories/movingExpense';
import { expenseTypeLabels } from 'constants/ppmExpenseTypes';

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
    [0, 0],
    [15000, 15000],
    [null, 0],
    [undefined, 0],
    ['not a number', 0],
  ])(
    `calculates net weight properly | emptyWeight: %s | fullWeight: %s | constructedWeight: %s | expectedNetWeight: %s`,
    (weight, expectedNetWeight) => {
      const proGearWeightTicket = createCompleteProGearWeightTicket(
        {},
        {
          weight,
        },
      );

      expect(calculateNetWeightForProGearWeightTicket(proGearWeightTicket)).toEqual(expectedNetWeight);
    },
  );
});

describe('calculateTotalNetWeightForProGearWeightTickets', () => {
  it.each([
    [[{ weight: 0 }], 0],
    [[{ weight: 0 }, { weight: 15000 }], 15000],
    [[{ weight: null }], 0],
    [[{ weight: null }, { weight: 15000 }], 15000],
    [[{ weight: undefined }], 0],
    [[{ weight: undefined }, { weight: 15000 }], 15000],
    [[{ weight: 'not a number' }], 0],
    [[{ weight: 'not a number' }, { weight: 15000 }], 15000],
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

describe('formatExpenseItems', () => {
  it.each([
    [
      [createCompleteMovingExpense()],
      {
        movingExpenseType: expenseTypeLabels.PACKING_MATERIALS,
        description: 'Medium and large boxes',
        amount: '$75.00',
      },
    ],
    [
      [createCompleteSITMovingExpense()],
      { movingExpenseType: expenseTypeLabels.STORAGE, description: 'Storage while away', amount: '$75.00' },
    ],
  ])(`formats moving expense for review`, (movingExpense, expectedMovingExpense) => {
    const formattedExpense = formatExpenseItems(movingExpense, '', {}, () => {})[0];

    expect(formattedExpense.rows[0].value).toEqual(expectedMovingExpense.movingExpenseType);
    expect(formattedExpense.rows[1].value).toEqual(expectedMovingExpense.description);
    expect(formattedExpense.rows[2].value).toEqual(expectedMovingExpense.amount);

    if (expectedMovingExpense.movingExpenseType === expenseTypeLabels.STORAGE) {
      expect(formattedExpense.rows[3].value).toEqual(6);
    }
  });
});
