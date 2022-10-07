import React from 'react';
import { render, waitFor, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import ProGearForm from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm';

const defaultProps = {
  onBack: jest.fn(),
  onSubmit: jest.fn(),
};

const selfProGearProps = {
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

      expect(screen.getByRole('button', { name: 'Return To Homepage' })).toBeInTheDocument();
      expect(screen.getByRole('button', { name: 'Save & Continue' })).toBeEnabled();
    });

    it('does not select a radio when selfProGear is null', () => {
      render(<ProGearForm {...defaultProps} />);
      expect(screen.getByLabelText('Me')).not.toBeChecked();
      expect(screen.getByLabelText('My spouse')).not.toBeChecked();
    });

    it('selects "Me" radio when selfProGear is true', () => {
      render(<ProGearForm {...defaultProps} {...selfProGearProps} />);
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
    it('calls the onSubmit callback with selfProGear set', async () => {
      const expectedPayload = {
        selfProGear: 'true',
      };
      render(<ProGearForm {...defaultProps} {...selfProGearProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Save & Continue' }));

      await waitFor(() => {
        expect(defaultProps.onSubmit).toHaveBeenCalledWith(expectedPayload, expect.anything());
      });
    });
    it('calls the onBack prop when the Return To Homepage button is clicked', async () => {
      render(<ProGearForm {...defaultProps} />);

      await userEvent.click(screen.getByRole('button', { name: 'Return To Homepage' }));

      await waitFor(() => {
        expect(defaultProps.onBack).toHaveBeenCalled();
      });
    });
  });
});
