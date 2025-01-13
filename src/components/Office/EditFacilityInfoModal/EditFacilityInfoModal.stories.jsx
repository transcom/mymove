import React from 'react';
import { Provider } from 'react-redux';

import { EditFacilityInfoModal } from './EditFacilityInfoModal';

import { configureStore } from 'shared/store';

const mockStore = configureStore({});

const storageFacility = {
  address: {
    streetAddress1: '123 Fake Street',
    streetAddress2: '',
    city: 'Pasadena',
    state: 'CA',
    postalCode: '90210',
    county: 'Los Angeles',
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
    county: 'Los Angeles',
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
    <Provider store={mockStore.store}>
      <EditFacilityInfoModal
        onSubmit={() => {}}
        onClose={() => {}}
        serviceOrderNumber="12345"
        storageFacility={storageFacility}
        shipmentType="HHG_INTO_NTS"
      />
    </Provider>
  </div>
);

export const WithInfoMissing = () => (
  <div className="officeApp">
    <Provider store={mockStore.store}>
      <EditFacilityInfoModal
        onSubmit={() => {}}
        onClose={() => {}}
        serviceOrderNumber="12345"
        storageFacility={storageFacilityInfoMissing}
        shipmentType="HHG_INTO_NTS"
      />
    </Provider>
  </div>
);
