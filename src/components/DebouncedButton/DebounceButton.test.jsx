import React from 'react';
import { render, screen } from '@testing-library/react';

import DebounceButton from './DebounceButton';

describe('DebounceButton component', () => {
  it('renders the DebounceButton component without errors', () => {
    render(<DebounceButton>Button Text</DebounceButton>);

    expect(screen.queryByTestId('debounce-button')).toBeInTheDocument();
  });

  // it('calls onChange prop on option selection', async () => {
  //   const onChangeSpy = jest.fn();
  //   render(
  //     <ButtonDropdown onChange={onChangeSpy}>
  //       <option>- Select - </option>
  //       <option value="value1">Option A</option>
  //       <option value="value2">Option B</option>
  //       <option value="value3">Option C</option>
  //     </ButtonDropdown>,
  //   );
  //   const dropdown = screen.getByRole('combobox');
  //   await userEvent.selectOptions(dropdown, 'value1');

  //   expect(onChangeSpy).toHaveBeenCalled();
  // });
});
