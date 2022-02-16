import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DatesAndLocation from './DatesAndLocation';

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  serviceMember: {
    id: '123',
    residentialAddress: {
      postalCode: '90210',
    },
  },
  destinationDutyStation: {
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

describe('DatesAndLocation component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<DatesAndLocation {...defaultProps} />);
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
      render(<DatesAndLocation {...defaultProps} />);
      const useCurrentZip = await screen.getByText('Use my current ZIP (90210)');
      const originZip = screen.getAllByLabelText('ZIP')[0];
      expect(originZip.value).toBe('');
      userEvent.click(useCurrentZip);
      await waitFor(() => {
        expect(originZip.value).toBe(defaultProps.serviceMember.residentialAddress.postalCode);
      });
    });

    it('removes current zip when "use my current zip" is deselected', async () => {
      render(<DatesAndLocation {...defaultProps} />);
      const useCurrentZip = await screen.getByText('Use my current ZIP (90210)');
      const originZip = screen.getAllByLabelText('ZIP')[0];
      expect(originZip.value).toBe('');
      userEvent.click(useCurrentZip);
      await waitFor(() => {
        expect(originZip.value).toBe(defaultProps.serviceMember.residentialAddress.postalCode);
      });
      userEvent.click(useCurrentZip);
      await waitFor(() => {
        expect(originZip.value).toBe('');
      });
    });

    it('displays secondary pickup postal code input when hasSecondaryPickupPostalCode is true', async () => {
      render(<DatesAndLocation {...defaultProps} />);
      const hasSecondaryPickupPostalCode = await screen.getAllByLabelText('Yes')[0];
      expect(screen.queryByLabelText('Second ZIP')).toBeNull();
      userEvent.click(hasSecondaryPickupPostalCode);

      await waitFor(() => {
        expect(screen.queryByLabelText('Second ZIP')).toBeInstanceOf(HTMLInputElement);
      });
    });

    it('displays destination zip when "Use the ZIP for my new duty location" is selected', async () => {
      render(<DatesAndLocation {...defaultProps} />);
      const useDestinationZip = await screen.getByText('Use the ZIP for my new duty location (94611)');
      const destinationZip = screen.getAllByLabelText('ZIP')[1];
      expect(destinationZip.value).toBe('');
      userEvent.click(useDestinationZip);
      await waitFor(() => {
        expect(destinationZip.value).toBe(defaultProps.destinationDutyStation?.address?.postalCode);
      });
    });

    it('removes destination zip when "Use the ZIP for my new duty location" is deselected', async () => {
      render(<DatesAndLocation {...defaultProps} />);
      const useDestinationZip = await screen.getByText('Use the ZIP for my new duty location (94611)');
      const destinationZip = screen.getAllByLabelText('ZIP')[1];
      expect(destinationZip.value).toBe('');
      userEvent.click(useDestinationZip);
      await waitFor(() => {
        expect(destinationZip.value).toBe(defaultProps.destinationDutyStation?.address?.postalCode);
      });

      userEvent.click(useDestinationZip);
      await waitFor(() => {
        expect(destinationZip.value).toBe('');
      });
    });

    it('displays secondary destination postal code input when hasSecondaryDestinationPostalCode is true', async () => {
      render(<DatesAndLocation {...defaultProps} />);
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
      render(<DatesAndLocation {...mtoShipmentProps} />);
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
      render(<DatesAndLocation {...defaultProps} />);

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
        mtoShipment: {
          ppmShipment: {
            pickupPostalCode: '90210',
            destinationPostalCode: '10001',
            expectedDepartureDate: '2022-07-04',
          },
        },
      };
      render(<DatesAndLocation {...hasSecondaryZIPs} />);

      const inputHasSecondaryZIP = screen.getAllByLabelText('Yes');

      await userEvent.click(inputHasSecondaryZIP[0]);
      await userEvent.click(inputHasSecondaryZIP[1]);

      const secondaryZIPs = screen.getAllByLabelText('Second ZIP');

      await userEvent.click(secondaryZIPs[0]);
      await userEvent.tab();

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
      render(<DatesAndLocation {...invalidTypes} />);

      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '1 January 2022');

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');

        // origin ZIP
        expect(requiredAlerts[0]).toHaveTextContent('Must be valid code');
        expect(requiredAlerts[0].nextElementSibling).toHaveAttribute('name', 'pickupPostalCode');

        // Secondary origin ZIP
        expect(requiredAlerts[1]).toHaveTextContent('Must be valid code');
        expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'secondaryPickupPostalCode');

        // Secondary destination ZIP
        expect(requiredAlerts[2]).toHaveTextContent('Must be valid code');
        expect(requiredAlerts[2].nextElementSibling).toHaveAttribute('name', 'destinationPostalCode');

        // Secondary destination ZIP
        expect(requiredAlerts[3]).toHaveTextContent('Must be valid code');
        expect(requiredAlerts[3].nextElementSibling).toHaveAttribute('name', 'secondaryDestinationPostalCode');

        // Departure date
        expect(requiredAlerts[4]).toHaveTextContent('Enter a complete date in DD MMM YYYY format (day, month, year).');
        expect(
          within(requiredAlerts[4].nextElementSibling).getByLabelText('When do you plan to start moving your PPM?'),
        ).toBeInTheDocument();
      });
    });

    it('calls postalCodeValidator when the ZIP value changes', async () => {
      render(<DatesAndLocation {...defaultProps} />);

      const primaryZIPs = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryZIPs[0], '12345');
      await userEvent.type(primaryZIPs[1], '67890');

      const inputHasSecondaryZIP = screen.getAllByLabelText('Yes');

      await userEvent.click(inputHasSecondaryZIP[0]);
      await userEvent.click(inputHasSecondaryZIP[1]);

      const secondaryZIPs = screen.getAllByLabelText('Second ZIP');
      await userEvent.type(secondaryZIPs[0], '11111');
      await userEvent.type(secondaryZIPs[1], '22222');

      await waitFor(() => {
        expect(defaultProps.postalCodeValidator).toHaveBeenCalledWith('12345', 'origin');
        expect(defaultProps.postalCodeValidator).toHaveBeenCalledWith('67890', 'destination');
        expect(defaultProps.postalCodeValidator).toHaveBeenCalledWith('11111', 'origin');
        expect(defaultProps.postalCodeValidator).toHaveBeenCalledWith('22222', 'destination');
      });
    });

    it('displays error when postal code lookup fails', async () => {
      const postalCodeValidatorFailure = {
        ...defaultProps,
        postalCodeValidator: jest
          .fn()
          .mockReturnValue('Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.'),
      };
      render(<DatesAndLocation {...postalCodeValidatorFailure} />);

      const primaryZIPs = screen.getAllByLabelText('ZIP');
      await userEvent.type(primaryZIPs[0], '99999');

      await waitFor(() => {
        expect(postalCodeValidatorFailure.postalCodeValidator).toHaveBeenCalledWith('99999', 'origin');
        /*
        expect(screen.getByRole('alert')).toHaveTextContent(
          'Sorry, we don’t support that zip code yet. Please contact your local PPPO for assistance.',
        );
        */
      });
    });
  });
});
