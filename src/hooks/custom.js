import { useMemo } from 'react';

import { shipmentStatuses } from 'constants/shipments';
import returnLowestValue from 'utils/returnLowestValue';

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

export const useCalculatedWeightRequested = (mtoShipments) => {
  return useMemo(() => {
    return (
      mtoShipments
        ?.filter((s) => includedStatusesForCalculatingWeights(s.status) && (s.primeActualWeight || s.reweigh?.weight))
        .reduce((prev, current) => {
          return prev + returnLowestValue(current.primeActualWeight, current.reweigh?.weight);
        }, 0) || null
    );
  }, [mtoShipments]);
};

export const useCalculatedEstimatedWeight = (mtoShipments) => {
  return useMemo(() => {
    if (mtoShipments?.some((s) => includedStatusesForCalculatingWeights(s.status) && s.primeEstimatedWeight)) {
      return mtoShipments
        ?.filter((s) => includedStatusesForCalculatingWeights(s.status) && s.primeEstimatedWeight)
        .reduce((prev, current) => {
          return prev + current.primeEstimatedWeight;
        }, 0);
    }
    return null;
  }, [mtoShipments]);
};
