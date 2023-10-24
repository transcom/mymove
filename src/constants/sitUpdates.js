export const SIT_ADDRESS_UPDATE_STATUS = {
  REQUESTED: 'REQUESTED',
  REJECTED: 'REJECTED',
  APPROVED: 'APPROVED',
};

// allowing edit of address for Domestic destination SIT delivery (DDDSIT)
// allowing edit of SIT entry date for Domestic destination 1st day SIT (DDFSIT)
// allowing edit of SIT entry date for Domestic origin 1st day SIT (DOFSIT)
export const ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES = ['DDDSIT', 'DOFSIT', 'DDFSIT'];
