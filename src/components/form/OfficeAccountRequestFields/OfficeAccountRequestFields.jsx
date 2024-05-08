import React from 'react';
import { func } from 'prop-types';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DutyLocationInput } from 'components/form/fields';
import { searchTransportationOfficesOpen } from 'services/ghcApi';

export const OfficeAccountRequestFields = ({ render }) => {
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
          <TextField label="First Name" name={firstNameFieldName} id="officeAccountRequestFirstName" />
          <TextField
            label="Middle Initial"
            name={middleInitialFieldName}
            id="officeAccountRequestMiddleInitial"
            labelHint="optional"
          />
          <TextField label="Last Name" name={lastNameFieldName} id="officeAccountRequestLastName" />
          <TextField label="Email" name={emailField} id="officeAccountRequestEmail" />
          <MaskedTextField
            label="Telephone"
            id="officeAccountRequestTelephone"
            name={telephoneFieldName}
            type="tel"
            minimum="12"
            mask="000{-}000{-}0000"
          />
          <TextField
            label="DODID#"
            labelHint="10 digit number"
            name={edipiFieldName}
            id="officeAccountRequestEdipi"
            maxLength="10"
            inputMode="numeric"
            data-testid="officeAccountRequestEdipi"
          />
          <TextField
            label="Other Unique ID"
            labelHint="If not using DODID#"
            name={otherUniqueIdName}
            id="officeAccountRequestOtherUniqueId"
            data-testid="officeAccountRequestOtherUniqueId"
          />
          <DutyLocationInput
            name={transportationOfficeDropDown}
            label="Transportation Office"
            searchLocations={searchTransportationOfficesOpen}
          />
          <h4>Requested Role(s)</h4>
          <CheckboxField
            id="transportationOrderingOfficerCheckBox"
            data-testid="transportationOrderingOfficerCheckBox"
            name="transportationOrderingOfficerCheckBox"
            label="Transportation Ordering Officer"
          />
          <CheckboxField
            id="transportationInvoicingOfficerCheckBox"
            data-testid="transportationInvoicingOfficerCheckBox"
            name="transportationInvoicingOfficerCheckBox"
            label="Transportation Invoicing Officer"
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
            id="qualityAssuranceAndCustomerSupportCheckBox"
            data-testid="qualityAssuranceAndCustomerSupportCheckBox"
            name="qualityAssuranceAndCustomerSupportCheckBox"
            label="Quality Assurance & Customer Support"
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
