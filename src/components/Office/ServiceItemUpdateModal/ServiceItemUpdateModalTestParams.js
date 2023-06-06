const finalAddress = {
  city: 'Fairfax',
  state: 'VA',
  postalCode: '12345',
  streetAddress1: '123 Fake Street',
  streetAddress2: '',
  streetAddress3: '',
  country: 'USA',
};
const oldAddress = {
  city: 'Richmond',
  state: 'VA',
  postalCode: '12508',
  streetAddress1: '345 Faker Rd',
  streetAddress2: '',
  streetAddress3: '',
  country: 'USA',
};
export const newAddress = {
  city: 'Alexandria',
  state: 'VA',
  postalCode: '12867',
  streetAddress1: '555 Fakest Dr',
  streetAddress2: 'Unit 133',
  streetAddress3: '',
  country: 'USA',
};

const originalAddress = {
  city: 'Alexandria',
  state: 'VA',
  postalCode: '12867',
  streetAddress1: '333 Most Fake Blvd',
  streetAddress2: '',
  streetAddress3: '',
  country: 'USA',
};

export const sitAddressUpdates = [
  {
    id: '123xyz',
    mtoServiceItemID: 'abc123',
    officeRemarks: '',
    contractorRemarks: 'Customer wishes to be closer to family',
    distance: 500,
    status: 'APPROVED',
    oldAddress,
    newAddress,
    createdAt: '2020-11-20',
    updatedAt: '2020-11-20',
  },
  {
    id: '456efg',
    mtoServiceItemID: 'abc123',
    officeRemarks: '',
    contractorRemarks: '',
    distance: 500,
    status: 'APPROVED',
    oldAddress: originalAddress,
    newAddress: oldAddress,
    createdAt: '2020-11-20',
    updatedAt: '2020-11-20',
  },
];

export const domesticDestinationSitServiceItem = {
  id: 'abc123',
  code: 'DDDSIT',
  submittedAt: '2020-11-20',
  serviceItem: 'Domestic destination SIT',
  mtoServiceItemID: '678hij',
  details: {
    reason: "Customer's housing at base is not ready",
    firstCustomerContact: { timeMilitary: '1200Z', firstAvailableDeliveryDate: '2020-09-15' },
    secondCustomerContact: { timeMilitary: '2300Z', firstAvailableDeliveryDate: '2020-09-21' },
    serviceItem: 'Domestic Destination SIT',
  },
  sitDestinationFinalAddress: finalAddress,
  sitDestinationOriginalAddress: originalAddress,
};

export const dddSitWithAddressUpdate = {
  ...domesticDestinationSitServiceItem,
  sitAddressUpdates,
};
