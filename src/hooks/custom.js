import { useMemo } from 'react';

import { shipmentStatuses } from 'constants/shipments';
import { calculateShipmentNetWeight } from 'utils/shipmentWeights';

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
        ?.filter((s) => includedStatusesForCalculatingWeights(s.status))
        .reduce((prev, current) => {
          return prev + (calculateShipmentNetWeight(current) || 0);
        }, 0) || null
    );
  }, [mtoShipments]);
};

export const calculateEstimatedWeight = (mtoShipments) => {
  if (mtoShipments?.some((s) => includedStatusesForCalculatingWeights(s.status) && s.primeEstimatedWeight)) {
    return mtoShipments
      ?.filter((s) => includedStatusesForCalculatingWeights(s.status) && s.primeEstimatedWeight)
      .reduce((prev, current) => {
        return prev + current.primeEstimatedWeight;
      }, 0);
  }
  return null;
};

export const useCalculatedEstimatedWeight = (mtoShipments) => {
  return useMemo(() => {
    return calculateEstimatedWeight(mtoShipments);
  }, [mtoShipments]);
};
