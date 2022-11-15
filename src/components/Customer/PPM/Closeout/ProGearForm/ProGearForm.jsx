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

const proGearDocumentRef = createRef();

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
  const { selfProGear, document, proGearWeight, description, missingWeightTicket } = proGear || {};

  const validationSchema = Yup.object().shape({
    selfProGear: Yup.bool().required('Required'),
    proGearDocument: Yup.object().nullable(),
    proGearWeight: Yup.number()
      .required('Required')
      .min(0, 'Enter a weight 0 lbs or greater')
      .when('selfProGear', (selfProGearField, schema) => {
        if (selfProGearField === null) {
          return schema;
        }
        if (selfProGearField) {
          return schema.max(proGearEntitlements.proGear);
        }
        return schema.max(proGearEntitlements.proGearSpouse);
      }),
    description: Yup.string().required('Required'),
    missingWeightTicket: Yup.string(),
  });

  let proGearValue;
  if (selfProGear === true) {
    proGearValue = 'true';
  }
  if (selfProGear === false) {
    proGearValue = 'false';
  }

  const initialValues = {
    selfProGear: proGearValue,
    proGearDocument: document?.uploads || [],
    proGearWeight: proGearWeight ? `${proGearWeight}` : '',
    description: description ? `${description}` : '',
    missingWeightTicket: missingWeightTicket ? `${missingWeightTicket}` : '',
  };

  const jtr = (
    <Link href="https://www.defensetravel.dod.mil/Docs/perdiem/JTR.pdf" target="_blank" rel="noopener">
      Joint Travel Regulations (JTR)
    </Link>
  );

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ handleSubmit, isValid, isSubmitting, setFieldValue, values, ...formikProps }) => {
        const getEntitlement = () => {
          if (values.selfProGear === null) {
            return null;
          }
          if (values.selfProGear === 'true') {
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
                    <label htmlFor="selfProGear" className={classnames('usa-label', styles.descriptionTextField)}>
                      Who does this pro-gear belong to?
                      <Hint className={styles.hint}>You have to separate yours and your spouse&apos;s pro-gear.</Hint>
                    </label>
                    <Field
                      as={Radio}
                      id="ownerOfProGearSelf"
                      label="Me"
                      name="selfProGear"
                      value="true"
                      checked={values.selfProGear === 'true'}
                      data-testid="selfProGear"
                    />
                    <Field
                      as={Radio}
                      id="ownerOfProGearSpouse"
                      label="My spouse"
                      name="selfProGear"
                      value="false"
                      checked={values.selfProGear === 'false'}
                      data-testid="spouseProGear"
                    />
                  </Fieldset>
                  {(values.selfProGear === 'true' || values.selfProGear === 'false') && (
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
                        name="proGearWeight"
                        label="Shipment's pro-gear weight"
                        labelHint={
                          <Hint className={styles.hint}>
                            Your maximum allowance is {formatWeight(getEntitlement())}.
                          </Hint>
                        }
                        id="proGearWeight"
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
                      {values.missingWeightTicket ? (
                        <div>
                          <WeightTicketUpload
                            fieldName="missingProGearWeightDocument"
                            missingWeightTicket
                            onCreateUpload={onCreateUpload}
                            onUploadComplete={onUploadComplete}
                            onUploadDelete={onUploadDelete}
                            fileUploadRef={proGearDocumentRef}
                            values={values}
                            formikProps={formikProps}
                          />
                        </div>
                      ) : (
                        <div>
                          <WeightTicketUpload
                            fieldName="proGearDocument"
                            onCreateUpload={onCreateUpload}
                            onUploadComplete={onUploadComplete}
                            onUploadDelete={onUploadDelete}
                            fileUploadRef={proGearDocumentRef}
                            values={values}
                            formikProps={formikProps}
                          />
                        </div>
                      )}
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
    selfProGear: null,
  },
};

const mapStateToProps = (state) => {
  const proGearEntitlements = selectProGearEntitlements(state);
  return {
    proGearEntitlements,
  };
};

export default connect(mapStateToProps)(ProGearForm);
