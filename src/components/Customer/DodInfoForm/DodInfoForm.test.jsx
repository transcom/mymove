import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DodInfoForm from './DodInfoForm';

describe('DodInfoForm component', () => {
  const testProps = {
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    initialValues: { affiliation: '', edipi: '', rank: '' },
    onBack: jest.fn(),
  };

  it('renders the form inputs', async () => {
    const { getByLabelText } = render(<DodInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText('Branch of service')).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText('Branch of service')).toBeRequired();

      expect(getByLabelText('DOD ID number')).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText('DOD ID number')).toBeRequired();

      expect(getByLabelText('Rank')).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText('Rank')).toBeRequired();
    });
  });

  it('validates the DOD ID number on blur', async () => {
    const { getByLabelText, getByText } = render(<DodInfoForm {...testProps} />);

    userEvent.type(getByLabelText('DOD ID number'), 'not a valid ID number');
    userEvent.tab();

    await waitFor(() => {
      expect(getByLabelText('DOD ID number')).not.toBeValid();
      expect(getByText('Enter a 10-digit DOD ID number')).toBeInTheDocument();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, getAllByText } = render(<DodInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByText('Required').length).toBe(3);
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<DodInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    userEvent.selectOptions(getByLabelText('Branch of service'), ['NAVY']);
    userEvent.type(getByLabelText('DOD ID number'), '1234567890');
    userEvent.selectOptions(getByLabelText('Rank'), ['E_5']);

    userEvent.click(submitBtn);

    await waitFor(() => {
      expect(testProps.onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({ affiliation: 'NAVY', edipi: '1234567890', rank: 'E_5' }),
        expect.anything(),
      );
    });
  });

  it('implements the onBack handler when the Back button is clicked', async () => {
    const { getByRole } = render(<DodInfoForm {...testProps} />);
    const backBtn = getByRole('button', { name: 'Back' });

    userEvent.click(backBtn);

    await waitFor(() => {
      expect(testProps.onBack).toHaveBeenCalled();
    });
  });

  afterEach(jest.resetAllMocks);
});
