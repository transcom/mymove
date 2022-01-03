import React from 'react';
import { object } from '@storybook/addon-knobs';

import CustomerInfoList from './CustomerInfoList';

export default {
  title: 'Office Components/CustomerInfoList',
  component: CustomerInfoList,
};

const info = {
  name: 'Smith, Kerry',
  dodId: '9999999999',
  phone: '+1 999-999-9999',
  email: 'ksmith@email.com',
  currentAddress: {
    streetAddress1: '812 S 129th St',
    city: 'San Antonio',
    state: 'TX',
    postalCode: '78234',
  },
  backupContact: {
    name: 'Quinn Ocampo',
    email: 'quinnocampo@myemail.com',
    phone: '999-999-9999',
  },
};

export const Basic = () => (
  <div className="officeApp">
    <CustomerInfoList customerInfo={object('info', info)} />
  </div>
);
