import React, { useEffect, useRef, useState, useMemo } from 'react';
import { func, array } from 'prop-types';
import { ErrorMessage, Fieldset, Label } from '@trussworks/react-uswds';
import { useFormikContext } from 'formik';

import RequiredAsterisk, { requiredAsteriskMessage } from '../RequiredAsterisk';

import styles from './OfficeAccountRequestFields.module.scss';

import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { CheckboxField, DutyLocationInput } from 'components/form/fields';
import { searchTransportationOfficesOpen } from 'services/ghcApi';
import { isBooleanFlagEnabledUnauthenticated } from 'utils/featureFlags';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import { elevatedPrivilegeTypes } from 'constants/userPrivileges';
import { roleTypes } from 'constants/userRoles';

export const OfficeAccountRequestFields = ({ render, rolesWithPrivs = [], privileges = [] }) => {
  const { values, errors, touched, setFieldTouched, validateField } = useFormikContext();
  const [edipiRequired, setEdipiRequired] = useState(false);
  const [uniqueIdRequired, setUniqueIdRequired] = useState(false);
  const [enableRequestAccountPrivileges, setEnableRequestAccountPrivileges] = useState(false);

  const filteredPrivileges = privileges.filter((privilege) => {
    if (privilege.privilegeType === elevatedPrivilegeTypes.SAFETY) {
      return false;
    }
    return true;
  });

  const availableRoles = rolesWithPrivs.filter((r) => r.roleType !== 'prime' && r.roleType !== roleTypes.CUSTOMER);

  const hasAnyRoleSelected = useMemo(
    () => availableRoles.some(({ roleType }) => !!values[`${roleType}Checkbox`]),
    [availableRoles, values],
  );
  const hasTransportConflict = !!values.task_ordering_officerCheckbox && !!values.task_invoicing_officerCheckbox;

  const showRequestedRolesError = touched.requestedRolesGroup && !hasAnyRoleSelected;

  const showTransportConflictError = touched.transportationOfficerRoleConflict && hasTransportConflict;

  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabledUnauthenticated(FEATURE_FLAG_KEYS.REQUEST_ACCOUNT_PRIVILEGES)?.then((enabled) => {
        setEnableRequestAccountPrivileges(enabled);
      });
    };
    fetchData();
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
    if (!firstInteractionOccurred.current) {
      if (hasAnyRoleSelected) {
        firstInteractionOccurred.current = true;
      }
      return;
    }
    setFieldTouched('requestedRolesGroup', true, false);
    validateField('requestedRolesGroup');
  }, [hasAnyRoleSelected, setFieldTouched, validateField]);

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
                const isTransportRole = roleType === roleTypes.TOO || roleType === roleTypes.TIO;

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
                  <Label data-testid="requestedPrivilegesHeading">Privilege(s)</Label>
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
            </fieldset>
          </div>
        </>,
      )}
    </Fieldset>
  );
};

OfficeAccountRequestFields.propTypes = {
  render: func,
  rolesWithPrivs: array,
  privileges: array,
};

OfficeAccountRequestFields.defaultProps = {
  render: (fields) => fields,
  rolesWithPrivs: [],
  privileges: [],
};

export default OfficeAccountRequestFields;
