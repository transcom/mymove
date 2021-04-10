import React from 'react';

import ContactInfoDisplay from './ContactInfoDisplay';

import { MockProviders } from 'testUtils';

export default {
  title: 'Customer Components / Profile / ContactInfoDisplay',
  component: ContactInfoDisplay,
  decorators: [
    (Story) => (
      <MockProviders>
        <Story />
      </MockProviders>
    ),
  ],
};

const baseProps = {
  telephone: '703-555-4578',
  personalEmail: 'test@example.com',
  emailIsPreferred: true,
  residentialAddress: {
    street_address_1: '1292 Orchard Terrace',
    street_address_2: 'Building C, Unit 10',
    city: 'El Paso',
    state: 'TX',
    postal_code: '79912',
  },
  backupMailingAddress: {
    street_address_1: '448 Washington Blvd NE',
    street_address_2: '',
    city: 'El Paso',
    state: 'TX',
    postal_code: '79936',
  },
  backupContact: {
    name: 'Gabriela Sáenz Perez',
    telephone: '206-555-8989',
    email: 'gsp@example.com',
  },
  editURL: '/moves/review/edit-profile',
};

export const DefaultState = () => <ContactInfoDisplay {...baseProps} />;

export const WithAltPhone = () => <ContactInfoDisplay {...baseProps} secondaryTelephone="619-555-3000" />;
