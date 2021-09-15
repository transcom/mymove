/* eslint-disable import/prefer-default-export */
import PropTypes from 'prop-types';

import { SIT_EXTENSION_REASON, SIT_EXTENSION_STATUS } from '../constants/sitExtensions';

export const SITExtensionReasons = PropTypes.oneOf([
  SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER,
  SIT_EXTENSION_REASON.SERIOUS_ILLNESS_DEPENDENT,
  SIT_EXTENSION_REASON.IMPENDING_ASSIGNEMENT,
  SIT_EXTENSION_REASON.DIRECTED_TEMPORARY_DUTY,
  SIT_EXTENSION_REASON.NONAVAILABILITY_OF_CIVILIAN_HOUSING,
  SIT_EXTENSION_REASON.AWAITING_COMPLETION_OF_RESIDENCE,
  SIT_EXTENSION_REASON.OTHER,
]);

export const SITExtensionShape = PropTypes.shape({
  mtoShipmentID: PropTypes.string,
  requestReason: PropTypes.string,
  contractorRemarks: PropTypes.string,
  requestedDays: PropTypes.number,
  status: PropTypes.oneOf([SIT_EXTENSION_STATUS.APPROVED, SIT_EXTENSION_STATUS.PENDING, SIT_EXTENSION_STATUS.DENIED]),
  approvedDays: PropTypes.number,
  decisionDate: PropTypes.string,
  officeRemarks: PropTypes.string,
  createdAt: PropTypes.string,
  updatedAt: PropTypes.string,
  eTag: PropTypes.string,
});
