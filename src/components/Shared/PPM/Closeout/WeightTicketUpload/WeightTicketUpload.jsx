import React from 'react';
import classnames from 'classnames';
import { FormGroup, Label, Link, Alert } from '@trussworks/react-uswds';
import { string, bool, func, shape } from 'prop-types';

import styles from './WeightTicketUpload.module.scss';

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
    <Alert type="info">
      If you do not upload legible certified weight tickets, your PPM incentive could be affected.
    </Alert>
    <p>Download the official government spreadsheet to calculate constructed weight.</p>
    <Link
      className={classnames('usa-button', 'usa-button--secondary', styles.constructedWeightLink)}
      href="https://www.ustranscom.mil/dp3/weightestimator.cfm"
      target="_blank"
      rel="noopener"
    >
      Go to download page
    </Link>
    <p className={styles.bold}>Enter the sum of your constructed weight and the empty weight as the full weight.</p>
  </>
);

const rentalAgreement = (
  <>
    <Alert type="info">
      If you do not upload legible certified weight tickets, your PPM incentive could be affected.
    </Alert>
    <p>Enter the PPM vehicle&apos;s weight as the empty weight. Your vehicle&apos;s weight can be obtained from:</p>
    <ul>
      <li>The Branham Automobile Reference Book</li>
      <li>
        <Link href="https://www.jdpower.com/cars" target="_blank" rel="noopener">
          National Automobile Dealers Association (NADA) Official Used Car Guide
        </Link>
      </li>
      <li>Your owner&apos;s manual</li>
      <li>Other appropriate reference sources of manufacturer&apos;s weight</li>
    </ul>
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
  formikProps: { setFieldTouched, setFieldValue },
}) => {
  const weightTicketRentalAgreement = fieldName === 'emptyDocument' && missingWeightTicket;

  const weightTicketUploadLabel = (name, isMissingWeightTicket) => {
    if (isMissingWeightTicket || name === 'missingProGearWeightDocument') {
      if (weightTicketRentalAgreement) {
        return `Since you do not have a certified weight ticket, upload the registration or rental agreement for the vehicle used
        during the PPM`;
      }
      return 'Upload your completed constructed weight spreadsheet';
    }

    if (name === 'emptyDocument') {
      return 'Upload empty weight ticket';
    }

    if (name === 'document') {
      return "Upload your pro-gear's weight tickets";
    }

    if (name === 'gunSafeDocument') {
      return "Upload your gun safe's weight tickets";
    }

    return 'Upload full weight ticket';
  };

  const weightTicketUploadHint = () => {
    return missingWeightTicket && !weightTicketRentalAgreement
      ? SpreadsheetUploadInstructions
      : DocumentAndImageUploadInstructions;
  };

  return (
    <div className={styles.WeightTicketUpload}>
      {missingWeightTicket && weightTicketRentalAgreement && rentalAgreement}
      {missingWeightTicket && !weightTicketRentalAgreement && constructedWeightDownload}
      <UploadsTable
        className={styles.uploadsTable}
        uploads={values[`${fieldName}`]}
        onDelete={(uploadId) => onUploadDelete(uploadId, fieldName, setFieldTouched, setFieldValue)}
      />
      <FormGroup>
        <div className="labelWrapper">
          <Label htmlFor={fieldName}> {weightTicketUploadLabel(fieldName, missingWeightTicket)} </Label>
        </div>
        <Hint className={styles.uploadTypeHint}> {weightTicketUploadHint()} </Hint>
        <FileUpload
          name={fieldName}
          className={fieldName}
          labelIdle={UploadDropZoneLabel}
          labelIdleMobile={UploadDropZoneLabelMobile}
          createUpload={(file) => onCreateUpload(fieldName, file, setFieldTouched)}
          onChange={(err, upload) => {
            onUploadComplete(err);
            fileUploadRef?.current?.removeFile(upload.id);
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
