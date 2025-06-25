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
      county: 'Muscogee',
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
      streetAddress1: '658 West Ave',
      streetAddress2: '',
      county: 'Muscogee',
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
    it('renders blank form on load and asterisks for required fields', async () => {
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...defaultProps} />
        </Provider>,
      );
      expect(document.querySelector('#reqAsteriskMsg')).toHaveTextContent('Fields marked with * are required.');
      expect(screen.getByRole('heading', { level: 2, name: 'Pickup Address' })).toBeInTheDocument();
      const postalCodes = screen.getAllByTestId('ZIP');
      const locationLookups = screen.getAllByLabelText(/Location Lookup/);
      const address1 = screen.getAllByLabelText(/Address 1 */);
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
      expect(locationLookups[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('Yes')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getAllByLabelText('No')[0]).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Delivery Address' })).toBeInTheDocument();
      expect(address1[1]).toBeInstanceOf(HTMLInputElement);
      expect(address2[1]).toBeInstanceOf(HTMLInputElement);
      expect(address3[1]).toBeInstanceOf(HTMLInputElement);
      expect(state[1]).toBeInstanceOf(HTMLLabelElement);
      expect(city[1]).toBeInstanceOf(HTMLLabelElement);
      expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
      expect(locationLookups[1]).toBeInstanceOf(HTMLInputElement);
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

    it('disables the save button if the move has been locked by an office user', async () => {
      await act(async () => {
        const defaultPropsWithLock = {
          ...defaultProps,
          mtoShipment: {
            ppmShipment: {},
          },
          isMoveLocked: true,
        };

        render(
          <Provider store={mockStore.store}>
            <DateAndLocationForm {...defaultPropsWithLock} />
          </Provider>,
        );

        await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM?/), '1 January 2022');
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
      });
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
        expect(
          screen.getAllByText(
            `${defaultProps.serviceMember.residential_address.city}, ${defaultProps.serviceMember.residential_address.state} ${defaultProps.serviceMember.residential_address.postalCode} (${defaultProps.serviceMember.residential_address.county})`,
          ),
        );
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
        expect(screen.getAllByText('Start typing a Zip or City, State Zip').length).toBe(1);
      });

      await waitFor(() => {
        expect(
          screen.getAllByText(
            `${defaultProps.serviceMember.residential_address.city}, ${defaultProps.serviceMember.residential_address.state} ${defaultProps.serviceMember.residential_address.postalCode} (${defaultProps.serviceMember.residential_address.county})`,
          ),
        );
      });

      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await waitFor(() => {
        expect(postalCodes[0]).toHaveTextContent('');
        expect(screen.getAllByText('Start typing a Zip or City, State Zip').length).toBe(2);
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
        const locationLookups = screen.getAllByLabelText(/Location Lookup/);
        const address1 = screen.getAllByLabelText(/Address 1 */, { exact: false });
        const address2 = screen.getAllByLabelText('Address 2', { exact: false });
        const state = screen.getAllByTestId(/State/);
        const city = screen.getAllByTestId(/City/);
        await waitFor(() => {
          expect(address1[1]).toBeInstanceOf(HTMLInputElement);
          expect(address2[1]).toBeInstanceOf(HTMLInputElement);
          expect(city[1]).toBeInstanceOf(HTMLLabelElement);
          expect(state[1]).toBeInstanceOf(HTMLLabelElement);
          expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
          expect(locationLookups[1]).toBeInstanceOf(HTMLInputElement);
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
        expect(
          screen.getAllByText(
            `${defaultProps.destinationDutyLocation.address.city}, ${defaultProps.destinationDutyLocation.address.state} ${defaultProps.destinationDutyLocation.address.postalCode} (${defaultProps.destinationDutyLocation.address.county})`,
          ),
        );
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

      await userEvent.click(screen.getByLabelText('Use my current delivery address'));

      const postalCodes = screen.getAllByTestId(/ZIP/);
      const address1 = screen.getAllByLabelText(/Address 1 */, { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);
      expect(address1[1]).toHaveValue(defaultProps.destinationDutyLocation.address.streetAddress1);
      expect(address2[1]).toHaveValue('');
      expect(city[1]).toHaveTextContent(defaultProps.destinationDutyLocation.address.city);
      expect(state[1]).toHaveTextContent(defaultProps.destinationDutyLocation.address.state);
      expect(postalCodes[1]).toHaveTextContent(defaultProps.destinationDutyLocation.address.postalCode);
      expect(
        screen.getAllByText(
          `${defaultProps.destinationDutyLocation.address.city}, ${defaultProps.destinationDutyLocation.address.state} ${defaultProps.destinationDutyLocation.address.postalCode} (${defaultProps.destinationDutyLocation.address.county})`,
        ),
      );

      const hasSecondaryDestinationAddress = await screen.getAllByLabelText('Yes')[1];
      await userEvent.click(hasSecondaryDestinationAddress);
      const secondaryPostalCodes = screen.getAllByTestId(/ZIP/);
      const locationLookups = screen.getAllByLabelText(/Location Lookup/);
      const secondaryAddress1 = screen.getAllByLabelText(/Address 1 */, { exact: false });
      const secondaryAddress2 = screen.getAllByLabelText('Address 2', { exact: false });
      const secondaryAddress3 = screen.getAllByLabelText('Address 3', { exact: false });
      const secondaryState = screen.getAllByTestId(/State/);
      const secondaryCity = screen.getAllByTestId(/City/);

      await waitFor(() => {
        expect(secondaryAddress1[2]).toBeInstanceOf(HTMLInputElement);
        expect(secondaryAddress2[2]).toBeInstanceOf(HTMLInputElement);
        expect(secondaryAddress3[2]).toBeInstanceOf(HTMLInputElement);
        expect(secondaryState[2]).toBeInstanceOf(HTMLLabelElement);
        expect(secondaryCity[2]).toBeInstanceOf(HTMLLabelElement);
        expect(secondaryPostalCodes[2]).toBeInstanceOf(HTMLLabelElement);
        expect(locationLookups[2]).toBeInstanceOf(HTMLInputElement);
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
      expect(screen.getByLabelText(/Which closeout office should review your PPM\? */)).toBeInTheDocument();
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
      expect(screen.getByLabelText(/Which closeout office should review your PPM\? */)).toBeInTheDocument();
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
      expect(screen.queryByLabelText(/Which closeout office should review your PPM\? */)).not.toBeInTheDocument();
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
      await userEvent.click(screen.getByLabelText(/Which closeout office should review your PPM\? */));
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

      await userEvent.type(screen.getByLabelText(/When do you plan to start moving your PPM? */), '1 January 2022');
      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();

        const requiredAlerts = screen.getAllByRole('alert');
        expect(requiredAlerts.length).toBe(1);

        // Departure date
        expect(requiredAlerts[0]).toHaveTextContent('Enter a complete date in DD MMM YYYY format (day, month, year).');
        expect(
          within(requiredAlerts[0].nextElementSibling).getByLabelText(/When do you plan to start moving your PPM? */),
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

        // E-05732: for PPMs, the destination address street 1 is now optional except for closeout
        // this field is usually always required other than PPMs
        const labelsWithAsterisk = screen.queryAllByText((content) => content.trim().endsWith('*'));
        expect(labelsWithAsterisk.length).toBe(9);
        // 'Required' labelHint on address display. expecting a total of 7(2 for pickup address and 3 delivery address with 2 misc).
        // This is to verify Required labelHints are displayed correctly for PPM onboarding/edit for the delivery address
        // street 1 is now OPTIONAL. If this fails it means addtional labelHints have been introduced elsewhere within the control.
        const hints = document.getElementsByClassName('usa-hint');
        expect(hints.length).toBe(15);
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
      const locationLookups = screen.getAllByLabelText(/Location Lookup/);
      const address1 = screen.getAllByLabelText(/Address 1 */, { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);

      await waitFor(() => {
        expect(address1[1]).toBeInstanceOf(HTMLInputElement);
        expect(address2[1]).toBeInstanceOf(HTMLInputElement);
        expect(city[1]).toBeInstanceOf(HTMLLabelElement);
        expect(state[1]).toBeInstanceOf(HTMLLabelElement);
        expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
        expect(locationLookups[1]).toBeInstanceOf(HTMLInputElement);
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
      const address1 = screen.getAllByLabelText(/Address 1 */, { exact: false });
      const address2 = screen.getAllByLabelText('Address 2', { exact: false });
      const address3 = screen.getAllByLabelText('Address 3', { exact: false });
      const state = screen.getAllByTestId(/State/);
      const city = screen.getAllByTestId(/City/);
      const locationLookup = screen.getAllByLabelText('Location Lookup', { exact: false });

      await waitFor(() => {
        expect(address1[1]).toBeInstanceOf(HTMLInputElement);
        expect(address2[1]).toBeInstanceOf(HTMLInputElement);
        expect(address3[1]).toBeInstanceOf(HTMLInputElement);
        expect(state[1]).toBeInstanceOf(HTMLLabelElement);
        expect(city[1]).toBeInstanceOf(HTMLLabelElement);
        expect(postalCodes[1]).toBeInstanceOf(HTMLLabelElement);
        expect(locationLookup[1]).toBeInstanceOf(HTMLInputElement);
      });
    });
  });

  it('remove Required alert when secondary pickup/delivery streetAddress1 is cleared but the toggle is switched to No', async () => {
    await act(async () => {
      const newPPM = {
        ...defaultProps,
        mtoShipment: {
          ppmShipment: {
            secondaryPickupAddress: {
              streetAddress1: '777 Test Street',
              city: 'ELIZABETHTOWN',
              state: 'KY',
              postalCode: '42702',
              county: 'Hardin',
            },
            secondaryDestinationAddress: {
              streetAddress1: '68 West Elm',
              city: 'Fort Benning',
              state: 'GA',
              postalCode: '94611',
              county: 'Muscogee',
            },
            expectedDepartureDate: '2025-03-08',
          },
        },
      };
      const navyServiceMember = {
        ...defaultProps.serviceMember,
        affiliation: SERVICE_MEMBER_AGENCIES.NAVY,
      };
      render(
        <Provider store={mockStore.store}>
          <DateAndLocationForm {...newPPM} serviceMember={navyServiceMember} />
        </Provider>,
      );
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current pickup address'));
      });

      await userEvent.click(screen.getByTitle('Yes, I have a second pickup address'));
      await act(async () => {
        await userEvent.click(screen.getByLabelText('Use my current delivery address'));
      });
      await userEvent.click(screen.getByTitle('Yes, I have a second delivery address'));

      const address1 = screen.getAllByLabelText(/Address 1 */, { exact: false });
      const locationLookups = screen.getAllByLabelText(/Location Lookup */);

      // verify pickup address is populated
      expect(address1[0]).toHaveValue('123 Main');
      expect(screen.getByText('Fort Benning, GA 90210 (Muscogee)'));

      await waitFor(() => {
        expect(address1[1]).toBeInstanceOf(HTMLInputElement);
        expect(locationLookups[1]).toBeInstanceOf(HTMLInputElement);
      });

      // verify 2nd pickup is populated
      expect(screen.getByRole('heading', { level: 4, name: 'Second Pickup Address' })).toBeInTheDocument();
      expect(address1[1]).toHaveValue('777 Test Street');
      expect(screen.getByText('ELIZABETHTOWN, KY 42702 (Hardin)'));

      // verify delivery address is populated
      expect(address1[2]).toHaveValue('658 West Ave');
      expect(screen.getAllByText('Fort Benning, GA 94611 (Muscogee)')[0]);

      await waitFor(() => {
        expect(address1[3]).toBeInstanceOf(HTMLInputElement);
        expect(locationLookups[3]).toBeInstanceOf(HTMLInputElement);
      });

      // verify 2nd delivery address is populated
      expect(screen.getByRole('heading', { level: 4, name: 'Second Delivery Address' })).toBeInTheDocument();

      expect(address1[3]).toHaveValue('68 West Elm');
      expect(screen.getAllByText('Fort Benning, GA 94611 (Muscogee)')[1]);

      // now clear out 2nd pickup address1 text, should raise required alert
      await userEvent.clear(document.querySelector('input[name="secondaryPickupAddress.address.streetAddress1"]'));
      await userEvent.keyboard('[Tab]');

      await waitFor(() => {
        const requiredAlerts = screen.queryAllByRole('alert');
        expect(requiredAlerts.length).toBe(1);
        requiredAlerts.forEach((alert) => {
          expect(alert).toHaveTextContent('Required');
        });
      });

      // toggle second pickup address to No, should get rid of Required error
      await userEvent.click(screen.getByTitle('No, I do not have a second pickup address'));

      const alerts = screen.queryAllByRole('alert');
      expect(alerts.length).toBe(0);

      // now clear out 2nd delivery address1 text, should raise required alert
      await userEvent.clear(document.querySelector('input[name="secondaryDestinationAddress.address.streetAddress1"]'));
      await userEvent.keyboard('[Tab]');

      await waitFor(() => {
        const requiredAlerts = screen.queryAllByRole('alert');
        expect(requiredAlerts.length).toBe(1);
        requiredAlerts.forEach((alert) => {
          expect(alert).toHaveTextContent('Required');
        });
      });

      // toggle second delivery address to No, should get rid of Required error
      await userEvent.click(screen.getByTitle('No, I do not have a second delivery address'));

      const newAlerts = screen.queryAllByRole('alert');
      expect(newAlerts.length).toBe(0);
    });
  });
});
