import React from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, Form, FormGroup, Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import * as Yup from 'yup';

import styles from './AboutForm.module.scss';

import ppmStyles from 'components/Shared/PPM/PPM.module.scss';
import closingPageStyles from 'components/Shared/PPM/Closeout/Closeout.module.scss';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { DatePickerInput } from 'components/form/fields';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { ShipmentShape } from 'types/shipment';
import { formatCentsTruncateWhole } from 'utils/formatters';
import { requiredW2AddressSchema, requiredAddressSchema } from 'utils/validation';
import { AddressFields } from 'components/form/AddressFields/AddressFields';
import { OptionalAddressSchema } from 'components/Shared/MtoShipmentForm/validationSchemas';
import { APP_NAME } from 'constants/apps';
import { PPM_TYPES } from 'shared/constants';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const AboutForm = ({ mtoShipment, onBack, onSubmit, isSubmitted, appName }) => {
  const isCustomerPage = appName === APP_NAME.MYMOVE;
  const formFieldsName = 'w2Address';
  const today = new Date();

  const validationSchema = Yup.object().shape({
    actualMoveDate: Yup.date()
      .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
      .required('Required'),
    pickupAddress: requiredAddressSchema,
    destinationAddress: requiredAddressSchema,
    secondaryPickupAddress: OptionalAddressSchema,
    secondaryDestinationAddress: OptionalAddressSchema,
    hasReceivedAdvance: Yup.boolean().required('Required'),
    advanceAmountReceived: Yup.number().when('hasReceivedAdvance', {
      is: true,
      then: (schema) =>
        schema
          .required('Required')
          .min(1, "The minimum advance request is $1. If you don't want an advance, select No."),
    }),
    w2Address: requiredW2AddressSchema.required(),
  });

  const ppmShipment = mtoShipment?.ppmShipment || {};
  const {
    ppmType,
    pickupAddress,
    secondaryPickupAddress,
    destinationAddress,
    secondaryDestinationAddress,
    hasSecondaryPickupAddress,
    hasSecondaryDestinationAddress,
    actualMoveDate,
    hasReceivedAdvance,
    advanceAmountReceived,
  } = ppmShipment;

  const initialValues = {
    actualMoveDate: actualMoveDate || '',
    pickupAddress,
    secondaryPickupAddress: hasSecondaryPickupAddress ? secondaryPickupAddress : {},
    destinationAddress,
    secondaryDestinationAddress: hasSecondaryDestinationAddress ? secondaryDestinationAddress : {},
    hasSecondaryPickupAddress: 'false',
    hasSecondaryDestinationAddress: 'false',
    hasReceivedAdvance: hasReceivedAdvance ? 'true' : 'false',
    advanceAmountReceived: hasReceivedAdvance ? formatCentsTruncateWhole(advanceAmountReceived) : '',
    [formFieldsName]: {
      streetAddress1: mtoShipment?.ppmShipment?.w2Address?.streetAddress1 || '',
      streetAddress2: mtoShipment?.ppmShipment?.w2Address?.streetAddress2 || '',
      streetAddress3: mtoShipment?.ppmShipment?.w2Address?.streetAddress3 || '',
      city: mtoShipment?.ppmShipment?.w2Address?.city || '',
      state: mtoShipment?.ppmShipment?.w2Address?.state || '',
      postalCode: mtoShipment?.ppmShipment?.w2Address?.postalCode || '',
      county: mtoShipment?.ppmShipment?.w2Address?.county || '',
      usPostRegionCitiesID: mtoShipment?.ppmShipment?.w2Address?.usPostRegionCitiesID || '',
    },
  };

  return (
    <>
      <div className={classnames(closingPageStyles['closing-section'], closingPageStyles['about-ppm'])}>
        <p>Finish moving this PPM before you start documenting it.</p>
        <h2>How to complete your PPM</h2>
        <p>To complete your PPM, you will:</p>
        <ul>
          <li>Upload weight tickets for each trip</li>
          <li>Upload receipts to document any expenses</li>
          <li>Upload receipts if you used short-term storage, so you can request reimbursement</li>
          <li>Upload any other documentation (such as proof of ownership for a trailer, if you used your own)</li>
          <li>Complete your PPM to send it to a counselor for review</li>
        </ul>
        <h2>About your final payment</h2>
        <p>Your final payment will be:</p>
        <ul>
          <li>based on your final incentive</li>
          <li>modified by expenses submitted (authorized expenses reduce your tax burden)</li>
          <li>minus any taxes withheld (the IRS considers your incentive to be taxable income)</li>
          <li>plus any reimbursements you receive</li>
        </ul>
        <p>
          Verified expenses reduce the taxable income you report to the IRS on form W-2. They may not be claimed again
          as moving expenses. Federal tax withholding will be deducted from the profit (entitlement less eligible
          operating expenses.)
        </p>
      </div>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
          return (
            <div className={classnames(ppmStyles.formContainer, styles.AboutForm)}>
              <Form className={classnames(formStyles.form, ppmStyles.form, styles.W2Address)} data-testid="aboutForm">
                <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                  <h2>{ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped Date' : 'Departure date'}</h2>
                  {requiredAsteriskMessage}
                  <DatePickerInput
                    disabledDays={{ after: today }}
                    className={classnames(styles.actualMoveDate, 'usa-input')}
                    name="actualMoveDate"
                    label={
                      ppmType === PPM_TYPES.SMALL_PACKAGE
                        ? 'When did you ship your package?'
                        : 'When did you leave your origin?'
                    }
                    showRequiredAsterisk
                    required
                  />
                  <Hint className={ppmStyles.hint}>
                    {ppmType === PPM_TYPES.SMALL_PACKAGE
                      ? 'If you shipped multiple packages, use the first day.'
                      : 'If it took you more than one day to move out, use the first day.'}
                  </Hint>
                  <h2>Locations</h2>
                  {ppmType !== PPM_TYPES.SMALL_PACKAGE && (
                    <p>
                      If you picked things up or dropped things off from other places a long way from your start or end
                      ZIPs, ask your counselor if you should add another PPM shipment.
                    </p>
                  )}
                  <AddressFields
                    name="pickupAddress"
                    legend={ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Shipped from Address' : 'Pickup Address'}
                    labelHint="Required"
                    formikProps={formikProps}
                    className={styles.AddressFieldSet}
                    render={(fields) => (
                      <>
                        {fields}
                        <h4>Second Pickup Address</h4>
                        <FormGroup>
                          <p>
                            Will you pick up any belongings from a second address? (Must be near the pickup address.
                            Subject to approval.)
                          </p>
                          <div className={formStyles.radioGroup}>
                            <Field
                              as={Radio}
                              id="has-secondary-pickup"
                              data-testid="has-secondary-pickup"
                              label="Yes"
                              name="hasSecondaryPickupAddress"
                              value="true"
                              title="Yes, there is a second pickup address"
                              checked={values.hasSecondaryPickupAddress === 'true'}
                            />
                            <Field
                              as={Radio}
                              id="no-secondary-pickup"
                              data-testid="no-secondary-pickup"
                              label="No"
                              name="hasSecondaryPickupAddress"
                              value="false"
                              title="No, there is not a second pickup address"
                              checked={values.hasSecondaryPickupAddress !== 'true'}
                            />
                          </div>
                        </FormGroup>
                        {values.hasSecondaryPickupAddress === 'true' && (
                          <AddressFields name="secondaryPickupAddress" labelHint="Required" formikProps={formikProps} />
                        )}
                      </>
                    )}
                  />
                  <AddressFields
                    name="destinationAddress"
                    legend={ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination Address' : 'Delivery Address'}
                    className={styles.AddressFieldSet}
                    labelHint="Required"
                    formikProps={formikProps}
                    render={(fields) => (
                      <>
                        {fields}
                        <h4>Second {ppmType === PPM_TYPES.SMALL_PACKAGE ? 'Destination' : 'Delivery'} Address</h4>
                        <FormGroup>
                          <p>
                            Will you deliver any belongings to a second address? (Must be near the delivery address.
                            Subject to approval.)
                          </p>
                          <div className={formStyles.radioGroup}>
                            <Field
                              as={Radio}
                              data-testid="has-secondary-destination"
                              id="has-secondary-destination"
                              label="Yes"
                              name="hasSecondaryDestinationAddress"
                              value="true"
                              title="Yes, there is a second delivery address"
                              checked={values.hasSecondaryDestinationAddress === 'true'}
                            />
                            <Field
                              as={Radio}
                              data-testid="no-secondary-destination"
                              id="no-secondary-destination"
                              label="No"
                              name="hasSecondaryDestinationAddress"
                              value="false"
                              title="No, there is not a second delivery address"
                              checked={values.hasSecondaryDestinationAddress !== 'true'}
                            />
                          </div>
                        </FormGroup>
                        {values.hasSecondaryDestinationAddress === 'true' && (
                          <AddressFields
                            name="secondaryDestinationAddress"
                            labelHint="Required"
                            formikProps={formikProps}
                          />
                        )}
                      </>
                    )}
                  />
                  <h2>Advance (AOA)</h2>
                  <FormGroup>
                    <Fieldset className={styles.advanceFieldset}>
                      <legend className="usa-label" aria-label="Did you receive an advance for this PPM?">
                        <span>Did you receive an advance for this PPM?</span>
                      </legend>
                      <Field
                        as={Radio}
                        id="yes-has-received-advance"
                        data-testid="yes-has-received-advance"
                        label="Yes"
                        name="hasReceivedAdvance"
                        value="true"
                        checked={values.hasReceivedAdvance === 'true'}
                      />
                      <Field
                        as={Radio}
                        id="no-has-received-advance"
                        data-testid="no-has-received-advance"
                        label="No"
                        name="hasReceivedAdvance"
                        value="false"
                        checked={values.hasReceivedAdvance === 'false'}
                      />
                      <Hint className={ppmStyles.hint}>
                        If you requested an advance but did not receive it, select No.
                      </Hint>
                      {values.hasReceivedAdvance === 'true' && requiredAsteriskMessage}
                      {values.hasReceivedAdvance === 'true' && (
                        <MaskedTextField
                          label="How much did you receive?"
                          name="advanceAmountReceived"
                          id="advanceAmountReceived"
                          showRequiredAsterisk
                          required
                          defaultValue="0"
                          mask={Number}
                          scale={0} // digits after point, 0 for integers
                          signed={false} // disallow negative
                          thousandsSeparator=","
                          lazy={false} // immediate masking evaluation
                          prefix="$"
                          hintClassName={ppmStyles.innerHint}
                        />
                      )}
                    </Fieldset>
                  </FormGroup>
                  <h2>W-2 address</h2>
                  <p>What is the address on your W-2?</p>
                  <AddressFields
                    name={formFieldsName}
                    className={styles.AddressFieldSet}
                    labelHint="Required"
                    formikProps={formikProps}
                    includePOBoxes
                  />
                </SectionWrapper>
                <div
                  className={`${
                    isCustomerPage ? ppmStyles.buttonContainer : `${formStyles.formActions} ${ppmStyles.buttonGroup}`
                  }`}
                >
                  <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                    {`${isCustomerPage ? 'Return To Homepage' : 'Cancel'}`}
                  </Button>
                  <Button
                    className={ppmStyles.saveButton}
                    type="button"
                    onClick={handleSubmit}
                    disabled={!isValid || isSubmitted || isSubmitting}
                  >
                    Save & Continue
                  </Button>
                </div>
              </Form>
            </div>
          );
        }}
      </Formik>
    </>
  );
};

AboutForm.propTypes = {
  mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default AboutForm;
