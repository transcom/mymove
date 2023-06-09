import { useEffect, useMemo } from 'react';
import { useLocation } from 'react-router-dom';
import { capitalize } from 'lodash';

import { isAdminSite, isMilmoveSite, isOfficeSite } from '../shared/constants';

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
        ?.filter((s) => {
          return includedStatusesForCalculatingWeights(s.status) && s.calculatedBillableWeight;
        })
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
