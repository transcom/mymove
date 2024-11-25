import { useEffect, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { capitalize } from 'lodash';

import { SHIPMENT_OPTIONS, isAdminSite, isMilmoveSite, isOfficeSite } from '../shared/constants';

import { ADMIN_BASE_PAGE_TITLE, MILMOVE_BASE_PAGE_TITLE, OFFICE_BASE_PAGE_TITLE } from 'constants/titles';
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

const addressesMatch = (address1, address2) => {
  // Null or undefined check. This resolves I-12397
  if (!address1 || !address2) {
    return false;
  }

  return (
    address1.city === address2.city &&
    address1.postalCode === address2.postalCode &&
    address1.state === address2.state &&
    address1.streetAddress1 === address2.streetAddress1
  );
};

// This allows us to take all shipments and identify the next one in the diversion chain easily
// A child diversion's pickup address is the delivery address of the parent
const findChildShipmentByAddress = (currentShipment, allShipments) => {
  // Find a shipment whose pickup address matches the current shipment's delivery address
  return allShipments.find(
    (shipment) => addressesMatch(shipment.pickupAddress, currentShipment.destinationAddress) && shipment.diversion,
  );
};

// This allows us to group diverted shipments
// For example, if there are 2 shipments in an order that the
// TOO decided to mark for diversion, each of these shipments will now have
// their own child shipments. This allows us to split them up into their
// respective chains for processing.
const groupDivertedShipmentsByAddress = (shipments) => {
  const chains = [];
  const remainingUnchainedShipments = new Set(shipments.map((s) => s.id));
  const shipmentMap = new Map(shipments.map((s) => [s.id, s]));

  // Create each chain
  shipments.forEach((shipment) => {
    if (shipment.diversion && remainingUnchainedShipments.has(shipment.id)) {
      const chain = [];
      let currentShipment = shipment;
      // Loop over the shipments inside of the remainingUnchainedShipments
      // Keep identifying the next child shipment in the chain and pushing it accordingly
      // Stop looping when it can no longer find child shipment by address or if
      // the shipment found is not incide of the remianing unchained shipments
      while (currentShipment && remainingUnchainedShipments.has(currentShipment.id)) {
        chain.push(currentShipment);
        remainingUnchainedShipments.delete(currentShipment.id);
        currentShipment = findChildShipmentByAddress(currentShipment, Array.from(shipmentMap.values()));
      }
      if (chain.length > 0) {
        chains.push(chain);
      }
    }
  });
  return chains;
};

const getEstimatedLowestShipmentWeight = (shipments) => {
  return shipments.reduce((lowest, shipment) => {
    const estimatedWeight = getShipmentEstimatedWeight(shipment);
    return estimatedWeight < lowest ? estimatedWeight : lowest;
  }, Number.MAX_SAFE_INTEGER);
};

const getLowestShipmentNetWeight = (shipments) => {
  return shipments.reduce((lowest, shipment) => {
    const currentNetWeight = calculateShipmentNetWeight(shipment);
    return currentNetWeight < lowest ? currentNetWeight : lowest;
  }, Number.MAX_SAFE_INTEGER);
};

/**
 * This function calculates the total Billable Weight of the move,
 * by adding up all of the calculatedBillableWeight fields of all shipments with the required statuses.
 * It has unique calculations to also only count the lowest weight from a diverted shipment "Chain".
 * It is chained by a shipment having its diversion parameter set to true and the delivery address
 * of the parent shipment matching the pickup address of the child shipment.
 *
 * This function does **NOT** include PPM net weights in the calculation.
 * @param mtoShipments An array of MTO Shipments
 * @return {int|null} The calculated total billable weight
 */
export const useCalculatedTotalBillableWeight = (mtoShipments, weightAdjustment = 1.0) => {
  return useMemo(() => {
    if (mtoShipments?.length) {
      // Separate diverted shipments and other eligible shipments for weight calculations
      // This is done because a diverted shipment only has one true weight, but when it gets diverted
      // it is entered as a whole new shipment. This causes the sum to be counted twice for its weight,
      // we filter to include only the lowest weight from the diverted shipments here to prevent that.
      const divertedEligibleShipments = mtoShipments.filter(
        (s) => s.diversion && includedStatusesForCalculatingWeights(s.status) && s.calculatedBillableWeight,
      );
      const otherEligibleShipments = mtoShipments.filter(
        (s) => !s.diversion && includedStatusesForCalculatingWeights(s.status) && s.calculatedBillableWeight,
      );
      // In order to properly sum the lowest weight of the diverted shipments, we must first put them into
      // their correct "chains". Please see comments for groupDivertedShipments for more details.
      const chains = groupDivertedShipmentsByAddress(divertedEligibleShipments);
      // Grab the lowest weight from each chain
      const chainWeights = chains.map((chain) =>
        chain.reduce((lowest, shipment) => {
          return shipment.calculatedBillableWeight < lowest ? shipment.calculatedBillableWeight : lowest;
        }, Number.MAX_SAFE_INTEGER),
      );
      // Now that we have the lowest weight from each chain, get the sum
      const sumChainWeights = chainWeights.reduce((total, weight) => total + weight, 0);

      // Sum non-diverted eligible billable weights
      const sumOtherEligibleWeights = otherEligibleShipments.reduce((total, current) => {
        let currentWeight =
          current.calculatedBillableWeight < current.primeEstimatedWeight * weightAdjustment
            ? current.calculatedBillableWeight
            : current.primeEstimatedWeight * weightAdjustment;

        if (current.shipmentType === SHIPMENT_OPTIONS.NTSR) {
          currentWeight =
            current.calculatedBillableWeight < current.ntsRecordedWeight * weightAdjustment
              ? current.calculatedBillableWeight
              : current.ntsRecordedWeight * weightAdjustment;
        }
        return total + currentWeight;
      }, 0);

      return sumOtherEligibleWeights + sumChainWeights > 0 ? sumOtherEligibleWeights + sumChainWeights : null;
    }
    return null;
  }, [mtoShipments, weightAdjustment]);
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
    // Separate diverted shipments and other eligible shipments for weight calculations
    // This is done because a diverted shipment only has one true weight, but when it gets diverted
    // it is entered as a whole new shipment. This causes the sum to be counted twice for its weight,
    // we filter to include only the lowest weight from the diverted shipments here to prevent that.
    const divertedEligibleShipments = mtoShipments.filter(
      (s) => s.diversion && includedStatusesForCalculatingWeights(s.status) && calculateShipmentNetWeight(s),
    );
    const otherEligibleShipments = mtoShipments.filter(
      (s) => !s.diversion && includedStatusesForCalculatingWeights(s.status) && calculateShipmentNetWeight(s),
    );

    // In order to properly sum the lowest weight of the diverted shipments, we must first put them into
    // their correct "chains". Please see comments for groupDivertedShipments for more details.
    const chains = groupDivertedShipmentsByAddress(divertedEligibleShipments);
    // Grab the lowest weight from each chain
    const chainWeights = chains.map((chain) => getLowestShipmentNetWeight(chain));
    // Now that we have the lowest weight from each chain, get the sum
    const sumChainWeights = chainWeights.reduce((total, weight) => total + weight, 0);

    const sumOtherEligibleWeights = otherEligibleShipments.reduce((total, current) => {
      return total + (calculateShipmentNetWeight(current) || 0);
    }, 0);

    return sumOtherEligibleWeights + sumChainWeights > 0 ? sumOtherEligibleWeights + sumChainWeights : 0;
  }
  return null;
};

export const useCalculatedWeightRequested = (mtoShipments) => {
  return useMemo(() => {
    return calculateWeightRequested(mtoShipments);
  }, [mtoShipments]);
};

export const calculateEstimatedWeight = (mtoShipments, shipmentType) => {
  if (mtoShipments?.some((s) => includedStatusesForCalculatingWeights(s.status) && getShipmentEstimatedWeight(s))) {
    // Separate diverted shipments and other eligible shipments for weight calculations
    // This is done because a diverted shipment only has one true weight, but when it gets diverted
    // it is entered as a whole new shipment. This causes the sum to be counted twice for its weight,
    // we filter to include only the lowest weight from the diverted shipments here to prevent that.
    const divertedEligibleShipments = mtoShipments.filter(
      (s) =>
        s.diversion &&
        (shipmentType ? s.shipmentType === shipmentType : true) &&
        includedStatusesForCalculatingWeights(s.status) &&
        getShipmentEstimatedWeight(s),
    );
    const otherEligibleShipments = mtoShipments.filter(
      (s) =>
        !s.diversion &&
        (shipmentType ? s.shipmentType === shipmentType : true) &&
        includedStatusesForCalculatingWeights(s.status) &&
        getShipmentEstimatedWeight(s),
    );

    // In order to properly sum the lowest weight of the diverted shipments, we must first put them into
    // their correct "chains". Please see comments for groupDivertedShipments for more details.
    const chains = groupDivertedShipmentsByAddress(divertedEligibleShipments);
    // Grab the lowest weight from each chain
    const chainWeights = chains.map((chain) => getEstimatedLowestShipmentWeight(chain));
    // Now that we have the lowest weight from each chain, get the sum
    const sumChainWeights = chainWeights.reduce((total, weight) => total + weight, 0);

    // Sum non diverted shipments
    const sumOtherEligibleWeights = otherEligibleShipments.reduce(
      (total, shipment) => total + getShipmentEstimatedWeight(shipment),
      0,
    );

    return sumOtherEligibleWeights + sumChainWeights > 0 ? sumOtherEligibleWeights + sumChainWeights : null;
  }

  return null;
};

export const useCalculatedEstimatedWeight = (mtoShipments) => {
  return useMemo(() => {
    return calculateEstimatedWeight(mtoShipments);
  }, [mtoShipments]);
};

/**
 * This function generates a page subtitle from the path,
 * by splitting the path at slashes, dashes and underscores, capitalizing, and joining with spaces and dashes.
 * e.g. "my-favorite_path/{pathId}/details" becomes "My Favorite Path - {pathId} - Details"
 *
 * @param path The path to convert
 * @return {string} The generated subtitle
 */
export function convertPathToSubtitle(path) {
  return (
    path &&
    path
      .split('/')
      .filter((parameter) => parameter)
      .map((segment) =>
        segment
          .split(/[-_]/)
          .map((word) => capitalize(word))
          .join(' '),
      )
      .join(' - ')
  );
}

/**
 * @func getBasePageTitle
 * @desc Creates a string using the variables with the `_BASE_PAGE_TITLE` suffix defined in and exported from `src/shared/constants`.
 * @returns {string} baseTitle - A base title which may or may not be blank.
 */
export function getBasePageTitle() {
  let baseTitle = '';
  if (isAdminSite) {
    baseTitle = ADMIN_BASE_PAGE_TITLE;
  }
  if (isMilmoveSite) {
    baseTitle = MILMOVE_BASE_PAGE_TITLE;
  }
  if (isOfficeSite) {
    baseTitle = OFFICE_BASE_PAGE_TITLE;
  }
  return baseTitle;
}

/**
 * @func generatePageTitle
 * @desc A function that generates the page title using a provided string appended to the output of `getBasePageTitle`.
 * @param {string} [string] - A string that is appended to the base title for a page.
 * @returns {string} A generated page title.
 */
export function generatePageTitle(string) {
  const baseTitle = getBasePageTitle();
  return baseTitle + (string ? ` - ${string}` : '');
}

/**
 * @func announcePageTitle
 * @desc A function that sets the document's title announcer element to the title
 * @param {string} title - A string that the title announcer's textContent is set to
 */
export function announcePageTitle(title) {
  const titleAnnouncer = document.getElementById('title-announcer');
  if (titleAnnouncer) {
    titleAnnouncer.textContent = title;
  }
}

/**
 * @func useTitle
 * @desc This function generates a subtitle using the pathname from the React Router DOM useLocation function unless a string is passed into the function and assigns the value to `subtitle`. It then calls useEffect to update the document's title attribute and an aria-live element using the `generatePageTitle` function with the internally created subtitle.
 * @param {string} [string] - A string value to be used for the title instead of using the URL path.
 */
export function useTitle(string) {
  const { pathname } = useLocation();
  const subtitle = string || convertPathToSubtitle(pathname);
  useEffect(() => {
    const title = generatePageTitle(subtitle);
    document.title = title;
    announcePageTitle(title);
  }, [subtitle]);
}
