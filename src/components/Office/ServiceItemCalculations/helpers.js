import { SERVICE_ITEM_CALCULATION_LABELS, SERVICE_ITEM_CODES, SERVICE_ITEM_PARAM_KEYS } from 'constants/serviceItems';
import { LONGHAUL_MIN_DISTANCE } from 'constants/shipments';
import { formatDateWithUTC } from 'shared/dates';
import {
  formatCents,
  formatWeight,
  formatWeightCWTFromLbs,
  formatDollarFromMillicents,
  toDollarString,
  formatDistanceUnitMiles,
} from 'utils/formatters';
import { SHIPMENT_OPTIONS } from 'shared/constants';

const calculation = (value, label, ...details) => {
  return {
    value,
    label,
    details: [...details],
  };
};

const getParamValue = (key, params) => {
  return params?.find((param) => param?.key === key)?.value;
};

const peak = (params) => {
  return `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;
};

const serviceAreaOrigin = (params) => {
  return `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin,
    params,
  )}`;
};

const sitServiceAreaOrigin = (params) => {
  return `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.SITServiceAreaOrigin]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.SITServiceAreaOrigin,
    params,
  )}`;
};

const serviceAreaDest = (params) => {
  return `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServiceAreaDest]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServiceAreaDest,
    params,
  )}`;
};

const requestedPickupDateLabel = (shipmentType) => {
  switch (shipmentType) {
    case SHIPMENT_OPTIONS.NTSR:
      return SERVICE_ITEM_CALCULATION_LABELS.NTSReleaseReferenceDate;

    default:
      return SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ReferenceDate];
  }
};

const referenceDate = (params, shipmentType) => {
  const label = requestedPickupDateLabel(shipmentType);
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.ReferenceDate, params);
  return `${label}: ${formatDateWithUTC(value, 'DD MMM YYYY')}`;
};

const cratingDate = (params) => {
  return `${SERVICE_ITEM_CALCULATION_LABELS.CratingDate}: ${formatDateWithUTC(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ReferenceDate, params),
    'DD MMM YYYY',
  )}`;
};

const unCratingDate = (params) => {
  return `${SERVICE_ITEM_CALCULATION_LABELS.UncratingDate}: ${formatDateWithUTC(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ReferenceDate, params),
    'DD MMM YYYY',
  )}`;
};

const getPriceRateOrFactor = (params) => {
  return getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params) || '';
};

const formatDetail = (detail, styles = {}) => {
  return {
    text: detail,
    styles,
  };
};

const formatMileage = (detail) => {
  return formatDistanceUnitMiles(detail, false);
};

// billable weight calculation
const formatWeightFromParams = (params, key) => {
  return formatWeight(parseInt(getParamValue(key, params), 10));
};

const formatWeightDetailText = (params, key) => {
  const value = getParamValue(key, params);
  const paramValue = value ? formatWeightFromParams(params, key) : '';
  const detailText = `${SERVICE_ITEM_CALCULATION_LABELS[key]}: ${paramValue}`;
  return paramValue ? detailText : '';
};

const billableWeight = (params) => {
  const value = formatWeightCWTFromLbs(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightBilled, params));
  const label = SERVICE_ITEM_CALCULATION_LABELS.BillableWeight;

  const details = [];
  const boldStyles = { fontWeight: 'bold' };

  const weightAdjustedDetail = formatWeightDetailText(params, SERVICE_ITEM_PARAM_KEYS.WeightAdjusted);
  if (weightAdjustedDetail) {
    // The weight adjusted detail should always be bolded
    details.push(formatDetail(weightAdjustedDetail, boldStyles));
  }

  const weightReweighDetail = formatWeightDetailText(params, SERVICE_ITEM_PARAM_KEYS.WeightReweigh);
  const weightOriginalDetail = formatWeightDetailText(params, SERVICE_ITEM_PARAM_KEYS.WeightOriginal);

  // If the reweigh weight exists, figure out if the reweigh or the original weight should be bolded.
  if (weightReweighDetail && weightOriginalDetail) {
    const weightReweighValue = parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightReweigh, params), 10);
    const weightOriginalValue = parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightOriginal, params), 10);

    let reweighStyles = {};
    let originalStyles = {};

    // Only bold the weights if there is no adjusted weight
    if (!weightAdjustedDetail) {
      // if the reweigh weight matches the original weight, bold the reweigh weight
      if (weightReweighValue <= weightOriginalValue) {
        reweighStyles = boldStyles;
      } else {
        originalStyles = boldStyles;
      }
    }

    details.push(formatDetail(weightReweighDetail, reweighStyles));
    details.push(formatDetail(weightOriginalDetail, originalStyles));
  } else if (weightOriginalDetail) {
    // Otherwise, always have the original weight as bolded if there is no weight adjusted.
    let originalStyles = {};
    if (!weightAdjustedDetail) {
      originalStyles = boldStyles;
    }
    details.push(formatDetail(weightOriginalDetail, originalStyles));
  }

  const weightEstimatedDetail = formatWeightDetailText(params, SERVICE_ITEM_PARAM_KEYS.WeightEstimated);
  if (weightEstimatedDetail) {
    details.push(formatDetail(weightEstimatedDetail));
  }

  const fscWeightBasedDistanceMultiplier = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier, params)}`;
  if (fscWeightBasedDistanceMultiplier.length > 2) {
    details.push(formatDetail(fscWeightBasedDistanceMultiplier));
  }

  return calculation(value, label, ...details);
};

const shuttleBillableWeight = (params) => {
  const value = formatWeightCWTFromLbs(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightBilled, params));
  const label = SERVICE_ITEM_CALCULATION_LABELS.BillableWeight;
  const weightReweighValue = parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightReweigh, params), 10);
  const weightOriginalValue = parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightOriginal, params), 10);
  let lowestActualWeight;
  // In order to grab the lower of the two integers, we need to make sure that both are in fact numbers first
  // If NaN comes back from parseInt, we know to use the other value
  if (!Number.isNaN(weightReweighValue) && !Number.isNaN(weightOriginalValue)) {
    lowestActualWeight = Math.min(weightReweighValue, weightOriginalValue);
  } else if (!Number.isNaN(weightReweighValue)) {
    lowestActualWeight = weightReweighValue;
  } else if (!Number.isNaN(weightOriginalValue)) {
    lowestActualWeight = weightOriginalValue;
  }
  const weightBilledDetail = `${SERVICE_ITEM_CALCULATION_LABELS.ShuttleWeight}: ${formatWeight(lowestActualWeight)}`;

  const weightEstimated = getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightEstimated, params);
  const weightEstimatedDetail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.WeightEstimated]}: ${
    weightEstimated ? formatWeight(parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightEstimated, params), 10)) : ''
  }`;
  return calculation(value, label, formatDetail(weightBilledDetail), formatDetail(weightEstimatedDetail));
};

const mileageZip = (params) => {
  const value = `${formatMileage(getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZip, params))}`;
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress,
    params,
  )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipDestAddress,
    params,
  )}`;

  return calculation(value, label, formatDetail(detail));
};

const mileageZipSIT = (params, itemCode) => {
  let label;
  let distanceZip;
  let detail;

  switch (itemCode) {
    case SERVICE_ITEM_CODES.DOSFSC:
      label = SERVICE_ITEM_CALCULATION_LABELS.MileageIntoSIT;
      distanceZip = SERVICE_ITEM_PARAM_KEYS.DistanceZipSITOrigin;
      detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
        SERVICE_ITEM_PARAM_KEYS.ZipSITOriginHHGOriginalAddress,
        params,
      )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
        SERVICE_ITEM_PARAM_KEYS.ZipSITOriginHHGActualAddress,
        params,
      )}`;
      break;

    case SERVICE_ITEM_CODES.DDSFSC:
      label = SERVICE_ITEM_CALCULATION_LABELS.MileageOutOfSIT;
      distanceZip = SERVICE_ITEM_PARAM_KEYS.DistanceZipSITDest;
      detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
        SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGOriginalAddress,
        params,
      )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
        SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGFinalAddress,
        params,
      )}`;
      break;

    default:
      label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
      distanceZip = SERVICE_ITEM_PARAM_KEYS.DistanceZipSITOrigin;
      detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
        SERVICE_ITEM_PARAM_KEYS.ZipSITOriginHHGOriginalAddress,
        params,
      )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
        SERVICE_ITEM_PARAM_KEYS.ZipSITOriginHHGActualAddress,
        params,
      )}`;
  }

  const value = formatMileage(getParamValue(distanceZip, params));

  return calculation(value, label, formatDetail(detail));
};

const baselineLinehaulPrice = (params, shipmentType) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.BaselineLinehaulPrice;

  return calculation(
    value,
    label,
    formatDetail(peak(params)),
    formatDetail(serviceAreaOrigin(params)),
    formatDetail(referenceDate(params, shipmentType)),
  );
};

const baselineShorthaulPrice = (params, shipmentType) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.BaselineShorthaulPrice;

  return calculation(
    value,
    label,
    formatDetail(peak(params)),
    formatDetail(serviceAreaOrigin(params)),
    formatDetail(referenceDate(params, shipmentType)),
  );
};
const dddSITmileageZip5 = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZipSITDest, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGOriginalAddress]
  } ${getParamValue(SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGOriginalAddress, params)} to ${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGFinalAddress]
  } ${getParamValue(SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGFinalAddress, params)}`;

  return calculation(value, label, formatDetail(detail));
};

// There is no param representing the orgin price as available in the re_domestic_service_area_prices table
// A param to return the service schedule is also not being created
const originPrice = (params, shipmentType, serviceCode) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.OriginPrice;

  // First day origin sit utilizes a SIT specific service area origin service param
  let serviceArea = serviceAreaOrigin(params);
  if (serviceCode === SERVICE_ITEM_CODES.DOFSIT) {
    serviceArea = sitServiceAreaOrigin(params);
  }
  return calculation(
    value,
    label,
    formatDetail(serviceArea),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const shuttleOriginPriceDomestic = (params) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.OriginPrice;

  const serviceSchedule = `${SERVICE_ITEM_CALCULATION_LABELS.ServiceSchedule}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin,
    params,
  )}`;

  const pickupDate = `${SERVICE_ITEM_CALCULATION_LABELS.PickupDate}: ${formatDateWithUTC(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ReferenceDate, params),
    'DD MMM YYYY',
  )}`;

  return calculation(
    value,
    label,
    formatDetail(serviceSchedule),
    formatDetail(pickupDate),
    formatDetail(SERVICE_ITEM_CALCULATION_LABELS.Domestic),
  );
};

// There is no param representing the destination price as available in the re_domestic_service_area_prices table
// A param to return the service schedule is also not being created
const destinationPrice = (params, shipmentType) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.DestinationPrice;

  return calculation(
    value,
    label,
    formatDetail(serviceAreaDest(params)),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const shuttleDestinationPriceDomestic = (params) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.DestinationPrice;

  const serviceSchedule = `${SERVICE_ITEM_CALCULATION_LABELS.ServiceSchedule}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServicesScheduleDest,
    params,
  )}`;

  const deliveryDate = `${SERVICE_ITEM_CALCULATION_LABELS.DeliveryDate}: ${formatDateWithUTC(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ReferenceDate, params),
    'DD MMM YYYY',
  )}`;

  return calculation(
    value,
    label,
    formatDetail(serviceSchedule),
    formatDetail(deliveryDate),
    formatDetail(SERVICE_ITEM_CALCULATION_LABELS.Domestic),
  );
};

const priceEscalationFactor = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params)
    ? getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params)
    : '';
  const label = SERVICE_ITEM_CALCULATION_LABELS.PriceEscalationFactor;

  const contractYearName = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ContractYearName]}: ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ContractYearName, params) || ''
  }`;

  return calculation(value, label, formatDetail(contractYearName));
};

const priceEscalationFactorWithoutContractYear = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params)
    ? getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params)
    : '';
  const label = SERVICE_ITEM_CALCULATION_LABELS.PriceEscalationFactor;

  return calculation(value, label);
};

const fuelSurchargePrice = (params, itemCode) => {
  // to get the Fuel surcharge price (per mi), multiply FSCWeightBasedDistanceMultiplier by distanceZip
  // which gets the value in Cents to the tenths decimal place
  let distanceZip;
  switch (itemCode) {
    case SERVICE_ITEM_CODES.DDSFSC:
      distanceZip = SERVICE_ITEM_PARAM_KEYS.DistanceZipSITDest;
      break;
    case SERVICE_ITEM_CODES.DOSFSC:
      distanceZip = SERVICE_ITEM_PARAM_KEYS.DistanceZipSITOrigin;
      break;
    default:
      distanceZip = SERVICE_ITEM_PARAM_KEYS.DistanceZip;
  }
  const value = parseFloat(
    String(
      getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier, params) *
        getParamValue(distanceZip, params),
    ),
  ).toFixed(1);
  const label =
    itemCode === SERVICE_ITEM_CODES.DOSFSC || itemCode === SERVICE_ITEM_CODES.DDSFSC
      ? SERVICE_ITEM_CALCULATION_LABELS.SITFuelSurchargePrice
      : SERVICE_ITEM_CALCULATION_LABELS.FuelSurchargePrice;

  const eiaFuelPrice = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice]
  }: ${formatDollarFromMillicents(getParamValue(SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice, params))}`;

  const fuelRateAdjustment = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.FSCPriceDifferenceInCents]
  }: ${formatCents(getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCPriceDifferenceInCents, params), 1, 1)}`;

  const fscWeightBasedDistanceMultiplier = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier, params)}`;

  const actualPickupDate = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ActualPickupDate]
  }: ${formatDateWithUTC(getParamValue(SERVICE_ITEM_PARAM_KEYS.ActualPickupDate, params), 'DD MMM YYYY')}`;

  const eiaFuelPrice = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice]
  }: ${formatDollarFromMillicents(getParamValue(SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice, params), 3)}`;

  const baselineRateDifference = `${SERVICE_ITEM_CALCULATION_LABELS.FSCPriceDifferenceInCents}: ${formatCents(
    parseFloat(getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCPriceDifferenceInCents, params)),
    1,
    1,
  )} \u00A2`;

  return calculation(
    value,
    label,
    formatDetail(actualPickupDate),
    formatDetail(eiaFuelPrice),
    formatDetail(baselineRateDifference),
  );
};

const packPrice = (params, shipmentType) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.PackPrice;
  const originServiceSchedule = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin, params)}`;

  return calculation(
    value,
    label,
    formatDetail(originServiceSchedule),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const ntsPackingFactor = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.NTSPackingFactor, params) || '';
  const label = SERVICE_ITEM_CALCULATION_LABELS.NTSPackingFactor;

  return calculation(value, label);
};

const unpackPrice = (params, shipmentType) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.UnpackPrice;
  const destServiceSchedule = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServicesScheduleDest]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.ServicesScheduleDest, params)}`;

  return calculation(
    value,
    label,
    formatDetail(destServiceSchedule),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const additionalDayOriginSITPrice = (params, shipmentType) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.AdditionalDaySITPrice;

  return calculation(
    value,
    label,
    formatDetail(sitServiceAreaOrigin(params)),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const additionalDayDestinationSITPrice = (params, shipmentType) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.AdditionalDaySITPrice;

  return calculation(
    value,
    label,
    formatDetail(serviceAreaDest(params)),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const sitDeliveryPrice = (params, shipmentType) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.SITDeliveryPrice;

  const sitScheduleDestination = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.SITScheduleDest]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.SITScheduleDest, params)}`;

  return calculation(
    value,
    label,
    formatDetail(sitScheduleDestination),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const sitDeliveryPriceShorthaulDifferentZIP3 = (params, shipmentType) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.SITDeliveryPrice;

  const sitScheduleDestination = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.SITScheduleDest]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.SITScheduleDest, params)}`;

  return calculation(
    value,
    label,
    formatDetail(sitScheduleDestination),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
    formatDetail('<=50 miles'),
  );
};

const daysInSIT = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.NumberDaysSIT, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.DaysInSIT;

  return calculation(value, label);
};

const pickupSITPrice = (params, shipmentType) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.PickupSITPrice;

  const originSITSchedule = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.SITScheduleOrigin]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.SITScheduleOrigin, params)}`;

  return calculation(
    value,
    label,
    formatDetail(originSITSchedule),
    formatDetail(referenceDate(params, shipmentType)),
    formatDetail(peak(params)),
  );
};

const cratingPrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.CratingPrice;

  const serviceSchedule = `${SERVICE_ITEM_CALCULATION_LABELS.ServiceSchedule}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin,
    params,
  )}`;

  return calculation(
    value,
    label,
    formatDetail(serviceSchedule),
    formatDetail(cratingDate(params)),
    formatDetail(SERVICE_ITEM_CALCULATION_LABELS.Domestic),
  );
};

const unCratingPrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.UncratingPrice;

  const serviceSchedule = `${SERVICE_ITEM_CALCULATION_LABELS.ServiceSchedule}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServicesScheduleDest,
    params,
  )}`;

  return calculation(
    value,
    label,
    formatDetail(serviceSchedule),
    formatDetail(unCratingDate(params)),
    formatDetail(SERVICE_ITEM_CALCULATION_LABELS.Domestic),
  );
};

const cratingSize = (params, mtoParams) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.CubicFeetBilled, params);
  const length = getParamValue(SERVICE_ITEM_PARAM_KEYS.DimensionLength, params);
  const height = getParamValue(SERVICE_ITEM_PARAM_KEYS.DimensionHeight, params);
  const width = getParamValue(SERVICE_ITEM_PARAM_KEYS.DimensionWidth, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.CubicFeetBilled;

  const description = `${SERVICE_ITEM_CALCULATION_LABELS.Description}: ${mtoParams.description}`;

  const formattedDimensions = `${SERVICE_ITEM_CALCULATION_LABELS.Dimensions}: ${length}x${width}x${height} in`;

  return calculation(value, label, formatDetail(description), formatDetail(formattedDimensions));
};

// totalAmountRequested is not a service item param
const totalAmountRequested = (totalAmount) => {
  const value = toDollarString(formatCents(totalAmount));
  const label = SERVICE_ITEM_CALCULATION_LABELS.FuelRateAdjustment;
  const detail = '';

  return calculation(value, label, formatDetail(detail));
};

export default function makeCalculations(itemCode, totalAmount, params, mtoParams, shipmentType) {
  let result = [];
  switch (itemCode) {
    case SERVICE_ITEM_CODES.DDDSIT: {
      const mileage = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZipSITDest, params);
      const startZip = getParamValue(SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGOriginalAddress, params)?.slice(0, 3);
      const endZip = getParamValue(SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGFinalAddress, params)?.slice(0, 3);
      // Mileage does not factor into the pricing for distances less than 50 miles and non-matching
      // zip3, so we won't display mileage
      if (mileage <= LONGHAUL_MIN_DISTANCE && startZip !== endZip) {
        result = [
          billableWeight(params),
          sitDeliveryPriceShorthaulDifferentZIP3(params, shipmentType), // Display under mileage threshold
          priceEscalationFactor(params),
          totalAmountRequested(totalAmount),
        ];
      } else {
        result = [
          billableWeight(params),
          dddSITmileageZip5(params),
          sitDeliveryPrice(params, shipmentType),
          priceEscalationFactor(params),
          totalAmountRequested(totalAmount),
        ];
      }
      break;
    }
    // Domestic longhaul
    case SERVICE_ITEM_CODES.DLH:
      result = [
        billableWeight(params),
        mileageZip(params),
        baselineLinehaulPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Fuel surcharge
    case SERVICE_ITEM_CODES.FSC:
      result = [
        billableWeight(params),
        mileageZip(params),
        mileageFactor(params, itemCode),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin SIT fuel surcharge
    case SERVICE_ITEM_CODES.DOSFSC:
      result = [
        billableWeight(params),
        mileageZipSIT(params, itemCode),
        mileageFactor(params, itemCode),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic destination SIT fuel surcharge
    case SERVICE_ITEM_CODES.DDSFSC:
      result = [
        billableWeight(params),
        mileageZipSIT(params, itemCode),
        mileageFactor(params, itemCode),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin price
    case SERVICE_ITEM_CODES.DOP:
      result = [
        billableWeight(params),
        originPrice(params, shipmentType, SERVICE_ITEM_CODES.DOP),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin 1st day SIT
    case SERVICE_ITEM_CODES.DOFSIT:
      result = [
        billableWeight(params),
        originPrice(params, shipmentType, SERVICE_ITEM_CODES.DOFSIT),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic destination 1st day SIT
    case SERVICE_ITEM_CODES.DDFSIT:
      result = [
        billableWeight(params),
        destinationPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic packing
    case SERVICE_ITEM_CODES.DPK:
      result = [
        billableWeight(params),
        packPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic NTS packing
    case SERVICE_ITEM_CODES.DNPK:
      result = [
        billableWeight(params),
        packPrice(params, shipmentType),
        ntsPackingFactor(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic shorthaul
    case SERVICE_ITEM_CODES.DSH:
      result = [
        billableWeight(params),
        mileageZip(params),
        baselineShorthaulPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic destination
    case SERVICE_ITEM_CODES.DDP:
      result = [
        billableWeight(params),
        destinationPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin additional SIT
    case SERVICE_ITEM_CODES.DOASIT:
      result = [
        billableWeight(params),
        daysInSIT(params),
        additionalDayOriginSITPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic Origin SIT Pickup
    case SERVICE_ITEM_CODES.DOPSIT:
      result = [
        billableWeight(params),
        mileageZipSIT(params, itemCode),
        pickupSITPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin shuttle service
    case SERVICE_ITEM_CODES.DOSHUT:
      result = [
        shuttleBillableWeight(params),
        shuttleOriginPriceDomestic(params),
        priceEscalationFactorWithoutContractYear(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic Destination Additional Days SIT
    case SERVICE_ITEM_CODES.DDASIT:
      result = [
        billableWeight(params),
        daysInSIT(params),
        additionalDayDestinationSITPrice(params, shipmentType),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic unpacking
    case SERVICE_ITEM_CODES.DUPK:
      result = [
        billableWeight(params),
        unpackPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic destination shuttle service
    case SERVICE_ITEM_CODES.DDSHUT:
      result = [
        shuttleBillableWeight(params),
        shuttleDestinationPriceDomestic(params),
        priceEscalationFactorWithoutContractYear(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic crating
    case SERVICE_ITEM_CODES.DCRT:
      result = [
        cratingSize(params, mtoParams),
        cratingPrice(params),
        priceEscalationFactorWithoutContractYear(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic uncrating
    case SERVICE_ITEM_CODES.DUCRT:
      result = [
        cratingSize(params, mtoParams),
        unCratingPrice(params),
        priceEscalationFactorWithoutContractYear(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    default:
      break;
  }
  return result;
}

export { makeCalculations };
