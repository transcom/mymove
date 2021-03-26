const SERVICE_ITEM_STATUSES = {
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
};

// TODO - refactor to order keys in alphabet order
const SERVICE_ITEM_PARAM_KEYS = {
  ContractYearName: 'ContractYearName',
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
  ServicesScheduleOrigin: 'ServicesScheduleOrigin',
  RequestedPickupDate: 'RequestedPickupDate',
  ActualPickupDate: 'ActualPickupDate',
  EscalationCompounded: 'EscalationCompounded',
  EIAFuelPrice: 'EIAFuelPrice',
  FSCWeightBasedDistanceMultiplier: 'FSCWeightBasedDistanceMultiplier',
  OriginPrice: 'OriginPrice',
  ServiceSchedule: 'ServiceSchedule',
  NumberDaysSIT: 'NumberDaysSIT',
};

const SERVICE_ITEM_CALCULATION_LABELS = {
  BillableWeight: 'Billable weight (cwt)',
  Mileage: 'Mileage',
  BaselineLinehaulPrice: 'Baseline linehaul price',
  BaselineShorthaulPrice: 'Baseline shorthaul price',
  PriceEscalationFactor: 'Price escalation factor',
  TotalAmountRequested: 'Total amount requested',
  FuelSurchargePrice: 'Fuel surcharge price (per mi)',
  AdditionalDaySITPrice: 'Additional day SIT price',
  DaysInSIT: 'Days in SIT',
  ServiceSchedule: 'Service schedule',
  ServiceArea: 'Service area',
  RequestedPickup: 'Requested pickup',
  PackPrice: 'Pack price',
  [SERVICE_ITEM_PARAM_KEYS.ContractYearName]: 'Base year',
  [SERVICE_ITEM_PARAM_KEYS.WeightBilledActual]: 'Shipment weight',
  [SERVICE_ITEM_PARAM_KEYS.WeightActual]: 'Shipment weight',
  [SERVICE_ITEM_PARAM_KEYS.WeightEstimated]: 'Estimated',
  [SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]: 'Zip',
  [SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]: 'Zip',
  // Domestic non-peak or Domestic peak
  [SERVICE_ITEM_PARAM_KEYS.IsPeak]: 'Domestic',
  [SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]: 'Origin service area',
  [SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin]: 'Origin service schedule',
  [SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]: 'Requested pickup',
  [SERVICE_ITEM_PARAM_KEYS.ActualPickupDate]: 'Pickup date',
  [SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice]: 'EIA diesel',
  [SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier]: 'Weight-based distance multiplier',
  [SERVICE_ITEM_PARAM_KEYS.OriginPrice]: 'Origin price',
};

const SERVICE_ITEM_CODES = {
  DLH: 'DLH',
  DPK: 'DPK',
  FSC: 'FSC',
  DOP: 'DOP',
  DSH: 'DSH',
  DOASIT: 'DOASIT',
  DOFSIT: 'DOFSIT',
};

// TODO - temporary, will remove once all service item calculations are implemented
const allowedServiceItemCalculations = [
  SERVICE_ITEM_CODES.DLH,
  SERVICE_ITEM_CODES.FSC,
  SERVICE_ITEM_CODES.DPK,
  SERVICE_ITEM_CODES.DSH,
  SERVICE_ITEM_CODES.DOP,
  SERVICE_ITEM_CODES.DOASIT,
  SERVICE_ITEM_CODES.DOFSIT,
];

export {
  SERVICE_ITEM_STATUSES as default,
  SERVICE_ITEM_PARAM_KEYS,
  SERVICE_ITEM_CALCULATION_LABELS,
  SERVICE_ITEM_CODES,
  allowedServiceItemCalculations,
};
