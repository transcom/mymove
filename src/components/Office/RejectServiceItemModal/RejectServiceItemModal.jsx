import React from 'react';
import { Formik } from 'formik';
import { Modal, Button, ModalContainer, Overlay } from '@trussworks/react-uswds';
import classNames from 'classnames';
import PropTypes from 'prop-types';

import styles from './RejectServiceItemModal.module.scss';

import { Form } from 'components/form';
import { TextInput } from 'components/form/fields';
import { ReactComponent as XLightIcon } from 'shared/icon/x-light.svg';
import { formatDate } from 'shared/dates';

const RejectServiceItemModal = ({ serviceItem, onSubmit, onClose }) => {
  // eslint-disable-next-line no-unused-vars
  const { serviceItem: serviceItemName, code, submittedAt, details } = serviceItem;
  return (
    <>
      <Overlay />
      <ModalContainer>
        <Modal className="modal container container--popout">
          <div>
            <div className={styles.modalTopContainer}>
              <h4>Are you sure you want to reject this request?</h4>
              <button
                type="button"
                title="Close reject service item modal"
                onClick={() => onClose()}
                className={classNames(styles.rejectReasonClose, 'usa-button--unstyled')}
                data-testid="closeRejectServiceItem"
              >
                <XLightIcon />
              </button>
            </div>
            <Formik initialValues={{ rejectionReason: '' }} onSubmit={onSubmit}>
              <Form>
                <div className={('table--service-item', 'table--service-item--hasimg')}>
                  <table>
                    <thead className="table--small">
                      <tr>
                        <th>Service item</th>
                        <th>Details</th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td className={styles.nameAndDate}>
                          <p className={styles.codeName}>{serviceItemName}</p>
                          <p>{formatDate(submittedAt, 'DD MMM YYYY')}</p>
                        </td>
                        <td />
                      </tr>
                    </tbody>
                  </table>
                </div>
                <TextInput name="rejectionReason" label="Reason for rejection" type="text" />
                <div className={styles.modalActions}>
                  <Button type="submit">Submit</Button>
                  <Button secondary type="reset" onClick={() => onClose()}>
                    Back
                  </Button>
                </div>
              </Form>
            </Formik>
          </div>
        </Modal>
      </ModalContainer>
    </>
  );
};

RejectServiceItemModal.propTypes = {
  serviceItem: PropTypes.shape({
    id: PropTypes.string,
    code: PropTypes.string,
    serviceItem: PropTypes.string,
    submittedAt: PropTypes.string,
    details: {
      description: PropTypes.string,
      pickupPostalCode: PropTypes.string,
      reason: PropTypes.string,
      itemDimensions: PropTypes.shape({ length: PropTypes.number, width: PropTypes.number, height: PropTypes.number }),
      crateDimensions: PropTypes.shape({ length: PropTypes.number, width: PropTypes.number, height: PropTypes.number }),
      firstCustomerContact: PropTypes.shape({
        timeMilitary: PropTypes.string,
        firstAvailableDeliveryDate: PropTypes.string,
      }),
    },
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};

RejectServiceItemModal.defaultProps = {};

export default RejectServiceItemModal;
