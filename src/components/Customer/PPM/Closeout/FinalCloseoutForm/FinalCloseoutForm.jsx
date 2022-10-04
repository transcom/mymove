import React from 'react';
import { func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';

import styles from './FinalCloseoutForm.module.scss';

import W2AddressForm from 'components/Customer/PPM/Closeout/W2AddressForm/W2AddressForm';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCents, formatWeight } from 'utils/formatters';
import {
  calculateTotalMovingExpensesAmount,
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForWeightTickets,
} from 'utils/ppmCloseout';

const FinalCloseoutForm = ({ mtoShipment, onBack, onSubmit }) => {
  const totalNetWeight = calculateTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets);

  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets,
  );

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses);

  const isValid = false;
  const isSubmitting = true;

  const formFieldsName = 'w2_address';
  const initialValues = {
    [formFieldsName]: {
      streetAddress1: '',
      streetAddress2: '',
      city: '',
      state: '',
      postalCode: '',
    },
  };

  return (
    <div className={styles.FinalCloseoutForm}>
      <h2>Your final estimated incentive: ${formatCents(mtoShipment?.ppmShipment?.finalEstimatedIncentive || 0)}</h2>
      <div className={styles.incentiveFactors}>
        <p className={styles.listDescription}>Your incentive is calculated using:</p>
        <dl>
          <div className={styles.definitionWrapper}>
            <dt>weight</dt>
            <dd>verified net weight of your completed PPM</dd>
          </div>
          <div className={styles.definitionWrapper}>
            <dt>distance</dt>
            <dd>starting and ending ZIP codes</dd>
          </div>
          <div className={styles.definitionWrapper}>
            <dt>date</dt>
            <dd>when you started moving your PPM</dd>
          </div>
          <div className={styles.definitionWrapper}>
            <dt>allowances</dt>
            <dd>
              your total weight allowance for your whole move, including all shipments, both PPMs and government-funded
              (such as HHGs)
            </dd>
          </div>
        </dl>
      </div>

      <div className={styles.shipmentTotals}>
        <h3>This PPM includes:</h3>
        <ul>
          <li>{formatWeight(totalNetWeight)} total net weight</li>
          <li>{formatWeight(totalProGearWeight)} of pro-gear</li>
          <li>${formatCents(totalExpensesClaimed)} in expenses claimed</li>
        </ul>
      </div>

      <h2>Your actual payment will probably be lower</h2>
      <div className={styles.finalPaymentFactors}>
        <p className={styles.listDescription}>Your final payment will be:</p>
        <ul>
          <li>based on your final incentive</li>
          <li>modified by expenses submitted (authorized expenses reduce your tax burden)</li>
          <li>minus any taxes withheld (the IRS considers your incentive to be taxable income)</li>
          <li>
            minus any advances you were given (you received $
            {formatCents(mtoShipment?.ppmShipment?.advanceAmountReceived || 0)})
          </li>
          <li>plus any reimbursements you receive</li>
        </ul>
        <p>
          Verified expenses reduce the taxable income you report to the IRS on form W-2. They may not be claimed again
          as moving expenses. Federal tax withholding will be deducted from the profit (entitlement less eligible
          operating expenses.)
        </p>
      </div>

      <W2AddressForm formFieldsName={formFieldsName} initialValues={initialValues} />

      <div className={ppmStyles.buttonContainer}>
        <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
          Return To Homepage
        </Button>
        <Button className={ppmStyles.saveButton} type="button" onClick={onSubmit} disabled={!isValid || isSubmitting}>
          Submit PPM Documentation
        </Button>
      </div>
    </div>
  );
};

FinalCloseoutForm.prototypes = {
  mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default FinalCloseoutForm;
