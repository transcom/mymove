import React from 'react';
import { render, screen, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import MobileHomeShipmentForm from './MobileHomeShipmentForm';

const mtoShipment = {
  mobileHomeShipment: {
    year: '2022',
    make: 'Skyline Homes',
    model: 'Crown',
    lengthInInches: 288, // 24 feet
    widthInInches: 102, // 8 feet 6 inches
    heightInInches: 84, // 7 feet
  },
};
const emptyMobileHomeInfo = {
  mobileHomeShipment: {
    year: '',
    make: '',
    model: '',
    lengthInInches: 0, // 24 feet
    widthInInches: 0, // 8 feet 6 inches
    heightInInches: 0, // 7 feet
  },
};

const defaultProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  mtoShipment,
};

const emptyInfoProps = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  emptyMobileHomeInfo,
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe('MobileHomeShipmentForm component', () => {
  describe('displays form', () => {
    it('renders filled form on load', async () => {
      render(<MobileHomeShipmentForm {...defaultProps} />);
      expect(screen.getByTestId('year')).toHaveValue(mtoShipment.mobileHomeShipment.year);
      expect(screen.getByTestId('make')).toHaveValue(mtoShipment.mobileHomeShipment.make);
      expect(screen.getByTestId('model')).toHaveValue(mtoShipment.mobileHomeShipment.model);
      expect(screen.getByTestId('lengthFeet')).toHaveValue('24');
      expect(screen.getByTestId('lengthInches')).toHaveValue('0');
      expect(screen.getByTestId('widthFeet')).toHaveValue('8');
      expect(screen.getByTestId('widthInches')).toHaveValue('6');
      expect(screen.getByTestId('heightFeet')).toHaveValue('7');
      expect(screen.getByTestId('heightInches')).toHaveValue('0');
      expect(
        screen.getByLabelText(
          'Are there things about this mobile home shipment that your counselor or movers should know or discuss with you?',
        ),
      ).toBeVisible();
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks required inputs when left empty', async () => {
      render(<MobileHomeShipmentForm {...emptyInfoProps} />);

      const requiredFields = [
        'year',
        'make',
        'model',
        'lengthFeet',
        'lengthInches',
        'widthFeet',
        'widthInches',
        'heightFeet',
        'heightInches',
      ];

      await act(async () => {
        requiredFields.forEach(async (field) => {
          const input = screen.getByTestId(field);
          await userEvent.clear(input);
          // await userEvent.click(input);
          fireEvent.blur(input);
        });
      });

      expect(screen.getAllByTestId('errorMessage').length).toBe(requiredFields.length);
    });
  });

  describe('form submission', () => {
    it('submits the form with valid data', async () => {
      render(<MobileHomeShipmentForm {...defaultProps} />);

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(defaultProps.onSubmit).toHaveBeenCalled();
    });

    it('does not submit the form with invalid data', async () => {
      render(<MobileHomeShipmentForm {...defaultProps} />);

      await act(async () => {
        await userEvent.clear(screen.getByTestId('year'));
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(defaultProps.onSubmit).not.toHaveBeenCalled();
    });
  });
});
