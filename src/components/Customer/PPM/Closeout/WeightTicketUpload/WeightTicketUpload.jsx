import React from 'react';
import classnames from 'classnames';
import { ErrorMessage, FormGroup, Label, Link } from '@trussworks/react-uswds';
import { string, bool, func, shape } from 'prop-types';

import styles from 'components/Customer/PPM/Closeout/WeightTicketUpload/WeightTicketUpload.module.scss';
import Hint from 'components/Hint';
import FileUpload from 'components/FileUpload/FileUpload';
import UploadsTable from 'components/UploadsTable/UploadsTable';
import {
  DocumentAndImageUploadInstructions,
  SpreadsheetUploadInstructions,
  UploadDropZoneLabel,
  UploadDropZoneLabelMobile,
} from 'content/uploads';

export const acceptableFileTypes = [
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
    if (showConstructedWeight || name === 'missingProGearWeightDocument') {
      return 'Upload constructed weight spreadsheet';
    }

    if (name === 'emptyDocument') {
      return 'Upload empty weight ticket';
    }

    if (name === 'document') {
      return "Upload your pro-gear's weight tickets";
    }

    return 'Upload full weight ticket';
  };

  const weightTicketUploadHint = (showConstructedWeight) => {
    return showConstructedWeight ? SpreadsheetUploadInstructions : DocumentAndImageUploadInstructions;
  };

  const showError = touched[`${fieldName}`] && errors[`${fieldName}`];

  return (
    <div className={styles.WeightTicketUpload}>
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
          className={fieldName}
          labelIdle={UploadDropZoneLabel}
          labelIdleMobile={UploadDropZoneLabelMobile}
          createUpload={(file) => onCreateUpload(fieldName, file, setFieldTouched)}
          onChange={(err, upload) => {
            onUploadComplete(err);
            fileUploadRef.current.removeFile(upload.id);
          }}
          acceptedFileTypes={acceptableFileTypes}
          labelFileTypeNotAllowed="Upload a supported file type"
          fileValidateTypeLabelExpectedTypes="Supported file types: PDF, JPG, PNG, XLS, or XLSX"
          ref={fileUploadRef}
        />
      </FormGroup>
    </div>
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

export default WeightTicketUpload;
