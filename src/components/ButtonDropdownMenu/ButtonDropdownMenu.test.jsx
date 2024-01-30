import React from 'react';
import { render, fireEvent } from '@testing-library/react';
import '@testing-library/jest-dom/extend-expect'; // For expect assertions

import ButtonDropdownMenu from './ButtonDropdownMenu';

describe('ButtonDropdownMenu', () => {
  const items = [
    { id: 1, value: 'Item 1' },
    { id: 2, value: 'Item 2' },
  ];

  it('calls onItemClick callback when an item is clicked', () => {
    const onItemClickMock = jest.fn();
    const { getByText } = render(<ButtonDropdownMenu title="Test" items={items} onItemClick={onItemClickMock} />);

    // Open the dropdown
    fireEvent.click(getByText('Test'));

    // Click on the first item
    fireEvent.click(getByText('Item 1'));

    // Ensure that onItemClick is called with the correct item
    expect(onItemClickMock).toHaveBeenCalledWith(items[0]);

    // Close the dropdown (optional)
    fireEvent.click(getByText('Test'));
  });
});
