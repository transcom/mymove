import { formatCents } from 'utils/formatters';
import { MOVE_DOC_TYPE } from 'shared/constants';

export const getNextPage = (nextPage, lastPage, pageToRevisit) => {
  if (lastPage && lastPage.pathname.includes(pageToRevisit)) {
    return lastPage.pathname;
  }
  return nextPage;
};

export const formatExpenseType = (expenseType) => {
  if (typeof expenseType !== 'string') return '';
  let type = expenseType.toLowerCase().replace('_', ' ');
  return type.charAt(0).toUpperCase() + type.slice(1);
};

export const formatExpenseDocs = (expenseDocs) => {
  return expenseDocs.map((expense) => ({
    id: expense.id,
    amount: formatCents(expense.requested_amount_cents),
    type: formatExpenseType(expense.moving_expense_type),
    paymentMethod: expense.payment_method,
    uploads: expense.document.uploads,
  }));
};

export const calcNetWeight = (documents) =>
  documents.reduce((accum, { move_document_type, full_weight, empty_weight }) => {
    if (move_document_type === MOVE_DOC_TYPE.WEIGHT_TICKET_SET) {
      return accum + (full_weight - empty_weight);
    }
    return accum;
  }, 0);
