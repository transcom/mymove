import React, { useEffect, useState } from 'react';
import { func } from 'prop-types';
import { Fieldset, Label } from '@trussworks/react-uswds';
import { useFormikContext } from 'formik';

import styles from './OfficeAccountRequestFields.module.scss';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DutyLocationInput } from 'components/form/fields';
import { searchTransportationOfficesOpen } from 'services/ghcApi';

export const OfficeAccountRequestFields = ({ render }) => {
  const { values } = useFormikContext();
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
          <TextField
            label="First Name"
            name={firstNameFieldName}
            id="officeAccountRequestFirstName"
            data-testid="officeAccountRequestFirstName"
            showRequiredAsterisk
          />
          <TextField
            label="Middle Initial"
            name={middleInitialFieldName}
            id="officeAccountRequestMiddleInitial"
            data-testid="officeAccountRequestMiddleInitial"
            labelHint="optional"
          />
          <TextField
            label="Last Name"
            name={lastNameFieldName}
            id="officeAccountRequestLastName"
            data-testid="officeAccountRequestLastName"
            showRequiredAsterisk
          />
          <TextField
            label="Email"
            name={emailField}
            id="officeAccountRequestEmail"
            data-testid="officeAccountRequestEmail"
            showRequiredAsterisk
          />
          <TextField
            label="Confirm Email"
            name="emailConfirmation"
            id="emailConfirmation"
            data-testid="emailConfirmation"
            disablePaste
            showRequiredAsterisk
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
          />
          <div className={styles.section}>
            <div className={styles.inputContainer}>
              <TextField
                label="DODID#"
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
                name={otherUniqueIdName}
                id="officeAccountRequestOtherUniqueId"
                data-testid="officeAccountRequestOtherUniqueId"
              />
            </div>
            <div className={styles.inputContainer}>
              <TextField
                label="Confirm Other Unique ID"
                name="otherUniqueIdConfirmation"
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
          />
          <Label data-testid="requestedRolesHeading">
            Requested Role(s)
            <span data-testid="requiredAsterisk" className={styles.requiredAsterisk}>
              *
            </span>
          </Label>
          <CheckboxField
            id="headquartersCheckBox"
            data-testid="headquartersCheckBox"
            name="headquartersCheckBox"
            label="Headquarters"
          />
          <CheckboxField
            id="taskOrderingOfficerCheckBox"
            data-testid="taskOrderingOfficerCheckBox"
            name="taskOrderingOfficerCheckBox"
            label="Task Ordering Officer"
          />
          <CheckboxField
            id="taskInvoicingOfficerCheckBox"
            data-testid="taskInvoicingOfficerCheckBox"
            name="taskInvoicingOfficerCheckBox"
            label="Task Invoicing Officer"
          />
          <CheckboxField
            id="transportationContractingOfficerCheckBox"
            data-testid="transportationContractingOfficerCheckBox"
            name="transportationContractingOfficerCheckBox"
            label="Contracting Officer"
          />
          <CheckboxField
            id="servicesCounselorCheckBox"
            data-testid="servicesCounselorCheckBox"
            name="servicesCounselorCheckBox"
            label="Services Counselor"
          />
          <CheckboxField
            id="qualityAssuranceEvaluatorCheckBox"
            data-testid="qualityAssuranceEvaluatorCheckBox"
            name="qualityAssuranceEvaluatorCheckBox"
            label="Quality Assurance Evaluator"
          />
          <CheckboxField
            id="customerSupportRepresentativeCheckBox"
            data-testid="customerSupportRepresentativeCheckBox"
            name="customerSupportRepresentativeCheckBox"
            label="Customer Support Representative"
          />
          <CheckboxField
            id="governmentSurveillanceRepresentativeCheckbox"
            data-testid="governmentSurveillanceRepresentativeCheckbox"
            name="governmentSurveillanceRepresentativeCheckbox"
            label="Government Surveillance Representative"
          />
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
