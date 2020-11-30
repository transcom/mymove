/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { withKnobs, text } from '@storybook/addon-knobs';

import Contact from './index';

export const Basic = () => (
  <Contact
    dutyStationName={text('Duty Station Name', 'Fort Knox')}
    header={text('Header', 'Contacts')}
    officeType={text('Office type', 'Origin Transportation Office')}
    telephone={text('Telephone', '(777) 777-7777')}
  />
);

export const MoveSubmitted = () => (
  <Contact
    dutyStationName={text('Duty Station Name', 'Fort Knox')}
    header={text('Header', 'Contacts')}
    moveSubmitted
    officeType={text('Office type', 'Origin Transportation Office')}
    telephone={text('Telephone', '(777) 777-7777')}
  />
);

export const missingPhone = () => (
  <Contact
    dutyStationName={text('Duty Station Name', 'Fort Knox')}
    header={text('Header', 'Contacts')}
    officeType={text('Office type', 'Origin Transportation Office')}
  />
);

export default {
  title: 'Customer Components | Contact',
  decorators: [
    withKnobs,
    (Story) => (
      <div className="grid-container">
        <Story />
      </div>
    ),
  ],
};
