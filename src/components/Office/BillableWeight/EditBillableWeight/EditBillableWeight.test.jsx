import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

import EditBillableWeight from './EditBillableWeight';

import { formatWeight } from 'utils/formatters';

jest.mock('formik', () => ({
  ...jest.requireActual('formik'),
}));

describe('EditBillableWeight', () => {
  it('renders weight and edit button intially', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      editEntity: () => {},
    };

    render(<EditBillableWeight {...defaultProps} />);

    expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
    expect(screen.getByText(defaultProps.title)).toBeInTheDocument();
    expect(screen.queryByText('Remarks')).toBeNull();
    expect(screen.getByText(formatWeight(defaultProps.maxBillableWeight))).toBeInTheDocument();
  });

  it('should show fields are required when empty', () => {});

  it('renders billable weight justification', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      editEntity: () => {},
      billableWeightJustification: 'Reduced billable weight to cap at 110% of estimated.',
    };

    render(<EditBillableWeight {...defaultProps} />);
    expect(screen.getByText(defaultProps.billableWeightJustification)).toBeInTheDocument();
  });

  it('renders max billable weight view', async () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      editEntity: () => {},
    };

    render(<EditBillableWeight {...defaultProps} />);
    userEvent.click(await screen.findByRole('button', { name: 'Edit' }));
    expect(await screen.findByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
    expect(screen.getByText('| weight allowance')).toBeInTheDocument();
    expect(screen.getByText('| 110% of total estimated weight')).toBeInTheDocument();
  });

  it('renders edit billable weight view', async () => {
    const defaultProps = {
      title: 'Billable weight',
      originalWeight: 10000,
      estimatedWeight: 13000,
      maxBillableWeight: 6000,
      billableWeight: 14400,
      totalBillableWeight: 11000,
      editEntity: () => {},
    };

    render(<EditBillableWeight {...defaultProps} />);
    userEvent.click(await screen.findByRole('button', { name: 'Edit' }));
    expect(await screen.findByText(formatWeight(defaultProps.originalWeight))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
    expect(
      screen.getByText(
        formatWeight(defaultProps.maxBillableWeight - defaultProps.totalBillableWeight + defaultProps.billableWeight),
      ),
    ).toBeInTheDocument();
    expect(screen.getByText('| original weight')).toBeInTheDocument();
    expect(screen.getByText('| 110% of total estimated weight')).toBeInTheDocument();
    expect(screen.getByText('| to fit within max billable weight')).toBeInTheDocument();
  });

  describe('hint text for max billable weight', () => {
    it('should not show the 110% of total estimated weight hint text if the estimated weight is missing', async () => {
      const defaultProps = {
        title: 'Max billable weight',
        weightAllowance: 8000,
        maxBillableWeight: 10000,
        editEntity: () => {},
      };

      render(<EditBillableWeight {...defaultProps} />);
      userEvent.click(screen.getByRole('button', { name: 'Edit' }));
      await waitFor(() => {
        expect(screen.getByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();
      });
      expect(screen.getByText('| weight allowance')).toBeInTheDocument();
      expect(screen.queryByText('| 110% of total estimated weight')).not.toBeInTheDocument();
    });
  });

  describe('hint text for billable weight', () => {
    describe('110% of total estimated weight hint text', () => {
      it('should show if the billable weight is greater than the estimated weight * 110%', async () => {
        const defaultProps = {
          title: 'Billable weight',
          originalWeight: 10000,
          estimatedWeight: 13000,
          maxBillableWeight: 6000,
          billableWeight: 14600,
          totalBillableWeight: 11000,
          editEntity: () => {},
        };

        render(<EditBillableWeight {...defaultProps} />);
        userEvent.click(screen.getByRole('button', { name: 'Edit' }));
        expect(await screen.findByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
        expect(screen.getByText('| 110% of total estimated weight')).toBeInTheDocument();
      });

      it('should not show if the estimated weight is missing', async () => {
        const defaultProps = {
          title: 'Billable weight',
          originalWeight: 10000,
          maxBillableWeight: 6000,
          billableWeight: 14600,
          totalBillableWeight: 11000,
          editEntity: () => {},
        };

        render(<EditBillableWeight {...defaultProps} />);
        userEvent.click(screen.getByRole('button', { name: 'Edit' }));
        await waitFor(() => {
          expect(screen.queryByText('| 110% of total estimated weight')).not.toBeInTheDocument();
        });
      });

      it('should not show if the billable weight is less than the estimated weight * 110%', async () => {
        const defaultProps = {
          title: 'Billable weight',
          originalWeight: 10000,
          estimatedWeight: 13000,
          maxBillableWeight: 6000,
          billableWeight: 14000,
          totalBillableWeight: 11000,
          editEntity: () => {},
        };

        render(<EditBillableWeight {...defaultProps} />);
        userEvent.click(screen.getByRole('button', { name: 'Edit' }));
        await waitFor(() => {
          expect(screen.queryByText(formatWeight(defaultProps.estimatedWeight * 1.1))).not.toBeInTheDocument();
        });
        expect(screen.queryByText('| 110% of total estimated weight')).not.toBeInTheDocument();
      });
      it('should not show if this is an NTS-release shipment', async () => {
        const defaultProps = {
          title: 'Billable weight',
          originalWeight: 10000,
          estimatedWeight: 13000,
          maxBillableWeight: 6000,
          billableWeight: 14600,
          totalBillableWeight: 11000,
          editEntity: () => {},
          isNTSRShipment: true,
        };

        render(<EditBillableWeight {...defaultProps} />);
        userEvent.click(screen.getByRole('button', { name: 'Edit' }));
        await waitFor(() => {
          expect(screen.queryByText(formatWeight(defaultProps.estimatedWeight * 1.1))).not.toBeInTheDocument();
        });
        expect(screen.queryByText('| 110% of total estimated weight')).not.toBeInTheDocument();
      });
    });

    describe('to fit within the max billable weight hint text', () => {
      it('should show if the billable weight is greater than the max billable weight and greater than the estimated weight * 110%', async () => {
        const defaultProps = {
          title: 'Billable weight',
          originalWeight: 11000,
          estimatedWeight: 13000,
          maxBillableWeight: 6000,
          billableWeight: 15000,
          totalBillableWeight: 11000,
          editEntity: () => {},
        };

        const fitWithinValue = formatWeight(
          defaultProps.maxBillableWeight - defaultProps.totalBillableWeight + defaultProps.billableWeight,
        );

        render(<EditBillableWeight {...defaultProps} />);
        userEvent.click(screen.getByRole('button', { name: 'Edit' }));
        expect(await screen.findByText(fitWithinValue)).toBeInTheDocument();
        expect(screen.getByText('| to fit within max billable weight')).toBeInTheDocument();
      });

      it('should not show if the billable weight is less than the max billable weight and less than the estimated weight * 110%', async () => {
        const defaultProps = {
          title: 'Billable weight',
          originalWeight: 10000,
          estimatedWeight: 13000,
          maxBillableWeight: 6000,
          billableWeight: 12000,
          totalBillableWeight: 11000,
          editEntity: () => {},
        };

        const fitWithinValue = formatWeight(
          defaultProps.maxBillableWeight - defaultProps.totalBillableWeight + defaultProps.billableWeight,
        );

        render(<EditBillableWeight {...defaultProps} />);
        userEvent.click(screen.getByRole('button', { name: 'Edit' }));
        await waitFor(() => {
          expect(screen.queryByText(fitWithinValue)).not.toBeInTheDocument();
        });
        expect(screen.queryByText('| to fit within max billable weight')).not.toBeInTheDocument();
      });
    });
  });

  it('clicking edit button shows different view', async () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      editEntity: () => {},
    };

    render(<EditBillableWeight {...defaultProps} />);

    userEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.queryByText('Edit')).toBeNull();
    // weights
    expect(await screen.findByText(formatWeight(defaultProps.weightAllowance))).toBeInTheDocument();
    expect(screen.getByText(formatWeight(defaultProps.estimatedWeight * 1.1))).toBeInTheDocument();
    // buttons
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();
  });

  it('should be able to toggle between views', () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      editEntity: () => {},
    };

    render(<EditBillableWeight {...defaultProps} />);
    userEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.queryByText('Edit')).toBeNull();
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancel' })).toBeInTheDocument();

    userEvent.click(screen.getByRole('button', { name: 'Cancel' }));
    expect(screen.queryByText('Edit')).toBeInTheDocument();
    expect(screen.queryByText('Save changes')).toBeNull();
    expect(screen.queryByText('Cancel')).toBeNull();
  });

  it('should call editEntity with data', async () => {
    const mockEditEntity = jest.fn();
    const newBillableWeight = 5000;
    const newBillableWeightJustification = 'some remarks';
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      editEntity: mockEditEntity,
    };

    render(<EditBillableWeight {...defaultProps} />);
    userEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.queryByText('Edit')).toBeNull();

    const textInput = await screen.findByTestId('textInput');
    userEvent.clear(textInput);
    userEvent.type(textInput, '5000');

    const remarksInput = await screen.findByTestId('remarks');
    userEvent.type(remarksInput, newBillableWeightJustification);
    userEvent.click(await screen.findByRole('button', { name: 'Save changes' }));

    await waitFor(() => {
      expect(mockEditEntity.mock.calls.length).toBe(1);
    });
    expect(mockEditEntity.mock.calls[0][0].billableWeight).toBe(String(newBillableWeight));
    expect(mockEditEntity.mock.calls[0][0].billableWeightJustification).toBe(newBillableWeightJustification);
  });

  it('should disable save button if remarks or billable weight fields are empty', async () => {
    const defaultProps = {
      title: 'Max billable weight',
      weightAllowance: 8000,
      estimatedWeight: 13750,
      maxBillableWeight: 10000,
      billableWeightJustification: 'some remarks',
      editEntity: () => {},
    };

    render(<EditBillableWeight {...defaultProps} />);
    userEvent.click(screen.getByRole('button', { name: 'Edit' }));
    expect(screen.queryByText('Edit')).toBeNull();
    userEvent.clear(screen.getByTestId('textInput'));
    userEvent.clear(screen.getByTestId('remarks'));
    (await screen.findByTestId('remarks')).blur();
    expect(screen.getByText('Required')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Save changes' })).toBeDisabled();
  });
});
