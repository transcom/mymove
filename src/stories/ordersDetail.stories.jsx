import React from 'react';
import { action } from '@storybook/addon-actions';
import { storiesOf } from '@storybook/react';
import { OrdersDetailForm } from 'components/Office/OrdersDetailForm';

const currentDutyStation = {
  address: {
    city: 'Dover AFB',
    id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
    postal_code: '19902',
    state: 'DE',
    street_address_1: '',
  },
  address_id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
  affiliation: 'AIR_FORCE',
  created_at: '2018-10-04T22:54:46.589Z',
  id: '071f6286-8255-4e35-b8ac-0e7fe1d10aa4',
  name: 'Dover AFB',
  updated_at: '2018-10-04T22:54:46.589Z',
};
const newDutyStation = {
  address: {
    city: 'Scott Air Force Base',
    id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
    postal_code: '62225',
    state: 'IL',
    street_address_1: '',
  },
  address_id: '9f8b0fad-afe1-4a44-bb28-296a335c1141',
  affiliation: 'AIR_FORCE',
  created_at: '2018-10-04T22:54:46.589Z',
  id: '071f6286-8255-4e35-b8ac-0e7fe1d10aa4',
  name: 'Scott AFB',
  updated_at: '2018-10-04T22:54:46.589Z',
};

const deptIndicatorOptions = Object.entries({
  NAVY_AND_MARINES: '17 Navy and Marine Corps',
  ARMY: '21 Army',
  AIR_FORCE: '57 Air Force',
  COAST_GUARD: '70 Coast Guard',
});

const ordersTypeOptions = Object.entries({
  PERMANENT_CHANGE_OF_STATION: 'Permanent Change Of Station',
});

const ordersTypeDetailOptions = Object.entries({
  HHG_PERMITTED: 'Shipment of HHG Permitted',
  PCS_TDY: 'PCS with TDY Enroute',
  HHG_RESTRICTED_PROHIBITED: 'Shipment of HHG Restricted or Prohibited',
  HHG_RESTRICTED_AREA: 'HHG Restricted Area-HHG Prohibited',
  INSTRUCTION_20_WEEKS: 'Course of Instruction 20 Weeks or More',
  HHG_PROHIBITED_20_WEEKS: 'Shipment of HHG Prohibited but Authorized within 20 weeks',
  DELAYED_APPROVAL: 'Delayed Approval 20 Weeks or More',
});

const OrdersDetail = () => (
  <div>
    <OrdersDetailForm
      initialValues={{
        currentDutyStation,
        newDutyStation,
        dateIssued: '08 Mar 2020',
        reportByDate: '01 Apr 2020',
        departmentIndicator: 'NAVY_AND_MARINES',
        ordersNumber: '999999999',
        ordersType: 'PERMANENT_CHANGE_OF_STATION',
        ordersTypeDetail: 'HHG_PERMITTED',
        tac: 'Tac',
        sac: 'Sac',
      }}
      deptIndicatorOptions={deptIndicatorOptions}
      ordersTypeOptions={ordersTypeOptions}
      ordersTypeDetailOptions={ordersTypeDetailOptions}
      onSubmit={action('Orders Detail Submit')}
      onReset={action('Orders Detail Cancel')}
    />
  </div>
);

storiesOf('TOO/TIO Components|OrdersDetailForm', module).add('with buttons to edit', () => (
  <div style={{ padding: `20px`, background: `#f0f0f0` }}>
    <OrdersDetail />
  </div>
));
