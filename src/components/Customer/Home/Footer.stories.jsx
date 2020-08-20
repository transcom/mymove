/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';
import { withKnobs, text } from '@storybook/addon-knobs';

import Footer from './Footer';

export const Basic = () => (
  <div className="grid-container">
    <Footer
      header={text('Header', 'This is the header')}
      dutyStationName={text('Duty Station Name', 'Some duty station')}
      officeType={text('Office type', 'Some office type')}
      telephone={text('Telephone', '(777) 777-7777')}
    />
  </div>
);

export default {
  title: 'Customer Components | Footer',
  decorators: [withKnobs],
};
