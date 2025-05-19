import React, { createRef } from 'react';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { func, number } from 'prop-types';
import { Button, Form, FormGroup, Link } from '@trussworks/react-uswds';
import classnames from 'classnames';

import closingPageStyles from 'components/Shared/PPM/Closeout/Closeout.module.scss';
import styles from 'components/Shared/PPM/Closeout/GunSafeForm/GunSafeForm.module.scss';
import Fieldset from 'shared/Fieldset';
import { GunSafeTicketShape } from 'types/shipment';
import { CheckboxField } from 'components/form/fields/CheckboxField';
import WeightTicketUpload from 'components/Shared/PPM/Closeout/WeightTicketUpload/WeightTicketUpload';
import Hint from 'components/Hint';
import TextField from 'components/form/fields/TextField/TextField';
import formStyles from 'styles/form.module.scss';
import ppmStyles from 'components/Shared/PPM/PPM.module.scss';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { uploadShape } from 'types/uploads';

const GunsafeForm = ({
  gunsafe,
  setNumber,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
  onBack,
  onSubmit,
  isSubmitted,
}) => {
  const { document, weight, description, hasWeightTickets } = gunsafe || {};
  const maxWeight = 500;

  const validationSchema = Yup.object().shape({
    belongsToSelf: Yup.bool().required('Required'),
    document: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
    weight: Yup.number()
      .required('Required')
      .min(1, 'Enter a weight greater than 0 lbs.')
      .when('belongsToSelf', (schema) => {
        return schema.max(maxWeight, `Pro gear weight must be less than or equal to ${maxWeight}.`);
      }),
    description: Yup.string().required('Required'),
    missingWeightTicket: Yup.string().required(),
  });

  const initialValues = {
    document: document?.uploads || [],
    weight: weight ? `${weight}` : '',
    description: description ? `${description}` : '',
    missingWeightTicket: hasWeightTickets === false,
  };

  const documentRef = createRef();

  const jtr = (
    <Link href="https://www.defensetravel.dod.mil/Docs/perdiem/JTR.pdf" target="_blank" rel="noopener">
      Joint Travel Regulations (JTR)
    </Link>
  );

  return (
    <>
      <div className={closingPageStyles['closing-section']}>
        <p>
          If you moved pro-gear for yourself or your spouse as part of this PPM, document the total weight here.
          Reminder: This pro-gear should be included in your total weight moved.
        </p>
      </div>
      <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
        {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
          return (
            <div className={classnames(ppmStyles.formContainer, styles.GunsafeForm)}>
              <Form className={classnames(ppmStyles.form, styles.form)}>
                <SectionWrapper className={formStyles.formSection}>
                  <h2>Set {setNumber}</h2>
                  <FormGroup error={formikProps.touched?.belongsToSelf && formikProps.errors?.belongsToSelf}>
                    {
                      <Fieldset>
                        <h3>Description</h3>
                        <TextField
                          className={styles.descriptionTextField}
                          label="Brief description of the pro-gear"
                          id="description"
                          name="description"
                        />
                        <Hint className={styles.hint}>
                          Examples of pro-gear include specialized apparel and government&ndash;issued equipment.
                          <br />
                          Check the {jtr} for examples of pro-gear.
                        </Hint>
                        <h3>Weight</h3>
                        <MaskedTextField
                          containerClassName={styles.weightField}
                          defaultValue="0"
                          name="weight"
                          label="Shipment's pro-gear weight"
                          labelHint={<Hint className={styles.hint}>Your maximum allowance is {maxWeight} lbs.</Hint>}
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
                    }
                  </FormGroup>
                </SectionWrapper>
                <div className={`${`${formStyles.formActions} ${ppmStyles.buttonGroup}`}`}>
                  <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                    Cancel
                  </Button>
                  <Button
                    className={ppmStyles.saveButton}
                    type="button"
                    onClick={handleSubmit}
                    disabled={!isValid || isSubmitting || isSubmitted}
                  >
                    Save &amp; Continue
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

GunsafeForm.propTypes = {
  setNumber: number,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  gunsafe: GunSafeTicketShape,
};

GunsafeForm.defaultProps = {
  setNumber: 1,
  gunsafe: {},
};

export default GunsafeForm;
