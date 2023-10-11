import React from 'react';
import classnames from 'classnames';
import { Formik, Field, useField } from 'formik';
import PropTypes from 'prop-types';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import moment from 'moment';
import * as Yup from 'yup';

import styles from './EditSitEntryDateModal.module.scss';

import DataTableWrapper from 'components/DataTableWrapper/index';
import DataTable from 'components/DataTable/index';
import { DatePickerInput } from 'components/form/fields';
import { Form } from 'components/form';
import { ModalContainer, Overlay } from 'components/MigratedModal/MigratedModal';
import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal';
import { formatDateForDatePicker, swaggerDateFormat } from 'shared/dates';

const SitEntryDateForm = ({ onChange }) => (
  <DatePickerInput name="sitEntryDate" label="" id="sitEntryDate" onChange={onChange} />
);
const DisabledSitEntryDateForm = () => (
  <DatePickerInput disabled name="prevSitEntryDate" label="" id="prevSitEntryDate" />
);

/**
 * @description This component contains the calendar pop outs and also sets the value when
 * the user changes the date in the proposed new SIT entry date change.
 */
const SitDatePickers = () => {
  // setting up a helper in order to update the values of form when the date is changed
  const entryDateHelper = useField({ name: 'sitEntryDate', id: 'sitEntryDate' })[2];
  const handleSitStartDateChange = (startDate) => {
    // Update form values
    entryDateHelper.setValue(startDate);
  };

  return (
    <div className={styles.tableContainer}>
      <DataTable
        columnHeaders={[`Original SIT entry date`, `New SIT entry date`]}
        dataRow={[
          <DisabledSitEntryDateForm />,
          <SitEntryDateForm
            onChange={(value) => {
              handleSitStartDateChange(value);
            }}
          />,
        ]}
        custClass={styles.currentLocation}
      />
    </div>
  );
};

/**
 * @description This component contains a form that can be viewed from the MTO page
 * when a user clicks "Edit" next to a service item that contains a SIT entry date
 */
const EditSitEntryDateModal = ({ onClose, onSubmit, serviceItem }) => {
  // setting initial values that requires some formatting for display requirements
  const initialValues = {
    sitEntryDate: formatDateForDatePicker(moment(serviceItem.sitEntryDate, swaggerDateFormat)),
    prevSitEntryDate: formatDateForDatePicker(moment(serviceItem.sitEntryDate, swaggerDateFormat)),
  };

  // right now the office remarks are just for show
  // TODO add change of SIT entry date to audit logs
  const editSitEntryDateModalSchema = Yup.object().shape({
    officeRemarks: Yup.string().required('Required'),
  });

  return (
    <div>
      <Overlay />
      <ModalContainer>
        <Modal className={styles.ReviewSITExtensionModal}>
          <ModalClose handleClick={() => onClose()} />
          <ModalTitle>
            <h2>Edit SIT Entry Date</h2>
          </ModalTitle>
          <div>
            <Formik
              onSubmit={(values) => onSubmit(serviceItem.id, values.sitEntryDate)}
              initialValues={{ ...initialValues, officeRemarks: '' }}
              validationSchema={editSitEntryDateModalSchema}
              validateOnMount
            >
              {({ isValid }) => {
                return (
                  <Form>
                    <DataTableWrapper
                      className={classnames('maxw-tablet', styles.sitDisplayForm)}
                      testID="sitExtensions"
                    >
                      <SitDatePickers serviceItem={serviceItem} />
                    </DataTableWrapper>
                    <Label htmlFor="officeRemarks">Office remarks</Label>
                    <Field
                      as={Textarea}
                      data-testid="officeRemarks"
                      label="No"
                      name="officeRemarks"
                      id="officeRemarks"
                    />
                    <ModalActions>
                      <Button type="submit" disabled={!isValid}>
                        Save
                      </Button>
                      <Button
                        type="button"
                        onClick={() => onClose()}
                        data-testid="modalCancelButton"
                        outline
                        className={styles.CancelButton}
                      >
                        Cancel
                      </Button>
                    </ModalActions>
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

EditSitEntryDateModal.propTypes = {
  onClose: PropTypes.func.isRequired,
  onSubmit: PropTypes.func.isRequired,
};
export default connectModal(EditSitEntryDateModal);
