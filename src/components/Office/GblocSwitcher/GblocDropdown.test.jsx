import React, { render, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import GblocDropdown, { gblocDropdownTestId } from './GblocDropdown';

describe('GblocDropdown', () => {
  it('displays the proper attributes with gbloc options', async () => {
    const changeHandler = jest.fn();

    const KKFA = 'KKFA';
    const USMC = 'USMC';
    const officeGblocs = [KKFA, USMC];

    render(<GblocDropdown ariaLabel="the label" handleChange={changeHandler} gblocs={officeGblocs} />);

    const user = userEvent.setup();
    const gblocCombobox = await screen.findByTestId(gblocDropdownTestId);

    expect(gblocCombobox).toBeInstanceOf(HTMLSelectElement);
    expect((await within(gblocCombobox).findAllByRole('option')).length).toBeGreaterThan(0);

    await user.selectOptions(gblocCombobox, KKFA);

    expect(changeHandler).toHaveBeenCalledWith(KKFA);
  });
});
