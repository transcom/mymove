import React from 'react';
import PropTypes from 'prop-types';
import { Button, Label, Textarea } from '@trussworks/react-uswds';
import { Formik, Field } from 'formik';

import { ServiceItemDetailsShape } from '../../../types/serviceItems';
import ServiceItemDetails from '../ServiceItemDetails/ServiceItemDetails';

import styles from './ServiceItemUpdateModal.module.scss';

import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal';
import { Form } from 'components/form';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { formatDateFromIso } from 'utils/formatters';

const ServiceItemDetail = ({ serviceItem }) => {
  const { id, code, submittedAt, details } = serviceItem;
  return (
    <table className={styles.serviceItemDetails}>
      <tr key={id}>
        <td className={styles.nameDateContainer}>
          <p className={styles.serviceItemName}>{serviceItem.serviceItem}</p>
          <p>{formatDateFromIso(submittedAt, 'DD MMM YYYY')}</p>
        </td>
        <td className={styles.detailsContainer}>
          <ServiceItemDetails code={code} details={details} />
        </td>
      </tr>
    </table>
  );
};
/**
 * @description This componment is the modal used for when a TOO edits the address for a Service item
 * or reviews a service item request from a the prime.
 */
export const ServiceItemUpdateModal = ({ onSave, closeModal, title, serviceItem, content }) => {
  const initialValues = {
    officeRemarks: '',
    newAddress: {
      streetAddress1: '',
      streetAddress2: '',
      city: '',
      state: '',
      postalCode: '',
    },
  };
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
      <Formik onSubmit={(e) => onSave(e)} initialValues={initialValues}>
        {({ isValid }) => {
          return (
            <Form>
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
  content: PropTypes.element,
};

ServiceItemUpdateModal.defaultProps = {
  content: <div />,
};

ServiceItemUpdateModal.displayName = 'ServiceItemUpdateModal';
export default connectModal(ServiceItemUpdateModal);
