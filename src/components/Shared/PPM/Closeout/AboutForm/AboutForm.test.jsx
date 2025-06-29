import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Provider } from 'react-redux';

import AboutForm from 'components/Shared/PPM/Closeout/AboutForm/AboutForm';
import { configureStore } from 'shared/store';
import { APP_NAME } from 'constants/apps';
import { PPM_TYPES } from 'shared/constants';

jest.mock('utils/featureFlags', () => ({
  ...jest.requireActual('utils/featureFlags'),
  isBooleanFlagEnabled: jest.fn().mockImplementation(() => Promise.resolve(false)),
}));

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  mtoShipment: {
    ppmShipment: {},
  },
};

const shipmentProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  mtoShipment: {
    ppmShipment: {
      actualMoveDate: '31 May 2022',
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        streetAddress3: '',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '78234',
        usPostRegionCitiesID: '',
      },
      destinationAddress: {
        streetAddress1: '441 SW Rio de la Plata Drive',
        streetAddress2: '',
        streetAddress3: '',
        city: 'Tacoma',
        state: 'WA',
        postalCode: '98421',
        usPostRegionCitiesID: '',
      },
      secondaryPickupAddress: {},
      secondaryDestinationAddress: {},
      hasSecondaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
      hasReceivedAdvance: 'true',
      advanceAmountReceived: '100000',
      w2Address: {
        streetAddress1: '11 NE Elm Road',
        streetAddress2: '',
        streetAddress3: '',
        city: 'Jacksonville',
        state: 'FL',
        postalCode: '32217',
        county: 'Duval',
        usPostRegionCitiesID: '',
      },
    },
  },
};

const mockStore = configureStore({});

describe('AboutForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load - Customer Page', async () => {
      render(
        <Provider store={mockStore.store}>
          <AboutForm {...defaultProps} appName={APP_NAME.MYMOVE} />
        </Provider>,
      );

      await waitFor(() => {
        expect(screen.getByText('Finish moving this PPM before you start documenting it.')).toBeInTheDocument();
      });
      const headings = screen.getAllByRole('heading', { level: 2 });
      expect(headings[0]).toHaveTextContent('How to complete your PPM');
      expect(headings[1]).toHaveTextContent('About your final payment');

      // renders form content
      expect(headings[2]).toHaveTextContent('Departure date');
      expect(headings[3]).toHaveTextContent('Locations');
      expect(headings[4]).toHaveTextContent('Advance (AOA)');
      expect(headings[5]).toHaveTextContent('W-2 address');

      expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();

      expect(screen.getByLabelText('When did you leave your origin?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Locations' })).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      expect(screen.getByTestId('yes-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('no-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('no-has-received-advance')).toBeChecked();

      expect(screen.getAllByLabelText(/Address 1/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByTestId('City')[0]).toHaveTextContent('');
      expect(screen.getAllByTestId('State')[0]).toHaveTextContent('');
      expect(screen.getAllByTestId('ZIP')[0]).toHaveTextContent('');

      expect(screen.getAllByLabelText(/Address 1/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByTestId('City')[1]).toHaveTextContent('');
      expect(screen.getAllByTestId('State')[1]).toHaveTextContent('');
      expect(screen.getAllByTestId('ZIP')[1]).toHaveTextContent('');
      expect(screen.getAllByLabelText(/Location Lookup/).length).toBe(3);

      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('renders blank form on load - Office Page', async () => {
      render(
        <Provider store={mockStore.store}>
          <AboutForm {...defaultProps} appName={APP_NAME.OFFICE} />
        </Provider>,
      );

      await waitFor(() => {
        expect(screen.getByText('Finish moving this PPM before you start documenting it.')).toBeInTheDocument();
      });
      const headings = screen.getAllByRole('heading', { level: 2 });
      expect(headings[0]).toHaveTextContent('How to complete your PPM');
      expect(headings[1]).toHaveTextContent('About your final payment');

      // renders form content
      expect(headings[2]).toHaveTextContent('Departure date');
      expect(headings[3]).toHaveTextContent('Locations');
      expect(headings[4]).toHaveTextContent('Advance (AOA)');
      expect(headings[5]).toHaveTextContent('W-2 address');

      expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      expect(screen.getByLabelText('When did you leave your origin?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Locations' })).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      expect(screen.getByTestId('yes-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('no-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('no-has-received-advance')).toBeChecked();

      expect(screen.getAllByLabelText(/Address 1/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');

      expect(screen.getAllByLabelText(/Address 1/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText(/Location Lookup/).length).toBe(3);

      expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    describe('validates form fields and displays error messages', () => {
      it('marks all required fields when form is submitted', async () => {
        render(
          <Provider store={mockStore.store}>
            <AboutForm {...defaultProps} />
          </Provider>,
        );

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
        });

        const requiredAlerts = screen.getAllByRole('alert');

        expect(requiredAlerts[0]).toHaveTextContent('Required');
        expect(
          within(requiredAlerts[0].nextElementSibling).getByLabelText('When did you leave your origin?'),
        ).toBeInTheDocument();

        expect(requiredAlerts[1]).toHaveTextContent('Required');
        expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'w2Address.streetAddress1');
        expect(requiredAlerts[2]).toHaveTextContent('Required');
        expect(requiredAlerts[2].nextElementSibling).toHaveAttribute('aria-label', 'w2Address.city');
        expect(requiredAlerts[3]).toHaveTextContent('Required');
        expect(requiredAlerts[3].nextElementSibling).toHaveAttribute('aria-label', 'w2Address.state');
        expect(requiredAlerts[4]).toHaveTextContent('Required');
        expect(requiredAlerts[4].nextElementSibling).toHaveAttribute('aria-label', 'w2Address.postalCode');

        await userEvent.click(screen.getByTestId('yes-has-received-advance'));
      });
    });

    it('populates appropriate field values', async () => {
      render(
        <Provider store={mockStore.store}>
          <AboutForm {...shipmentProps} />
        </Provider>,
      );

      await waitFor(() => {
        expect(screen.getByLabelText('When did you leave your origin?')).toHaveDisplayValue('31 May 2022');
      });

      expect(screen.getByTestId('yes-has-received-advance')).toBeChecked();
      expect(screen.getByTestId('no-has-received-advance')).not.toBeChecked();
      expect(screen.getByLabelText('How much did you receive?')).toHaveDisplayValue('1,000');

      expect(screen.getAllByLabelText(/Address 1/)[2]).toHaveDisplayValue('11 NE Elm Road');
      expect(screen.getAllByLabelText(/Address 2/)[2]).toHaveDisplayValue('');
      expect(screen.getAllByTestId(/City/)[2]).toHaveTextContent('Jacksonville');
      expect(screen.getAllByTestId(/State/)[2]).toHaveTextContent('FL');
      expect(screen.getAllByTestId(/ZIP/)[2]).toHaveTextContent('32217');
      expect(screen.getAllByTestId(/County/)[2]).toHaveTextContent('Duval');
      expect(screen.getByText('Jacksonville, FL 32217 (Duval)'));
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('PPM destination street1 is required', async () => {
      render(
        <Provider store={mockStore.store}>
          <AboutForm {...shipmentProps} />
        </Provider>,
      );
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

      // Start controlled test case to verify everything is working.
      const input = document.querySelector('input[name="destinationAddress.streetAddress1"]');
      expect(input).toBeInTheDocument();
      // clear
      await userEvent.clear(input);
      await userEvent.tab();
      // verify Required alert is displayed
      const requiredAlerts = screen.getByRole('alert');
      expect(requiredAlerts).toHaveTextContent('Required');

      // verify validation disables save button. destination street 1 is required only in PPM doc upload while
      // it's OPTIONAL during onboarding..etc...
      expect(screen.getByRole('button', { name: 'Save & Continue' })).not.toBeEnabled();

      // verify save is enabled
      await userEvent.type(input, '123 Street');
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

      // 'Optional' labelHint on address display. expecting a total of 9(3 for pickup address, 3 delivery address, 3 w2 address).
      // This is to verify Required labelHints are displayed correctly for PPM doc uploading for the delivery address
      // street 1 is now OPTIONAL for onboarding but required for PPM doc upload. If this fails it means addtional labelHints
      // have been introduced elsewhere within the control.
      const hints = document.getElementsByClassName('usa-hint');
      expect(hints.length).toBe(18);
      // verify labelHints are actually 'Optional'
      for (let i = 0; i < hints.length; i += 1) {
        expect(hints[i]).toHaveTextContent('Required');
      }
    });

    it('displays type error messages for invalid input', async () => {
      render(
        <Provider store={mockStore.store}>
          <AboutForm {...defaultProps} />
        </Provider>,
      );

      await userEvent.type(screen.getByLabelText('When did you leave your origin?'), '1 January 2022');
      await userEvent.tab();
    });

    it('displays error when advance received is below 1 dollar minimum', async () => {
      render(
        <Provider store={mockStore.store}>
          <AboutForm {...defaultProps} />
        </Provider>,
      );

      await userEvent.click(screen.getByTestId('yes-has-received-advance'));

      await userEvent.type(screen.getByLabelText('How much did you receive?'), '0');

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(
          "The minimum advance request is $1. If you don't want an advance, select No.",
        );
      });
    });

    describe('calls button event handlers', () => {
      it('calls onBack handler when "Return To Homepage" is pressed - Customer page', async () => {
        render(
          <Provider store={mockStore.store}>
            <AboutForm {...defaultProps} appName={APP_NAME.MYMOVE} />
          </Provider>,
        );

        await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));

        await waitFor(() => {
          expect(defaultProps.onBack).toHaveBeenCalled();
        });
      });

      it('calls onBack handler when "Cancel" is pressed - Office page', async () => {
        render(
          <Provider store={mockStore.store}>
            <AboutForm {...defaultProps} appName={APP_NAME.OFFICE} />
          </Provider>,
        );

        await userEvent.click(screen.getByRole('button', { name: 'Cancel' }));

        await waitFor(() => {
          expect(defaultProps.onBack).toHaveBeenCalled();
        });
      });

      it('calls onSubmit handler when "Save & Continue" is pressed', async () => {
        render(
          <Provider store={mockStore.store}>
            <AboutForm {...shipmentProps} />
          </Provider>,
        );

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          expect(shipmentProps.onSubmit).toHaveBeenCalledWith(
            {
              actualMoveDate: '31 May 2022',
              pickupAddress: {
                streetAddress1: '812 S 129th St',
                streetAddress2: '#123',
                streetAddress3: '',
                city: 'San Antonio',
                state: 'TX',
                postalCode: '78234',
                usPostRegionCitiesID: '',
              },
              destinationAddress: {
                streetAddress1: '441 SW Rio de la Plata Drive',
                streetAddress2: '',
                streetAddress3: '',
                city: 'Tacoma',
                state: 'WA',
                postalCode: '98421',
                usPostRegionCitiesID: '',
              },
              secondaryPickupAddress: {},
              secondaryDestinationAddress: {},
              hasSecondaryPickupAddress: 'false',
              hasSecondaryDestinationAddress: 'false',
              hasReceivedAdvance: 'true',
              advanceAmountReceived: '1000',
              w2Address: {
                streetAddress1: '11 NE Elm Road',
                streetAddress2: '',
                streetAddress3: '',
                city: 'Jacksonville',
                state: 'FL',
                postalCode: '32217',
                county: 'Duval',
                usPostRegionCitiesID: '',
              },
            },
            expect.anything(),
          );
        });
      });
    });

    describe('AboutForm - when ppmType is SMALL_PACKAGE', () => {
      const smallPackageMtoShipment = {
        ppmShipment: {
          ppmType: PPM_TYPES.SMALL_PACKAGE,
          actualMoveDate: '01 Jan 2022',
          pickupAddress: {
            streetAddress1: '123 Small Package St',
            streetAddress2: '',
            streetAddress3: '',
            city: 'Smalltown',
            state: 'SP',
            postalCode: '12345',
            usPostRegionCitiesID: '',
          },
          destinationAddress: {
            streetAddress1: '456 Destination Ave',
            streetAddress2: '',
            streetAddress3: '',
            city: 'Destination City',
            state: 'SP',
            postalCode: '67890',
            usPostRegionCitiesID: '',
          },
          hasReceivedAdvance: false,
          w2Address: {
            streetAddress1: '',
            streetAddress2: '',
            streetAddress3: '',
            city: '',
            state: '',
            postalCode: '',
            usPostRegionCitiesID: '',
          },
        },
      };

      const smallPackageProps = {
        onSubmit: jest.fn(),
        onBack: jest.fn(),
        mtoShipment: smallPackageMtoShipment,
      };

      it('renders "Shipped Date" heading and small package labels', async () => {
        render(
          <Provider store={mockStore.store}>
            <AboutForm {...smallPackageProps} appName={APP_NAME.MYMOVE} />
          </Provider>,
        );

        const headings = screen.getAllByRole('heading', { level: 2 });
        expect(headings[2]).toHaveTextContent('Shipped Date');

        expect(screen.getByLabelText('When did you ship your package?')).toBeInTheDocument();

        expect(
          screen.queryByText(/If you picked things up or dropped things off from other places/),
        ).not.toBeInTheDocument();

        expect(screen.queryByText('Destination Address')).toBeInTheDocument();
        expect(screen.queryByText('Delivery Address')).not.toBeInTheDocument();

        expect(screen.getByText('W-2 address')).toBeInTheDocument();

        expect(screen.getByRole('heading', { level: 2, name: /Locations/ })).toBeInTheDocument();

        expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      });
    });
  });
});
