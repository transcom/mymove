import React from 'react';
import { withKnobs, object } from '@storybook/addon-knobs';

import CustomerInfoTable from '../components/Office/CustomerInfoTable';

const info = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    street_address_1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postal_code: '78234',
  },
  destinationAddress: {
    street_address_1: '441 SW Rio de la Plata Drive',
    city: 'Tacoma',
    state: 'WA',
    postal_code: '98421',
  },
  backupContactName: 'Quinn Ocampo',
  backupContactPhone: '+1 999-999-9999',
  backupContactEmail: 'quinnocampo@myemail.com',
};

export default {
  title: 'TOO/TIO Components|CustomerInfoTable',
  decorator: withKnobs,
};

export const Default = () => <CustomerInfoTable customerInfo={object('customerInfo', info)} />;
