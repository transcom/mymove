import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import { Button, Form, Radio } from '@trussworks/react-uswds';
import classnames from 'classnames';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import { MtoShipmentShape } from 'types/customerShapes';
import formStyles from 'styles/form.module.scss';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Hint from 'components/Hint';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';
import { CheckboxField } from 'components/form/fields';
import { calculateMaxAdvanceAndFormatAdvanceAndIncentive, getFormattedMaxAdvancePercentage } from 'utils/incentives';
import { formatCentsTruncateWhole } from 'utils/formatters';

const validationSchema = (maxAdvance, formattedMaxAdvance) => {
  return Yup.object().shape({
    advanceRequested: Yup.boolean().required('Required'),
    amountRequested: Yup.number().when('advanceRequested', {
      is: true,
      then: (schema) =>
        schema
          .required('Required')
          .min(1, "The minimum advance request is $1. If you don't want an advance, select No.")
          .max(maxAdvance, `Enter an amount $${formattedMaxAdvance} or less`),
    }),
    agreeToTerms: Yup.boolean().when('advanceRequested', {
      is: true,
      then: (schema) => schema.oneOf([true], 'Required'),
    }),
  });
};

const AdvanceForm = ({ mtoShipment, onSubmit, onBack }) => {
  const { advance, estimatedIncentive } = mtoShipment?.ppmShipment || {};
  const initialValues = {
    amountRequested: advance ? formatCentsTruncateWhole(advance) : '',
    advanceRequested: advance ? 'true' : 'false',
    agreeToTerms: false,
  };

  const { maxAdvance, formattedMaxAdvance, formattedIncentive } =
    calculateMaxAdvanceAndFormatAdvanceAndIncentive(estimatedIncentive);

  return (
    <Formik
      initialValues={initialValues}
      validationSchema={() => validationSchema(maxAdvance, formattedMaxAdvance)}
      onSubmit={onSubmit}
    >
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={ppmStyles.formContainer}>
            <Form className={(formStyles.form, ppmStyles.form)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>{`You can ask for up to $${formattedMaxAdvance} as an advance`}</h2>
                <p>{`That's ${getFormattedMaxAdvancePercentage()} of $${formattedIncentive}, the estimated incentive for your PPM.`}</p>
                <p>
                  You can use an advance to pay for up-front moving expenses. Your service will award and distribute any
                  advance. Different services handle advances differently. Talk to your counselor for specifics.
                </p>
                <p>If you take an advance, it will be deducted from your incentive payment.</p>
                <Fieldset>
                  <legend className="usa-label">Would you like to request an advance on your incentive?</legend>
                  <Field
                    as={Radio}
                    id="advanceRequestedYes"
                    label="Yes"
                    name="advanceRequested"
                    value="true"
                    checked={values.advanceRequested === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="advanceRequestedNo"
                    label="No"
                    name="advanceRequested"
                    value="false"
                    checked={values.advanceRequested === 'false'}
                  />
                </Fieldset>
                {values.advanceRequested === 'true' && (
                  <>
                    <MaskedTextField
                      defaultValue="0"
                      name="amountRequested"
                      label="Amount requested"
                      id="amountRequested"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      prefix="$"
                      hintClassName={ppmStyles.innerHint}
                    />
                    <Hint>
                      Your move counselor will discuss next steps with you and let you know how you&apos;ll receive your
                      advance.
                    </Hint>
                    <CheckboxField
                      id="agreeToTerms"
                      name="agreeToTerms"
                      label="I acknowledge that any advance I'm given will be deducted from my final incentive payment. If my advance ends up being more than my incentive, I will need to repay the difference."
                    />
                  </>
                )}
              </SectionWrapper>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
                </Button>
                <Button
                  className={ppmStyles.saveButton}
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

AdvanceForm.propTypes = {
  mtoShipment: MtoShipmentShape,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

AdvanceForm.defaultProps = {
  mtoShipment: undefined,
};

export default AdvanceForm;
