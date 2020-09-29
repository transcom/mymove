/* eslint-disable react/jsx-props-no-spreading */
import React from 'react';

import SelectableCard from './SelectableCard';

const defaultProps = {
  id: 'PPM',
  label: 'A great choice for you and your family',
  value: 'PPM',
  name: 'shipmentType',
  cardText:
    "Maroon wherry swing the lead spanker Brethren of the Coast aft heave down shrouds grapple ballast. Crow's nest hardtack yardarm lee driver spirits Admiral of the Black take a caulk crimp chandler. ",
  onChange: () => {
    console.log('changed!'); // eslint-disable-line no-console
  },
  checked: false,
};

export default {
  title: 'Customer Components | SelectableCard',
};

export const Basic = () => (
  <div>
    <SelectableCard {...defaultProps} />
  </div>
);
