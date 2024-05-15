import React from 'react';
import { render, waitFor, screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import AboutForm from 'components/Customer/PPM/Closeout/AboutForm/AboutForm';

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
      actualPickupPostalCode: '78234',
      actualDestinationPostalCode: '98421',
      pickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#123',
        city: 'San Antonio',
        state: 'TX',
        postalCode: '78234',
      },
      secondaryPickupAddress: {
        streetAddress1: '812 S 129th St',
        streetAddress2: '#124',
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
      secondaryDestinationAddress: {
        streetAddress1: '442 SW Rio de la Plata Drive',
        city: 'Tacoma',
        state: 'WA',
        postalCode: '98421',
      },
      hasSecondaryPickupAddress: 'false',
      hasSecondaryDestinationAddress: 'false',
      hasReceivedAdvance: true,
      advanceAmountReceived: 123456,
      w2Address: {
        streetAddress1: '11 NE Elm Road',
        streetAddress2: '',
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
  await userEvent.paste('31 May 2022');

  within(form).getAllByLabelText('Address 1')[0].focus();
  await userEvent.paste('812 S 129th St');

  within(form)
    .getAllByLabelText(/Address 2/)[0]
    .focus();
  await userEvent.paste('#123');

  within(form).getAllByLabelText('City')[0].focus();
  await userEvent.paste('San Antonio');

  await userEvent.selectOptions(within(form).getAllByLabelText('State')[0], 'TX');

  within(form).getAllByLabelText('ZIP')[0].focus();
  await userEvent.paste('78232');

  within(form).getAllByLabelText('Address 1')[1].focus();
  await userEvent.paste('441 SW Rio de la Plata Drive');

  within(form).getAllByLabelText('City')[1].focus();
  await userEvent.paste('Tacoma');

  await userEvent.selectOptions(within(form).getAllByLabelText('State')[1], 'WA');

  within(form).getAllByLabelText('ZIP')[1].focus();
  await userEvent.paste('98421');

  within(form).getAllByLabelText('Address 1')[2].focus();
  await userEvent.paste('11 NE Elm Road');

  within(form).getAllByLabelText('City')[2].focus();
  await userEvent.paste('Jacksonville');

  await userEvent.selectOptions(within(form).getAllByLabelText('State')[2], 'FL');

  within(form).getAllByLabelText('ZIP')[2].focus();
  await userEvent.paste('32217');
};

describe('AboutForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentProps} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      });
      expect(screen.getByLabelText('When did you leave your origin?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Locations' })).toBeInTheDocument();

      expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      expect(screen.getByTestId('yes-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('no-has-received-advance')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByTestId('no-has-received-advance')).toBeChecked(); // Has advance received is set to No by default

      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue('');
      expect(screen.getAllByLabelText('State')[0]).toHaveValue('');
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue('');

      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('');

      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates appropriate field values', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('When did you leave your origin?')).toHaveDisplayValue('19 May 2022');
      });

      expect(screen.getByTestId('yes-has-received-advance')).toBeChecked();
      expect(screen.getByTestId('no-has-received-advance')).not.toBeChecked();
      expect(screen.getByLabelText('How much did you receive?')).toHaveDisplayValue('1,234');

      expect(screen.getAllByLabelText('Address 1')[2]).toHaveDisplayValue('11 NE Elm Road');
      expect(screen.getAllByLabelText(/Address 2/)[2]).toHaveDisplayValue('');
      expect(screen.getAllByLabelText('City')[2]).toHaveDisplayValue('Jacksonville');
      expect(screen.getAllByLabelText('State')[2]).toHaveDisplayValue('FL');
      expect(screen.getAllByLabelText('ZIP')[2]).toHaveDisplayValue('32217');

      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks all required fields when form is submitted', async () => {
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
      expect(requiredAlerts[1].nextElementSibling).toHaveAttribute('name', 'pickupAddress.streetAddress1');
      expect(requiredAlerts[2]).toHaveTextContent('Required');
      expect(requiredAlerts[2].nextElementSibling).toHaveAttribute('name', 'pickupAddress.city');
      expect(requiredAlerts[3]).toHaveTextContent('Required');
      expect(requiredAlerts[3].nextElementSibling).toHaveAttribute('name', 'pickupAddress.state');
      expect(requiredAlerts[4]).toHaveTextContent('Required');
      expect(requiredAlerts[4].nextElementSibling).toHaveAttribute('name', 'pickupAddress.postalCode');

      expect(requiredAlerts[5]).toHaveTextContent('Required');
      expect(requiredAlerts[5].nextElementSibling).toHaveAttribute('name', 'destinationAddress.streetAddress1');
      expect(requiredAlerts[6]).toHaveTextContent('Required');
      expect(requiredAlerts[6].nextElementSibling).toHaveAttribute('name', 'destinationAddress.city');
      expect(requiredAlerts[7]).toHaveTextContent('Required');
      expect(requiredAlerts[7].nextElementSibling).toHaveAttribute('name', 'destinationAddress.state');
      expect(requiredAlerts[8]).toHaveTextContent('Required');
      expect(requiredAlerts[8].nextElementSibling).toHaveAttribute('name', 'destinationAddress.postalCode');

      await userEvent.click(screen.getByTestId('yes-has-received-advance'));

      await waitFor(() => {
        expect(screen.getByLabelText('How much did you receive?')).toBeInTheDocument();
      });

      requiredAlerts = screen.getAllByRole('alert');

      expect(requiredAlerts[10]).toHaveTextContent('Required');
      expect(
        within(requiredAlerts[19].nextElementSibling).getByLabelText('How much did you receive?'),
      ).toBeInTheDocument();

      expect(requiredAlerts[11]).toHaveTextContent('Required');
      expect(requiredAlerts[11].nextElementSibling).toHaveAttribute('name', 'w2Address.streetAddress1');
      expect(requiredAlerts[12]).toHaveTextContent('Required');
      expect(requiredAlerts[12].nextElementSibling).toHaveAttribute('name', 'w2Address.city');
      expect(requiredAlerts[13]).toHaveTextContent('Required');
      expect(requiredAlerts[13].nextElementSibling).toHaveAttribute('name', 'w2Address.state');
      expect(requiredAlerts[14]).toHaveTextContent('Required');
      expect(requiredAlerts[14].nextElementSibling).toHaveAttribute('name', 'w2Address.postalCode');
    });

    it('displays type error messages for invalid input', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.type(screen.getByLabelText('When did you leave your origin?'), '1 January 2022');
      await userEvent.tab();
    });

    it('displays error when advance received is below 1 dollar minimum', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.click(screen.getByTestId('yes-has-received-advance'));

      await userEvent.type(screen.getByLabelText('How much did you receive?'), '0');

      await waitFor(() => {
        expect(screen.getByRole('alert')).toHaveTextContent(
          "The minimum advance request is $1. If you don't want an advance, select No.",
        );
      });
    });
  });

  describe('calls button event handlers', () => {
    it('calls onBack handler when "Return To Homepage" is pressed', async () => {
      render(<AboutForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });

    it('calls onSubmit handler when "Save & Continue" is pressed', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentProps} />);

      await fillOutBasicForm();

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(
          {
            actualMoveDate: '2022-05-31',
            actualPickupPostalCode: '78234',
            actualDestinationPostalCode: '98421',
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
            hasSecondaryPickupAddress: false,
            hasSecondaryDestinationAddress: false,
            hasReceivedAdvance: 'true',
            advanceAmountReceived: '1234',
            w2Address: {
              streetAddress1: '11 NE Elm Road',
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
