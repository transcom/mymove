import { SERVICE_ITEM_CALCULATION_LABELS, SERVICE_ITEM_CODES, SERVICE_ITEM_PARAM_KEYS } from 'constants/serviceItems';
import { formatWeight, formatCents, toDollarString } from 'shared/formatters';
import { formatDate } from 'shared/dates';
import { formatWeightCWTFromLbs } from 'utils/formatters';

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

// billable weight calculation
const billableWeight = (params) => {
  const value = formatWeightCWTFromLbs(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightActual, params));
  const label = SERVICE_ITEM_CALCULATION_LABELS.BillableWeight;
  const detail1 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.WeightActual]}: ${formatWeight(
    parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightActual, params), 10),
  )}`;
  const detail2 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.WeightEstimated]}: ${formatWeight(
    parseInt(getParamValue(SERVICE_ITEM_PARAM_KEYS.WeightEstimated, params), 10),
  )}`;

  return calculation(value, label, detail1, detail2);
};

// mileage calculation
const mileage = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.DistanceZip3, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.Mileage;
  const detail = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipPickupAddress, // take the zip 3
    params,
  ).slice(2)} to ${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.ZipDestAddress]} ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.ZipDestAddress,
    params,
  ).slice(2)}`;

  return calculation(value, label, detail);
};

const baselineLinehaulPrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.PriceRateOrFactor, params);
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

const priceEscalationFactor = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.EscalationCompounded, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.PriceEscalationFactor;
  const detail = '';

  return calculation(value, label, detail);
};

const fuelSurchargePrice = (params) => {
  const value = getParamValue(SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice, params);
  const label = SERVICE_ITEM_CALCULATION_LABELS.FuelSurchargePrice;
  const detail1 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice]}: ${getParamValue(
    SERVICE_ITEM_PARAM_KEYS.EIAFuelPrice,
    params,
  )}`;
  const detail2 = `${
    SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier]
  }: ${getParamValue(SERVICE_ITEM_PARAM_KEYS.FSCWeightBasedDistanceMultiplier, params)}`;
  const detail3 = `${SERVICE_ITEM_CALCULATION_LABELS[SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate]}: ${formatDate(
    getParamValue(SERVICE_ITEM_PARAM_KEYS.RequestedPickupDate, params),
    'DD MMM YYYY',
  )}`;

  return calculation(value, label, detail1, detail2, detail3);
};

// totalAmountRequested is not a service item param
const totalAmountRequested = (totalAmount) => {
  const value = toDollarString(formatCents(totalAmount));
  const label = SERVICE_ITEM_CALCULATION_LABELS.TotalAmountRequested;
  const detail = '';

  return calculation(value, label, detail);
};

const makeCalculations = (itemCode, totalAmount, params) => {
  let result;

  switch (itemCode) {
    case SERVICE_ITEM_CODES.DLH:
      result = [
        billableWeight(params),
        mileage(params),
        baselineLinehaulPrice(params),
        priceEscalationFactor(params),
        totalAmountRequested(totalAmount),
      ];
      break;
    case SERVICE_ITEM_CODES.FSC:
      result = [billableWeight(params), mileage(params), fuelSurchargePrice(params), totalAmountRequested(totalAmount)];
      break;
    default:
      break;
  }
  return result;
};

export { makeCalculations as default, makeCalculations };
