const SERVICE_ITEM_STATUSES = {
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
};

const SERVICE_ITEM_PARAM_KEYS = {
  WeightActual: 'WeightActual',
  WeightEstimated: 'WeightEstimated',
  DistanceZip3: 'DistanceZip3',
  ZipDestAddress: 'ZipDestAddress',
  ZipPickupAddress: 'ZipPickupAddress',
  PriceRateOrFactor: 'PriceRateOrFactor',
  IsPeak: 'IsPeak',
  ServiceAreaOrigin: 'ServiceAreaOrigin',
  RequestedPickupDate: 'RequestedPickupDate',
  EscalationCompounded: 'EscalationCompounded',
};

const SERVICE_ITEM_CALCULATION_LABELS = {
  BillableWeight: 'Billable weight (cwt)',
  Mileage: 'Mileage',
  BaselineLinehaulPrice: 'Baseline linehaul price',
  PriceEscalationFactor: 'Price escalation factor',
  TotalAmountRequested: 'Total amount requested',
  [SERVICE_ITEM_PARAM_KEYS.WeightActual]: 'Shipment weight',
  [SERVICE_ITEM_PARAM_KEYS.WeightEstimated]: 'Estimated',
  [SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]: 'Zip',
  [SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]: 'Zip',
  // Domestic non-peak or Domestic peak
  [SERVICE_ITEM_PARAM_KEYS.IsPeak]: 'Domestic',
  [SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]: 'Origin service area',
  [SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]: 'Pickup date',
};

const SERVICE_ITEM_CODES = {
  DLH: 'DLH',
};

// TODO - temporary, will remove once all service item calculations are implemented
const allowedServiceItemCalculations = [SERVICE_ITEM_CODES.DLH];

export {
  SERVICE_ITEM_STATUSES as default,
  SERVICE_ITEM_PARAM_KEYS,
  SERVICE_ITEM_CALCULATION_LABELS,
  SERVICE_ITEM_CODES,
  allowedServiceItemCalculations,
};
