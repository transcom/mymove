import React from 'react';
import { PropTypes, func } from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { Formik } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';

import styles from './FinalCloseoutForm.module.scss';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCents, formatWeight } from 'utils/formatters';
import {
  calculateTotalMovingExpensesAmount,
  calculateTotalNetWeightForProGearWeightTickets,
  calculateTotalNetWeightForWeightTickets,
} from 'utils/ppmCloseout';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { Form } from 'components/form/Form';
import formStyles from 'styles/form.module.scss';
import { requiredW2AddressSchema } from 'utils/validation';
import { W2AddressShape } from 'types/address';

const FinalCloseoutForm = ({ mtoShipment, onBack }) => {
  const totalNetWeight = calculateTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets);

  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets,
  );

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses);

  const formFieldsName = 'w2_address';
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredW2AddressSchema.required(),
  });
  const initialValues = {
    [formFieldsName]: {
      streetAddress1: mtoShipment?.ppmShipment?.w2Address?.streetAddress1 || '',
      streetAddress2: mtoShipment?.ppmShipment?.w2Address?.streetAddress2 || '',
      city: mtoShipment?.ppmShipment?.w2Address?.city || '',
      state: mtoShipment?.ppmShipment?.w2Address?.state || '',
      postalCode: mtoShipment?.ppmShipment?.w2Address?.postalCode || '',
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

      <Formik initialValues={initialValues} validationSchema={validationSchema}>
        {({ isValid, isSubmitting, handleSubmit }) => {
          return (
            <>
              <div className={classnames(ppmStyles.formContainer)}>
                <Form className={classnames(ppmStyles.form, formStyles.form, styles.W2AddressForm)}>
                  <SectionWrapper className={classnames(formStyles.formSection)}>
                    <h2>W-2 address</h2>
                    <p>What is the address on your W-2?</p>
                    <AddressFields name={formFieldsName} className={classnames(styles.AddressFieldSet)} />
                  </SectionWrapper>
                </Form>
              </div>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Finish Later
                </Button>
                <Button
                  className={ppmStyles.saveButton}
                  type="button"
                  onClick={handleSubmit}
                  disabled={!isValid || isSubmitting}
                >
                  Submit PPM Documentation
                </Button>
              </div>
            </>
          );
        }}
      </Formik>
    </div>
  );
};

FinalCloseoutForm.prototypes = {
  mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  formFieldsName: PropTypes.string.isRequired,
  initialValues: W2AddressShape.isRequired,
};

export default FinalCloseoutForm;
