import React, { useEffect, useRef, useState } from 'react';
import { func } from 'prop-types';
import { ErrorMessage, Fieldset } from '@trussworks/react-uswds';
import { useFormikContext } from 'formik';

import RequiredAsterisk, { requiredAsteriskMessage } from '../RequiredAsterisk';

import styles from './OfficeAccountRequestFields.module.scss';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DutyLocationInput } from 'components/form/fields';
import { searchTransportationOfficesOpen } from 'services/ghcApi';

export const OfficeAccountRequestFields = ({ render }) => {
  const { values, errors, touched, setFieldTouched, validateField } = useFormikContext();
  const [edipiRequired, setEdipiRequired] = useState(false);
  const [uniqueIdRequired, setUniqueIdRequired] = useState(false);

  useEffect(() => {
    if (values.officeAccountRequestEdipi !== '') {
      setEdipiRequired(true);
    } else {
      setEdipiRequired(false);
    }
    if (values.officeAccountRequestOtherUniqueId !== '') {
      setUniqueIdRequired(true);
    } else {
      setUniqueIdRequired(false);
    }
  }, [values.officeAccountRequestEdipi, values.officeAccountRequestOtherUniqueId]);

  const firstInteractionOccurred = useRef(false);
  useEffect(() => {
    const anyChecked = [
      values.taskOrderingOfficerCheckBox,
      values.taskInvoicingOfficerCheckBox,
      values.servicesCounselorCheckBox,
      values.transportationContractingOfficerCheckBox,
      values.qualityAssuranceEvaluatorCheckBox,
      values.headquartersCheckBox,
      values.customerSupportRepresentativeCheckBox,
      values.governmentSurveillanceRepresentativeCheckbox,
    ].some(Boolean);

    // only start marking the field as touched after initial mount
    if (!firstInteractionOccurred.current) {
      if (anyChecked) {
        firstInteractionOccurred.current = true;
      }
      return;
    }

    setFieldTouched('requestedRolesGroup', true, false);
    validateField('requestedRolesGroup');
  }, [
    values.taskOrderingOfficerCheckBox,
    values.taskInvoicingOfficerCheckBox,
    values.servicesCounselorCheckBox,
    values.transportationContractingOfficerCheckBox,
    values.qualityAssuranceEvaluatorCheckBox,
    values.headquartersCheckBox,
    values.customerSupportRepresentativeCheckBox,
    values.governmentSurveillanceRepresentativeCheckbox,
    setFieldTouched,
    validateField,
  ]);

  const transportationOfficerTouched = useRef(false);
  useEffect(() => {
    const bothChecked = values.taskOrderingOfficerCheckBox || values.taskInvoicingOfficerCheckBox;

    if (!transportationOfficerTouched.current) {
      if (bothChecked) {
        transportationOfficerTouched.current = true;
      }
      return;
    }

    setFieldTouched('transportationOfficerRoleConflict', true, false);
    validateField('transportationOfficerRoleConflict');
  }, [values.taskOrderingOfficerCheckBox, values.taskInvoicingOfficerCheckBox, setFieldTouched, validateField]);

  const firstNameFieldName = 'officeAccountRequestFirstName';
  const middleInitialFieldName = 'officeAccountRequestMiddleInitial';
  const lastNameFieldName = 'officeAccountRequestLastName';
  const emailField = 'officeAccountRequestEmail';
  const telephoneFieldName = 'officeAccountRequestTelephone';
  const edipiFieldName = 'officeAccountRequestEdipi';
  const otherUniqueIdName = 'officeAccountRequestOtherUniqueId';
  const transportationOfficeDropDown = 'officeAccountTransportationOffice';

  return (
    <Fieldset>
      {render(
        <>
          {requiredAsteriskMessage}
          <TextField
            label="First Name"
            name={firstNameFieldName}
            id="officeAccountRequestFirstName"
            data-testid="officeAccountRequestFirstName"
            showRequiredAsterisk
            required
          />
          <TextField
            label="Middle Initial"
            name={middleInitialFieldName}
            id="officeAccountRequestMiddleInitial"
            data-testid="officeAccountRequestMiddleInitial"
          />
          <TextField
            label="Last Name"
            name={lastNameFieldName}
            id="officeAccountRequestLastName"
            data-testid="officeAccountRequestLastName"
            showRequiredAsterisk
            required
          />
          <TextField
            label="Email"
            name={emailField}
            id="officeAccountRequestEmail"
            data-testid="officeAccountRequestEmail"
            showRequiredAsterisk
            required
          />
          <TextField
            label="Confirm Email"
            name="emailConfirmation"
            id="emailConfirmation"
            data-testid="emailConfirmation"
            disablePaste
            showRequiredAsterisk
            required
          />
          <MaskedTextField
            label="Telephone"
            id="officeAccountRequestTelephone"
            data-testid="officeAccountRequestTelephone"
            name={telephoneFieldName}
            type="tel"
            minimum="12"
            mask="000{-}000{-}0000"
            showRequiredAsterisk
            required
          />
          <div className={styles.section}>
            <div className={styles.inputContainer}>
              <TextField
                label="DODID#"
                aria-label="D O D I D # is required if not using other unique identifier"
                labelHint="10 digit number"
                name={edipiFieldName}
                id="officeAccountRequestEdipi"
                data-testid="officeAccountRequestEdipi"
                maxLength="10"
                inputMode="numeric"
              />
            </div>
            <div className={styles.inputContainer}>
              <TextField
                label="Confirm DODID#"
                aria-label="Confirm D O D I D # is required if D O D I D # is being used"
                name="edipiConfirmation"
                id="edipiConfirmation"
                data-testid="edipiConfirmation"
                maxLength="10"
                disablePaste
                showRequiredAsterisk={edipiRequired}
              />
            </div>
          </div>
          <div className={styles.section}>
            <div className={styles.inputContainer}>
              <TextField
                label="Other Unique ID"
                labelHint="If not using DODID#"
                aria-label="Other Unique ID is required if not using D O D I D #"
                name={otherUniqueIdName}
                id="officeAccountRequestOtherUniqueId"
                data-testid="officeAccountRequestOtherUniqueId"
              />
            </div>
            <div className={styles.inputContainer}>
              <TextField
                label="Confirm Other Unique ID"
                name="otherUniqueIdConfirmation"
                aria-label="Confirm Other Unique ID is required if using Other Unique ID"
                id="otherUniqueIdConfirmation"
                data-testid="otherUniqueIdConfirmation"
                disablePaste
                showRequiredAsterisk={uniqueIdRequired}
              />
            </div>
          </div>
          <DutyLocationInput
            data-testid="transportationOfficeSelector"
            name={transportationOfficeDropDown}
            label="Transportation Office"
            searchLocations={searchTransportationOfficesOpen}
            showRequiredAsterisk
            required
          />
          <div className="margin-top-2">
            <fieldset>
              <legend className="usa-label" aria-label="At least one requested role is required.">
                <span data-testid="requestedRolesHeadingSpan">
                  Requested Role(s) <RequiredAsterisk />
                </span>
              </legend>

              {errors.requestedRolesGroup && touched.requestedRolesGroup && (
                <ErrorMessage
                  id="requestedRolesGroupError"
                  className={styles.errorText}
                  data-testid="requestedRolesGroupError"
                >
                  {errors.requestedRolesGroup}
                </ErrorMessage>
              )}
              <CheckboxField
                id="headquartersCheckBox"
                data-testid="headquartersCheckBox"
                name="headquartersCheckBox"
                label="Headquarters"
                aria-describedby={errors.requestedRolesGroup ? 'requestedRolesGroupError' : undefined}
                aria-invalid={!!errors.requestedRolesGroup}
              />
              {errors.transportationOfficerRoleConflict && touched.transportationOfficerRoleConflict && (
                <ErrorMessage
                  id="transportationOfficerRoleConflictError"
                  className={styles.errorText}
                  data-testid="transportationOfficerRoleConflictError"
                >
                  {errors.transportationOfficerRoleConflict}
                </ErrorMessage>
              )}
              <CheckboxField
                id="taskOrderingOfficerCheckBox"
                data-testid="taskOrderingOfficerCheckBox"
                name="taskOrderingOfficerCheckBox"
                label="Task Ordering Officer"
                aria-describedby={[
                  errors.requestedRolesGroup && touched.requestedRolesGroup ? 'requestedRolesGroupError' : null,
                  errors.transportationOfficerRoleConflict && touched.transportationOfficerRoleConflict
                    ? 'transportationOfficerRoleConflictError'
                    : null,
                ]
                  .filter(Boolean)
                  .join(' ')}
                aria-invalid={!!errors.requestedRolesGroup || !!errors.transportationOfficerRoleConflict}
              />
              <CheckboxField
                id="taskInvoicingOfficerCheckBox"
                data-testid="taskInvoicingOfficerCheckBox"
                name="taskInvoicingOfficerCheckBox"
                label="Task Invoicing Officer"
                aria-describedby={[
                  errors.requestedRolesGroup && touched.requestedRolesGroup ? 'requestedRolesGroupError' : null,
                  errors.transportationOfficerRoleConflict && touched.transportationOfficerRoleConflict
                    ? 'transportationOfficerRoleConflictError'
                    : null,
                ]
                  .filter(Boolean)
                  .join(' ')}
                aria-invalid={!!errors.requestedRolesGroup || !!errors.transportationOfficerRoleConflict}
              />

              <CheckboxField
                id="transportationContractingOfficerCheckBox"
                data-testid="transportationContractingOfficerCheckBox"
                name="transportationContractingOfficerCheckBox"
                label="Contracting Officer"
                aria-describedby={errors.requestedRolesGroup ? 'requestedRolesGroupError' : undefined}
                aria-invalid={!!errors.requestedRolesGroup}
              />
              <CheckboxField
                id="servicesCounselorCheckBox"
                data-testid="servicesCounselorCheckBox"
                name="servicesCounselorCheckBox"
                label="Services Counselor"
                aria-describedby={errors.requestedRolesGroup ? 'requestedRolesGroupError' : undefined}
                aria-invalid={!!errors.requestedRolesGroup}
              />
              <CheckboxField
                id="qualityAssuranceEvaluatorCheckBox"
                data-testid="qualityAssuranceEvaluatorCheckBox"
                name="qualityAssuranceEvaluatorCheckBox"
                label="Quality Assurance Evaluator"
                aria-describedby={errors.requestedRolesGroup ? 'requestedRolesGroupError' : undefined}
                aria-invalid={!!errors.requestedRolesGroup}
              />
              <CheckboxField
                id="customerSupportRepresentativeCheckBox"
                data-testid="customerSupportRepresentativeCheckBox"
                name="customerSupportRepresentativeCheckBox"
                label="Customer Support Representative"
                aria-describedby={errors.requestedRolesGroup ? 'requestedRolesGroupError' : undefined}
                aria-invalid={!!errors.requestedRolesGroup}
              />
              <CheckboxField
                id="governmentSurveillanceRepresentativeCheckbox"
                data-testid="governmentSurveillanceRepresentativeCheckbox"
                name="governmentSurveillanceRepresentativeCheckbox"
                label="Government Surveillance Representative"
                aria-describedby={errors.requestedRolesGroup ? 'requestedRolesGroupError' : undefined}
                aria-invalid={!!errors.requestedRolesGroup}
              />
            </fieldset>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

OfficeAccountRequestFields.propTypes = {
  render: func,
};

OfficeAccountRequestFields.defaultProps = {
  render: (fields) => fields,
};

export default OfficeAccountRequestFields;
