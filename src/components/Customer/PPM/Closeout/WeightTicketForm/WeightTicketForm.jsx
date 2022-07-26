import * as Yup from 'yup';
import React, { createRef } from 'react';
import { Field, Formik } from 'formik';
import classnames from 'classnames';
import { Button, ErrorMessage, Form, FormGroup, Label, Link, Radio } from '@trussworks/react-uswds';
import { string, bool, func, shape } from 'prop-types';

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
  UploadDropZoneLabelMobile,
} from 'content/uploads';
import { uploadShape } from 'types/uploads';

const validationSchema = Yup.object().shape({
  vehicleDescription: Yup.string().required('Required'),
  emptyWeight: Yup.number().min(0, 'Enter a weight 0 lbs or greater').required('Required'),
  missingEmptyWeightTicket: Yup.boolean(),
  emptyDocument: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  fullWeight: Yup.number()
    .min(0, 'Enter a weight 0 lbs or greater')
    .required('Required')
    .when('emptyWeight', (emptyWeight, schema) => {
      return emptyWeight != null
        ? schema.min(emptyWeight + 1, 'The full weight must be greater than the empty weight')
        : schema;
    }),
  missingFullWeightTicket: Yup.boolean(),
  fullDocument: Yup.array().of(uploadShape).min(1, 'At least one upload is required'),
  ownsTrailer: Yup.boolean().required('Required'),
  trailerMeetsCriteria: Yup.boolean(),
  proofOfTrailerOwnershipDocument: Yup.array()
    .of(uploadShape)
    .when(['ownsTrailer', 'trailerMeetsCriteria'], (ownsTrailer, trailerMeetsCriteria, schema) => {
      return ownsTrailer && trailerMeetsCriteria ? schema.min(1, 'At least one upload is required') : schema;
    }),
});

const acceptableFileTypes = [
  'image/jpeg',
  'image/png',
  'application/pdf',
  'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
  'application/vnd.ms-excel',
];

const constructedWeightDownload = (
  <>
    <p>Download the official government spreadsheet to calculate constructed weight.</p>
    <Link
      className={classnames('usa-button', 'usa-button--secondary', styles.constructedWeightLink)}
      href="https://www.ustranscom.mil/dp3/weightestimator.cfm"
      target="_blank"
      rel="noopener"
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

const WeightTicketUpload = ({
  fieldName,
  missingWeightTicket,
  onCreateUpload,
  onUploadComplete,
  onUploadDelete,
  fileUploadRef,
  values,
  formikProps: { touched, errors, setFieldTouched, setFieldValue },
}) => {
  const weightTicketUploadLabel = (name, showConstructedWeight) => {
    if (name === 'emptyDocument') {
      return showConstructedWeight ? 'Upload constructed weight spreadsheet' : 'Upload empty weight ticket';
    }

    return showConstructedWeight ? 'Upload constructed weight spreadsheet' : 'Upload full weight ticket';
  };

  const weightTicketUploadHint = (showConstructedWeight) => {
    return showConstructedWeight ? SpreadsheetUploadInstructions : DocumentAndImageUploadInstructions;
  };

  const showError = touched[`${fieldName}`] && errors[`${fieldName}`];

  return (
    <>
      {missingWeightTicket && constructedWeightDownload}
      <UploadsTable
        className={styles.uploadsTable}
        uploads={values[`${fieldName}`]}
        onDelete={(uploadId) => onUploadDelete(uploadId, fieldName, setFieldTouched, setFieldValue)}
      />
      <FormGroup error={showError}>
        <div className="labelWrapper">
          <Label error={showError} htmlFor={fieldName}>
            {weightTicketUploadLabel(fieldName, missingWeightTicket)}
          </Label>
        </div>
        {showError && <ErrorMessage>{errors[`${fieldName}`]}</ErrorMessage>}
        <Hint className={styles.uploadTypeHint}>{weightTicketUploadHint(missingWeightTicket)}</Hint>
        <FileUpload
          name={fieldName}
          labelIdle={UploadDropZoneLabel}
          labelIdleMobile={UploadDropZoneLabelMobile}
          createUpload={(file) => onCreateUpload(fieldName, file)}
          onChange={(err, upload) => {
            setFieldTouched(fieldName, true);
            onUploadComplete(err);
            fileUploadRef.current.removeFile(upload.id);
          }}
          acceptedFileTypes={acceptableFileTypes}
          labelFileTypeNotAllowed="Upload a supported file type"
          fileValidateTypeLabelExpectedTypes="Supported file types: PDF, JPG, PNG, XLS, or XLSX"
          ref={fileUploadRef}
        />
      </FormGroup>
    </>
  );
};

WeightTicketUpload.propTypes = {
  fieldName: string.isRequired,
  missingWeightTicket: bool,
  onCreateUpload: func.isRequired,
  onUploadComplete: func.isRequired,
  onUploadDelete: func.isRequired,
  fileUploadRef: shape({ current: shape({}) }).isRequired,
  values: shape({}).isRequired,
  formikProps: shape({
    touched: shape({}),
    errors: shape({}),
    setFieldTouched: func,
    setFieldValue: func,
  }).isRequired,
};

WeightTicketUpload.defaultProps = {
  missingWeightTicket: false,
};

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
    emptyDocument,
    fullWeight,
    missingFullWeightTicket,
    fullDocument,
    ownsTrailer,
    trailerMeetsCriteria,
    proofOfTrailerOwnershipDocument,
  } = weightTicket || {};

  const initialValues = {
    vehicleDescription: vehicleDescription || '',
    missingEmptyWeightTicket: !!missingEmptyWeightTicket,
    emptyWeight: emptyWeight ? `${emptyWeight}` : '',
    emptyDocument: emptyDocument?.uploads || [],
    fullWeight: fullWeight ? `${fullWeight}` : '',
    missingFullWeightTicket: !!missingFullWeightTicket,
    fullDocument: fullDocument?.uploads || [],
    ownsTrailer: ownsTrailer ? 'true' : 'false',
    trailerMeetsCriteria: trailerMeetsCriteria ? 'true' : 'false',
    proofOfTrailerOwnershipDocument: proofOfTrailerOwnershipDocument?.uploads || [],
  };

  const emptyDocumentRef = createRef();
  const fullDocumentRef = createRef();
  const proofOfTrailerOwnershipDocumentRef = createRef();

  return (
    <Formik initialValues={initialValues} validationSchema={validationSchema} onSubmit={onSubmit}>
      {({ isValid, isSubmitting, handleSubmit, values, ...formikProps }) => {
        return (
          <div className={classnames(ppmStyles.formContainer, styles.WeightTicketForm)}>
            <Form className={classnames(formStyles.form, ppmStyles.form)}>
              <SectionWrapper className={classnames(formStyles.formSection, styles.weightTicketSectionWrapper)}>
                <h2>{`Trip ${tripNumber}`}</h2>
                <h3>Vehicle</h3>
                <TextField label="Vehicle description" name="vehicleDescription" id="vehicleDescription" />
                <Hint className={ppmStyles.hint}>Car make and model, type of truck or van, etc.</Hint>
                <h3>Empty weight</h3>
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
                  <WeightTicketUpload
                    fieldName="emptyDocument"
                    missingWeightTicket={values.missingEmptyWeightTicket}
                    onCreateUpload={onCreateUpload}
                    onUploadComplete={onUploadComplete}
                    onUploadDelete={onUploadDelete}
                    fileUploadRef={emptyDocumentRef}
                    values={values}
                    formikProps={formikProps}
                  />
                </div>
                <h3>Full weight</h3>
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
                  <WeightTicketUpload
                    fieldName="fullDocument"
                    missingWeightTicket={values.missingFullWeightTicket}
                    onCreateUpload={onCreateUpload}
                    onUploadComplete={onUploadComplete}
                    onUploadDelete={onUploadDelete}
                    fileUploadRef={fullDocumentRef}
                    values={values}
                    formikProps={formikProps}
                  />
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
                      id="yesOwnsTrailer"
                      label="Yes"
                      name="ownsTrailer"
                      value="true"
                      checked={values.ownsTrailer === 'true'}
                    />
                    <Field
                      as={Radio}
                      id="noOwnsTrailer"
                      label="No"
                      name="ownsTrailer"
                      value="false"
                      checked={values.ownsTrailer === 'false'}
                    />
                  </Fieldset>
                  {values.ownsTrailer === 'true' && (
                    <Fieldset className={styles.trailerClaimedFieldset}>
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
                        id="yestrailerMeetsCriteria"
                        label="Yes"
                        name="trailerMeetsCriteria"
                        value="true"
                        checked={values.trailerMeetsCriteria === 'true'}
                      />
                      <Field
                        as={Radio}
                        id="notrailerMeetsCriteria"
                        label="No"
                        name="trailerMeetsCriteria"
                        value="false"
                        checked={values.trailerMeetsCriteria === 'false'}
                      />
                      {values.trailerMeetsCriteria === 'true' ? (
                        <>
                          <p>You can claim the weight of this trailer one time during your move.</p>
                          <div>
                            <UploadsTable
                              className={styles.uploadsTable}
                              uploads={values.proofOfTrailerOwnershipDocument}
                              onDelete={(uploadId) =>
                                onUploadDelete(
                                  uploadId,
                                  'proofOfTrailerOwnershipDocument',
                                  formikProps.setFieldTouched,
                                  formikProps.setFieldValue,
                                )
                              }
                            />
                            <FormGroup
                              error={
                                formikProps.touched?.proofOfTrailerOwnershipDocument &&
                                formikProps.errors?.proofOfTrailerOwnershipDocument
                              }
                            >
                              <div className="labelWrapper">
                                <Label
                                  error={
                                    formikProps.touched?.proofOfTrailerOwnershipDocument &&
                                    formikProps.errors?.proofOfTrailerOwnershipDocument
                                  }
                                  htmlFor="proofOfTrailerOwnershipDocument"
                                >
                                  Upload proof of ownership
                                </Label>
                              </div>
                              {formikProps.touched?.proofOfTrailerOwnershipDocument &&
                                formikProps.errors?.proofOfTrailerOwnershipDocument && (
                                  <ErrorMessage>{formikProps.errors?.proofOfTrailerOwnershipDocument}</ErrorMessage>
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
                                name="proofOfTrailerOwnershipDocument"
                                createUpload={(file) => onCreateUpload('proofOfTrailerOwnershipDocument', file)}
                                labelIdle={UploadDropZoneLabel}
                                labelIdleMobile={UploadDropZoneLabelMobile}
                                onChange={(err, upload) => {
                                  formikProps.setFieldTouched('proofOfTrailerOwnershipDocument', true);
                                  onUploadComplete(err);
                                  proofOfTrailerOwnershipDocumentRef.current.removeFile(upload.id);
                                }}
                                acceptedFileTypes={acceptableFileTypes}
                                ref={proofOfTrailerOwnershipDocumentRef}
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
  tripNumber: string,
  onCreateUpload: func.isRequired,
  onUploadComplete: func.isRequired,
  onUploadDelete: func,
  onBack: func.isRequired,
  onSubmit: func.isRequired,
};

WeightTicketForm.defaultProps = {
  weightTicket: undefined,
  onUploadDelete: undefined,
  tripNumber: '1',
};

export default WeightTicketForm;
