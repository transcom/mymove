import * as Yup from 'yup';
import React from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, Form, FormGroup, Label, Radio } from '@trussworks/react-uswds';
import { func, number } from 'prop-types';

import ppmStyles from 'components/Customer/PPM/PPM.module.scss';
import styles from 'components/Customer/PPM/Closeout/WeightTicketForm/WeightTicketForm.module.scss';
import formStyles from 'styles/form.module.scss';
import SectionWrapper from 'components/Customer/SectionWrapper';
import { CheckboxField } from 'components/form/fields';
import Hint from 'components/Hint';
import TextField from 'components/form/fields/TextField/TextField';
import Fieldset from 'shared/Fieldset';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { WeightTicketShape } from 'types/shipment';
import FileUpload from 'components/FileUpload/FileUpload';
import { formatWeight } from 'utils/formatters';

const validationSchema = Yup.object().shape({
  vehicleDescription: Yup.string().required('Required'),
  emptyWeight: Yup.number().required('Required'),
  missingEmptyWeightTicket: Yup.boolean(),
  emptyWeightTickets: Yup.string(),
  fullWeight: Yup.number().required('Required'),
  missingFullWeightTicket: Yup.boolean(),
  fullWeightTickets: Yup.string(),
  hasOwnTrailer: Yup.boolean().required('Required'),
  hasClaimedTrailer: Yup.boolean(),
  trailerOwnershipDocs: Yup.string(),
});

const WeightTicketForm = ({ weightTicket, tripNumber, onBack, onSubmit }) => {
  // const { id: mtoShipmentId } = mtoShipment;

  const {
    // id: weightTicketId,
    vehicleDescription,
    missingEmptyWeightTicket,
    emptyWeight,
    // emptyWeightTickets,
    fullWeight,
    missingFullWeightTicket,
    // fullWeightTickets,
    hasOwnTrailer,
    hasClaimedTrailer,
    // trailerOwnershipDocs,
  } = weightTicket || {};

  const initialValues = {
    vehicleDescription: vehicleDescription || '',
    missingEmptyWeightTicket,
    emptyWeight: emptyWeight ? `${emptyWeight}` : '',
    emptyWeightTickets: [],
    fullWeight: fullWeight ? `${fullWeight}` : '',
    missingFullWeightTicket,
    fullWeightTickets: [],
    hasOwnTrailer: hasOwnTrailer === true ? 'true' : 'false',
    hasClaimedTrailer: hasClaimedTrailer === true ? 'true' : 'false',
    trailerOwnershipDocs: [],
  };

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values }) => {
        return (
          <div className={classnames(ppmStyles.formContainer, styles.WeightTicketForm)}>
            <Form className={classnames(formStyles.form, ppmStyles.form)}>
              <SectionWrapper className={classnames(ppmStyles.sectionWrapper, formStyles.formSection)}>
                <h2>{`Trip ${tripNumber}`}</h2>
                <h3>Vehicle</h3>
                <TextField label="Vehicle description" name="vehicleDescription" id="vehicleDescription" />
                <Hint className={ppmStyles.hint}>Car make and model, type of truck or van, etc.</Hint>
                <h3>Empty Weight</h3>
                <MaskedTextField
                  defaultValue="0"
                  name="emptyWeight"
                  label="Empty weight"
                  id="emptyWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />
                <CheckboxField
                  id="missingEmptyWeightTicket"
                  name="missingEmptyWeightTicket"
                  label="I don't have this weight ticket"
                />
                <div>
                  <Label htmlFor="emptyWeightTickets">Upload empty weight ticket</Label>
                  <Hint>PDF, JPG, or PNG only. Maximum file size 25 MB. Each page must be clear and legible.</Hint>
                  <FileUpload
                    name="emptyWeightTickets"
                    labelIdle='Drag files here or <span class="filepond--label-action">choose from folder</span>'
                  />
                </div>
                <h3>Full Weight</h3>
                <MaskedTextField
                  defaultValue="0"
                  name="fullWeight"
                  label="Full weight"
                  id="fullWeight"
                  mask={Number}
                  scale={0} // digits after point, 0 for integers
                  signed={false} // disallow negative
                  thousandsSeparator=","
                  lazy={false} // immediate masking evaluation
                  suffix="lbs"
                />
                <CheckboxField
                  id="missingFullWeightTicket"
                  name="missingFullWeightTicket"
                  label="I don't have this weight ticket"
                />
                <div>
                  <Label htmlFor="emptyWeightTickets">Upload empty weight ticket</Label>
                  <Hint>PDF, JPG, or PNG only. Maximum file size 25 MB. Each page must be clear and legible.</Hint>
                  <FileUpload
                    name="fullWeightTickets"
                    labelIdle='Drag files here or <span class="filepond--label-action">choose from folder</span>'
                  />
                </div>
                {values.fullWeight > 0 && values.emptyWeight > 0 ? (
                  <h3>{`Trip weight: ${formatWeight(values.fullWeight - values.emptyWeight)}`}</h3>
                ) : (
                  <h3>Trip weight:</h3>
                )}
                <h3>Trailer</h3>
                <FormGroup>
                  <Fieldset>
                    <legend className="usa-label">On this trip, were you using a trailer that you own?</legend>
                    <Field
                      as={Radio}
                      id="yesHasOwnTrailer"
                      label="Yes"
                      name="hasOwnTrailer"
                      value="true"
                      checked={values.hasOwnTrailer === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="noHasOwnTrailer"
                      label="No"
                      name="hasOwnTrailer"
                      value="false"
                      checked={values.hasOwnTrailer === 'false'}
                    />
                  </Fieldset>
                  {values.hasOwnTrailer === 'true' && (
                    <Fieldset>
                      <legend className="usa-label">Does your trailer meet all of these criteria?</legend>
                      <ul>
                        <li>Single axle</li>
                        <li>No more than 12 feet long from rear to trailer hitch</li>
                        <li>No more than 8 feet wide from outside tire to outside tire</li>
                        <li>Side rails and body no higher than 28 inches (unless detachable)</li>
                        <li>Ramp or gate no higher than 4 feet (unless detachable)</li>
                        <li className="text-bold">
                          Trailer weight has not already been claimed on another trip in this move
                        </li>
                      </ul>
                      <Field
                        as={Radio}
                        id="yesHasClaimTrailer"
                        label="Yes"
                        name="hasClaimedTrailer"
                        value="true"
                        checked={values.hasClaimedTrailer === 'true'}
                      />
                      <Field
                        as={Radio}
                        id="noHasClaimTrailer"
                        label="No"
                        name="hasClaimedTrailer"
                        value="false"
                        checked={values.hasClaimedTrailer === 'false'}
                      />
                      {values.hasClaimedTrailer === 'true' ? (
                        <>
                          <p>You can claim the weight of this trailer one time during your move.</p>
                          <div>
                            <Label htmlFor="trailerOwnershipDocs">Upload proof of ownership</Label>
                            <Hint>
                              <p>Examples include a registration or bill of sale.</p>
                              <p>
                                If you don’t have that documentation, upload a signed, dated statement certifying that
                                you or your spouse own this trailer.
                              </p>
                              <p>
                                PDF, JPG, or PNG only. Maximum file size 25 MB. Each page must be clear and legible.
                              </p>
                            </Hint>
                            <FileUpload
                              name="trailerOwnershipDocs"
                              labelIdle='Drag files here or <span class="filepond--label-action">choose from folder</span>'
                            />
                          </div>
                        </>
                      ) : (
                        <p className="text-bold">Do not claim the weight of this trailer for this trip.</p>
                      )}
                    </Fieldset>
                  )}
                </FormGroup>
              </SectionWrapper>
              <div className={ppmStyles.buttonContainer}>
                <Button className={ppmStyles.backButton} type="button" onClick={onBack} secondary outline>
                  Finish Later
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

WeightTicketForm.propTypes = {
  weightTicket: WeightTicketShape,
  tripNumber: number,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

WeightTicketForm.defaultProps = {
  weightTicket: undefined,
  tripNumber: 1,
};

export default WeightTicketForm;
