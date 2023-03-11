import { useMemo } from 'react';

import { shipmentStatuses } from 'constants/shipments';
import { calculateShipmentNetWeight, getShipmentEstimatedWeight } from 'utils/shipmentWeights';

// only sum estimated/actual/reweigh weights for shipments in these statuses
export const includedStatusesForCalculatingWeights = (status) => {
  return (
    status === shipmentStatuses.APPROVED ||
    status === shipmentStatuses.DIVERSION_REQUESTED ||
    status === shipmentStatuses.CANCELLATION_REQUESTED
  );
};

export const useCalculatedTotalBillableWeight = (mtoShipments) => {
  return useMemo(() => {
    return (
      mtoShipments
        ?.filter((s) => includedStatusesForCalculatingWeights(s.status) && s.calculatedBillableWeight)
        .reduce((prev, current) => {
          return prev + current.calculatedBillableWeight;
        }, 0) || null
    );
  }, [mtoShipments]);
};

export const calculateWeightRequested = (mtoShipments) => {
  if (mtoShipments?.some((s) => includedStatusesForCalculatingWeights(s.status) && calculateShipmentNetWeight(s))) {
    return (
      mtoShipments
        ?.filter((s) => includedStatusesForCalculatingWeights(s.status))
        .reduce((prev, current) => {
          return prev + (calculateShipmentNetWeight(current) || 0);
        }, 0) || null
    );
  }
  return null;
};

export const useCalculatedWeightRequested = (mtoShipments) => {
  return useMemo(() => {
    return calculateWeightRequested(mtoShipments);
  }, [mtoShipments]);
};

export const calculateEstimatedWeight = (mtoShipments) => {
  if (mtoShipments?.some((s) => includedStatusesForCalculatingWeights(s.status) && getShipmentEstimatedWeight(s))) {
    return mtoShipments
      ?.filter((s) => includedStatusesForCalculatingWeights(s.status) && getShipmentEstimatedWeight(s))
      .reduce((prev, current) => {
        return prev + getShipmentEstimatedWeight(current);
      }, 0);
  }
  return null;
};

export const useCalculatedEstimatedWeight = (mtoShipments) => {
  return useMemo(() => {
    return calculateEstimatedWeight(mtoShipments);
  }, [mtoShipments]);
};
