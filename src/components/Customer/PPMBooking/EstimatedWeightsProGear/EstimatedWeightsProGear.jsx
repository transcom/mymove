import React from 'react';
import { func } from 'prop-types';
import * as Yup from 'yup';
import { Formik, Field } from 'formik';
import { Button, Form, Radio, FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from './EstimatedWeightsProGear.module.scss';

import formStyles from 'styles/form.module.scss';
import { EntitlementShape } from 'types/order';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import Hint from 'components/Hint/index';
import SectionWrapper from 'components/Customer/SectionWrapper';
import Fieldset from 'shared/Fieldset';

const validationSchema = Yup.object().shape({
  estimatedPPMWeight: Yup.number().required('Required'),
  hasProGear: Yup.boolean().required('Required'),
  estimatedProGearWeight: Yup.number().when('hasProGear', {
    is: true,
    then: (schema) => schema.required('Required'),
  }),
  estimatedSpouseProGearWeight: Yup.number().when('hasProGear', {
    is: true,
    then: (schema) => schema.required('Required'),
  }),
});

const EstimatedWeightsProGear = ({ entitlement, onSubmit, onBack }) => {
  const initialValues = {
    estimatedPPMWeight: '',
    hasProGear: 'false',
    estimatedProGearWeight: '',
    estimatedSpouseProGearWeight: '',
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={styles.EstimatedWeightsProGearForm}>
            <Form className={(formStyles.form, styles.form)}>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
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
                  lower. The services have an official weight estimation calculator you can use for a more accurate
                  estimate. (Link opens a new window.)
                </p>
              </SectionWrapper>
              <SectionWrapper className={classnames(styles.sectionWrapper, formStyles.formSection)}>
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
                <Hint className={styles.hint}>
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
                    />
                  </>
                )}
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

EstimatedWeightsProGear.propTypes = {
  entitlement: EntitlementShape.isRequired,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

export default EstimatedWeightsProGear;
