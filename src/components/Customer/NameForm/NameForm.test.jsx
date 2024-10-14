import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import NameForm from './NameForm';

describe('NameForm component', () => {
  it('renders the form inputs', async () => {
    const { getByLabelText } = render(
      <NameForm
        onSubmit={jest.fn()}
        onBack={jest.fn()}
        initialValues={{ first_name: '', middle_name: '', last_name: '', suffix: '' }}
      />,
    );
    await waitFor(() => {
      expect(getByLabelText(/First name/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/First name/)).toBeRequired();

      expect(getByLabelText(/Middle name/)).toBeInstanceOf(HTMLInputElement);

      expect(getByLabelText(/Last name/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/Last name/)).toBeRequired();

      expect(getByLabelText(/Suffix/)).toBeInstanceOf(HTMLInputElement);
    });
  });

  it('shows an error message and disables submit when fields are invalid', async () => {
    const onSubmit = jest.fn();
    const { getByRole, getAllByTestId, getByLabelText } = render(
      <NameForm
        onSubmit={onSubmit}
        onBack={jest.fn()}
        initialValues={{ first_name: '', middle_name: '', last_name: '', suffix: 'Mrs.' }}
      />,
    );
    await userEvent.clear(getByLabelText(/First name/));
    await userEvent.clear(getByLabelText(/Last name/));
    await userEvent.tab();

    const submitBtn = getByRole('button', { name: 'Next' });
    await waitFor(() => {
      expect(submitBtn).not.toBeEnabled();
    });

    await waitFor(() => {
      expect(getAllByTestId('errorMessage').length).toBe(2);
    });
    expect(onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const onSubmit = jest.fn();
    const { getByRole, getByLabelText } = render(
      <NameForm
        onSubmit={onSubmit}
        onBack={jest.fn()}
        initialValues={{ first_name: '', middle_name: '', last_name: '', suffix: '' }}
      />,
    );
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.type(getByLabelText(/First name/), 'Leo');
    await userEvent.type(getByLabelText(/Last name/), 'Spaceman');

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalled();
    });
  });

  it('uses the onBack handler when the back button is clicked', async () => {
    const onBack = jest.fn();
    const { getByRole } = render(
      <NameForm
        onSubmit={jest.fn()}
        onBack={onBack}
        initialValues={{ first_name: '', middle_name: '', last_name: '', suffix: 'Miss.' }}
      />,
    );
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
