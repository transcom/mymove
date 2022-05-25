import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AboutForm from 'components/Customer/PPM/Closeout/AboutForm/AboutForm';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  postalCodeValidator: jest.fn(),
  mtoShipment: {
    ppmShipment: {},
  },
};

const mtoShipmentProps = {
  mtoShipment: {
    ppmShipment: {
      actualMoveDate: '2022-05-19',
      actualPickupPostalCode: '10001',
      actualDestinationPostalCode: '60652',
      hasReceivedAdvance: true,
      advanceAmountReceived: 123456,
    },
  },
};

describe('AboutForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<AboutForm {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      });
      expect(screen.getByLabelText('When did you leave your origin?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Locations' })).toBeInTheDocument();
      expect(screen.getByLabelText('Starting ZIP')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('Ending ZIP')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeChecked(); // Has advance received is set to No by default
      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates all field values', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('When did you leave your origin?')).toHaveDisplayValue('19 May 2022');
      });
      expect(screen.getByLabelText('Starting ZIP')).toHaveDisplayValue('10001');
      expect(screen.getByLabelText('Ending ZIP')).toHaveDisplayValue('60652');
      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('How much did you receive?')).toHaveDisplayValue('1,234');
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks all required fields when form is submitted', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });

      let requiredAlerts = screen.getAllByRole('alert');

      expect(requiredAlerts[0]).toHaveTextContent('Required');
      expect(
        within(requiredAlerts[0].nextElementSibling).getByLabelText('When did you leave your origin?'),
      ).toBeInTheDocument();

      expect(requiredAlerts[1]).toHaveTextContent('Required');
      expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'actualPickupPostalCode');

      expect(requiredAlerts[2]).toHaveTextContent('Required');
      expect(requiredAlerts[2].nextElementSibling).toHaveAttribute('name', 'actualDestinationPostalCode');

      await userEvent.click(screen.getByLabelText('Yes'));

      await waitFor(() => {
        expect(screen.getByLabelText('How much did you receive?')).toBeInTheDocument();
      });

      requiredAlerts = screen.getAllByRole('alert');

      expect(requiredAlerts[3]).toHaveTextContent('Required');
      expect(
        within(requiredAlerts[3].nextElementSibling).getByLabelText('How much did you receive?'),
      ).toBeInTheDocument();
    });

    it('displays type error messages for invalid input', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.type(screen.getByLabelText('When did you leave your origin?'), '1 January 2022');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(
          'Enter a complete date in DD MMM YYYY format (day, month, year).',
        );
      });

      await userEvent.type(screen.getByLabelText('Starting ZIP'), '1');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getAllByRole('alert')[1]).toHaveTextContent('Enter a 5-digit ZIP code');
      });

      await userEvent.type(screen.getByLabelText('Ending ZIP'), '2');
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getAllByRole('alert')[2]).toHaveTextContent('Enter a 5-digit ZIP code');
      });
    });

    it('displays error when advance received is below 1 dollar minimum', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.click(screen.getByLabelText('Yes'));

      await userEvent.type(screen.getByLabelText('How much did you receive?'), '0');

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(
          "The minimum advance request is $1. If you don't want an advance, select No.",
        );
      });
    });

    it('calls the postal code validator for starting and ending ZIP inputs', async () => {
      const postalCodeValidatorProps = {
        postalCodeValidator: jest.fn().mockReturnValue(UnsupportedZipCodePPMErrorMsg),
      };
      render(<AboutForm {...defaultProps} {...postalCodeValidatorProps} />);

      userEvent.type(screen.getByLabelText('Starting ZIP'), '10000');

      await waitFor(() => {
        expect(postalCodeValidatorProps.postalCodeValidator).toHaveBeenCalledWith(
          '10000',
          'origin',
          UnsupportedZipCodePPMErrorMsg,
        );
      });

      userEvent.type(screen.getByLabelText('Ending ZIP'), '20000');

      await waitFor(() => {
        expect(postalCodeValidatorProps.postalCodeValidator).toHaveBeenCalledWith(
          '20000',
          'destination',
          UnsupportedZipCodePPMErrorMsg,
        );
      });
    });
  });

  describe('calls button event handlers', () => {
    it('calls onBack handler when "Finish Later" is pressed', async () => {
      render(<AboutForm {...defaultProps} />);

      userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });

    it('calls onSubmit handler when "Save & Continue" is pressed', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentProps} />);

      userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(
          {
            actualMoveDate: '2022-05-19',
            actualPickupPostalCode: '10001',
            actualDestinationPostalCode: '60652',
            hasReceivedAdvance: 'true',
            advanceAmountReceived: '1234',
          },
          expect.anything(),
        );
      });
    });
  });
});
