import React from 'react';
import { render, screen, act, fireEvent } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditPPMNetWeight from './EditPPMNetWeight';

import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
}));

const shipments = [
  // moveWeightTotal = 8000
  { primeActualWeight: 1000, reweigh: null, status: 'APPROVED' },
  { primeActualWeight: 2000, reweigh: { weight: 3000 }, status: 'APPROVED' },
  {
    ppmShipment: {
      weightTickets: [
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
      ],
    },
    status: 'APPROVED',
  },
  {
    ppmShipment: {
      weightTickets: [
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
      ],
    },
    status: 'APPROVED',
  },
];

const defaultProps = {
  netWeightRemarks: '',
  weightAllowance: 9000,
  weightTicket: { fullWeight: 1200, emptyWeight: 200 },
  shipments,
};

const excessWeight = {
  ...defaultProps,
  weightTicket: { fullWeight: 2700, emptyWeight: 200 },
  shipments: [
    ...shipments,
    {
      ppmShipment: {
        weightTickets: [createCompleteWeightTicket({}, { fullWeight: 2700, emptyWeight: 200 })],
      },
      status: 'APPROVED',
    },
  ],
  // moveWeightTotal = 9500
};

const reduceWeight = {
  ...defaultProps, // moveWeightTotal = 10000
  shipments: [
    ...shipments,
    {
      ppmShipment: {
        weightTickets: [
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
          createCompleteWeightTicket({}, { fullWeight: 1200, emptyWeight: 200 }),
        ],
      },
      status: 'APPROVED',
    },
  ],
};

describe('EditNetPPMWeight', () => {
  it('renders weights and edit button initially', async () => {
    await render(<EditPPMNetWeight {...defaultProps} />);
    expect(screen.getByText('Net weight')).toBeInTheDocument();
    expect(screen.getByText('| original weight')).toBeInTheDocument();
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.getAllByRole('button', { name: 'Edit' }));
  });
  it('renders correct labels and no calculations if there is no excess weight', async () => {
    await render(<EditPPMNetWeight {...defaultProps} />);
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.queryByText('| to fit within weight allowance')).not.toBeInTheDocument();
    expect(screen.queryByText('| to reduce excess weight')).not.toBeInTheDocument();
  });
  it('renders correct labels and calculations if there is excess weight', async () => {
    await render(<EditPPMNetWeight {...excessWeight} />);
    expect(screen.getByText('| to fit within weight allowance')).toBeInTheDocument();
    expect(screen.getAllByText('2,500 lbs')).toHaveLength(2);
    expect(screen.getByText('-500 lbs')).toBeInTheDocument();
  });
  it('renders correct labels and calculations if there is excess weight and reducing the full weight ticket is necessary', async () => {
    await render(<EditPPMNetWeight {...reduceWeight} />);
    expect(screen.getByText('| to reduce excess weight')).toBeInTheDocument();
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.getByText('-1,000 lbs')).toBeInTheDocument();
  });
  it('renders warning when if move total is higher than weight allowance', async () => {
    await render(<EditPPMNetWeight {...excessWeight} />);
    expect(await screen.findByTestId('warning')).toBeInTheDocument();
  });
  describe('when editing PPM Net weight', () => {
    it('renders editing form when the edit button is clicked', async () => {
      render(<EditPPMNetWeight {...defaultProps} />);
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      // Net weight form
      expect(screen.getByText('Net weight')).toBeInTheDocument();
      expect(screen.getByText('| original weight')).toBeInTheDocument();
      expect(screen.getAllByText('1,000 lbs')).toHaveLength(1);
      // Buttons
      expect(screen.getByRole('button', { name: 'Save changes' }));
      expect(screen.getByRole('button', { name: 'Cancel' }));
      // Input
      expect(await screen.findByTestId('weightInput')).toBeInTheDocument();
    });
    it('renders additional hints for excess weight', async () => {
      render(<EditPPMNetWeight {...excessWeight} />);
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      // Calculations for excess weight
      expect(screen.getByText('Move weight (total)')).toBeInTheDocument();
      expect(screen.getByText('9,500 lbs')).toBeInTheDocument();
      expect(screen.getByText('Weight allowance')).toBeInTheDocument();
      expect(screen.getByText('9,000 lbs')).toBeInTheDocument();
      expect(screen.getByText('Excess weight (total)')).toBeInTheDocument();
      expect(screen.getByText('500 lbs')).toBeInTheDocument();
    });
    it('disables the save button if save remarks or weight field is empty', async () => {
      render(<EditPPMNetWeight {...excessWeight} />);
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      expect(await screen.getByRole('button', { name: 'Save changes' })).toBeDisabled();
    });
    it('shows an error if there is incomplete fields in the form', async () => {
      render(<EditPPMNetWeight {...excessWeight} />);
      // Weight Input
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      const textInput = await screen.findByTestId('weightInput');
      await act(() => userEvent.clear(textInput));
      expect(screen.getByText('Required'));
      await act(() => userEvent.type(textInput, '1000'));
      // Remarks
      const remarksField = await screen.findByTestId('formRemarks');
      await act(() => userEvent.clear(remarksField));
      await fireEvent.blur(remarksField);
      expect(await screen.findByTestId('errorIndicator')).toBeInTheDocument();
      expect(screen.getByText('Required'));
    });
    // it('saves changes when editing the form', async () => {
    //   render(<EditPPMNetWeight {...reduceWeight} />);
    //   await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
    //   const textInput = await screen.findByTestId('weightInput');
    //   await act(() => userEvent.type(textInput, '1000'));
    //   const remarksField = await screen.findByTestId('formRemarks');
    //   await act(() => userEvent.type(remarksField, 'Reduced by as much as I can'));
    //   await act( () => userEvent.click(screen.getByRole('button', { name: 'Save changes' })))
    //   expect(screen.getByText('0 lbs')).toBeInTheDocument();
    //   expect(screen.getByText('Remarks')).toBeInTheDocument();
    //   expect(screen.getByText('Reduced by as much as I can')).toBeInTheDocument();
    // });
  });
});
