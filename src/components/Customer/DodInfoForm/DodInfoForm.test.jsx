import React from 'react';
import { render, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DodInfoForm from './DodInfoForm';

import { isBooleanFlagEnabled } from 'utils/featureFlags';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

describe('DodInfoForm component', () => {
  const testProps = {
    onSubmit: jest.fn().mockImplementation(() => Promise.resolve()),
    initialValues: { affiliation: '', edipi: '1234567890' },
    onBack: jest.fn(),
  };

  const coastGuardTestProps = {
    ...testProps,
    initialValues: { affiliation: 'COAST_GUARD', edipi: '6546546541' },
  };

  it('renders the form inputs', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(true));
    const { getByLabelText } = render(<DodInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText(/Branch of service/)).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText(/Branch of service/)).toBeRequired();

      expect(getByLabelText(/DOD ID number/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/DOD ID number/)).toBeDisabled();
    });
  });

  it('renders the form inputs but enables editing of DOD ID when flag is on', async () => {
    isBooleanFlagEnabled.mockImplementation(() => Promise.resolve(false));
    const { getByLabelText } = render(<DodInfoForm {...testProps} />);

    await waitFor(() => {
      expect(getByLabelText(/Branch of service/)).toBeInstanceOf(HTMLSelectElement);
      expect(getByLabelText(/Branch of service/)).toBeRequired();

      expect(getByLabelText(/DOD ID number/)).toBeInstanceOf(HTMLInputElement);
      expect(getByLabelText(/DOD ID number/)).toBeEnabled();
    });
  });

  it('shows an error message if trying to submit an invalid form', async () => {
    const { getByRole, getAllByTestId, getByLabelText } = render(<DodInfoForm {...testProps} />);
    await userEvent.click(getByLabelText(/Branch of service/));
    await userEvent.click(getByLabelText(/DOD ID number/));

    const submitBtn = getByRole('button', { name: 'Next' });
    await userEvent.click(submitBtn);

    await waitFor(() => {
      expect(getAllByTestId('errorMessage').length).toBe(1);
      expect(submitBtn).toBeDisabled();
    });
    expect(testProps.onSubmit).not.toHaveBeenCalled();
  });

  it('submits the form when its valid', async () => {
    const { getByRole, getByLabelText } = render(<DodInfoForm {...testProps} />);
    const submitBtn = getByRole('button', { name: 'Next' });

    await userEvent.selectOptions(getByLabelText(/Branch of service/), ['NAVY']);
    await userEvent.type(getByLabelText(/DOD ID number/), '1234567890');

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

  describe('Coast Guard Customers', () => {
    it('shows an error message if EMPLID not present ', async () => {
      const { getByRole, getAllByTestId, getByLabelText } = render(<DodInfoForm {...coastGuardTestProps} />);
      await userEvent.click(getByLabelText(/Branch of service/));
      await userEvent.click(getByLabelText(/DOD ID number/));
      await userEvent.click(getByLabelText(/EMPLID/));

      const submitBtn = getByRole('button', { name: 'Next' });
      await userEvent.click(submitBtn);

      await waitFor(() => {
        expect(getAllByTestId('errorMessage').length).toBe(1);
        expect(submitBtn).toBeDisabled();
      });
      expect(testProps.onSubmit).not.toHaveBeenCalled();
    });

    it('submits the form when its valid', async () => {
      const { getByRole, getByLabelText } = render(<DodInfoForm {...testProps} />);
      const submitBtn = getByRole('button', { name: 'Next' });

      await userEvent.selectOptions(getByLabelText(/Branch of service/), ['COAST_GUARD']);
      await userEvent.type(getByLabelText(/DOD ID number/), '1234567890');
      await userEvent.type(getByLabelText(/EMPLID/), '1234567');

      await userEvent.click(submitBtn);

      await waitFor(() => {
        expect(testProps.onSubmit).toHaveBeenCalledWith(
          expect.objectContaining({ affiliation: 'COAST_GUARD', edipi: '1234567890', emplid: '1234567' }),
          expect.anything(),
        );
      });
    });
  });

  afterEach(jest.resetAllMocks);
});
