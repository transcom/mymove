/*  react/jsx-props-no-spreading */
import React from 'react';
import { withKnobs, text } from '@storybook/addon-knobs';

import Contact from '.';

export const Basic = () => (
  <div className="grid-container">
    <Contact
      header={text('Header', 'This is the header')}
      dutyStationName={text('Duty Station Name', 'Some duty station')}
      officeType={text('Office type', 'Some office type')}
      telephone={text('Telephone', '(777) 777-7777')}
    />
  </div>
);

export default {
  title: 'Customer Components | Contact',
  decorators: [withKnobs],
};
