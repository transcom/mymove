import React, { useEffect, useState } from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, Alert } from '@trussworks/react-uswds';
import classnames from 'classnames';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from 'components/Customer/PPM/Booking/EstimatedWeightsProGearForm/EstimatedWeightsProGearForm.module.scss';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Hint from 'components/Hint';
import Fieldset from 'shared/Fieldset';
import formStyles from 'styles/form.module.scss';
import { OrdersShape } from 'types/customerShapes';
import { ShipmentShape } from 'types/shipment';
import { formatWeight } from 'utils/formatters';
import LoadingButton from 'components/LoadingButton/LoadingButton';
import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

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
  hasGunSafe: Yup.boolean().required('Required'),
  gunSafeWeight: Yup.number().when('hasGunSafe', {
    is: true,
    then: (schema) =>
      schema.min(1, 'Enter a weight 1 lb or greater').max(500, 'Enter a weight 500 lbs or less').required('Required'),
  }),
});

const EstimatedWeightsProGearForm = ({ orders, mtoShipment, onSubmit, onBack }) => {
  const initialValues = {
    estimatedWeight: mtoShipment?.ppmShipment?.estimatedWeight?.toString() || '',
    hasProGear: mtoShipment?.ppmShipment?.hasProGear?.toString() || 'false',
    proGearWeight: mtoShipment?.ppmShipment?.proGearWeight?.toString() || '',
    spouseProGearWeight: mtoShipment?.ppmShipment?.spouseProGearWeight?.toString() || '',
    hasGunSafe: mtoShipment?.ppmShipment?.hasGunSafe?.toString() || 'false',
    gunSafeWeight: mtoShipment?.ppmShipment?.gunSafeWeight?.toString() || '',
  };

  const [isGunSafeEnabled, setIsGunSafeEnabled] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      setIsGunSafeEnabled(await isBooleanFlagEnabled(FEATURE_FLAG_KEYS.GUN_SAFE));
    };
    fetchData();
  }, []);

  const weightAuthorized = orders.authorizedWeight;

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={classnames(styles.EstimatedWeightsProGearForm, ppmStyles.formContainer)}>
            <Form className={(formStyles.form, ppmStyles.form)}>
              <Alert headingLevel="h4" type="info">{`Remember: Your standard weight allowance is: ${formatWeight(
                weightAuthorized,
              )}. If you are moving to an administratively restricted HHG weight location this amount may be less. You will not be reimbursed for any excess weight you move.`}</Alert>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>PPM</h2>
                <p>
                  Estimate the full weight of your PPM, including everything you plan to move. If you are moving
                  pro-gear {isGunSafeEnabled && 'and/or a gun safe'} in this PPM, include that weight in this estimate.
                </p>
                <p className={formStyles.pBeforeFormGroup}>
                  Do not count weight twice, though. Do not include weight in your estimate that will be moved in other
                  shipments.
                </p>
                {requiredAsteriskMessage}
                <MaskedTextField
                  defaultValue="0"
                  name="estimatedWeight"
                  label="Estimated weight of this PPM shipment"
                  showRequiredAsterisk
                  required
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
                {requiredAsteriskMessage}
                <Fieldset>
                  <legend
                    className="usa-label"
                    aria-label="Required: Do you or your spouse have pro-gear that you'll move in this PPM?"
                  >
                    <span required>
                      Do you or your spouse have pro-gear that you&apos;ll move in this PPM? <RequiredAsterisk />
                    </span>
                  </legend>
                  <Field
                    as={Radio}
                    id="hasProGearYes"
                    data-testid="hasProGearYes"
                    label="Yes"
                    name="hasProGear"
                    value="true"
                    checked={values.hasProGear === 'true'}
                  />
                  <Field
                    as={Radio}
                    id="hasProGearNo"
                    data-testid="hasProGearNo"
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
                      required
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
                      required
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
              {isGunSafeEnabled && (
                <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                  <h2>Gun safe</h2>
                  {requiredAsteriskMessage}
                  <Fieldset>
                    <legend
                      className="usa-label"
                      aria-label="Required: Do you have a gun safe that you'll move in this PPM?"
                    >
                      <span required>
                        Do you have a gun safe that you&apos;ll move in this PPM? <RequiredAsterisk />
                      </span>
                    </legend>
                    <Field
                      as={Radio}
                      id="hasGunSafeYes"
                      data-testid="hasGunSafeYes"
                      label="Yes"
                      name="hasGunSafe"
                      value="true"
                      checked={values.hasGunSafe === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="hasGunSafeNo"
                      data-testid="hasGunSafeNo"
                      label="No"
                      name="hasGunSafe"
                      value="false"
                      checked={values.hasGunSafe === 'false'}
                    />
                  </Fieldset>
                  {values.hasGunSafe === 'true' && (
                    <div>
                      <MaskedTextField
                        defaultValue="0"
                        name="gunSafeWeight"
                        label="Estimated weight of your gun safe"
                        showRequiredAsterisk
                        required
                        id="gunSafeWeight"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        suffix="lbs"
                      />
                      <Hint>
                        The government authorizes the shipment of a gun safe up to 500 lbs. This is not charged against
                        the authorized weight entitlement. The weight entitlement is charged for any weight over 500
                        lbs. The gun safe weight cannot be added to overall entitlement for O-6 and higher ranks.
                      </Hint>
                    </div>
                  )}
                </SectionWrapper>
              )}
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Back
                </Button>
                <LoadingButton
                  buttonClassName={ppmStyles.saveButton}
                  type="button"
                  onClick={handleSubmit}
                  disabled={isSubmitting || !isValid}
                  isLoading={isSubmitting}
                  labelText="Save & Continue"
                  loadingText="Saving"
                />
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
