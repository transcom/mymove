import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { act } from 'react-dom/test-utils';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
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
  postalCodeValidator: jest.fn(),
  ...serviceMember,
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe('DateAndLocationForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      expect(await screen.getByRole('heading', { level: 2, name: 'Origin' })).toBeInTheDocument();
      const postalCodes = screen.getAllByLabelText('ZIP');
      const address1 = screen.getAllByLabelText('Address 1');
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const address3 = screen.getAllByLabelText('Address 3', { exact: false });
      const state = screen.getAllByLabelText('State');
      const city = screen.getAllByLabelText('City');

      expect(address1[0]).toBeInstanceOf(HTMLInputElement);
      expect(address2[0]).toBeInstanceOf(HTMLInputElement);
      expect(address3[0]).toBeInstanceOf(HTMLInputElement);
      expect(state[0]).toBeInstanceOf(HTMLSelectElement);
      expect(city[0]).toBeInstanceOf(HTMLInputElement);
      expect(postalCodes[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Destination' })).toBeInTheDocument();
      expect(address1[1]).toBeInstanceOf(HTMLInputElement);
      expect(address2[1]).toBeInstanceOf(HTMLInputElement);
      expect(address3[1]).toBeInstanceOf(HTMLInputElement);
      expect(state[1]).toBeInstanceOf(HTMLSelectElement);
      expect(city[1]).toBeInstanceOf(HTMLInputElement);
      expect(postalCodes[1]).toBeInstanceOf(HTMLInputElement);
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
    it('displays current address when "Use my current origin address" is selected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const postalCodes = screen.getAllByLabelText('ZIP');
      expect(postalCodes[0].value).toBe('');
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current origin address'));
      });
      await waitFor(() => {
        expect(postalCodes[0].value).toBe(defaultProps.serviceMember.residential_address.postalCode);
      });
    });

    it('removes current Address when "Use my current origin address" is deselected', async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current origin address'));
      });
      const postalCodes = screen.getAllByLabelText('ZIP');

      await waitFor(() => {
        expect(postalCodes[0].value).toBe(defaultProps.serviceMember.residential_address.postalCode);
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current origin address'));
      });

      await waitFor(() => {
        expect(postalCodes[0].value).toBe('');
      });
    });

    it('displays secondary pickup Address input when hasSecondaryPickupAddress is true', async () => {
      await act(async () => {
        render(<DateAndLocationForm {...defaultProps} />);
        const hasSecondaryPickupAddress = await screen.getAllByLabelText('Yes')[1];

        await userEvent.click(hasSecondaryPickupAddress);
        const postalCodes = screen.getAllByLabelText('ZIP');
        const address1 = screen.getAllByLabelText('Address 1', { exact: false });
        const address2 = screen.getAllByLabelText('Address 2', { exact: false });
        const state = screen.getAllByLabelText('State');
        const city = screen.getAllByLabelText('City');
        await waitFor(() => {
          expect(address1[1]).toBeInstanceOf(HTMLInputElement);
          expect(address2[1]).toBeInstanceOf(HTMLInputElement);
          expect(city[1]).toBeInstanceOf(HTMLInputElement);
          expect(state[1]).toBeInstanceOf(HTMLSelectElement);
          expect(postalCodes[1]).toBeInstanceOf(HTMLInputElement);
        });
      });
    });

    it('displays destination address when "Use my current destination address" is selected', async () => {
      await act(async () => {
        render(<DateAndLocationForm {...defaultProps} />);
        await userEvent.click(screen.getByLabelText('Use my current destination address'));
        const postalCodes = screen.getAllByLabelText('ZIP');
        const address1 = screen.getAllByLabelText('Address 1', { exact: false });
        const address2 = screen.getAllByLabelText('Address 2', { exact: false });
        const state = screen.getAllByLabelText('State');
        const city = screen.getAllByLabelText('City');
        expect(await address1[1]).toHaveValue(defaultProps.destinationDutyLocation.address.streetAddress1);
        expect(address2[1]).toHaveValue('');
        expect(city[1]).toHaveValue(defaultProps.destinationDutyLocation.address.city);
        expect(state[1]).toHaveValue(defaultProps.destinationDutyLocation.address.state);
        expect(postalCodes[1]).toHaveValue(defaultProps.destinationDutyLocation.address.postalCode);
      });
    });
  });

  it('displays secondary destination Address input when hasSecondaryDestinationAddress is true', async () => {
    await act(async () => {
      render(<DateAndLocationForm {...defaultProps} />);
      const hasSecondaryDestinationAddress = await screen.getAllByLabelText('Yes')[1];

      await userEvent.click(hasSecondaryDestinationAddress);
      const postalCodes = screen.getAllByLabelText('ZIP');
      const address1 = screen.getAllByLabelText('Address 1', { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const address3 = screen.getAllByLabelText('Address 3', { exact: false });
      const state = screen.getAllByLabelText('State');
      const city = screen.getAllByLabelText('City');

      await waitFor(() => {
        expect(address1[2]).toBeInstanceOf(HTMLInputElement);
        expect(address2[2]).toBeInstanceOf(HTMLInputElement);
        expect(address3[2]).toBeInstanceOf(HTMLInputElement);
        expect(state[2]).toBeInstanceOf(HTMLSelectElement);
        expect(city[2]).toBeInstanceOf(HTMLInputElement);
        expect(postalCodes[2]).toBeInstanceOf(HTMLInputElement);
      });
    });
  });

  it('displays the closeout office select when the service member is in the Army', async () => {
    await act(async () => {
      const armyServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.ARMY,
      };
      render(<DateAndLocationForm {...defaultProps} serviceMember={armyServiceMember} />);

      expect(screen.getByText('Closeout Office')).toBeInTheDocument();
      expect(screen.getByLabelText('Which closeout office should review your PPM?')).toBeInTheDocument();
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
    });
  });

  it('displays the closeout office select when the service member is in the Air Force', async () => {
    await act(async () => {
      const airForceServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.AIR_FORCE,
      };
      render(<DateAndLocationForm {...defaultProps} serviceMember={airForceServiceMember} />);

      expect(screen.getByText('Closeout Office')).toBeInTheDocument();
      expect(screen.getByLabelText('Which closeout office should review your PPM?')).toBeInTheDocument();
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
    });
  });

  it('5', async () => {
    await act(async () => {
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
});

describe('validates form fields and displays error messages', () => {
  it('marks required inputs when left empty', async () => {
    await act(async () => {
      render(<DateAndLocationForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');

        // Departure date
        expect(requiredAlerts[0]).toHaveTextContent('Required');
        expect(
          within(requiredAlerts[0].nextElementSibling).getByLabelText('When do you plan to start moving your PPM?'),
        ).toBeInTheDocument();
      });
    });
  });
  it('displays type errors when input fails validation schema', async () => {
    await act(async () => {
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

      await userEvent.type(document.querySelector('input[name="pickupAddress.address.postalCode"]'), '1000');

      await userEvent.type(document.querySelector('input[name="destinationAddress.address.postalCode"]'), '1000');

      await userEvent.type(screen.getByLabelText('When do you plan to start moving your PPM?'), '1 January 2022');
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');
        expect(requiredAlerts.length).toBe(3);

        // Departure date
        expect(requiredAlerts[2]).toHaveTextContent('Enter a complete date in DD MMM YYYY format (day, month, year).');
        expect(
          within(requiredAlerts[2].nextElementSibling).getByLabelText('When do you plan to start moving your PPM?'),
        ).toBeInTheDocument();
      });
    });
  });
});
