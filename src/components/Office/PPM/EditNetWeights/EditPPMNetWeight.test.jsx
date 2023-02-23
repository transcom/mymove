import React from 'react';
import { render, screen, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { debug } from 'jest-preview';

import EditPPMNetWeight from './EditPPMNetWeight';

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
}));

const defaultProps = {
  maxBillableWeight: 19500,
  originalWeight: 4500,
  weightAllowance: 18000,
  totalBillableWeight: 11000,
  estimatedWeight: 13000,
  billableWeight: 7000,
  ppmNetWeightRemarks: 'Removed the pet elephant to reduce the weight',
};

const excessWeight = {
  maxBillableWeight: 19500,
  originalWeight: 4500,
  weightAllowance: 18000,
  totalBillableWeight: 21000,
  estimatedWeight: 13000,
  billableWeight: 7000,
  ppmNetWeightRemarks: 'Too much weight',
};

const noRemarks = {
  ...defaultProps,
  ppmNetWeightRemarks: '',
};

describe('EditNetPPMWeight', () => {
  it('renders weights and edit button initially', async () => {
    await render(<EditPPMNetWeight {...defaultProps} />);
    expect(screen.getByText('Net weight')).toBeInTheDocument();
    expect(screen.getByText('| original weight')).toBeInTheDocument();
    expect(screen.getByText('| to fit within weight allowance')).toBeInTheDocument();
    expect(screen.getByText('15,500 lbs')).toBeInTheDocument();
    expect(screen.getAllByRole('button', { name: 'Edit' }));
  });
  describe('when editing PPM Net weight', () => {
    it('renders editting form when the edit button is clicked', async () => {
      render(<EditPPMNetWeight {...defaultProps} />);
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      // Net weight form
      expect(screen.getByText('Net weight')).toBeInTheDocument();
      expect(screen.getByText('| original weight')).toBeInTheDocument();
      expect(screen.getByText('4,500 lbs')).toBeInTheDocument();
      expect(screen.getByText('| to fit within weight allowance')).toBeInTheDocument();
      expect(screen.getByText('Remarks')).toBeInTheDocument();
      expect(screen.getByText('Removed the pet elephant to reduce the weight')).toBeInTheDocument();
      // Buttons
      expect(screen.getByRole('button', { name: 'Save changes' }));
      expect(screen.getByRole('button', { name: 'Cancel' }));
    });
    it('renders additional hints for excess weight', async () => {
      render(<EditPPMNetWeight {...excessWeight} />);
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      debug();
      // Calculations for excess weight
      expect(screen.getByText('Move weight (total)')).toBeInTheDocument();
      expect(screen.getByText('21,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Weight allowance')).toBeInTheDocument();
      expect(screen.getByText('18,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Excess weight (total)')).toBeInTheDocument();
      expect(screen.getByText('3,000 lbs')).toBeInTheDocument();
    });
    it('disables the save button if save remarks or weight field is empty', async () => {
      render(<EditPPMNetWeight {...noRemarks} />);
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      expect(await screen.getByRole('button', { name: 'Save changes' })).toBeDisabled();
    });
    it('shows an error if there is incomplete fields in the form', async () => {
      render(<EditPPMNetWeight {...defaultProps} />);
      // Weight Input
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      const textInput = await screen.findByTestId('weightInput');
      await act(() => userEvent.clear(textInput));
      expect(screen.getByText('Required'));
      await act(() => userEvent.type(textInput, '0'));
      expect(screen.getByText('Authorized weight must be greater than or equal to 1'));
      await act(() => userEvent.type(textInput, '1000'));
      // Remarks
      const remarksField = await screen.findByTestId('formRemarks');
      await act(() => userEvent.clear(remarksField));
      await fireEvent.blur(remarksField);
      expect(await screen.findByTestId('errorIndicator')).toBeInTheDocument();
      expect(screen.getByText('Required'));
    });
    // it('saves changes when editing the form', () => {
    //   render(<EditPPMNetWeight />);
    // });
  });
});
