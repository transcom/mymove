import React from 'react';
import PropTypes from 'prop-types';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import { Formik, Field } from 'formik';
import * as Yup from 'yup';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';
import ServiceItemDetails from '../ServiceItemDetails/ServiceItemDetails';

import styles from './ServiceItemUpdateModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal';
import { Form } from 'components/form/Form';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { formatDateFromIso } from 'utils/formatters';
import formStyles from 'styles/form.module.scss';

const ServiceItemDetail = ({ serviceItem }) => {
  const { id, code, submittedAt, details } = serviceItem;
  return (
    <table data-testid="sitAddressUpdateDetailTable" className={styles.serviceItemDetails}>
      <tbody>
        <tr key={`sid-${id}`}>
          <td className={styles.nameDateContainer}>
            <p className={styles.serviceItemName}>{serviceItem.serviceItem}</p>
            <p>{formatDateFromIso(submittedAt, 'DD MMM YYYY')}</p>
          </td>
          <td className={styles.detailsContainer}>
            <ServiceItemDetails code={code} details={details} />
          </td>
        </tr>
      </tbody>
    </table>
  );
};
/**
 * @description This componment is the modal used for when a TOO edits the address for a Service item
 * or reviews a service item request from a the prime.
 * @param {function} onSave saves the form values
 * @param {function} closeModal closes the modal without saving changes
 * @param {string} title title of the modal
 * @param {ServiceItemDetailsShape} serviceItem
 * @param {element} content the form specific to either Edit or Reviewing address updates
 * @param {object} initialValues the initialValues for the form
 * @param {object} validations Form validaitons specific to the modal.
 */
export const ServiceItemUpdateModal = ({
  onSave,
  closeModal,
  title,
  serviceItem,
  children: content,
  initialValues,
  validations,
}) => {
  const serviceItemUpdateModalSchema = Yup.object().shape({
    officeRemarks: Yup.string().nullable(),
    ...validations,
  });
  return (
    <Modal className={styles.serviceItemUpdateModal}>
      <div>
        <ShipmentTag shipmentType="HHG" />
        <ModalClose handleClick={() => closeModal()} />
      </div>
      <ModalTitle className={styles.titleSection}>
        <h2>{title}</h2>
      </ModalTitle>
      <ServiceItemDetail serviceItem={serviceItem} />
      <Formik
        onSubmit={(formValues) => onSave(serviceItem.id, formValues)}
        initialValues={initialValues}
        validationSchema={serviceItemUpdateModalSchema}
      >
        {({ isValid }) => {
          return (
            <Form className={formStyles.form}>
              <h3>SIT delivery address</h3>
              {content}
              <Label htmlFor="officeRemarks">Office remarks</Label>
              <Field as={Textarea} data-testid="officeRemarks" label="No" name="officeRemarks" id="officeRemarks" />
              <ModalActions>
                <Button type="submit">Save</Button>
                <Button type="button" secondary onClick={closeModal} disabled={!isValid}>
                  Cancel
                </Button>
              </ModalActions>
            </Form>
          );
        }}
      </Formik>
    </Modal>
  );
};

ServiceItemUpdateModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
  title: PropTypes.string.isRequired,
  serviceItem: ServiceItemDetailsShape.isRequired,
};

ServiceItemUpdateModal.displayName = 'ServiceItemUpdateModal';
export default connectModal(ServiceItemUpdateModal);
