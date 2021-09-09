import { useMemo } from 'react';

import { shipmentStatuses } from 'constants/shipments';
import returnLowestValue from 'utils/returnLowestValue';

export const includedStatuses = (status) => {
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
        ?.filter((s) => includedStatuses(s.status) && s.calculatedBillableWeight)
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
        ?.filter((s) => includedStatuses(s.status) && (s.primeActualWeight || s.reweigh?.weight))
        .reduce((prev, current) => {
          return prev + returnLowestValue(current.primeActualWeight, current.reweigh?.weight);
        }, 0) || null
    );
  }, [mtoShipments]);
};
