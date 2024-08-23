import React, { useState } from 'react';
import { func, string, bool } from 'prop-types';

import styles from './OrdersDetailForm.module.scss';

import { dropdownInputOptions, formatLabelReportByDate } from 'utils/formatters';
import { CheckboxField, DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownArrayOf } from 'types/form';

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
  touched,
  ntsLongLineOfAccounting,
}) => {
  const [formOrdersType, setFormOrdersType] = useState(ordersType);
  const reportDateRowLabel = formatLabelReportByDate(formOrdersType);
  // The text/placeholder are different if the customer is retiring or separating.
  const isRetirementOrSeparation = ['RETIREMENT', 'SEPARATION'].includes(formOrdersType);
  return (
    <div className={styles.OrdersDetailForm}>
      <DutyLocationInput
        name="originDutyLocation"
        label="Current duty location"
        displayAddress={false}
        isDisabled={formIsDisabled}
        touched={touched}
      />
      <DutyLocationInput
        name="newDutyLocation"
        label={isRetirementOrSeparation ? 'HOR, HOS or PLEAD' : 'New duty location'}
        displayAddress={false}
        placeholder={isRetirementOrSeparation ? 'Enter a city or ZIP' : 'Start typing a duty location...'}
        isDisabled={formIsDisabled}
        touched={touched}
      />
      <DropdownInput
        data-testid="payGradeInput"
        name="payGrade"
        label="Pay grade"
        id="payGradeInput"
        options={payGradeOptions}
        showDropdownPlaceholderText={false}
        isDisabled={formIsDisabled}
      />
      <DatePickerInput name="issueDate" label="Date issued" disabled={formIsDisabled} />
      <DatePickerInput name="reportByDate" label={reportDateRowLabel} disabled={formIsDisabled} />
      {showDepartmentIndicator && (
        <DropdownInput
          name="departmentIndicator"
          label="Department indicator"
          options={deptIndicatorOptions}
          isDisabled={formIsDisabled}
        />
      )}
      {showOrdersNumber && (
        <TextField name="ordersNumber" label="Orders number" id="ordersNumberInput" isDisabled={formIsDisabled} />
      )}
      <DropdownInput
        name="ordersType"
        label="Orders type"
        options={formOrdersType === 'SAFETY' ? dropdownInputOptions({ SAFETY: 'Safety' }) : ordersTypeOptions}
        onChange={(e) => {
          setFormOrdersType(e.target.value);
          setFieldValue('ordersType', e.target.value);
        }}
        isDisabled={formIsDisabled || formOrdersType === 'SAFETY'}
      />
      {showOrdersTypeDetail && (
        <DropdownInput
          name="ordersTypeDetail"
          label="Orders type detail"
          options={ordersTypeDetailOptions}
          isDisabled={formIsDisabled}
        />
      )}

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
          isDisabled={formIsDisabled}
        />
      )}
      {showHHGSac && (
        <TextField
          name="sac"
          label="SAC"
          id="hhgSacInput"
          data-testid="hhgSacInput"
          isDisabled={formIsDisabled}
          maxLength="80"
          optional
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
          optional
        />
      )}
      {showNTSSac && (
        <TextField
          name="ntsSac"
          label="SAC"
          id="ntsSacInput"
          isDisabled={formIsDisabled}
          data-testid="ntsSacInput"
          maxLength="80"
          optional
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
  payGradeOptions: DropdownArrayOf,
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
  payGradeOptions: null,
  formIsDisabled: false,
  hhgLongLineOfAccounting: '',
  ntsLongLineOfAccounting: '',
};

export default OrdersDetailForm;
