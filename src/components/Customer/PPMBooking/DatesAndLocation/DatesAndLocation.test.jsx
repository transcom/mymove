import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
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
});
