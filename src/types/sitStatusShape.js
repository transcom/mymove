/* eslint-disable import/prefer-default-export */

import PropTypes from 'prop-types';

import { SERVICE_ITEM_CODES } from 'constants/serviceItems';

export const LOCATION_TYPES = {
  ORIGIN: 'ORIGIN',
  DESTINATION: 'DESTINATION',
};

const serviceItemStatuses = {
  SUBMITTED: 'SUBMITTED',
  APPROVED: 'APPROVED',
  REJECTED: 'REJECTED',
};

export const LOCATION_TYPES_ONE_OF = PropTypes.oneOf([LOCATION_TYPES.ORIGIN, LOCATION_TYPES.DESTINATION]);

const serviceItemStatusesOneOf = PropTypes.oneOf([
  serviceItemStatuses.SUBMITTED,
  serviceItemStatuses.APPROVED,
  serviceItemStatuses.REJECTED,
]);

const serviceItemCodesOneOf = PropTypes.oneOf([
  SERVICE_ITEM_CODES.DDASIT,
  SERVICE_ITEM_CODES.DDDSIT,
  SERVICE_ITEM_CODES.DDFSIT,
  SERVICE_ITEM_CODES.DDP,
  SERVICE_ITEM_CODES.DLH,
  SERVICE_ITEM_CODES.DOASIT,
  SERVICE_ITEM_CODES.DOFSIT,
  SERVICE_ITEM_CODES.DOP,
  SERVICE_ITEM_CODES.DOPSIT,
  SERVICE_ITEM_CODES.DOSHUT,
  SERVICE_ITEM_CODES.DPK,
  SERVICE_ITEM_CODES.DSH,
  SERVICE_ITEM_CODES.DUPK,
  SERVICE_ITEM_CODES.FSC,
  SERVICE_ITEM_CODES.DDSHUT,
  SERVICE_ITEM_CODES.DCRT,
  SERVICE_ITEM_CODES.DUCRT,
]);

const PastServiceItemsShape = PropTypes.shape({
  SITPostalCode: PropTypes.string,
  createdAt: PropTypes.string,
  deletedAt: PropTypes.string,
  description: PropTypes.string,
  eTag: PropTypes.string,
  id: PropTypes.string,
  moveTaskOrderID: PropTypes.string,
  mtoShipmentID: PropTypes.string,
  pickupPostalCode: PropTypes.string,
  reServiceCode: serviceItemCodesOneOf,
  reServiceID: PropTypes.string,
  reServiceName: PropTypes.string,
  reason: PropTypes.string,
  sitDepartureDate: PropTypes.string,
  sitEntryDate: PropTypes.string,
  status: serviceItemStatusesOneOf,
  submittedAt: PropTypes.string,
  updatedAt: PropTypes.string,
});

export const SitStatusShape = PropTypes.shape({
  location: LOCATION_TYPES_ONE_OF,
  totalSITDaysUsed: PropTypes.number,
  totalDaysRemaining: PropTypes.number,
  daysInSIT: PropTypes.number,
  sitEntryDate: PropTypes.string,
  sitDepartureDate: PropTypes.string,
  pastSITServiceItems: PropTypes.arrayOf(PastServiceItemsShape),
});
