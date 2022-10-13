import React from 'react';
import { PropTypes, func, shape, string } from 'prop-types';
import { Button, TextInput, Label, FormGroup, Fieldset, ErrorMessage, Grid, Alert } from '@trussworks/react-uswds';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';
import classnames from 'classnames';

import styles from './FinalCloseoutForm.module.scss';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCents, formatWeight, formatSwaggerDate } from 'utils/formatters';
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

const FinalCloseoutForm = ({ mtoShipment, onBack, error, onSubmit }) => {
  const totalNetWeight = calculateTotalNetWeightForWeightTickets(mtoShipment?.ppmShipment?.weightTickets);

  const totalProGearWeight = calculateTotalNetWeightForProGearWeightTickets(
    mtoShipment?.ppmShipment?.proGearWeightTickets,
  );

  const totalExpensesClaimed = calculateTotalMovingExpensesAmount(mtoShipment?.ppmShipment?.movingExpenses);

  const formFieldsName = 'w2_address';
  const validationSchema = Yup.object().shape({
    [formFieldsName]: requiredW2AddressSchema.required(),
    signature: Yup.string().required('Required'),
    date: Yup.date().required(),
  });
  const initialValues = {
    [formFieldsName]: {
      streetAddress1: mtoShipment?.ppmShipment?.w2Address?.streetAddress1 || '',
      streetAddress2: mtoShipment?.ppmShipment?.w2Address?.streetAddress2 || '',
      city: mtoShipment?.ppmShipment?.w2Address?.city || '',
      state: mtoShipment?.ppmShipment?.w2Address?.state || '',
      postalCode: mtoShipment?.ppmShipment?.w2Address?.postalCode || '',
    },
    signature: '',
    date: formatSwaggerDate(new Date()),
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

      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {({ isValid, isSubmitting, handleSubmit, errors, touched }) => {
          const showSignatureError = !!(errors.signature && touched.signature);
          return (
            <div className={classnames(ppmStyles.formContainer)}>
              <Form className={classnames(ppmStyles.form, formStyles.form, styles.W2AddressForm)}>
                <SectionWrapper className={classnames(formStyles.formSection)}>
                  <h2>W-2 address</h2>
                  <p>What is the address on your W-2?</p>
                  <AddressFields name={formFieldsName} className={classnames(styles.AddressFieldSet)} />
                </SectionWrapper>

                <SectionWrapper className={classnames(formStyles.formSection, styles.signatureBox)}>
                  <h2>Customer agreement</h2>
                  <p>
                    I certify that any expenses claimed in this application were legitimately incurred during my PPM.
                  </p>
                  <p>
                    Failure to furnish data may result in partial or total denial of claim and/or improper tax
                    application.
                  </p>
                  <p>
                    I understand the penalty for willfully making a false statement of claim is a maximum fine of
                    $10,000, maximum imprisonment of five years, or both (U.S.C, Title 18, Section 287).
                  </p>
                  <Fieldset className={styles.signatureFieldSet}>
                    <Grid row gap>
                      <Grid tablet={{ col: 'fill' }}>
                        <FormGroup error={showSignatureError}>
                          <Label htmlFor="signature">Signature</Label>
                          {showSignatureError && (
                            <ErrorMessage id="signature-error-message">{errors.signature}</ErrorMessage>
                          )}
                          <Field
                            as={TextInput}
                            name="signature"
                            id="signature"
                            required
                            aria-describedby={showSignatureError ? 'signature-error-message' : null}
                          />
                        </FormGroup>
                      </Grid>
                      <Grid tablet={{ col: 'auto' }}>
                        <FormGroup>
                          <Label htmlFor="date">Date</Label>
                          <Field as={TextInput} name="date" id="date" disabled />
                        </FormGroup>
                      </Grid>
                    </Grid>
                  </Fieldset>
                  {error && (
                    <Alert type="error" headingLevel="h4" heading="Server Error">
                      There was a problem saving your signature.
                    </Alert>
                  )}
                </SectionWrapper>

                <div className={ppmStyles.buttonContainer}>
                  <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                    Return To Homepage
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
              </Form>
            </div>
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
  initialValues: shape({
    W2AddressShape,
    signature: string,
    date: string,
  }).isRequired,
};

FinalCloseoutForm.defaultProps = {
  error: false,
};

export default FinalCloseoutForm;
