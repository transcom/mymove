import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import { UnsupportedZipCodePPMErrorMsg } from 'utils/validation';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';

const serviceMember = {
  serviceMember: {
    id: '123',
    current_location: {
      name: 'Fort Drum',
    },
    residential_address: {
      city: 'Fort Benning',
      state: 'GA',
      postalCode: '90210',
      streetAddress1: '123 Main',
      streetAddress2: '',
    },
    affiliation: SERVICE_MEMBER_AGENCIES.ARMY,
  },
};

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  destinationDutyLocation: {
    address: {
      city: 'Fort Benning',
      state: 'GA',
      postalCode: '94611',
      streetAddress1: '123 Main',
      streetAddress2: '',
    },
  },
  currentResidence: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  originDutyLocationAddress: {
    city: 'Fort Benning',
    state: 'GA',
    postalCode: '31905',
    streetAddress1: '123 Main',
    streetAddress2: '',
  },
  postalCodeValidator: jest.fn(),
  ...serviceMember,
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
      expect(screen.getByName('serviceMember.residential_address.streetAddress1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.residential_address.streetAddress2')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.residential_address.streetAddress3')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.residential_address.state')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByName('serviceMember.residential_address.city')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.residential_address.postalCode')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Destination' })).toBeInTheDocument();
      expect(screen.getByName('serviceMember.destination_address.streetAddress1')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.destination_address.streetAddress2')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.destination_address.streetAddress3')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.destination_address.state')).toBeInstanceOf(HTMLSelectElement);
      expect(screen.getByName('serviceMember.destination_address.city')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByName('serviceMember.destination_address.postalCode')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Closeout Office' })).toBeInTheDocument();
      expect(screen.getByLabelText('Which closeout office should review your PPM?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Storage' })).toBeInTheDocument();
      expect(screen.getAllByLabelText('Yes')[2]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[2]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      expect(screen.getByLabelText('When do you plan to start moving your PPM?')).toBeInstanceOf(HTMLInputElement);
    });
  });

  describe('displays conditional inputs', () => {
    it('displays current address when "use my current address" is selected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      await user.click(screen.getByLabelText('Use current address'));

      expect((await screen.getByName('serviceMember.residential_address.streetAddress1'))).toHaveValue(
        defaultProps.currentResidence.streetAddress1,
      );
      expect(screen.getByName('serviceMember.residential_address.streetAddress2')).toHaveValue(defaultProps.currentResidence.streetAddress2);
      expect(screen.getByName('serviceMember.residential_address.state')).toHaveValue(defaultProps.currentResidence.state);
      expect(screen.getByName('serviceMember.residential_address.city')).toHaveValue(defaultProps.currentResidence.city);
      expect(screen.getByName('serviceMember.residential_address.postalCode')).toHaveValue(defaultProps.currentResidence.postalCode);
    });

    it('removes current Address when "use my current Address" is deselected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      await user.click(screen.getByLabelText('Use current address'));

      expect((await screen.getByName('serviceMember.residential_address.streetAddress1'))).toHaveValue('');
      expect(screen.getByName('serviceMember.residential_address.streetAddress2')).toHaveValue('');
      expect(screen.getByName('serviceMember.residential_address.state')).toHaveValue('');
      expect(screen.getByName('serviceMember.residential_address.city')).toHaveValue('');
      expect(screen.getByName('serviceMember.residential_address.postalCode')).toHaveValue('');
    });

    it('displays secondary pickup Address input when hasSecondaryPickupAddress is true', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const hasSecondaryPickupAddress = await screen.getAllByLabelText('Yes')[1];

      await userEvent.click(hasSecondaryPickupAddress);

      await waitFor(() => {
        expect(screen.getByName('mtoShipment.secondaryPickupAddress.streetAddress1')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryPickupAddress.streetAddress2')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryPickupAddress.city')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryPickupAddress.state')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryPickupAddress.postalCode')).toBeInstanceOf(HTMLInputElement);
      });
    });

    it('displays destination address when "Use my current destination address" is selected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      await user.click(screen.getByLabelText('Use my current destination address'));
      expect((await screen.getByName('serviceMember.destination_address.streetAddress1'))).toHaveValue(
        defaultProps.destinationDutyLocation.address.streetAddress1,
      );
      expect(screen.getByName('serviceMember.destination_address.streetAddress2')).toHaveValue('');
      expect(screen.getByName('serviceMember.destination_address.city')).toHaveValue(defaultProps.destinationDutyLocation.address.city);
      expect(screen.getByName('serviceMember.destination_address.state')).toHaveValue(defaultProps.destinationDutyLocation.address.state);
      expect(screen.getByName('serviceMember.destination_address.postalCode')).toHaveValue(defaultProps.destinationDutyLocation.address.postalCode);
      });
    });

    it('removes destination Address when "Use my current destination address" is deselected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      await user.click(screen.getByLabelText('Use my current destination address'));

      expect((await screen.getByName('mtoShipment.secondaryDestinationAddress.streetAddress1'))).toHaveValue('');
      expect(screen.getByName('mtoShipment.secondaryDestinationAddress.streetAddress2')).toHaveValue('');
      expect(screen.getByName('mtoShipment.secondaryDestinationAddress.city')).toHaveValue('');
      expect(screen.getByName('mtoShipment.secondaryDestinationAddress.state')).toHaveValue('');
      expect(screen.getByName('mtoShipment.secondaryDestinationAddress.postalCode')).toHaveValue('');
    });

    it('displays secondary destination Address input when hasSecondaryDestinationAddress is true', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const hasSecondaryDestinationAddress = await screen.getAllByLabelText('Yes')[1];
      
      await userEvent.click(hasSecondaryDestinationAddress);

      await waitFor(() => {
        expect(screen.getByName('mtoShipment.secondaryDestinationAddress.streetAddress1')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryDestinationAddress.streetAddress2')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryDestinationAddress.city')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryDestinationAddress.state')).toBeInstanceOf(HTMLInputElement);
        expect(screen.getByName('mtoShipment.secondaryDestinationAddress.postalCode')).toBeInstanceOf(HTMLInputElement);
      });
    });

    it('displays the closeout office select when the service member is in the Army', async () => {
      const armyServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.ARMY,
      };
      render(<DateAndLocationForm {...defaultProps} serviceMember={armyServiceMember} />);

      expect(screen.getByText('Closeout Office')).toBeInTheDocument();
      expect(screen.getByLabelText('Which closeout office should review your PPM?')).toBeInTheDocument();
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
    });

    it('displays the closeout office select when the service member is in the Air Force', async () => {
      const airForceServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.AIR_FORCE,
      };
      render(<DateAndLocationForm {...defaultProps} serviceMember={airForceServiceMember} />);

      expect(screen.getByText('Closeout Office')).toBeInTheDocument();
      expect(screen.getByLabelText('Which closeout office should review your PPM?')).toBeInTheDocument();
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
    });

    it('does not display the closeout office select when the service member is not in the Army/Air-Force', async () => {
      const navyServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.NAVY,
      };

      render(<DateAndLocationForm {...defaultProps} serviceMember={navyServiceMember} />);
      expect(screen.queryByText('Closeout Office')).not.toBeInTheDocument();
      expect(screen.queryByLabelText('Which closeout office should review your PPM?')).not.toBeInTheDocument();
      expect(screen.queryByText('Start typing a closeout office...')).not.toBeInTheDocument();
    });
  });

  describe('pull values from the ppm shipment when available', () => {
    it('renders blank form on load', async () => {
      render(<DateAndLocationForm {...mtoShipmentProps} />);
      expect(await screen.getAllByLabelText('ZIP')[0].value).toHaveValue(
        mtoShipmentProps.mtoShipment.ppmShipment.pickupPostalCode,
      );
      expect(screen.getAllByLabelText('Yes')[0].value).toBe('true');
      expect(screen.getAllByLabelText('ZIP')[2].value).toHaveValue(
        mtoShipmentProps.mtoShipment.ppmShipment.destinationPostalCode,
      );
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
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
        expect(requiredAlerts[0].nextElementSibling).toHaveAttribute('name', 'serviceMember.residential_address.postalCode');

        // Destination ZIP
        expect(requiredAlerts[1]).toHaveTextContent('Required');
        expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'serviceMember.destination_address.postalCode');

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

      const secondaryZIPs = screen.getAllByLabelText('ZIP');

      await userEvent.click(secondaryZIPs[1]);
      await userEvent.click(secondaryZIPs[3]);
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

      const zipInputs = screen.getAllByLabelText('ZIP');
      await userEvent.click(zipInputs[0]);
      await userEvent.click(zipInputs[1]);
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');
        expect(requiredAlerts.length).toBe(5);

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

      await userEvent.type(primaryZIPs[0], '12345');
      await userEvent.type(primaryZIPs[1], '67890');

      const inputHasSecondaryZIP = screen.getAllByLabelText('Yes');

      await userEvent.click(inputHasSecondaryZIP[0]);
      await userEvent.click(inputHasSecondaryZIP[1]);

      const secondaryZIPs = screen.getAllByLabelText('ZIP');
      await userEvent.type(secondaryZIPs[1], '11111');
      await userEvent.type(secondaryZIPs[3], '22222');

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
      await userEvent.type(primaryZIPs[0], '99999');

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
;
