import React from 'react';
import classnames from 'classnames';
import { Formik, Field } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import moment from 'moment';

import styles from './ConvertSITToCustomerExpenseModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { LOCATION_TYPES } from 'types/sitStatusShape';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';
import { formatDate } from 'utils/formatters';
import { DEFAULT_EMPTY_VALUE } from 'shared/constants';
import RequiredAsterisk, { requiredAsteriskMessage } from 'components/form/RequiredAsterisk';

const SitStatusTables = ({ sitStatus, shipment }) => {
  const { daysInSIT } = sitStatus.currentSIT;
  const sitEntryDate = moment(sitStatus.currentSIT.sitEntryDate, swaggerDateFormat);
  const sitDepartureDate =
    formatDate(sitStatus.currentSIT?.sitDepartureDate, swaggerDateFormat, 'DD MMM YYYY') || DEFAULT_EMPTY_VALUE;
  // Currently active SIT
  const currentLocation = sitStatus.currentSIT.location === LOCATION_TYPES.ORIGIN ? 'origin SIT' : 'destination SIT';

  const currentDaysInSit = <p>{daysInSIT}</p>;
  const currentDateEnteredSit = <p>{formatDateForDatePicker(sitEntryDate)}</p>;
  const sitEndDate = moment(sitStatus.currentSIT?.sitAuthorizedEndDate, swaggerDateFormat);
  const sitEndDateString = sitEndDate.isValid() ? formatDateForDatePicker(sitEndDate) : '—';

  const totalDaysRemaining = () => {
    const daysRemaining = sitStatus ? sitStatus.totalDaysRemaining : shipment.sitDaysAllowance;
    if (daysRemaining > 0) {
      return daysRemaining;
    }
    return 'Expired';
  };

  return (
    <>
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT)</p>
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        {/* Sit Total days table */}
        <DataTable
          custClass={styles.totalDaysTable}
          columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
          dataRow={[shipment.sitDaysAllowance, sitStatus.totalSITDaysUsed, totalDaysRemaining()]}
        />
      </div>
      <div className={styles.tableContainer} data-testid="sitStartAndEndTable">
        {/* Sit Start and End table */}
        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
        <DataTable
          columnHeaders={[`SIT start date`, 'SIT authorized end date', 'Calculated total SIT days']}
          dataRow={[currentDateEnteredSit, sitEndDateString, sitStatus.calculatedTotalDaysInSIT]}
          custClass={styles.currentLocation}
        />
      </div>
      <div className={styles.tableContainer}>
        {/* Total days at current location */}
        <DataTable
          columnHeaders={[`Total days in ${currentLocation}`, `SIT departure date`]}
          dataRow={[currentDaysInSit, sitDepartureDate]}
        />
      </div>
    </>
  );
};

const ConvertSITToCustomerExpenseModal = ({ shipment, sitStatus, onClose, onSubmit }) => {
  let sitStartDate = sitStatus?.sitEntryDate;
  if (!sitStartDate) {
    sitStartDate = shipment.mtoServiceItems?.reduce((item, acc) => {
      if (item.sitEntryDate < acc.sitEntryDate) {
        return item;
      }
      return acc;
    }).sitEntryDate;
  }

  const initialValues = {
    remarks: '',
    convertToCustomersExpense: true,
  };
  const convertSITToCustomerExpenseSchema = Yup.object().shape({
    remarks: Yup.string().required('Required'),
    convertToCustomersExpense: Yup.boolean().required('Required'),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.ConvertSITToCustomerExpenseModal} onClose={() => onClose()}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Convert SIT To Customer Expense</h2>
          </ModalTitle>
          <Formik
            validationSchema={convertSITToCustomerExpenseSchema}
            onSubmit={(e) => onSubmit(true, e.remarks)}
            initialValues={initialValues}
            validateOnMount
          >
            {({ isValid }) => {
              return (
                <Form>
                  <DataTableWrapper className={classnames('maxw-tablet', styles.sitDisplayForm)} testID="sitExtensions">
                    <SitStatusTables sitStatus={sitStatus} shipment={shipment} />
                  </DataTableWrapper>
                  {requiredAsteriskMessage}
                  <Label htmlFor="remarks" required>
                    <span required>
                      Remarks <RequiredAsterisk />
                    </span>
                  </Label>
                  <Field as={Textarea} data-testid="remarks" label="No" name="remarks" id="remarks" required />
                  <ModalActions>
                    <Button
                      type="button"
                      onClick={() => onClose()}
                      data-testid="modalCancelButton"
                      secondary
                      className={styles.CancelButton}
                    >
                      Cancel
                    </Button>
                    <Button type="submit" disabled={!isValid}>
                      Save
                    </Button>
                  </ModalActions>
                </Form>
              );
            }}
          </Formik>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ConvertSITToCustomerExpenseModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default ConvertSITToCustomerExpenseModal;
