import * as Yup from 'yup';
import React, { createRef } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, ErrorMessage, Form, FormGroup, Label, Link, Radio } from '@trussworks/react-uswds';
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
import UploadsTable from 'components/UploadsTable/UploadsTable';
import {
  DocumentAndImageUploadInstructions,
  SpreadsheetUploadInstructions,
  UploadDropZoneLabel,
} from 'content/uploads';
import uploadShape from 'types/uploads';

const validationSchema = Yup.object().shape({
  vehicleDescription: Yup.string().required('Required'),
  emptyWeight: Yup.number().min(0, 'Enter a weight 0 lbs or greater').required('Required'),
  missingEmptyWeightTicket: Yup.boolean(),
  emptyWeightTickets: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  fullWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .required('Required')
    .when('emptyWeight', (emptyWeight, schema) => {
      return emptyWeight
        ? schema.min(emptyWeight + 1, 'The full weight must be greater than the empty weight')
        : schema;
    }),
  missingFullWeightTicket: Yup.boolean(),
  fullWeightTickets: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  hasOwnTrailer: Yup.boolean().required('Required'),
  hasClaimedTrailer: Yup.boolean(),
  trailerOwnershipDocs: Yup.array()
    .of(uploadShape)
    .when('hasClaimedTrailer', (hasClaimedTrailer, schema) => {
      return hasClaimedTrailer ? schema.min(1, 'At least one upload is required') : schema;
    }),
});

const acceptableFileTypes = [
  'image/jpeg',
  'image/png',
  'application/pdf',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'application/vnd.ms-excel',
];

const WeightTicketForm = ({
  weightTicket,
  tripNumber,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
  onBack,
  onSubmit,
}) => {
  // const { id: mtoShipmentId } = mtoShipment;

  const {
    // id: weightTicketId,
    vehicleDescription,
    missingEmptyWeightTicket,
    emptyWeight,
    emptyWeightTickets,
    fullWeight,
    missingFullWeightTicket,
    fullWeightTickets,
    hasOwnTrailer,
    hasClaimedTrailer,
    trailerOwnershipDocs,
  } = weightTicket || {};

  const initialValues = {
    vehicleDescription: vehicleDescription || '',
    missingEmptyWeightTicket,
    emptyWeight: emptyWeight ? `${emptyWeight}` : '',
    emptyWeightTickets: emptyWeightTickets || [],
    fullWeight: fullWeight ? `${fullWeight}` : '',
    missingFullWeightTicket,
    fullWeightTickets: fullWeightTickets || [],
    hasOwnTrailer: hasOwnTrailer === true ? 'true' : 'false',
    hasClaimedTrailer: hasClaimedTrailer === true ? 'true' : 'false',
    trailerOwnershipDocs: trailerOwnershipDocs || [],
  };

  const constructedWeightDownload = (
    <>
      <p>Download the official government spreadsheet to calculate constructed weight.</p>
      <Link
        className={classnames('usa-button', 'usa-button--secondary', styles.constructedWeightLink)}
        to="https://www.ustranscom.mil/dp3/weightestimator.cfm"
      >
        Go to download page
      </Link>
      <p>
        Enter the constructed weight you calculated.
        <br />
        Upload a completed copy of the spreadsheet.
      </p>
    </>
  );

  const emptyWeightTicketUploadLabel = (showConstructedWeight) => {
    return showConstructedWeight ? 'Upload constructed weight spreadsheet' : 'Upload empty weight ticket';
  };

  const fullWeightTicketUploadLabel = (showConstructedWeight) => {
    return showConstructedWeight ? 'Upload constructed weight spreadsheet' : 'Upload full weight ticket';
  };

  const weightTicketUploadHint = (showConstructedWeight) => {
    return showConstructedWeight ? SpreadsheetUploadInstructions : DocumentAndImageUploadInstructions;
  };

  const emptyWeightTicketsRef = createRef();
  const fullWeightTicketsRef = createRef();
  const trailerOwnershipDocsRef = createRef();

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values, touched, errors, setFieldTouched, setFieldValue }) => {
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
                  {values.missingEmptyWeightTicket && constructedWeightDownload}
                  <UploadsTable
                    className={styles.uploadsTable}
                    uploads={values.emptyWeightTickets}
                    onDelete={(uploadId) =>
                      onUploadDelete(uploadId, 'emptyWeightTickets', values, setFieldTouched, setFieldValue)
                    }
                  />
                  <FormGroup error={touched?.emptyWeightTickets && errors?.emptyWeightTickets}>
                    <div className="labelWrapper">
                      <Label
                        error={touched?.emptyWeightTickets && errors?.emptyWeightTickets}
                        htmlFor="emptyWeightTickets"
                      >
                        {emptyWeightTicketUploadLabel(values.missingEmptyWeightTicket)}
                      </Label>
                    </div>
                    {touched?.emptyWeightTickets && errors?.emptyWeightTickets && (
                      <ErrorMessage>{errors?.emptyWeightTickets}</ErrorMessage>
                    )}
                    <Hint className={styles.uploadTypeHint}>
                      {weightTicketUploadHint(values.missingEmptyWeightTicket)}
                    </Hint>
                    <FileUpload
                      name="emptyWeightTickets"
                      labelIdle={UploadDropZoneLabel}
                      createUpload={onCreateUpload}
                      onChange={(err, upload) => {
                        setFieldTouched('emptyWeightTickets', true);
                        onUploadComplete(upload, err, 'emptyWeightTickets', values, setFieldValue);
                        emptyWeightTicketsRef.current.removeFile(upload.id);
                      }}
                      acceptedFileTypes={acceptableFileTypes}
                      ref={emptyWeightTicketsRef}
                    />
                  </FormGroup>
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
                  {values.missingFullWeightTicket && constructedWeightDownload}
                  <UploadsTable
                    className={styles.uploadsTable}
                    uploads={values.fullWeightTickets}
                    onDelete={onUploadDelete}
                  />
                  <FormGroup error={touched?.fullWeightTickets && errors?.fullWeightTickets}>
                    <div className="labelWrapper">
                      <Label
                        error={touched?.fullWeightTickets && errors?.fullWeightTickets}
                        htmlFor="emptyWeightTickets"
                      >
                        {fullWeightTicketUploadLabel(values.missingFullWeightTicket)}
                      </Label>
                    </div>
                    {touched?.fullWeightTickets && errors?.fullWeightTickets && (
                      <ErrorMessage>{errors?.fullWeightTickets}</ErrorMessage>
                    )}
                    <Hint className={styles.uploadTypeHint}>
                      {weightTicketUploadHint(values.missingFullWeightTicket)}
                    </Hint>
                    <FileUpload
                      name="fullWeightTickets"
                      labelIdle={UploadDropZoneLabel}
                      createUpload={onCreateUpload}
                      onChange={(err, upload) => {
                        setFieldTouched('fullWeightTickets', true);
                        onUploadComplete(upload, err, 'fullWeightTickets', values, setFieldValue);
                        fullWeightTicketsRef.current.removeFile(upload.id);
                      }}
                      acceptedFileTypes={acceptableFileTypes}
                      ref={fullWeightTicketsRef}
                    />
                  </FormGroup>
                </div>
                {values.fullWeight > 0 && values.emptyWeight > 0 ? (
                  <h3>{`Trip weight: ${formatWeight(values.fullWeight - values.emptyWeight)}`}</h3>
                ) : (
                  <h3>Trip weight:</h3>
                )}
                <h3>Trailer</h3>
                <FormGroup>
                  <Fieldset className={styles.trailerOwnershipFieldset}>
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
                            <UploadsTable
                              className={styles.uploadsTable}
                              uploads={values.trailerOwnershipDocs}
                              onDelete={onUploadDelete}
                            />
                            <FormGroup error={touched?.trailerOwnershipDocs && errors?.trailerOwnershipDocs}>
                              <div className="labelWrapper">
                                <Label
                                  error={touched?.trailerOwnershipDocs && errors?.trailerOwnershipDocs}
                                  htmlFor="trailerOwnershipDocs"
                                >
                                  Upload proof of ownership
                                </Label>
                              </div>
                              {touched?.trailerOwnershipDocs && errors?.trailerOwnershipDocs && (
                                <ErrorMessage>{errors?.trailerOwnershipDocs}</ErrorMessage>
                              )}
                              <Hint>
                                <p>Examples include a registration or bill of sale.</p>
                                <p>
                                  If you don’t have that documentation, upload a signed, dated statement certifying that
                                  you or your spouse own this trailer.
                                </p>
                                <p className={styles.uploadTypeHint}>{DocumentAndImageUploadInstructions}</p>
                              </Hint>
                              <FileUpload
                                name="trailerOwnershipDocs"
                                createUpload={onCreateUpload}
                                labelIdle={UploadDropZoneLabel}
                                onChange={(err, upload) => {
                                  setFieldTouched('trailerOwnershipDocs', true);
                                  onUploadComplete(upload, err, 'trailerOwnershipDocs', values, setFieldValue);
                                  trailerOwnershipDocsRef.current.removeFile(upload.id);
                                }}
                                acceptedFileTypes={acceptableFileTypes}
                                ref={trailerOwnershipDocsRef}
                              />
                            </FormGroup>
                          </div>
                        </>
                      ) : (
                        <p className={styles.doNotClaimTrailerWeight}>
                          Do not claim the weight of this trailer for this trip.
                        </p>
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
  onCreateUpload: func.isRequired,
  onUploadComplete: func.isRequired,
  onUploadDelete: func,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

WeightTicketForm.defaultProps = {
  weightTicket: undefined,
  onUploadDelete: undefined,
  tripNumber: 1,
};

export default WeightTicketForm;
