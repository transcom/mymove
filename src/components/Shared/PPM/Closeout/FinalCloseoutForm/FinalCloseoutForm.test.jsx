import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { v4 } from 'uuid';

import FinalCloseoutForm from 'components/Shared/PPM/Closeout/FinalCloseoutForm/FinalCloseoutForm';
import { createPPMShipmentWithFinalIncentive } from 'utils/test/factories/ppmShipment';
import { createCompleteMovingExpense } from 'utils/test/factories/movingExpense';
import { createCompleteProGearWeightTicket } from 'utils/test/factories/proGearWeightTicket';
import { createCompleteWeightTicket } from 'utils/test/factories/weightTicket';
import { APP_NAME } from 'constants/apps';

beforeEach(() => {
  jest.clearAllMocks();
});

const defaultPropsOffice = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  initialValues: { date: '2022-11-01', signature: '' },
  affiliation: 'ARMY',
  selectedMove: {
    closeout_office: {
      name: 'Altus AFB',
    },
  },
  appName: APP_NAME.OFFICE,
};

const defaultPropsCustomer = {
  onSubmit: jest.fn(),
  onBack: jest.fn(),
  initialValues: { date: '2022-11-01', signature: '' },
  affiliation: 'ARMY',
  selectedMove: {
    closeout_office: {
      name: 'Altus AFB',
    },
  },
  appName: APP_NAME.MYMOVE,
};

describe('FinalCloseoutForm component', () => {
  const prepListSearchForItem = (listOfItems) => (text) => listOfItems.find((item) => item.textContent === text);

  it('displays final incentive and shipment totals', () => {
    const serviceMemberId = v4();

    const weightTicket = createCompleteWeightTicket({ serviceMemberId }, { emptyWeight: 14000, fullWeight: 18000 });
    const movingExpense = createCompleteMovingExpense({ serviceMemberId }, { amount: 30000 });
    const proGearWeightTicket = createCompleteProGearWeightTicket({ serviceMemberId }, { weight: 1500 });

    const mtoShipment = createPPMShipmentWithFinalIncentive({
      ppmShipment: {
        advanceAmountReceived: 90000000,
        finalIncentive: 200000000,
        weightTickets: [weightTicket],
        movingExpenses: [movingExpense],
        proGearWeightTickets: [proGearWeightTicket],
      },
    });

    render(<FinalCloseoutForm mtoShipment={mtoShipment} {...defaultPropsOffice} />);

    expect(
      screen.getByRole('heading', { level: 2, name: 'Your final estimated incentive: $2,000,000.00' }),
    ).toBeInTheDocument();

    expect(screen.getByRole('heading', { level: 3, name: 'This PPM includes:' })).toBeInTheDocument();

    const findListItemWithText = prepListSearchForItem(screen.getAllByRole('listitem'));

    expect(findListItemWithText('4,000 lbs total net weight')).toBeInTheDocument();
    expect(findListItemWithText('1,500 lbs of pro-gear')).toBeInTheDocument();
    expect(findListItemWithText('$300.00 in expenses claimed')).toBeInTheDocument();

    expect(
      screen.getByRole('heading', { level: 2, name: 'Your actual payment will probably be lower' }),
    ).toBeInTheDocument();

    expect(findListItemWithText('minus any advances you were given (you received $900,000.00)')).toBeInTheDocument();

    expect(screen.getByText('Altus AFB', { exact: false })).toBeInTheDocument();
  });

  it('properly handles multiple weight tickets, pro gear weight tickets, and moving expenses', () => {
    const serviceMemberId = v4();

    const weightTickets = [
      { emptyWeight: 14000, fullWeight: 18000 },
      { emptyWeight: 14000, fullWeight: 17000 },
    ].map((fieldOverrides) => createCompleteWeightTicket({ serviceMemberId }, fieldOverrides));

    const movingExpenses = [{ amount: 30000 }, { amount: 50000 }, { amount: 40000 }].map((fieldOverrides) =>
      createCompleteMovingExpense({ serviceMemberId }, fieldOverrides),
    );

    const proGearWeightTickets = [{ weight: 750 }, { weight: 750 }].map((fieldOverrides) =>
      createCompleteProGearWeightTicket({ serviceMemberId }, fieldOverrides),
    );

    const mtoShipment = createPPMShipmentWithFinalIncentive({
      ppmShipment: {
        weightTickets,
        movingExpenses,
        proGearWeightTickets,
      },
    });

    render(<FinalCloseoutForm mtoShipment={mtoShipment} {...defaultPropsOffice} />);

    const findListItemWithText = prepListSearchForItem(screen.getAllByRole('listitem'));

    expect(findListItemWithText('7,000 lbs total net weight')).toBeInTheDocument();
    expect(findListItemWithText('1,500 lbs of pro-gear')).toBeInTheDocument();
    expect(findListItemWithText('$1,200.00 in expenses claimed')).toBeInTheDocument();
  });

  it('calls onBack func when "Back" button is clicked', async () => {
    const mtoShipment = createPPMShipmentWithFinalIncentive();

    render(<FinalCloseoutForm mtoShipment={mtoShipment} {...defaultPropsOffice} />);

    await userEvent.click(screen.getByRole('button', { name: 'Back' }));

    expect(defaultPropsOffice.onBack).toHaveBeenCalled();
  });

  describe('Customer side specific tests', () => {
    it('displays signature box and calls onSubmit func when "Submit PPM Documentation" button is clicked', async () => {
      const mtoShipment = createPPMShipmentWithFinalIncentive();
      const modifiedProps = {
        ...defaultPropsCustomer,
        initialValues: {
          ...defaultPropsCustomer.initialValues,
          signature: 'Grace Griffin',
        },
      };

      render(<FinalCloseoutForm mtoShipment={mtoShipment} {...modifiedProps} />);

      const signatureField = screen.getByRole('textbox', { name: 'Signature' });
      await waitFor(() => expect(signatureField).toHaveValue('Grace Griffin'));

      const saveButton = screen.getByRole('button', { name: 'Submit PPM Documentation' });
      await userEvent.click(saveButton);
      await waitFor(() => expect(modifiedProps.onSubmit).toHaveBeenCalled());
    });
  });
});
