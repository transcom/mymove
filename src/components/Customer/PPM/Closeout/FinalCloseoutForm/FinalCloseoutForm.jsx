import React from 'react';
import { func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCentsTruncateWhole, formatWeight } from 'utils/formatters';
import {
  calculateTotalMovingExpensesAmount,
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForWeightTickets,
} from 'utils/ppmCloseout';

const FinalCloseoutForm = ({ mtoShipment, onBack, onSubmit }) => {
  const totalNetWeight = calculateTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets || []);

  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets || [],
  );

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses || []);

  const isValid = false;
  const isSubmitting = true;

  return (
    <>
      <h2>
        Your final estimated incentive: $
        {formatCentsTruncateWhole(mtoShipment?.ppmShipment?.finalEstimatedIncentive || 0)}
      </h2>
      <div>
        <p>Your incentive is calculated using:</p>
        <ul>
          <li>weight: verified net weight of your completed PPM</li>
          <li>distance: starting and ending ZIP codes</li>
          <li>date: when you started moving your PPM</li>
          <li>
            allowances: your total weight allowance for your whole move, including all shipments, both PPMs and
            government-funded (such as HHGs)
          </li>
        </ul>
      </div>

      <h2>This PPM includes:</h2>
      <ul>
        <li>{formatWeight(totalNetWeight)} total net weight</li>
        <li>{formatWeight(totalProGearWeight)} of pro-gear</li>
        <li>${formatCentsTruncateWhole(totalExpensesClaimed)} in expenses claimed</li>
      </ul>

      <h2>Your actual payment will probably be lower</h2>
      <div>
        <p>Your final payment will be:</p>
        <ul>
          <li>based on your final incentive</li>
          <li>modified by expenses submitted (authorized expenses reduce your tax burden)</li>
          <li>minus any taxes withheld (the IRS considers your incentive to be taxable income)</li>
          <li>
            minus any advances you were given (you received $
            {formatCentsTruncateWhole(mtoShipment?.ppmShipment?.advanceAmountReceived || 0)})
          </li>
          <li>plus any reimbursements you receive</li>
        </ul>
        <p>
          Verified expenses reduce the taxable income you report to the IRS on form W-2. They may not be claimed again
          as moving expenses. Federal tax withholding will be deducted from the profit (entitlement less eligible
          operating expenses.)
        </p>
      </div>

      <div className={ppmStyles.buttonContainer}>
        <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
          Finish Later
        </Button>
        <Button className={ppmStyles.saveButton} type="button" onClick={onSubmit} disabled={!isValid || isSubmitting}>
          Submit PPM Documentation
        </Button>
      </div>
    </>
  );
};

FinalCloseoutForm.prototypes = {
  mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default FinalCloseoutForm;
