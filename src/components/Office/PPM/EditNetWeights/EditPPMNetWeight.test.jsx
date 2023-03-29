import React from 'react';
import { render, screen, act, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditPPMNetWeight from './EditPPMNetWeight';

import { patchWeightTicket } from 'services/ghcApi';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { ReactQueryWrapper } from 'testUtils';

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
}));

jest.mock('services/ghcApi', () => ({
  ...jest.requireActual('services/ghcApi'),
  patchWeightTicket: jest.fn(),
}));

const shipments = [
  // moveWeightTotal = 7000
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
  // moveWeightTotal = 7000
  netWeightRemarks: '',
  weightAllowance: 9000,
  weightTicket: { fullWeight: 1200, emptyWeight: 200 },
  shipments,
};

const excessWeight = {
  // moveWeightTotal = 9500
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
};

const reduceWeight = {
  // moveWeightTotal = 10000
  ...defaultProps,
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
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...defaultProps} />
      </ReactQueryWrapper>,
    );
    expect(screen.getByText('Net weight')).toBeInTheDocument();
    expect(screen.getByText('| original weight')).toBeInTheDocument();
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.getAllByRole('button', { name: 'Edit' }));
  });
  it('renders correct labels and no calculations if there is no excess weight', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...defaultProps} />
      </ReactQueryWrapper>,
    );
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.queryByText('| to fit within weight allowance')).not.toBeInTheDocument();
    expect(screen.queryByText('| to reduce excess weight')).not.toBeInTheDocument();
  });
  it('renders correct labels and calculations if there is excess weight', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...excessWeight} />
      </ReactQueryWrapper>,
    );
    expect(screen.getByText('| to fit within weight allowance')).toBeInTheDocument();
    expect(screen.getAllByText('2,500 lbs')).toHaveLength(2);
    expect(screen.getByText('-500 lbs')).toBeInTheDocument();
  });
  it('renders correct labels and calculations if there is excess weight and reducing the full weight ticket is necessary', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...reduceWeight} />
      </ReactQueryWrapper>,
    );
    expect(screen.getByText('| to reduce excess weight')).toBeInTheDocument();
    expect(screen.getAllByText('1,000 lbs')).toHaveLength(2);
    expect(screen.getByText('-1,000 lbs')).toBeInTheDocument();
  });
  it('renders warning when if move total is higher than weight allowance', async () => {
    await render(
      <ReactQueryWrapper>
        <EditPPMNetWeight {...excessWeight} />
      </ReactQueryWrapper>,
    );
    expect(await screen.findByTestId('warning')).toBeInTheDocument();
  });
  describe('when editing PPM Net weight', () => {
    it('renders editing form when the edit button is clicked', async () => {
      await render(
        <ReactQueryWrapper>
          <EditPPMNetWeight {...defaultProps} />
        </ReactQueryWrapper>,
      );
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
      await render(
        <ReactQueryWrapper>
          <EditPPMNetWeight {...excessWeight} />
        </ReactQueryWrapper>,
      );
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      // Calculations for excess weight
      expect(screen.getByText('Move weight (total)')).toBeInTheDocument();
      expect(screen.getByText('9,500 lbs')).toBeInTheDocument();
      expect(screen.getByText('Weight allowance')).toBeInTheDocument();
      expect(screen.getByText(/9,000 lbs/)).toBeInTheDocument();
      expect(screen.getByText('Excess weight (total)')).toBeInTheDocument();
      expect(screen.getByText('500 lbs')).toBeInTheDocument();
    });
    it('disables the save button if save remarks or weight field is empty', async () => {
      await render(
        <ReactQueryWrapper>
          <EditPPMNetWeight {...excessWeight} />
        </ReactQueryWrapper>,
      );
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      expect(await screen.getByRole('button', { name: 'Save changes' })).toBeDisabled();
    });
    it('shows validation errors in the form', async () => {
      await render(
        <ReactQueryWrapper>
          <EditPPMNetWeight {...excessWeight} />
        </ReactQueryWrapper>,
      );
      // Weight Input is required
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      const textInput = await screen.findByTestId('weightInput');
      await act(() => userEvent.clear(textInput));
      expect(screen.getByText('Required'));

      // Weight input cannot be greater than full weight
      await act(() => userEvent.type(textInput, '10000'));
      expect(screen.getByText('Net weight must be less than or equal to the full weight'));
      await act(() => userEvent.type(textInput, '1000'));
      // Remarks Input is required
      const remarksField = await screen.findByTestId('formRemarks');
      await act(() => userEvent.clear(remarksField));
      await fireEvent.blur(remarksField);
      expect(await screen.findByTestId('errorIndicator')).toBeInTheDocument();
      expect(screen.getByText('Required'));
    });
    it('saves changes when editing the form', async () => {
      const netWeightRemarks = 'Reduced by as much as I can';
      const adjustedNetWeight = 0;
      render(
        <ReactQueryWrapper>
          <EditPPMNetWeight {...reduceWeight} />
        </ReactQueryWrapper>,
      );

      await act(() => userEvent.click(screen.getByRole('button', { name: 'Edit' })));
      const textInput = await screen.findByTestId('weightInput');
      await act(() => userEvent.clear(textInput));
      await act(() => userEvent.type(textInput, adjustedNetWeight.toString(10)));
      const remarksField = await screen.findByTestId('formRemarks');
      await act(() => userEvent.type(remarksField, netWeightRemarks));
      await act(() => userEvent.click(screen.getByRole('button', { name: 'Save changes' })));

      await waitFor(() => {
        expect(patchWeightTicket.mock.calls.length).toBe(1);
      });
      expect(patchWeightTicket.mock.calls[0][0].payload.adjustedNetWeight).toBe(adjustedNetWeight);
      expect(patchWeightTicket.mock.calls[0][0].payload.netWeightRemarks).toBe(netWeightRemarks);
    });
  });
});
