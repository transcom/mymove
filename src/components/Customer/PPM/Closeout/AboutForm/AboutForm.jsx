import React, { useState } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, Form, FormGroup, Radio } from '@trussworks/react-uswds';
import { func } from 'prop-types';
import * as Yup from 'yup';

import ppmBookingStyles from 'components/Customer/PPM/PPMBooking.module.scss';
import formStyles from 'styles/form.module.scss';
import { DatePickerInput } from 'components/form/fields';
import Hint from 'components/Hint';
import { MtoShipmentShape } from 'types/customerShapes';
import SectionWrapper from 'components/Customer/SectionWrapper';
import TextField from 'components/form/fields/TextField/TextField';
import { InvalidZIPTypeError, UnsupportedZipCodePPMErrorMsg, ZIP5_CODE_REGEX } from 'utils/validation';
import Fieldset from 'shared/Fieldset';
import { formatCentsTruncateWhole } from 'utils/formatters';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const validationSchema = Yup.object().shape({
  actualDepartureDate: Yup.date()
    .typeError('Enter a complete date in DD MMM YYYY format (day, month, year).')
    .required('Required'),
  actualPickupPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  actualDestinationPostalCode: Yup.string().matches(ZIP5_CODE_REGEX, InvalidZIPTypeError).required('Required'),
  hasReceivedAdvance: Yup.boolean().required('Required'),
  advanceAmountReceived: Yup.number().when('hasReceivedAdvance', {
    is: true,
    then: (schema) =>
      schema.required('Required').min(1, "The minimum advance request is $1. If you don't want an advance, select No."),
  }),
});

const AboutForm = ({ mtoShipment, onBack, onSubmit, postalCodeValidator }) => {
  const [postalCodeValid, setPostalCodeValid] = useState({});

  const {
    actualDepartureDate,
    actualPickupPostalCode,
    actualDestinationPostalCode,
    hasReceivedAdvance,
    advanceAmountReceived,
  } = mtoShipment?.ppmShipment || {};

  const initialValues = {
    actualDepartureDate: actualDepartureDate || '',
    actualPickupPostalCode: actualPickupPostalCode || '',
    actualDestinationPostalCode: actualDestinationPostalCode || '',
    hasReceivedAdvance: hasReceivedAdvance ? 'true' : 'false',
    advanceAmountReceived: hasReceivedAdvance ? formatCentsTruncateWhole(advanceAmountReceived) : '',
  };

  const postalCodeValidate = async (value, location, name) => {
    if (value?.length !== 5) {
      return undefined;
    }
    // only revalidate if the value has changed, editing other fields will re-validate unchanged ones
    if (postalCodeValid[`${name}`]?.value !== value) {
      const response = await postalCodeValidator(value, location, UnsupportedZipCodePPMErrorMsg);
      setPostalCodeValid((state) => {
        return {
          ...state,
          [name]: { value, isValid: !response },
        };
      });
      return response;
    }
    return postalCodeValid[`${name}`]?.isValid ? undefined : UnsupportedZipCodePPMErrorMsg;
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={ppmBookingStyles.formContainer}>
            <Form className={classnames(formStyles.form, ppmBookingStyles.form)}>
              <SectionWrapper className={classnames(ppmBookingStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Departure date</h2>
                <DatePickerInput name="actualDepartureDate" label="When did you leave your origin?" />
                <Hint className={ppmBookingStyles.hint}>
                  If it took you more than one day to move out, use the first day.
                </Hint>
                <h2>Locations</h2>
                <p>
                  If you picked things up or dropped things off from other places a long way from your start or end
                  ZIPs, ask your counselor if you should add another PPM shipment.
                </p>
                <TextField
                  label="Starting ZIP"
                  id="actualPickupPostalCode"
                  name="actualPickupPostalCode"
                  maxLength={5}
                  validate={(value) => postalCodeValidate(value, 'origin', 'actualPickupPostalCode')}
                />
                <Hint className={ppmBookingStyles.hint}>The ZIP for the address you moved away from.</Hint>
                <TextField
                  label="Ending ZIP"
                  id="actualDestinationPostalCode"
                  name="actualDestinationPostalCode"
                  maxLength={5}
                  validate={(value) => postalCodeValidate(value, 'destination', 'actualDestinationPostalCode')}
                />
                <Hint className={ppmBookingStyles.hint}>The ZIP for your new permanent address.</Hint>
                <h2>Advance (AOA)</h2>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">Did you receive an advance for this PPM?</legend>
                    <Field
                      as={Radio}
                      id="yes-has-received-advance"
                      label="Yes"
                      name="hasReceivedAdvance"
                      value="true"
                      checked={values.hasReceivedAdvance === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="no-has-received-advance"
                      label="No"
                      name="hasReceivedAdvance"
                      value="false"
                      checked={values.hasReceivedAdvance === 'false'}
                    />
                    <Hint className={ppmBookingStyles.hint}>
                      If you requested an advance but did not receive it, select No.
                    </Hint>
                    {values.hasReceivedAdvance === 'true' && (
                      <MaskedTextField
                        label="How much did you receive?"
                        name="advanceAmountReceived"
                        id="advanceAmountReceived"
                        defaultValue="0"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        prefix="$"
                        hintClassName={ppmBookingStyles.innerHint}
                      />
                    )}
                  </Fieldset>
                </FormGroup>
              </SectionWrapper>
              <div className={ppmBookingStyles.buttonContainer}>
                <Button className={ppmBookingStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Finish later
                </Button>
                <Button
                  className={ppmBookingStyles.saveButton}
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

AboutForm.propTypes = {
  mtoShipment: MtoShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  postalCodeValidator: func.isRequired,
};

export default AboutForm;
