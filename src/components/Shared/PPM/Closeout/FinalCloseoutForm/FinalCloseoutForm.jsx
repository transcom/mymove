import React, { useEffect, useState } from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';

import styles from './FinalCloseoutForm.module.scss';

import ppmStyles from 'components/Shared/PPM/PPM.module.scss';
import formStyles from 'styles/form.module.scss';
import { ShipmentShape } from 'types/shipment';
import { MoveShape } from 'types/customerShapes';
import { formatCents, formatWeight } from 'utils/formatters';
import {
  calculateTotalMovingExpensesAmount,
  getNonProGearWeightSPR,
  getProGearWeightSPR,
  getTotalPackageWeightSPR,
} from 'utils/ppmCloseout';
import {
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForGunSafeWeightTickets,
  getTotalNetWeightForWeightTickets,
} from 'utils/shipmentWeights';
import affiliations from 'content/serviceMemberAgencies';
import { APP_NAME } from 'constants/apps';
import { FEATURE_FLAG_KEYS, PPM_TYPES } from 'shared/constants';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { isBooleanFlagEnabled } from 'utils/featureFlags';

const FinalCloseoutForm = ({ initialValues, mtoShipment, onBack, onSubmit, affiliation, move, appName }) => {
  const [gunSafeEnabled, setGunSafeEnabled] = useState(false);
  useEffect(() => {
    const fetchData = async () => {
      setGunSafeEnabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.GUN_SAFE));
    };
    fetchData();
  }, []);
  const isCustomerPage = appName === APP_NAME.MYMOVE;
  const closeoutOfficeName = move?.closeoutOffice?.name || move?.closeout_office?.name || '';

  const validationSchema = Yup.object().shape({
    signature: isCustomerPage ? Yup.string().required('Required') : Yup.string(),
    date: Yup.string(),
  });
  const totalNetWeight = getTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets);
  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets,
  );
  const totalGunSafeWeight = calculateTotalNetWeightForGunSafeWeightTickets(
    mtoShipment?.ppmShipment?.gunSafeWeightTickets,
  );

  const canChoosePPMLocation =
    affiliation === affiliations.ARMY ||
    affiliation === affiliations.AIR_FORCE ||
    affiliation === affiliations.SPACE_FORCE;

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses);
  const ppmShipment = mtoShipment?.ppmShipment || {};
  const { ppmType } = ppmShipment;

  const totalNonProGearWeightSPR = getNonProGearWeightSPR(mtoShipment?.ppmShipment?.movingExpenses);
  const totalProGearWeightSPR = getProGearWeightSPR(mtoShipment?.ppmShipment?.movingExpenses);
  const totalWeightSPR = getTotalPackageWeightSPR(mtoShipment?.ppmShipment?.movingExpenses);

  return (
    <Formik validationSchema={validationSchema} initialValues={initialValues} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit }) => (
        <div className={styles.FinalCloseoutForm}>
          {ppmType !== PPM_TYPES.SMALL_PACKAGE && (
            <>
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
            </>
          )}

          <div className={styles.shipmentTotals}>
            <h3>This PPM includes:</h3>
            <ul>
              {ppmType === PPM_TYPES.SMALL_PACKAGE ? (
                <>
                  <li>${formatCents(totalExpensesClaimed)} in expenses claimed</li>
                  <li>{formatWeight(totalNonProGearWeightSPR)} total expense weight</li>
                  <li>{formatWeight(totalProGearWeightSPR)} total pro-gear weight</li>
                  <li>{formatWeight(totalWeightSPR)} in total weight</li>
                </>
              ) : (
                <>
                  <li>{formatWeight(totalNetWeight)} total net weight</li>
                  <li>{formatWeight(totalProGearWeight)} of pro-gear</li>
                  {gunSafeEnabled && <li>{formatWeight(totalGunSafeWeight)} of gun safe weight</li>}
                  <li>${formatCents(totalExpensesClaimed)} in expenses claimed</li>
                </>
              )}
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
                Your closeout office for your PPM(s) is {closeoutOfficeName}. This is where your PPM paperwork will be
                reviewed before you can submit it to finance to receive your incentive.
              </p>
            )}
          </div>

          {appName === APP_NAME.MYMOVE && (
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
              <div>
                <div className={styles.signatureField}>
                  <TextField label="Signature" id="signature" name="signature" />
                </div>
                <div className={styles.dateField}>
                  <TextField label="Date" id="date" name="date" disabled />
                </div>
              </div>
            </SectionWrapper>
          )}

          <div
            className={`${
              isCustomerPage ? ppmStyles.buttonContainer : `${formStyles.formActions} ${ppmStyles.buttonGroup}`
            }`}
          >
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
  move: MoveShape,
};

FinalCloseoutForm.defaultProps = {
  move: {},
};

export default FinalCloseoutForm;
