import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';

import styles from './FinalCloseoutForm.module.scss';

import ppmStyles from 'components/Shared/PPM/PPM.module.scss';
import { ShipmentShape } from 'types/shipment';
import { MoveShape } from 'types/customerShapes';
import { formatCents, formatWeight } from 'utils/formatters';
import { calculateTotalMovingExpensesAmount } from 'utils/ppmCloseout';
import {
  calculateTotalNetWeightForProGearWeightTickets,
  getTotalNetWeightForWeightTickets,
} from 'utils/shipmentWeights';
import affiliations from 'content/serviceMemberAgencies';
import { APP_NAME } from 'constants/apps';

const FinalCloseoutForm = ({ initialValues, mtoShipment, onBack, onSubmit, affiliation, selectedMove, appName }) => {
  const totalNetWeight = getTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets);
  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets,
  );

  const canChoosePPMLocation =
    affiliation === affiliations.ARMY ||
    affiliation === affiliations.AIR_FORCE ||
    affiliation === affiliations.SPACE_FORCE;

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses);

  return (
    <Formik initialValues={initialValues} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit }) => (
        <div className={styles.FinalCloseoutForm}>
          <h2>Your final estimated incentive: ${formatCents(mtoShipment?.ppmShipment?.finalIncentive || 0)}</h2>
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
            {canChoosePPMLocation && (
              <p>
                Your closeout office for your PPM(s) is{' '}
                {selectedMove?.closeout_office?.name ? selectedMove.closeout_office.name : ''}. This is where your PPM
                paperwork will be reviewed before you can submit it to finance to receive your incentive.
              </p>
            )}
          </div>

          <div className={ppmStyles.buttonContainer}>
            <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
              {appName === APP_NAME.OFFICE ? 'Back' : 'Return To Homepage'}
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
    affiliation: PropTypes.string,
    signature: PropTypes.string,
    date: PropTypes.string,
  }).isRequired,
  mtoShipment: ShipmentShape.isRequired,
  onBack: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  affiliation: PropTypes.string.isRequired,
  selectedMove: MoveShape,
};

FinalCloseoutForm.defaultProps = {
  selectedMove: {},
};

export default FinalCloseoutForm;
