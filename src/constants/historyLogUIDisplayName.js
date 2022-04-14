import PropTypes from 'prop-types';

export const HistoryLogValuesShape = PropTypes.object;
export const HistoryLogContextShape = PropTypes.object;

export const HistoryLogRecordShape = PropTypes.shape({
  action: PropTypes.string,
  changedValues: HistoryLogValuesShape,
  context: HistoryLogContextShape,
  eventName: PropTypes.string,
  oldValues: HistoryLogValuesShape,
  tableName: PropTypes.string,
});

/*
export const modelToDisplayName = {
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
  serviceOrderNumber: 'Service order number',
  tacType: 'TAC type',
  sacType: 'SAC type',
};

export const dbFieldsToModel = {
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
  service_order_number: 'serviceOrderNumber',
  tac_type: 'tacType',
  sac_type: 'sacType',
};
*/

export const dbFieldToDisplayName = {
  updated_at: 'Updated at',
  diversion: 'Diversion',
  billable_weight_cap: 'Billable weight cap',
  uses_external_vendor: 'Uses external vendor',
  requested_delivery_date: 'Requested delivery date',
  scheduled_pickup_date: 'Scheduled pickup date',
  status: 'Status',
  customer_remarks: 'Customer remarks',
  approved_date: 'Approved date',
  actual_pickup_date: 'Actual pickup date',
  prime_estimated_weight: 'Prime estimated weight',
  counselor_remarks: 'Counselor remarks',
  service_order_number: 'Service order number',
  tac_type: 'TAC type',
  sac_type: 'SAC type',
  authorized_weight: 'Authorized weight',
};

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
  updateBillableWeight: 'Updated move',
};

export const eventNamesWithServiceItemDetails = {
  updateMTOServiceItem: 'Updated service item', // ghc.yaml
  createMTOServiceItem: 'Requested service item', // prime.yaml
};
