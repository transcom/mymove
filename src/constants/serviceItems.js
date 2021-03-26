const SERVICE_ITEM_STATUSES = {
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
};

// TODO - refactor to order keys in alphabet order
const SERVICE_ITEM_PARAM_KEYS = {
  WeightBilledActual: 'WeightBilledActual',
  WeightActual: 'WeightActual',
  WeightEstimated: 'WeightEstimated',
  DistanceZip3: 'DistanceZip3',
  DistanceZip5: 'DistanceZip5',
  ZipDestAddress: 'ZipDestAddress',
  ZipPickupAddress: 'ZipPickupAddress',
  PriceRateOrFactor: 'PriceRateOrFactor',
  IsPeak: 'IsPeak',
  ServiceAreaOrigin: 'ServiceAreaOrigin',
  RequestedPickupDate: 'RequestedPickupDate',
  ActualPickupDate: 'ActualPickupDate',
  EscalationCompounded: 'EscalationCompounded',
  EIAFuelPrice: 'EIAFuelPrice',
  FSCWeightBasedDistanceMultiplier: 'FSCWeightBasedDistanceMultiplier',
};

const SERVICE_ITEM_CALCULATION_LABELS = {
  BillableWeight: 'Billable weight (cwt)',
  Mileage: 'Mileage',
  BaselineLinehaulPrice: 'Baseline linehaul price',
  BaselineShorthaulPrice: 'Baseline shorthaul price',
  PriceEscalationFactor: 'Price escalation factor',
  TotalAmountRequested: 'Total amount requested',
  FuelSurchargePrice: 'Fuel surcharge price (per mi)',
  [SERVICE_ITEM_PARAM_KEYS.WeightBilledActual]: 'Shipment weight',
  [SERVICE_ITEM_PARAM_KEYS.WeightActual]: 'Shipment weight',
  [SERVICE_ITEM_PARAM_KEYS.WeightEstimated]: 'Estimated',
  [SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]: 'Zip',
  [SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]: 'Zip',
  // Domestic non-peak or Domestic peak
  [SERVICE_ITEM_PARAM_KEYS.IsPeak]: 'Domestic',
  [SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]: 'Origin service area',
  [SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]: 'Pickup date',
  [SERVICE_ITEM_PARAM_KEYS.ActualPickupDate]: 'Pickup date',
  [SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice]: 'EIA diesel',
  [SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier]: 'Weight-based distance multiplier',
};

const SERVICE_ITEM_CODES = {
  DLH: 'DLH',
  FSC: 'FSC',
  DSH: 'DSH',
};

// TODO - temporary, will remove once all service item calculations are implemented
const allowedServiceItemCalculations = [SERVICE_ITEM_CODES.DLH, SERVICE_ITEM_CODES.FSC, SERVICE_ITEM_CODES.DSH];

export {
  SERVICE_ITEM_STATUSES as default,
  SERVICE_ITEM_PARAM_KEYS,
  SERVICE_ITEM_CALCULATION_LABELS,
  SERVICE_ITEM_CODES,
  allowedServiceItemCalculations,
};
