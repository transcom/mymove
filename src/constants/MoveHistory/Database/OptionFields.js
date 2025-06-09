import {
  ORDERS_BRANCH_OPTIONS,
  ORDERS_TYPE_DETAILS_OPTIONS,
  ORDERS_TYPE_OPTIONS,
  ORDERS_DEPARTMENT_INDICATOR,
} from 'constants/orders';
import { shipmentDestinationTypes } from 'constants/shipments';
import { SIT_EXTENSION_REASONS } from 'constants/sitExtensions';

// This is to map the human-readable text to the options
export default {
  ...ORDERS_BRANCH_OPTIONS,
  ...ORDERS_TYPE_DETAILS_OPTIONS,
  ...ORDERS_TYPE_OPTIONS,
  ...ORDERS_DEPARTMENT_INDICATOR,
  ...shipmentDestinationTypes,
  ...SIT_EXTENSION_REASONS,
};
