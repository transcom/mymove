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

/**
 * This function calculates the total Billable Weight of the move,
 * by adding up all of the calculatedBillableWeight fields of all shipments with the required statuses.
 *
 * This function does **NOT** include PPM net weights in the calculation.
 * @param mtoShipments An array of MTO Shipments
 * @return {int|null} The calculated total billable weight
 */
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

/**
 * This function calculates the weight requested of a move,
 * by adding up all of the net weights of all shipments with the required statuses.
 *
 * This function includes PPM net weights in its calculation. In order to calculate the PPM net weights,
 * the corresponding weight tickets must be attached to the PPM shipments.
 * @see useAddWeightTicketsToPPMShipments in hooks/queries for information on adding weight tickets to PPM shipments
 * @param mtoShipments An array of MTO Shipments
 * @return {int|null} The total weight requested
 */
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
