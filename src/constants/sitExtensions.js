/* eslint-disable import/prefer-default-export */

export const SIT_EXTENSION_REASON = {
  SERIOUS_ILLNESS_MEMBER: 'SERIOUS_ILLNESS_MEMBER',
  SERIOUS_ILLNESS_DEPENDENT: 'SERIOUS_ILLNESS_DEPENDENT',
  IMPENDING_ASSIGNEMENT: 'IMPENDING_ASSIGNEMENT',
  DIRECTED_TEMPORARY_DUTY: 'DIRECTED_TEMPORARY_DUTY',
  NONAVAILABILITY_OF_CIVILIAN_HOUSING: 'NONAVAILABILITY_OF_CIVILIAN_HOUSING',
  AWAITING_COMPLETION_OF_RESIDENCE: 'AWAITING_COMPLETION_OF_RESIDENCE',
  OTHER: 'OTHER',
};

export const SIT_EXTENSION_STATUS = {
  PENDING: 'PENDING',
  APPROVED: 'APPROVED',
  DENIED: 'DENIED',
  REMOVED: 'REMOVED',
};

export const sitExtensionReasons = {
  [SIT_EXTENSION_REASON.SERIOUS_ILLNESS_MEMBER]: 'Serious illness of the member',
  [SIT_EXTENSION_REASON.SERIOUS_ILLNESS_DEPENDENT]: 'Serious illness or death of a dependent',
  [SIT_EXTENSION_REASON.IMPENDING_ASSIGNEMENT]: 'Impending assignment to government quarters',
  [SIT_EXTENSION_REASON.DIRECTED_TEMPORARY_DUTY]: 'Directed temporary duty after arrival at permanent duty location',
  [SIT_EXTENSION_REASON.NONAVAILABILITY_OF_CIVILIAN_HOUSING]: 'Nonavailability of suitable civilian housing',
  [SIT_EXTENSION_REASON.AWAITING_COMPLETION_OF_RESIDENCE]: 'Awaiting completion of residence under construction',
  [SIT_EXTENSION_REASON.OTHER]: 'Other reason',
};

export const SIT_EXTENSION_STATUSES = {
  [SIT_EXTENSION_STATUS.PENDING]: 'Pending',
  [SIT_EXTENSION_STATUS.APPROVED]: 'Approved',
  [SIT_EXTENSION_STATUS.DENIED]: 'Denied',
};

export const SIT_EXTENSION_REASONS = {
  SERIOUS_ILLNESS_MEMBER: 'Serious illness of the member',
  SERIOUS_ILLNESS_DEPENDENT: 'Serious illness or death of a dependent',
  IMPENDING_ASSIGNEMENT: 'Impending assignment to government quarters',
  DIRECTED_TEMPORARY_DUTY: 'Directed temporary duty after arrival at permanent duty location',
  NONAVAILABILITY_OF_CIVILIAN_HOUSING: 'Nonavailability of suitable civilian housing',
  AWAITING_COMPLETION_OF_RESIDENCE: 'Awaiting completion of residence under construction',
  OTHER: 'Other reason',
};
