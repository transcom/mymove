import { get, filter } from 'lodash';
import reduceReducers from 'reduce-reducers';

import { GetPpmIncentive, GetExpenseSummary } from './api.js';

import * as ReduxHelpers from 'shared/ReduxHelpers';
import { MOVE_DOC_STATUS } from 'shared/constants';

const GET_PPM_INCENTIVE = 'GET_PPM_INCENTIVE';
const GET_PPM_EXPENSE_SUMMARY = 'GET_PPM_EXPENSE_SUMMARY';
const CLEAR_PPM_INCENTIVE = 'CLEAR_PPM_INCENTIVE';
export const getIncentiveActionType = ReduxHelpers.generateAsyncActionTypes(GET_PPM_INCENTIVE);
export const getPpmIncentive = ReduxHelpers.generateAsyncActionCreator(GET_PPM_INCENTIVE, GetPpmIncentive);

export const getExpenseSummaryActionType = ReduxHelpers.generateAsyncActionTypes(GET_PPM_EXPENSE_SUMMARY);
export const getPpmExpenseSummary = ReduxHelpers.generateAsyncActionCreator(GET_PPM_EXPENSE_SUMMARY, GetExpenseSummary);
const summaryReducer = ReduxHelpers.generateAsyncReducer(GET_PPM_EXPENSE_SUMMARY, (v) => {
  return {
    summary: { ...v },
  };
});

export const clearPpmIncentive = () => ({ type: CLEAR_PPM_INCENTIVE });

export const getTabularExpenses = (expenseData, movingExpenseSchema) => {
  if (!expenseData || !movingExpenseSchema) return [];
  const expenses = movingExpenseSchema.enum.map((type) => {
    const item = expenseData.categories.find((item) => item.category === type);
    if (!item)
      return {
        type: get(movingExpenseSchema['x-display-value'], type),
        GTCC: null,
        other: null,
        total: null,
      };
    return {
      type: get(movingExpenseSchema['x-display-value'], type),
      GTCC: get(item, 'payment_methods.GTCC', null),
      other: get(item, 'payment_methods.OTHER', null),
      total: item.total,
    };
  });
  expenses.push({
    type: 'Total',
    GTCC: get(expenseData, 'grand_total.payment_method_totals.GTCC'),
    other: get(expenseData, 'grand_total.payment_method_totals.OTHER'),
    total: get(expenseData, 'grand_total.total'),
  });
  return expenses;
};

export const getDocsByStatusAndType = (documents, statusToExclude, typeToExclude) => {
  return filter(documents, (expense) => {
    if (!statusToExclude) {
      return expense.move_document_type !== typeToExclude;
    }
    if (!typeToExclude) {
      return ![MOVE_DOC_STATUS.EXCLUDE, statusToExclude].includes(expense.status);
    }

    return (
      ![MOVE_DOC_STATUS.EXCLUDE, statusToExclude].includes(expense.status) &&
      expense.move_document_type !== typeToExclude
    );
  });
};

function clearReducer(state, action) {
  if (action.type === CLEAR_PPM_INCENTIVE) return { ...state, calculation: null };
  return state;
}
const incentiveReducer = ReduxHelpers.generateAsyncReducer(GET_PPM_INCENTIVE, (v) => ({
  calculation: { ...v },
}));

export default reduceReducers(clearReducer, incentiveReducer, summaryReducer);
