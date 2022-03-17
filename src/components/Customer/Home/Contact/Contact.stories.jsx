/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { text } from '@storybook/addon-knobs';

import Contact from './index';

export const Basic = () => (
  <Contact
    dutyLocationName={text('Duty Location Name', 'Fort Knox')}
    header={text('Header', 'Contacts')}
    officeType={text('Office type', 'Origin Transportation Office')}
    telephone={text('Telephone', '(777) 777-7777')}
  />
);

export const missingPhone = () => (
  <Contact
    dutyLocationName={text('Duty Location Name', 'Fort Knox')}
    header={text('Header', 'Contacts')}
    officeType={text('Office type', 'Origin Transportation Office')}
  />
);

export const nonInstallation = () => (
  <Contact
    header={text('Header', 'Contacts')}
    dutyLocationName={text('Duty Station', '')}
    officeType={text('Office type', '')}
    telephone={text('Telephone', '')}
  />
);

export default {
  title: 'Customer Components / Contact',
  decorators: [
    (Story) => (
      <div className="grid-container">
        <Story />
      </div>
    ),
  ],
};
