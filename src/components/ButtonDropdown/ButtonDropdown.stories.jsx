import React from 'react';

import ButtonDropdown from './ButtonDropdown';

export default {
  title: 'Components/ButtonDropdown',
  component: ButtonDropdown,
};

export const defaultDropdown = () => (
  <ButtonDropdown onChange={() => {}} ariaLabel="Shipment selection">
    <option>Add a new shipment</option>
    <option value="optionA">Option A</option>
    <option value="optionB">Option B</option>
    <option value="optionC">Option C</option>
  </ButtonDropdown>
);
