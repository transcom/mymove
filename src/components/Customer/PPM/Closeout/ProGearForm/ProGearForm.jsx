import React, { createRef } from 'react';
import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import { func, number } from 'prop-types';
import { Button, Form, Link, Radio } from '@trussworks/react-uswds';
import { FormGroup } from '@material-ui/core';
import classnames from 'classnames';

import Fieldset from 'shared/Fieldset';
import { ProGearTicketShape } from 'types/shipment';
import { CheckboxField } from 'components/form/fields/CheckboxField';
import WeightTicketUpload from 'components/Customer/PPM/Closeout/WeightTicketUpload/WeightTicketUpload';
import Hint from 'components/Hint';
import TextField from 'components/form/fields/TextField/TextField';
import formStyles from 'styles/form.module.scss';
import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';

const validationSchema = Yup.object().shape({
  selfProGear: Yup.bool().required('Required'),
});

const proGearDocumentRef = createRef();

const ProGearForm = ({ proGear, setNumber, onSubmit, onBack, onCreateUpload, onUploadComplete, onUploadDelete }) => {
  const { selfProGear, document, proGearWeight } = proGear || {};
  const initialValues = {
    selfProGear,
    proGearDocument: document?.uploads || [],
    proGearWeight: proGearWeight ? `${proGearWeight}` : '',
  };

  const jtr = (
    <Link
      className={classnames('string')}
      href="https://www.defensetravel.dod.mil/Docs/perdiem/JTR.pdf"
      target="_blank"
      rel="noopener"
    >
      Joint Travel Regulations (JTR)
    </Link>
  );
  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ handleSubmit, isValid, isSubmitting, values, ...formikProps }) => {
        return (
          <div className={ppmStyles.formContainer}>
            <Form className={ppmStyles.form}>
              <SectionWrapper className={formStyles.formSection}>
                <h2>Set {setNumber}</h2>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label margin-bottom-0">Who does this pro-gear belongs to?</legend>
                    <Hint className={ppmStyles.hint}>You have to separate yours and your spouse&apos;s pro-gear.</Hint>
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
                  {values.selfProGear === 'true' && (
                    <Fieldset>
                      <h3>Description</h3>
                      <TextField label="Brief description of the pro-gear" id="description" name="description" />
                      <Hint className={ppmStyles.hint}>
                        Examples of pro-gear includes specialized apparel and government issued equiptment. Check the{' '}
                        {jtr} for examples of pro-gear.
                      </Hint>
                      <h3>Weight</h3>
                      <MaskedTextField
                        defaultValue="0"
                        name="proGearWeight"
                        label="Shipment's pro-gear weight"
                        id="proGearWeight"
                        mask={Number}
                        scale={0} // digits after point, 0 for integers
                        signed={false} // disallow negative
                        thousandsSeparator=","
                        lazy={false} // immediate masking evaluation
                        suffix="lbs"
                      />
                      <Hint className={ppmStyles.hint}>Your maximum allowance is X,XXX lbs.</Hint>
                      <CheckboxField
                        id="missingWeightTicket"
                        name="missingWeightTicket"
                        label="I don't have weight tickets"
                      />
                      {values.missingWeightTicket ? (
                        <>
                          <p>Download the official government spreadsheet to calculate the constructed weight.</p>
                          <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                            Go to Download Page
                          </Button>
                          <p>Enter the constructed weight you calculated.</p>
                          <p>Upload a completed copy of the spreadsheet.</p>
                          <div>
                            <WeightTicketUpload
                              fieldName="missingProGearWeightDocument"
                              onCreateUpload={onCreateUpload}
                              onUploadComplete={onUploadComplete}
                              onUploadDelete={onUploadDelete}
                              fileUploadRef={proGearDocumentRef}
                              values={values}
                              formikProps={formikProps}
                            />
                          </div>
                        </>
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
    selfProGear: true,
  },
};

export default ProGearForm;
