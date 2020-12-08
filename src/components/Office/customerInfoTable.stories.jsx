import React from 'react';
import { object } from '@storybook/addon-knobs';

import CustomerInfoTable from './CustomerInfoTable';

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
  backupContact: {
    name: 'Quinn Ocampo',
    email: 'quinnocampo@myemail.com',
    phone: '999-999-9999',
  },
};

export default {
  title: 'TOO/TIO Components/CustomerInfoTable',
};

export const Default = () => <CustomerInfoTable customerInfo={object('customerInfo', info)} />;
