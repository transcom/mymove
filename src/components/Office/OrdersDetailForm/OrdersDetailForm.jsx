import React, { useState } from 'react';
import { func, string, bool } from 'prop-types';

import styles from './OrdersDetailForm.module.scss';

import { dropdownInputOptions, formatLabelReportByDate } from 'utils/formatters';
import { CheckboxField, DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownArrayOf } from 'types/form';
import { SPECIAL_ORDERS_TYPES } from 'constants/orders';

const OrdersDetailForm = ({
  deptIndicatorOptions,
  ordersTypeOptions,
  ordersTypeDetailOptions,
  hhgTacWarning,
  hhgLoaWarning,
  ntsLoaWarning,
  ntsTacWarning,
  validateHHGTac,
  validateNTSTac,
  validateHHGLoa,
  validateNTSLoa,
  showDepartmentIndicator,
  showOrdersNumber,
  showOrdersTypeDetail,
  showHHGTac,
  showHHGSac,
  showNTSTac,
  showNTSSac,
  showHHGLoa,
  showNTSLoa,
  showOrdersAcknowledgement,
  ordersType,
  setFieldValue,
  payGradeOptions,
  formIsDisabled,
  hhgLongLineOfAccounting,
  ntsLongLineOfAccounting,
}) => {
  const [formOrdersType, setFormOrdersType] = useState(ordersType);
  const reportDateRowLabel = formatLabelReportByDate(formOrdersType);
  const noStarOrQuote = (value) => (/^[^*"]*$/.test(value) ? undefined : 'SAC cannot contain * or " characters');
  // The text/placeholder are different if the customer is retiring or separating.
  const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(formOrdersType);
  return (
    <div className={styles.OrdersDetailForm}>
      {requiredAsteriskMessage}
      <DutyLocationInput
        name="originDutyLocation"
        label="Current duty location"
        displayAddress={false}
        isDisabled={formIsDisabled}
        showRequiredAsterisk
        required
      />
      <DutyLocationInput
        name="newDutyLocation"
        label={isRetirementOrSeparation ? 'HOR, HOS or PLEAD' : 'New duty location'}
        displayAddress={false}
        placeholder={isRetirementOrSeparation ? 'Enter a city or ZIP' : 'Start typing a duty location...'}
        isDisabled={formIsDisabled}
        showRequiredAsterisk
        required
      />
      <DropdownInput
        data-testid="payGradeInput"
        name="payGrade"
        label="Pay grade"
        id="payGradeInput"
        options={payGradeOptions}
        showDropdownPlaceholderText={false}
        isDisabled={formIsDisabled}
        showRequiredAsterisk
        required
      />
      <DatePickerInput name="issueDate" label="Date issued" showRequiredAsterisk required disabled={formIsDisabled} />
      <DatePickerInput
        name="reportByDate"
        label={reportDateRowLabel}
        showRequiredAsterisk
        required
        disabled={formIsDisabled}
      />
      {showDepartmentIndicator && (
        <DropdownInput
          name="departmentIndicator"
          label="Department indicator"
          showRequiredAsterisk
          required
          options={deptIndicatorOptions}
          isDisabled={formIsDisabled}
        />
      )}
      {showOrdersNumber && (
        <TextField
          name="ordersNumber"
          label="Orders number"
          id="ordersNumberInput"
          showRequiredAsterisk
          required
          isDisabled={formIsDisabled}
        />
      )}
      <DropdownInput
        name="ordersType"
        label="Orders type"
        showRequiredAsterisk
        required
        options={
          formOrdersType === SPECIAL_ORDERS_TYPES.SAFETY_NON_LABEL || formOrdersType === SPECIAL_ORDERS_TYPES.BLUEBARK
            ? dropdownInputOptions({ SAFETY: 'Safety', BLUEBARK: 'Bluebark' })
            : ordersTypeOptions
        }
        onChange={(e) => {
          setFormOrdersType(e.target.value);
          setFieldValue('ordersType', e.target.value);
        }}
        isDisabled={
          formIsDisabled ||
          formOrdersType === SPECIAL_ORDERS_TYPES.SAFETY_NON_LABEL ||
          formOrdersType === SPECIAL_ORDERS_TYPES.BLUEBARK
        }
      />
      {showOrdersTypeDetail && (
        <DropdownInput
          name="ordersTypeDetail"
          label="Orders type detail"
          options={ordersTypeDetailOptions}
          showRequiredAsterisk
          required
          isDisabled={formIsDisabled}
        />
      )}
      <div className={styles.wrappedCheckbox}>
        <CheckboxField
          id="dependentsAuthorizedInput"
          data-testid="dependentsAuthorizedInput"
          name="dependentsAuthorized"
          label="Dependents authorized"
          isDisabled={formIsDisabled}
        />
      </div>
      {showHHGTac && showHHGSac && <h3>HHG accounting codes</h3>}
      {showHHGTac && (
        <MaskedTextField
          name="tac"
          label="TAC"
          id="hhgTacInput"
          mask="****"
          inputTestId="hhgTacInput"
          warning={hhgTacWarning}
          validate={validateHHGTac}
          showRequiredAsterisk
          required
          isDisabled={formIsDisabled}
        />
      )}
      {showHHGSac && (
        <MaskedTextField
          name="sac"
          label="SAC"
          mask={/[A-Za-z0-9]*/}
          id="hhgSacInput"
          inputTestId="hhgSacInput"
          data-testid="hhgSacInput"
          isDisabled={formIsDisabled}
          maxLength="80"
          validate={noStarOrQuote}
        />
      )}
      {showHHGTac && showHHGLoa && (
        <TextField
          name="hhgLoa"
          label="LOA"
          id="hhgLoaTextField"
          mask="****"
          inputTestId="hhgLoaTextField"
          data-testid="hhgLoaTextField"
          warning={hhgLoaWarning}
          validate={validateHHGLoa}
          value={hhgLongLineOfAccounting}
          isDisabled
        />
      )}

      {showNTSTac && showNTSSac && <h3>NTS accounting codes</h3>}
      {showNTSTac && (
        <MaskedTextField
          name="ntsTac"
          label="TAC"
          id="ntsTacInput"
          mask="****"
          inputTestId="ntsTacInput"
          warning={ntsTacWarning}
          validate={validateNTSTac}
          isDisabled={formIsDisabled}
        />
      )}
      {showNTSSac && (
        <MaskedTextField
          name="ntsSac"
          label="SAC"
          id="ntsSacInput"
          mask={/[A-Za-z0-9]*/}
          isDisabled={formIsDisabled}
          inputTestId="ntsSacInput"
          data-testid="ntsSacInput"
          maxLength="80"
          validate={noStarOrQuote}
        />
      )}
      {showNTSTac && showNTSLoa && (
        <TextField
          name="ntsLoa"
          label="LOA"
          id="ntsLoaTextField"
          mask="****"
          inputTestId="ntsLoaTextField"
          warning={ntsLoaWarning}
          data-testid="ntsLoaTextField"
          validate={validateNTSLoa}
          value={ntsLongLineOfAccounting}
          isDisabled
        />
      )}

      {showOrdersAcknowledgement && (
        <div className={styles.wrappedCheckbox}>
          <CheckboxField
            id="ordersAcknowledgementInput"
            name="ordersAcknowledgement"
            label="I have read the new orders"
            isDisabled={formIsDisabled}
          />
        </div>
      )}
    </div>
  );
};

OrdersDetailForm.propTypes = {
  ordersTypeOptions: DropdownArrayOf.isRequired,
  deptIndicatorOptions: DropdownArrayOf,
  ordersTypeDetailOptions: DropdownArrayOf,
  hhgTacWarning: string,
  ntsTacWarning: string,
  hhgLoaWarning: string,
  ntsLoaWarning: string,
  validateHHGTac: func,
  validateNTSTac: func,
  validateHHGLoa: func,
  validateNTSLoa: func,
  showDepartmentIndicator: bool,
  showOrdersNumber: bool,
  showOrdersTypeDetail: bool,
  showHHGTac: bool,
  showHHGSac: bool,
  showNTSTac: bool,
  showHHGLoa: bool,
  showNTSLoa: bool,
  showNTSSac: bool,
  showOrdersAcknowledgement: bool,
  ordersType: string.isRequired,
  setFieldValue: func.isRequired,
  formIsDisabled: bool,
  hhgLongLineOfAccounting: string,
  ntsLongLineOfAccounting: string,
};

OrdersDetailForm.defaultProps = {
  hhgTacWarning: '',
  ntsTacWarning: '',
  hhgLoaWarning: '',
  ntsLoaWarning: '',
  deptIndicatorOptions: null,
  ordersTypeDetailOptions: null,
  validateHHGTac: null,
  validateNTSTac: null,
  validateHHGLoa: null,
  validateNTSLoa: null,
  showDepartmentIndicator: true,
  showOrdersNumber: true,
  showOrdersTypeDetail: true,
  showHHGTac: true,
  showHHGSac: true,
  showNTSTac: true,
  showHHGLoa: true,
  showNTSLoa: true,
  showNTSSac: true,
  showOrdersAcknowledgement: false,
  formIsDisabled: false,
  hhgLongLineOfAccounting: '',
  ntsLongLineOfAccounting: '',
};

export default OrdersDetailForm;
