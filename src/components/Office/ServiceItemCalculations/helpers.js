import { SERVICE_ITEM_CALCULATION_LABELS, SERVICE_ITEM_CODES, SERVICE_ITEM_PARAM_KEYS } from 'constants/serviceItems';
import { formatWeight, formatCents, toDollarString } from 'shared/formatters';
import { formatDate } from 'shared/dates';
import { formatWeightCWTFromLbs, formatDollarFromMillicents } from 'utils/formatters';

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

const getPriceRateOrFactor = (params) => {
  return getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params) || '';
};

// billable weight calculation
const billableWeight = (params) => {
  const value = formatWeightCWTFromLbs(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightBilledActual, params));
  const label = SERVICE_ITEM_CALCULATION_LABELS.BillableWeight;

  const weightBilledActualDetail = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.WeightBilledActual]
  }: ${formatWeight(parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightBilledActual, params), 10))}`;

  const weightEstimated = getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightEstimated, params);
  const weightEstimatedDetail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.WeightEstimated]}: ${
    weightEstimated ? formatWeight(parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightEstimated, params), 10)) : ''
  }`;
  return calculation(value, label, weightBilledActualDetail, weightEstimatedDetail);
};

// mileage calculation
const mileageZIP3 = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZip3, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress, // take the zip 3
    params,
  )?.slice(2)} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipDestAddress,
    params,
  )?.slice(2)}`;

  return calculation(value, label, detail);
};

const mileageZip5 = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZip5, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress,
    params,
  )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipDestAddress,
    params,
  )}`;

  return calculation(value, label, detail);
};

const mileageZipSITOrigin = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZipSITOrigin, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipSITOriginHHGOriginalAddress,
    params,
  )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipSITOriginHHGActualAddress,
    params,
  )}`;

  return calculation(value, label, detail);
};

const baselineLinehaulPrice = (params) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.BaselineLinehaulPrice;
  const detail1 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;
  const detail2 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin,
    params,
  )}`;
  const detail3 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]}: ${formatDate(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params),
    'DD MMM YYYY',
  )}`;

  return calculation(value, label, detail1, detail2, detail3);
};

const baselineShorthaulPrice = (params) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.BaselineShorthaulPrice;
  const detail1 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;
  const detail2 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin,
    params,
  )}`;
  const detail3 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]}: ${formatDate(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params),
    'DD MMM YYYY',
  )}`;

  return calculation(value, label, detail1, detail2, detail3);
};

const dddSITmileageZip5 = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZipSITDest, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipDestAddress,
    params,
  )} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGFinalAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipSITDestHHGFinalAddress,
    params,
  )}`;

  return calculation(value, label, detail);
};

// There is no param representing the orgin price or destination price as available in the re_domestic_service_area_prices table
// A param to return the service schedule is also not being created
const originOrDestinationPrice = (params, isOrigin = true) => {
  const value = getPriceRateOrFactor(params);
  const serviceAreaKey = isOrigin ? SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin : SERVICE_ITEM_PARAM_KEYS.ServiceAreaDest;
  const serviceAreaVal = getParamValue(serviceAreaKey, params) ? getParamValue(serviceAreaKey, params) : '';
  const requestedPickupDateVal = getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params) || '';
  const label = isOrigin
    ? SERVICE_ITEM_CALCULATION_LABELS.OriginPrice
    : SERVICE_ITEM_CALCULATION_LABELS.DestinationPrice;

  const serviceArea = `${SERVICE_ITEM_CALCULATION_LABELS.ServiceArea}: ${serviceAreaVal}`;
  const requestedPickupDate = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]
  }: ${formatDate(requestedPickupDateVal, 'DD MMM YYYY')}`;

  const isPeak = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;

  return calculation(value, label, serviceArea, requestedPickupDate, isPeak);
};

const priceEscalationFactor = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params)
    ? getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params)
    : '';
  const label = SERVICE_ITEM_CALCULATION_LABELS.PriceEscalationFactor;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ContractYearName]}: ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ContractYearName, params) || ''
  }`;

  return calculation(value, label, detail);
};

const fuelSurchargePrice = (params) => {
  // to get the Fuel surcharge price (per mi), multiply FSCWeightBasedDistanceMultiplier by DistanceZip3
  // which gets the dollar value
  const value = parseFloat(
    String(
      getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier, params) *
        getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZip3, params),
    ),
  ).toFixed(2);
  const label = SERVICE_ITEM_CALCULATION_LABELS.FuelSurchargePrice;
  const detail1 = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice]
  }: ${formatDollarFromMillicents(getParamValue(SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice, params))}`;
  const detail2 = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier, params)}`;
  const detail3 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ActualPickupDate]}: ${formatDate(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.ActualPickupDate, params),
    'DD MMM YYYY',
  )}`;

  return calculation(value, label, detail1, detail2, detail3);
};

const packPrice = (params) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.PackPrice;
  const originServiceSchedule = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.ServicesScheduleOrigin, params)}`;
  const requestedPickup = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]
  }: ${formatDate(getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params), 'DD MMM YYYY')}`;
  const domesticNonPeak = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;

  return calculation(value, label, originServiceSchedule, requestedPickup, domesticNonPeak);
};

const unpackPrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.UnpackPrice;
  const destServiceSchedule = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServicesScheduleDest]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.ServicesScheduleDest, params)}`;
  const requestedPickup = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]
  }: ${formatDate(getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params), 'DD MMM YYYY')}`;
  const domesticNonPeak = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;

  return calculation(value, label, destServiceSchedule, requestedPickup, domesticNonPeak);
};

const additionalDaySITPrice = (params) => {
  const value = getPriceRateOrFactor(params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.AdditionalDaySITPrice;
  const serviceArea = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ServiceAreaOrigin,
    params,
  )}`;
  const requestedPickupDate = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]
  }: ${formatDate(getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params), 'DD MMM YYYY')}`;
  const peak = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;

  return calculation(value, label, serviceArea, requestedPickupDate, peak);
};

const sitDeliveryPrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.SITDeliveryPrice;
  const sitScheduleDestination = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.SITScheduleDest]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.SITScheduleDest, params)}`;
  const requestedPickup = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]
  }: ${formatDate(getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params), 'DD MMM YYYY')}`;
  const domesticNonPeak = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;

  return calculation(value, label, sitScheduleDestination, requestedPickup, domesticNonPeak);
};

const daysInSIT = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.NumberDaysSIT, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.DaysInSIT;

  return calculation(value, label);
};

const pickupSITPrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.PickupSITPrice;

  const originSITSchedule = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.SITScheduleOrigin]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.SITScheduleOrigin, params)}`;

  const requestedPickupDate = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]
  }: ${formatDate(getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params), 'DD MMM YYYY')}`;

  const peak = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.IsPeak]} ${
    getParamValue(SERVICE_ITEM_PARAM_KEYS.IsPeak, params)?.toLowerCase() === 'true' ? 'peak' : 'non-peak'
  }`;

  return calculation(value, label, originSITSchedule, requestedPickupDate, peak);
};

// totalAmountRequested is not a service item param
const totalAmountRequested = (totalAmount) => {
  const value = toDollarString(formatCents(totalAmount));
  const label = SERVICE_ITEM_CALCULATION_LABELS.TotalAmountRequested;
  const detail = '';

  return calculation(value, label, detail);
};

const makeCalculations = (itemCode, totalAmount, params) => {
  let result = [];

  switch (itemCode) {
    case SERVICE_ITEM_CODES.DDDSIT:
      result = [
        billableWeight(params),
        dddSITmileageZip5(params),
        sitDeliveryPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic longhaul
    case SERVICE_ITEM_CODES.DLH:
      result = [
        billableWeight(params),
        mileageZIP3(params),
        baselineLinehaulPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Fuel surcharge
    case SERVICE_ITEM_CODES.FSC:
      result = [
        billableWeight(params),
        mileageZIP3(params),
        fuelSurchargePrice(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin price
    case SERVICE_ITEM_CODES.DOP:
      result = [
        billableWeight(params),
        originOrDestinationPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin 1st day SIT
    case SERVICE_ITEM_CODES.DOFSIT:
      result = [
        billableWeight(params),
        originOrDestinationPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic destination 1st day SIT
    case SERVICE_ITEM_CODES.DDFSIT:
      result = [
        billableWeight(params),
        originOrDestinationPrice(params, false),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic packing
    case SERVICE_ITEM_CODES.DPK:
      result = [
        billableWeight(params),
        packPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic shorthaul
    case SERVICE_ITEM_CODES.DSH:
      result = [
        billableWeight(params),
        mileageZip5(params),
        baselineShorthaulPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic destination
    case SERVICE_ITEM_CODES.DDP:
      result = [
        billableWeight(params),
        originOrDestinationPrice(params, false),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    // Domestic origin additional SIT
    case SERVICE_ITEM_CODES.DOASIT:
      result = [
        billableWeight(params),
        daysInSIT(params),
        additionalDaySITPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    case SERVICE_ITEM_CODES.DOPSIT:
      result = [
        billableWeight(params),
        mileageZipSITOrigin(params),
        pickupSITPrice(params),
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

    default:
      break;
  }
  return result;
};

export { makeCalculations as default, makeCalculations };
