import PropTypes from 'prop-types';

export const HistoryLogValuesShape = PropTypes.object;
export const HistoryLogContextShape = PropTypes.arrayOf(PropTypes.object);

export const HistoryLogRecordShape = PropTypes.shape({
  action: PropTypes.string,
  changedValues: HistoryLogValuesShape,
  context: HistoryLogContextShape,
  eventName: PropTypes.string,
  oldValues: HistoryLogValuesShape,
  tableName: PropTypes.string,
});

export const dbActions = {
  UPDATE: 'UPDATE',
  INSERT: 'INSERT',
  DELETE: 'DELETE',
};

export const dbTables = {
  moves: 'moves',
  mto_shipments: 'mto_shipments',
  orders: 'orders',
  mto_service_items: 'mto_service_items',
  entitlements: 'entitlements',
  payment_requests: 'payment_requests',
  mto_agents: 'mto_agents',
  reweighs: 'reweighs',
  addresses: 'addresses',
};

export const dbFieldToDisplayName = {
  updated_at: 'Updated at',
  diversion: 'Diversion',
  billable_weight_cap: 'Billable weight',
  billable_weight_justification: 'Billable weight remarks',
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
  report_by_date: 'Report by date',
  issue_date: 'Orders date',
  orders_type: 'Orders type',
  orders_type_detail: 'Orders type detail',
  origin_duty_location_name: 'Origin duty location name',
  new_duty_location_name: 'New duty location name',
  orders_number: 'Orders number',
  tac: 'HHG TAC',
  sac: 'HHG SAC',
  nts_tac: 'NTS TAC',
  nts_sac: 'NTS SAC',
  department_indicator: 'Dept. indicator',
  amended_orders_acknowledged_at: 'Amended orders acknowledged at',
  storage_in_transit: 'Storage in transit (SIT)',
  dependents_authorized: 'Dependents',
  pro_gear_weight: 'Pro-gear',
  pro_gear_weight_spouse: 'Spouse pro-gear',
  required_medical_equipment_weight: 'RME',
  organizational_clothing_and_individual_equipment: 'OCIE',
  requested_pickup_date: 'Requested pickup date',
  grade: 'Rank',
  shipment_type: 'Shipment type',
  pickup_address: 'Origin address',
  destination_address: 'Destination address',
  receiving_agent: 'Receiving agent',
  releasing_agent: 'Releasing agent',
  tio_remarks: 'Max billable weight remark',
};

export const dbWeightFields = {
  billable_weight_cap: 'billable_weight_cap',
  prime_estimated_weight: 'prime_estimated_weight',
  authorized_weight: 'authorized_weight',
  pro_gear_weight: 'pro_gear_weight',
  pro_gear_weight_spouse: 'pro_gear_weight_spouse',
  required_medical_equipment_weight: 'required_medical_equipment_weight',
};

export const dbDateFields = {
  requested_delivery_date: 'requested_delivery_date',
  scheduled_pickup_date: 'scheduled_pickup_date',
  approved_date: 'approved_date',
  actual_pickup_date: 'actual_pickup_date',
  report_by_date: 'report_by_date',
  issue_date: 'issue_date',
  requested_pickup_date: 'requested_pickup_date',
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
