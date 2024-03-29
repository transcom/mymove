import React from 'react';
import { func, node, string } from 'prop-types';
import { Fieldset } from '@trussworks/react-uswds';

import TextField from 'components/form/fields/TextField/TextField';
// import { DropdownInput } from 'components/form/fields/DropdownInput';
// import { CheckboxField, DutyLocationInput } from 'components/form/fields';
import { CheckboxField } from 'components/form/fields';
// import { searchTransportationOffices } from 'services/ghcApi';

export const OfficeAccountRequestFields = ({ legend, className, render }) => {
  const firstNameFieldName = 'officeAccountRequestFirstName';
  const middleInitialFieldName = 'officeAccountRequestMiddleInitial';
  const lastNameFieldName = 'officeAccountRequestLastName';
  const emailField = 'officeAccountRequestEmail';
  const telephoneFieldName = 'officeAccountRequestTelephone';
  const edipiFieldName = 'officeAccountRequestEdipi';
  const otherUniqueIdName = 'officeAccountRequestOtherUniqueId';
  // const transportationOfficeDropDown = 'officeAccountTransportationOffice';

  return (
    <Fieldset legend={legend} className={className}>
      {render(
        <>
          <TextField label="First Name" name={firstNameFieldName} id="officeAccountRequestFirstName" />
          <TextField
            label="Middle Initial"
            name={middleInitialFieldName}
            id="officeAccountRequestMiddleInitial"
            optional
          />
          <TextField label="Last Name" name={lastNameFieldName} id="officeAccountRequestLastName" />
          <TextField label="Email" name={emailField} id="officeAccountRequestEmail" />
          <TextField label="Telephone" name={telephoneFieldName} id="officeAccountRequestTelephone" />
          <TextField
            label="DoD ID number | EDIPI"
            labelHint="10 digit number"
            name={edipiFieldName}
            id="officeAccountRequestEdipi"
            maxLength="10"
            inputMode="numeric"
          />
          <TextField
            label="Other unique identifier"
            labelHint="If using PIV"
            name={otherUniqueIdName}
            id="officeAccountRequestOtherUniqueId"
            maxLength="10"
            inputMode="numeric"
          />
          {/* <DropdownInput
            name={transportationOfficeDropDown}
            id="officeAccountRequestTransportationOfficeDropdown"
            label="Transportation Office"
            placeHolderText="Select Transportation Office"
          /> */}
          {/* <DutyLocationInput
            name={transportationOfficeDropDown}
            label="Transportation Office"
            placeHolderText="Select Transportation Office"
            searchLocations={searchTransportationOffices}
          /> */}
          <h4>Requested Role(s)</h4>
          <CheckboxField
            id="transportationOrderingOfficerCheckBox"
            name="transportationOrderingOfficerCheckBox"
            label="Transportation Ordering Officer"
          />
          <CheckboxField
            id="transportationInvoicingOfficerCheckBox"
            name="transportationInvoicingOfficerCheckBox"
            label="Transportation Invoicing Officer"
          />
          <CheckboxField
            id="transportationContractingOfficerCheckBox"
            name="transportationContractingOfficerCheckBox"
            label="Contracting Officer"
          />
          <CheckboxField id="servicesCounselorCheckBox" name="servicesCounselorCheckBox" label="Services Counselor" />
          <CheckboxField
            id="qualityAssuranceAndCustomerSupportCheckBox"
            name="qualityAssuranceAndCustomerSupportCheckBox"
            label="Quality Assurance & Customer Support"
          />
        </>,
      )}
    </Fieldset>
  );
};

OfficeAccountRequestFields.propTypes = {
  legend: node,
  className: string,
  render: func,
};

OfficeAccountRequestFields.defaultProps = {
  legend: '',
  className: '',
  render: (fields) => fields,
};

export default OfficeAccountRequestFields;
