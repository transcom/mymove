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
    streetAddress1: '1292 Orchard Terrace',
    streetAddress2: 'Building C, Unit 10',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79912',
    county: 'El Paso',
  },
  backupMailingAddress: {
    streetAddress1: '448 Washington Blvd NE',
    streetAddress2: '',
    city: 'El Paso',
    state: 'TX',
    postalCode: '79936',
    county: 'El Paso',
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
