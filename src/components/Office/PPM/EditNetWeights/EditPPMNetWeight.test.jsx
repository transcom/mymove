import React from 'react';
import { render, screen, act } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditPPMNetWeight from './EditPPMNetWeight';

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
}));

const defaultProps = {
  weightAllowance: 18000,
  moveWeight: 25000,
  originalWeight: 4500,
};

describe('EditNetPPMWeight', () => {
  it('renders weights and edit button initially', () => {
    render(<EditPPMNetWeight />);
    expect(screen.queryByText('Net Weight')).toBeInTheDocument();
    expect(screen.queryByText('original weight')).toBeInTheDocument();
    expect(screen.queryByText('to reduce excess weight')).toBeInTheDocument();
    expect(screen.queryByText('4500lbs')).toBeInTheDocument();
    expect(screen.getAllByRole('button', { name: 'Edit' }));
  });
  describe('when editing PPM Net weight', () => {
    it('renders editting form when the edit button is clicked', async () => {
      render(<EditPPMNetWeight />);
      await act(async () => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      // Calculations for excess weight
      expect(screen.queryByText('Move weight (total)')).toBeInTheDocument();
      expect(screen.queryByText(defaultProps.moveWeight)).toBeInTheDocument();
      expect(screen.queryByText('Weight allowance')).toBeInTheDocument();
      expect(screen.queryByText(defaultProps.weightAllowance)).toBeInTheDocument();
      expect(screen.queryByText('Excess weight (total)')).toBeInTheDocument();
      // Net weight form
      expect(screen.queryByText('Net Weight')).toBeInTheDocument();
      expect(screen.queryByText('original weight')).toBeInTheDocument();
      expect(screen.queryByText(defaultProps.originalWeight)).toBeInTheDocument();
      expect(screen.queryByText('to reduce excess wight')).toBeInTheDocument();
      expect(screen.queryByText('Remarks')).toBeInTheDocument();
      expect(screen.queryByText('Removed the pet elephant to reduce the weight')).toBeInTheDocument();
      // Buttons
      expect(screen.getAllByRole('button', { name: 'Save Changes' }));
      expect(screen.getAllByRole('button', { name: 'Cancel' }));
    });
    // it('disables the save button if save remarks or weight field is empty', () => {
    //   render(<EditPPMNetWeight />);
    // });
    // it('saves changes when editing the form', () => {
    //   render(<EditPPMNetWeight />);
    // });
    // it('shows an error if there is incomplete fields in the form', () => {
    //   render(<EditPPMNetWeight />);
    // });
  });
});
