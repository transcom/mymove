// allowing edit of SIT entry date for Domestic destination 1st day SIT (DDFSIT)
// allowing edit of SIT entry date for Domestic origin 1st day SIT (DOFSIT)
export const ALLOWED_SIT_UPDATE_SI_CODES = ['DOFSIT', 'DDFSIT'];

// allowing display of old service item details for following SIT types which can be resubmitted
export const ALLOWED_RESUBMISSION_SI_CODES = [
  'DDDSIT',
  'DOFSIT',
  'DDFSIT',
  'DOASIT',
  'DOPSIT',
  'DOSFSC',
  'DDASIT',
  'DDSFSC',
];
