import PropTypes from 'prop-types';

export const HistoryLogValuesShape = PropTypes.arrayOf(
  PropTypes.shape({
    columnName: PropTypes.string.isRequired,
    columnValue: PropTypes.string,
  }),
);

/*
const modelToDisplayName = {
  updatedAt: 'Updated at',
  diversion: 'Diversion',
  billableWeightCap: 'Billable weight cap',
  usesExternalVendor: 'Uses external vendor',
  requestedDeliveryDate: 'Requested delivery date',
  scheduledPickupDate: 'Scheduled pickup date',
  status: 'Status',
  customerRemarks: 'Customer remarks',
  approvedDate: 'Approved date',
  actualPickupDate: 'Actual pickup date',
  primeEstimatedEeight: 'Prime estimated weight',
  counselorRemarks: 'Counselor remarks',
};

const dbFieldsToModel = {
  updated_at: 'updatedAt',
  diversion: 'diversion',
  billable_weight_cap: 'billableWeightCap',
  uses_external_vendor: 'usesExternalVendor',
  requested_delivery_date: 'requesteDeliveryDate',
  scheduled_pickup_date: 'scheduledPickupDate',
  status: 'status',
  customer_remarks: 'customerRemarks',
  approved_date: 'approvedDate',
  actual_pickup_date: 'actualPickupDate',
  prime_estimated_weight: 'primeEstimatedWeight',
  counselor_remarks: 'counselorRemarks',
  street_address_1: 'streetAddress1',
  street_address_2: 'streetAddress2',
  street_address_3: 'streetAddress3',
  postal_code: 'postalCode',
  city: 'city',
  country: 'country',
};
*/

export const eventNamesWithLabelledDetails = {
  counselingUpdateOrder: 'Updated orders', // ghc.yaml
  updateOrder: 'Updated orders', // ghc.yaml
  updateAllowance: 'Updated allowances', // ghc.yaml
  counselingUpdateAllowance: 'Updated allowances', // ghc.yaml
  updateMoveTaskOrder: 'Updated move', // ghc.yaml
  updateMTOShipment: 'Updated shipment', // ghc.yaml internal.yaml prime.
  updateOrders: 'Updated orders', // internal.yaml
  submitAmendedOrders: 'Updated orders', // internal.yaml
  updateMTOShipmentAddress: 'Updated shipment', // prime.yaml
};

export const eventNamesWithServiceItemDetails = {
  updateMTOServiceItem: 'Updated service item', // ghc.yaml
  createMTOServiceItem: 'Requested service item', // prime.yaml
};

export const eventNamesWithEmptyDetails = {
  createOrders: 'Submitted orders', // internal.yaml
  uploadAmendedOrders: 'Updated orders', // internal.yaml
  submitMoveForApproval: 'Submitted move', // internal.yaml
  createMTOShipment: 'Submitted/Requested shipments', // internal.yaml prime.yaml
};

export const retrieveValue = (nameToFind, values) => {
  return values.find((value) => value.columnName === nameToFind).columnValue;
};

export const eventNamePlainTextToDisplay = {
  approveShipment: () => 'Approved shipment',
  updateMTOServiceItemStatus: () => 'Service item status', // ghc.yaml Need to check status as well
  requestShipmentDiversion: () => 'Requested diversion', // ghc.yaml
  setFinancialReviewFlag: (changedValues) => {
    const financialReviewFlag = retrieveValue('financial_review_flag', changedValues);
    return financialReviewFlag === 'true' ? 'Move flagged for financial review' : 'Move unflagged for financial review';
  },
  requestShipmentCancellation: () => 'Shipment cancelled',
  updateMoveTaskOrderStatus: (changedValues) => {
    const status = retrieveValue('status', changedValues);
    return status === 'APPROVED' ? 'Created Move Task Order (MTO)' : 'Rejected Move Task Order (MTO)';
  },
};

export const eventNamesWithPlainTextDetails = {
  approveShipment: 'Approved shipment', // ghc.yaml
  requestShipmentDiversion: 'Requested diversion', // ghc.yaml
  updateMTOServiceItemStatus: 'Service item status', // ghc.yaml Need to check status as well
  setFinancialReviewFlag: 'Flagged move', // ghc.yaml
  requestShipmentCancellation: 'Updated shipment', // ghc.yaml
  updateMoveTaskOrderStatus: 'Move task order status', // ghc.yaml Need to check status as well
};

export const historyLogEventNameDisplay = {
  // operationId, UI Display
  counselingUpdateOrder: 'Updated orders', // ghc.yaml
  updateOrder: 'Updated orders', // ghc.yaml
  updateAllowance: 'Updated allowances', // ghc.yaml
  counselingUpdateAllowance: 'Updated allowances', // ghc.yaml
  updateMoveTaskOrder: 'Updated move', // ghc.yaml
  updateMTOShipment: 'Updated shipment', // ghc.yaml internal.yaml prime.yaml
  approveShipment: 'Approved shipment', // ghc.yaml
  requestShipmentDiversion: 'Requested diversion', // ghc.yaml
  updateMTOServiceItem: 'Updated service item', // ghc.yaml
  updateMTOServiceItemStatus: '', // ghc.yaml
  updateMoveTaskOrderStatus: '', // ghc.yaml
  setFinancialReviewFlag: 'Flagged move', // ghc.yaml
  requestShipmentCancellation: 'Updated shipment', // ghc.yaml
  createOrders: 'Submitted orders', // internal.yaml
  updateOrders: 'Updated orders', // internal.yaml
  uploadAmendedOrders: 'Updated orders', // internal.yaml
  submitMoveForApproval: 'Submitted move', // internal.yaml
  submitAmendedOrders: 'Updated orders', // internal.yaml
  createMTOShipment: 'Submitted/Requested shipments', // internal.yaml prime.yaml
  updateMTOShipmentAddress: 'Updated shipment', // prime.yaml
  createMTOServiceItem: 'Requested service item', // prime.yaml
};

export function getHistoryLogEventNameDisplay({ eventName /* operationId */, changedValues }) {
  switch (eventName) {
    case 'updateMTOServiceItemStatus': {
      // find 'columnName' with 'columnValue'
      const status = changedValues?.find((changedValue) => changedValue?.columnName === 'status');
      switch (status?.columnValue) {
        case 'APPROVED':
          return 'Approved service item';
        case 'REJECTED':
          return 'Rejected service item';
        default:
          return '';
      }
    }
    case 'updateMoveTaskOrderStatus': {
      // find 'columnName' with 'columnValue'
      const status = changedValues?.find((changedValue) => changedValue?.columnName === 'status');
      switch (status?.columnValue) {
        case 'APPROVED':
          return 'Move approved';
        case 'REJECTED':
          return 'Move rejected';
        default:
          return '';
      }
    }
    default:
      return historyLogEventNameDisplay[eventName];
  }
}

export default getHistoryLogEventNameDisplay;
