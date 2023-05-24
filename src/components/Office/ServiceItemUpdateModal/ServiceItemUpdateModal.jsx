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

/**
 * @function
 * @description Return the details component that has the service item name, date, and additonal information specific to teh service item type.
 * @param {Object} Props
 * @param {ServiceItemDetailsShape} Props.serviceItem
 * @returns {React.ReactElement}
 */
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
            <ServiceItemDetails id={`sid-${id}`} code={code} details={details} />
          </td>
        </tr>
      </tbody>
    </table>
  );
};

/**
 * @component
 * @description This componment is the modal used for when a TOO edits the address for a Service item
 * or reviews a service item request from a the prime.
 * @param {ServiceItemUpdateModalProps}
 *
 * @returns {React.ReactElement}
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
  /**
   * @description The validation schema takes in the validations specific to the modal.
   * Since office remarks is shared is already declared in the schema.
   */
  const serviceItemUpdateModalSchema = Yup.object().shape({
    officeRemarks: Yup.string().required('Required'),
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
        validationSchema={serviceItemUpdateModalSchema}
        onSubmit={(formValues) => onSave(serviceItem.id, formValues)}
        // add Office remarks to the initial values as it's shared by all modals
        initialValues={{ ...initialValues, officeRemarks: '' }}
        validateOnMount
      >
        {({ isValid }) => {
          return (
            <Form className={formStyles.form}>
              <div className={styles.sitPanelForm}>
                <h3 className={styles.modalReviewHeader}>SIT delivery address</h3>
                {content}
                <Label htmlFor="officeRemarks">Office remarks</Label>
                <Field
                  as={Textarea}
                  data-testid="officeRemarks"
                  label="No"
                  name="officeRemarks"
                  id="officeRemarks"
                  className={styles.officeRemarks}
                />
              </div>
              <ModalActions>
                <Button type="submit" disabled={!isValid}>
                  Save
                </Button>
                <Button type="button" secondary onClick={closeModal}>
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

/**
 * @typedef {object} ServiceItemUpdateModalProps
 * @prop {function} onSave saves the form values
 * @prop {function} closeModal closes the modal without saving changes
 * @prop {string} title title of the modal
 * @prop {ServiceItemDetailsShape} serviceItem the current service item selected
 * @prop {element} content additional form contents specific to either Edit or Reviewing address updates
 * @prop {object} initialValues the initialValues for the form
 * @prop {object} validations Form validaitons specific to the modal.
 * @extends {ServiceItemUpdateModal<ServiceItemUpdateModalProps>}
 */

ServiceItemUpdateModal.propTypes = {
  closeModal: PropTypes.func.isRequired,
  onSave: PropTypes.func.isRequired,
  title: PropTypes.string.isRequired,
  serviceItem: ServiceItemDetailsShape.isRequired,
  initialValues: PropTypes.object,
  validations: PropTypes.object,
};

ServiceItemUpdateModal.defaultProps = {
  initialValues: {},
  validations: {},
};

ServiceItemUpdateModal.displayName = 'ServiceItemUpdateModal';
export default connectModal(ServiceItemUpdateModal);
