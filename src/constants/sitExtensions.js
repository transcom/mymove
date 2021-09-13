/* eslint-disable import/prefer-default-export */
import { SIT_EXTENSION_REASON, SIT_EXTENSION_STATUS } from 'shared/constants';

export const sitExtensionReasons = {
  [SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER]: 'Serious illness of the member',
  [SIT_EXTENSION_REASON.SERIOUS_ILLNESS_DEPENDENT]: 'Serious illness or death of a dependent',
  [SIT_EXTENSION_REASON.IMPENDING_ASSIGNEMENT]: 'Impending assignment to government quarters',
  [SIT_EXTENSION_REASON.DIRECTED_TEMPORARY_DUTY]: 'Directed temporary duty after arrival at permanent duty station',
  [SIT_EXTENSION_REASON.NONAVAILABILITY_OF_CIVILIAN_HOUSING]: 'Nonavailability of suitable civilian housing',
  [SIT_EXTENSION_REASON.AWAITING_COMPLETION_OF_RESIDENCE]: 'Awaiting completion of residence under construction',
  [SIT_EXTENSION_REASON.OTHER]: 'Other reason',
};

export const SIT_EXTENSION_STATUSES = {
  [SIT_EXTENSION_STATUS.PENDING]: 'Pending',
  [SIT_EXTENSION_STATUS.APPROVED]: 'Approved',
  [SIT_EXTENSION_STATUS.DENIED]: 'Denied',
};
