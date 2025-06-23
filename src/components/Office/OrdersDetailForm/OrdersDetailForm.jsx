import React, { useEffect, useState } from 'react';
import { func, string, bool } from 'prop-types';

import styles from './OrdersDetailForm.module.scss';

import { dropdownInputOptions, formatLabelReportByDate } from 'utils/formatters';
import { CheckboxField, DropdownInput, DatePickerInput, DutyLocationInput } from 'components/form/fields';
import { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';
import TextField from 'components/form/fields/TextField/TextField';
import MaskedTextField from 'components/form/fields/MaskedTextField/MaskedTextField';
import { DropdownArrayOf } from 'types/form';
import { SPECIAL_ORDERS_TYPES } from 'constants/orders';
import { getRankOptions } from 'services/ghcApi';
import { sortRankOptions } from 'shared/utils';

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
  affiliation,
  handleChange,
  currentGrade,
}) => {
  const [formOrdersType, setFormOrdersType] = useState(ordersType);
  const reportDateRowLabel = formatLabelReportByDate(formOrdersType);
  const noStarOrQuote = (value) => (/^[^*"]*$/.test(value) ? undefined : 'SAC cannot contain * or " characters');
  const [grade, setGrade] = useState(currentGrade);

  const [rankOptions, setRankOptions] = useState([]);
  useEffect(() => {
    const fetchRankGradeOptions = async () => {
      // setShowLoadingSpinner(true, 'Loading Rank/Grade options');
      try {
        const fetchedRanks = await getRankOptions(affiliation, grade || currentGrade);
        if (fetchedRanks) {
          const formattedOptions = sortRankOptions(fetchedRanks.body);
          setRankOptions(formattedOptions);
        }
      } catch (error) {
        // const { message } = error;
        // milmoveLogger.error({ message, info: null });
        // retryPageLoading(error);
      }
      // setShowLoadingSpinner(false, null);
    };

    fetchRankGradeOptions();
  }, [affiliation, grade, currentGrade]);
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
      />
      <DutyLocationInput
        name="newDutyLocation"
        label={isRetirementOrSeparation ? 'HOR, HOS or PLEAD' : 'New duty location'}
        displayAddress={false}
        placeholder={isRetirementOrSeparation ? 'Enter a city or ZIP' : 'Start typing a duty location...'}
        isDisabled={formIsDisabled}
        showRequiredAsterisk
      />
      <DropdownInput
        data-testid="payGradeInput"
        label="Pay grade"
        name="grade"
        id="payGradeInput"
        options={payGradeOptions}
        showDropdownPlaceholderText={false}
        isDisabled={formIsDisabled}
        showRequiredAsterisk
        onChange={(e) => {
          setGrade(e.target.value);
          handleChange(e);
        }}
      />
      {grade !== '' ? (
        <DropdownInput label="Rank" name="rank" id="rankInput" required options={rankOptions} showRequiredAsterisk />
      ) : null}
      <DatePickerInput name="issueDate" label="Date issued" showRequiredAsterisk disabled={formIsDisabled} />
      <DatePickerInput name="reportByDate" label={reportDateRowLabel} showRequiredAsterisk disabled={formIsDisabled} />
      {showDepartmentIndicator && (
        <DropdownInput
          name="departmentIndicator"
          label="Department indicator"
          showRequiredAsterisk
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
          isDisabled={formIsDisabled}
        />
      )}
      <DropdownInput
        name="ordersType"
        label="Orders type"
        showRequiredAsterisk
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
