import React from 'react';
import { render, screen, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import BoatShipmentForm from './BoatShipmentForm';

const mtoShipment = {
  boatShipment: {
    year: '2022',
    make: 'Yamaha',
    model: '242X',
    lengthInInches: 288, // 24 feet
    widthInInches: 102, // 8 feet 6 inches
    heightInInches: 84, // 7 feet
    hasTrailer: true,
    isRoadworthy: true,
  },
};
const emptyBoatInfo = {
  boatShipment: {
    year: '',
    make: '',
    model: '',
    lengthInInches: 0, // 24 feet
    widthInInches: 0, // 8 feet 6 inches
    heightInInches: 0, // 7 feet
    hasTrailer: false,
    isRoadworthy: false,
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
  emptyBoatInfo,
};

beforeEach(() => {
  jest.clearAllMocks();
});

describe('BoatShipmentForm component', () => {
  describe('displays form', () => {
    it('renders filled form on load', async () => {
      render(<BoatShipmentForm {...defaultProps} />);
      expect(await screen.getByTestId('year')).toHaveValue(mtoShipment.boatShipment.year);
      expect(screen.getByTestId('make')).toHaveValue(mtoShipment.boatShipment.make);
      expect(screen.getByTestId('model')).toHaveValue(mtoShipment.boatShipment.model);
      expect(screen.getByTestId('lengthFeet')).toHaveValue('24');
      expect(screen.getByTestId('lengthInches')).toHaveValue('0');
      expect(screen.getByTestId('widthFeet')).toHaveValue('8');
      expect(screen.getByTestId('widthInches')).toHaveValue('6');
      expect(screen.getByTestId('heightFeet')).toHaveValue('7');
      expect(screen.getByTestId('heightInches')).toHaveValue('0');
      expect(screen.getByTestId('hasTrailerYes').checked).toBe(true);
      expect(screen.getByTestId('hasTrailerNo').checked).toBe(false);
      expect(screen.getByTestId('isRoadworthyYes').checked).toBe(true);
      expect(screen.getByTestId('isRoadworthyNo').checked).toBe(false);
      expect(
        screen.getByLabelText(
          'Are there things about this boat shipment that your counselor or movers should know or discuss with you?',
        ),
      ).toBeVisible();
    });
  });

  describe('displays conditional inputs', () => {
    it('displays and hides trailer roadworthy options based on hasTrailer selection', async () => {
      render(<BoatShipmentForm {...defaultProps} />);
      expect(screen.getByText('Is the trailer roadworthy?')).toBeInTheDocument();

      await act(async () => {
        await userEvent.click(screen.getByTestId('hasTrailerNo'));
      });

      expect(screen.queryByText('Is the trailer roadworthy?')).not.toBeInTheDocument();
    });
  });

  describe('validates form fields and displays error messages', () => {
    it('marks required inputs when left empty', async () => {
      render(<BoatShipmentForm {...emptyInfoProps} />);

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
      render(<BoatShipmentForm {...defaultProps} />);

      await act(async () => {
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(defaultProps.onSubmit).toHaveBeenCalled();
    });

    it('does not submit the form with invalid data', async () => {
      render(<BoatShipmentForm {...defaultProps} />);

      await act(async () => {
        await userEvent.clear(screen.getByTestId('year'));
        await userEvent.click(screen.getByRole('button', { name: 'Continue' }));
      });

      expect(defaultProps.onSubmit).not.toHaveBeenCalled();
    });
  });
});
