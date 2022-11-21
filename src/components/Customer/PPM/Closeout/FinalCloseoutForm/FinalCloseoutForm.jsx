import React from 'react';
import PropTypes from 'prop-types';
import { Button, Grid } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import styles from './FinalCloseoutForm.module.scss';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCents, formatWeight } from 'utils/formatters';
import {
  calculateTotalMovingExpensesAmount,
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForWeightTickets,
} from 'utils/ppmCloseout';
import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';

const validationSchema = Yup.object().shape({
  signature: Yup.string().required('Required'),
  date: Yup.string(),
});

const FinalCloseoutForm = ({ initialValues, mtoShipment, onBack, onSubmit }) => {
  const totalNetWeight = calculateTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets);

  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets,
  );

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses);

  return (
    <Formik validationSchema={validationSchema} initialValues={initialValues} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit }) => (
        <div className={styles.FinalCloseoutForm}>
          <h2>
            Your final estimated incentive: ${formatCents(mtoShipment?.ppmShipment?.finalEstimatedIncentive || 0)}
          </h2>
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
                  your total weight allowance for your whole move, including all shipments, both PPMs and
                  government-funded (such as HHGs)
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
              Verified expenses reduce the taxable income you report to the IRS on form W-2. They may not be claimed
              again as moving expenses. Federal tax withholding will be deducted from the profit (entitlement less
              eligible operating expenses.)
            </p>
          </div>

          <SectionWrapper>
            <h2>Customer agreement</h2>
            <p>I certify that any expenses claimed in this application were legitimately incurred during my PPM.</p>
            <p>
              Failure to furnish data may result in partial or total denial of claim and/or improper tax application.
            </p>
            <p>
              I understand the penalty for willfully making a false statement of claim is a maximum fine of $10,000,
              maximum imprisonment of five years, or both (U.S.C, Title 18, Section 287).
            </p>
            <Grid row>
              <Grid desktop={{ col: 6 }}>
                <TextField label="Signature" id="signature" name="signature" />
              </Grid>
              <Grid desktop={{ col: 4, offset: 1 }}>
                <TextField label="Date" id="date" name="date" disabled />
              </Grid>
            </Grid>
          </SectionWrapper>

          <div className={ppmStyles.buttonContainer}>
            <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
              Return To Homepage
            </Button>
            <Button
              className={ppmStyles.saveButton}
              type="submit"
              onClick={handleSubmit}
              disabled={!isValid || isSubmitting}
            >
              Submit PPM Documentation
            </Button>
          </div>
        </div>
      )}
    </Formik>
  );
};

FinalCloseoutForm.propTypes = {
  initialValues: PropTypes.shape({
    signature: PropTypes.string,
    date: PropTypes.string,
  }).isRequired,
  mtoShipment: ShipmentShape.isRequired,
  onBack: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};

export default FinalCloseoutForm;
