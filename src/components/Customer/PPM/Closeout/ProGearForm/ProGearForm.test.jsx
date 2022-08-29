import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';

const defaultProps = {
  onBack: jest.fn(),
  onSubmit: jest.fn(),
  proGear: {
    selfProGear: true,
  },
};

const spouseProGearProps = {
  proGear: {
    selfProGear: false,
  },
};

describe('ProGearForm component', () => {
  describe('displays form', () => {
    it('renders blank form on load with defaults', () => {
      render(<ProGearForm {...defaultProps} />);

      expect(screen.getByRole('heading', { level: 2, name: 'Set 1' })).toBeInTheDocument();
      expect(screen.getByText('Pro-gear belongs to')).toBeInstanceOf(HTMLLegendElement);
      expect(screen.getByLabelText('Me')).toBeInstanceOf(HTMLInputElement);
      expect(screen.getByLabelText('My spouse')).toBeInstanceOf(HTMLInputElement);

      expect(screen.getByRole('button', { name: 'Finish Later' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('selects "Me" radio when selfProGear is true', () => {
      render(<ProGearForm {...defaultProps} />);
      expect(screen.getByLabelText('Me')).toBeChecked();
      expect(screen.getByLabelText('My spouse')).not.toBeChecked();
    });

    it('selects "My spouse" radio when selfProGear is false', () => {
      render(<ProGearForm {...defaultProps} {...spouseProGearProps} />);
      expect(screen.getByLabelText('My spouse')).toBeChecked();
      expect(screen.getByLabelText('Me')).not.toBeChecked();
    });
  });
  describe('attaches button handler callbacks', () => {
    it('calls the onSubmit callback', async () => {
      render(<ProGearForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalled();
      }, expect.anything());
    });
    it('calls the onBack prop when the Finish Later button is clicked', async () => {
      render(<ProGearForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Finish Later' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });
});
