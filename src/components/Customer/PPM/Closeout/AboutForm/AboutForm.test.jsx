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
      hasSecondaryPickupAddress: 'true',
      hasSecondaryDestinationAddress: 'true',
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

const mtoShipmentWithZips = {
  mtoShipment: {
    ppmShipment: {
      pickupPostalCode: '78234',
      destinationPostalCode: '98421',
    },
  },
};

describe('AboutForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load (except zips)', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentWithZips} />);

      await waitFor(() => {
        expect(screen.getByRole('heading', { level: 2, name: 'Departure date' })).toBeInTheDocument();
      });
      expect(screen.getByLabelText('When did you leave your origin?')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByRole('heading', { level: 2, name: 'Locations' })).toBeInTheDocument();

      const startingZip = screen.getByLabelText('Starting ZIP');
      expect(startingZip).toBeInstanceOf(HTMLInputElement);
      expect(startingZip).toHaveDisplayValue('42442');

      const endingZip = screen.getByLabelText('Ending ZIP');
      expect(endingZip).toBeInstanceOf(HTMLInputElement);
      expect(endingZip).toHaveDisplayValue('42444');

      expect(screen.getByRole('heading', { level: 2, name: 'Advance (AOA)' })).toBeInTheDocument();
      expect(screen.getByLabelText('Yes')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('No')).toBeChecked(); // Has advance received is set to No by default

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

      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[2]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[2]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[2]).toHaveValue('');
      expect(screen.getAllByLabelText('State')[2]).toHaveValue('');
      expect(screen.getAllByLabelText('ZIP')[2]).toHaveValue('');

      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[3]).toHaveValue('');
      expect(screen.getAllByLabelText(/Address 2/)[3]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[3]).toHaveValue('');
      expect(screen.getAllByLabelText('State')[3]).toHaveValue('');
      expect(screen.getAllByLabelText('ZIP')[3]).toHaveValue('');

      expect(screen.getAllByLabelText('Address 1')[4]).toHaveDisplayValue('');
      expect(screen.getAllByLabelText(/Address 2/)[4]).toHaveDisplayValue('');
      expect(screen.getAllByLabelText('City')[4]).toHaveDisplayValue('');
      expect(screen.getAllByLabelText('State')[4]).toHaveValue('');
      expect(screen.getAllByLabelText('ZIP')[4]).toHaveDisplayValue('');

      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('populates all field values', async () => {
      render(<AboutForm {...defaultProps} {...mtoShipmentProps} />);

      await waitFor(() => {
        expect(screen.getByLabelText('When did you leave your origin?')).toHaveDisplayValue('19 May 2022');
      });

      expect(screen.getByLabelText('Yes')).toBeChecked();
      expect(screen.getByLabelText('No')).not.toBeChecked();
      expect(screen.getByLabelText('How much did you receive?')).toHaveDisplayValue('1,234');

      expect(screen.getAllByLabelText('Address 1')[0]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[0]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[0]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText('State')[0]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[0]).toHaveValue('78234');

      expect(screen.getAllByLabelText('Address 1')[1]).toHaveValue('812 S 129th St');
      expect(screen.getAllByLabelText(/Address 2/)[1]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[1]).toHaveValue('San Antonio');
      expect(screen.getAllByLabelText('State')[1]).toHaveValue('TX');
      expect(screen.getAllByLabelText('ZIP')[1]).toHaveValue('78234');

      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[2]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[2]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[2]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[2]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[2]).toHaveValue('98421');

      expect(screen.getByTitle('Yes, I know my delivery address')).toBeChecked();
      expect(screen.getAllByLabelText('Address 1')[3]).toHaveValue('441 SW Rio de la Plata Drive');
      expect(screen.getAllByLabelText(/Address 2/)[3]).toHaveValue('');
      expect(screen.getAllByLabelText('City')[3]).toHaveValue('Tacoma');
      expect(screen.getAllByLabelText('State')[3]).toHaveValue('WA');
      expect(screen.getAllByLabelText('ZIP')[3]).toHaveValue('98421');

      expect(screen.getAllByLabelText('Address 1')[4]).toHaveDisplayValue('11 NE Elm Road');
      expect(screen.getAllByLabelText(/Address 2/)[4]).toHaveDisplayValue('');
      expect(screen.getAllByLabelText('City')[4]).toHaveDisplayValue('Jacksonville');
      expect(screen.getAllByLabelText('State')[4]).toHaveDisplayValue('FL');
      expect(screen.getAllByLabelText('ZIP')[4]).toHaveDisplayValue('32217');

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

      expect(requiredAlerts[9]).toHaveTextContent('Required');
      expect(requiredAlerts[9].nextElementSibling).toHaveAttribute('name', 'actualDestinationPostalCode');

      await userEvent.click(screen.getByLabelText('Yes'));

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

      await userEvent.click(screen.getByLabelText('Yes'));

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

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(
          {
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
            hasSecondaryPickupAddress: 'true',
            hasSecondaryDestinationAddress: 'true',
            hasReceivedAdvance: 'true',
            advanceAmountReceived: '1234',
            w2Address: {
              streetAddress1: '11 NE Elm Road',
              streetAddress2: '',
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
