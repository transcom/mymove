import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from 'components/Customer/PPM/Booking/EstimatedWeightsProGearForm/EstimatedWeightsProGearForm.module.scss';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { OrdersShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { formatWeight } from 'utils/formatters';
import RequiredTag from 'components/form/RequiredTag';

const validationSchema = Yup.object().shape({
  estimatedWeight: Yup.number().min(1, 'Enter a weight greater than 0 lbs').required('Required'),
  hasProGear: Yup.boolean().required('Required'),
  proGearWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .when(['hasProGear', 'spouseProGearWeight'], {
      is: (hasProGear, spouseProGearWeight) => hasProGear && !spouseProGearWeight,
      then: (schema) =>
        schema
          .required(`Enter a weight into at least one pro-gear field. If you won't have pro-gear, select No above.`)
          .max(2000, 'Enter a weight 2,000 lbs or less'),
      otherwise: (schema) =>
        schema.min(0, 'Enter a weight 0 lbs or greater').max(2000, 'Enter a weight 2,000 lbs or less'),
    }),
  spouseProGearWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .max(500, 'Enter a weight 500 lbs or less'),
});

const EstimatedWeightsProGearForm = ({ orders, mtoShipment, onSubmit, onBack }) => {
  const initialValues = {
    estimatedWeight: mtoShipment?.ppmShipment?.estimatedWeight?.toString() || '',
    hasProGear: mtoShipment?.ppmShipment?.hasProGear?.toString() || 'false',
    proGearWeight: mtoShipment?.ppmShipment?.proGearWeight?.toString() || '',
    spouseProGearWeight: mtoShipment?.ppmShipment?.spouseProGearWeight?.toString() || '',
  };

  const weightAuthorized = orders.authorizedWeight;

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={classnames(styles.EstimatedWeightsProGearForm, ppmStyles.formContainer)}>
            <Form className={(formStyles.form, ppmStyles.form)}>
              <Alert headingLevel="h4" type="info">{`Total weight allowance for your move: ${formatWeight(
                weightAuthorized,
              )}`}</Alert>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>PPM</h2>
                <p>
                  Estimate the full weight of your PPM, including everything you plan to move. If you are moving
                  pro-gear in this PPM, include that weight in this estimate.
                </p>
                <p className={formStyles.pBeforeFormGroup}>
                  Do not count weight twice, though. Do not include weight in your estimate that will be moved in other
                  shipments.
                </p>
                <MaskedTextField
                  defaultValue="0"
                  name="estimatedWeight"
                  label="Estimated weight of this PPM shipment"
                  labelHint="Required"
                  id="estimatedWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                  warning={
                    values.estimatedWeight > weightAuthorized
                      ? 'This weight is more than your weight allowance. Talk to your counselor about what that could mean for your move.'
                      : ''
                  }
                  hintClassName={ppmStyles.innerHint}
                />

                <p>
                  This estimate can give you an idea of what you could earn for your PPM incentive. It&apos;s OK if you
                  end up moving more or less weight than this estimate.
                </p>

                <p>
                  Your final incentive amount will be determined by your finance office, based on certified weight
                  tickets that document the actual weight you moved in your PPM.
                </p>
                <h3>Need help estimating your PPM&apos;s weight?</h3>
                <p>A good guideline: Estimate 1,000 to 1,500 lbs per room.</p>

                <p>
                  If you own a lot of personal property, estimate on the higher side. If you own less, estimate lower.
                  The Services have an official{' '}
                  <a
                    href="https://www.ustranscom.mil/dp3/weightestimator.cfm"
                    target="_blank"
                    rel="noreferrer noopener"
                  >
                    weight estimation calculator
                  </a>{' '}
                  <FontAwesomeIcon icon="external-link-alt" /> tool to discover the average weight of standard household
                  items and get a calculation to compare with your entitlement allowance. (Link opens in a new window.)
                </p>
              </SectionWrapper>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Pro-gear</h2>
                <p>
                  Pro-gear, or PBP&E, includes books, papers, and equipment you need for your official duties. Service
                  members can move up to 2,000 lbs of pro-gear. Additionally, spouses can move up to 500 lbs.
                </p>

                <p>You get paid for moving pro-gear, but it does not count against your total weight allowance.</p>
                <Fieldset>
                  <legend className="usa-label">
                    Do you or your spouse have pro-gear that you&apos;ll move in this PPM?
                  </legend>
                  <RequiredTag />
                  <Field
                    as={Radio}
                    id="hasProGearYes"
                    label="Yes"
                    name="hasProGear"
                    value="true"
                    checked={values.hasProGear === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="hasProGearNo"
                    label="No"
                    name="hasProGear"
                    value="false"
                    checked={values.hasProGear === 'false'}
                  />
                </Fieldset>
                <Hint className={ppmStyles.hint}>
                  If you are not sure, select yes and your counselor can help you figure it out.
                </Hint>
                {values.hasProGear === 'true' && (
                  <>
                    <MaskedTextField
                      defaultValue="0"
                      name="proGearWeight"
                      label="Estimated weight of your pro-gear"
                      labelHint="Required"
                      id="proGearWeight"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                    />
                    <MaskedTextField
                      defaultValue="0"
                      name="spouseProGearWeight"
                      label="Estimated weight of your spouse’s pro-gear"
                      labelHint="Required"
                      id="spouseProGearWeight"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                    />
                    <Hint>
                      Talk to your counselor about requirements for documenting pro-gear included in your PPM.
                    </Hint>
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

EstimatedWeightsProGearForm.propTypes = {
  orders: OrdersShape.isRequired,
  mtoShipment: ShipmentShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default EstimatedWeightsProGearForm;
