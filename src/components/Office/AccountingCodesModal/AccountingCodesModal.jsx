import React from 'react';
import PropTypes from 'prop-types';
import { Formik } from 'formik';
import { Button } from '@trussworks/react-uswds';

import styles from './AccountingCodesModal.module.scss';

import AccountingCodeSection from 'components/Office/AccountingCodeSection/AccountingCodeSection';
import Modal, { ModalActions, ModalClose, ModalTitle, connectModal } from 'components/Modal/Modal';
import ShipmentTag from 'components/ShipmentTag/ShipmentTag';
import { Form } from 'components/form';
import { shipmentTypes } from 'constants/shipments';
import { AccountingCodesShape } from 'types/accountingCodes';

const AccountingCodesModal = ({ onClose, onSubmit, onEditCodesClick, shipmentType, TACs, SACs, tacType, sacType }) => {
  const handleFormSubmit = (values) => onSubmit(values);

  return (
    <Modal data-testid="AccountingCodes">
      <ModalClose handleClick={onClose} />

      <ModalTitle>
        <ShipmentTag shipmentType={shipmentType} />
        <h2 className={styles.Title}>Edit accounting codes</h2>
      </ModalTitle>

      <Formik initialValues={{ tacType, sacType }} onSubmit={handleFormSubmit}>
        <Form>
          <AccountingCodeSection
            label="TAC"
            emptyMessage="No TAC code entered."
            fieldName="tacType"
            shipmentTypes={TACs}
          />

          <AccountingCodeSection
            label="SAC (optional)"
            emptyMessage="No SAC code entered."
            fieldName="sacType"
            shipmentTypes={SACs}
          />

          <div>
            <button type="button" onClick={onEditCodesClick} className={styles.EditCodes}>
              Add or edit codes
            </button>
          </div>
          <ModalActions>
            <Button type="submit">Save</Button>
            <Button type="button" secondary onClick={onClose}>
              Cancel
            </Button>
          </ModalActions>
        </Form>
      </Formik>
    </Modal>
  );
};

AccountingCodesModal.propTypes = {
  onClose: PropTypes.func,
  onSubmit: PropTypes.func,
  onEditCodesClick: PropTypes.func,
  shipmentType: PropTypes.oneOf(Object.keys(shipmentTypes)).isRequired,
  TACs: AccountingCodesShape,
  SACs: AccountingCodesShape,
  tacType: PropTypes.string,
  sacType: PropTypes.string,
};

AccountingCodesModal.defaultProps = {
  onClose: () => {},
  onSubmit: () => {},
  onEditCodesClick: () => {},
  TACs: {},
  SACs: {},
  tacType: '',
  sacType: '',
};

export default connectModal(AccountingCodesModal);
