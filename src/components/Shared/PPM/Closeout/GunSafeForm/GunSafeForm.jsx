import React, { createRef } from 'react';
import * as Yup from 'yup';
import { Formik } from 'formik';
import { func, number } from 'prop-types';
import { Button, Form, FormGroup } from '@trussworks/react-uswds';
import classnames from 'classnames';

import styles from 'components/Shared/PPM/Closeout/ProGearForm/ProGearForm.module.scss';
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
  gunSafe,
  setNumber,
  entitlements,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
  onBack,
  onSubmit,
  isSubmitted,
}) => {
  const { document, weight, description, hasWeightTickets } = gunSafe || {};
  const maxWeight = entitlements?.gunSafeWeight ?? 500;

  const validationSchema = Yup.object().shape({
    gunSafeDocument: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
    weight: Yup.number()
      .required('Required')
      .min(1, 'Enter a weight greater than 0 lbs.')
      .max(maxWeight, `Weight must be lower than ${maxWeight} lbs.`),
    description: Yup.string().required('Required'),
    missingWeightTicket: Yup.string().required(),
  });

  const initialValues = {
    gunSafeDocument: document?.uploads || [],
    weight: weight ? `${weight}` : '',
    description: description ? `${description}` : '',
    missingWeightTicket: hasWeightTickets === false,
  };

  const documentRef = createRef();

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
        return (
          <div className={classnames(ppmStyles.formContainer, styles.GunsafeForm)}>
            <Form className={classnames(ppmStyles.form, styles.form)}>
              <SectionWrapper className={formStyles.formSection}>
                <h2>Gun Safe {setNumber}</h2>
                <FormGroup>
                  {
                    <Fieldset>
                      <h3>Description</h3>
                      <TextField
                        className={styles.descriptionTextField}
                        label="Brief description of the gun safe"
                        id="description"
                        name="description"
                      />
                      <h3>Weight</h3>
                      <MaskedTextField
                        containerClassName={styles.weightField}
                        defaultValue="0"
                        name="weight"
                        label="Shipment's gun safe weight"
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
                          fieldName="gunSafeDocument"
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
  );
};

GunsafeForm.propTypes = {
  setNumber: number,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
  gunSafe: GunSafeTicketShape,
};

GunsafeForm.defaultProps = {
  setNumber: 1,
  gunSafe: {},
};

export default GunsafeForm;
