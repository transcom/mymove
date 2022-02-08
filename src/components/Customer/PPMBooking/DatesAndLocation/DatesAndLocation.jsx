import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './DatesAndLocation.module.scss';

import formStyles from 'styles/form.module.scss';
import { MtoShipmentShape, ServiceMemberShape } from 'types/customerShapes';
import { ZIP_CODE_REGEX } from 'utils/validation';
import TextField from 'components/form/fields/TextField/TextField';
import { CheckboxField, DatePickerInput } from 'components/form/fields';
import Hint from 'components/Hint/index';
import { DutyStationShape } from 'types';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';

const validationSchema = Yup.object().shape({
  pickupPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  useResidentialAddressZIP: Yup.boolean(),
  hasSecondaryPickupPostalCode: Yup.boolean().required('Required'),
  secondaryPickupPostalCode: Yup.string().when('hasSecondaryPickupPostalCode', {
    is: true,
    then: (schema) => schema.matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  }),
  useDestinationDutyLocationZIP: Yup.boolean(),
  destinationPostalCode: Yup.string().matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  hasSecondaryDestinationPostalCode: Yup.boolean().required('Required'),
  secondaryDestinationPostalCode: Yup.string().when('hasSecondaryDestinationPostalCode', {
    is: true,
    then: (schema) => schema.matches(ZIP_CODE_REGEX, 'Must be valid code').required('Required'),
  }),
  sitExpected: Yup.boolean().required('Required'),
  expectedDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
});

const setZip = (setFieldValue, postalCodeField, postalCode, isChecked, isCheckedField) => {
  setFieldValue(isCheckedField, isChecked === 'true' ? 'false' : 'true');
  setFieldValue(postalCodeField, isChecked === 'true' ? '' : postalCode);
};

const DatesAndLocation = ({
  mtoShipment,
  destinationDutyStation,
  serviceMember,
  onBack,
  onSubmit,
  postalCodeValidator,
}) => {
  const initialValues = {
    pickupPostalCode: mtoShipment?.ppmShipment?.pickupPostalCode || '',
    useResidentialAddressZIP: '',
    hasSecondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode ? 'true' : 'false',
    secondaryPickupPostalCode: mtoShipment?.ppmShipment?.secondaryPickupPostalCode || '',
    useDestinationDutyLocationZIP: '',
    destinationPostalCode: mtoShipment?.ppmShipment?.destinationPostalCode || '',
    hasSecondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode ? 'true' : 'false',
    secondaryDestinationPostalCode: mtoShipment?.ppmShipment?.secondaryDestinationPostalCode || '',
    sitExpected: mtoShipment?.ppmShipment?.sitExpected || 'false',
    expectedDepartureDate: mtoShipment?.ppmShipment?.expectedDepartureDate || '',
  };

  // TODO: async validation call to validate postal codes are valid for rate engine

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, setFieldValue, values }) => {
        return (
          <div className={styles.DatesAndLocationForm}>
            <Form className={(formStyles.form, styles.form)}>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
                <h2>Origin</h2>
                <TextField
                  label="ZIP"
                  id="pickupPostalCode"
                  name="pickupPostalCode"
                  maxLength={10}
                  validate={(value) => postalCodeValidator(value, 'origin')}
                />
                <CheckboxField
                  id="useResidentialAddressZIP"
                  name="useResidentialAddressZIP"
                  label={`Use my current ZIP (${serviceMember?.residentialAddress?.postalCode})`}
                  onChange={() =>
                    setZip(
                      setFieldValue,
                      'pickupPostalCode',
                      serviceMember?.residentialAddress?.postalCode,
                      values.useResidentialAddressZIP,
                      'useResidentialAddressZIP',
                    )
                  }
                  checked={values.useResidentialAddressZIP === 'true'}
                />
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">
                      Will you add items to your PPM from a place in a different ZIP code?
                    </legend>
                    <Field
                      as={Radio}
                      data-testid="yes-secondary-pickup-postal-code"
                      id="yes-secondary-pickup-postal-code"
                      label="Yes"
                      name="hasSecondaryPickupPostalCode"
                      value="true"
                      checked={values.hasSecondaryPickupPostalCode === 'true'}
                    />
                    <Field
                      as={Radio}
                      data-testid="no-secondary-pickup-postal-code"
                      id="no-secondary-pickup-postal-code"
                      label="No"
                      name="hasSecondaryPickupPostalCode"
                      value="false"
                      checked={values.hasSecondaryPickupPostalCode === 'false'}
                    />
                  </Fieldset>
                </FormGroup>
                {values.hasSecondaryPickupPostalCode === 'true' && (
                  <>
                    <TextField
                      label="Second ZIP"
                      id="secondaryPickupPostalCode"
                      name="secondaryPickupPostalCode"
                      maxLength={10}
                      validate={(value) => postalCodeValidator(value, 'origin')}
                    />
                    <Hint className={styles.hint}>
                      <p>A second origin ZIP could mean that your final incentive is lower than your estimate.</p>
                      <p>
                        Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to
                        your move counselor for more detailed information.
                      </p>
                    </Hint>
                  </>
                )}
              </SectionWrapper>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
                <h2>Destination</h2>
                <TextField
                  label="ZIP"
                  id="destinationPostalCode"
                  name="destinationPostalCode"
                  maxLength={10}
                  validate={(value) => postalCodeValidator(value, 'destination')}
                />
                <CheckboxField
                  id="useDestinationDutyLocationZIP"
                  name="useDestinationDutyLocationZIP"
                  label={`Use the ZIP for my new duty location (${destinationDutyStation?.address?.postalCode})`}
                  onChange={() =>
                    setZip(
                      setFieldValue,
                      'destinationPostalCode',
                      destinationDutyStation?.address?.postalCode,
                      values.useDestinationDutyLocationZIP,
                      'useDestinationDutyLocationZIP',
                    )
                  }
                  checked={values.useDestinationDutyLocationZIP === 'true'}
                />
                <Hint className={styles.hint}>
                  Use the ZIP for your new address if you know it. Use the ZIP for your new duty location if you
                  don&apos;t have a new address yet.
                </Hint>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">
                      Will you deliver part of your PPM to another place in a different ZIP code?
                    </legend>
                    <Field
                      as={Radio}
                      id="hasSecondaryDestinationPostalCodeYes"
                      label="Yes"
                      name="hasSecondaryDestinationPostalCode"
                      value="true"
                      checked={values.hasSecondaryDestinationPostalCode === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="hasSecondaryDestinationPostalCodeNo"
                      label="No"
                      name="hasSecondaryDestinationPostalCode"
                      value="false"
                      checked={values.hasSecondaryDestinationPostalCode === 'false'}
                    />
                  </Fieldset>
                </FormGroup>
                {values.hasSecondaryDestinationPostalCode === 'true' && (
                  <>
                    <TextField
                      label="Second ZIP"
                      id="secondaryDestinationPostalCode"
                      name="secondaryDestinationPostalCode"
                      maxLength={10}
                      validate={(value) => postalCodeValidator(value, 'destination')}
                    />
                    <Hint className={styles.hint}>
                      <p>A second destination ZIP could mean that your final incentive is lower than your estimate.</p>
                      <p>
                        Get separate weight tickets for each leg of the trip to show how the weight changes. Talk to
                        your move counselor for more detailed information.
                      </p>
                    </Hint>
                  </>
                )}
              </SectionWrapper>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
                <h2>Storage</h2>
                <Fieldset>
                  <legend className="usa-label">Do you plan to store items from your PPM?</legend>
                  <Field
                    as={Radio}
                    id="sitExpectedYes"
                    label="Yes"
                    name="sitExpected"
                    value="true"
                    checked={values.sitExpected === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="sitExpectedNo"
                    label="No"
                    name="sitExpected"
                    value="false"
                    checked={values.sitExpected === 'false'}
                  />
                </Fieldset>
                {values.sitExpected === 'false' ? (
                  <Hint className={styles.hint}>
                    You can be reimbursed for up to 90 days of temporary storage (SIT).
                  </Hint>
                ) : (
                  <Hint>
                    <p>You can be reimbursed for up to 90 days of temporary storage (SIT).</p>
                    <p>
                      Your reimbursement amount is limited to the Government&apos;s Constructed Cost — what the
                      government would have paid to store your belongings.
                    </p>
                    <p>
                      You will need to pay for the storage yourself, then submit receipts and request reimbursement
                      after your PPM is complete.
                    </p>
                    <p>Your move counselor can give you more information about additional requirements.</p>
                  </Hint>
                )}
              </SectionWrapper>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
                <h2>Departure date</h2>
                <DatePickerInput name="expectedDepartureDate" label="When do you plan to start moving your PPM?" />
                <Hint className={styles.hint}>
                  Enter the first day you expect to move things. It&apos;s OK if the actual date is different. We will
                  ask for your actual departure date when you document and complete your PPM.
                </Hint>
              </SectionWrapper>
              <div className={styles.buttonContainer}>
                <Button className={styles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
                </Button>
                <Button
                  className={styles.saveButton}
                  type="button"
                  onClick={handleSubmit}
                  disabled={!isValid || isSubmitting}
                >
                  Save & Continue
                </Button>
              </div>
            </Form>
          </div>
        );
      }}
    </Formik>
  );
};

DatesAndLocation.propTypes = {
  mtoShipment: MtoShipmentShape,
  serviceMember: ServiceMemberShape.isRequired,
  destinationDutyStation: DutyStationShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  postalCodeValidator: func.isRequired,
};

DatesAndLocation.defaultProps = {
  mtoShipment: undefined,
};

export default DatesAndLocation;
