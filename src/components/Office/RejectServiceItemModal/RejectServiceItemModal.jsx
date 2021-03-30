import React from 'react';
import { Formik } from 'formik';
import { Modal, Button, ModalContainer, Overlay } from '@trussworks/react-uswds';
import classNames from 'classnames';
import PropTypes from 'prop-types';
import * as Yup from 'yup';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import styles from './RejectServiceItemModal.module.scss';

import { Form } from 'components/form';
import TextField from 'components/form/fields/TextField';
import ServiceItemDetails from 'components/Office/ServiceItemDetails/ServiceItemDetails';
import { formatDateFromIso } from 'shared/formatters';
import { SERVICE_ITEM_STATUS } from 'shared/constants';

const rejectionSchema = Yup.object().shape({
  rejectionReason: Yup.string().required('Required'),
});

const RejectServiceItemModal = ({ serviceItem, onSubmit, onClose }) => {
  const { serviceItem: serviceItemName, id, mtoShipmentID, code, status, createdAt, approvedAt, details } = serviceItem;
  return (
    <>
      <Overlay />
      <ModalContainer>
        <Modal className={classNames(styles.RejectServiceItemModal, 'modal', 'container', 'container--popout')}>
          <div>
            <div className={styles.modalTopContainer}>
              <h4>Are you sure you want to reject this request?</h4>
              <Button
                type="button"
                title="Close reject service item modal"
                onClick={() => onClose()}
                className={classNames(styles.rejectReasonClose, 'usa-button--unstyled')}
                data-testid="closeRejectServiceItem"
              >
                <FontAwesomeIcon icon="times" title="Close" aria-label="Close" />
              </Button>
            </div>
            <Formik
              initialValues={{ rejectionReason: '' }}
              validationSchema={rejectionSchema}
              onSubmit={(values) => {
                onSubmit(id, mtoShipmentID, SERVICE_ITEM_STATUS.REJECTED, values.rejectionReason);
              }}
            >
              {({ handleChange, values, isValid, dirty }) => {
                return (
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
                              <p>
                                {formatDateFromIso(
                                  status === SERVICE_ITEM_STATUS.SUBMITTED ? createdAt : approvedAt,
                                  'DD MMM YYYY',
                                )}
                              </p>
                            </td>
                            <td className={styles.detail}>
                              <ServiceItemDetails id={id} code={code} details={details} />
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                    <TextField
                      id="rejectionReason"
                      name="rejectionReason"
                      label="Reason for rejection"
                      type="text"
                      value={values.rejectionReason}
                      onChange={handleChange}
                    />
                    <div className={styles.modalActions}>
                      <Button type="submit" disabled={!isValid || !dirty} data-testid="submitButton">
                        Submit
                      </Button>
                      <Button secondary type="reset" onClick={() => onClose()} data-testid="backButton">
                        Back
                      </Button>
                    </div>
                  </Form>
                );
              }}
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
    mtoShipmentID: PropTypes.string,
    code: PropTypes.string,
    status: PropTypes.string,
    serviceItem: PropTypes.string,
    createdAt: PropTypes.string,
    rejectedAt: PropTypes.string,
    approvedAt: PropTypes.string,
    details: PropTypes.shape({
      description: PropTypes.string,
      pickupPostalCode: PropTypes.string,
      reason: PropTypes.string,
      itemDimensions: PropTypes.shape({ length: PropTypes.number, width: PropTypes.number, height: PropTypes.number }),
      crateDimensions: PropTypes.shape({ length: PropTypes.number, width: PropTypes.number, height: PropTypes.number }),
      firstCustomerContact: PropTypes.shape({
        timeMilitary: PropTypes.string,
        firstAvailableDeliveryDate: PropTypes.string,
      }),
    }),
  }).isRequired,
  onSubmit: PropTypes.func.isRequired,
  onClose: PropTypes.func.isRequired,
};

RejectServiceItemModal.defaultProps = {};

export default RejectServiceItemModal;
