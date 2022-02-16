import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import ppmBookingStyles from '../PPMBooking.module.scss';

import styles from './EstimatedWeightsProGearForm.module.scss';

import { MtoShipmentShape } from 'types/customerShapes';
import formStyles from 'styles/form.module.scss';
import { EntitlementShape } from 'types/order';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Hint from 'components/Hint/index';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';
import { formatWeight } from 'utils/formatters';

const validationSchema = Yup.object().shape({
  estimatedPPMWeight: Yup.number().required('Required'),
  hasProGear: Yup.boolean().required('Required'),
  estimatedProGearWeight: Yup.number().when(['hasProGear', 'estimatedSpouseProGearWeight'], {
    is: (hasProGear, estimatedSpouseProGearWeight) => hasProGear && !estimatedSpouseProGearWeight,
    then: (schema) =>
      schema
        .required(`Enter a weight into at least one pro-gear field. If you won't have pro-gear, select No above.`)
        .max(2000, 'Enter a weight less than 2,000 lbs'),
    otherwise: Yup.number().max(2000, 'Enter a weight less than 2,000 lbs'),
  }),
  estimatedSpouseProGearWeight: Yup.number().max(500, 'Enter a weight less than 500 lbs'),
});

const EstimatedWeightsProGearForm = ({ entitlement, mtoShipment, onSubmit, onBack }) => {
  const initialValues = {
    estimatedPPMWeight: mtoShipment?.ppmShipment?.estimatedWeight || '',
    hasProGear: mtoShipment?.ppmShipment?.hasProGear || 'false',
    estimatedProGearWeight: mtoShipment?.ppmShipment?.estimatedProGearWeight || '',
    estimatedSpouseProGearWeight: mtoShipment?.ppmShipment?.estimatedSpouseProGearWeight || '',
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={classnames(styles.EstimatedWeightsProGearForm, ppmBookingStyles.formContainer)}>
            <Form className={(formStyles.form, ppmBookingStyles.form)}>
              <Alert type="info">{`Total weight allowance for your move: ${formatWeight(
                entitlement?.authorizedWeight,
              )}`}</Alert>
              <SectionWrapper className={classnames(ppmBookingStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Full PPM</h2>
                <p>
                  Estimate the full weight of your PPM, including everything you plan to move. If you’re moving pro-gear
                  in this PPM, include that weight in this estimate. Do not count weight twice, though. Do not include
                  weight in your estimate that will be moved in other shipments.{' '}
                </p>
                <MaskedTextField
                  defaultValue="0"
                  name="estimatedPPMWeight"
                  label="Estimated weight of this PPM shipment"
                  id="estimatedPPMWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                  // formatWeight will display 0lbs if the weight is undefined.
                  // therefore if entitlement is undefined display the overweight warning
                  warning={
                    values.estimatedPPMWeight > entitlement?.authorizedWeight || !entitlement
                      ? 'This weight is more than your weight allowance. Talk to your counselor about what that could mean for your move.'
                      : ''
                  }
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
                  If you own a lot of things for your space, estimate on the higher side. If you own less, estimate
                  lower. The services have an official{' '}
                  <a href="https://www.ustranscom.mil/dp3/weightestimator.cfm" target="_blank" rel="noreferrer">
                    weight estimation calculator
                  </a>{' '}
                  <FontAwesomeIcon icon="external-link-alt" /> you can use for a more accurate estimate. (Link opens a
                  new window.)
                </p>
              </SectionWrapper>
              <SectionWrapper className={classnames(ppmBookingStyles.sectionWrapper, formStyles.formSection)}>
                <h2>Pro-gear</h2>
                <p>
                  Pro-gear, or PBP&E, includes books, papers, and equipment you need for your official duties. Service
                  members can move up to 2,000 lbs of pro-gear. Spouses can move up to 500 lbs.
                </p>

                <p>You get paid for moving pro-gear, but it does not count against your total weight allowance.</p>
                <Fieldset>
                  <legend className="usa-label">
                    Do you or your spouse have pro-gear that you&apos;ll move in this PPM?
                  </legend>
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
                <Hint className={ppmBookingStyles.hint}>
                  If you&apos;re not sure, select yes and your counselor can help you figure it out.
                </Hint>
                {values.hasProGear === 'true' && (
                  <>
                    <MaskedTextField
                      defaultValue="0"
                      name="estimatedProGearWeight"
                      label="Estimated weight of your pro-gear"
                      id="estimatedProGearWeight"
                      mask={Number}
                      scale={0} // digits after point, 0 for integers
                      signed={false} // disallow negative
                      thousandsSeparator=","
                      lazy={false} // immediate masking evaluation
                      suffix="lbs"
                    />
                    <MaskedTextField
                      defaultValue="0"
                      name="estimatedSpouseProGearWeight"
                      label="Estimated weight of your spouse’s pro-gear"
                      id="estimatedSpouseProGearWeight"
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
              <div className={ppmBookingStyles.buttonContainer}>
                <Button className={ppmBookingStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
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

EstimatedWeightsProGearForm.propTypes = {
  entitlement: EntitlementShape.isRequired,
  mtoShipment: MtoShipmentShape,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

EstimatedWeightsProGearForm.defaultProps = {
  mtoShipment: undefined,
};

export default EstimatedWeightsProGearForm;
