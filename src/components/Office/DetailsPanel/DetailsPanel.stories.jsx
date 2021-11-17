import React from 'react';
import { object } from '@storybook/addon-knobs';

import CustomerInfoList from '../DefinitionLists/CustomerInfoList';

import DetailsPanel from './DetailsPanel';

export default {
  title: 'Office Components/DetailsPanel',
  component: DetailsPanel,
};

export const Basic = () => (
  <DetailsPanel title="Details panel">
    <p>Child content!</p>
  </DetailsPanel>
);

export const WithEditButton = () => (
  <DetailsPanel
    title="Details panel with edit button"
    editButton={
      <a href="#" className="usa-button usa-button--secondary">
        Edit
      </a>
    }
  >
    <p>Child content!</p>
  </DetailsPanel>
);

export const WithClassname = () => (
  <DetailsPanel title="Details panel with added CSS class" className="border-2px">
    <p>I have a border class added via props!</p>
  </DetailsPanel>
);

export const WithTag = () => (
  <DetailsPanel title="Details panel with tag" tag="NEW">
    <p>I have a tag added via props!</p>
  </DetailsPanel>
);

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

export const WithDefinitionListComponent = () => (
  <DetailsPanel
    title="Customer info"
    editButton={
      <a href="#" className="usa-button usa-button--secondary">
        Edit customer info
      </a>
    }
  >
    <CustomerInfoList customerInfo={object('customerInfo', info)} />
  </DetailsPanel>
);
