import React from 'react';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import TextBoxFilter from './TextBoxFilter';

describe('Table/Filters/TextBoxFilter', () => {
  it('renders without crashing', () => {
    render(
      <TextBoxFilter
        column={{
          filterValue: 'test value',
          setFilter: jest.fn(),
        }}
      />,
    );

    expect(screen.getByRole('textbox')).toHaveValue('test value');
  });

  it('triggers on setFilter on blur', async () => {
    const setFilter = jest.fn();
    render(
      <TextBoxFilter
        column={{
          filterValue: '',
          setFilter,
        }}
      />,
    );

    const textbox = screen.getByRole('textbox');
    await userEvent.type(textbox, 'Test Value');
    expect(setFilter).not.toBeCalled();
    await userEvent.tab();
    expect(setFilter).toBeCalledWith('Test Value');
  });

  it('triggers on setFilter on enter', async () => {
    const setFilter = jest.fn();
    render(
      <TextBoxFilter
        column={{
          filterValue: '',
          setFilter,
        }}
      />,
    );

    const textbox = screen.getByRole('textbox');
    await userEvent.type(textbox, 'Test Value{enter}');
    expect(setFilter).toBeCalledWith('Test Value');
  });

  it('triggers setFilter with undefined given empty input', async () => {
    const setFilter = jest.fn();
    render(
      <TextBoxFilter
        column={{
          filterValue: '',
          setFilter,
        }}
      />,
    );

    const textbox = screen.getByRole('textbox');
    await userEvent.click(textbox);
    expect(setFilter).not.toBeCalled();
    await userEvent.tab();
    expect(setFilter).toBeCalledWith(undefined);
  });
});
