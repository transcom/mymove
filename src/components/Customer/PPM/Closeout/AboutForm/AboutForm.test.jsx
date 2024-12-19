import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AboutForm from 'components/Customer/PPM/Closeout/AboutForm/AboutForm';

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
      actualPickupPostalCode: '',
      actualDestinationPostalCode: '',
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        streetAddress3: '',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '78234',
      },
      destinationAddress: {
        streetAddress1: '441 SW Rio de la Plata Drive',
        streetAddress2: '',
        streetAddress3: '',
        city: 'Tacoma',
        state: 'WA',
        postalCode: '98421',
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
      },
    },
  },
};

const fillOutBasicForm = async () => {
  let form;
  await waitFor(() => {
    form = screen.getByTestId('aboutForm');
  });

  within(form).getByLabelText('When did you leave your origin?').focus();
  await userEvent.paste('2022-05-31');

  within(form)
    .getAllByLabelText(/Address 1/)[0]
    .focus();
  await userEvent.paste('812 S 129th St');

  within(form)
    .getAllByLabelText(/Address 2/)[0]
    .focus();
  await userEvent.paste('#123');

  within(form)
    .getAllByLabelText(/Address 3/)[0]
    .focus();
  await userEvent.paste('');

  within(form).getAllByLabelText(/City/)[0].focus();
  await userEvent.paste('San Antonio');

  await userEvent.selectOptions(within(form).getAllByLabelText(/State/)[0], 'TX');

  within(form).getAllByLabelText(/ZIP/)[0].focus();
  await userEvent.paste('78234');

  within(form)
    .getAllByLabelText(/Address 1/)[1]
    .focus();
  await userEvent.paste('441 SW Rio de la Plata Drive');

  within(form)
    .getAllByLabelText(/Address 2/)[1]
    .focus();
  await userEvent.paste('');

  within(form)
    .getAllByLabelText(/Address 3/)[1]
    .focus();
  await userEvent.paste('');

  within(form).getAllByLabelText(/City/)[1].focus();
  await userEvent.paste('Tacoma');

  await userEvent.selectOptions(within(form).getAllByLabelText(/State/)[1], 'WA');

  within(form).getAllByLabelText(/ZIP/)[1].focus();
  await userEvent.paste('98421');

  within(form)
    .getAllByLabelText(/Address 1/)[2]
    .focus();
  await userEvent.paste('11 NE Elm Road');

  within(form).getAllByLabelText(/City/)[2].focus();
  await userEvent.paste('Jacksonville');

  await userEvent.selectOptions(within(form).getAllByLabelText(/State/)[2], 'FL');

  within(form).getAllByLabelText(/ZIP/)[2].focus();
  await userEvent.paste('32217');

  await userEvent.click(within(form).getAllByLabelText('Yes')[2]);

  within(form).getByLabelText('How much did you receive?').focus();
  await userEvent.paste('1000');
};

describe('AboutForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<AboutForm {...defaultProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      });
      await expect(screen.getByLabelText('When did you leave your origin?')).toBeInstanceOf(HTMLInputElement);
      await expect(screen.getByRole('heading', { level: 2, name: 'Locations' })).toBeInTheDocument();

      await expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      await expect(screen.getByTestId('yes-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      await expect(screen.getByTestId('no-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      await expect(screen.getByTestId('no-has-received-advance')).toBeChecked();

      await expect(screen.getAllByLabelText(/Address 1/)[0]).toHaveValue('');
      await expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      await expect(screen.getAllByLabelText(/City/)[0]).toHaveValue('');
      await expect(screen.getAllByLabelText(/State/)[0]).toHaveValue('');
      await expect(screen.getAllByLabelText(/ZIP/)[0]).toHaveValue('');

      await expect(screen.getAllByLabelText(/Address 1/)[1]).toHaveValue('');
      await expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      await expect(screen.getAllByLabelText(/City/)[1]).toHaveValue('');
      await expect(screen.getAllByLabelText(/State/)[1]).toHaveValue('');
      await expect(screen.getAllByLabelText(/ZIP/)[1]).toHaveValue('');

      await expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
      await expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    describe('validates form fields and displays error messages', () => {
      it('marks all required fields when form is submitted', async () => {
        render(<AboutForm {...defaultProps} />);

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeDisabled();
        });

        const requiredAlerts = screen.getAllByRole('alert');

        await expect(requiredAlerts[0]).toHaveTextContent('Required');
        await expect(
          within(requiredAlerts[0].nextElementSibling).getByLabelText('When did you leave your origin?'),
        ).toBeInTheDocument();

        await expect(requiredAlerts[1]).toHaveTextContent('Required');
        await expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'w2Address.streetAddress1');
        await expect(requiredAlerts[2]).toHaveTextContent('Required');
        await expect(requiredAlerts[2].nextElementSibling).toHaveAttribute('name', 'w2Address.city');
        await expect(requiredAlerts[3]).toHaveTextContent('Required');
        await expect(requiredAlerts[3].nextElementSibling).toHaveAttribute('name', 'w2Address.state');
        await expect(requiredAlerts[4]).toHaveTextContent('Required');
        await expect(requiredAlerts[4].nextElementSibling).toHaveAttribute('name', 'w2Address.postalCode');

        await userEvent.click(screen.getByTestId('yes-has-received-advance'));
      });
    });

    it('populates appropriate field values', async () => {
      render(<AboutForm {...shipmentProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('When did you leave your origin?')).toHaveDisplayValue('31 May 2022');
      });

      await expect(screen.getByTestId('yes-has-received-advance')).toBeChecked();
      await expect(screen.getByTestId('no-has-received-advance')).not.toBeChecked();
      await expect(screen.getByLabelText('How much did you receive?')).toHaveDisplayValue('1,000');

      await expect(screen.getAllByLabelText(/Address 1/)[2]).toHaveDisplayValue('11 NE Elm Road');
      await expect(screen.getAllByLabelText(/Address 2/)[2]).toHaveDisplayValue('');
      await expect(screen.getAllByLabelText(/City/)[2]).toHaveDisplayValue('Jacksonville');
      await expect(screen.getAllByLabelText(/State/)[2]).toHaveDisplayValue('FL');
      await expect(screen.getAllByLabelText(/ZIP/)[2]).toHaveDisplayValue('32217');

      await expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('PPM destination street1 is required', async () => {
      render(<AboutForm {...shipmentProps} />);
      await expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

      // Start controlled test case to verify everything is working.
      const input = await document.querySelector('input[name="destinationAddress.streetAddress1"]');
      expect(input).toBeInTheDocument();
      // clear
      await userEvent.clear(input);
      await userEvent.tab();
      // verify Required alert is displayed
      const requiredAlerts = screen.getByRole('alert');
      expect(requiredAlerts).toHaveTextContent('Required');

      // verify validation disables save button. destination street 1 is required only in PPM doc upload while
      // it's OPTIONAL during onboarding..etc...
      await expect(screen.getByRole('button', { name: 'Save & Continue' })).not.toBeEnabled();

      // verify save is enabled
      await userEvent.type(input, '123 Street');
      await expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();

      // 'Optional' labelHint on address display. expecting a total of 9(3 for pickup address, 3 delivery address, 3 w2 address).
      // This is to verify Required labelHints are displayed correctly for PPM doc uploading for the delivery address
      // street 1 is now OPTIONAL for onboarding but required for PPM doc upload. If this fails it means addtional labelHints
      // have been introduced elsewhere within the control.
      const hints = document.getElementsByClassName('usa-hint');
      expect(hints.length).toBe(9);
      // verify labelHints are actually 'Optional'
      for (let i = 0; i < hints.length; i += 1) {
        expect(hints[i]).toHaveTextContent('Required');
      }
    });

    it('displays type error messages for invalid input', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.type(screen.getByLabelText('When did you leave your origin?'), '1 January 2022');
      await userEvent.tab();
    });

    it('displays error when advance received is below 1 dollar minimum', async () => {
      await render(<AboutForm {...defaultProps} />);

      await userEvent.click(screen.getByTestId('yes-has-received-advance'));

      await userEvent.type(screen.getByLabelText('How much did you receive?'), '0');

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(
          "The minimum advance request is $1. If you don't want an advance, select No.",
        );
      });
    });

    describe('calls button event handlers', () => {
      it('calls onBack handler when "Return To Homepage" is pressed', async () => {
        await render(<AboutForm {...defaultProps} />);

        await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));

        await waitFor(() => {
          expect(defaultProps.onBack).toHaveBeenCalled();
        });
      });

      it('calls onSubmit handler when "Save & Continue" is pressed', async () => {
        await render(<AboutForm {...defaultProps} />);

        await fillOutBasicForm();

        await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

        await waitFor(() => {
          expect(defaultProps.onSubmit).toHaveBeenCalledWith(
            {
              actualMoveDate: '31 May 2022',
              actualPickupPostalCode: '',
              actualDestinationPostalCode: '',
              pickupAddress: {
                streetAddress1: '812 S 129th St',
                streetAddress2: '#123',
                city: 'San Antonio',
                state: 'TX',
                postalCode: '78234',
              },
              destinationAddress: {
                streetAddress1: '441 SW Rio de la Plata Drive',
                city: 'Tacoma',
                state: 'WA',
                postalCode: '98421',
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
              },
            },
            expect.anything(),
          );
        });
      });
    });
  });
});
