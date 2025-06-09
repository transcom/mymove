import React, { useEffect, useRef, useState } from 'react';
import { func } from 'prop-types';
import { ErrorMessage, Fieldset, Label } from '@trussworks/react-uswds';
import { useFormikContext } from 'formik';

import RequiredAsterisk from '../RequiredAsterisk';
import OptionalTag from '../OptionalTag';

import styles from './OfficeAccountRequestFields.module.scss';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DutyLocationInput } from 'components/form/fields';
import { searchTransportationOfficesOpen } from 'services/ghcApi';
import { isBooleanFlagEnabledUnauthenticatedOffice } from 'utils/featureFlags';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { useRolesPrivilegesQueriesOfficeApp } from 'hooks/queries';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';

export const OfficeAccountRequestFields = ({ render }) => {
  const { values, errors, touched, setFieldTouched, validateField } = useFormikContext();
  const [edipiRequired, setEdipiRequired] = useState(false);
  const [uniqueIdRequired, setUniqueIdRequired] = useState(false);
  const [enableRequestAccountPrivileges, setEnableRequestAccountPrivileges] = useState(false);

  const { result } = useRolesPrivilegesQueriesOfficeApp();
  const { privileges, rolesWithPrivs } = result;

  const availableRoles = rolesWithPrivs.filter((r) => r.roleType !== 'prime' && r.roleType !== 'customer');
  const hasAnyRoleSelected = React.useMemo(
    () => availableRoles.some(({ roleType }) => !!values[`${roleType}Checkbox`]),
    [availableRoles, values],
  );
  const hasTransportConflict = !!values.task_ordering_officerCheckbox && !!values.task_invoicing_officerCheckbox;

  const showRequestedRolesError = touched.requestedRolesGroup && !hasAnyRoleSelected;

  const showTransportConflictError = touched.transportationOfficerRoleConflict && hasTransportConflict;

  useEffect(() => {
    isBooleanFlagEnabledUnauthenticatedOffice(FEATURE_FLAG_KEYS.REQUEST_ACCOUNT_PRIVILEGES)?.then((enabled) => {
      setEnableRequestAccountPrivileges(enabled);
    });
  }, []);

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
      values.task_ordering_officerCheckbox,
      values.task_invoicing_officerCheckbox,
      values.services_counselorCheckbox,
      values.contracting_officerCheckbox,
      values.qaeCheckbox,
      values.headquartersCheckbox,
      values.customer_services_representativeCheckBox,
      values.gsrCheckbox,
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
    values.task_ordering_officerCheckbox,
    values.task_invoicing_officerCheckbox,
    values.services_counselorCheckbox,
    values.contracting_officerCheckbox,
    values.qaeCheckbox,
    values.headquartersCheckbox,
    values.customer_services_representativeCheckBox,
    values.gsrCheckbox,
    setFieldTouched,
    validateField,
  ]);

  const transportationOfficerTouched = useRef(false);
  useEffect(() => {
    const bothChecked = values.task_ordering_officerCheckbox || values.task_invoicing_officerCheckbox;

    if (!transportationOfficerTouched.current) {
      if (bothChecked) {
        transportationOfficerTouched.current = true;
      }
      return;
    }

    setFieldTouched('transportationOfficerRoleConflict', true, false);
    validateField('transportationOfficerRoleConflict');
  }, [values.task_ordering_officerCheckbox, values.task_invoicing_officerCheckbox, setFieldTouched, validateField]);

  const firstNameFieldName = 'officeAccountRequestFirstName';
  const middleInitialFieldName = 'officeAccountRequestMiddleInitial';
  const lastNameFieldName = 'officeAccountRequestLastName';
  const emailField = 'officeAccountRequestEmail';
  const telephoneFieldName = 'officeAccountRequestTelephone';
  const edipiFieldName = 'officeAccountRequestEdipi';
  const otherUniqueIdName = 'officeAccountRequestOtherUniqueId';
  const transportationOfficeDropDown = 'officeAccountTransportationOffice';

  const filteredPrivileges = privileges.filter((privilege) => {
    if (privilege.privilegeType === elevatedPrivilegeTypes.SAFETY) {
      return false;
    }
    return true;
  });

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
            required
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
            required
          />
          <Label data-testid="requestedRolesHeading">
            Requested Role(s)
            <RequiredAsterisk />
          </Label>
          {showRequestedRolesError && (
            <ErrorMessage
              id="requestedRolesGroupError"
              className={styles.errorText}
              data-testid="requestedRolesGroupError"
            >
              {errors.requestedRolesGroup}
            </ErrorMessage>
          )}

          {showTransportConflictError && (
            <ErrorMessage
              id="transportationOfficerRoleConflictError"
              className={styles.errorText}
              data-testid="transportationOfficerRoleConflictError"
            >
              {errors.transportationOfficerRoleConflict}
            </ErrorMessage>
          )}
          {availableRoles.map(({ roleType, roleName }) => {
            const fieldName = `${roleType}Checkbox`;
            const isTransportRole = roleType === 'task_ordering_officer' || roleType === 'task_invoicing_officer';

            const describedBy = [
              showRequestedRolesError && 'requestedRolesGroupError',
              isTransportRole && showTransportConflictError && 'transportationOfficerRoleConflictError',
            ]
              .filter(Boolean)
              .join(' ');

            return (
              <CheckboxField
                key={fieldName}
                id={fieldName}
                data-testid={fieldName}
                name={fieldName}
                label={roleName}
                aria-describedby={describedBy || undefined}
                aria-invalid={showRequestedRolesError || (isTransportRole && showTransportConflictError)}
              />
            );
          })}
          {enableRequestAccountPrivileges && (
            <>
              <Label data-testid="requestedPrivilegesHeading">
                Privilege(s)
                <OptionalTag />
              </Label>
              {filteredPrivileges.map(({ privilegeType, privilegeName }) => (
                <CheckboxField
                  id={`${privilegeType}PrivilegeCheckbox`}
                  data-testid={`${privilegeType}PrivilegeCheckbox`}
                  name={`${privilegeType}PrivilegeCheckbox`}
                  label={privilegeName}
                />
              ))}
            </>
          )}
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
