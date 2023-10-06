export const SIT_ADDRESS_UPDATE_STATUS = {
  REQUESTED: 'REQUESTED',
  REJECTED: 'REJECTED',
  APPROVED: 'APPROVED',
};

// allowing edit of Domestic origin SIT pickup (DOPSIT)
// allowing edit of Domestic destination SIT delivery (DDDSIT)
// allowing edit of Domestic destination 1st day SIT (DDFSIT)
export const ALLOWED_SIT_ADDRESS_UPDATE_SI_CODES = ['DDDSIT', 'DOPSIT', 'DDFSIT'];
