import { calculateTotalMovingExpensesAmount, formatExpenseItems } from 'utils/ppmCloseout';
import { createCompleteMovingExpense, createCompleteSITMovingExpense } from 'utils/test/factories/movingExpense';
import { expenseTypeLabels } from 'constants/ppmExpenseTypes';

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
