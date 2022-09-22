import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  serviceMember: {
    id: '123',
    residential_address: {
      postalCode: '90210',
    },
  },
  destinationDutyLocation: {
    address: {
      postalCode: '94611',
    },
  },
  postalCodeValidator: jest.fn(),
};

const mtoShipmentProps = {
  ...defaultProps,
  mtoShipment: {
    id: '123',
    ppmShipment: {
      id: '123',
      pickupPostalCode: '12345',
      secondaryPickupPostalCode: '34512',
      destinationPostalCode: '94611',
      secondaryDestinationPostalCode: '90210',
      sitExpected: true,
      expectedDepartureDate: '2022-09-23',
    },
  },
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe('DateAndLocationForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      expect(await screen.getByRole('heading', { level: 2, name: 'Origin' })).toBeInTheDocument();
      expect(screen.getAllByLabelText('ZIP')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Destination' })).toBeInTheDocument();
      expect(screen.getAllByLabelText('ZIP')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Storage' })).toBeInTheDocument();
      expect(screen.getAllByLabelText('Yes')[2]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[2]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      expect(screen.getByLabelText('When do you plan to start moving your PPM?')).toBeInstanceOf(HTMLInputElement);
    });
  });

  describe('displays conditional inputs', () => {
    it('displays current zip when "use my current zip" is selected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const useCurrentZip = await screen.getByText('Use my current ZIP (90210)');
      const originZip = screen.getAllByLabelText('ZIP')[0];
      expect(originZip.value).toBe('');
      userEvent.click(useCurrentZip);
      await waitFor(() => {
        expect(originZip.value).toBe(defaultProps.serviceMember.residential_address.postalCode);
      });
    });

    it('removes current zip when "use my current zip" is deselected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const useCurrentZip = await screen.getByText('Use my current ZIP (90210)');
      const originZip = screen.getAllByLabelText('ZIP')[0];
      expect(originZip.value).toBe('');
      userEvent.click(useCurrentZip);
      await waitFor(() => {
        expect(originZip.value).toBe(defaultProps.serviceMember.residential_address.postalCode);
      });
      userEvent.click(useCurrentZip);
      await waitFor(() => {
        expect(originZip.value).toBe('');
      });
    });

    it('displays secondary pickup postal code input when hasSecondaryPickupPostalCode is true', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const hasSecondaryPickupPostalCode = await screen.getAllByLabelText('Yes')[0];
      expect(screen.queryByLabelText('Second ZIP')).toBeNull();
      userEvent.click(hasSecondaryPickupPostalCode);

      await waitFor(() => {
        expect(screen.queryByLabelText('Second ZIP')).toBeInstanceOf(HTMLInputElement);
      });
    });

    it('displays destination zip when "Use the ZIP for my new duty location" is selected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const useDestinationZip = await screen.getByText('Use the ZIP for my new duty location (94611)');
      const destinationZip = screen.getAllByLabelText('ZIP')[1];
      expect(destinationZip.value).toBe('');
      userEvent.click(useDestinationZip);
      await waitFor(() => {
        expect(destinationZip.value).toBe(defaultProps.destinationDutyLocation?.address?.postalCode);
      });
    });

    it('removes destination zip when "Use the ZIP for my new duty location" is deselected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const useDestinationZip = await screen.getByText('Use the ZIP for my new duty location (94611)');
      const destinationZip = screen.getAllByLabelText('ZIP')[1];
      expect(destinationZip.value).toBe('');
      userEvent.click(useDestinationZip);
      await waitFor(() => {
        expect(destinationZip.value).toBe(defaultProps.destinationDutyLocation?.address?.postalCode);
      });

      userEvent.click(useDestinationZip);
      await waitFor(() => {
        expect(destinationZip.value).toBe('');
      });
    });

    it('displays secondary destination postal code input when hasSecondaryDestinationPostalCode is true', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const hasSecondaryDestinationPostalCode = await screen.getAllByLabelText('Yes')[0];
      expect(screen.queryByLabelText('Second ZIP')).toBeNull();
      userEvent.click(hasSecondaryDestinationPostalCode);

      await waitFor(() => {
        expect(screen.queryByLabelText('Second ZIP')).toBeInstanceOf(HTMLInputElement);
      });
    });
  });

  describe('pull values from the ppm shipment when available', () => {
    it('renders blank form on load', async () => {
      render(<DateAndLocationForm {...mtoShipmentProps} />);
      expect(await screen.getAllByLabelText('ZIP')[0].value).toBe(
        mtoShipmentProps.mtoShipment.ppmShipment.pickupPostalCode,
      );
      expect(screen.getAllByLabelText('Yes')[0].value).toBe('true');
      expect(screen.getAllByLabelText('ZIP')[1].value).toBe(
        mtoShipmentProps.mtoShipment.ppmShipment.destinationPostalCode,
      );
      expect(screen.getAllByLabelText('Yes')[1].value).toBe('true');
      expect(screen.getAllByLabelText('Yes')[2].value).toBe('true');
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks required inputs when left empty', async () => {
      render(<DateAndLocationForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');

        // Origin ZIP
        expect(requiredAlerts[0]).toHaveTextContent('Required');
        expect(requiredAlerts[0].nextElementSibling).toHaveAttribute('name', 'pickupPostalCode');

        // Destination ZIP
        expect(requiredAlerts[1]).toHaveTextContent('Required');
        expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'destinationPostalCode');

        // Departure date
        expect(requiredAlerts[2]).toHaveTextContent('Required');
        expect(
          within(requiredAlerts[2].nextElementSibling).getByLabelText('When do you plan to start moving your PPM?'),
        ).toBeInTheDocument();
      });
    });
    it('marks secondary ZIP fields as required when conditionally displayed', async () => {
      const hasSecondaryZIPs = {
        ...defaultProps,
        postalCodeValidator: jest.fn(),
      };
      render(<DateAndLocationForm {...hasSecondaryZIPs} />);

      const inputHasSecondaryZIP = screen.getAllByLabelText('Yes');

      await userEvent.click(inputHasSecondaryZIP[0]);
      await userEvent.click(inputHasSecondaryZIP[1]);

      const secondaryZIPs = screen.getAllByLabelText('Second ZIP');

      await userEvent.click(secondaryZIPs[0]);

      await userEvent.click(secondaryZIPs[1]);
      await userEvent.tab();

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');

        // Secondary origin ZIP
        expect(requiredAlerts[0]).toHaveTextContent('Required');
        expect(requiredAlerts[0].nextElementSibling).toHaveAttribute('name', 'secondaryPickupPostalCode');

        // Secondary destination ZIP
        expect(requiredAlerts[1]).toHaveTextContent('Required');
        expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'secondaryDestinationPostalCode');
      });
    });
    it('displays type errors when input fails validation schema', async () => {
      const invalidTypes = {
        ...defaultProps,
        mtoShipment: {
          ppmShipment: {
            pickupPostalCode: '1000',
            secondaryPickupPostalCode: '2000',
            destinationPostalCode: '3000',
            secondaryDestinationPostalCode: '4000',
          },
        },
      };
      render(<DateAndLocationForm {...invalidTypes} />);

      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '1 January 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');

        // origin ZIP
        expect(requiredAlerts[0]).toHaveTextContent('Enter a 5-digit ZIP code');
        expect(requiredAlerts[0].nextElementSibling).toHaveAttribute('name', 'pickupPostalCode');

        // Secondary origin ZIP
        expect(requiredAlerts[1]).toHaveTextContent('Enter a 5-digit ZIP code');
        expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'secondaryPickupPostalCode');

        // Secondary destination ZIP
        expect(requiredAlerts[2]).toHaveTextContent('Enter a 5-digit ZIP code');
        expect(requiredAlerts[2].nextElementSibling).toHaveAttribute('name', 'destinationPostalCode');

        // Secondary destination ZIP
        expect(requiredAlerts[3]).toHaveTextContent('Enter a 5-digit ZIP code');
        expect(requiredAlerts[3].nextElementSibling).toHaveAttribute('name', 'secondaryDestinationPostalCode');

        // Departure date
        expect(requiredAlerts[4]).toHaveTextContent('Enter a complete date in DD MMM YYYY format (day, month, year).');
        expect(
          within(requiredAlerts[4].nextElementSibling).getByLabelText('When do you plan to start moving your PPM?'),
        ).toBeInTheDocument();
      });
    });
    it('calls postalCodeValidator when the ZIP value changes', async () => {
      const validatorProps = {
        ...defaultProps,
        postalCodeValidator: jest.fn(),
      };
      render(<DateAndLocationForm {...validatorProps} />);
      const primaryZIPs = screen.getAllByLabelText('ZIP');

      userEvent.type(primaryZIPs[0], '12345');

      userEvent.type(primaryZIPs[1], '67890');

      const inputHasSecondaryZIP = screen.getAllByLabelText('Yes');

      userEvent.click(inputHasSecondaryZIP[0]);
      userEvent.click(inputHasSecondaryZIP[1]);

      const secondaryZIPs = screen.getAllByLabelText('Second ZIP');
      userEvent.type(secondaryZIPs[0], '11111');
      userEvent.type(secondaryZIPs[1], '22222');

      await waitFor(async () => {
        expect(validatorProps.postalCodeValidator).toHaveBeenCalledWith(
          '12345',
          'origin',
          UnsupportedZipCodePPMErrorMsg,
        );
        expect(validatorProps.postalCodeValidator).toHaveBeenCalledWith(
          '67890',
          'destination',
          UnsupportedZipCodePPMErrorMsg,
        );
        expect(validatorProps.postalCodeValidator).toHaveBeenCalledWith(
          '11111',
          'origin',
          UnsupportedZipCodePPMErrorMsg,
        );
        expect(validatorProps.postalCodeValidator).toHaveBeenCalledWith(
          '22222',
          'destination',
          UnsupportedZipCodePPMErrorMsg,
        );
      });
    });

    it('displays error when postal code lookup fails', async () => {
      const postalCodeValidatorFailure = {
        ...defaultProps,
        postalCodeValidator: jest
          .fn()
          .mockReturnValue('Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.'),
      };
      render(<DateAndLocationForm {...postalCodeValidatorFailure} />);

      const primaryZIPs = screen.getAllByLabelText('ZIP');
      userEvent.type(primaryZIPs[0], '99999');

      await waitFor(() => {
        expect(postalCodeValidatorFailure.postalCodeValidator).toHaveBeenCalledWith(
          '99999',
          'origin',
          UnsupportedZipCodePPMErrorMsg,
        );
        /*
        expect(screen.getByRole('alert')).toHaveTextContent(
          'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.',
        );
       */
      });
    });
  });
});
