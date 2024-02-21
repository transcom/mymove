import React from 'react';

import ButtonDropdownMenu from './ButtonDropdownMenu';

export default {
  title: 'Components/ButtonDropdownMenu',
  component: ButtonDropdownMenu,
};

const dropdownMenuItems = [
  {
    id: 1,
    value: 'PCS Orders',
  },
  {
    id: 2,
    value: 'PPM Packet',
  },
];

export const defaultDropdown = () => <ButtonDropdownMenu title="Download" items={dropdownMenuItems} outline />;
