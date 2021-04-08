import React from 'react';

import ContactInfo from 'components/Customer/Profile/ContactInfo/ContactInfo';
import SectionWrapper from 'components/Customer/SectionWrapper';

export default {
  title: 'Customer Components / Profile / ContactInfo',
  component: ContactInfo,
  argTypes: {
    onEditClick: 'go to edit page',
  },
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
};

export const DefaultState = (argTypes) => (
  <SectionWrapper>
    <ContactInfo {...baseProps} onEditClick={argTypes.onEditClick} />
  </SectionWrapper>
);

export const WithAltPhone = (argTypes) => (
  <SectionWrapper>
    <ContactInfo {...baseProps} secondaryTelephone="619-555-3000" onEditClick={argTypes.onEditClick} />
  </SectionWrapper>
);
