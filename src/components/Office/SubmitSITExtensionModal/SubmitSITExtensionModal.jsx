import React from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import moment from 'moment';

import styles from './SubmitSITExtensionModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { DropdownInput } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { sitExtensionReasons } from 'constants/sitExtensions';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';
import {
  calculateEndDate,
  calculateSitDaysAllowance,
  calculateDaysInPreviousSIT,
  calculateSITEndDate,
  calculateSITTotalDaysRemaining,
  CurrentSITDateData,
  formatSITAuthorizedEndDate,
  formatSITDepartureDate,
  formatSITEntryDate,
  getSITCurrentLocation,
  SitDaysAllowanceForm,
  SitEndDateForm,
} from 'utils/sitFormatters';

const SitStatusTables = ({ sitStatus, shipment }) => {
  const { totalSITDaysUsed } = sitStatus;
  const { daysInSIT } = sitStatus.currentSIT;
  const sitDepartureDate = formatSITDepartureDate(sitStatus.currentSIT.sitDepartureDate);
  const sitEntryDate = formatSITEntryDate(sitStatus.currentSIT.sitEntryDate);
  const daysInPreviousSIT = calculateDaysInPreviousSIT(totalSITDaysUsed, daysInSIT);
  const currentLocation = getSITCurrentLocation(sitStatus);
  const totalDaysRemaining = calculateSITTotalDaysRemaining(sitStatus, shipment);

  const sitAllowanceHelper = useField({ name: 'daysApproved', id: 'daysApproved' })[2];
  const endDateHelper = useField({ name: 'sitEndDate', id: 'sitEndDate' })[2];

  const currentDateEnteredSit = <p>{sitEntryDate}</p>;

  const handleSitEndDateChange = (endDate) => {
    const endDay = calculateEndDate(sitEntryDate, endDate);
    const calculatedSitDaysAllowance = calculateSitDaysAllowance(sitEntryDate, daysInPreviousSIT, endDay);

    // Update form values
    endDateHelper.setValue(endDate);
    sitAllowanceHelper.setValue(String(calculatedSitDaysAllowance));
  };

  const handleDaysAllowanceChange = (daysApproved) => {
    sitAllowanceHelper.setValue(daysApproved);
    const calculatedSITEndDate = calculateSITEndDate(sitEntryDate, daysApproved, daysInPreviousSIT);
    endDateHelper.setTouched(true);
    endDateHelper.setValue(calculatedSITEndDate);
  };

  return (
    <>
      <div className={styles.title}>
        <p>SIT (STORAGE IN TRANSIT)</p>
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        <DataTable
          custClass={styles.totalDaysTable}
          columnHeaders={['Total days of SIT approved', 'Total days used', 'Total days remaining']}
          dataRow={[
            <SitDaysAllowanceForm onChange={(e) => handleDaysAllowanceChange(e.target.value)} />,
            sitStatus.totalSITDaysUsed,
            totalDaysRemaining,
          ]}
        />
      </div>
      <div className={styles.tableContainer} data-testid="sitStartAndEndTable">
        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
        <DataTable
          columnHeaders={[`SIT start date`, 'SIT authorized end date', 'Calculated total SIT days']}
          dataRow={[
            currentDateEnteredSit,
            <SitEndDateForm
              onChange={(value) => {
                handleSitEndDateChange(value);
              }}
            />,
            sitStatus.calculatedTotalDaysInSIT,
          ]}
          custClass={styles.currentLocation}
        />
      </div>
      <div className={styles.tableContainer}>
        <CurrentSITDateData
          currentLocation={currentLocation}
          daysInSIT={daysInSIT}
          sitDepartureDate={sitDepartureDate}
        />
      </div>
    </>
  );
};

const SubmitSITExtensionModal = ({ shipment, sitStatus, onClose, onSubmit }) => {
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
    requestReason: '',
    officeRemarks: '',
    daysApproved: String(shipment.sitDaysAllowance),
    sitEndDate: formatSITAuthorizedEndDate(sitStatus),
  };
  const minimumDaysAllowed = sitStatus.calculatedTotalDaysInSIT - sitStatus.currentSIT.daysInSIT + 1;
  const sitEntryDate = moment(sitStatus.currentSIT.sitEntryDate, swaggerDateFormat);
  const reviewSITExtensionSchema = Yup.object().shape({
    requestReason: Yup.string().required('Required'),
    officeRemarks: Yup.string().nullable(),
    daysApproved: Yup.number()
      .min(minimumDaysAllowed, `Total days of SIT approved must be ${minimumDaysAllowed} or more.`)
      .required('Required'),
    sitEndDate: Yup.date().min(
      formatDateForDatePicker(sitEntryDate.add(1, 'days')),
      'The end date must occur after the start date. Please select a new date.',
    ),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.SubmitSITExtensionModal} onClose={() => onClose()}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Edit SIT authorization</h2>
          </ModalTitle>
          <Formik
            validationSchema={reviewSITExtensionSchema}
            onSubmit={(e) => onSubmit(e)}
            initialValues={initialValues}
          >
            {({ isValid }) => {
              return (
                <Form>
                  <DataTableWrapper className={classnames('maxw-tablet', styles.sitDisplayForm)} testID="sitExtensions">
                    <SitStatusTables sitStatus={sitStatus} shipment={shipment} />
                  </DataTableWrapper>
                  <div className={styles.reasonDropdown}>
                    <DropdownInput
                      label="Reason for edit"
                      name="requestReason"
                      data-testid="reasonDropdown"
                      options={dropdownInputOptions(sitExtensionReasons)}
                    />
                  </div>
                  <Label htmlFor="officeRemarks">Office remarks</Label>
                  <Field as={Textarea} data-testid="officeRemarks" label="No" name="officeRemarks" id="officeRemarks" />
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

SubmitSITExtensionModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default SubmitSITExtensionModal;
