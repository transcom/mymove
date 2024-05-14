import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DodInfoForm from './DodInfoForm';

describe('DodInfoForm component', () => {
  const testProps = {
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    initialValues: { affiliation: '', edipi: '' },
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<DodInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Branch of service')).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText('Branch of service')).toBeRequired();

      expect(getByLabelText('DOD ID number')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('DOD ID number')).toBeRequired();
    });
  });

  it('validates the DOD ID number on blur', async () => {
    const { getByLabelText, getByText } = render(<DodInfoForm {...testProps} />);

    await userEvent.type(getByLabelText('DOD ID number'), 'not a valid ID number');
    await userEvent.tab();

    await waitFor(() => {
      expect(getByLabelText('DOD ID number')).not.toBeValid();
      expect(getByText('Enter a 10-digit DOD ID number')).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, getAllByText, getByLabelText } = render(<DodInfoForm {...testProps} />);
    await userEvent.click(getByLabelText('Branch of service'));
    await userEvent.click(getByLabelText('DOD ID number'));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(2);
      expect(submitBtn).toBeDisabled();
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<DodInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.selectOptions(getByLabelText('Branch of service'), ['NAVY']);
    await userEvent.type(getByLabelText('DOD ID number'), '1234567890');

    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({ affiliation: 'NAVY', edipi: '1234567890' }),
        expect.anything(),
      );
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<DodInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    await userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
