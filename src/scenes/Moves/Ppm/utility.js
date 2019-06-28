import { formatCents } from 'shared/formatters';

// maps int to int with ordinal 1 -> 1st, 2 -> 2nd, 3rd ...
export const intToOrdinal = n => {
  const s = ['th', 'st', 'nd', 'rd'];
  const v = n % 100;
  // eslint-disable-next-line security/detect-object-injection
  return n + (s[(v - 20) % 10] || s[v] || s[0]);
};

export const getNextPage = (nextPage, lastPage, pageToRevisit) => {
  if (lastPage && lastPage.pathname.includes(pageToRevisit)) {
    return lastPage.pathname;
  }
  return nextPage;
};

export const formatExpenseType = expenseType => {
  if (typeof expenseType !== 'string') return '';
  let type = expenseType.toLowerCase().replace('_', ' ');
  return type.charAt(0).toUpperCase() + type.slice(1);
};

export const formatExpenseDocs = expenseDocs => {
  return expenseDocs.map(expense => ({
    id: expense.id,
    amount: formatCents(expense.requested_amount_cents),
    type: formatExpenseType(expense.moving_expense_type),
    paymentMethod: expense.payment_method,
  }));
};
