import PropTypes from 'prop-types';

import { shipmentTypes } from 'constants/shipments';

export const HistoryLogValuesShape = PropTypes.object;
export const HistoryLogContextShape = PropTypes.object;

export const HistoryLogRecordShape = PropTypes.shape({
  context: HistoryLogContextShape,
  eventName: PropTypes.string,
  changedValues: HistoryLogValuesShape,
  oldValues: HistoryLogValuesShape,
});

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

export const eventNamesWithLabeledDetails = {
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

export const shipmentOptionToDisplay = {
  HHG_OUTOF_NTS_DOMESTIC: 'NTS-release',
  HHG_INTO_NTS_DOMESTIC: 'NTS',
  HHG: 'HHG',
  PPM: 'PPM',
  HHG_SHORTHAUL_DOMESTIC: 'HHG_SHORTHAUL_DOMESTIC',
};

export const eventNamePlainTextToDisplay = {
  approveShipment: (_, oldValues) => `${shipmentTypes[oldValues.shipment_type]} shipment`,
  approveShipmentDiversion: (_, oldValues) => `${shipmentTypes[oldValues.shipment_type]} shipment`,
  updateMTOServiceItemStatus: (historyRecord) =>
    `${shipmentOptionToDisplay[historyRecord.context.shipment_type]} shipment, ${historyRecord.context.name}`,
  requestShipmentDiversion: (historyRecord) =>
    `Requested diversion for ${shipmentOptionToDisplay[historyRecord.oldValues.shipment_type]} shipment`, // ghc.yaml
  setFinancialReviewFlag: (historyRecord) => {
    return historyRecord.changedValues.financial_review_flag === 'true'
      ? 'Move flagged for financial review'
      : 'Move unflagged for financial review';
  },
  requestShipmentCancellation: (historyRecord) =>
    `Requested cancellation for ${shipmentOptionToDisplay[historyRecord.oldValues.shipment_type]} shipment`,
  updateMoveTaskOrderStatus: (historyRecord) => {
    return historyRecord.changedValues.status === 'APPROVED'
      ? 'Created Move Task Order (MTO)'
      : 'Rejected Move Task Order (MTO)';
  },
};

export const eventNamesWithPlainTextDetails = {
  approveShipment: 'Approved shipment', // ghc.yaml
  approveShipmentDiversion: 'Approved shipment',
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
  approveShipmentDiversion: 'Approved shipment',
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
      switch (changedValues.status) {
        case 'APPROVED':
          return 'Approved service item';
        case 'REJECTED':
          return 'Rejected service item';
        default:
          return '';
      }
    }
    case 'updateMoveTaskOrderStatus': {
      switch (changedValues.status) {
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
