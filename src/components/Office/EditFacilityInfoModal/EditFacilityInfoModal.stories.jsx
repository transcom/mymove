import React from 'react';

import { EditFacilityInfoModal } from './EditFacilityInfoModal';

const storageFacility = {
  address: {
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    city: 'Pasadena',
    state: 'CA',
    postalCode: '90210',
  },
  lotNumber: '11232',
  facilityName: 'My Facility',
  phone: '915-555-2942',
  email: 'my@email.com',
};

const storageFacilityInfoMissing = {
  address: {
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    city: 'Pasadena',
    state: 'CA',
    postalCode: '90210',
  },
  lotNumber: '11232',
  facilityName: '',
  phone: '915-555-2942',
  email: 'my@email.com',
};

export default {
  title: 'Office Components/EditFacilityInfoModal',
  component: EditFacilityInfoModal,
};

export const Basic = () => (
  <div className="officeApp">
    <EditFacilityInfoModal
      onSubmit={() => {}}
      onClose={() => {}}
      serviceOrderNumber="12345"
      storageFacility={storageFacility}
      shipmentType="HHG_INTO_NTS_DOMESTIC"
    />
  </div>
);

export const WithInfoMissing = () => (
  <div className="officeApp">
    <EditFacilityInfoModal
      onSubmit={() => {}}
      onClose={() => {}}
      serviceOrderNumber="12345"
      storageFacility={storageFacilityInfoMissing}
      shipmentType="HHG_INTO_NTS_DOMESTIC"
    />
  </div>
);
