import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ButtonDropdown from './ButtonDropdown';

describe('ButtonDropdown component', () => {
  it('renders the USWDS dropdown component without errors', () => {
    render(
      <ButtonDropdown onChange={() => {}}>
        <option>- Select - </option>
        <option value="value1">Option A</option>
        <option value="value2">Option B</option>
        <option value="value3">Option C</option>
      </ButtonDropdown>,
    );

    expect(screen.queryByTestId('dropdown')).toBeInTheDocument();
  });

  it('calls onChange prop on option selection', async () => {
    const onChangeSpy = jest.fn();
    render(
      <ButtonDropdown onChange={onChangeSpy}>
        <option>- Select - </option>
        <option value="value1">Option A</option>
        <option value="value2">Option B</option>
        <option value="value3">Option C</option>
      </ButtonDropdown>,
    );
    const dropdown = screen.getByRole('combobox');
    await userEvent.selectOptions(dropdown, 'value1');

    expect(onChangeSpy).toHaveBeenCalled();
  });
});
