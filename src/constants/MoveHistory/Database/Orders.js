import {
  ORDERS_BRANCH_OPTIONS,
  ORDERS_RANK_OPTIONS,
  ORDERS_TYPE_DETAILS_OPTIONS,
  ORDERS_TYPE_OPTIONS,
} from 'constants/orders';

// This is to map the human-readable text to the options
export default {
  ...ORDERS_BRANCH_OPTIONS,
  ...ORDERS_TYPE_DETAILS_OPTIONS,
  ...ORDERS_TYPE_OPTIONS,
  ...ORDERS_RANK_OPTIONS,
};
