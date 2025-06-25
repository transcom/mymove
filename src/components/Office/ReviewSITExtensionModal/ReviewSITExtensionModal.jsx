import React, { useState } from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { Button, Radio, FormGroup, Label, Tag, Textarea, Fieldset } from '@trussworks/react-uswds';
import moment from 'moment';

import { SITExtensionShape } from '../../../types/sitExtensions';

import styles from './ReviewSITExtensionModal.module.scss';
import ConfirmCustomerExpenseModal from './ConfirmCustomerExpenseModal/ConfirmCustomerExpenseModal';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import { DropdownInput, CheckboxField } from 'components/form/fields';
import { dropdownInputOptions } from 'utils/formatters';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle } from 'components/Modal/Modal';
import { sitExtensionReasons } from 'constants/sitExtensions';
import { formatDateForDatePicker } from 'shared/dates';
import { SitStatusShape } from 'types/sitStatusShape';
import { ShipmentShape } from 'types';
import {
  calculateEndDate,
  calculateSitDaysAllowance,
  calculateDaysInPreviousSIT,
  calculateSITEndDate,
  CurrentSITDateData,
  formatSITDepartureDate,
  formatSITEntryDate,
  getSITCurrentLocation,
  SitDaysAllowanceForm,
  SitEndDateForm,
  calculateApprovedAndRequestedDaysCombined,
  calculateSITTotalDaysRemaining,
  calculateApprovedAndRequestedDatesCombined,
  formatEndDate,
  SITHistoryItemHeaderDays,
  SITHistoryItemHeaderDate,
} from 'utils/sitFormatters';

const SitStatusTables = ({ sitStatus, sitExtension, shipment }) => {
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

  const approvedAndRequestedDaysCombined = calculateApprovedAndRequestedDaysCombined(shipment, sitExtension);
  const approvedAndRequestedDatesCombined = calculateApprovedAndRequestedDatesCombined(
    sitExtension,
    totalDaysRemaining,
  );

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
        <Tag>SIT EXTENSION REQUESTED</Tag>
      </div>
      <div className={styles.tableContainer} data-testid="sitStatusTable">
        <DataTable
          custClass={styles.totalDaysTable}
          columnHeaders={[
            <SITHistoryItemHeaderDays
              title="Total days of SIT proposed"
              approved={shipment.sitDaysAllowance}
              requested={sitExtension.requestedDays}
              value={approvedAndRequestedDaysCombined}
            />,
            'Total days used',
            'Proposed total days remaining (if extension request is approved)',
          ]}
          dataRow={[
            <SitDaysAllowanceForm onChange={(e) => handleDaysAllowanceChange(e.target.value)} />,
            sitStatus.totalSITDaysUsed,
            approvedAndRequestedDaysCombined - sitStatus.totalSITDaysUsed,
          ]}
        />
      </div>
      <div className={styles.tableContainer} data-testid="sitStartAndEndTable">
        <p className={styles.sitHeader}>Current location: {currentLocation}</p>
        <DataTable
          columnHeaders={[
            `SIT start date`,
            <SITHistoryItemHeaderDate
              title="Proposed SIT authorized end date"
              endDate={sitStatus.currentSIT.sitAuthorizedEndDate}
              requested={sitExtension.requestedDays}
              value={approvedAndRequestedDatesCombined}
            />,
            'Calculated total SIT days',
          ]}
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

/**
 * @description This component contains a form that can be viewed from the SIT
 * Display on the MTO page when the Prime submits a SIT Extension for review of
 * the TOO.
 */
const ReviewSITExtensionsModal = ({ onClose, sitExtension, shipment, sitStatus, onSubmit }) => {
  const [showConfirmCustomerExpenseModal, setShowConfirmCustomerExpenseModal] = useState(false);
  const approvedAndRequestedDaysCombined = calculateApprovedAndRequestedDaysCombined(shipment, sitExtension);
  const initialValues = {
    acceptExtension: '',
    convertToCustomerExpense: false,
    daysApproved: String(approvedAndRequestedDaysCombined),
    requestReason: sitExtension.requestReason,
    officeRemarks: '',
    sitEndDate: formatEndDate(sitStatus.currentSIT.sitAuthorizedEndDate, sitExtension.requestedDays),
  };
  const minimumDaysAllowed = shipment.sitDaysAllowance;
  const sitEntryDate = formatSITEntryDate(sitStatus.currentSIT.sitEntryDate);

  const reviewSITExtensionSchema = Yup.object().shape({
    acceptExtension: Yup.mixed().oneOf(['yes', 'no']).required('Required'),
    convertToCustomerExpense: Yup.boolean().default(false),
    requestReason: Yup.string().required('Required'),
    officeRemarks: Yup.string().when('acceptExtension', {
      is: 'no',
      then: () => Yup.string().required('Required'),
      otherwise: () => Yup.string().nullable(),
    }),
    daysApproved: Yup.number().when('acceptExtension', {
      is: 'yes',
      then: () =>
        Yup.number()
          .min(minimumDaysAllowed, `Total days of SIT approved must be ${minimumDaysAllowed} or more.`)
          .required('Required'),
    }),
    sitEndDate: Yup.date().min(
      formatDateForDatePicker(moment(sitEntryDate).add(1, 'days')),
      'The end date must occur after the start date. Please select a new date.',
    ),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.ReviewSITExtensionModal} onClose={() => onClose()}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Review SIT Extension Request</h2>
          </ModalTitle>
          <div>
            <Formik
              validationSchema={reviewSITExtensionSchema}
              onSubmit={(e) => onSubmit(sitExtension.id, e)}
              initialValues={initialValues}
            >
              {({ isValid, values, setValues }) => {
                const handleRadioSelection = (e) => {
                  if (e.target.value === 'yes') {
                    setValues({
                      ...values,
                      acceptExtension: 'yes',
                      convertToCustomerExpense: false,
                    });
                  } else if (e.target.value === 'no') {
                    setValues({
                      ...values,
                      acceptExtension: 'no',
                    });
                  }
                };
                const handleCheckBoxClick = (e) => {
                  if (e.target.value === 'false') {
                    setShowConfirmCustomerExpenseModal(true);
                  } else {
                    setValues({
                      ...values,
                      convertToCustomerExpense: false,
                    });
                  }
                };
                return (
                  <Form>
                    <DataTableWrapper
                      className={classnames('maxw-tablet', styles.sitDisplayForm)}
                      testID="sitExtensions"
                    >
                      <SitStatusTables sitStatus={sitStatus} sitExtension={sitExtension} shipment={shipment} />
                    </DataTableWrapper>
                    <div className={styles.ModalPanel}>
                      <dl className={styles.SITSummary}>
                        <div>
                          <dt>Days requested for SIT extension:</dt>
                          <dd>{sitExtension.requestedDays}</dd>
                        </div>
                        <div>
                          <dt>Reason:</dt>
                          <dd>{sitExtensionReasons[sitExtension.requestReason]}</dd>
                        </div>
                        <div>
                          <dt>Contractor remarks:</dt>
                          <dd>{sitExtension.contractorRemarks}</dd>
                        </div>
                      </dl>
                      <FormGroup>
                        <Fieldset legend="Accept request for extension?">
                          <Field
                            as={Radio}
                            label="Yes"
                            id="acceptExtension"
                            name="acceptExtension"
                            value="yes"
                            title="Yes, accept extension"
                            type="radio"
                            onChange={handleRadioSelection}
                          />
                          <Field
                            as={Radio}
                            label="No"
                            id="denyExtension"
                            name="acceptExtension"
                            value="no"
                            title="No, deny extension"
                            type="radio"
                            onChange={handleRadioSelection}
                          />
                        </Fieldset>
                      </FormGroup>
                      {values.acceptExtension === 'yes' && (
                        <div className={styles.reasonDropdown}>
                          <DropdownInput
                            label="Reason for edit"
                            name="requestReason"
                            data-testid="reasonDropdown"
                            options={dropdownInputOptions(sitExtensionReasons)}
                          />
                        </div>
                      )}
                      {values.acceptExtension === 'no' && (
                        <div className={styles.convertRadio} data-testid="convertToCustomerExpense">
                          <CheckboxField
                            id="convertToCustomerExpense"
                            label="Convert to Customer Expense"
                            name="convertToCustomerExpense"
                            onChange={handleCheckBoxClick}
                          />
                        </div>
                      )}
                      <Label htmlFor="officeRemarks">Office remarks</Label>
                      <Field
                        as={Textarea}
                        data-testid="officeRemarks"
                        label="No"
                        name="officeRemarks"
                        id="officeRemarks"
                      />
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
                      {showConfirmCustomerExpenseModal && (
                        <>
                          <Overlay />
                          <ModalContainer>
                            <ConfirmCustomerExpenseModal
                              setShowConfirmModal={setShowConfirmCustomerExpenseModal}
                              values={values}
                              setValues={setValues}
                            />
                          </ModalContainer>
                        </>
                      )}
                    </div>
                  </Form>
                );
              }}
            </Formik>
          </div>
        </Modal>
      </ModalContainer>
    </div>
  );
};

ReviewSITExtensionsModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
  sitExtension: SITExtensionShape.isRequired,
  sitStatus: SitStatusShape.isRequired,
  shipment: ShipmentShape.isRequired,
};
export default ReviewSITExtensionsModal;
