import React, { createRef } from 'react';
import { connect } from 'react-redux';
import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import { func, number } from 'prop-types';
import { Button, Form, Link, Radio } from '@trussworks/react-uswds';
import { FormGroup } from '@material-ui/core';
import classnames from 'classnames';

import { formatWeight } from 'utils/formatters';
import { selectProGearEntitlements } from 'store/entities/selectors';
import Fieldset from 'shared/Fieldset';
import { ProGearTicketShape } from 'types/shipment';
import { CheckboxField } from 'components/form/fields/CheckboxField';
import WeightTicketUpload from 'components/Customer/PPM/Closeout/WeightTicketUpload/WeightTicketUpload';
import Hint from 'components/Hint';
import TextField from 'components/form/fields/TextField/TextField';
import styles from 'components/Customer/PPM/Closeout/ProGearForm/ProGearForm.module.scss';
import formStyles from 'styles/form.module.scss';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const documentRef = createRef();

const ProGearForm = ({
  proGear,
  setNumber,
  onSubmit,
  onBack,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
  proGearEntitlements,
}) => {
  const { belongsToSelf, document, weight, description, hasWeightTickets } = proGear || {};

  const validationSchema = Yup.object().shape({
    belongsToSelf: Yup.bool().required('Required'),
    document: Yup.object().nullable(),
    weight: Yup.number()
      .required('Required')
      .min(0, 'Enter a weight 0 lbs or greater.')
      .when('belongsToSelf', (belongsToSelfField, schema) => {
        if (belongsToSelfField === null) {
          return schema;
        }
        let maximum;
        if (belongsToSelfField) {
          maximum = proGearEntitlements.proGear;
        } else {
          maximum = proGearEntitlements.proGearSpouse;
        }
        return schema.max(maximum, `Pro gear weight must be less than or equal to ${formatWeight(maximum)}.`);
      }),
    description: Yup.string().required('Required'),
    hasWeightTickets: Yup.string(),
  });

  let proGearValue;
  if (belongsToSelf === true) {
    proGearValue = 'true';
  }
  if (belongsToSelf === false) {
    proGearValue = 'false';
  }

  const initialValues = {
    belongsToSelf: proGearValue,
    document: document?.uploads || [],
    weight: weight ? `${weight}` : '',
    description: description ? `${description}` : '',
    missingWeightTicket: hasWeightTickets === 'false',
  };

  const jtr = (
    <Link href="https://www.defensetravel.dod.mil/Docs/perdiem/JTR.pdf" target="_blank" rel="noopener">
      Joint Travel Regulations (JTR)
    </Link>
  );

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
        const getEntitlement = () => {
          if (values.belongsToSelf === null) {
            return null;
          }
          if (values.belongsToSelf === 'true') {
            return proGearEntitlements.proGear;
          }
          return proGearEntitlements.proGearSpouse;
        };
        return (
          <div className={classnames(ppmStyles.formContainer, styles.ProGearForm)}>
            <Form className={classnames(ppmStyles.form, styles.form)}>
              <SectionWrapper className={formStyles.formSection}>
                <h2>Set {setNumber}</h2>
                <FormGroup>
                  <Fieldset>
                    <label htmlFor="belongsToSelf" className={classnames('usa-label', styles.descriptionTextField)}>
                      Who does this pro-gear belong to?
                      <Hint className={styles.hint}>You have to separate yours and your spouse&apos;s pro-gear.</Hint>
                    </label>
                    <Field
                      as={Radio}
                      id="ownerOfProGearSelf"
                      label="Me"
                      name="belongsToSelf"
                      value="true"
                      checked={values.belongsToSelf === 'true'}
                      data-testid="belongsToSelf"
                    />
                    <Field
                      as={Radio}
                      id="ownerOfProGearSpouse"
                      label="My spouse"
                      name="belongsToSelf"
                      value="false"
                      checked={values.belongsToSelf === 'false'}
                      data-testid="spouseProGear"
                    />
                  </Fieldset>
                  {(values.belongsToSelf === 'true' || values.belongsToSelf === 'false') && (
                    <Fieldset>
                      <h3>Description</h3>
                      <TextField
                        className={styles.descriptionTextField}
                        label="Brief description of the pro-gear"
                        labelHint={
                          <Hint className={styles.hint}>
                            Examples of pro-gear include specialized apparel and government&ndash;issued equipment.
                            <br />
                            Check the {jtr} for examples of pro-gear.
                          </Hint>
                        }
                        id="description"
                        name="description"
                      />
                      <h3>Weight</h3>
                      <MaskedTextField
                        containerClassName={styles.weightField}
                        defaultValue="0"
                        name="weight"
                        label="Shipment's pro-gear weight"
                        labelHint={
                          <Hint className={styles.hint}>
                            Your maximum allowance is {formatWeight(getEntitlement())}.
                          </Hint>
                        }
                        id="weight"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        suffix="lbs"
                      />
                      <CheckboxField
                        id="missingWeightTicket"
                        name="missingWeightTicket"
                        label="I don't have weight tickets"
                      />
                      <div>
                        <WeightTicketUpload
                          fieldName="document"
                          missingWeightTicket={values.missingWeightTicket}
                          onCreateUpload={onCreateUpload}
                          onUploadComplete={onUploadComplete}
                          onUploadDelete={onUploadDelete}
                          fileUploadRef={documentRef}
                          values={values}
                          formikProps={formikProps}
                        />
                      </div>
                    </Fieldset>
                  )}
                </FormGroup>
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
                  Save &amp; Continue
                </Button>
              </div>
            </Form>
          </div>
        );
      }}
    </Formik>
  );
};

ProGearForm.propTypes = {
  setNumber: number,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  proGear: ProGearTicketShape,
};

ProGearForm.defaultProps = {
  setNumber: 1,
  proGear: {
    belongsToSelf: null,
  },
};

const mapStateToProps = (state) => {
  const proGearEntitlements = selectProGearEntitlements(state);
  return {
    proGearEntitlements,
  };
};

export default connect(mapStateToProps)(ProGearForm);
