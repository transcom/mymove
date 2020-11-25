/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { withKnobs, text } from '@storybook/addon-knobs';

import Contact from '.';

export const Basic = () => (
  <div className="grid-container">
    <Contact
      dutyStationName={text('Duty Station Name', 'Fort Knox')}
      header={text('Header', 'Contacts')}
      officeType={text('Office type', 'Origin Transportation Office')}
      telephone={text('Telephone', '(777) 777-7777')}
    />
  </div>
);

export const MoveSubmitted = () => (
  <div className="grid-container">
    <Contact
      dutyStationName={text('Duty Station Name', 'Fort Knox')}
      header={text('Header', 'Contacts')}
      moveSubmitted
      officeType={text('Office type', 'Origin Transportation Office')}
      telephone={text('Telephone', '(777) 777-7777')}
    />
  </div>
);

export default {
  title: 'Customer Components | Contact',
  decorators: [withKnobs],
};
