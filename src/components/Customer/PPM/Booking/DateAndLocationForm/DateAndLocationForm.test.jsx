import React from 'react';
import { render, waitFor, screen, within, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import DateAndLocationForm from 'components/Customer/PPM/Booking/DateAndLocationForm/DateAndLocationForm';
import SERVICE_MEMBER_AGENCIES from 'content/serviceMemberAgencies';
import { configureStore } from 'shared/store';

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

const mockStore = configureStore({});

describe('DateAndLocationForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );
      expect(screen.getByRole('heading', { level: 2, name: 'Pickup Address' })).toBeInTheDocument();
      const postalCodes = screen.getAllByTestId('ZIP');
      const address1 = screen.getAllByLabelText(/Address 1/);
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const address3 = screen.getAllByLabelText('Address 3', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);

      expect(address1[0]).toBeInstanceOf(HTMLInputElement);
      expect(address2[0]).toBeInstanceOf(HTMLInputElement);
      expect(address3[0]).toBeInstanceOf(HTMLInputElement);
      expect(state[0]).toBeInstanceOf(HTMLLabelElement);
      expect(city[0]).toBeInstanceOf(HTMLLabelElement);
      expect(postalCodes[0]).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Delivery Address' })).toBeInTheDocument();
      expect(address1[1]).toBeInstanceOf(HTMLInputElement);
      expect(address2[1]).toBeInstanceOf(HTMLInputElement);
      expect(address3[1]).toBeInstanceOf(HTMLInputElement);
      expect(state[1]).toBeInstanceOf(HTMLLabelElement);
      expect(city[1]).toBeInstanceOf(HTMLLabelElement);
      expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
      expect(screen.getAllByLabelText('Yes')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[1]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Closeout Office' })).toBeInTheDocument();
      expect(screen.getByLabelText(/Which closeout office should review your PPM\?/)).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Storage' })).toBeInTheDocument();
      expect(screen.getAllByLabelText('Yes')[2]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[2]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      expect(screen.getByLabelText(/When do you plan to start moving your PPM?/)).toBeInstanceOf(HTMLInputElement);
    });
  });

  describe('displays conditional inputs', () => {
    it('displays current address when "Use my current pickup address" is selected', async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );
      const postalCodes = screen.getAllByTestId(/ZIP/);
      expect(postalCodes[0]).toHaveTextContent('');
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });
      await waitFor(() => {
        expect(postalCodes[0]).toHaveTextContent(defaultProps.serviceMember.residential_address.postalCode);
      });
    });

    it('removes current Address when "Use my current pickup address" is deselected', async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });
      const postalCodes = screen.getAllByTestId(/ZIP/);

      await waitFor(() => {
        expect(postalCodes[0]).toHaveTextContent(defaultProps.serviceMember.residential_address.postalCode);
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await waitFor(() => {
        expect(postalCodes[0]).toHaveTextContent('');
      });
    });

    it('displays secondary pickup Address input when hasSecondaryPickupAddress is true', async () => {
      await act(async () => {
        render(
          <Provider store={mockStore.store}>
            <DateAndLocationForm {...defaultProps} />
          </Provider>,
        );
        const hasSecondaryPickupAddress = await screen.getAllByLabelText('Yes')[1];

        await userEvent.click(hasSecondaryPickupAddress);
        const postalCodes = screen.getAllByTestId(/ZIP/);
        const address1 = screen.getAllByLabelText(/Address 1/, { exact: false });
        const address2 = screen.getAllByLabelText('Address 2', { exact: false });
        const state = screen.getAllByTestId(/State/);
        const city = screen.getAllByTestId(/City/);
        await waitFor(() => {
          expect(address1[1]).toBeInstanceOf(HTMLInputElement);
          expect(address2[1]).toBeInstanceOf(HTMLInputElement);
          expect(city[1]).toBeInstanceOf(HTMLLabelElement);
          expect(state[1]).toBeInstanceOf(HTMLLabelElement);
          expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
        });
      });
    });

    it('displays delivery address when "Use my current delivery address" is selected', async () => {
      await act(async () => {
        render(
          <Provider store={mockStore.store}>
            <DateAndLocationForm {...defaultProps} />
          </Provider>,
        );
        await userEvent.click(screen.getByLabelText('Use my current delivery address'));
        const postalCodes = screen.getAllByTestId(/ZIP/);
        const address1 = screen.getAllByLabelText(/Address 1/, { exact: false });
        const address2 = screen.getAllByLabelText('Address 2', { exact: false });
        const state = screen.getAllByTestId(/State/);
        const city = screen.getAllByTestId(/City/);
        expect(address1[1]).toHaveValue(defaultProps.destinationDutyLocation.address.streetAddress1);
        expect(address2[1]).toHaveValue('');
        expect(city[1]).toHaveTextContent(defaultProps.destinationDutyLocation.address.city);
        expect(state[1]).toHaveTextContent(defaultProps.destinationDutyLocation.address.state);
        expect(postalCodes[1]).toHaveTextContent(defaultProps.destinationDutyLocation.address.postalCode);
      });
    });
  });

  it('displays secondary delivery address input when hasSecondaryDestinationAddress is true', async () => {
    await act(async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );
      const hasSecondaryDestinationAddress = await screen.getAllByLabelText('Yes')[1];

      await userEvent.click(hasSecondaryDestinationAddress);
      const postalCodes = screen.getAllByTestId(/ZIP/);
      const address1 = screen.getAllByLabelText(/Address 1/, { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const address3 = screen.getAllByLabelText('Address 3', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);

      await waitFor(() => {
        expect(address1[2]).toBeInstanceOf(HTMLInputElement);
        expect(address2[2]).toBeInstanceOf(HTMLInputElement);
        expect(address3[2]).toBeInstanceOf(HTMLInputElement);
        expect(state[2]).toBeInstanceOf(HTMLLabelElement);
        expect(city[2]).toBeInstanceOf(HTMLLabelElement);
        expect(postalCodes[2]).toBeInstanceOf(HTMLLabelElement);
      });
    });
  });

  it('displays the closeout office select when the service member is in the Army', async () => {
    await act(async () => {
      const armyServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.ARMY,
      };
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} serviceMember={armyServiceMember} />
        </Provider>,
      );

      expect(screen.getByText('Closeout Office')).toBeInTheDocument();
      expect(screen.getByLabelText(/Which closeout office should review your PPM\?/)).toBeInTheDocument();
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
    });
  });

  it('displays the closeout office select when the service member is in the Air Force', async () => {
    await act(async () => {
      const airForceServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.AIR_FORCE,
      };
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} serviceMember={airForceServiceMember} />
        </Provider>,
      );
      expect(screen.getByText('Closeout Office')).toBeInTheDocument();
      expect(screen.getByLabelText(/Which closeout office should review your PPM\?/)).toBeInTheDocument();
      expect(screen.getByText('Start typing a closeout office...')).toBeInTheDocument();
    });
  });

  it('displays the closeout office select when the service member is in the Navy', async () => {
    await act(async () => {
      const navyServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.NAVY,
      };

      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} serviceMember={navyServiceMember} />
        </Provider>,
      );
      expect(screen.queryByText('Closeout Office')).not.toBeInTheDocument();
      expect(screen.queryByLabelText(/Which closeout office should review your PPM\?/)).not.toBeInTheDocument();
      expect(screen.queryByText('Start typing a closeout office...')).not.toBeInTheDocument();
    });
  });
});

describe('validates form fields and displays error messages', () => {
  it('marks required inputs when left empty', async () => {
    render(
      <Provider store={mockStore.store}>
        <DateAndLocationForm {...defaultProps} />
      </Provider>,
    );

    await act(async () => {
      await userEvent.click(screen.getByLabelText(/Which closeout office should review your PPM\?/));
      await userEvent.keyboard('{backspace}[Tab]');
    });

    expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
    await userEvent.click(screen.getByText('Start typing a closeout office...'));
    expect(screen.getByTestId('errorMessage')).toBeVisible();
  });
  it('displays type errors when input fails validation schema', async () => {
    await act(async () => {
      const invalidTypes = {
        ...defaultProps,
        mtoShipment: {
          ppmShipment: {},
        },
      };

      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...invalidTypes} />
        </Provider>,
      );

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '1 January 2022');
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');
        expect(requiredAlerts.length).toBe(1);

        // Departure date
        expect(requiredAlerts[0]).toHaveTextContent('Enter a complete date in DD MMM YYYY format (day, month, year).');
        expect(
          within(requiredAlerts[0].nextElementSibling).getByLabelText(/When do you plan to start moving your PPM?/),
        ).toBeInTheDocument();
      });
    });
  });

  it('delivery address 1 is empty passes validation schema - destination street 1 is OPTIONAL', async () => {
    await act(async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );

      // type something in for delivery address 1
      await userEvent.type(
        document.querySelector('input[name="destinationAddress.address.streetAddress1"]'),
        '1234 Street',
      );
      // now clear out text, should not raise required alert because street is OPTIONAL in DateAndLocationForm context.
      await userEvent.clear(document.querySelector('input[name="destinationAddress.address.streetAddress1"]'));

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        // 'Required' labelHint on address display. expecting a total of 7(2 for pickup address and 3 delivery address with 2 misc).
        // This is to verify Required labelHints are displayed correctly for PPM onboarding/edit for the delivery address
        // street 1 is now OPTIONAL. If this fails it means addtional labelHints have been introduced elsewhere within the control.
        const hints = document.getElementsByClassName('usa-hint');
        expect(hints.length).toBe(11);
        // verify labelHints are actually 'Optional'
        for (let i = 0; i < hints.length; i += 1) {
          expect(hints[i]).toHaveTextContent('Required');
        }
      });
    });
  });

  it('displays tertiary pickup Address input when hasTertiaryPickupAddress is true', async () => {
    await act(async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );

      const hasTertiaryPickupAddress = screen.getAllByLabelText('Yes')[2];

      await userEvent.click(hasTertiaryPickupAddress);
      const postalCodes = screen.getAllByTestId(/ZIP/);
      const address1 = screen.getAllByLabelText(/Address 1/, { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);
      await waitFor(() => {
        expect(address1[1]).toBeInstanceOf(HTMLInputElement);
        expect(address2[1]).toBeInstanceOf(HTMLInputElement);
        expect(city[1]).toBeInstanceOf(HTMLLabelElement);
        expect(state[1]).toBeInstanceOf(HTMLLabelElement);
        expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
      });
    });
  });
  it('displays tertiary delivery address input when hasTertiaryDestinationAddress is true', async () => {
    await act(async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );
      const hasTertiaryDestinationAddress = screen.getAllByLabelText('Yes')[2];

      await userEvent.click(hasTertiaryDestinationAddress);
      const postalCodes = screen.getAllByTestId(/ZIP/);
      const address1 = screen.getAllByLabelText(/Address 1/, { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const address3 = screen.getAllByLabelText('Address 3', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);

      await waitFor(() => {
        expect(address1[1]).toBeInstanceOf(HTMLInputElement);
        expect(address2[1]).toBeInstanceOf(HTMLInputElement);
        expect(address3[1]).toBeInstanceOf(HTMLInputElement);
        expect(state[1]).toBeInstanceOf(HTMLLabelElement);
        expect(city[1]).toBeInstanceOf(HTMLLabelElement);
        expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
      });
    });
  });
});
