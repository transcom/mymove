import React from 'react';

import EditFacilityInfoModal from './EditFacilityInfoModal';

const storageFacilityAddress = {
  address: {
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    city: 'Pasadena',
    state: 'CA',
    postalCode: '90210',
  },
  lotNumber: '11232',
};
const storageFacility = {
  facilityName: 'My Facility',
  phone: '1235553434',
  email: 'my@email.com',
  serviceOrderNumber: '12345',
};

export default {
  title: 'Office Components/EditFacilityInfoModal',
  component: EditFacilityInfoModal,
};

export const Basic = () => (
  <EditFacilityInfoModal
    onSubmit={() => {}}
    onClose={() => {}}
    storageFacility={storageFacility}
    storageFacilityAddress={storageFacilityAddress}
  />
);

export const WithInfoMissing = () => (
  <EditFacilityInfoModal
    onSubmit={() => {}}
    onClose={() => {}}
    storageFacility={{
      facilityName: '',
      phone: '1235553434',
      email: 'my@email.com',
      serviceOrderNumber: '12345',
    }}
    storageFacilityAddress={storageFacilityAddress}
  />
);
