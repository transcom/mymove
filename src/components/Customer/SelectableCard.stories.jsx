/*  react/jsx-props-no-spreading */
import React from 'react';
import { action } from '@storybook/addon-actions';

import SelectableCard from './SelectableCard';

const defaultProps = {
  id: 'PPM',
  label: 'A great choice for you and your family',
  value: 'PPM',
  name: 'shipmentType',
  cardText:
    "Maroon wherry swing the lead spanker Brethren of the Coast aft heave down shrouds grapple ballast. Crow's nest hardtack yardarm lee driver spirits Admiral of the Black take a caulk crimp chandler. ",
  onChange: () => {
    console.log('changed!'); // -line no-console
  },
  checked: false,
};

const selectedProps = {
  checked: true,
};

const disabledProps = {
  disabled: true,
};

export default {
  title: 'Customer Components | SelectableCard',
  component: SelectableCard,
  decorators: [
    (Story) => (
      <div style={{ marginTop: '3em' }}>
        <Story />
      </div>
    ),
  ],
};

export const Unselected = () => (
  <div>
    <SelectableCard {...defaultProps} />
  </div>
);

export const Selected = () => {
  const props = { ...defaultProps, ...selectedProps };
  return (
    <div>
      <SelectableCard {...props} />
    </div>
  );
};

export const Disabled = () => {
  const props = { ...defaultProps, ...disabledProps };
  return (
    <div>
      <SelectableCard {...props} />
    </div>
  );
};

export const WithHelpButton = () => (
  <div>
    <SelectableCard {...defaultProps} onHelpClick={action('Open help')} />
  </div>
);
